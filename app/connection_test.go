package app

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.NoError(err, "Setup: Failed to initialize in-memory DB")

	t.Cleanup(func() {
		err := db.Close()
		r.NoError(err, "Teardown: Failed to close DB")
	})
	return db
}

// createTestConnection adds a connection to the metadata DB and returns the full Connection object including its ID.
func createTestConnection(t *testing.T, db *sql.DB, conn Connection) Connection {
	t.Helper()
	r := require.New(t)

	err := createConnection(db, conn)
	r.NoError(err, "Setup: Failed to create test connection %q", conn.Name)

	// Retrieve the connection to get its generated ID
	connections, err := getConnections(db)
	r.NoError(err, "Setup: Failed to retrieve connections after creation")
	// Find the specific connection we just added (useful if helper is used multiple times in one test)
	for _, c := range connections {
		if c.Name == conn.Name && c.Type == conn.Type && c.ConnectionString == conn.ConnectionString {
			return c
		}
	}
	r.FailNow("Setup: Could not find the newly created connection %q", conn.Name)
	return Connection{} // Should not be reached
}

var testConnectionName = "DBisous Test"
var testPostgresConnectionString = "postgres://postgres:postgres@localhost:5432/dbisous_test?sslmode=disable"
var testMariaConnectionString = "root:mysql@tcp(localhost:3306)/dbisous_test"
var testMysqlConnectionString = "root:mysql@tcp(localhost:33306)/dbisous_test"

func TestCreateConnection(t *testing.T) {
	r := require.New(t)
	db := setupTestDB(t)

	testConn := Connection{Type: PostgreSQL, Name: testConnectionName, ConnectionString: testPostgresConnectionString}
	err := createConnection(db, testConn)
	r.NoError(err)

	connections, err := getConnections(db)
	r.NoError(err)
	r.Equal(1, len(connections), "Expected 1 connection after creation")

	connection := connections[0]
	r.NotEmpty(connection.ID, "Connection ID should not be empty")
	r.Equal(testConnectionName, connection.Name, "Wrong connection name")
	r.Equal(PostgreSQL, connection.Type, "Wrong connection type")
	r.Equal(testPostgresConnectionString, connection.ConnectionString, "Wrong connection string")
}

func TestUpdateConnection(t *testing.T) {
	r := require.New(t)
	db := setupTestDB(t)

	// Create initial connection
	initialDetails := Connection{Type: MySQL, Name: "Old Name", ConnectionString: "old/conn/string"}
	createdConn := createTestConnection(t, db, initialDetails)

	// Prepare updated details
	updatedConn := createdConn // Copy ID and timestamps
	updatedConn.Name = testConnectionName
	updatedConn.Type = PostgreSQL
	updatedConn.ConnectionString = testPostgresConnectionString

	// Perform update
	err := updateConnection(db, updatedConn)
	r.NoError(err)

	// Verify update
	connections, err := getConnections(db)
	r.NoError(err)
	r.Equal(1, len(connections), "Expected 1 connection after update")

	fetchedConn := connections[0]
	r.Equal(createdConn.ID, fetchedConn.ID, "Connection ID should not change on update")
	r.Equal(testConnectionName, fetchedConn.Name, "Wrong connection name after update")
	r.Equal(PostgreSQL, fetchedConn.Type, "Wrong connection type after update")
	r.Equal(testPostgresConnectionString, fetchedConn.ConnectionString, "Wrong connection string after update")
	// Optionally check updated_at timestamp if relevant
}

func TestTestConnection(t *testing.T) {
	r := require.New(t)
	// Note: These tests rely on external DBs being available as configured
	// Consider using testcontainers for more robust integration testing.

	testCases := []struct {
		name          string
		dbType        ConnectionType
		connStr       string
		expectSuccess bool
	}{
		{"SQLite Valid Memory", SQLite, ":memory:", true},
		// {"SQLite Valid File", SQLite, "./test_connection.db", true}, // Creates the file
		// {"SQLite Invalid Path", SQLite, "/invalid/path/nonexistent.db", false}, // Behavior depends on permissions/OS

		{"PostgreSQL Valid", PostgreSQL, testPostgresConnectionString, true},
		{"PostgreSQL Invalid DB", PostgreSQL, "postgres://postgres:postgres@localhost:5432/inexisting?sslmode=disable", false},
		{"PostgreSQL Invalid Host", PostgreSQL, "postgres://postgres:postgres@invalidhost:5432/dbisous_test?sslmode=disable", false},
		{"PostgreSQL Invalid Credentials", PostgreSQL, "postgres://invalid:user@localhost:5432/dbisous_test?sslmode=disable", false},

		{"MariaDB Valid", MySQL, testMariaConnectionString, true},
		{"MariaDB Invalid DB", MySQL, "root:mysql@tcp(localhost:3306)/inexisting", false},
		{"MariaDB Invalid Host", MySQL, "root:mysql@tcp(invalidhost:3306)/dbisous_test", false},
		{"MariaDB Invalid Credentials", MySQL, "invalid:user@tcp(localhost:3306)/dbisous_test", false},

		{"MySQL Valid", MySQL, testMysqlConnectionString, true},
		{"MySQL Invalid DB", MySQL, "root:mysql@tcp(localhost:33306)/inexisting", false},
		{"MySQL Invalid Host", MySQL, "root:mysql@tcp(invalidhost:33306)/dbisous_test", false},
		{"MySQL Invalid Credentials", MySQL, "invalid:user@tcp(localhost:33306)/dbisous_test", false},

		{"Unsupported Type", "unknown", "some_string", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testConnection(tc.dbType, tc.connStr)
			if tc.expectSuccess {
				r.NoError(err, "Expected connection test to succeed, but it failed")
			} else {
				r.Error(err, "Expected connection test to fail, but it succeeded")
			}
		})
	}
	// os.Remove("./test_connection.db") // Clean up if file test is enabled
}

func TestConnectConnection(t *testing.T) {
	r := require.New(t)
	db := setupTestDB(t)
	activeConnections := make(map[string]*sql.DB) // Simulate App's active connections

	// Create a connection entry in metadata
	connDetails := Connection{Type: PostgreSQL, Name: testConnectionName, ConnectionString: testPostgresConnectionString}
	createdConn := createTestConnection(t, db, connDetails)

	// Test successful connection
	dbClient, err := connect(activeConnections, db, createdConn.ID)
	r.NoError(err, "Failed to connect to a valid connection ID")
	r.NotNil(dbClient, "Returned DB client should not be nil on success")
	_, ok := activeConnections[createdConn.ID]
	r.True(ok, "Connection ID should be present in activeConnections map after successful connect")

	// Test connection to non-existent ID
	_, err = connect(activeConnections, db, "non-existent-id")
	r.Error(err, "Expected an error when connecting to a non-existent ID")
	// Optionally: Check the specific error type or message

	// Test connection to existing ID but invalid connection string (if testConnection wasn't mandatory before connect)
	// This requires modifying createTestConnection or adding another test case if needed.
}

func TestDisconnectConnection(t *testing.T) {
	r := require.New(t)
	db := setupTestDB(t)
	activeConnections := make(map[string]*sql.DB) // Simulate App's active connections

	// Create and connect a connection
	connDetails := Connection{Type: PostgreSQL, Name: testConnectionName, ConnectionString: testPostgresConnectionString}
	createdConn := createTestConnection(t, db, connDetails)
	dbClient, err := connect(activeConnections, db, createdConn.ID)
	r.NoError(err, "Setup: Failed to connect for disconnect test")
	r.NotNil(dbClient) // Ensure connection actually happened

	// Test successful disconnect
	err = disconnect(activeConnections, createdConn.ID)
	r.NoError(err, "Failed to disconnect a valid, active connection ID")
	_, ok := activeConnections[createdConn.ID]
	r.False(ok, "Connection ID should be removed from activeConnections map after successful disconnect")

	// Test disconnect of already disconnected / non-existent ID
	err = disconnect(activeConnections, createdConn.ID) // Try disconnecting again
	r.Error(err, "Expected an error when disconnecting an ID not in the active map")

	err = disconnect(activeConnections, "non-existent-id")
	r.Error(err, "Expected an error when disconnecting a non-existent ID")
}

func TestDeleteConnection(t *testing.T) {
	r := require.New(t)
	db := setupTestDB(t)

	// Create a connection to delete
	connDetails := Connection{Type: PostgreSQL, Name: testConnectionName, ConnectionString: testPostgresConnectionString}
	createdConn := createTestConnection(t, db, connDetails)

	// Ensure it exists before delete
	connections, err := getConnections(db)
	r.NoError(err)
	r.Equal(1, len(connections), "Expected 1 connection before delete")

	// Perform delete
	err = deleteConnection(db, createdConn.ID)
	r.NoError(err)

	// Verify deletion
	connections, err = getConnections(db)
	r.NoError(err)
	r.Equal(0, len(connections), "Expected 0 connections after delete")

	// Test deleting non-existent ID (optional, depends on deleteConnection behavior)
	// err = deleteConnection(db, "non-existent-id")
	// r.NoError(err) // Or r.Error(err) depending on desired behavior (e.g., SQL not finding row is often not an error)
}
