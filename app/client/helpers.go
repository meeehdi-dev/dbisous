package client

import (
	"database/sql"
	"time"
)

func fetchRows(rows *sql.Rows) (QueryResult, error) {
	columns, err := rows.Columns()
	if err != nil {
		return QueryResult{}, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return QueryResult{}, err
	}

	var columnMetadata []ColumnMetadata
	for i, col := range columns {
		nullable, _ := columnTypes[i].Nullable()
		columnMetadata = append(columnMetadata, ColumnMetadata{
			Name:     col,
			Type:     columnTypes[i].DatabaseTypeName(),
			Nullable: nullable,
		})
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
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return QueryResult{
		Rows:    results,
		Columns: columnMetadata,
	}, nil
}

func executeQuery(db *sql.DB, query string, args ...interface{}) (QueryResult, error) {
	start := time.Now()
	rows, err := db.Query(query, args...)
	sqlDuration := time.Since(start).String()
	if err != nil {
		return QueryResult{}, err
	}
	defer rows.Close()

	result, err := fetchRows(rows)
	if err != nil {
		return QueryResult{}, err
	}

	result.SqlDuration = sqlDuration
	result.TotalDuration = time.Since(start).String()
	return result, nil
}
