package pgxscan

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/randallmlough/pgxscan/internal/test"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func Test_Row(t *testing.T) {
	suite.Run(t, new(rowTest))
}

type rowTest struct {
	suite.Suite
	db *pgx.Conn
}

func (rt *rowTest) SetupSuite() {
	conn, err := test.NewConnection()
	rt.NoError(err)
	rt.db = conn
}

func (rt *rowTest) Test_row_Variadic() {

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

	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	row := rt.db.QueryRow(context.Background(), stmt, 1)

	err := NewScanner(row).Scan(
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
func (rt *rowTest) Test_row_VariadicOfPointers() {
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

	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	row := rt.db.QueryRow(context.Background(), stmt, 1)

	err := NewScanner(row).Scan(
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

func (rt *rowTest) Test_row_InterfaceSliceOfTypes() {
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
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	row := rt.db.QueryRow(context.Background(), stmt, 1)

	err := NewScanner(row).Scan(dst)
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

func (rt *rowTest) Test_scanner_Scan() {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	var ts test.TestStruct
	row := rt.db.QueryRow(context.Background(), stmt, 1)
	err := NewScanner(row).Scan(
		&ts.ID,
		&ts.Int,
		&ts.Float32,
		&ts.String,
		&ts.Time,
		&ts.Bool,
		&ts.Bytes,
		&ts.StringSlice,
		&ts.JSONB)
	if err != nil {
		rt.NoError(err)
	}

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
