package app

import (
	"database/sql"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
)

var metadataDB *sql.DB

const path = "dbisous/metadata.db"

func InitMetadataDB() error {
	dataFilePath, err := xdg.DataFile(path)
	if err != nil {
		return err
	}

	metadataDB, err = sql.Open("sqlite3", dataFilePath)
	if err != nil {
		return err
	}

	err = metadataDB.Ping()
	if err != nil {
		return err
	}

	err = createMetadataTable()
	if err != nil {
		return err
	}

	return nil
}

func CloseMetadataDB() {
	if metadataDB != nil {
		metadataDB.Close()
	}
}

func createMetadataTable() error {
	_, err := metadataDB.Exec(`
CREATE TABLE IF NOT EXISTS connection (
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  connection_string TEXT NOT NULL
)`)
	if err != nil {
		return err
	}

	return nil
}
