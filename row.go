package pgxscan

import (
	"github.com/jackc/pgx/v4"
)

type row struct {
	row     pgx.Row
	columns []string
}

func (r *row) Scan(i ...interface{}) error {
	if i == nil {
		return nil
	} else if ii, ok := i[0].([]interface{}); ok {
		if err := r.row.Scan(ii...); err != nil {
			return err
		}
	} else {
		if err := r.row.Scan(i...); err != nil {
			return err
		}
	}

	return nil
}

func (r *row) SetCols(cols ...string) Scanner {
	r.columns = cols
	return r
}
