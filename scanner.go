package pgxscan

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/randallmlough/sqlmaper"
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
func NewScanner(src Scanner, opts ...Option) Scanner {
	cfg := &Config{
		ReturnErrNoRowsForRows:  true,
		MatchAllColumnsToStruct: true,
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	switch s := src.(type) {
	case pgx.Rows:
		return &rows{rows: s, cfg: cfg}
	case pgx.Row:
		return &row{row: s, cfg: cfg}
	}
	return nil
}

type Config struct {
	ReturnErrNoRowsForRows  bool
	MatchAllColumnsToStruct bool
}

type Option interface {
	apply(*Config)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Config)

func (f optionFunc) apply(s *Config) {
	f(s)
}

// ErrNoRowsQuery sets whether or not a pgx.ErrNoRows error should be returned on a query that has no rows
func ErrNoRowsQuery(b bool) Option {
	return optionFunc(func(cfg *Config) {
		cfg.ReturnErrNoRowsForRows = b
	})
}

// MatchAllColumns sets whether or not a unableToFindFieldError error
// should be returned on a query that has more columns than fields in the struct
func MatchAllColumns(b bool) Option {
	return optionFunc(func(cfg *Config) {
		cfg.MatchAllColumnsToStruct = b
	})
}

var ErrNoCols = errors.New("columns can not be nil")

// ScanStruct will scan the current row into i.
// When matchAllColumnsToStruct is false, it will not complain about extra columns
// in the result set that are not mapped to the columns in the struct, or, said
// another way, it will allow unmapped items, which can, sometimes, be convenient
func ScanStruct(scan scannerFunc, i interface{}, cols []string, matchAllColumnsToStruct bool) error {
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
		case strings.HasPrefix(col, QueryColumnNotatePrefix):
			// notated columns are always skipped
			scans[idx] = new(int8)
		case !ok:
			if matchAllColumnsToStruct {
				return unableToFindFieldError(col)
			}
		default:
			scans[idx] = reflect.New(data.GoType).Interface()
		}
	}

	if err := scan(scans...); err != nil {
		// identify the offending field in case types do not match, very useful
		// when using this library
		var scanErr pgx.ScanArgError
		if errors.As(err, &scanErr) {
			return fmt.Errorf("can't scan into dest[%d] (field '%s'): %s", scanErr.ColumnIndex, cols[scanErr.ColumnIndex], scanErr.Err)
		}
		return err
	}

	record := make(map[string]interface{}, len(cols))
	for index, col := range cols {
		if cm[col].Optional {
			scanVal := reflect.ValueOf(scans[index])
			// If the type is optional, then we selectively unwind the pointer chain
			if !scanVal.Elem().IsNil() {
				record[col] = scanVal.Elem().Interface()
			}
		} else {
			record[col] = scans[index]
		}
	}

	sqlmaper.AssignStructVals(i, record, cm)

	return nil
}

func validate(i interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("destination must be a pointer")
	}

	if !val.Elem().CanSet() {
		return reflect.Value{}, errors.New("destination must be initialized. Don't use var foo *Foo. Use foo := new(Foo) or foo := &Foo{}")
	}

	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			v := reflect.New(val.Type().Elem())
			// TODO: refactoring required to handle non initialized nil values like, `var foo *Foo`
			// the previous recursion call below worked, however, it limits the possibility of doing post processing after close,
			// such as rows returned, which is more useful than accepting nil values.
			//if err := scan.Scan(v.Interface()); err != nil {
			//	return reflect.Value{}, err
			//}
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
