package app

import (
	"context"
	"fmt"

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

func (a *App) Execute(id string, query string) error {
	dbClient, exists := dbClients[id]
	if !exists {
		return fmt.Errorf("no database client for database ID: %s", id)
	}

	err := dbClient.Execute(query)
	if err != nil {
		return err
	}

	return nil
}
