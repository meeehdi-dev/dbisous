package app

import (
	"bufio"
	"dbisous/app/client"
	"fmt"
	"io"
	"os"
)

func exportDatabase(file string, id string, options client.ExportOptions) (string, error) {
	// TODO: savefiledialog before exporting to make use of buffered writes and avoid memory issues
	dbClient, exists := dbClients[id]
	if !exists {
		return "", fmt.Errorf("no database client for database ID: %s", id)
	}

	contents, err := dbClient.Export(options)
	if err != nil {
		return "", err
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

func importDatabase(file string, id string) (string, error) {
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
