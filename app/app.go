package app

import (
	"context"
	"dbisous/app/client"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	Ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx
	err := InitMetadataDB()
	if err != nil {
		runtime.MessageDialog(a.Ctx, runtime.MessageDialogOptions{Title: err.Error()})
	}
}

func (a *App) Shutdown(ctx context.Context) {
	CloseMetadataDB()
}

func (a *App) SelectFile() (string, error) {
	file, err := runtime.OpenFileDialog(a.Ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return "", err
	}

	return file, nil
}

func (a *App) GetConnections() ([]Connection, error) {
	return GetConnections()
}

func (a *App) CreateConnection(connection Connection) error {
	return CreateConnection(connection)
}

func (a *App) UpdateConnection(connection Connection) error {
	return UpdateConnection(connection)
}

func (a *App) DeleteConnection(id string) error {
	return DeleteConnection(id)
}

func (a *App) Connect(id string) error {
	return Connect(id)
}

func (a *App) Disconnect(id string) error {
	return Disconnect(id)
}

func (a *App) GetDatabaseSchemas(id string, page int, itemsPerPage int) (client.QueryResult, error) {
	limit, offset := itemsPerPage, (page-1)*itemsPerPage
	return GetDatabaseSchemas(id, limit, offset)
}

func (a *App) GetDatabaseInfo(id string, page int, itemsPerPage int) (client.QueryResult, error) {
	limit, offset := itemsPerPage, (page-1)*itemsPerPage
	return GetDatabaseInfo(id, limit, offset)
}

func (a *App) GetSchemaTables(id string, page int, itemsPerPage int, schema string) (client.QueryResult, error) {
	limit, offset := itemsPerPage, (page-1)*itemsPerPage
	return GetSchemaTables(id, limit, offset, schema)
}

func (a *App) GetSchemaInfo(id string, page int, itemsPerPage int, schema string) (client.QueryResult, error) {
	limit, offset := itemsPerPage, (page-1)*itemsPerPage
	return GetSchemaInfo(id, limit, offset, schema)
}

func (a *App) GetTableRows(id string, page int, itemsPerPage int, schema string, table string) (client.QueryResult, error) {
	limit, offset := itemsPerPage, (page-1)*itemsPerPage
	return GetTableRows(id, limit, offset, schema, table)
}

func (a *App) GetTableInfo(id string, page int, itemsPerPage int, schema string, table string) (client.QueryResult, error) {
	limit, offset := itemsPerPage, (page-1)*itemsPerPage
	return GetTableInfo(id, limit, offset, schema, table)
}

func (a *App) ExecuteQuery(id string, query string) (client.QueryResult, error) {
	return ExecuteQuery(id, query)
}
