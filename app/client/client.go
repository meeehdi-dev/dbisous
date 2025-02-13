package client

type DatabaseClient interface {
	GetSchemas() (QueryResult, error)
	GetDatabaseInfo() (QueryResult, error)
	GetTables(schema string) (QueryResult, error)
	GetSchemaInfo(schema string) (QueryResult, error)
	GetTableRows(schema string, table string) (QueryResult, error)
	GetTableInfo(schema string, table string) (QueryResult, error)
	ExecuteQuery(query string, args ...interface{}) (QueryResult, error)
}

type ColumnMetadata struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

type Row map[string]interface{}

type QueryResult struct {
	// NOTE: Rows should be []Row (fixed in wails v3?)
	Rows          interface{}      `json:"rows"`
	Columns       []ColumnMetadata `json:"columns"`
	SqlDuration   string           `json:"sql_duration"`
	TotalDuration string           `json:"total_duration"`
}
