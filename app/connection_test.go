package app

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

func TestConnection(t *testing.T) {
	db, err := InitMetadataDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var activeConnections = make(map[string]*sql.DB)

	err = createConnection(db, Connection{Type: PostgreSQL, Name: "Inexisting DB", ConnectionString: "postgres://postgres:postgres@localhost:5432/inexisting?sslmode=disable"})
	assert.Equal(t, err, nil, "Couldn't create connection")

	connections, err := getConnections(db)
	assert.Equal(t, len(connections), 1, "Couldn't get connection")
	connection := connections[0]
	assert.Equal(t, connection.Name, "Inexisting DB", "Wrong connection name")
	assert.Equal(t, connection.Type, PostgreSQL, "Wrong connection type")
	assert.Equal(t, connection.ConnectionString, "postgres://postgres:postgres@localhost:5432/inexisting?sslmode=disable", "Wrong connection string")

	connection.Name = "DBisous Test"
	connection.Type = PostgreSQL
	connection.ConnectionString = "postgres://postgres:postgres@localhost:5432/dbisous_test?sslmode=disable"

	err = updateConnection(db, connection)
	assert.Equal(t, err, nil, "Couldn't update connection")

	connections, err = getConnections(db)
	assert.Equal(t, len(connections), 1, "Couldn't get connection")
	connection = connections[0]
	assert.Equal(t, connection.Name, "DBisous Test", "Wrong connection name")
	assert.Equal(t, connection.Type, PostgreSQL, "Wrong connection type")
	assert.Equal(t, connection.ConnectionString, "postgres://postgres:postgres@localhost:5432/dbisous_test?sslmode=disable", "Wrong connection string")

	// SQLite
	err = testConnection(SQLite, ":memory:")
	assert.Equal(t, err, nil, "Couldn't test connection")
	// err = testConnection(SQLite, "./inexisting.db") // doesn't work bc sqlite creates db file automatically
	// assert.NotEqual(t, err, nil, "Tested inexisting connection")
	// PostgreSQL
	err = testConnection(PostgreSQL, "postgres://postgres:postgres@localhost:5432/dbisous_test?sslmode=disable")
	assert.Equal(t, err, nil, "Couldn't test connection")
	err = testConnection(PostgreSQL, "postgres://postgres:postgres@localhost:5432/inexisting?sslmode=disable")
	assert.NotEqual(t, err, nil, "Tested inexisting connection")
	// MySQL
	err = testConnection(MySQL, "root:mysql@tcp(localhost:3306)/dbisous_test")
	assert.Equal(t, err, nil, "Couldn't test connection")
	err = testConnection(MySQL, "root:mysql@tcp(localhost:3306)/inexisting")
	assert.NotEqual(t, err, nil, "Tested inexisting connection")

	_, err = connect(activeConnections, db, connection.ID)
	assert.Equal(t, err, nil, "Couldn't connect to database")
	_, err = connect(activeConnections, db, "")
	assert.NotEqual(t, err, nil, "Connected to inexisting database")
	err = disconnect(activeConnections, connection.ID)
	assert.Equal(t, err, nil, "Couldn't disconnect from database")
	err = disconnect(activeConnections, "")
	assert.NotEqual(t, err, nil, "Disconnected from inexisting database")

	err = deleteConnection(db, connection.ID)
	assert.Equal(t, err, nil, "Couldn't delete connection")

	connections, err = getConnections(db)
	assert.Equal(t, len(connections), 0, "Connection still exists")
}
