package client

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
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

	rows, err := c.Db.Query("SELECT name FROM sqlite_master WHERE type LIKE 'table'")
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

func (c *SqliteClient) fetchColumnsMetadata(table string, columns []string) ([]ColumnMetadata, error) {
	columnsMetadata := make([]ColumnMetadata, 0)

	cColumns := ""
	if len(columns) > 0 {
		cColumns = fmt.Sprintf(" WHERE name IN (%s)", "'"+strings.Join(columns, "', '")+"'")
	}

	queryColumns, err := c.Db.Query(fmt.Sprintf("SELECT name, type, COALESCE(dflt_value, 'NULL') AS default_value, CASE \"notnull\" WHEN 1 THEN false ELSE true END nullable, pk AS primary_key FROM pragma_table_info('%s')%s", table, cColumns))
	if err != nil {
		return columnsMetadata, err
	}
	columnsMetadata, err = fetchColumns(queryColumns)
	if err != nil {
		return columnsMetadata, err
	}

	return columnsMetadata, nil
}

func (c *SqliteClient) executeSelectQuery(query string, params QueryParams) (QueryResult, error) {
	queryParts := strings.Split(query, " ")
	table := queryParts[0]

	result, err := executeSelectQuery(c.Db, query, params)
	if err != nil {
		return result, err
	}

	columns := []string{}
	aliases := make(map[string]string)
	for _, col := range params.Columns {
		tokens := strings.Split(col, " AS ")
		columns = append(columns, tokens[0])
		if len(tokens) > 1 {
			aliases[tokens[0]] = tokens[1]
		}
	}

	columnsMetadata, err := c.fetchColumnsMetadata(table, columns)
	if err != nil {
		return result, err
	}
	result.Columns = columnsMetadata

	// handle aliases
	for i, col := range result.Columns {
		result.Columns[i].OriginalName = col.Name
		if aliases[col.Name] != "" {
			result.Columns[i].Name = aliases[col.Name]
		}
	}

	return result, err
}

func (c *SqliteClient) GetConnectionDatabases(params QueryParams) (QueryResult, error) {
	rows := make([]Row, 0)
	row := make(Row)
	row["name"] = "main"
	rows = append(rows, row)
	columns := make([]ColumnMetadata, 0)
	columns = append(columns, ColumnMetadata{Name: "name"})
	return QueryResult{
		Rows:    rows,
		Columns: columns,
	}, nil
}

func (c *SqliteClient) GetDatabaseSchemas(params QueryParams) (QueryResult, error) {
	params.Columns = []string{"name"}
	return c.executeSelectQuery("sqlite_master WHERE type LIKE 'table'", params)
}

func (c *SqliteClient) GetSchemaTables(params QueryParams, schema string) (QueryResult, error) {
	params.Columns = []string{"name"}
	return c.executeSelectQuery(fmt.Sprintf("sqlite_master WHERE type LIKE 'table' AND name = '%s'", schema), params)
}

func (c *SqliteClient) GetTableRows(params QueryParams, schema string, table string) (QueryResult, error) {
	return c.executeSelectQuery(table, params)
}

func (c *SqliteClient) ExecuteQuery(query string) (QueryResult, error) {
	return executeQuery(c.Db, query)
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

	tableColumnsMap := make(map[string][]string)

	// NOTE: STEP 1 => Create tables
	// TODO: refactor and use helper function to avoid duplicate code
	if options.SchemaOnly || options.DropTable != DoNothing {
		for i, entity := range options.Selected {
			if strings.Count(entity, ".") == 0 {
				// NOTE: ignore default schema
			}
			if strings.Count(entity, ".") == 1 {
				parts := strings.Split(entity, ".")
				// schema := parts[0]
				table := parts[1]
				if table != currentTable {
					var err error
					currentTableMetadata, err = c.fetchColumnsMetadata(table, []string{})
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
				// schema := parts[0]
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
					return "", fmt.Errorf("invalid column name: %s", column)
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
					// NOTE: only part differing from others as we can use the schema here
					if strings.HasPrefix(next, "main."+table+".") {
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

func (c *SqliteClient) Import(contents string) error {
	err := c.Execute(contents)
	if err != nil {
		if strings.Contains(contents, "BEGIN;") || strings.Contains(contents, "BEGIN TRANSACTION;") {
			c.Execute("ROLLBACK;")
		}
		return err
	}
	return nil
}
