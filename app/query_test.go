package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConnectionDatabases(t *testing.T) {
	r := require.New(t)
	db, err := InitMetadataDB(":memory:")
	r.Equal(err, nil, err)
	defer db.Close()

	// TODO:
}
