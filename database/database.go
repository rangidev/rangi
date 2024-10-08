package database

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	db     *sqlx.DB
	dbType DatabaseType
}

type DatabaseType string

const (
	DatabaseTypeSqlite3  = DatabaseType("sqlite3")
	DatabaseTypePostgres = DatabaseType("postgres")
)

var (
	ErrorUnknownDatabaseType = errors.New("unknown database type")
)
