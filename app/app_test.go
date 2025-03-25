package app

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInitMetadata(t *testing.T) {
	db, err := InitMetadataDB("file:test.db?mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("Couldn't instantiate test database")
	}
}
