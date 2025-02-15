package client

import (
	"database/sql"
	"fmt"
)

type MysqlClient struct {
	Db *sql.DB
}

func (c *MysqlClient) GetSchemas() (Result, error) {
	var result Result

	data, err := executeSelectQuery(c.Db, "information_schema.schemata")
	if err != nil {
		return result, err
	}
	result.Data = data

	info, err := executeSelectQuery(c.Db, "information_schema.columns WHERE table_schema = 'information_schema' AND table_name = 'schemata'")
	if err != nil {
		return result, err
	}
	result.Info = info

	return result, nil
}

func (c *MysqlClient) GetTables(schema string) (Result, error) {
	var result Result

	data, err := executeSelectQuery(c.Db, "information_schema.tables WHERE table_schema = ?", schema)
	if err != nil {
		return result, err
	}
	result.Data = data

	info, err := executeSelectQuery(c.Db, "information_schema.columns WHERE table_schema = 'information_schema' AND table_name = 'tables'")
	if err != nil {
		return result, err
	}
	result.Info = info

	return result, nil
}

func (c *MysqlClient) GetTable(schema string, table string) (Result, error) {
	var result Result

	data, err := executeSelectQuery(c.Db, fmt.Sprintf("`%s`.`%s`", schema, table))
	if err != nil {
		return result, err
	}
	result.Data = data

	info, err := executeSelectQuery(c.Db, "information_schema.columns WHERE table_schema = ? AND table_name = ?", schema, table)
	if err != nil {
		return result, err
	}
	result.Info = info

	return result, nil
}

func (c *MysqlClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
