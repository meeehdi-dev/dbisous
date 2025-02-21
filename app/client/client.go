package client

type DatabaseClient interface {
	GetDatabaseSchemas(int, int) (QueryResult, error)
	GetSchemaTables(int, int, string) (QueryResult, error)
	GetTableRows(int, int, string, string) (QueryResult, error)
	ExecuteQuery(string, ...interface{}) (QueryResult, error)
}

type ColumnMetadata struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"default_value"`
	Nullable     bool   `json:"nullable"`
}

type Row map[string]interface{}

// NOTE: Rows should be []Row (fixed in wails v3?)
type QueryResult struct {
	Table      string           `json:"table"`
	PrimaryKey string           `json:"primary_key"`
	Rows       interface{}      `json:"rows"`
	Columns    []ColumnMetadata `json:"columns"`
	Total      int              `json:"total"`
	Duration   string           `json:"duration"`
}
