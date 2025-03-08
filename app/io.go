package app

import (
	"dbisous/app/client"
	"fmt"
)

func Export(id string, options client.ExportOptions) (string, error) {
	dbClient, exists := dbClients[id]
	if !exists {
		return "", fmt.Errorf("no database client for database ID: %s", id)
	}

	contents, err := dbClient.Export(options)
	if err != nil {
		return "", err
	}

	return contents, nil
}

func Import(id string, contents string) error {
	dbClient, exists := dbClients[id]
	if !exists {
		return fmt.Errorf("no database client for database ID: %s", id)
	}

	err := dbClient.Import(contents)
	if err != nil {
		return err
	}

	return nil
}
