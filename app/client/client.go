package client

type DatabaseClient interface {
	GetDatabaseMetadata() (DatabaseMetadata, error)
	GetDatabaseSchemas(QueryParams) (QueryResult, error)
	GetSchemaTables(QueryParams, string) (QueryResult, error)
	GetTableRows(QueryParams, string, string) (QueryResult, error)
	ExecuteQuery(string, ...any) (QueryResult, error)
	Execute(string) error
}

type ColumnMetadata struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"default_value"`
	Nullable     bool   `json:"nullable"`
	PrimaryKey   bool   `json:"primary_key"`
}

type Row map[string]any

type QueryResult struct {
	Rows     []Row            `json:"rows"`
	Columns  []ColumnMetadata `json:"columns"`
	Total    int              `json:"total"`
	Duration string           `json:"duration"`
}

type DatabaseMetadata struct {
	Columns map[string]map[string][]string `json:"columns"`
}

type OrderDirection string

const (
	Ascending  OrderDirection = "ASC"
	Descending OrderDirection = "DESC"
)

var OrderDirections = []struct {
	Value  OrderDirection
	TSName string
}{
	{Ascending, "Ascending"},
	{Descending, "Descending"},
}

type QueryOrder struct {
	Column    string         `json:"column"`
	Direction OrderDirection `json:"direction"`
}

type QueryFilter struct {
	Column string `json:"column"`
	Value  string `json:"value"`
}

type QueryParams struct {
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
	Order  []QueryOrder  `json:"order"`
	Filter []QueryFilter `json:"filter"`
}
