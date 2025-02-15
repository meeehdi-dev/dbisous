package client

import (
	"database/sql"
	"fmt"
)

type MysqlClient struct {
	Db *sql.DB
}

func (c *MysqlClient) GetDatabaseSchemas(limit int, offset int) (QueryResult, error) {
	return executeSelectQuery(c.Db, "information_schema.schemata", limit, offset)
}

func (c *MysqlClient) GetDatabaseInfo(limit int, offset int) (QueryResult, error) {
	return executeSelectQuery(c.Db, "information_schema.columns WHERE table_schema = 'information_schema' AND table_name = 'schemata'", limit, offset)
}

func (c *MysqlClient) GetSchemaTables(limit int, offset int, schema string) (QueryResult, error) {
	return executeSelectQuery(c.Db, "information_schema.tables WHERE table_schema = ?", limit, offset, schema)
}

func (c *MysqlClient) GetSchemaInfo(limit int, offset int, schema string) (QueryResult, error) {
	return executeSelectQuery(c.Db, "information_schema.columns WHERE table_schema = 'information_schema' AND table_name = 'tables'", limit, offset)
}

func (c *MysqlClient) GetTableRows(limit int, offset int, schema string, table string) (QueryResult, error) {
	return executeSelectQuery(c.Db, fmt.Sprintf("`%s`.`%s`", schema, table), limit, offset)
}

func (c *MysqlClient) GetTableInfo(limit int, offset int, schema string, table string) (QueryResult, error) {
	return executeSelectQuery(c.Db, "information_schema.columns WHERE table_schema = ? AND table_name = ?", limit, offset, schema, table)
}

func (c *MysqlClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
