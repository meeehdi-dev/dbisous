package app

import (
	"context"
	"database/sql"
	"log"

	"github.com/adrg/xdg"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	Ctx context.Context
}

func NewApp() *App {
	return &App{}
}

var metadataDB *sql.DB

func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx

	dataFilePath, err := xdg.DataFile("DBisous/metadata.db")
	if err != nil {
		log.Fatal(err)
	}

	metadataDB, err = InitMetadataDB(dataFilePath)
	if err != nil {
		runtime.MessageDialog(a.Ctx, runtime.MessageDialogOptions{Title: err.Error()})
		log.Fatal(err)
	}
}

func (a *App) Shutdown(ctx context.Context) {
	CloseMetadataDB()
}
