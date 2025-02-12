package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"

	"dbisous/app/client"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var activeConnections = make(map[string]*sql.DB)
var dbClients = make(map[string]client.DatabaseClient)

func CreateDatabase(dbInfo Database) error {
	if dbInfo.ID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		dbInfo.ID = id.String()
	}

	query := `
        INSERT INTO databases (id, name, type, connection_string)
        VALUES (?, ?, ?, ?)
    `
	_, err := metadataDB.Exec(query, dbInfo.ID, dbInfo.Name, dbInfo.Type, dbInfo.ConnectionString)
	return err
}

func GetDatabases() ([]Database, error) {
	query := `SELECT id, name, type, connection_string FROM databases`
	rows, err := metadataDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []Database
	for rows.Next() {
		var dbInfo Database
		err := rows.Scan(&dbInfo.ID, &dbInfo.Name, &dbInfo.Type, &dbInfo.ConnectionString)
		if err != nil {
			return nil, err
		}
		databases = append(databases, dbInfo)
	}
	return databases, nil
}

func UpdateDatabase(dbInfo Database) error {
	query := `
        UPDATE databases
        SET name = ?, type = ?, connection_string = ?
        WHERE id = ?
    `
	_, err := metadataDB.Exec(query, dbInfo.Name, dbInfo.Type, dbInfo.ConnectionString, dbInfo.ID)
	return err
}

func DeleteDatabase(id string) error {
	query := `DELETE FROM databases WHERE id = ?`
	_, err := metadataDB.Exec(query, id)
	return err
}

func ConnectToDatabase(id string) error {
	query := `SELECT type, connection_string FROM databases WHERE id = ?`
	var dbType, connectionString string
	err := metadataDB.QueryRow(query, id).Scan(&dbType, &connectionString)
	if err != nil {
		return err
	}

	var conn *sql.DB
	switch dbType {
	case "sqlite":
		conn, err = sql.Open("sqlite3", connectionString)
		dbClients[id] = &client.SqliteClient{}
	case "mysql":
		conn, err = sql.Open("mysql", connectionString)
		dbClients[id] = &client.MysqlClient{}
	case "postgres":
		conn, err = sql.Open("postgres", connectionString)
		dbClients[id] = &client.PostgresClient{}
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}
	if err != nil {
		return err
	}

	activeConnections[id] = conn
	return nil
}

func DisconnectFromDatabase(id string) error {
	if conn, exists := activeConnections[id]; exists {
		delete(dbClients, id)
		return conn.Close()
	}
	return fmt.Errorf("no active connection for database ID: %s", id)
}

func GetSchemas(id string) (client.QueryResult, error) {
	conn, exists := activeConnections[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no active connection for database ID: %s", id)
	}

	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetSchemas(conn)
}

func GetTables(id string, schema string) (client.QueryResult, error) {
	conn, exists := activeConnections[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no active connection for database ID: %s", id)
	}

	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetTables(conn, schema)
}

func GetTableRows(id string, schema string, table string) (client.QueryResult, error) {
	conn, exists := activeConnections[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no active connection for database ID: %s", id)
	}

	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetTableRows(conn, schema, table)
}

func ExecuteQuery(id string, query string) (client.QueryResult, error) {
	conn, exists := activeConnections[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no active connection for database ID: %s", id)
	}

	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.ExecuteQuery(conn, query)
}

type Database struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	ConnectionString string `json:"connection_string"`
}
