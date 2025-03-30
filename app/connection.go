package app

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"

	"dbisous/app/client"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Connection struct {
	ID               string         `json:"id"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
	Name             string         `json:"name"`
	Type             ConnectionType `json:"type"`
	ConnectionString string         `json:"connection_string"`
}

type ConnectionType string

const (
	SQLite     ConnectionType = "sqlite"
	PostgreSQL ConnectionType = "postgresql"
	MySQL      ConnectionType = "mysql"
)

var AllConnectionTypes = []struct {
	Value  ConnectionType
	TSName string
}{
	{SQLite, "SQLite"},
	{PostgreSQL, "PostgreSQL"},
	{MySQL, "MySQL"},
}

var activeConnections = make(map[string]*sql.DB)
var dbClients = make(map[string]client.DatabaseClient)

func getConnections(db *sql.DB) ([]Connection, error) {
	rows, err := db.Query(`SELECT id, created_at, updated_at, name, type, connection_string FROM connection`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	connections := make([]Connection, 0)
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

func createConnection(db *sql.DB, connection Connection) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	connection.ID = id.String()

	_, err = db.Exec(`INSERT INTO connection (id, name, type, connection_string)
  VALUES (?, ?, ?, ?)`, connection.ID, connection.Name, connection.Type, connection.ConnectionString)

	return err
}

func updateConnection(db *sql.DB, connection Connection) error {
	_, err := db.Exec(`UPDATE connection
  SET name = ?, type = ?, connection_string = ?, updated_at = CURRENT_TIMESTAMP
  WHERE id = ?`, connection.Name, connection.Type, connection.ConnectionString, connection.ID)
	return err
}

func deleteConnection(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM connection WHERE id = ?`, id)
	return err
}

func testConnection(dbType ConnectionType, connectionString string) error {
	var db *sql.DB
	var err error

	switch dbType {
	case SQLite:
		db, err = sql.Open("sqlite3", connectionString)
	case MySQL:
		db, err = sql.Open("mysql", connectionString)
	case PostgreSQL:
		db, err = sql.Open("postgres", connectionString)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func connect(activeConnections map[string]*sql.DB, db *sql.DB, id string) (client.DatabaseMetadata, error) {
	var databaseMetadata client.DatabaseMetadata

	var dbType, connectionString string
	err := db.QueryRow(`SELECT type, connection_string FROM connection WHERE id = ?`, id).Scan(&dbType, &connectionString)
	if err != nil {
		return databaseMetadata, err
	}

	var connectionDb *sql.DB
	switch dbType {
	case string(SQLite):
		connectionDb, err = sql.Open("sqlite3", connectionString)
		dbClients[id] = &client.SqliteClient{Db: connectionDb}
	case string(MySQL):
		connectionDb, err = sql.Open("mysql", connectionString)
		dbClients[id] = &client.MysqlClient{Db: connectionDb}
	case string(PostgreSQL):
		connectionDb, err = sql.Open("postgres", connectionString)
		dbClients[id] = &client.PostgresClient{Db: connectionDb}
	default:
		return databaseMetadata, fmt.Errorf("unsupported database type: %s", dbType)
	}
	if err != nil {
		return databaseMetadata, err
	}

	err = connectionDb.Ping()
	if err != nil {
		return databaseMetadata, err
	}

	activeConnections[id] = connectionDb

	return dbClients[id].GetDatabaseMetadata()
}

func disconnect(activeConnections map[string]*sql.DB, id string) error {
	conn, exists := activeConnections[id]
	if !exists {
		return fmt.Errorf("no active connection for database ID: %s", id)
	}

	delete(dbClients, id)
	delete(activeConnections, id) // Add this line
	return conn.Close()
}
