package client

import (
	"database/sql"
	"fmt"
)

type SqliteClient struct{}

func (c *SqliteClient) GetSchemas(db *sql.DB) (QueryResult, error) {
	query := "SELECT * FROM sqlite_master WHERE type='table'"
	return executeQuery(db, query)
}

func (c *SqliteClient) GetTables(db *sql.DB, schema string) (QueryResult, error) {
	query := "SELECT * FROM sqlite_master WHERE type='table' AND name = ?"
	return executeQuery(db, query, schema)
}

func (c *SqliteClient) GetTableRows(db *sql.DB, schema string, table string) (QueryResult, error) {
	query := fmt.Sprintf("SELECT * FROM `%s`", table)
	return executeQuery(db, query)
}

func (c *SqliteClient) ExecuteQuery(db *sql.DB, query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(db, query, args...)
}
