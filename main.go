package main

import (
	"dbisous/app"
	"embed"
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	dbisous := app.NewApp()
	dbisous.Startup()
	defer dbisous.Shutdown()

	startHidden := false
	env := os.Environ()
	for _, e := range env {
		split := strings.Split(e, "=")
		if len(split) > 0 {
			key := split[0]
			if key == "devserver" {
				startHidden = true
				break
			}
		}
	}

	wailsApp := application.New(application.Options{
		Name:        "DBisous",
		Description: "DBisous Database Client",
		Services: []application.Service{
			application.NewService(dbisous),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		// Enum generation happens at build time in v3 via wails.json configuration or bindings generator
	})

	// Register Enums if v3 has runtime enum registration
	// Wails v3 usually uses bindings generator for enums, we'll keep them out of runtime API

	windowOptions := application.WebviewWindowOptions{
		Title:  "DBisous",
		Width:  1024,
		Height: 768,
		URL:    "/",
		Hidden: startHidden,
	}

	wailsApp.Window.NewWithOptions(windowOptions)

	err := wailsApp.Run()
	if err != nil {
		panic(err)
	}
}
