package sqlmaper

import (
	"database/sql"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type reflectTest struct {
	suite.Suite
}

func (rt *reflectTest) TestColumnRename() {
	// different key names are used each time to circumvent the caching that happens
	// it seems like a solid assumption that when people use this feature,
	// they would simply set a renaming function once at startup,
	// and not change between requests like this

	// changing rename function
	SetColumnRenameFunction(camelCase)
	camelAnon := struct {
		FirstCamel string
		LastCamel  string
	}{}
	camelColumnMap, camelErr := GetColumnMap(&camelAnon)
	rt.NoError(camelErr)

	var camelKeys []string
	for key := range camelColumnMap {
		camelKeys = append(camelKeys, key)
	}
	rt.Contains(camelKeys, "first_camel")
	rt.Contains(camelKeys, "last_camel")

	// changing rename function
	SetColumnRenameFunction(lowerCaseColumnRenameFunction)

	lowerAnon := struct {
		FirstLower string
		LastLower  string
	}{}
	lowerColumnMap, lowerErr := GetColumnMap(&lowerAnon)
	rt.NoError(lowerErr)

	var lowerKeys []string
	for key := range lowerColumnMap {
		lowerKeys = append(lowerKeys, key)
	}
	rt.Contains(lowerKeys, "firstlower")
	rt.Contains(lowerKeys, "lastlower")

	// changing rename function
	SetColumnRenameFunction(strings.ToUpper)

	upperAnon := struct {
		FirstUpper string
		LastUpper  string
	}{}
	upperColumnMap, upperErr := GetColumnMap(&upperAnon)
	rt.NoError(upperErr)

	var upperKeys []string
	for key := range upperColumnMap {
		upperKeys = append(upperKeys, key)
	}
	rt.Contains(upperKeys, "FIRSTUPPER")
	rt.Contains(upperKeys, "LASTUPPER")

	SetColumnRenameFunction(defaultColumnRenameFunction)
}

func (rt *reflectTest) TestParallelGetColumnMap() {

	type item struct {
		id   uint
		name string
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		i := item{id: 1, name: "bob"}
		m, err := GetColumnMap(i)
		rt.NoError(err)
		rt.NotNil(m)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		i := item{id: 2, name: "sally"}
		m, err := GetColumnMap(i)
		rt.NoError(err)
		rt.NotNil(m)
		wg.Done()
	}()

	wg.Wait()
}

func (rt *reflectTest) TestAssignStructVals_withStruct() {

	type TestStruct struct {
		Str    string
		Int    int64
		Bool   bool
		Valuer sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	data := map[string]interface{}{
		"str":    "string",
		"int":    int64(10),
		"bool":   true,
		"valuer": sql.NullString{String: "null_str", Valid: true},
	}

	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		Str:    "string",
		Int:    10,
		Bool:   true,
		Valuer: sql.NullString{String: "null_str", Valid: true},
	})
}

func (rt *reflectTest) TestAssignStructVals_withStructWithPointerVals() {
	type TestStruct struct {
		Str    string
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"str":    "string",
		"int":    int64(10),
		"bool":   true,
		"valuer": &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		Str:    "string",
		Int:    10,
		Bool:   true,
		Valuer: ns,
	})
}

func (rt *reflectTest) TestAssignStructVals_withStructWithEmbeddedStruct() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		EmbeddedStruct
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"str":    "string",
		"int":    int64(10),
		"bool":   true,
		"valuer": &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}

func (rt *reflectTest) TestAssignStructVals_withStructWithEmbeddedStructPointer() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		*EmbeddedStruct
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"str":    "string",
		"int":    int64(10),
		"bool":   true,
		"valuer": &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: &EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithEmbeddedStructWithFollowTag() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		EmbeddedStruct
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"str":    "string",
		"int":    int64(10),
		"bool":   true,
		"valuer": &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}

func (rt *reflectTest) TestAssignStructVals_withStructWithEmbeddedStructPointerWithFollowTag() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		*EmbeddedStruct
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"str":    "string",
		"int":    int64(10),
		"bool":   true,
		"valuer": &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: &EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithEmbeddedStructWithMultiLevelFollow() {

	type Follow struct {
		End string
	}
	type EmbeddedStruct struct {
		Follow Follow `scan:"follow"`
		Str    string
	}
	type TestStruct struct {
		EmbeddedStruct `db:",follow"`
		Int            int64
		Bool           bool
		Valuer         *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"str":    "string",
		"end":    "fin",
		"int":    int64(10),
		"bool":   true,
		"valuer": &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: EmbeddedStruct{Str: "string", Follow: Follow{End: "fin"}},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithTaggedEmbeddedStruct() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		EmbeddedStruct `db:"embedded"`
		Int            int64
		Bool           bool
		Valuer         *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"embedded.str": "string",
		"int":          int64(10),
		"bool":         true,
		"valuer":       &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}

func (rt *reflectTest) TestAssignStructVals_withStructWithTaggedEmbeddedPointer() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		*EmbeddedStruct `db:"embedded"`
		Int             int64
		Bool            bool
		Valuer          *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"embedded.str": "string",
		"int":          int64(10),
		"bool":         true,
		"valuer":       &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: &EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withNamedStructFieldWithFollowTaggedStructField() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		Embedded EmbeddedStruct `db:"name_isnt_evaluated_due_to_follow_tag,follow"`
		Int      int64
		Bool     bool
		Valuer   *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"str":    "string",
		"int":    int64(10),
		"bool":   true,
		"valuer": &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		Embedded: EmbeddedStruct{Str: "string"},
		Int:      10,
		Bool:     true,
		Valuer:   ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withNamedPointerStructFieldWithFollowTaggedStructField() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		Embedded *EmbeddedStruct `db:"name_isnt_evaluated_due_to_follow_tag" scan:"follow"`
		Int      int64
		Bool     bool
		Valuer   *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"str":    "string",
		"int":    int64(10),
		"bool":   true,
		"valuer": &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		Embedded: &EmbeddedStruct{Str: "string"},
		Int:      10,
		Bool:     true,
		Valuer:   ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithTaggedStructField() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		Embedded EmbeddedStruct `db:"embedded" scan:"notate"`
		Int      int64
		Bool     bool
		Valuer   *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"embedded.str": "string",
		"int":          int64(10),
		"bool":         true,
		"valuer":       &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		Embedded: EmbeddedStruct{Str: "string"},
		Int:      10,
		Bool:     true,
		Valuer:   ns,
	})
}

func (rt *reflectTest) TestAssignStructVals_withStructWithTaggedPointerField() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		Embedded *EmbeddedStruct `db:"embedded" scan:"notate"`
		Int      int64
		Bool     bool
		Valuer   *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"embedded.str": "string",
		"int":          int64(10),
		"bool":         true,
		"valuer":       &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		Embedded: &EmbeddedStruct{Str: "string"},
		Int:      10,
		Bool:     true,
		Valuer:   ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithAnonEmbeddedStructAnnotationOff() {

	NotatedByDefault(false)

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		EmbeddedStruct
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"str":    "string",
		"int":    int64(10),
		"bool":   true,
		"valuer": &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithNamedAnonEmbeddedStructAnnotationOff() {

	NotatedByDefault(false)

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		EmbeddedStruct `db:"embedded"`
		Int            int64
		Bool           bool
		Valuer         *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"embedded.str": "string",
		"int":          int64(10),
		"bool":         true,
		"valuer":       &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_StructWithNamedEmbeddedStructWithAnnotateTagAndAnnotationOff() {

	NotatedByDefault(false)

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		EmbeddedStruct EmbeddedStruct `db:"embedded" scan:"notate"`
		Int            int64
		Bool           bool
		Valuer         *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"embedded.str": "string",
		"int":          int64(10),
		"bool":         true,
		"valuer":       &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}

// this test is for when you turn annotation off, maybe randomly, but checks to make sure that the
// behavior of embedded tag is the same. Should do the exact same thing.
func (rt *reflectTest) TestAssignStructVals_StructWithNamedEmbeddedStructWithEmbedTagAndAnnotationOff() {

	NotatedByDefault(false)

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		EmbeddedStruct EmbeddedStruct `db:"embedded"`
		Int            int64
		Bool           bool
		Valuer         *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"embedded": &EmbeddedStruct{Str: "string"},
		"int":      int64(10),
		"bool":     true,
		"valuer":   &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithNamedEmbeddedStructAnnotationOff() {

	NotatedByDefault(false)

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		EmbeddedStruct EmbeddedStruct `db:"embedded"`
		Int            int64
		Bool           bool
		Valuer         *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"embedded": &EmbeddedStruct{Str: "string"},
		"int":      int64(10),
		"bool":     true,
		"valuer":   &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}

func (rt *reflectTest) TestAssignStructVals_withStructWithUnNamedEmbeddedStructPointerAnnotationOff() {
	NotatedByDefault(false)

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		*EmbeddedStruct
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"str":    "string",
		"int":    int64(10),
		"bool":   true,
		"valuer": &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		EmbeddedStruct: &EmbeddedStruct{Str: "string"},
		Int:            10,
		Bool:           true,
		Valuer:         ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithTaggedStructFieldOfAsIs() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		Embedded EmbeddedStruct `db:"embedded"`
		Int      int64
		Bool     bool
		Valuer   *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"embedded": EmbeddedStruct{Str: "string"},
		"int":      int64(10),
		"bool":     true,
		"valuer":   &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		Embedded: EmbeddedStruct{Str: "string"},
		Int:      10,
		Bool:     true,
		Valuer:   ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithTaggedPointerStructFieldOfAsIs() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		Embedded *EmbeddedStruct `db:"embedded"`
		Int      int64
		Bool     bool
		Valuer   *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	es := &EmbeddedStruct{Str: "string"}
	data := map[string]interface{}{
		"embedded": &es,
		"int":      int64(10),
		"bool":     true,
		"valuer":   &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		Embedded: es,
		Int:      10,
		Bool:     true,
		Valuer:   ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithEmbeddedStructWithAsIsTag() {
	type AsIsStruct struct {
		Int int
	}
	type EmbeddedStruct struct {
		Str  string
		AsIs AsIsStruct `db:"as_is"`
	}
	type TestStruct struct {
		Embedded EmbeddedStruct `db:"embedded" scan:"notate"`
		Int      int64
		Bool     bool
		Valuer   *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	data := map[string]interface{}{
		"embedded.str":   "string",
		"embedded.as_is": AsIsStruct{Int: 5},
		"int":            int64(10),
		"bool":           true,
		"valuer":         &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		Embedded: EmbeddedStruct{Str: "string", AsIs: AsIsStruct{Int: 5}},
		Int:      10,
		Bool:     true,
		Valuer:   ns,
	})
}
func (rt *reflectTest) TestAssignStructVals_withStructWithEmbeddedPointerStructWithAsIsTag() {
	type AsIsStruct struct {
		Int int
	}
	type EmbeddedStruct struct {
		Str  string
		AsIs *AsIsStruct `db:"as_is"`
	}
	type TestStruct struct {
		Embedded *EmbeddedStruct `db:"embedded" scan:"notate"`
		Int      int64
		Bool     bool
		Valuer   *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	ns := &sql.NullString{String: "null_str1", Valid: true}
	es := &EmbeddedStruct{Str: "string", AsIs: &AsIsStruct{Int: 5}}
	data := map[string]interface{}{
		"embedded.str":   "string",
		"embedded.as_is": &es.AsIs,
		"int":            int64(10),
		"bool":           true,
		"valuer":         &ns,
	}
	AssignStructVals(&ts, data, cm)
	rt.Equal(ts, TestStruct{
		Embedded: es,
		Int:      10,
		Bool:     true,
		Valuer:   ns,
	})
}
func (rt *reflectTest) TestGetColumnMap_withStruct() {

	type TestStruct struct {
		Str    string
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"str":    {ColumnName: "str", FieldIndex: []int{0}, GoType: reflect.TypeOf("")},
		"int":    {ColumnName: "int", FieldIndex: []int{1}, GoType: reflect.TypeOf(int64(1))},
		"bool":   {ColumnName: "bool", FieldIndex: []int{2}, GoType: reflect.TypeOf(true)},
		"valuer": {ColumnName: "valuer", FieldIndex: []int{3}, GoType: reflect.TypeOf(&sql.NullString{})},
	}, cm)
}
func (rt *reflectTest) TestGetColumnMap_withStructWithTag() {

	type TestStruct struct {
		Str    string          `db:"s"`
		Int    int64           `db:"i"`
		Bool   bool            `db:"b"`
		Valuer *sql.NullString `db:"v"`
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"s": {ColumnName: "s", FieldIndex: []int{0}, GoType: reflect.TypeOf("")},
		"i": {ColumnName: "i", FieldIndex: []int{1}, GoType: reflect.TypeOf(int64(1))},
		"b": {ColumnName: "b", FieldIndex: []int{2}, GoType: reflect.TypeOf(true)},
		"v": {ColumnName: "v", FieldIndex: []int{3}, GoType: reflect.TypeOf(&sql.NullString{})},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withOmitemptyTag() {

	type TestStruct struct {
		Str    *string `db:"str"`
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	var ps *string
	rt.Equal(ColumnMap{
		"str":    {ColumnName: "str", FieldIndex: []int{0}, GoType: reflect.TypeOf(ps)},
		"int":    {ColumnName: "int", FieldIndex: []int{1}, GoType: reflect.TypeOf(int64(1))},
		"bool":   {ColumnName: "bool", FieldIndex: []int{2}, GoType: reflect.TypeOf(true)},
		"valuer": {ColumnName: "valuer", FieldIndex: []int{3}, GoType: reflect.TypeOf(&sql.NullString{})},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withStructWithTransientFields() {

	type TestStruct struct {
		Str    string
		Int    int64
		Bool   bool
		Valuer *sql.NullString `db:"-"`
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"str":  {ColumnName: "str", FieldIndex: []int{0}, GoType: reflect.TypeOf("")},
		"int":  {ColumnName: "int", FieldIndex: []int{1}, GoType: reflect.TypeOf(int64(1))},
		"bool": {ColumnName: "bool", FieldIndex: []int{2}, GoType: reflect.TypeOf(true)},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withSliceOfStructs() {

	type TestStruct struct {
		Str    string
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts []TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"str":    {ColumnName: "str", FieldIndex: []int{0}, GoType: reflect.TypeOf("")},
		"int":    {ColumnName: "int", FieldIndex: []int{1}, GoType: reflect.TypeOf(int64(1))},
		"bool":   {ColumnName: "bool", FieldIndex: []int{2}, GoType: reflect.TypeOf(true)},
		"valuer": {ColumnName: "valuer", FieldIndex: []int{3}, GoType: reflect.TypeOf(&sql.NullString{})},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withNonStruct() {

	var v int64
	_, err := GetColumnMap(&v)
	rt.EqualError(err, "cannot scan into this type: int64")

}

func (rt *reflectTest) TestGetColumnMap_withStructWithEmbeddedStruct() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		EmbeddedStruct
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"str":    {ColumnName: "str", FieldIndex: []int{0, 0}, GoType: reflect.TypeOf("")},
		"int":    {ColumnName: "int", FieldIndex: []int{1}, GoType: reflect.TypeOf(int64(1))},
		"bool":   {ColumnName: "bool", FieldIndex: []int{2}, GoType: reflect.TypeOf(true)},
		"valuer": {ColumnName: "valuer", FieldIndex: []int{3}, GoType: reflect.TypeOf(&sql.NullString{})},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withStructWithEmbeddedStructPointer() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		*EmbeddedStruct
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"str":    {ColumnName: "str", FieldIndex: []int{0, 0}, GoType: reflect.TypeOf("")},
		"int":    {ColumnName: "int", FieldIndex: []int{1}, GoType: reflect.TypeOf(int64(1))},
		"bool":   {ColumnName: "bool", FieldIndex: []int{2}, GoType: reflect.TypeOf(true)},
		"valuer": {ColumnName: "valuer", FieldIndex: []int{3}, GoType: reflect.TypeOf(&sql.NullString{})},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withIgnoredEmbeddedStruct() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		EmbeddedStruct `db:"-"`
		Int            int64
		Bool           bool
		Valuer         *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"int":    {ColumnName: "int", FieldIndex: []int{1}, GoType: reflect.TypeOf(int64(1))},
		"bool":   {ColumnName: "bool", FieldIndex: []int{2}, GoType: reflect.TypeOf(true)},
		"valuer": {ColumnName: "valuer", FieldIndex: []int{3}, GoType: reflect.TypeOf(&sql.NullString{})},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withIgnoredEmbeddedPointerStruct() {

	type EmbeddedStruct struct {
		Str string
	}
	type TestStruct struct {
		*EmbeddedStruct `db:"-"`
		Int             int64
		Bool            bool
		Valuer          *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"int":    {ColumnName: "int", FieldIndex: []int{1}, GoType: reflect.TypeOf(int64(1))},
		"bool":   {ColumnName: "bool", FieldIndex: []int{2}, GoType: reflect.TypeOf(true)},
		"valuer": {ColumnName: "valuer", FieldIndex: []int{3}, GoType: reflect.TypeOf(&sql.NullString{})},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withPrivateFields() {

	type TestStruct struct {
		str    string // nolint:structcheck,unused
		Int    int64
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"int":    {ColumnName: "int", FieldIndex: []int{1}, GoType: reflect.TypeOf(int64(1))},
		"bool":   {ColumnName: "bool", FieldIndex: []int{2}, GoType: reflect.TypeOf(true)},
		"valuer": {ColumnName: "valuer", FieldIndex: []int{3}, GoType: reflect.TypeOf(&sql.NullString{})},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withPrivateEmbeddedFields() {

	type TestEmbedded struct {
		str string // nolint:structcheck,unused
		Int int64
	}

	type TestStruct struct {
		TestEmbedded
		Bool   bool
		Valuer *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"int":    {ColumnName: "int", FieldIndex: []int{0, 1}, GoType: reflect.TypeOf(int64(1))},
		"bool":   {ColumnName: "bool", FieldIndex: []int{1}, GoType: reflect.TypeOf(true)},
		"valuer": {ColumnName: "valuer", FieldIndex: []int{2}, GoType: reflect.TypeOf(&sql.NullString{})},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withEmbeddedTaggedStruct() {

	type TestEmbedded struct {
		Bool   bool
		Valuer *sql.NullString
	}

	type TestStruct struct {
		TestEmbedded `db:"test_embedded" scan:"notate"`
		Bool         bool
		Valuer       *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"test_embedded.bool": {
			ColumnName: "test_embedded.bool",
			FieldIndex: []int{0, 0},
			GoType:     reflect.TypeOf(true),
		},
		"test_embedded.valuer": {
			ColumnName: "test_embedded.valuer",
			FieldIndex: []int{0, 1},
			GoType:     reflect.TypeOf(&sql.NullString{}),
		},
		"bool": {
			ColumnName: "bool",
			FieldIndex: []int{1},
			GoType:     reflect.TypeOf(true),
		},
		"valuer": {
			ColumnName: "valuer",
			FieldIndex: []int{2},
			GoType:     reflect.TypeOf(&sql.NullString{}),
		},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withEmbeddedTaggedStructPointer() {

	type TestEmbedded struct {
		Bool   bool
		Valuer *sql.NullString
	}

	type TestStruct struct {
		*TestEmbedded `db:"test_embedded" scan:"notate"`
		Bool          bool
		Valuer        *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"test_embedded.bool": {
			ColumnName: "test_embedded.bool",
			FieldIndex: []int{0, 0},
			GoType:     reflect.TypeOf(true),
		},
		"test_embedded.valuer": {
			ColumnName: "test_embedded.valuer",
			FieldIndex: []int{0, 1},
			GoType:     reflect.TypeOf(&sql.NullString{}),
		},
		"bool": {
			ColumnName: "bool",
			FieldIndex: []int{1},
			GoType:     reflect.TypeOf(true),
		},
		"valuer": {
			ColumnName: "valuer",
			FieldIndex: []int{2},
			GoType:     reflect.TypeOf(&sql.NullString{}),
		},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withTaggedStructField() {

	type TestEmbedded struct {
		Bool   bool
		Valuer *sql.NullString
	}

	type TestStruct struct {
		Embedded TestEmbedded `db:"test_embedded" scan:"notate"`
		Bool     bool
		Valuer   *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"test_embedded.bool": {
			ColumnName: "test_embedded.bool",
			FieldIndex: []int{0, 0},
			GoType:     reflect.TypeOf(true),
		},
		"test_embedded.valuer": {
			ColumnName: "test_embedded.valuer",
			FieldIndex: []int{0, 1},
			GoType:     reflect.TypeOf(&sql.NullString{}),
		},
		"bool": {
			ColumnName: "bool",
			FieldIndex: []int{1},
			GoType:     reflect.TypeOf(true),
		},
		"valuer": {
			ColumnName: "valuer",
			FieldIndex: []int{2},
			GoType:     reflect.TypeOf(&sql.NullString{}),
		},
	}, cm)
}

func (rt *reflectTest) TestGetColumnMap_withTaggedStructPointerField() {

	type TestEmbedded struct {
		Bool   bool
		Valuer *sql.NullString
	}

	type TestStruct struct {
		Embedded *TestEmbedded `db:"test_embedded" scan:"notate"`
		Bool     bool
		Valuer   *sql.NullString
	}
	var ts TestStruct
	cm, err := GetColumnMap(&ts)
	rt.NoError(err)
	rt.Equal(ColumnMap{
		"test_embedded.bool": {
			ColumnName: "test_embedded.bool",
			FieldIndex: []int{0, 0},
			GoType:     reflect.PtrTo(reflect.TypeOf(true)),
			Optional:   true,
		},
		"test_embedded.valuer": {
			ColumnName: "test_embedded.valuer",
			FieldIndex: []int{0, 1},
			GoType:     reflect.PtrTo(reflect.TypeOf(&sql.NullString{})),
			Optional:   true,
		},
		"bool": {
			ColumnName: "bool",
			FieldIndex: []int{1},
			GoType:     reflect.TypeOf(true),
		},
		"valuer": {
			ColumnName: "valuer",
			FieldIndex: []int{2},
			GoType:     reflect.TypeOf(&sql.NullString{}),
		},
	}, cm)
}

func (rt *reflectTest) TestGetTypeInfo() {
	var a int64
	var b []int64
	var c []*time.Time

	t, k := GetTypeInfo(&a, reflect.ValueOf(a))
	rt.Equal(reflect.TypeOf(a), t)
	rt.Equal(reflect.Int64, k)

	t, k = GetTypeInfo(&b, reflect.ValueOf(a))
	rt.Equal(reflect.TypeOf(a), t)
	rt.Equal(reflect.Int64, k)

	t, k = GetTypeInfo(c, reflect.ValueOf(c))
	rt.Equal(reflect.TypeOf(time.Time{}), t)
	rt.Equal(reflect.Struct, k)
}

func (rt *reflectTest) TestSafeGetFieldByIndex() {
	type TestEmbedded struct {
		FieldA int
	}
	type TestEmbeddedPointerStruct struct {
		*TestEmbedded
		FieldB string
	}
	type TestEmbeddedStruct struct {
		TestEmbedded
		FieldB string
	}
	v := reflect.ValueOf(TestEmbeddedPointerStruct{})
	f, isAvailable := SafeGetFieldByIndex(v, []int{0, 0})
	rt.False(isAvailable)
	rt.False(f.IsValid())
	f, isAvailable = SafeGetFieldByIndex(v, []int{1})
	rt.True(isAvailable)
	rt.True(f.IsValid())
	rt.Equal(reflect.String, f.Type().Kind())
	f, isAvailable = SafeGetFieldByIndex(v, []int{})
	rt.True(isAvailable)
	rt.Equal(v, f)

	v = reflect.ValueOf(TestEmbeddedPointerStruct{TestEmbedded: &TestEmbedded{}})
	f, isAvailable = SafeGetFieldByIndex(v, []int{0, 0})
	rt.True(isAvailable)
	rt.True(f.IsValid())
	rt.Equal(reflect.Int, f.Type().Kind())
	f, isAvailable = SafeGetFieldByIndex(v, []int{1})
	rt.True(isAvailable)
	rt.True(f.IsValid())
	rt.Equal(reflect.String, f.Type().Kind())
	f, isAvailable = SafeGetFieldByIndex(v, []int{})
	rt.True(isAvailable)
	rt.Equal(v, f)

	v = reflect.ValueOf(TestEmbeddedStruct{})
	f, isAvailable = SafeGetFieldByIndex(v, []int{0, 0})
	rt.True(isAvailable)
	rt.True(f.IsValid())
	rt.Equal(reflect.Int, f.Type().Kind())
	f, isAvailable = SafeGetFieldByIndex(v, []int{1})
	rt.True(isAvailable)
	rt.True(f.IsValid())
	rt.Equal(reflect.String, f.Type().Kind())
	f, isAvailable = SafeGetFieldByIndex(v, []int{})
	rt.True(isAvailable)
	rt.Equal(v, f)

	v = reflect.ValueOf(TestEmbeddedStruct{TestEmbedded: TestEmbedded{}})
	f, isAvailable = SafeGetFieldByIndex(v, []int{0, 0})
	rt.True(isAvailable)
	rt.True(f.IsValid())
	f, isAvailable = SafeGetFieldByIndex(v, []int{1})
	rt.True(isAvailable)
	rt.True(f.IsValid())
	rt.Equal(reflect.String, f.Type().Kind())
	f, isAvailable = SafeGetFieldByIndex(v, []int{})
	rt.True(isAvailable)
	rt.Equal(v, f)
}

func (rt *reflectTest) TestSafeSetFieldByIndex() {
	type TestEmbedded struct {
		FieldA int
	}
	type TestEmbeddedPointerStruct struct {
		*TestEmbedded
		FieldB string
	}
	type TestEmbeddedStruct struct {
		TestEmbedded
		FieldB string
	}
	var teps TestEmbeddedPointerStruct
	v := reflect.ValueOf(&teps)
	f := SafeSetFieldByIndex(v, []int{}, nil)
	rt.Equal(TestEmbeddedPointerStruct{}, f.Interface())

	f = SafeSetFieldByIndex(v, []int{0, 0}, 1)
	rt.Equal(TestEmbeddedPointerStruct{
		TestEmbedded: &TestEmbedded{FieldA: 1},
	}, f.Interface())

	f = SafeSetFieldByIndex(v, []int{1}, "hello")
	rt.Equal(TestEmbeddedPointerStruct{
		TestEmbedded: &TestEmbedded{FieldA: 1},
		FieldB:       "hello",
	}, f.Interface())
	rt.Equal(TestEmbeddedPointerStruct{
		TestEmbedded: &TestEmbedded{FieldA: 1},
		FieldB:       "hello",
	}, teps)

	var tes TestEmbeddedStruct
	v = reflect.ValueOf(&tes)
	f = SafeSetFieldByIndex(v, []int{}, nil)
	rt.Equal(TestEmbeddedStruct{}, f.Interface())

	f = SafeSetFieldByIndex(v, []int{0, 0}, 1)
	rt.Equal(TestEmbeddedStruct{
		TestEmbedded: TestEmbedded{FieldA: 1},
	}, f.Interface())

	f = SafeSetFieldByIndex(v, []int{1}, "hello")
	rt.Equal(TestEmbeddedStruct{
		TestEmbedded: TestEmbedded{FieldA: 1},
		FieldB:       "hello",
	}, f.Interface())
	rt.Equal(TestEmbeddedStruct{
		TestEmbedded: TestEmbedded{FieldA: 1},
		FieldB:       "hello",
	}, tes)
}

func (rt *reflectTest) TestGetSliceElementType() {
	type MyStruct struct{}

	tests := []struct {
		slice interface{}
		want  reflect.Type
	}{
		{
			slice: []int{},
			want:  reflect.TypeOf(1),
		},
		{
			slice: []*int{},
			want:  reflect.TypeOf(1),
		},
		{
			slice: []MyStruct{},
			want:  reflect.TypeOf(MyStruct{}),
		},
		{
			slice: []*MyStruct{},
			want:  reflect.TypeOf(MyStruct{}),
		},
	}

	for _, tt := range tests {
		sliceVal := reflect.ValueOf(tt.slice)
		elementType := GetSliceElementType(sliceVal)

		rt.Equal(tt.want, elementType)
	}
}

func (rt *reflectTest) TestAppendSliceElement() {
	type MyStruct struct{}

	sliceVal := reflect.Indirect(reflect.ValueOf(&[]MyStruct{}))
	AppendSliceElement(sliceVal, reflect.ValueOf(&MyStruct{}))

	rt.Equal([]MyStruct{{}}, sliceVal.Interface())

	sliceVal = reflect.Indirect(reflect.ValueOf(&[]*MyStruct{}))
	AppendSliceElement(sliceVal, reflect.ValueOf(&MyStruct{}))

	rt.Equal([]*MyStruct{{}}, sliceVal.Interface())
}

func TestReflectSuite(t *testing.T) {
	suite.Run(t, new(reflectTest))
}
