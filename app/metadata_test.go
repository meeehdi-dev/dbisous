package app

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestMetadata(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB("file:test.db?mode=memory")
	r.Equal(err, nil, err)
	defer db.Close()

	err = db.Ping()
	r.Equal(err, nil, err)

	err = createConnectionTable(db)
	r.Equal(err, nil, err)

	err = createPastQueryTable(db)
	r.Equal(err, nil, err)
}
