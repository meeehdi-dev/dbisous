package app

import (
	"database/sql"
	"dbisous/app/client"
	"fmt"
)

func (a *App) GetConnectionDatabases(id string, params client.QueryParams) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetConnectionDatabases(params)
}

func (a *App) UseDatabase(id string, connectionString string) error {
	var dbType string
	err := metadataDB.QueryRow(`SELECT type FROM connection WHERE id = ?`, id).Scan(&dbType)
	if err != nil {
		return err
	}

	var db *sql.DB
	switch dbType {
	case string(MySQL):
		db, err = sql.Open("mysql", connectionString)
		dbClients[id] = &client.MysqlClient{Db: db}
	case string(PostgreSQL):
		db, err = sql.Open("postgres", connectionString)
		dbClients[id] = &client.PostgresClient{Db: db}
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}
	if err != nil {
		return err
	}

	return nil
}

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
