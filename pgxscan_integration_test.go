//go:build integration
// +build integration

package pgxscan_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

const defaultDbURI = "postgres://root:pass@pgxscan_postgres:5432/pgxscan?sslmode=disable"

var (
	ctxb   = context.Background()
	testDB *pgxpool.Pool
)

func TestMain(m *testing.M) {

	dbURI := os.Getenv("PG_URI")
	if dbURI == "" {
		dbURI = defaultDbURI
	}

	var err error
	testDB, err = pgxpool.Connect(ctxb, dbURI)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func newTestDB(tb testing.TB) *pgxpool.Conn {
	conn, err := testDB.Acquire(ctxb)
	if err != nil {
		tb.Fatalf("error aquiring a new connection: %v", err)
	}

	tb.Cleanup(func() {
		conn.Release()
	})

	return conn
}
