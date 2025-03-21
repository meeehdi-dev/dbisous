package app

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreate(t *testing.T) {
	db, err := InitMetadataDB("file:test.db?mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	connections, err := getConnections(db)

	if len(connections) != 0 {
		t.Errorf("nope!")
	}
}
