# PGX Scan
A simple scanning library to extend [PGX's](https://github.com/jackc/pgx) awesome capabilities.

`pgxscan` supports scanning to structs (including things like join tables and JSON columns), slices of structs, scanning from interface slices and variadic arguments.

## How to use
### For the Row interface
Scanning to a row (ie. by calling `QueryRow()`) which returns the row interface only exposes the scan method. Currently `pgxscan` or for that matter, `pgx`, doesnt have a way to expose the columns returned from the row query. Because of this `pgxscan` can only scan to pre defined types.
To scan to a struct by passing in a struct, use the rows interface  (ie. `Query()`).


#### Scan to standard types
```go
var (
    ID     uint32
    Int    int
    Float  float32
    String string
    Time   time.Time
    Bool   bool
    Bytes  []byte
    Slice  []string
)

conn, _ := pgx.ConnectConfig(ctxb, "postgres://postgres:@localhost:5432/pgxscan?sslmode=disable")
stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice" FROM "test" WHERE id = $1`
row := conn.QueryRow(context.Background(), stmt, 1)

_ := pgxscan.NewScanner(row).Scan(
    &ID,
    &Int,
    &Float,
    &String,
    &Time,
    &Bool,
    &Bytes,
    &Slice,
)
```

#### Scan to standard types from `[]interface{}`
```go
var (
    ID     uint32
    Int    int
    Float  float32
    String string
    Time   time.Time
    Bool   bool
    Bytes  []byte
    Slice  []string
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
}
conn, _ := pgx.ConnectConfig(ctxb, "postgres://postgres:@localhost:5432/pgxscan?sslmode=disable")
stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice" FROM "test" WHERE id = $1`
row := conn.QueryRow(context.Background(), stmt, 1)

_ := pgxscan.NewScanner(row).Scan(dst)
```

#### Scan to struct fields
```go
type TestStruct struct {
    ID uint32 `db:"id"`

    // builtin types
    Int     int       `db:"int"`
    Float32 float32   `db:"float_32"`
    String  string    `db:"string"`
    Bool    bool      `db:"bool"`
    Time    time.Time `db:"time"`
    Bytes   []byte    `db:"bytes"`
    StringSlice []string  `db:"string_slice"`
    JSONB JSON  `json:"json_b" db:"json_b"`
}
type JSON struct {
    Str      string         `json:"str"`
    Int      int            `json:"int"`
    Embedded EmbeddedStruct `json:"embedded"`
    Ignore   string         `json:"-"`
}
type EmbeddedStruct struct {
    Bool bool `json:"data"`
}

// scanning to pre defined struct fields
stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
var dst TestStruct
row := conn.QueryRow(context.Background(), stmt, 1)
_ := pgxscan.NewScanner(row).Scan(
    &dst.ID,
    &dst.Int,
    &dst.Float32,
    &dst.String,
    &dst.Time,
    &dst.Bool,
    &dst.Bytes,
    &dst.StringSlice,
    &dst.JSONB,
)
``` 

### For the Rows interface
The Rows interface exposes more data like, returned column names, which allows us to scan into a without pre defining the values first. 
But all the previous examples will also work for rows too.

#### Scan to a struct
```go
stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" WHERE id = $1`
rows, _ := conn.Query(context.Background(), stmt, 1)

var dst TestStruct
// pgxscan will take care of closing the rows and calling next()
if err := pgxscan.NewScanner(rows).Scan(&dst); err != nil {
    return err
}
``` 

#### Scan to slice of structs
```go
stmt := `SELECT "id", "int", "float_32", "string", "time", "bool", "bytes", "string_slice", "json_b" FROM "test" ORDER BY "id" ASC LIMIT 2`
rows, _ := conn.Query(context.Background(), stmt)

var dst []TestStruct
if err := pgxscan.NewScanner(rows).Scan(&dst); err != nil {
    return err
}
```

#### Scan to struct with join table
There's two ways to handle join tables. Either use the tag `scan:"table"` or `scan:"follow"`. `scan table` will annotate the struct to something like `"table_one.column"` this is particularly useful if joining tables that have column name conflicts. However, you will have to alias the sql column to match.
`scan follow` wont annotate and instead go into the struct and add the field names to the map. If you know you won't have column name conflicts this will work fine and no aliasing is required.

Example with aliasing
```go
stmt := `
WITH usr AS (
	SELECT
		"id", "name", "email"
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
	addresses
`
	// aliased address, line_1, and city notation.
	rows, _ := conn.Query(context.Background(), stmt, 1)

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
			Address Address `scan:"table"` // table annotates the struct
		}
	)
	var user User
	if err := NewScanner(rows).Scan(&user); err != nil {
		return err
	}
```

Example without aliasing
```go
stmt := `
WITH usr AS (
	SELECT
		"id", "name", "email"
	FROM
		"users"
	WHERE
		"id" = $1
),
addresses AS (
	SELECT
		"line_1",
		"city"
	FROM
		"address", usr
	WHERE
		"user_id" = usr."id"
)
SELECT
	usr.*, addresses.*
FROM
	usr,
	addresses
`
	// note that the "id" column is not being selected, which removes the naming conflict therefore no aliasing is necessary
	rows, _ := conn.Query(context.Background(), stmt, 1)

	type (
		Address struct {
			Line1 string `db:"line_1"`
			City  string
		}
		User struct {
			ID      uint32
			Name    string
			Email   string
			Address Address `scan:"follow"` // follow inspects the struct and adds the fields without being annotated. 
		}
	)
	var user User
	if err := NewScanner(rows).Scan(&user); err != nil {
		return err
	}
```

Checkout the many other tests for examples on scanning to different data types