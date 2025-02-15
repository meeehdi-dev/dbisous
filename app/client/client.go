package client

type DatabaseClient interface {
	GetDatabaseSchemas(int, int) (QueryResult, error)
	GetDatabaseInfo(int, int) (QueryResult, error)
	GetSchemaTables(int, int, string) (QueryResult, error)
	GetSchemaInfo(int, int, string) (QueryResult, error)
	GetTableRows(int, int, string, string) (QueryResult, error)
	GetTableInfo(int, int, string, string) (QueryResult, error)
	ExecuteQuery(string, ...interface{}) (QueryResult, error)
}

type ColumnMetadata struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

type Row map[string]interface{}

type QueryResult struct {
	// NOTE: Rows should be []Row (fixed in wails v3?)
	Rows     interface{}      `json:"rows"`
	Columns  []ColumnMetadata `json:"columns"`
	Total    int              `json:"total"`
	Duration string           `json:"duration"`
}
