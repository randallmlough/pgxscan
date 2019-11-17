package pgxscan

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/randallmlough/sqlmaper"
	"reflect"
	"time"
)

type (
	// Scanner knows how to scan sql.Rows into structs.
	Scanner interface {
		Scan(v ...interface{}) error
	}

	scannerFunc func(i ...interface{}) error
)

func unableToFindFieldError(col string) error {
	return fmt.Errorf(`unable to find corresponding field to column "%s" returned by query`, col)
}

// NewScanner takes in a scanner returns a scanner
// Since the pgx row and rows interface both have a `Scan(v ...interface{}) error` method,
// either one can be passed as the argument and scanner will take care of the rest.
func NewScanner(src Scanner) Scanner {
	switch s := src.(type) {
	case pgx.Rows:
		return &rows{rows: s}
	case pgx.Row:
		return &row{row: s}
	}
	return nil
}

var ErrNoCols = errors.New("columns can not be nil")

// ScanStruct will scan the current row into i.
func ScanStruct(scan scannerFunc, i interface{}, cols []string) error {
	if cols == nil {
		return ErrNoCols
	}
	cm, err := sqlmaper.GetColumnMap(i)
	if err != nil {
		return err
	}

	scans := make([]interface{}, len(cols))
	for idx, col := range cols {
		data, ok := cm[col]
		switch {
		case !ok:
			return unableToFindFieldError(col)
		default:
			scans[idx] = reflect.New(data.GoType).Interface()
		}
	}

	if err := scan(scans...); err != nil {
		return err
	}

	record := make(map[string]interface{}, len(cols))
	for index, col := range cols {
		record[col] = scans[index]
	}

	sqlmaper.AssignStructVals(i, record, cm)

	return nil
}

func validate(scan Scanner, i interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(i)
	if !sqlmaper.IsPointer(val.Kind()) {
		return reflect.Value{}, errors.New("destination must be a pointer")
	}

	if val.Kind() == reflect.Ptr && !val.Elem().CanSet() {
		return reflect.Value{}, errors.New("destination must be initialized")
	}

	val = reflect.Indirect(val)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			v := reflect.New(val.Type().Elem())
			if err := scan.Scan(v.Interface()); err != nil {
				return reflect.Value{}, err
			}
			val.Set(v)
		} else {
			val = val.Elem()
		}
	}
	return val, nil
}
func isVariadic(i ...interface{}) bool {
	switch len(i) {
	case 0:
		return false
	case 1:
		if isBuiltin(i[0]) {
			return true
		}
		return false
	default:
		return true
	}
}

func isBuiltin(i interface{}) bool {
	switch i.(type) {
	case
		string,
		uint, uint8, uint16, uint32, uint64,
		int, int8, int16, int32, int64,
		complex64, complex128,
		float32, float64,
		bool:
		return true
	case
		*string,
		*uint, *uint8, *uint16, *uint32, *uint64,
		*int, *int8, *int16, *int32, *int64,
		*complex64, *complex128,
		*float32, *float64,
		*bool:
		return true
	case
		[]string,
		[]uint, []uint8, []uint16, []uint32, []uint64,
		[]int, []int8, []int16, []int32, []int64,
		[]complex64, []complex128,
		[]float32, []float64,
		[]bool:
		return true
	case
		*[]string,
		*[]uint, *[]uint8, *[]uint16, *[]uint32, *[]uint64,
		*[]int, *[]int8, *[]int16, *[]int32, *[]int64,
		*[]complex64, *[]complex128,
		*[]float32, *[]float64,
		*[]bool:
		return true
	case time.Time, *time.Time:
		return true
	case []time.Time, *[]time.Time:
		return true
	default:
		return false
	}
}
