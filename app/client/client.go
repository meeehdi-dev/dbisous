package client

type DatabaseClient interface {
	GetSchemas() (Result, error)
	GetTables(schema string) (Result, error)
	GetTable(schema string, table string) (Result, error)
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
	Rows     interface{}      `json:"rows"`
	Columns  []ColumnMetadata `json:"columns"`
	Total    int              `json:"total"`
	Duration string           `json:"duration"`
}

type Result struct {
	Data QueryResult `json:"data"`
	Info QueryResult `json:"info"`
}
