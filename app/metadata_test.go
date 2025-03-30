package app

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestMetadata(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.NoError(err)
	defer db.Close()

	err = db.Ping()
	r.NoError(err)

	err = createConnectionTable(db)
	r.NoError(err)

	err = createPastQueryTable(db)
	r.NoError(err)
}
