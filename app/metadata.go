package app

import (
	"database/sql"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
)

var metadataDB *sql.DB

func InitMetadataDB(filepath string) error {
	var err error

	dataFilePath, err := xdg.DataFile("dbisous/metadata.db")

	metadataDB, err = sql.Open("sqlite3", dataFilePath)
	if err != nil {
		return err
	}

	err = metadataDB.Ping()
	if err != nil {
		return err
	}

	return createMetadataTable()
}

func CloseMetadataDB() {
	if metadataDB != nil {
		metadataDB.Close()
	}
}

func createMetadataTable() error {
	query := `
CREATE TABLE IF NOT EXISTS connection (
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  connection_string TEXT NOT NULL
)`
	_, err := metadataDB.Exec(query)
	return err
}
