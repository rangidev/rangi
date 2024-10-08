package database

import "github.com/rangidev/rangi/blueprint"

type SQLType string

const (
	// SQLite3
	SQLTypeSqlite3ID      = "INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT"
	SQLTypeSqlite3Text    = "TEXT"
	SQLTypeSqlite3Boolean = "INTEGER"
	SQLTypeSqlite3Integer = "INTEGER"
	SQLTypeSqlite3Blob    = "BLOB"
	// PostgreSQL
	SQLTypePostgresID      = "BIGSERIAL NOT NULL PRIMARY KEY"
	SQLTypePostgresUUID    = "CHARACTER(36)"
	SQLTypePostgresText    = "TEXT"
	SQLTypePostgresBoolean = "BOOLEAN"
	SQLTypePostgresBigInt  = "BIGINT"
	SQLTypePostgresJsonb   = "JSONB"
	// All
	SQLTypeReference = "REFERENCE"
)

func (db *DB) castToSQLType(typ blueprint.Type) (SQLType, bool) {
	switch db.dbType {
	case DatabaseTypeSqlite3:
		switch typ {
		case blueprint.TypeID:
			return SQLTypeSqlite3ID, true
		case blueprint.TypeUUID:
			return SQLTypeSqlite3Text, true
		case blueprint.TypeString:
			return SQLTypeSqlite3Text, true
		case blueprint.TypeBoolean:
			return SQLTypeSqlite3Boolean, true
		case blueprint.TypeInt:
			return SQLTypeSqlite3Integer, true
		case blueprint.TypeArray:
			return SQLTypeSqlite3Blob, true
		case blueprint.TypeObject:
			return SQLTypeSqlite3Blob, true
		case blueprint.TypeReference:
			return SQLTypeReference, true
		default:
			return "", false
		}
	case DatabaseTypePostgres:
		switch typ {
		case blueprint.TypeID:
			return SQLTypePostgresID, true
		case blueprint.TypeUUID:
			return SQLTypePostgresUUID, true
		case blueprint.TypeString:
			return SQLTypePostgresText, true
		case blueprint.TypeBoolean:
			return SQLTypePostgresBoolean, true
		case blueprint.TypeInt:
			return SQLTypePostgresBigInt, true
		case blueprint.TypeArray:
			return SQLTypePostgresJsonb, true
		case blueprint.TypeObject:
			return SQLTypePostgresJsonb, true
		case blueprint.TypeReference:
			return SQLTypeReference, true
		default:
			return "", false
		}
	default:
		return "", false
	}
}
