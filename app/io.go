package app

import (
	"bufio"
	"dbisous/app/client"
	"fmt"
	"io"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) Export(id string, options client.ExportOptions) (string, error) {
	// TODO: savefiledialog before exporting to make use of buffered writes and avoid memory issues
	dbClient, exists := dbClients[id]
	if !exists {
		return "", fmt.Errorf("no database client for database ID: %s", id)
	}

	contents, err := dbClient.Export(options)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}

	file, err := runtime.SaveFileDialog(a.Ctx, runtime.SaveDialogOptions{})
	if err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("No file selected")
	}

	f, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(contents)
	if err != nil {
		return "", err
	}
	w.Flush()

	return file, nil
}

func (a *App) Import(id string) (string, error) {
	// TODO: buffered read and do it step by step to avoid memory overload
	file, err := runtime.SaveFileDialog(a.Ctx, runtime.SaveDialogOptions{})
	if err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("No file selected")
	}

	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	contents := ""
	for {
		line, err := reader.ReadString('\n') // FIXME: replace with ";" ig?
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		contents += line
	}

	dbClient, exists := dbClients[id]
	if !exists {
		return "", fmt.Errorf("no database client for database ID: %s", id)
	}

	err = dbClient.Import(contents)
	if err != nil {
		return "", err
	}

	return file, nil
}
