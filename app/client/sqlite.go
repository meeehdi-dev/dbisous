package client

import (
	"database/sql"
	"strings"
)

type SqliteClient struct {
	Db *sql.DB
}

func (c *SqliteClient) GetDatabaseMetadata() (DatabaseMetadata, error) {
	var databaseMetadata DatabaseMetadata

	tables, err := c.getTables()
	if err != nil {
		return databaseMetadata, err
	}

	databaseMetadata.Columns = make(map[string]map[string][]string)
	databaseMetadata.Columns["main"] = make(map[string][]string)
	for _, table := range tables {
		columns, err := c.getColumns(table)
		if err != nil {
			continue
		}
		databaseMetadata.Columns["main"][table] = columns
	}

	return databaseMetadata, nil
}

func (c *SqliteClient) getColumns(table string) ([]string, error) {
	var columns []string

	rows, err := c.Db.Query("SELECT name FROM pragma_table_info(?)", table)
	if err != nil {
		return columns, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName string
		err := rows.Scan(&columnName)
		if err != nil {
			return columns, err
		}

		columns = append(columns, columnName)
	}

	return columns, nil
}

func (c *SqliteClient) getTables() ([]string, error) {
	var tables []string

	rows, err := c.Db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return tables, err
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return tables, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (c *SqliteClient) fetchColumnsMetadata(table string) ([]ColumnMetadata, error) {
	var columnsMetadata []ColumnMetadata

	columns, err := c.Db.Query("SELECT name, type, COALESCE(dflt_value, 'NULL') AS default_value, CASE \"notnull\" WHEN 1 THEN false ELSE true END nullable, pk AS primary_key FROM pragma_table_info(?)", table)
	if err != nil {
		return columnsMetadata, err
	}
	columnsMetadata, err = fetchColumns(columns)
	if err != nil {
		return columnsMetadata, err
	}

	return columnsMetadata, nil
}

func (c *SqliteClient) executeSelectQuery(query string, params QueryParams, args ...interface{}) (QueryResult, error) {
	queryParts := strings.Split(query, " ")
	table := queryParts[0]

	result, err := executeSelectQuery(c.Db, query, params, args...)
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

func (c *SqliteClient) GetDatabaseSchemas(params QueryParams) (QueryResult, error) {
	return c.executeSelectQuery("sqlite_master WHERE type = 'table'", params)
}

func (c *SqliteClient) GetSchemaTables(params QueryParams, schema string) (QueryResult, error) {
	return c.executeSelectQuery("sqlite_master WHERE type='table' AND name = ?", params, schema)
}

func (c *SqliteClient) GetTableRows(params QueryParams, schema string, table string) (QueryResult, error) {
	return c.executeSelectQuery(table, params)
}

func (c *SqliteClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}

func (c *SqliteClient) Execute(query string) error {
	return execute(c.Db, query)
}
