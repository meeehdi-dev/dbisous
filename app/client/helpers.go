package client

import (
	"database/sql"
	"fmt"
	"time"
)

func fetchColumns(rows *sql.Rows) ([]ColumnMetadata, error) {
	var columns []ColumnMetadata

	for rows.Next() {
		var column ColumnMetadata
		err := rows.Scan(&column.Name, &column.Type, &column.DefaultValue, &column.Nullable)
		if err != nil {
			return columns, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}

func fetchRows(rows *sql.Rows) (QueryResult, error) {
	columns, err := rows.Columns()
	if err != nil {
		return QueryResult{}, err
	}

	columnsMetadata := []ColumnMetadata{}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return QueryResult{}, err
	}

	for i, col := range columns {
		columnType := columnTypes[i]
		columnMetadata := ColumnMetadata{
			Name:         col,
			Type:         columnType.DatabaseTypeName(),
			DefaultValue: "",
			Nullable:     false,
		}
		columnsMetadata = append(columnsMetadata, columnMetadata)
	}

	var results []Row
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return QueryResult{}, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			value := values[i]
			switch v := value.(type) {
			case nil:
				value = "NULL"
			case []byte:
				value = string(v)
			default:
				value = v
			}
			row[col] = value
		}
		results = append(results, row)
	}

	return QueryResult{
		Rows:    results,
		Columns: columnsMetadata,
		Total:   len(results),
	}, nil
}

func executeQuery(db *sql.DB, query string, args ...interface{}) (QueryResult, error) {
	var result QueryResult

	start := time.Now()
	rows, err := db.Query(query, args...)
	duration := time.Since(start).String()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	result, err = fetchRows(rows)
	if err != nil {
		return result, err
	}

	result.Duration = duration

	return result, nil
}

func executeSelectQuery(db *sql.DB, query string, limit int, offset int, args ...interface{}) (QueryResult, error) {
	result, err := executeQuery(db, fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", query, limit, offset), args...)

	countRow := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", query), args...)
	err = countRow.Scan(&result.Total)
	if err != nil {
		return result, err
	}

	return result, nil
}
