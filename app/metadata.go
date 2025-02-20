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

	err = createConnectionTable()
	if err != nil {
		return err
	}

	err = createPastQueryTable()
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

type Connection struct {
	ID               string         `json:"id"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
	Name             string         `json:"name"`
	Type             ConnectionType `json:"type"`
	ConnectionString string         `json:"connection_string"`
}

func createConnectionTable() error {
	_, err := metadataDB.Exec(`
CREATE TABLE IF NOT EXISTS connection (
  id TEXT NOT NULL PRIMARY KEY,
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

type PastQuery struct {
	ID       string `json:"id"`
	Query    string `json:"query"`
	LastUsed string `json:"last_used"`
}

func createPastQueryTable() error {
	_, err := metadataDB.Exec(`
CREATE TABLE IF NOT EXISTS past_query (
  id TEXT NOT NULL PRIMARY KEY,
  query TEXT NOT NULL UNIQUE,
  last_used TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)`)
	if err != nil {
		return err
	}

	return nil
}
