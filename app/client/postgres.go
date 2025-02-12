package client

import (
	"database/sql"
	"fmt"
)

type PostgresClient struct{}

func (c *PostgresClient) GetSchemas(db *sql.DB) (QueryResult, error) {
	query := "SELECT * FROM information_schema.schemata"
	return executeQuery(db, query)
}

func (c *PostgresClient) GetTables(db *sql.DB, schema string) (QueryResult, error) {
	query := "SELECT * FROM information_schema.tables WHERE table_schema = $1"
	return executeQuery(db, query, schema)
}

func (c *PostgresClient) GetTableRows(db *sql.DB, schema string, table string) (QueryResult, error) {
	query := fmt.Sprintf("SELECT * FROM `%s`.`%s`", schema,table)
	return executeQuery(db, query)
}

func (c *PostgresClient) ExecuteQuery(db *sql.DB, query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(db, query, args...)
}
