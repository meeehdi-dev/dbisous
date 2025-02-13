package client

import (
	"database/sql"
	"fmt"
)

type PostgresClient struct {
	Db *sql.DB
}

func (c *PostgresClient) GetSchemas() (QueryResult, error) {
	query := "SELECT * FROM information_schema.schemata"
	return executeQuery(c.Db, query)
}

func (c *PostgresClient) GetTables(schema string) (QueryResult, error) {
	query := "SELECT * FROM information_schema.tables WHERE table_schema = $1"
	return executeQuery(c.Db, query, schema)
}

func (c *PostgresClient) GetTableRows(schema string, table string) (QueryResult, error) {
	query := fmt.Sprintf("SELECT * FROM `%s`.`%s`", schema, table)
	return executeQuery(c.Db, query)
}

func (c *PostgresClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
