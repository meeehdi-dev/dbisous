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

func (a *App) GetConnections() ([]Connection, error) {
	rows, err := metadataDB.Query(`SELECT id, created_at, updated_at, name, type, connection_string FROM connection`)
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

func (a *App) CreateConnection(connection Connection) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	connection.ID = id.String()

	_, err = metadataDB.Exec(`INSERT INTO connection (id, name, type, connection_string)
  VALUES (?, ?, ?, ?)`, connection.ID, connection.Name, connection.Type, connection.ConnectionString)

	return err
}

func (a *App) UpdateConnection(connection Connection) error {
	_, err := metadataDB.Exec(`UPDATE connection
  SET name = ?, type = ?, connection_string = ?, updated_at = CURRENT_TIMESTAMP
  WHERE id = ?`, connection.Name, connection.Type, connection.ConnectionString, connection.ID)
	return err
}

func (a *App) DeleteConnection(id string) error {
	_, err := metadataDB.Exec(`DELETE FROM connection WHERE id = ?`, id)
	return err
}

func (a *App) Connect(id string) (client.DatabaseMetadata, error) {
	var databaseMetadata client.DatabaseMetadata

	var dbType, connectionString string
	err := metadataDB.QueryRow(`SELECT type, connection_string FROM connection WHERE id = ?`, id).Scan(&dbType, &connectionString)
	if err != nil {
		return databaseMetadata, err
	}

	var db *sql.DB
	switch dbType {
	case string(SQLite):
		db, err = sql.Open("sqlite3", connectionString)
		dbClients[id] = &client.SqliteClient{Db: db}
	case string(MySQL):
		db, err = sql.Open("mysql", connectionString)
		dbClients[id] = &client.MysqlClient{Db: db}
	case string(PostgreSQL):
		db, err = sql.Open("postgres", connectionString)
		dbClients[id] = &client.PostgresClient{Db: db}
	default:
		return databaseMetadata, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return databaseMetadata, err
	}

	activeConnections[id] = db

	return dbClients[id].GetDatabaseMetadata()
}

func (a *App) Disconnect(id string) error {
	conn, exists := activeConnections[id]
	if !exists {
		return fmt.Errorf("no active connection for database ID: %s", id)
	}

	delete(dbClients, id)
	return conn.Close()
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
