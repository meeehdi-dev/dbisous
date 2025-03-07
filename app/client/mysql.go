package client

import (
	"database/sql"
	"fmt"
	"strings"
)

type MysqlClient struct {
	Db *sql.DB
}

func (c *MysqlClient) GetDatabaseMetadata() (DatabaseMetadata, error) {
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

func (c *MysqlClient) getColumns(schema string, table string) ([]string, error) {
	columns := make([]string, 0)

	rows, err := c.Db.Query("SELECT column_name FROM information_schema.columns WHERE table_schema = ? AND table_name = ?", schema, table)
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

func (c *MysqlClient) getTables(schema string) ([]string, error) {
	tables := make([]string, 0)

	rows, err := c.Db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = ?", schema)
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

func (c *MysqlClient) getSchemas() ([]string, error) {
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

func (c *MysqlClient) fetchColumnsMetadata(schema string, table string) ([]ColumnMetadata, error) {
	var columnsMetadata []ColumnMetadata

	columns, err := c.Db.Query("SELECT c.COLUMN_NAME AS name, c.DATA_TYPE AS type, COALESCE(c.COLUMN_DEFAULT, 'NULL') AS default_value, CASE c.IS_NULLABLE WHEN 'YES' THEN true ELSE false END nullable, COALESCE((SELECT TRUE FROM information_schema.TABLE_CONSTRAINTS tc LEFT JOIN information_schema.KEY_COLUMN_USAGE kcu ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME WHERE tc.TABLE_SCHEMA = ? AND tc.TABLE_NAME = ? AND tc.CONSTRAINT_TYPE = 'PRIMARY KEY' AND kcu.COLUMN_NAME = c.COLUMN_NAME GROUP BY tc.TABLE_SCHEMA, tc.TABLE_NAME, kcu.COLUMN_NAME), FALSE) AS primary_key FROM information_schema.COLUMNS c WHERE c.TABLE_SCHEMA = ? AND c.TABLE_NAME = ?", schema, table, schema, table)
	if err != nil {
		return columnsMetadata, err
	}
	columnsMetadata, err = fetchColumns(columns)
	if err != nil {
		return columnsMetadata, err
	}

	return columnsMetadata, nil
}

func (c *MysqlClient) executeSelectQuery(query string, params QueryParams, args ...any) (QueryResult, error) {
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

func (c *MysqlClient) GetDatabaseSchemas(params QueryParams) (QueryResult, error) {
	return c.executeSelectQuery("information_schema.schemata", params)
}

func (c *MysqlClient) GetSchemaTables(params QueryParams, schema string) (QueryResult, error) {
	return c.executeSelectQuery("information_schema.tables WHERE table_schema = ?", params, schema)
}

func (c *MysqlClient) GetTableRows(params QueryParams, schema string, table string) (QueryResult, error) {
	return c.executeSelectQuery(fmt.Sprintf("`%s`.`%s`", schema, table), params)
}

func (c *MysqlClient) ExecuteQuery(query string, args ...any) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}

func (c *MysqlClient) Execute(query string) error {
	return execute(c.Db, query)
}

func (c *MysqlClient) Export(options ExportOptions) (string, error) {
	contents := ""

	if options.WrapInTransaction {
		contents += "BEGIN;\n"
	}

	currentTable := ""
	currentTableMetadata := make([]ColumnMetadata, 0)

	// NOTE: STEP 1 => Create tables
	// TODO: refactor and use helper function to avoid duplicate code
	if options.DropTable != DoNothing {
		for i, entity := range options.Selected {
			if strings.Count(entity, ".") == 0 {
				// TODO: schema
			}
			parts := strings.Split(entity, ".")
			schema := parts[0]
			table := parts[1]
			if table != currentTable {
				currentTable = table
				var err error
				currentTableMetadata, err = c.fetchColumnsMetadata(schema, table)
				if err != nil {
					return "", err
				}
			}
			switch options.DropTable {
			case DropAndCreate:
				contents += fmt.Sprintf("DROP TABLE %s;\n", entity)
				contents += fmt.Sprintf("CREATE TABLE %s (\n", entity)
			case Create:
				contents += fmt.Sprintf("CREATE TABLE %s (\n", entity)
			case CreateIfNotExists:
				contents += fmt.Sprintf("CREATE IF NOT EXISTS TABLE %s (\n", entity)
			}
			if strings.Count(entity, ".") == 2 {
				parts := strings.Split(entity, ".")
				schema := parts[0]
				table := parts[1]
				column := parts[2]
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

	if options.WrapInTransaction {
		contents += "COMMIT;\n"
	}

	return contents, nil
}
