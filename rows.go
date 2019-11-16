package pgxscan

import (
	"github.com/jackc/pgx/v4"
	"github.com/randallmlough/sqlmapping"
	"reflect"
)

type rows struct {
	rows    pgx.Rows
	columns []string
}

// Next prepares the next row for Scanning. See sql.Rows#Next for more
// information.
func (r *rows) Next() bool {
	return r.rows.Next()
}

// Err returns the error, if any that was encountered during iteration. See
// sql.Rows#Err for more information.
func (r *rows) Err() error {
	return r.rows.Err()
}

// ScanStruct will scan the current row into i.
func (r *rows) Scan(i ...interface{}) error {
	if i == nil {
		return nil
	} else if isVariadic(i...) {
		return r.ScanVal(i...)
	} else if ii, ok := i[0].([]interface{}); ok {
		return r.ScanVal(ii...)
	}

	val, err := validate(r, i[0])
	if err != nil {
		return err
	}

	defer r.Close()
	switch val.Kind() {
	case reflect.Slice:
		sliceOf := sqlmapping.GetSliceElementType(val)
		for r.Next() {
			sliceVal := reflect.New(sliceOf)

			cols := make([]string, 0, len(r.rows.FieldDescriptions()))
			for _, field := range r.rows.FieldDescriptions() {
				cols = append(cols, string(field.Name))
			}
			if err := ScanStruct(r.rows.Scan, sliceVal.Interface(), cols); err != nil {
				return err
			}
			sqlmapping.AppendSliceElement(val, sliceVal)
		}
	case reflect.Struct:
		for r.Next() {
			if val.CanAddr() {
				cols := make([]string, 0, len(r.rows.FieldDescriptions()))
				for _, field := range r.rows.FieldDescriptions() {
					cols = append(cols, string(field.Name))
				}
				if err := ScanStruct(r.rows.Scan, val.Addr().Interface(), cols); err != nil {
					return err
				}
			}
		}
	}
	return r.Err()
}

// ScanVal will scan the current row and column into i.
func (r *rows) ScanVal(v ...interface{}) error {
	defer r.Close()
	for r.Next() {
		if err := r.rows.Scan(v...); err != nil {
			return err
		}
	}
	return r.rows.Err()
}

// Close closes the Rows, preventing further enumeration. See sql.Rows#Close
// for more info.
func (r *rows) Close() {
	r.rows.Close()
}

func (r *rows) SetCols(cols ...string) Scanner {
	r.columns = cols
	return r
}
