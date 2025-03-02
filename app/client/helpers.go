package client

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func fetchColumns(rows *sql.Rows) ([]ColumnMetadata, error) {
	var columns []ColumnMetadata

	for rows.Next() {
		var column ColumnMetadata
		err := rows.Scan(&column.Name, &column.Type, &column.DefaultValue, &column.Nullable, &column.PrimaryKey)
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
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return QueryResult{}, err
		}

		row := make(Row)
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

func executeQuery(db *sql.DB, query string, args ...any) (QueryResult, error) {
	var result QueryResult
	result.Rows = make([]Row, 0)
	result.Columns = make([]ColumnMetadata, 0)

	lower := strings.ToLower(query)
	// TODO: refactor?
	isReturning := strings.Contains(lower, "returning")
	isMutate := strings.Contains(lower, "insert") || strings.Contains(lower, "update") || strings.Contains(lower, "delete") || strings.Contains(lower, "upsert") || strings.Contains(lower, "create") || strings.Contains(lower, "alter") || strings.Contains(lower, "truncate") || strings.Contains(lower, "drop")

	if isMutate && !isReturning {
		start := time.Now()
		_, err := db.Exec(query, args...)
		duration := time.Since(start).String()
		if err != nil {
			return result, err
		}

		result.Duration = duration

		return result, nil
	} else {
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
}

func executeSelectQuery(db *sql.DB, query string, params QueryParams, args ...any) (QueryResult, error) {
	execQuery := fmt.Sprintf("SELECT * FROM %s", query)
	if len(params.Order) > 0 {
		execQuery += fmt.Sprintf(" ORDER BY %s %s", params.Order[0].Column, params.Order[0].Direction)
	}
	execQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, params.Offset)
	result, err := executeQuery(db, execQuery, args...)

	countRow := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", query), args...)
	err = countRow.Scan(&result.Total)
	if err != nil {
		return result, err
	}

	return result, nil
}

func execute(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
