package client

import (
	"database/sql"
	"fmt"
	"strings"
)

type PostgresClient struct {
	Db *sql.DB
}

func (c *PostgresClient) GetDatabaseMetadata() (DatabaseMetadata, error) {
	var databaseMetadata DatabaseMetadata

	schemas, err := c.getSchemas()
	if err != nil {
		return databaseMetadata, err
	}

	databaseMetadata.Columns = make(map[string]map[string][]string)
	for _, schema := range schemas {
		databaseMetadata.Columns[schema] = make(map[string][]string)
		tables, err := c.getTables(schema)
		if err != nil {
			continue
		}
		for _, table := range tables {
			columns, err := c.getColumns(schema, table)
			if err != nil {
				continue
			}
			databaseMetadata.Columns[schema][table] = columns
		}
	}

	return databaseMetadata, nil
}

func (c *PostgresClient) getColumns(schema string, table string) ([]string, error) {
	columns := make([]string, 0)

	rows, err := c.Db.Query("SELECT column_name FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2", schema, table)
	if err != nil {
		return columns, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName string
		err := rows.Scan(&columnName)
		if err != nil {
			return columns, err
		}

		columns = append(columns, columnName)
	}

	return columns, nil
}

func (c *PostgresClient) getTables(schema string) ([]string, error) {
	tables := make([]string, 0)

	rows, err := c.Db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = $1", schema)
	if err != nil {
		return tables, err
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return tables, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (c *PostgresClient) getSchemas() ([]string, error) {
	schemas := make([]string, 0)

	rows, err := c.Db.Query("SELECT schema_name FROM information_schema.schemata")
	if err != nil {
		return schemas, err
	}
	defer rows.Close()

	for rows.Next() {
		var schema string
		err := rows.Scan(&schema)
		if err != nil {
			return schemas, err
		}
		schemas = append(schemas, schema)
	}

	return schemas, nil
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

func (c *PostgresClient) executeSelectQuery(query string, params QueryParams, args ...any) (QueryResult, error) {
	queryParts := strings.Split(query, " ")
	table := queryParts[0]
	tableParts := strings.Split(table, ".")
	schema := "public"
	tableName := tableParts[0]
	if len(tableParts) > 1 {
		schema = strings.ReplaceAll(tableParts[0], "`", "")
		tableName = strings.ReplaceAll(tableParts[1], "`", "")
	}

	result, err := executeSelectQuery(c.Db, query, params, args...)
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

func (c *PostgresClient) GetDatabaseSchemas(params QueryParams) (QueryResult, error) {
	return c.executeSelectQuery("information_schema.schemata", params)
}

func (c *PostgresClient) GetSchemaTables(params QueryParams, schema string) (QueryResult, error) {
	return c.executeSelectQuery("information_schema.tables WHERE table_schema = $1", params, schema)
}

func (c *PostgresClient) GetTableRows(params QueryParams, schema string, table string) (QueryResult, error) {
	return c.executeSelectQuery(fmt.Sprintf("%s.%s", schema, table), params)
}

func (c *PostgresClient) ExecuteQuery(query string, args ...any) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}

func (c *PostgresClient) Execute(query string) error {
	return execute(c.Db, query)
}
