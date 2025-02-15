package app

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

func GetConnections() ([]Connection, error) {
	rows, err := metadataDB.Query(`SELECT id, created_at, updated_at, name, type, connection_string FROM connection`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []Connection
	for rows.Next() {
		var connection Connection
		err := rows.Scan(&connection.ID, &connection.CreatedAt, &connection.UpdatedAt, &connection.Name, &connection.Type, &connection.ConnectionString)
		if err != nil {
			return nil, err
		}
		connections = append(connections, connection)
	}

	return connections, nil
}

func CreateConnection(connection Connection) error {
	if connection.ID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		connection.ID = id.String()
	}

	_, err := metadataDB.Exec(`INSERT INTO connection (id, name, type, connection_string)
  VALUES (?, ?, ?, ?)`, connection.ID, connection.Name, connection.Type, connection.ConnectionString)

	return err
}

func UpdateConnection(connection Connection) error {
	_, err := metadataDB.Exec(`UPDATE connection
  SET name = ?, type = ?, connection_string = ?, updated_at = CURRENT_TIMESTAMP
  WHERE id = ?`, connection.Name, connection.Type, connection.ConnectionString, connection.ID)
	return err
}

func DeleteConnection(id string) error {
	_, err := metadataDB.Exec(`DELETE FROM connection WHERE id = ?`, id)
	return err
}

func Connect(id string) error {
	var dbType, connectionString string
	err := metadataDB.QueryRow(`SELECT type, connection_string FROM connection WHERE id = ?`, id).Scan(&dbType, &connectionString)
	if err != nil {
		return err
	}

	var db *sql.DB
	switch dbType {
	case "sqlite":
		db, err = sql.Open("sqlite3", connectionString)
		dbClients[id] = &client.SqliteClient{Db: db}
	case "mysql":
		db, err = sql.Open("mysql", connectionString)
		dbClients[id] = &client.MysqlClient{Db: db}
	case "postgresql":
		db, err = sql.Open("postgres", connectionString)
		dbClients[id] = &client.PostgresClient{Db: db}
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return err
	}

	activeConnections[id] = db
	return nil
}

func Disconnect(id string) error {
	conn, exists := activeConnections[id]
	if !exists {
		return fmt.Errorf("no active connection for database ID: %s", id)
	}

	delete(dbClients, id)
	return conn.Close()
}

func GetDatabaseSchemas(id string, limit int, offset int) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetDatabaseSchemas(limit, offset)
}

func GetDatabaseInfo(id string, limit int, offset int) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetDatabaseInfo(limit, offset)
}

func GetSchemaTables(id string, limit int, offset int, schema string) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]

	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetSchemaTables(limit, offset, schema)
}

func GetSchemaInfo(id string, limit int, offset int, schema string) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]

	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetSchemaInfo(limit, offset, schema)
}

func GetTableRows(id string, limit int, offset int, schema string, table string) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetTableRows(limit, offset, schema, table)
}

func GetTableInfo(id string, limit int, offset int, schema string, table string) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.GetTableInfo(limit, offset, schema, table)
}

func ExecuteQuery(id string, query string) (client.QueryResult, error) {
	dbClient, exists := dbClients[id]
	if !exists {
		return client.QueryResult{}, fmt.Errorf("no database client for database ID: %s", id)
	}

	return dbClient.ExecuteQuery(query)
}

type Connection struct {
	ID               string `json:"id"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	ConnectionString string `json:"connection_string"`
}
