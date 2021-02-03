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
func NewScanner(src Scanner, opts ...Option) Scanner {
	cfg := &Config{
		ReturnErrNoRowsForRows: true,
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
	ReturnErrNoRowsForRows bool
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
