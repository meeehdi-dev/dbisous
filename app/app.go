package app

import (
	"database/sql"
	"log"

	"github.com/adrg/xdg"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

var metadataDB *sql.DB

func (a *App) Startup() {
	dataFilePath, err := xdg.DataFile("DBisous/metadata.db")
	if err != nil {
		log.Fatal(err)
	}

	metadataDB, err = InitMetadataDB(dataFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) Shutdown() {
	CloseMetadataDB()
}

// Dialog is a simple helper function to show an error dialog
func (a *App) ErrorDialog(title, message string) {
	application.Get().Logger.Error(message)
	dialog := application.Get().Dialog.Error()
	dialog.SetTitle(title)
	dialog.SetMessage(message)
	dialog.Show()
}
