package test

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"os"
)

const (
	defaultDbURI = "postgres://postgres:@localhost:5435/pgxscan?sslmode=disable"
)

var ctxb = context.Background()

func NewConnection() (*pgx.Conn, error) {

	dbURI := os.Getenv("PG_URI")
	if dbURI == "" {
		dbURI = defaultDbURI
	}

	uri, err := pq.ParseURL(dbURI)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse db uri")
	}

	cfg, err := pgx.ParseConfig(uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create db config")
	}

	conn, err := pgx.ConnectConfig(ctxb, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connection to database")
	}

	if err := conn.Ping(ctxb); err != nil {
		return nil, errors.Wrap(err, "failed to ping db")
	}
	return conn, nil
}
