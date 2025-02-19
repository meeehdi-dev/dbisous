package main

import (
	"dbisous/app"
	"embed"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	dbisous := app.NewApp()

	startHidden := false
	env := os.Environ()
	for _, e := range env {
		split := strings.Split(e, "=")
		key := split[0]
		if key == "devserver" {
			startHidden = true
			break
		}
	}

	err := wails.Run(&options.App{
		Title:     "DBisous",
		MinWidth:  1024,
		MinHeight: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  dbisous.Startup,
		OnShutdown: dbisous.Shutdown,
		Bind: []interface{}{
			dbisous,
		},
    EnumBind: []interface{}{
      app.AllConnectionTypes,
    },
		StartHidden: startHidden,
	})

	if err != nil {
		runtime.MessageDialog(dbisous.Ctx, runtime.MessageDialogOptions{Title: err.Error()})
	}
}
