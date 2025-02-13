package client

import (
	"database/sql"
	"fmt"
)

type MysqlClient struct {
	Db *sql.DB
}

func (c *MysqlClient) GetSchemas() (QueryResult, error) {
	query := "SELECT * FROM information_schema.schemata"
	return executeQuery(c.Db, query)
}

func (c *MysqlClient) GetTables(schema string) (QueryResult, error) {
	query := "SELECT * FROM information_schema.tables WHERE table_schema = ?"
	return executeQuery(c.Db, query, schema)
}

func (c *MysqlClient) GetTableRows(schema string, table string) (QueryResult, error) {
	query := fmt.Sprintf("SELECT * FROM `%s`.`%s`", schema, table)
	return executeQuery(c.Db, query)
}

func (c *MysqlClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
