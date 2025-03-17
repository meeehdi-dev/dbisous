package client

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
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

func (c *PostgresClient) executeSelectQuery(query string, params QueryParams) (QueryResult, error) {
	queryParts := strings.Split(query, " ")
	table := queryParts[0]
	tableParts := strings.Split(table, ".")
	schema := "public"
	tableName := tableParts[0]
	if len(tableParts) > 1 {
		schema = strings.ReplaceAll(tableParts[0], "`", "")
		tableName = strings.ReplaceAll(tableParts[1], "`", "")
	}

	result, err := executeSelectQuery(c.Db, query, params)
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
	return c.executeSelectQuery(fmt.Sprintf("information_schema.tables WHERE table_schema = '%s'", schema), params)
}

func (c *PostgresClient) GetTableRows(params QueryParams, schema string, table string) (QueryResult, error) {
	return c.executeSelectQuery(fmt.Sprintf("%s.%s", schema, table), params)
}

func (c *PostgresClient) ExecuteQuery(query string) (QueryResult, error) {
	return executeQuery(c.Db, query)
}

func (c *PostgresClient) Execute(query string) error {
	return execute(c.Db, query)
}

func (c *PostgresClient) Export(options ExportOptions) (string, error) {
	contents := ""

	if options.WrapInTransaction {
		contents += "BEGIN;\n"
	}

	currentTable := ""
	currentTableMetadata := make([]ColumnMetadata, 0)

	tableColumnsMap := make(map[string][]string)

	// NOTE: STEP 1 => Create tables
	// TODO: refactor and use helper function to avoid duplicate code
	if options.DropTable != DoNothing {
		for i, entity := range options.Selected {
			if strings.Count(entity, ".") == 0 {
				// TODO: schema
			}
			if strings.Count(entity, ".") == 1 {
				parts := strings.Split(entity, ".")
				schema := parts[0]
				table := parts[1]
				if table != currentTable {
					var err error
					currentTableMetadata, err = c.fetchColumnsMetadata(schema, table)
					if err != nil {
						return "", err
					}
					currentTable = table
					tableColumnsMap[table] = make([]string, 0)
				}
				switch options.DropTable {
				case DropAndCreate:
					contents += fmt.Sprintf("DROP TABLE %s;\n", table)
					contents += fmt.Sprintf("CREATE TABLE %s (\n", table)
				case Create:
					contents += fmt.Sprintf("CREATE TABLE %s (\n", table)
				case CreateIfNotExists:
					contents += fmt.Sprintf("CREATE IF NOT EXISTS TABLE %s (\n", table)
				}
			}
			if strings.Count(entity, ".") == 2 {
				parts := strings.Split(entity, ".")
				schema := parts[0]
				table := parts[1]
				column := parts[2]
				tableColumnsMap[table] = append(tableColumnsMap[table], column)
				var currentColumn *ColumnMetadata = nil
				for _, col := range currentTableMetadata {
					if col.Name == column {
						currentColumn = &col
						break
					}
				}
				if currentColumn == nil {
					return "", fmt.Errorf("invalid column name: %s", entity)
				}

				nullable := ""
				defaultValue := ""
				primaryKey := ""
				if !currentColumn.Nullable {
					nullable = " NOT NULL"
				}
				if currentColumn.DefaultValue != "NULL" {
					defaultValue = fmt.Sprintf(" DEFAULT %s", currentColumn.DefaultValue)
				}
				if currentColumn.PrimaryKey {
					primaryKey = " PRIMARY KEY"
				}
				contents += fmt.Sprintf("    %s %s%s%s%s", currentColumn.Name, currentColumn.Type, nullable, defaultValue, primaryKey)
				if i+1 < len(options.Selected) {
					next := options.Selected[i+1]
					// NOTE: only part differing from sqlite as we can use the schema here
					if strings.HasPrefix(next, schema+"."+table+".") {
						contents += ","
					} else {
						contents += "\n);"
					}
				} else {
					contents += "\n);"
				}
				contents += "\n"
			}
		}
	}

	// NOTE: STEP 2 => Insert data
	// TODO: use helper function to avoid duplicate code
	if !options.SchemaOnly {
		contents += "\n"
		for table, columns := range tableColumnsMap {
			query := "SELECT "
			for i, column := range columns {
				query += fmt.Sprintf("\"%s\"", column)
				if i+1 < len(columns) {
					query += ", "
				}
			}
			query += fmt.Sprintf(" FROM %s;", table)
			result, err := c.ExecuteQuery(query)
			if err != nil {
				return "", err
			}
			if len(result.Rows) == 0 {
				continue
			}

			contents += fmt.Sprintf("INSERT INTO %s (", table)
			for i, column := range columns {
				contents += fmt.Sprintf("\"%s\"", column)
				if i+1 < len(columns) {
					contents += ", "
				}
			}
			contents += ") VALUES\n"

			for i, row := range result.Rows {
				contents += "    ("
				for j, column := range columns {
					value := row[column]
					switch v := value.(type) {
					case nil:
						contents += "NULL"
					case []byte:
						// TODO: check type directly from column metadata to correctly add quotes
						contents += fmt.Sprintf("'%s'", strings.ReplaceAll(string(v), "'", "''"))
					case bool:
						if v {
							contents += "TRUE"
						} else {
							contents += "FALSE"
						}
					case int:
						contents += fmt.Sprint(v)
					case int64:
						contents += fmt.Sprint(v)
					case float64:
						contents += fmt.Sprint(v)
					case string:
						contents += fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
					case time.Time:
						contents += fmt.Sprintf("'%s'", v.Format(time.RFC3339))
					default:
						return "", fmt.Errorf("invalid value type: %s (%s)", v, reflect.TypeOf(value))
					}
					if j+1 < len(columns) {
						contents += ", "
					} else {
						contents += ")"
					}
				}
				if i+1 < len(result.Rows) {
					contents += ","
				} else {
					contents += ";"
				}
				contents += "\n"
			}
		}

		contents += "\n"
	}

	if options.WrapInTransaction {
		contents += "COMMIT;\n"
	}

	return contents, nil
}

func (c *PostgresClient) Import(contents string) error {
	err := c.Execute(contents)
	if err != nil {
		if strings.Contains(contents, "BEGIN;") || strings.Contains(contents, "BEGIN TRANSACTION;") {
			c.Execute("ROLLBACK;")
		}
		return err
	}
	return nil
}
