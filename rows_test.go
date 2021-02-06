package pgxscan

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/randallmlough/pgxscan/internal/test"
	"github.com/stretchr/testify/suite"
	"sort"
	"strings"
	"testing"
	"time"
)

func Test_Rows(t *testing.T) {
	suite.Run(t, new(rowsTest))
}

type rowsTest struct {
	suite.Suite
	db *pgx.Conn
}

func (rt *rowsTest) SetupSuite() {
	conn, err := test.NewConnection()
	rt.NoError(err)
	rt.db = conn
}

func (rt *rowsTest) Test_rows_Variadic() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	var (
		ID     uint32
		Int    int
		Float  float32
		String string
		Time   time.Time
		Bool   bool
		Bytes  []byte
		Slice  []string
		Json   test.JSON
	)
	scanner := NewScanner(rows)
	err = scanner.Scan(
		&ID,
		&Int,
		&Float,
		&String,
		&Time,
		&Bool,
		&Bytes,
		&Slice,
		&Json,
	)

	rt.NoError(err)
	rt.Equal(test.TestRow1.ID, ID)
	rt.Equal(test.TestRow1.Int, Int)
	rt.Equal(test.TestRow1.Float32, Float)
	rt.Equal(test.TestRow1.String, String)
	rt.Equal(test.TestRow1.Time, Time)
	rt.Equal(test.TestRow1.Bool, Bool)
	rt.Equal(test.TestRow1.Bytes, Bytes)
	rt.Equal(test.TestRow1.StringSlice, Slice)
	rt.Equal(test.TestRow1.JSONB, Json)
}
func (rt *rowsTest) Test_rows_VariadicOfPointers() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	var (
		ID     *uint32
		Int    *int
		Float  *float32
		String *string
		Time   *time.Time
		Bool   *bool
		Bytes  *[]byte
		Slice  *[]string
		Json   *test.JSON
	)
	scanner := NewScanner(rows)
	err = scanner.Scan(
		&ID,
		&Int,
		&Float,
		&String,
		&Time,
		&Bool,
		&Bytes,
		&Slice,
		&Json,
	)
	rt.NoError(err)

	rt.Equal(test.TestRow1.ID, *ID)
	rt.Equal(test.TestRow1.Int, *Int)
	rt.Equal(test.TestRow1.Float32, *Float)
	rt.Equal(test.TestRow1.String, *String)
	rt.Equal(test.TestRow1.Time, *Time)
	rt.Equal(test.TestRow1.Bool, *Bool)
	rt.Equal(test.TestRow1.Bytes, *Bytes)
	rt.Equal(test.TestRow1.StringSlice, *Slice)
	rt.Equal(test.TestRow1.JSONB, *Json)

	rt.Equal(&test.TestRow1.ID, ID)
	rt.Equal(&test.TestRow1.Int, Int)
	rt.Equal(&test.TestRow1.Float32, Float)
	rt.Equal(&test.TestRow1.String, String)
	rt.Equal(&test.TestRow1.Time, Time)
	rt.Equal(&test.TestRow1.Bool, Bool)
	rt.Equal(&test.TestRow1.Bytes, Bytes)
	rt.Equal(&test.TestRow1.StringSlice, Slice)
	rt.Equal(&test.TestRow1.JSONB, Json)
}

func (rt *rowsTest) Test_rows_InterfaceSliceOfTypes() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	var (
		ID     uint32
		Int    int
		Float  float32
		String string
		Time   time.Time
		Bool   bool
		Bytes  []byte
		Slice  []string
		Json   test.JSON
	)

	dst := []interface{}{
		&ID,
		&Int,
		&Float,
		&String,
		&Time,
		&Bool,
		&Bytes,
		&Slice,
		&Json,
	}
	err = scanner.Scan(dst)
	rt.NoError(err)

	rt.Equal(test.TestRow1.ID, ID)
	rt.Equal(test.TestRow1.Int, Int)
	rt.Equal(test.TestRow1.Float32, Float)
	rt.Equal(test.TestRow1.String, String)
	rt.Equal(test.TestRow1.Time, Time)
	rt.Equal(test.TestRow1.Bool, Bool)
	rt.Equal(test.TestRow1.Bytes, Bytes)
	rt.Equal(test.TestRow1.StringSlice, Slice)
	rt.Equal(test.TestRow1.JSONB, Json)
}
func (rt *rowsTest) Test_rows_InterfaceSliceOfStructFields() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	var ts test.TestStruct
	dst := []interface{}{
		&ts.ID,
		&ts.Int,
		&ts.Float32,
		&ts.String,
		&ts.Time,
		&ts.Bool,
		&ts.Bytes,
		&ts.StringSlice,
		&ts.JSONB,
	}
	err = scanner.Scan(dst)
	rt.NoError(err)

	rt.Equal(test.TestStruct{
		ID:          test.TestRow1.ID,
		Int:         test.TestRow1.Int,
		Float32:     test.TestRow1.Float32,
		String:      test.TestRow1.String,
		Time:        test.TestRow1.Time,
		Bool:        test.TestRow1.Bool,
		Bytes:       test.TestRow1.Bytes,
		StringSlice: test.TestRow1.StringSlice,
		JSONB:       test.TestRow1.JSONB,
	}, ts)
}

func (rt *rowsTest) Test_rows_ScanStruct() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	var dst test.TestStruct
	err = scanner.Scan(&dst)
	rt.NoError(err)

	rt.Equal(test.TestStruct{
		ID:          test.TestRow1.ID,
		Int:         test.TestRow1.Int,
		Float32:     test.TestRow1.Float32,
		String:      test.TestRow1.String,
		Time:        test.TestRow1.Time,
		Bool:        test.TestRow1.Bool,
		Bytes:       test.TestRow1.Bytes,
		StringSlice: test.TestRow1.StringSlice,
		JSONB:       test.TestRow1.JSONB,
	}, dst)
}
func (rt *rowsTest) Test_rows_ScanStructSelectOrderDiffFromStructFieldOrder() {
	cols := []string{"id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b"}
	sort.Strings(cols)
	stmt := fmt.Sprintf(`SELECT %s FROM "test" WHERE id = $1`, strings.Join(cols, ","))
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	var dst test.TestStruct
	err = scanner.Scan(&dst)
	rt.NoError(err)

	rt.Equal(test.TestStruct{
		ID:          test.TestRow1.ID,
		Int:         test.TestRow1.Int,
		Float32:     test.TestRow1.Float32,
		String:      test.TestRow1.String,
		Time:        test.TestRow1.Time,
		Bool:        test.TestRow1.Bool,
		Bytes:       test.TestRow1.Bytes,
		StringSlice: test.TestRow1.StringSlice,
		JSONB:       test.TestRow1.JSONB,
	}, dst)
}
func (rt *rowsTest) Test_rows_ScanStructSelectLessThanStructFields() {
	cols := []string{"id", "int", "float_32", "string", "time", "bool"}
	sort.Strings(cols)
	stmt := fmt.Sprintf(`SELECT %s FROM "test" WHERE id = $1`, strings.Join(cols, ","))
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	dst := struct {
		ID      uint32    `db:"id"`
		Int     int       `db:"int"`
		Float32 float32   `db:"float_32"`
		String  string    `db:"string"`
		Time    time.Time `db:"time"`
		Bool    bool      `db:"bool"`
	}{}
	err = scanner.Scan(&dst)
	rt.NoError(err)
	rt.Equal(struct {
		ID      uint32    `db:"id"`
		Int     int       `db:"int"`
		Float32 float32   `db:"float_32"`
		String  string    `db:"string"`
		Time    time.Time `db:"time"`
		Bool    bool      `db:"bool"`
	}{
		ID:      test.TestRow1.ID,
		Int:     test.TestRow1.Int,
		Float32: test.TestRow1.Float32,
		String:  test.TestRow1.String,
		Time:    test.TestRow1.Time,
		Bool:    test.TestRow1.Bool,
	}, dst)
}
func (rt *rowsTest) Test_rows_WantErr_ScanNonPointerStruct() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	var dst test.TestStruct
	err = scanner.Scan(dst)
	rt.Error(err)
}

func (rt *rowsTest) Test_rows_ScanPointerStruct() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	dst := &test.TestStruct{}
	err = scanner.Scan(dst)
	rt.NoError(err)

	rt.Equal(&test.TestStruct{
		ID:          test.TestRow1.ID,
		Int:         test.TestRow1.Int,
		Float32:     test.TestRow1.Float32,
		String:      test.TestRow1.String,
		Time:        test.TestRow1.Time,
		Bool:        test.TestRow1.Bool,
		Bytes:       test.TestRow1.Bytes,
		StringSlice: test.TestRow1.StringSlice,
		JSONB:       test.TestRow1.JSONB,
	}, dst)
}
func (rt *rowsTest) Test_rows_ScanAddrToPointerStruct() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	var dst *test.TestStruct
	err = scanner.Scan(&dst)
	rt.NoError(err)

	rt.Equal(&test.TestStruct{
		ID:          test.TestRow1.ID,
		Int:         test.TestRow1.Int,
		Float32:     test.TestRow1.Float32,
		String:      test.TestRow1.String,
		Time:        test.TestRow1.Time,
		Bool:        test.TestRow1.Bool,
		Bytes:       test.TestRow1.Bytes,
		StringSlice: test.TestRow1.StringSlice,
		JSONB:       test.TestRow1.JSONB,
	}, dst)
}

func (rt *rowsTest) Test_rows_WantErr_ScanToNilPointerStruct() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	var dst *test.TestStruct
	err = scanner.Scan(dst)
	rt.Error(err)

}
func (rt *rowsTest) Test_rows_ScanNewPointerStruct() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	dst := new(test.TestStruct)
	err = scanner.Scan(dst)
	rt.NoError(err)

	rt.Equal(&test.TestStruct{
		ID:          test.TestRow1.ID,
		Int:         test.TestRow1.Int,
		Float32:     test.TestRow1.Float32,
		String:      test.TestRow1.String,
		Time:        test.TestRow1.Time,
		Bool:        test.TestRow1.Bool,
		Bytes:       test.TestRow1.Bytes,
		StringSlice: test.TestRow1.StringSlice,
		JSONB:       test.TestRow1.JSONB,
	}, dst)
}
func (rt *rowsTest) Test_rows_ScanAddrToNewPointerStruct() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	scanner := NewScanner(rows)
	dst := new(test.TestStruct)
	err = scanner.Scan(&dst)
	rt.NoError(err)

	rt.Equal(&test.TestStruct{
		ID:          test.TestRow1.ID,
		Int:         test.TestRow1.Int,
		Float32:     test.TestRow1.Float32,
		String:      test.TestRow1.String,
		Time:        test.TestRow1.Time,
		Bool:        test.TestRow1.Bool,
		Bytes:       test.TestRow1.Bytes,
		StringSlice: test.TestRow1.StringSlice,
		JSONB:       test.TestRow1.JSONB,
	}, dst)
}

func (rt *rowsTest) Test_rows_SliceStructScan() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" ORDER BY "id" ASC LIMIT 2`
	rows, err := rt.db.Query(context.Background(), stmt)
	rt.NoError(err)

	scanner := NewScanner(rows)
	var dst []test.TestStruct
	err = scanner.Scan(&dst)
	rt.NoError(err)
	rt.Equal([]test.TestStruct{
		{
			ID:          test.TestRow1.ID,
			Int:         test.TestRow1.Int,
			Float32:     test.TestRow1.Float32,
			String:      test.TestRow1.String,
			Time:        test.TestRow1.Time,
			Bool:        test.TestRow1.Bool,
			Bytes:       test.TestRow1.Bytes,
			StringSlice: test.TestRow1.StringSlice,
			JSONB:       test.TestRow1.JSONB,
		},
		{
			ID:          test.TestRow2.ID,
			Int:         test.TestRow2.Int,
			Float32:     test.TestRow2.Float32,
			String:      test.TestRow2.String,
			Time:        test.TestRow2.Time,
			Bool:        test.TestRow2.Bool,
			Bytes:       test.TestRow2.Bytes,
			StringSlice: test.TestRow2.StringSlice,
			JSONB:       test.TestRow2.JSONB,
		},
	}, dst)
}

func (rt *rowsTest) Test_rows_JoinTable() {
	stmt := `
WITH usr AS (
	SELECT
		*
	FROM
		"users"
	WHERE
		"id" = $1
),
addresses AS (
	SELECT
		"address"."id" AS "address.id",
		"line_1" AS "address.line_1",
		"city" AS "address.city"
	FROM
		"address", usr
	WHERE
		"user_id" = usr."id"
)
SELECT
	usr.*, addresses.*
FROM
	usr,
	addresses`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	type (
		Address struct {
			ID    uint32
			Line1 string `db:"line_1"`
			City  string
		}
		User struct {
			ID      uint32
			Name    string
			Email   string
			Address Address `scan:"notate"`
		}
	)
	var user User
	if err := NewScanner(rows).Scan(&user); err != nil {
		rt.NoError(err)
	}

	rt.Equal(User{
		ID:    1,
		Name:  "user01",
		Email: "user01@email.com",
		Address: Address{
			ID:    1,
			Line1: "line01_user01",
			City:  "city01",
		},
	}, user)
}

func (rt *rowsTest) Test_rows_JoinTableWithNotationColumn() {
	stmt := `
	SELECT users.*, 
	       0 as "notate:address", -- delimiter column
	       address.*
	FROM users, address
	WHERE users.id = $1 
	  AND address.user_id = users.id
	`
	rows, err := rt.db.Query(context.Background(), stmt, 1)
	rt.NoError(err)

	type (
		Address struct {
			ID    uint32
			Line1 string `db:"line_1"`
			City  string
		}
		User struct {
			ID      uint32
			Name    string
			Email   string
			Address Address `scan:"notate"`
		}
	)
	var user User
	scanner := NewScanner(rows, MatchAllColumns(false))
	if err := scanner.Scan(&user); err != nil {
		rt.NoError(err)
	}

	rt.Equal(User{
		ID:    1,
		Name:  "user01",
		Email: "user01@email.com",
		Address: Address{
			ID:    1,
			Line1: "line01_user01",
			City:  "city01",
		},
	}, user)
}

func (rt *rowsTest) Test_rows_JoinConflictTable() {
	stmt := `
      SELECT  123 as A,

              0 as "notate:c1",
              c1.*,
              
              -127 as "notate:c2",  -- the number, as long as it is [-127, 128] it does not matter
              c2.*,
              
              0 as "notate:", -- remove notations (just for the sake of testing)
              456 as B,
              
              3 as "notate:c3",
              c3.*
							
        FROM conflicting1 as c1
   LEFT JOIN conflicting2 as c2 ON c2.b = c1.b + 1
   LEFT JOIN conflicting3 as c3 on c3.a = c1.a + 2
       WHERE c1.a = 0
		`
	rows, err := rt.db.Query(context.Background(), stmt)
	rt.NoError(err)

	type (
		Conflicting1 struct {
			A uint32
			B uint32
		}
		Conflicting2 struct {
			B uint32
			C uint32
		}
		Conflicting3 struct {
			A uint32
			B uint32
			C uint32
		}
		Joined struct {
			A  uint32
			C1 Conflicting1 `scan:"notate" db:"c1"`
			C2 Conflicting2 `scan:"notate" db:"c2"`
			C3 Conflicting3 `scan:"notate" db:"c3"`
			B  uint32
		}
	)

	var joined Joined
	if err := NewScanner(rows).Scan(&joined); err != nil {
		rt.NoError(err)
	}

	rt.Equal(Joined{
		A:  123,
		C1: Conflicting1{A: 0, B: 0},
		C2: Conflicting2{B: 1, C: 1},
		B:  456,
		C3: Conflicting3{A: 2, B: 2, C: 2},
	}, joined)
}
