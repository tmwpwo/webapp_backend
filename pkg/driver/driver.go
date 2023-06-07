package driver

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB represents a database connection pool.
type DB struct {
	SQL *sql.DB
}

// Constants for configuring the database connection pool.
const (
	maxOpenDBConn = 10
	maxIdleDBConn = 5
	maxDBLifeTime = 5 * time.Minute
)

// ConnectSQL establishes a connection to a PostgreSQL database and returns a new DB object.
func ConnectSQL(dsn string) (*DB, error) {
	sqlDB, err := NewDatabase(dsn)
	if err != nil {
		return nil, errors.New("failed to connect to database: " + err.Error())
	}

	// Configure the connection pool.
	sqlDB.SetMaxOpenConns(maxOpenDBConn)
	sqlDB.SetMaxIdleConns(maxIdleDBConn)
	sqlDB.SetConnMaxLifetime(maxDBLifeTime)

	// Test the connection to the database.
	if err = TestDB(sqlDB); err != nil {
		return nil, errors.New("failed to connect to database: " + err.Error())
	}

	return &DB{SQL: sqlDB}, nil
}

// testDB pings the database to test the connection.
func TestDB(sqlDB *sql.DB) error {
	return sqlDB.Ping()
}

// newDatabase creates a new instance of the PostgreSQL database.
func NewDatabase(dsn string) (*sql.DB, error) {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return sqlDB, nil
}
