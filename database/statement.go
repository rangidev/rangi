package database

var (
	// Create table
	statementCreateTable                 = "CREATE TABLE IF NOT EXISTS %s (%s);"
	statementCreateReferenceTableSqlite3 = "CREATE TABLE IF NOT EXISTS %s (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, %s_id INTEGER NOT NULL, %s_id INTEGER NOT NULL, FOREIGN KEY(%s_id) REFERENCES %s(id), FOREIGN KEY(%s_id) REFERENCES %s(id));"
	// TODO: Fix this for Postgres
	statementCreateReferenceTablePostgres = "CREATE TABLE IF NOT EXISTS %s (id BIGSERIAL NOT NULL PRIMARY KEY, %s_id BIGINT NOT NULL, %s_id BIGINT NOT NULL, FOREIGN KEY(%s_id) REFERENCES %s(id), FOREIGN KEY(%s_id) REFERENCES %s(id));"
)
