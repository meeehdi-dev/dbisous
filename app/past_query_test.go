package app

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var pastQuery = "SELECT * FROM connection;"

type testSetup struct {
	r  *require.Assertions
	db *sql.DB
}

func setup(t *testing.T) *testSetup {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.NoError(err, "Failed to initialize metadata DB")

	return &testSetup{r: r, db: db}
}

func TestInsertPastQuery(t *testing.T) {
	s := setup(t)
	defer s.db.Close()

	err := insertPastQuery(s.db, pastQuery)
	s.r.NoError(err, "Failed to insert past query")

	queries, err := getPastQueries(s.db)
	s.r.NoError(err, "Failed to get past queries after insertion")
	s.r.Len(queries, 1, "Expected 1 past query after insertion, but got %d", len(queries))
	lastUsed, err := time.Parse(time.RFC3339, queries[0].LastUsed)
	s.r.NoError(err, "Failed to parse last_used timestamp")
	s.r.True(time.Since(lastUsed) < 1*time.Second, "last_used timestamp is not recent (older than 1 second)")
}

func TestGetPastQueries(t *testing.T) {
	s := setup(t)
	defer s.db.Close()

	err := insertPastQuery(s.db, pastQuery)
	s.r.NoError(err, "Failed to insert past query")

	queries, err := getPastQueries(s.db)
	s.r.NoError(err, "Failed to get past queries")
	s.r.Len(queries, 1, "Expected 1 past query, but got %d", len(queries))

	query := queries[0]
	s.r.Equal(query.Query, pastQuery, "Retrieved query does not match inserted query")
}

func TestDeletePastQuery(t *testing.T) {
	s := setup(t)
	defer s.db.Close()

	err := insertPastQuery(s.db, pastQuery)
	s.r.NoError(err, "Failed to insert past query")

	queries, err := getPastQueries(s.db)
	s.r.NoError(err, "Failed to get past queries")
	s.r.Len(queries, 1, "Expected 1 past query, but got %d", len(queries))
	query := queries[0]

	err = deletePastQuery(s.db, query.ID)
	s.r.NoError(err, "Failed to delete past query")

	queries, err = getPastQueries(s.db)
	s.r.NoError(err, "Failed to get past queries after deletion")
	s.r.Equal(len(queries), 0, "Expected 0 past queries after deletion, but got %d", len(queries))
}
