package client

import (
	"database/sql"
)

type SqliteClient struct {
	Db *sql.DB
}

func (c *SqliteClient) GetDatabaseSchemas(limit int, offset int) (QueryResult, error) {
	return executeSelectQuery(c.Db, "sqlite_master WHERE type = 'table'", limit, offset)
}

func (c *SqliteClient) GetDatabaseInfo(limit int, offset int) (QueryResult, error) {
	return executeSelectQuery(c.Db, "pragma_table_info('sqlite_master')", limit, offset)
}

func (c *SqliteClient) GetSchemaTables(limit int, offset int, schema string) (QueryResult, error) {
	return executeSelectQuery(c.Db, "sqlite_master WHERE type='table' AND name = ?", limit, offset, schema)
}

func (c *SqliteClient) GetSchemaInfo(limit int, offset int, schema string) (QueryResult, error) {
	return executeSelectQuery(c.Db, "pragma_table_info(?)", limit, offset, schema)
}

func (c *SqliteClient) GetTableRows(limit int, offset int, schema string, table string) (QueryResult, error) {
	return executeSelectQuery(c.Db, table, limit, offset)
}

func (c *SqliteClient) GetTableInfo(limit int, offset int, schema string, table string) (QueryResult, error) {
	return executeSelectQuery(c.Db, "pragma_table_info(?)", limit, offset, table)
}

func (c *SqliteClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
