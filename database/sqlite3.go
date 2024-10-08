package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewSqlite3Instance(databasePath string) (*DB, error) {
	db, err := sqlx.Connect("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", databasePath))
	if err != nil {
		return nil, err
	}
	return &DB{db: db, dbType: DatabaseTypeSqlite3}, nil
}
