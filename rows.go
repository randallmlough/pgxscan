package pgxscan

import (
	"errors"
	sqlmaper "github.com/randallmlough/pgxscan/internal/sqlmapper"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v4"
)

type rows struct {
	rows    pgx.Rows
	columns []string
	cfg     *Config
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

var ErrColumnNotateSyntax = errors.New("column notate syntax is invalid: expecting \"notate:name\"")
var QueryColumnNotatePrefix = "notate:"

// This function returns a list of column names from SQL columns.
//
// It renames SQL notated columns, when needed, for column names who start
// with 'notate:' prefix (you can override prefix with QueryColumnNotatePrefix
// variable)
//
// From columns:
//
//	["a", "b", "notate:whatever", "a", "b"]
//
// Will return:
//
//	["a", "b", "notate:whatever", "whatever.a", "whatever.b"]
//
// Background info:
//
// In order to allow complex queries and prevent having to expand on all column
// names for complex mappings, delimiter columns are used to notate results
// from postgres.
//
// Imagine we have `SELECT a.*, b.* FROM ...` if both _a_ and _b_ tables have
// a field name _id_ there would be no way for us to map it to a struct with
// a couple of nested structs, however, if  we use column notation we can
// rewrite the query as:
//
//	SELECT 0 as "notate:a",
//	       a.*,
//	       0 as "notate:b",
//	       b.*
//	FROM ...
//
// This way, everything that comes after column "notate:a" will be treated as if
// we would have defined an alias for each column named "a.<col>", and so on
//
// These notations allow zero (using "notate:" with nothing after colon),
// one level (like the example above) or many levels of notations (just do
// "notate:a.sub1.sub2")
//
// This helps map values to struct with simple queries without having to list
// all columns in the SQL.
func GetColumnNames(rows *pgx.Rows) ([]string, error) {
	cols := make([]string, 0, len((*rows).FieldDescriptions()))

	notatePrefix := ""
	for _, field := range (*rows).FieldDescriptions() {
		colName := string(field.Name)

		// if starts by 'notate:' use what comes after that as the prefix for
		// all column definitions moving forward
		if strings.HasPrefix(colName, QueryColumnNotatePrefix) {
			// "notate: a.b.c" -> ["notate:", " a.b.c"]
			splitted := strings.Split(colName, QueryColumnNotatePrefix)
			if len(splitted) != 2 {
				return nil, ErrColumnNotateSyntax
			}

			// "a.b.c" or ""
			trimmed := strings.TrimSpace(splitted[1])
			if len(trimmed) == 0 {
				notatePrefix = ""
			} else {
				notatePrefix = strings.TrimRight(trimmed, ".") + "."
			}
		} else {
			// notatePrefix can be empty, and thus, have no effect
			colName = notatePrefix + colName
		}

		cols = append(cols, colName)
	}

	return cols, nil
}

// ScanStruct will scan the current row into i.
func (r *rows) Scan(i ...interface{}) (err error) {
	if i == nil {
		return nil
	} else if isVariadic(i...) {
		return r.ScanVal(i...)
	} else if ii, ok := i[0].([]interface{}); ok {
		return r.ScanVal(ii...)
	}

	val, valErr := validate(i[0])
	if valErr != nil {
		err = valErr
		return
	}

	var rowCount int64
	defer func() {
		r.Close()
		if r.cfg.ReturnErrNoRowsForRows && err == nil && rowCount == 0 {
			err = pgx.ErrNoRows
		}
	}()
	switch val.Kind() {
	case reflect.Slice:
		sliceOf := sqlmaper.GetSliceElementType(val)
		for r.Next() {
			sliceVal := reflect.New(sliceOf)

			cols, colErr := GetColumnNames(&r.rows)
			if colErr != nil {
				err = colErr
				return
			}
			if ssErr := ScanStruct(r.rows.Scan, sliceVal.Interface(), cols, r.cfg.MatchAllColumnsToStruct); ssErr != nil {
				err = ssErr
				return
			}
			sqlmaper.AppendSliceElement(val, sliceVal)
			rowCount++
		}
	case reflect.Struct:
		for r.Next() {
			if val.CanAddr() {
				cols, colErr := GetColumnNames(&r.rows)
				if colErr != nil {
					err = colErr
					return
				}

				if ssErr := ScanStruct(r.rows.Scan, val.Addr().Interface(), cols, r.cfg.MatchAllColumnsToStruct); ssErr != nil {
					err = ssErr
					return
				}
			}
			rowCount++
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
