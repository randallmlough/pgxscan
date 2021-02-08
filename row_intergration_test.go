// +build integration

package pgxscan_test

import (
	"context"
	"testing"
	"time"

	"github.com/randallmlough/pgxscan/testdata"
	"github.com/stretchr/testify/require"

	"github.com/randallmlough/pgxscan"
)

func Test_row_Variadic(t *testing.T) {

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

	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	row := newTestDB(t).QueryRow(context.Background(), stmt, 1)

	err := pgxscan.NewScanner(row).Scan(
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
func Test_row_VariadicOfPointers(t *testing.T) {
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

	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	row := newTestDB(t).QueryRow(context.Background(), stmt, 1)

	err := pgxscan.NewScanner(row).Scan(
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

func Test_row_InterfaceSliceOfTypes(t *testing.T) {
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
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	row := newTestDB(t).QueryRow(context.Background(), stmt, 1)

	err := pgxscan.NewScanner(row).Scan(dst)
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

func Test_scanner_Scan(t *testing.T) {
	stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
	var ts testdata.TestStruct
	row := newTestDB(t).QueryRow(context.Background(), stmt, 1)
	err := pgxscan.NewScanner(row).Scan(
		&ts.ID,
		&ts.Int,
		&ts.Float32,
		&ts.String,
		&ts.Time,
		&ts.Bool,
		&ts.Bytes,
		&ts.StringSlice,
		&ts.JSONB)
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
