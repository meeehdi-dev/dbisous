package client

import (
	"database/sql"
)

type SqliteClient struct {
	Db *sql.DB
}

func (c *SqliteClient) GetSchemas() (Result, error) {
	var result Result

	data, err := executeSelectQuery(c.Db, "sqlite_master WHERE type = 'table'")
	if err != nil {
		return result, err
	}
	result.Data = data

	info, err := executeSelectQuery(c.Db, "pragma_table_info('sqlite_master')")
	if err != nil {
		return result, err
	}
	result.Info = info

	return result, nil
}

func (c *SqliteClient) GetTables(schema string) (Result, error) {
	var result Result

	data, err := executeSelectQuery(c.Db, "sqlite_master WHERE type='table' AND name = ?", schema)
	if err != nil {
		return result, err
	}
	result.Data = data

	info, err := executeSelectQuery(c.Db, "pragma_table_info(?)", schema)
	if err != nil {
		return result, err
	}
	result.Info = info

	return result, nil
}

func (c *SqliteClient) GetTable(schema string, table string) (Result, error) {
	var result Result

	data, err := executeSelectQuery(c.Db, table)
	if err != nil {
		return result, err
	}
	result.Data = data

	info, err := executeSelectQuery(c.Db, "pragma_table_info(?)", table)
	if err != nil {
		return result, err
	}
	result.Info = info

	return result, nil
}

func (c *SqliteClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
