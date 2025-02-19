package client

import (
	"database/sql"
	"fmt"
	"strings"
)

type MysqlClient struct {
	Db *sql.DB
}

func (c *MysqlClient) fetchColumnsMetadata(schema string, table string) ([]ColumnMetadata, error) {
	var columnsMetadata []ColumnMetadata

	columns, err := c.Db.Query("SELECT column_name AS name, data_type AS type, COALESCE(column_default, 'NULL') AS default_value, CASE \"is_nullable\" WHEN 'YES' THEN true ELSE false END nullable FROM information_schema.columns WHERE table_schema = ? AND table_name = ?", schema, table)
	if err != nil {
		return columnsMetadata, err
	}
	columnsMetadata, err = fetchColumns(columns)
	if err != nil {
		return columnsMetadata, err
	}

	return columnsMetadata, nil
}

func (c *MysqlClient) executeSelectQuery(query string, limit int, offset int, args ...interface{}) (QueryResult, error) {
	queryParts := strings.Split(query, " ")
	table := queryParts[0]
	tableParts := strings.Split(table, ".")
	schema := "public"
	tableName := tableParts[0]
	if len(tableParts) > 1 {
		schema = tableParts[0]
		tableName = tableParts[1]
	}

	result, err := executeSelectQuery(c.Db, query, limit, offset, args...)
	if err != nil {
		return result, err
	}

	columnsMetadata, err := c.fetchColumnsMetadata(schema, tableName)
	if err != nil {
		return result, err
	}
	result.Columns = columnsMetadata

	return result, err
}

func (c *MysqlClient) GetDatabaseSchemas(limit int, offset int) (QueryResult, error) {
	return c.executeSelectQuery("information_schema.schemata", limit, offset)
}

func (c *MysqlClient) GetSchemaTables(limit int, offset int, schema string) (QueryResult, error) {
	return c.executeSelectQuery("information_schema.tables WHERE table_schema = ?", limit, offset, schema)
}

func (c *MysqlClient) GetTableRows(limit int, offset int, schema string, table string) (QueryResult, error) {
	return c.executeSelectQuery(fmt.Sprintf("%s.%s", schema, table), limit, offset)
}

func (c *MysqlClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}
