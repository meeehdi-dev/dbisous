package client

import (
	"database/sql"
	"fmt"
	"strings"
)

type PostgresClient struct {
	Db *sql.DB
}

func (c *PostgresClient) fetchColumnsMetadata(schema string, table string) ([]ColumnMetadata, error) {
	var columnsMetadata []ColumnMetadata

	columns, err := c.Db.Query("SELECT c.column_name AS name, c.data_type AS type, COALESCE(c.column_default, 'NULL') AS default_value, CASE c.is_nullable WHEN 'YES' THEN true ELSE false END nullable, COALESCE((SELECT TRUE FROM information_schema.table_constraints tc LEFT JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name WHERE tc.table_schema = $1 AND tc.table_name = $2 AND tc.constraint_type = 'PRIMARY KEY' AND kcu.COLUMN_NAME = c.COLUMN_NAME GROUP BY tc.TABLE_SCHEMA, tc.TABLE_NAME, kcu.COLUMN_NAME), FALSE) AS primary_key FROM information_schema.columns c WHERE c.table_schema = $1 AND c.table_name = $2", schema, table)
	if err != nil {
		return columnsMetadata, err
	}
	columnsMetadata, err = fetchColumns(columns)
	if err != nil {
		return columnsMetadata, err
	}

	return columnsMetadata, nil
}

func (c *PostgresClient) executeSelectQuery(query string, limit int, offset int, args ...interface{}) (QueryResult, error) {
	queryParts := strings.Split(query, " ")
	table := queryParts[0]
	tableParts := strings.Split(table, ".")
	schema := "public"
	tableName := tableParts[0]
	if len(tableParts) > 1 {
		schema = strings.ReplaceAll(tableParts[0], "`", "")
		tableName = strings.ReplaceAll(tableParts[1], "`", "")
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

func (c *PostgresClient) GetDatabaseSchemas(limit int, offset int) (QueryResult, error) {
	return c.executeSelectQuery("information_schema.schemata", limit, offset)
}

func (c *PostgresClient) GetSchemaTables(limit int, offset int, schema string) (QueryResult, error) {
	return c.executeSelectQuery("information_schema.tables WHERE table_schema = $1", limit, offset, schema)
}

func (c *PostgresClient) GetTableRows(limit int, offset int, schema string, table string) (QueryResult, error) {
	return c.executeSelectQuery(fmt.Sprintf("%s.%s", schema, table), limit, offset)
}

func (c *PostgresClient) ExecuteQuery(query string, args ...interface{}) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}

func (c *PostgresClient) Execute(query string) error {
	return execute(c.Db, query)
}
