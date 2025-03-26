package app

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	db, err := InitMetadataDB("file:test.db?mode=memory")
	assert.Equal(t, err, nil, err)
	defer db.Close()

	err = db.Ping()
	assert.Equal(t, err, nil, err)

	err = createConnectionTable(db)
	assert.Equal(t, err, nil, err)

	err = createPastQueryTable(db)
	assert.Equal(t, err, nil, err)
}
