package app

import (
	"bufio"
	"context"
	"dbisous/app/client"
	"fmt"
	"io"
	"os"

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

func (a *App) Connect(id string) (client.DatabaseMetadata, error) {
	return Connect(id)
}

func (a *App) Disconnect(id string) error {
	return Disconnect(id)
}

func (a *App) GetDatabaseSchemas(id string, params client.QueryParams) (client.QueryResult, error) {
	return GetDatabaseSchemas(id, params)
}

func (a *App) GetSchemaTables(id string, params client.QueryParams, schema string) (client.QueryResult, error) {
	return GetSchemaTables(id, params, schema)
}

func (a *App) GetTableRows(id string, params client.QueryParams, schema string, table string) (client.QueryResult, error) {
	return GetTableRows(id, params, schema, table)
}

func (a *App) ExecuteQuery(id string, query string) (client.QueryResult, error) {
	return ExecuteQuery(id, query)
}

func (a *App) GetPastQueries() ([]PastQuery, error) {
	return GetPastQueries()
}

func (a *App) DeletePastQuery(id string) error {
	return DeletePastQuery(id)
}

func (a *App) Execute(id string, query string) error {
	return Execute(id, query)
}

func (a *App) Export(id string, options client.ExportOptions) (string, error) {
	// TODO: savefiledialog before exporting to make use of buffered writes and avoid memory issues
	contents, err := Export(id, options)
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

	err = Import(id, contents)
	if err != nil {
		return "", err
	}

	return file, nil
}
