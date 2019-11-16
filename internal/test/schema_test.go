package test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func Test_CreateTestDBSchema(t *testing.T) {
	if err := CreateTestDBSchema(); err != nil {
		t.Error(err)
		return
	}
}
func Test_CreateTestRows(t *testing.T) {
	if err := CreateTestRows(); err != nil {
		t.Error(err)
		return
	}
}
func Test_Select(t *testing.T) {
	conn, err := NewConnection()
	if err != nil {
		t.Error(err)
		return
	}

	var ts TestStruct
	stmt := `SELECT 
			"id",
			"int",
			"int_8",
			"int_16",
			"int_32",
			"int_64",
			"uint",
			"uint_8",
			"uint_16",
			"uint_32",
			"uint_64",
			"float_32",
			"float_64",
			"rune",
			"byte",
			"string",
			"bool",
			"time",
			"bytes",
			"string_slice",
			"bool_slice",
			"int_slice",
			"float_slice",
			"json",
			"json_b",
			"map"
	FROM "test" ORDER BY "id" ASC LIMIT 1`
	err = conn.QueryRow(context.Background(), stmt).Scan(
		&ts.ID,
		&ts.Int,
		&ts.Int8,
		&ts.Int16,
		&ts.Int32,
		&ts.Int64,
		&ts.Uint,
		&ts.Uint8,
		&ts.Uint16,
		&ts.Uint32,
		&ts.Uint64,
		&ts.Float32,
		&ts.Float64,
		&ts.Rune,
		&ts.Byte,
		&ts.String,
		&ts.Bool,
		&ts.Time,
		&ts.Bytes,
		&ts.StringSlice,
		&ts.BoolSlice,
		&ts.IntSlice,
		&ts.FloatSlice,
		&ts.JSON,
		&ts.JSONB,
		&ts.Map,
	)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(TestRow1, ts) {
		t.Error(fmt.Errorf("got: %v want: %v", TestRow1, ts))
		return
	}
}

func Test_SelectAll(t *testing.T) {
	conn, err := NewConnection()
	if err != nil {
		t.Error(err)
		return
	}

	var ts TestStruct
	stmt := `SELECT * FROM "test" ORDER BY "id" ASC LIMIT 1`
	err = conn.QueryRow(context.Background(), stmt).Scan(
		&ts.ID,
		&ts.Int,
		&ts.Int8,
		&ts.Int16,
		&ts.Int32,
		&ts.Int64,
		&ts.Uint,
		&ts.Uint8,
		&ts.Uint16,
		&ts.Uint32,
		&ts.Uint64,
		&ts.Float32,
		&ts.Float64,
		&ts.Rune,
		&ts.Byte,
		&ts.String,
		&ts.Bool,
		&ts.Time,
		&ts.Bytes,
		&ts.StringSlice,
		&ts.BoolSlice,
		&ts.IntSlice,
		&ts.FloatSlice,
		&ts.JSON,
		&ts.JSONB,
		&ts.Map,
	)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(TestRow1, ts) {
		t.Error(fmt.Errorf("got: %v want: %v", TestRow1, ts))
		return
	}
}

func Test_Insert(t *testing.T) {
	conn, err := NewConnection()
	if err != nil {
		t.Error(err)
		return
	}

	stmt := `INSERT INTO "test" ("bool_slice") VALUES ($1) RETURNING "id"`
	rows, err := conn.Query(context.Background(), stmt, []bool{true, false})
	if err != nil {
		t.Error(err)
		return
	}

	defer rows.Close()
	var id int
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			t.Error(err)
			return
		}
	}

	var ts TestStruct
	stmt = `SELECT "bool_slice" FROM "test" WHERE id=$1`
	if err = conn.QueryRow(ctxb, stmt, id).Scan(&ts.BoolSlice); err != nil {
		t.Error(err)
		return
	}
	want := []bool{true, false}
	if !reflect.DeepEqual(want, ts.BoolSlice) {
		t.Error(fmt.Errorf("got: %v want: %v", want, ts.BoolSlice))
		return
	}
}
