package sqlmaper

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type (
	ColumnData struct {
		ColumnName string
		FieldIndex []int
		GoType     reflect.Type
		Optional   bool
	}
	ColumnMap map[string]ColumnData
)

const (
	followTagName = "follow"
	embedTagName  = "embed"
	notateTagName = "notate"
)

func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Invalid:
		return true
	}
	return false
}

var structMapCache = make(map[interface{}]ColumnMap)
var structMapCacheLock = sync.Mutex{}

// notateByDefault will dot annotate an embedded struct
//
// Example:
//
//	type EmbeddedStruct struct {
//	   String string
//	}
//
//	type Struct struct {
//	    TableOne EmbeddedStruct `db:"table_one"`
//	}
//
// Output: "table_one"."string"
//
// This is true by default
// This can be changed by calling AnnotatedByDefault(bool)
var notateByDefault = false

func NotatedByDefault(notate bool) {
	notateByDefault = notate
}
func isNotated() bool {
	return notateByDefault
}

var camelCaseColumnRenameFunction = camelCase
var lowerCaseColumnRenameFunction = strings.ToLower

var defaultColumnRenameFunction = camelCaseColumnRenameFunction
var columnRenameFunction = defaultColumnRenameFunction

func SetColumnRenameFunction(newFunction func(string) string) {
	columnRenameFunction = newFunction
}

// GetSliceElementType returns the type for a slices elements.
func GetSliceElementType(val reflect.Value) reflect.Type {
	elemType := val.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	return elemType
}

// AppendSliceElement will append val to slice. Handles slice of pointers and
// not pointers. Val needs to be a pointer.
func AppendSliceElement(slice, val reflect.Value) {
	if slice.Type().Elem().Kind() == reflect.Ptr {
		slice.Set(reflect.Append(slice, val))
	} else {
		slice.Set(reflect.Append(slice, reflect.Indirect(val)))
	}
}

func GetTypeInfo(i interface{}, val reflect.Value) (reflect.Type, reflect.Kind) {
	var t reflect.Type
	valKind := val.Kind()
	if valKind == reflect.Slice {
		if reflect.ValueOf(i).Kind() == reflect.Ptr {
			t = reflect.TypeOf(i).Elem().Elem()
		} else {
			t = reflect.TypeOf(i).Elem()
		}
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		valKind = t.Kind()
	} else {
		t = val.Type()
	}
	return t, valKind
}

func SafeGetFieldByIndex(v reflect.Value, fieldIndex []int) (result reflect.Value, isAvailable bool) {
	switch len(fieldIndex) {
	case 0:
		return v, true
	case 1:
		return v.FieldByIndex(fieldIndex), true
	default:
		if f := reflect.Indirect(v.Field(fieldIndex[0])); f.IsValid() {
			return SafeGetFieldByIndex(f, fieldIndex[1:])
		}
	}
	return reflect.ValueOf(nil), false
}

func SafeSetFieldByIndex(v reflect.Value, fieldIndex []int, src interface{}) (result reflect.Value) {
	v = reflect.Indirect(v)
	switch len(fieldIndex) {
	case 0:
		return v
	case 1:
		f := v.FieldByIndex(fieldIndex)
		srcVal := reflect.ValueOf(src)
		f.Set(reflect.Indirect(srcVal))
	default:
		f := v.Field(fieldIndex[0])
		switch f.Kind() {
		case reflect.Ptr:
			s := f
			if f.IsNil() || !f.IsValid() {
				s = reflect.New(f.Type().Elem())
				f.Set(s)
			}
			SafeSetFieldByIndex(reflect.Indirect(s), fieldIndex[1:], src)
		case reflect.Struct:
			SafeSetFieldByIndex(f, fieldIndex[1:], src)
		}
	}
	return v
}

type rowData = map[string]interface{}

// AssignStructVals will assign the data from rd to i.
func AssignStructVals(i interface{}, rd rowData, cm ColumnMap) {
	val := reflect.Indirect(reflect.ValueOf(i))

	for name, data := range cm {
		src, ok := rd[name]
		if ok {
			SafeSetFieldByIndex(val, data.FieldIndex, src)
		}
	}
}

func GetColumnMap(i interface{}) (ColumnMap, error) {
	val := reflect.Indirect(reflect.ValueOf(i))
	t, valKind := GetTypeInfo(i, val)
	if valKind != reflect.Struct {
		return nil, fmt.Errorf("cannot scan into this type: %v", t) // #nosec
	}

	structMapCacheLock.Lock()
	defer structMapCacheLock.Unlock()
	if _, ok := structMapCache[t]; !ok {
		structMapCache[t] = createColumnMap(t, []int{}, []string{}, false)
	}
	return structMapCache[t], nil
}

func createColumnMap(t reflect.Type, fieldIndex []int, prefixes []string, optional bool) ColumnMap {
	cm, n := ColumnMap{}, t.NumField()
	var subColMaps []ColumnMap
	for i := 0; i < n; i++ {
		f := t.Field(i)
		dbTag := NewTag("db", f.Tag)
		if !dbTag.Ignore() {
			// get scan tags if present
			scanTag := NewTag("scan", f.Tag).Values()
			// merge dbTag options with scan tags for more manageable tag handling
			options := append(dbTag.Options(), scanTag...)
			var columnName string

			if !dbTag.IsNamed() {
				columnName = columnRenameFunction(f.Name)
			} else {
				columnName = dbTag.Name()
			}

			if (f.Anonymous || options.Contains(followTagName)) && IsUnderlyingStruct(f.Type) {
				subFieldIndexes := append(fieldIndex, f.Index...)

				if f.Type.Kind() == reflect.Ptr {
					f.Type = f.Type.Elem()
				}

				if dbTag.IsNamed() && !options.Contains(followTagName) {
					subPrefixes := append(prefixes, columnName)
					subColMaps = append(subColMaps, createColumnMap(f.Type, subFieldIndexes, subPrefixes, false))
				} else {
					subColMaps = append(subColMaps, createColumnMap(f.Type, subFieldIndexes, prefixes, false))
				}

			} else if !implementsScanner(f.Type) && (isNotated() || options.Contains(notateTagName)) && !options.Contains(embedTagName) {
				subFieldIndexes := append(fieldIndex, f.Index...)
				subPrefixes := append(prefixes, columnName)
				var subCm ColumnMap
				if f.Type.Kind() == reflect.Ptr {
					subCm = createColumnMap(f.Type.Elem(), subFieldIndexes, subPrefixes, true)
				} else {
					subCm = createColumnMap(f.Type, subFieldIndexes, subPrefixes, false)
				}
				if len(subCm) != 0 {
					subColMaps = append(subColMaps, subCm)
					continue
				}
			} else if f.PkgPath == "" {
				// if PkgPath is empty then it is an exported field
				columnName = strings.Join(append(prefixes, columnName), ".")
				var goType reflect.Type
				if optional {
					goType = reflect.PtrTo(f.Type)
				} else {
					goType = f.Type
				}
				cm[columnName] = ColumnData{
					ColumnName: columnName,
					FieldIndex: append(fieldIndex, f.Index...),
					GoType:     goType,
					Optional:   optional,
				}
			}
		}
	}
	for _, subCm := range subColMaps {
		for key, val := range subCm {
			if _, ok := cm[key]; !ok {
				cm[key] = val
			}
		}
	}
	return cm
}

func IsUnderlyingStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}
func (cm ColumnMap) Cols() []string {
	var structCols []string
	for key := range cm {
		structCols = append(structCols, key)
	}
	sort.Strings(structCols)
	return structCols
}

func implementsScanner(t reflect.Type) bool {
	if IsPointer(t.Kind()) {
		t = t.Elem()
	}
	if reflect.PtrTo(t).Implements(scannerType) {
		return true
	}
	if IsBuiltin(t) { // accounts for time.Time builtin struct
		return true
	}

	return false
}
