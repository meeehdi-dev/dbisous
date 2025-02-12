package client

import (
	"database/sql"
)

type DatabaseClient interface {
	GetSchemas(db *sql.DB) (QueryResult, error)
	GetTables(db *sql.DB, schema string) (QueryResult, error)
	GetTableRows(db *sql.DB, schema string, table string) (QueryResult, error)
	ExecuteQuery(db *sql.DB, query string, args ...interface{}) (QueryResult, error)
}

type Row map[string]interface{}

type QueryResult struct {
	// NOTE: Rows should be []Row (fixed in wails v3?)
	Rows          interface{}      `json:"rows"`
	Columns       []ColumnMetadata `json:"columns"`
	SqlDuration   string           `json:"sql_duration"`
	TotalDuration string           `json:"total_duration"`
}

type ColumnMetadata struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}
