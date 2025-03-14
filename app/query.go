package app

import (
	"dbisous/app/client"
	"fmt"
)

func (a *App) GetDatabaseSchemas(id string, params client.QueryParams) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetDatabaseSchemas(params)
}

func (a *App) GetSchemaTables(id string, params client.QueryParams, schema string) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]

	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetSchemaTables(params, schema)
}

func (a *App) GetTableRows(id string, params client.QueryParams, schema string, table string) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetTableRows(params, schema, table)
}

func (a *App) ExecuteQuery(id string, query string) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	result, err := dbClient.ExecuteQuery(query)
	if err != nil {
		return result, err
	}

	err = InsertPastQuery(query)
	if err != nil {
		return result, err
	}

	return result, nil
}
