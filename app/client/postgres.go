package client

import (
	"database/sql"
	"fmt"
)

type PostgresClient struct {
	Db *sql.DB
}

func (c *PostgresClient) GetDatabaseSchemas(limit int, offset int) (QueryResult, error) {
	return executeSelectQuery(c.Db, "information_schema.schemata", limit, offset)
}

func (c *PostgresClient) GetDatabaseInfo(limit int, offset int) (QueryResult, error) {
	return executeSelectQuery(c.Db, "information_schema.columns WHERE table_schema = 'information_schema' AND table_name = 'schemata'", limit, offset)
}

func (c *PostgresClient) GetSchemaTables(limit int, offset int, schema string) (QueryResult, error) {
	return executeSelectQuery(c.Db, "information_schema.tables WHERE table_schema = $1", limit, offset, schema)
}

func (c *PostgresClient) GetSchemaInfo(limit int, offset int, schema string) (QueryResult, error) {
	return executeSelectQuery(c.Db, "information_schema.columns WHERE table_schema = 'information_schema' AND table_name = 'tables'", limit, offset)
}

func (c *PostgresClient) GetTableRows(limit int, offset int, schema string, table string) (QueryResult, error) {
	return executeSelectQuery(c.Db, fmt.Sprintf("%s.%s", schema, table), limit, offset)
}

func (c *PostgresClient) GetTableInfo(limit int, offset int, schema string, table string) (QueryResult, error) {
	return executeSelectQuery(c.Db, "information_schema.columns WHERE table_schema = $1 AND table_name = $2", limit, offset, schema, table)
}

func (c *PostgresClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
