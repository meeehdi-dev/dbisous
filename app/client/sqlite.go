package client

import (
	"database/sql"
	"fmt"
)

type SqliteClient struct {
	Db *sql.DB
}

func (c *SqliteClient) GetSchemas() (QueryResult, error) {
	query := "SELECT * FROM sqlite_master WHERE type='table'"
	return executeQuery(c.Db, query)
}

func (c *SqliteClient) GetDatabaseInfo() (QueryResult, error) {
	query := "SELECT * FROM pragma_table_info('sqlite_master')"
	return executeQuery(c.Db, query)
}

func (c *SqliteClient) GetTables(schema string) (QueryResult, error) {
	query := "SELECT * FROM sqlite_master WHERE type='table' AND name = ?"
	return executeQuery(c.Db, query, schema)
}

func (c *SqliteClient) GetSchemaInfo(schema string) (QueryResult, error) {
	query := "SELECT * FROM pragma_table_info(?)"
	return executeQuery(c.Db, query, schema)
}

func (c *SqliteClient) GetTableRows(schema string, table string) (QueryResult, error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	return executeQuery(c.Db, query)
}

func (c *SqliteClient) GetTableInfo(schema string, table string) (QueryResult, error) {
	query := "SELECT * FROM pragma_table_info(?)"
	return executeQuery(c.Db, query, table)
}

func (c *SqliteClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
