package app

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

var testConnectionName = "DBisous Test"
var testPostgresConnectionString = "postgres://postgres:postgres@localhost:5432/dbisous_test?sslmode=disable"
var testMysqlConnectionString = "root:mysql@tcp(localhost:3306)/dbisous_test"

func TestCreateConnection(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	err = createConnection(db, Connection{Type: PostgreSQL, Name: testConnectionName, ConnectionString: testPostgresConnectionString})
	r.Equal(err, nil, err)

	connections, err := getConnections(db)
	r.Equal(len(connections), 1, "Couldn't get connection")
	connection := connections[0]
	r.Equal(connection.Name, testConnectionName, "Wrong connection name")
	r.Equal(connection.Type, PostgreSQL, "Wrong connection type")
	r.Equal(connection.ConnectionString, testPostgresConnectionString, "Wrong connection string")
}

func TestUpdateConnection(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	_ = createConnection(db, Connection{Type: PostgreSQL, Name: "Inexisting DB", ConnectionString: "postgres://postgres:postgres@localhost:5432/inexisting?sslmode=disable"})

	connections, err := getConnections(db)
	connection := connections[0]

	connection.Name = testConnectionName
	connection.Type = PostgreSQL
	connection.ConnectionString = testPostgresConnectionString

	err = updateConnection(db, connection)
	r.Equal(err, nil, err)

	connections, err = getConnections(db)
	r.Equal(len(connections), 1, "Couldn't get connection")
	connection = connections[0]
	r.Equal(connection.Name, testConnectionName, "Wrong connection name")
	r.Equal(connection.Type, PostgreSQL, "Wrong connection type")
	r.Equal(connection.ConnectionString, testPostgresConnectionString, "Wrong connection string")
}

func TestTestConnection(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	// SQLite
	err = testConnection(SQLite, ":memory:")
	r.Equal(err, nil, err)
	// err = testConnection(SQLite, "./inexisting.db") // doesn't work bc sqlite creates db file automatically
	// r.NotEqual(err, nil, "Tested inexisting connection")
	// PostgreSQL
	err = testConnection(PostgreSQL, testPostgresConnectionString)
	r.Equal(err, nil, err)
	err = testConnection(PostgreSQL, "postgres://postgres:postgres@localhost:5432/inexisting?sslmode=disable")
	r.NotEqual(err, nil, "Tested inexisting connection")
	// MySQL
	err = testConnection(MySQL, testMysqlConnectionString)
	r.Equal(err, nil, err)
	err = testConnection(MySQL, "root:mysql@tcp(localhost:3306)/inexisting")
	r.NotEqual(err, nil, "Tested inexisting connection")
}

func TestConnectConnection(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	var activeConnections = make(map[string]*sql.DB)

	_ = createConnection(db, Connection{Type: PostgreSQL, Name: testConnectionName, ConnectionString: testPostgresConnectionString})

	connections, err := getConnections(db)
	connection := connections[0]

	_, err = connect(activeConnections, db, connection.ID)
	r.Equal(err, nil, err)
	_, err = connect(activeConnections, db, "")
	r.NotEqual(err, nil, "Connected to inexisting database")
}

func TestDisconnectConnection(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	var activeConnections = make(map[string]*sql.DB)

	_ = createConnection(db, Connection{Type: PostgreSQL, Name: testConnectionName, ConnectionString: testPostgresConnectionString})

	connections, err := getConnections(db)
	connection := connections[0]

	_, _ = connect(activeConnections, db, connection.ID)

	err = disconnect(activeConnections, connection.ID)
	r.Equal(err, nil, err)
	err = disconnect(activeConnections, "")
	r.NotEqual(err, nil, "Disconnected from inexisting database")
}

func TestDeleteConnection(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	_ = createConnection(db, Connection{Type: PostgreSQL, Name: testConnectionName, ConnectionString: testPostgresConnectionString})

	connections, err := getConnections(db)
	connection := connections[0]

	err = deleteConnection(db, connection.ID)
	r.Equal(err, nil, err)

	connections, err = getConnections(db)
	r.Equal(len(connections), 0, "Connection still exists")
}
