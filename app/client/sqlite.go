package client

import (
	"database/sql"
	"strings"
)

type SqliteClient struct {
	Db *sql.DB
}

func (c *SqliteClient) fetchColumnsMetadata(table string) ([]ColumnMetadata, error) {
	var columnsMetadata []ColumnMetadata

	columns, err := c.Db.Query("SELECT name, type, COALESCE(dflt_value, 'NULL') AS default_value, CASE \"notnull\" WHEN 1 THEN false ELSE true END nullable FROM pragma_table_info(?)", table)
	if err != nil {
		return columnsMetadata, err
	}
	columnsMetadata, err = fetchColumns(columns)
	if err != nil {
		return columnsMetadata, err
	}

	return columnsMetadata, nil
}

func (c *SqliteClient) executeSelectQuery(query string, limit int, offset int, args ...interface{}) (QueryResult, error) {
	queryParts := strings.Split(query, " ")
	table := queryParts[0]

	result, err := executeSelectQuery(c.Db, query, limit, offset, args...)
	if err != nil {
		return result, err
	}

	columnsMetadata, err := c.fetchColumnsMetadata(table)
	if err != nil {
		return result, err
	}
	result.Columns = columnsMetadata

	return result, err
}

func (c *SqliteClient) GetDatabaseSchemas(limit int, offset int) (QueryResult, error) {
	return c.executeSelectQuery("sqlite_master WHERE type = 'table'", limit, offset)
}

func (c *SqliteClient) GetDatabaseInfo(limit int, offset int) (QueryResult, error) {
	return c.executeSelectQuery("pragma_table_info('sqlite_master')", limit, offset)
}

func (c *SqliteClient) GetSchemaTables(limit int, offset int, schema string) (QueryResult, error) {
	return c.executeSelectQuery("sqlite_master WHERE type='table' AND name = ?", limit, offset, schema)
}

func (c *SqliteClient) GetSchemaInfo(limit int, offset int, schema string) (QueryResult, error) {
	return c.executeSelectQuery("pragma_table_info(?)", limit, offset, schema)
}

func (c *SqliteClient) GetTableRows(limit int, offset int, schema string, table string) (QueryResult, error) {
	return c.executeSelectQuery(table, limit, offset)
}

func (c *SqliteClient) GetTableInfo(limit int, offset int, schema string, table string) (QueryResult, error) {
	return c.executeSelectQuery("pragma_table_info(?)", limit, offset, table)
}

func (c *SqliteClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
