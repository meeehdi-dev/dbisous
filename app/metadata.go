package app

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitMetadataDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createConnectionTable(db)
	if err != nil {
		return nil, err
	}

	err = createPastQueryTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseMetadataDB() {
	metadataDB.Close()
}

func createConnectionTable(db *sql.DB) error {
	_, err := db.Exec(`
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

func createPastQueryTable(db *sql.DB) error {
	_, err := db.Exec(`
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
