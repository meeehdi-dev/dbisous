package client

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func fetchColumns(rows *sql.Rows) ([]ColumnMetadata, error) {
	columns := make([]ColumnMetadata, 0)

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

	results := make([]Row, 0)
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

func executeQuery(db *sql.DB, query string) (QueryResult, error) {
	result := QueryResult{Query: query}
	result.Rows = make([]Row, 0)
	result.Columns = make([]ColumnMetadata, 0)

	lower := strings.ToLower(query)
	// FIXME: rewrite to avoid issues with columns named 'created_at' etc !!
	isReturning := strings.Contains(lower, "returning ")
	isMutate := strings.Contains(lower, "insert ") || strings.Contains(lower, "update ") || strings.Contains(lower, "delete ") || strings.Contains(lower, "upsert ") || strings.Contains(lower, "create ") || strings.Contains(lower, "alter ") || strings.Contains(lower, "truncate ") || strings.Contains(lower, "drop ")

	if isMutate && !isReturning {
		start := time.Now()
		_, err := db.Exec(query)
		duration := time.Since(start).String()
		if err != nil {
			return result, err
		}

		result.Duration = duration

		return result, nil
	} else {
		start := time.Now()
		rows, err := db.Query(query)
		duration := time.Since(start).String()
		if err != nil {
			return result, err
		}
		defer rows.Close()

		result, err = fetchRows(rows)
		if err != nil {
			return result, err
		}

		result.Query = query + ";"
		result.Duration = duration

		return result, nil
	}
}

func executeSelectQuery(db *sql.DB, query string, params QueryParams) (QueryResult, error) {
	execQuery := fmt.Sprintf("SELECT * FROM %s", query)

	if len(params.Filter) > 0 {
		if !strings.Contains(execQuery, "WHERE") {
			execQuery += " WHERE"
		} else {
			execQuery += " AND"
		}
		execQuery += " ("

		filters := make([]string, 0)
		for _, filter := range params.Filter {
			operator := "="
			if strings.HasPrefix(filter.Value, "'") {
				operator = "LIKE"
			}
			filters = append(filters, filter.Column+" "+operator+" "+filter.Value)
		}
		execQuery += strings.Join(filters, " AND ")

		execQuery += ")"
	}

	if len(params.Order) > 0 {
		execQuery += " ORDER BY "
		orders := make([]string, 0)
		for _, order := range params.Order {
			orders = append(orders, order.Column+" "+string(order.Direction))
		}
		execQuery += strings.Join(orders, ", ")
	}

	execQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, params.Offset)

	result, err := executeQuery(db, execQuery)

	countRow := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", query))
	err = countRow.Scan(&result.Total)
	if err != nil {
		return result, err
	}

	return result, nil
}

func execute(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	if err != nil {
		if strings.Contains(query, "BEGIN;") || strings.Contains(query, "BEGIN TRANSACTION;") {
			db.Exec("ROLLBACK;")
		}
		return err
	}

	return nil
}
