// +build integration

package pgxscan_test

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/randallmlough/pgxscan/testdata"

	"github.com/randallmlough/pgxscan"
	"github.com/stretchr/testify/require"
)

func Test_rows_Variadic(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	var (
		ID     uint32
		Int    int
		Float  float32
		String string
		Time   time.Time
		Bool   bool
		Bytes  []byte
		Slice  []string
		Json   testdata.JSON
	)
	err = pgxscan.NewScanner(rows).Scan(
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
	require.NoError(t, err)
	require.Equal(t, testdata.TestRow1.ID, ID)
	require.Equal(t, testdata.TestRow1.Int, Int)
	require.Equal(t, testdata.TestRow1.Float32, Float)
	require.Equal(t, testdata.TestRow1.String, String)
	require.Equal(t, testdata.TestRow1.Time, Time)
	require.Equal(t, testdata.TestRow1.Bool, Bool)
	require.Equal(t, testdata.TestRow1.Bytes, Bytes)
	require.Equal(t, testdata.TestRow1.StringSlice, Slice)
	require.Equal(t, testdata.TestRow1.JSONB, Json)
}

func Test_rows_VariadicOfPointers(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	var (
		ID     *uint32
		Int    *int
		Float  *float32
		String *string
		Time   *time.Time
		Bool   *bool
		Bytes  *[]byte
		Slice  *[]string
		Json   *testdata.JSON
	)

	err = pgxscan.NewScanner(rows).Scan(
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
	require.NoError(t, err)

	require.Equal(t, testdata.TestRow1.ID, *ID)
	require.Equal(t, testdata.TestRow1.Int, *Int)
	require.Equal(t, testdata.TestRow1.Float32, *Float)
	require.Equal(t, testdata.TestRow1.String, *String)
	require.Equal(t, testdata.TestRow1.Time, *Time)
	require.Equal(t, testdata.TestRow1.Bool, *Bool)
	require.Equal(t, testdata.TestRow1.Bytes, *Bytes)
	require.Equal(t, testdata.TestRow1.StringSlice, *Slice)
	require.Equal(t, testdata.TestRow1.JSONB, *Json)

	require.Equal(t, &testdata.TestRow1.ID, ID)
	require.Equal(t, &testdata.TestRow1.Int, Int)
	require.Equal(t, &testdata.TestRow1.Float32, Float)
	require.Equal(t, &testdata.TestRow1.String, String)
	require.Equal(t, &testdata.TestRow1.Time, Time)
	require.Equal(t, &testdata.TestRow1.Bool, Bool)
	require.Equal(t, &testdata.TestRow1.Bytes, Bytes)
	require.Equal(t, &testdata.TestRow1.StringSlice, Slice)
	require.Equal(t, &testdata.TestRow1.JSONB, Json)
}

func Test_rows_InterfaceSliceOfTypes(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	var (
		ID     uint32
		Int    int
		Float  float32
		String string
		Time   time.Time
		Bool   bool
		Bytes  []byte
		Slice  []string
		Json   testdata.JSON
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
	err = pgxscan.NewScanner(rows).Scan(dst)
	require.NoError(t, err)

	require.Equal(t, testdata.TestRow1.ID, ID)
	require.Equal(t, testdata.TestRow1.Int, Int)
	require.Equal(t, testdata.TestRow1.Float32, Float)
	require.Equal(t, testdata.TestRow1.String, String)
	require.Equal(t, testdata.TestRow1.Time, Time)
	require.Equal(t, testdata.TestRow1.Bool, Bool)
	require.Equal(t, testdata.TestRow1.Bytes, Bytes)
	require.Equal(t, testdata.TestRow1.StringSlice, Slice)
	require.Equal(t, testdata.TestRow1.JSONB, Json)
}
func Test_rows_InterfaceSliceOfStructFields(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	var ts testdata.TestStruct
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
	err = pgxscan.NewScanner(rows).Scan(dst)
	require.NoError(t, err)

	require.Equal(t, testdata.TestStruct{
		ID:          testdata.TestRow1.ID,
		Int:         testdata.TestRow1.Int,
		Float32:     testdata.TestRow1.Float32,
		String:      testdata.TestRow1.String,
		Time:        testdata.TestRow1.Time,
		Bool:        testdata.TestRow1.Bool,
		Bytes:       testdata.TestRow1.Bytes,
		StringSlice: testdata.TestRow1.StringSlice,
		JSONB:       testdata.TestRow1.JSONB,
	}, ts)
}

func Test_rows_ScanStruct(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	var dst testdata.TestStruct
	err = pgxscan.NewScanner(rows).Scan(&dst)
	require.NoError(t, err)

	require.Equal(t, testdata.TestStruct{
		ID:          testdata.TestRow1.ID,
		Int:         testdata.TestRow1.Int,
		Float32:     testdata.TestRow1.Float32,
		String:      testdata.TestRow1.String,
		Time:        testdata.TestRow1.Time,
		Bool:        testdata.TestRow1.Bool,
		Bytes:       testdata.TestRow1.Bytes,
		StringSlice: testdata.TestRow1.StringSlice,
		JSONB:       testdata.TestRow1.JSONB,
	}, dst)
}
func Test_rows_ScanStructSelectOrderDiffFromStructFieldOrder(t *testing.T) {
	cols := []string{"id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b"}
	sort.Strings(cols)

	stmt := fmt.Sprintf(`SELECT %s FROM "test" WHERE id = $1`, strings.Join(cols, ","))
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	var dst testdata.TestStruct
	err = pgxscan.NewScanner(rows).Scan(&dst)
	require.NoError(t, err)

	require.Equal(t, testdata.TestStruct{
		ID:          testdata.TestRow1.ID,
		Int:         testdata.TestRow1.Int,
		Float32:     testdata.TestRow1.Float32,
		String:      testdata.TestRow1.String,
		Time:        testdata.TestRow1.Time,
		Bool:        testdata.TestRow1.Bool,
		Bytes:       testdata.TestRow1.Bytes,
		StringSlice: testdata.TestRow1.StringSlice,
		JSONB:       testdata.TestRow1.JSONB,
	}, dst)
}

func Test_rows_ScanStructSelectLessThanStructFields(t *testing.T) {
	cols := []string{"id", "int", "float_32", "string", "time", "bool"}
	sort.Strings(cols)
	stmt := fmt.Sprintf(`SELECT %s FROM "test" WHERE id = $1`, strings.Join(cols, ","))
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	dst := struct {
		ID      uint32    `db:"id"`
		Int     int       `db:"int"`
		Float32 float32   `db:"float_32"`
		String  string    `db:"string"`
		Time    time.Time `db:"time"`
		Bool    bool      `db:"bool"`
	}{}
	err = pgxscan.NewScanner(rows).Scan(&dst)
	require.NoError(t, err)
	require.Equal(t, struct {
		ID      uint32    `db:"id"`
		Int     int       `db:"int"`
		Float32 float32   `db:"float_32"`
		String  string    `db:"string"`
		Time    time.Time `db:"time"`
		Bool    bool      `db:"bool"`
	}{
		ID:      testdata.TestRow1.ID,
		Int:     testdata.TestRow1.Int,
		Float32: testdata.TestRow1.Float32,
		String:  testdata.TestRow1.String,
		Time:    testdata.TestRow1.Time,
		Bool:    testdata.TestRow1.Bool,
	}, dst)
}
func Test_rows_WantErr_ScanNonPointerStruct(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	scanner := pgxscan.NewScanner(rows)
	var dst testdata.TestStruct
	err = scanner.Scan(dst)
	require.Error(t, err)
}

func Test_rows_ScanPointerStruct(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	dst := &testdata.TestStruct{}
	err = pgxscan.NewScanner(rows).Scan(dst)
	require.NoError(t, err)

	require.Equal(t, &testdata.TestStruct{
		ID:          testdata.TestRow1.ID,
		Int:         testdata.TestRow1.Int,
		Float32:     testdata.TestRow1.Float32,
		String:      testdata.TestRow1.String,
		Time:        testdata.TestRow1.Time,
		Bool:        testdata.TestRow1.Bool,
		Bytes:       testdata.TestRow1.Bytes,
		StringSlice: testdata.TestRow1.StringSlice,
		JSONB:       testdata.TestRow1.JSONB,
	}, dst)
}
func Test_rows_ScanAddrToPointerStruct(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	var dst *testdata.TestStruct
	err = pgxscan.NewScanner(rows).Scan(&dst)
	require.NoError(t, err)

	require.Equal(t, &testdata.TestStruct{
		ID:          testdata.TestRow1.ID,
		Int:         testdata.TestRow1.Int,
		Float32:     testdata.TestRow1.Float32,
		String:      testdata.TestRow1.String,
		Time:        testdata.TestRow1.Time,
		Bool:        testdata.TestRow1.Bool,
		Bytes:       testdata.TestRow1.Bytes,
		StringSlice: testdata.TestRow1.StringSlice,
		JSONB:       testdata.TestRow1.JSONB,
	}, dst)
}

func Test_rows_WantErr_ScanToNilPointerStruct(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	var dst *testdata.TestStruct
	err = pgxscan.NewScanner(rows).Scan(dst)
	require.Error(t, err)

}
func Test_rows_ScanNewPointerStruct(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	dst := new(testdata.TestStruct)
	err = pgxscan.NewScanner(rows).Scan(dst)
	require.NoError(t, err)

	require.Equal(t, &testdata.TestStruct{
		ID:          testdata.TestRow1.ID,
		Int:         testdata.TestRow1.Int,
		Float32:     testdata.TestRow1.Float32,
		String:      testdata.TestRow1.String,
		Time:        testdata.TestRow1.Time,
		Bool:        testdata.TestRow1.Bool,
		Bytes:       testdata.TestRow1.Bytes,
		StringSlice: testdata.TestRow1.StringSlice,
		JSONB:       testdata.TestRow1.JSONB,
	}, dst)
}
func Test_rows_ScanAddrToNewPointerStruct(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	dst := new(testdata.TestStruct)
	err = pgxscan.NewScanner(rows).Scan(&dst)
	require.NoError(t, err)

	require.Equal(t, &testdata.TestStruct{
		ID:          testdata.TestRow1.ID,
		Int:         testdata.TestRow1.Int,
		Float32:     testdata.TestRow1.Float32,
		String:      testdata.TestRow1.String,
		Time:        testdata.TestRow1.Time,
		Bool:        testdata.TestRow1.Bool,
		Bytes:       testdata.TestRow1.Bytes,
		StringSlice: testdata.TestRow1.StringSlice,
		JSONB:       testdata.TestRow1.JSONB,
	}, dst)
}

func Test_rows_SliceStructScan(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" ORDER BY "id" ASC LIMIT 2`
	rows, err := newTestDB(t).Query(context.Background(), stmt)
	require.NoError(t, err)

	var dst []testdata.TestStruct
	err = pgxscan.NewScanner(rows).Scan(&dst)
	require.NoError(t, err)
	require.Equal(t, []testdata.TestStruct{
		{
			ID:          testdata.TestRow1.ID,
			Int:         testdata.TestRow1.Int,
			Float32:     testdata.TestRow1.Float32,
			String:      testdata.TestRow1.String,
			Time:        testdata.TestRow1.Time,
			Bool:        testdata.TestRow1.Bool,
			Bytes:       testdata.TestRow1.Bytes,
			StringSlice: testdata.TestRow1.StringSlice,
			JSONB:       testdata.TestRow1.JSONB,
		},
		{
			ID:          testdata.TestRow2.ID,
			Int:         testdata.TestRow2.Int,
			Float32:     testdata.TestRow2.Float32,
			String:      testdata.TestRow2.String,
			Time:        testdata.TestRow2.Time,
			Bool:        testdata.TestRow2.Bool,
			Bytes:       testdata.TestRow2.Bytes,
			StringSlice: testdata.TestRow2.StringSlice,
			JSONB:       testdata.TestRow2.JSONB,
		},
	}, dst)
}

func Test_rows_JoinTable(t *testing.T) {
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
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

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
	err = pgxscan.NewScanner(rows).Scan(&user)
	require.NoError(t, err)

	require.Equal(t, User{
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

func Test_rows_JoinTableWithNotationColumn(t *testing.T) {
	stmt := `
	SELECT users.*,
	       0 as "notate:address", -- delimiter column
	       address.*
	FROM users, address
	WHERE users.id = $1
	  AND address.user_id = users.id
	`
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

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
	scanner := pgxscan.NewScanner(rows, pgxscan.MatchAllColumns(false))
	if err := scanner.Scan(&user); err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, User{
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

func Test_rows_JoinConflictTable(t *testing.T) {
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
	rows, err := newTestDB(t).Query(context.Background(), stmt)
	require.NoError(t, err)

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
	if err := pgxscan.NewScanner(rows).Scan(&joined); err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, Joined{
		A:  123,
		C1: Conflicting1{A: 0, B: 0},
		C2: Conflicting2{B: 1, C: 1},
		B:  456,
		C3: Conflicting3{A: 2, B: 2, C: 2},
	}, joined)
}

func Test_rows_IdentifyNullColumn(t *testing.T) {
	stmt := `
	SELECT users.*
	FROM users
	WHERE users.id = $1
	`
	// user 10 has NULL name
	rows, err := newTestDB(t).Query(context.Background(), stmt, 10)
	require.NoError(t, err)

	type User struct {
		ID    uint32
		Name  string // NULL cannot be handled here, thus, error
		Email string
	}

	var user User
	err = pgxscan.NewScanner(rows).Scan(&user)
	require.Equal(
		t,
		err.Error(),
		"can't scan into dest[1] (field 'name'): cannot assign NULL to *string",
	)
}

func Test_rows_IdentifyWrongTypeForColumn(t *testing.T) {
	stmt := `
	SELECT users.*
	FROM users
	WHERE users.id = $1
	`
	// user 10 has NULL name
	rows, err := newTestDB(t).Query(context.Background(), stmt, 1)
	require.NoError(t, err)

	type User struct {
		ID    uint32
		Name  string
		Email int // this should be string! it will fail
	}

	var user User
	err = pgxscan.NewScanner(rows).Scan(&user)
	require.Equal(
		t,
		err.Error(),
		"can't scan into dest[2] (field 'email'): unable to assign to *int",
	)
}

func Benchmark_rows_JoinPgxBaseline(b *testing.B) {
	type (
		Address struct {
			ID    uint32
			Line1 string
			City  string
		}
		User struct {
			ID      uint32
			Name    string
			Email   string
			Address Address
		}
	)

	stmt := `
WITH usr AS (
	SELECT
		id, name, email
	FROM
		"users"
	WHERE
		"id" = $1
),
addr AS (
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
	usr.*, addr.*
FROM
	usr, addr
`

	var user User
	testDb := newTestDB(b)
	for i := 0; i < b.N; i++ {
		rows, err := testDb.Query(context.Background(), stmt, 1)
		require.NoError(b, err)

		if rows.Next() {
			err := rows.Scan(
				&user.ID,
				&user.Name,
				&user.Email,
				&user.Address.ID,
				&user.Address.Line1,
				&user.Address.City,
			)
			require.NoError(b, err)
		}
	}
}

func Benchmark_rows_JoinTableWithStructNotation(b *testing.B) {
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

	var user User
	testDb := newTestDB(b)
	for i := 0; i < b.N; i++ {
		rows, err := testDb.Query(context.Background(), stmt, 1)
		require.NoError(b, err)

		scanner := pgxscan.NewScanner(rows)
		err = scanner.Scan(&user)
		require.NoError(b, err)
	}
}

func Benchmark_rows_JoinTableWithSqlNotation(b *testing.B) {
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

	stmt := `
	SELECT users.*,
	       0 as "notate:address", -- delimiter column
	       address.*
	FROM users, address
	WHERE users.id = $1
	  AND address.user_id = users.id
	`

	var user User
	testDb := newTestDB(b)
	for i := 0; i < b.N; i++ {
		rows, err := testDb.Query(context.Background(), stmt, 1)
		require.NoError(b, err)

		scanner := pgxscan.NewScanner(rows, pgxscan.MatchAllColumns(false))
		err = scanner.Scan(&user)
		require.NoError(b, err)
	}
}
