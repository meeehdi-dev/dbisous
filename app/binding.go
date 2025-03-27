package app

import (
	"dbisous/app/client"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetConnections() ([]Connection, error) {
	return getConnections(metadataDB)
}

func (a *App) CreateConnection(connection Connection) error {
	return createConnection(metadataDB, connection)
}

func (a *App) UpdateConnection(connection Connection) error {
	return updateConnection(metadataDB, connection)
}

func (a *App) DeleteConnection(id string) error {
	return deleteConnection(metadataDB, id)
}

func (a *App) Connect(id string) (client.DatabaseMetadata, error) {
	return connect(activeConnections, metadataDB, id)
}

func (a *App) Disconnect(id string) error {
	return disconnect(activeConnections, id)
}

func (a *App) TestConnection(dbType ConnectionType, connectionString string) error {
	return testConnection(dbType, connectionString)
}

func (a *App) GetPastQueries() ([]PastQuery, error) {
	return getPastQueries(metadataDB)
}

func (a *App) DeletePastQuery(id string) error {
	return deletePastQuery(metadataDB, id)
}

func (a *App) ExportDatabase(id string, options client.ExportOptions) (string, error) {
	file, err := runtime.SaveFileDialog(a.Ctx, runtime.SaveDialogOptions{})
	if err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("No file selected")
	}
	return exportDatabase(file, id, options)
}

func (a *App) ImportDatabase(id string) (string, error) {
	// TODO: buffered read and do it step by step to avoid memory overload
	file, err := runtime.SaveFileDialog(a.Ctx, runtime.SaveDialogOptions{})
	if err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("No file selected")
	}

	return importDatabase(file, id)
}

func (a *App) SelectFile() (string, error) {
	file, err := runtime.OpenFileDialog(a.Ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return "", err
	}

	return file, nil
}

func (a *App) GetConnectionDatabases(id string, params client.QueryParams) (client.QueryResult, error) {
	return getConnectionDatabases(id, params)
}

func (a *App) UseDatabase(id string, connectionString string) error {
	return useDatabase(id, connectionString)
}

func (a *App) GetDatabaseSchemas(id string, params client.QueryParams) (client.QueryResult, error) {
	return getDatabaseSchemas(id, params)
}

func (a *App) GetSchemaTables(id string, params client.QueryParams, schema string) (client.QueryResult, error) {
	return getSchemaTables(id, params, schema)
}

func (a *App) GetTableRows(id string, params client.QueryParams, schema string, table string) (client.QueryResult, error) {
	return getTableRows(id, params, schema, table)
}

func (a *App) ExecuteQuery(id string, query string) (client.QueryResult, error) {
	return executeQuery(id, query)
}

func (a *App) Execute(id string, query string) error {
	return execute(id, query)
}
