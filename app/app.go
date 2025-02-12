package app

import (
	"context"
	"dbisous/app/client"
	"dbisous/app/database"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	database.InitMetadataDB("metadata.db")
}

func (a *App) Shutdown(ctx context.Context) {
	database.CloseMetadataDB()
}

func (a *App) SelectFile() (string, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return "", err
	}

	return file, nil
}

func (a *App) GetConnections() ([]database.Connection, error) {
	return database.GetConnections()
}

func (a *App) CreateConnection(connection database.Connection) error {
	return database.CreateConnection(connection)
}

func (a *App) UpdateConnection(connection database.Connection) error {
	return database.UpdateConnection(connection)
}

func (a *App) DeleteConnection(id string) error {
	return database.DeleteConnection(id)
}

func (a *App) Connect(id string) error {
	return database.Connect(id)
}

func (a *App) Disconnect(id string) error {
	return database.Disconnect(id)
}

func (a *App) GetSchemas(id string) (client.QueryResult, error) {
	return database.GetSchemas(id)
}

func (a *App) GetTables(id string, schema string) (client.QueryResult, error) {
	return database.GetTables(id, schema)
}

func (a *App) GetTableRows(id string, schema string, table string) (client.QueryResult, error) {
	return database.GetTableRows(id, schema, table)
}

func (a *App) ExecuteQuery(id string, query string) (client.QueryResult, error) {
	return database.ExecuteQuery(id, query)
}
