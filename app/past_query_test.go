package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var pastQuery = "SELECT * FROM connection;"

func TestInsertPastQuery(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	err = insertPastQuery(db, pastQuery)
	r.Equal(err, nil, err)
}

func TestGetPastQueries(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	insertPastQuery(db, pastQuery)

	queries, err := getPastQueries(db)
	r.Equal(err, nil, err)

	query := queries[0]
	r.Equal(query.Query, pastQuery, "Wrong past query")
}

func TestDeletePastQuery(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	insertPastQuery(db, pastQuery)

	queries, err := getPastQueries(db)
	query := queries[0]

	err = deletePastQuery(db, query.ID)

	queries, err = getPastQueries(db)
	r.Equal(len(queries), 0, "Past query still exists")
}
