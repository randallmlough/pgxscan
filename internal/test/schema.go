package test

import (
	"time"
)

const schema = `
        DROP TABLE IF EXISTS "test";
        CREATE  TABLE "test" (
            "id" SERIAL PRIMARY KEY NOT NULL,
			"int" INT,
			"int_8" SMALLINT,
			"int_16" SMALLINT,
			"int_32" INTEGER,
			"int_64" BIGINT,
			"uint" INT,
			"uint_8" SMALLINT,
			"uint_16" SMALLINT,
			"uint_32" INTEGER,
			"uint_64" BIGINT,
			"float_32" NUMERIC,
			"float_64" NUMERIC,
			"rune" INTEGER,
			"byte" SMALLINT,
			"string" TEXT,
			"bool" BOOL,
			"time" TIMESTAMP,
			"bytes" VARCHAR(45),
			"string_slice" TEXT[],
			"bool_slice" BOOL[],
			"int_slice" INTEGER[],
			"float_slice" NUMERIC[],
			"json" json,
			"json_b" jsonb,
			"map" jsonb
		);
    `

var TestCols = []string{
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
	"map",
}

type (
	TestStruct struct {
		ID uint32 `db:"id"`

		// builtin types
		Int     int       `db:"int"`
		Int8    int8      `db:"int_8"`
		Int16   int16     `db:"int_16"`
		Int32   int32     `db:"int_32"`
		Int64   int64     `db:"int_64"`
		Uint    uint      `db:"uint"`
		Uint8   uint8     `db:"uint_8"`
		Uint16  uint16    `db:"uint_16"`
		Uint32  uint32    `db:"uint_32"`
		Uint64  uint64    `db:"uint_64"`
		Float32 float32   `db:"float_32"`
		Float64 float64   `db:"float_64"`
		Rune    rune      `db:"rune"`
		Byte    byte      `db:"byte"`
		String  string    `db:"string"`
		Bool    bool      `db:"bool"`
		Time    time.Time `db:"time"`
		Bytes   []byte    `db:"bytes"`

		// slices
		StringSlice []string  `db:"string_slice"`
		BoolSlice   []bool    `db:"bool_slice"`
		IntSlice    []int32   `db:"int_slice"`
		FloatSlice  []float32 `db:"float_slice"`

		// json data
		JSON  JSON                   `json:"json" db:"json"`
		JSONB JSON                   `json:"json_b" db:"json_b"`
		Map   map[string]interface{} `json:"map" db:"map"`
	}
	JSON struct {
		Str      string         `json:"str"`
		Int      int            `json:"int"`
		Embedded EmbeddedStruct `json:"embedded"`
		Ignore   string         `json:"-"`
	}
	EmbeddedStruct struct {
		Bool bool `json:"data"`
	}
)

var (
	MockTime = time.Date(2019, 01, 01, 01, 01, 01, 000000, time.UTC)
	TestRow1 = TestStruct{
		ID:          1,
		Int:         1,
		Int8:        121,
		Int16:       32761,
		Int32:       2147483641,
		Int64:       9223372036854775801,
		Uint:        11,
		Uint8:       121,
		Uint16:      32761,
		Uint32:      2147483641,
		Uint64:      9223372036854775801,
		Float32:     1.21,
		Float64:     9715.631,
		Rune:        '😀',
		Byte:        'a',
		String:      "Hello world",
		Bool:        true,
		Time:        MockTime,
		Bytes:       []byte(`first row`),
		StringSlice: []string{"cats", "dogs"},
		BoolSlice:   []bool{true, false, false, true},
		IntSlice:    []int32{1, 2, 3, 4, 5},
		FloatSlice:  []float32{1.21, 2.21, 3.21, 4.21},
		JSON: JSON{
			Str:      "I'm json",
			Int:      1,
			Embedded: EmbeddedStruct{},
		},
		JSONB: JSON{
			Str: "I'm json b",
			Int: 1,
			Embedded: EmbeddedStruct{
				Bool: true,
			},
		},
		Map: map[string]interface{}{
			"key": "value",
		},
	}

	TestRow2 = TestStruct{
		ID:          2,
		Int:         2,
		Int8:        122,
		Int16:       32762,
		Int32:       2147483642,
		Int64:       9223372036854775802,
		Uint:        12,
		Uint8:       252,
		Uint16:      32762,
		Uint32:      2147483642,
		Uint64:      9223372036854775802,
		Float32:     1.22,
		Float64:     9715.632,
		Rune:        '😂',
		Byte:        'b',
		String:      "foo bar",
		Bool:        true,
		Time:        MockTime,
		Bytes:       []byte(`second row`),
		StringSlice: []string{"john doe", "jane smith"},
		BoolSlice:   []bool{false, false, false, true},
		IntSlice:    []int32{6, 7, 8, 9, 10},
		FloatSlice:  []float32{5.21, 6.21, 7.21, 8.21},
		JSONB: JSON{
			Str: "Hi",
			Int: 2,
			Embedded: EmbeddedStruct{
				Bool: true,
			},
		},
		Map: map[string]interface{}{
			"marco": "polo",
		},
	}
)

func CreateTestDBSchema() error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}
	if _, err := conn.Exec(ctxb, schema); err != nil {
		return err
	}
	return nil
}

func CreateTestRows() error {
	conn, err := NewConnection()
	if err != nil {
		return err
	}

	stmt := `INSERT INTO "test" (
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
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25);`
	_, err = conn.Exec(ctxb, stmt,
		TestRow1.Int,
		TestRow1.Int8,
		TestRow1.Int16,
		TestRow1.Int32,
		TestRow1.Int64,
		TestRow1.Uint,
		TestRow1.Uint8,
		TestRow1.Uint16,
		TestRow1.Uint32,
		TestRow1.Uint64,
		TestRow1.Float32,
		TestRow1.Float64,
		TestRow1.Rune,
		TestRow1.Byte,
		TestRow1.String,
		TestRow1.Bool,
		TestRow1.Time,
		TestRow1.Bytes,
		TestRow1.StringSlice,
		TestRow1.BoolSlice,
		TestRow1.IntSlice,
		TestRow1.FloatSlice,
		TestRow1.JSON,
		TestRow1.JSONB,
		TestRow1.Map,
	)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctxb, stmt,
		TestRow2.Int,
		TestRow2.Int8,
		TestRow2.Int16,
		TestRow2.Int32,
		TestRow2.Int64,
		TestRow2.Uint,
		TestRow2.Uint8,
		TestRow2.Uint16,
		TestRow2.Uint32,
		TestRow2.Uint64,
		TestRow2.Float32,
		TestRow2.Float64,
		TestRow2.Rune,
		TestRow2.Byte,
		TestRow2.String,
		TestRow2.Bool,
		TestRow2.Time,
		TestRow2.Bytes,
		TestRow2.StringSlice,
		TestRow2.BoolSlice,
		TestRow2.IntSlice,
		TestRow2.FloatSlice,
		TestRow2.JSON,
		TestRow2.JSONB,
		TestRow2.Map,
	)
	if err != nil {
		return err
	}
	return nil
}
