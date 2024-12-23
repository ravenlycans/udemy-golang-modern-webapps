package driver

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// DB holds the database connection pool.
type DB struct {
	SQL *pgxpool.Pool
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxDbLifetime = 5 * time.Minute
const maxDbIdletime = 10 * time.Second

// ConnectSQL creates a db pool for postgres.
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// testDB tries to ping the database.
func testDB(d *pgxpool.Pool) error {
	err := d.Ping(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// NewDatabase creates a new database for the application.
func NewDatabase(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConnLifetime = maxDbLifetime
	config.MaxConns = maxOpenDbConn
	config.MaxConnIdleTime = maxDbIdletime

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}
