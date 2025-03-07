package client

import (
	"database/sql"
	"fmt"
	"strings"
)

type SqliteClient struct {
	Db *sql.DB
}

func (c *SqliteClient) GetDatabaseMetadata() (DatabaseMetadata, error) {
	var databaseMetadata DatabaseMetadata

	tables, err := c.getTables()
	if err != nil {
		return databaseMetadata, err
	}

	databaseMetadata.Columns = make(map[string]map[string][]string)
	databaseMetadata.Columns["main"] = make(map[string][]string)
	for _, table := range tables {
		columns, err := c.getColumns(table)
		if err != nil {
			continue
		}
		databaseMetadata.Columns["main"][table] = columns
	}

	return databaseMetadata, nil
}

func (c *SqliteClient) getColumns(table string) ([]string, error) {
	columns := make([]string, 0)

	rows, err := c.Db.Query("SELECT name FROM pragma_table_info(?)", table)
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

func (c *SqliteClient) getTables() ([]string, error) {
	tables := make([]string, 0)

	rows, err := c.Db.Query("SELECT name FROM sqlite_master WHERE type='table'")
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

func (c *SqliteClient) fetchColumnsMetadata(table string) ([]ColumnMetadata, error) {
	columnsMetadata := make([]ColumnMetadata, 0)

	columns, err := c.Db.Query("SELECT name, type, COALESCE(dflt_value, 'NULL') AS default_value, CASE \"notnull\" WHEN 1 THEN false ELSE true END nullable, pk AS primary_key FROM pragma_table_info(?)", table)
	if err != nil {
		return columnsMetadata, err
	}
	columnsMetadata, err = fetchColumns(columns)
	if err != nil {
		return columnsMetadata, err
	}

	return columnsMetadata, nil
}

func (c *SqliteClient) executeSelectQuery(query string, params QueryParams, args ...any) (QueryResult, error) {
	queryParts := strings.Split(query, " ")
	table := queryParts[0]

	result, err := executeSelectQuery(c.Db, query, params, args...)
	if err != nil {
		return result, err
	}

	columnsMetadata, err := c.fetchColumnsMetadata(table)
	if err != nil {
		return result, err
	}
	result.Columns = columnsMetadata

	return result, err
}

func (c *SqliteClient) GetDatabaseSchemas(params QueryParams) (QueryResult, error) {
	return c.executeSelectQuery("sqlite_master WHERE type = 'table'", params)
}

func (c *SqliteClient) GetSchemaTables(params QueryParams, schema string) (QueryResult, error) {
	return c.executeSelectQuery("sqlite_master WHERE type='table' AND name = ?", params, schema)
}

func (c *SqliteClient) GetTableRows(params QueryParams, schema string, table string) (QueryResult, error) {
	return c.executeSelectQuery(table, params)
}

func (c *SqliteClient) ExecuteQuery(query string, args ...any) (QueryResult, error) {
	return executeQuery(c.Db, query, args...)
}

func (c *SqliteClient) Execute(query string) error {
	return execute(c.Db, query)
}

func (c *SqliteClient) Export(options ExportOptions) (string, error) {
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
			if entity == "main" {
				continue
			}
			entity, found := strings.CutPrefix(entity, "main.")
			if !found {
				return "", fmt.Errorf("invalid entity name: %s", entity)
			}
			if !strings.Contains(entity, ".") {
				if entity != currentTable {
					currentTable = entity
					var err error
					currentTableMetadata, err = c.fetchColumnsMetadata(currentTable)
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
			}
			if strings.Contains(entity, ".") {
				parts := strings.Split(entity, ".")
				table := parts[0]
				column := parts[1]
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
					next, found = strings.CutPrefix(next, "main.")
					// FIXME: is there a way to avoid any default schema handling?
					if strings.HasPrefix(next, table+".") {
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
