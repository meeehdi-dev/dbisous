package client

type DatabaseClient interface {
	GetDatabaseMetadata() (DatabaseMetadata, error)
	GetDatabaseSchemas(int, int) (QueryResult, error)
	GetSchemaTables(int, int, string) (QueryResult, error)
	GetTableRows(int, int, string, string) (QueryResult, error)
	ExecuteQuery(string, ...interface{}) (QueryResult, error)
	Execute(string) error
}

type ColumnMetadata struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"default_value"`
	Nullable     bool   `json:"nullable"`
	PrimaryKey   bool   `json:"primary_key"`
}

type Row map[string]interface{}

// NOTE: Rows should be []Row (fixed in wails v3?)
type QueryResult struct {
	Rows     interface{}      `json:"rows"`
	Columns  []ColumnMetadata `json:"columns"`
	Total    int              `json:"total"`
	Duration string           `json:"duration"`
}

type DatabaseMetadata struct {
	Columns map[string]map[string][]string `json:"columns"`
}
