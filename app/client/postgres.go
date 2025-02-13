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

func (c *PostgresClient) GetDatabaseInfo() (QueryResult, error) {
	query := "SELECT * FROM information_schema.columns WHERE table_schema = 'information_schema' AND table_name = 'schemata'"
	return executeQuery(c.Db, query)
}

func (c *PostgresClient) GetTables(schema string) (QueryResult, error) {
	query := "SELECT * FROM information_schema.tables WHERE table_schema = $1"
	return executeQuery(c.Db, query, schema)
}

func (c *PostgresClient) GetSchemaInfo(schema string) (QueryResult, error) {
	query := "SELECT * FROM information_schema.columns WHERE table_schema = 'information_schema' AND table_name = 'tables'"
	return executeQuery(c.Db, query)
}

func (c *PostgresClient) GetTableRows(schema string, table string) (QueryResult, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s", schema, table)
	return executeQuery(c.Db, query)
}

func (c *PostgresClient) GetTableInfo(schema string, table string) (QueryResult, error) {
	query := "SELECT * FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2"
	return executeQuery(c.Db, query, schema, table)
}

func (c *PostgresClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
