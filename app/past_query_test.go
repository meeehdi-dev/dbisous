package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var pastQuery = "SELECT * FROM connection;"

func TestInsertPastQuery(t *testing.T) {
	db, err := InitMetadataDB(":memory:")
	assert.Equal(t, err, nil, err)
	defer db.Close()

	err = insertPastQuery(db, pastQuery)
	assert.Equal(t, err, nil, err)
}

func TestGetPastQueries(t *testing.T) {
	db, err := InitMetadataDB(":memory:")
	assert.Equal(t, err, nil, err)
	defer db.Close()

	insertPastQuery(db, pastQuery)

	queries, err := getPastQueries(db)
	assert.Equal(t, err, nil, err)

	query := queries[0]
	assert.Equal(t, query.Query, pastQuery, "Wrong past query")
}

func TestDeletePastQuery(t *testing.T) {
	db, err := InitMetadataDB(":memory:")
	assert.Equal(t, err, nil, err)
	defer db.Close()

	insertPastQuery(db, pastQuery)

	queries, err := getPastQueries(db)
	query := queries[0]

	err = deletePastQuery(db, query.ID)

	queries, err = getPastQueries(db)
	assert.Equal(t, len(queries), 0, "Past query still exists")
}
