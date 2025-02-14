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

func (a *App) GetSchemas(id string) (client.Result, error) {
	return GetSchemas(id)
}

func (a *App) GetTables(id string, schema string) (client.Result, error) {
	return GetTables(id, schema)
}

func (a *App) GetTable(id string, schema string, table string) (client.Result, error) {
	return GetTable(id, schema, table)
}

func (a *App) ExecuteQuery(id string, query string) (client.QueryResult, error) {
	return ExecuteQuery(id, query)
}
