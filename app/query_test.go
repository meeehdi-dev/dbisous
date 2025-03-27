package app

import (
	"dbisous/app/client"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConnectionDatabases(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	_ = createConnection(db, Connection{Type: PostgreSQL, Name: testConnectionName, ConnectionString: testPostgresConnectionString})
	connections, err := getConnections(db)
	connection := connections[0]
	connect(activeConnections, db, connection.ID)

	filters := make([]client.QueryFilter, 0)
	filters = append(filters, client.QueryFilter{Column: "datname", Value: "'dbisous_test'"})
	result, err := getConnectionDatabases(connection.ID, client.QueryParams{Offset: 0, Limit: 10, Filter: filters})
	r.Equal(nil, err, err)
	r.Equal(1, len(result.Rows), "Couldn't get databases")

	// TODO: add separate tests for each param (offset, limit, filters, order)
}
