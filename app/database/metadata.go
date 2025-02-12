package database

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var metadataDB *sql.DB

func InitMetadataDB(filepath string) {
    var err error
    metadataDB, err = sql.Open("sqlite3", filepath)
    if err != nil {
        log.Fatal(err)
    }
    createMetadataTable()
}

func CloseMetadataDB() {
    if metadataDB != nil {
        metadataDB.Close()
    }
}

func createMetadataTable() {
    query := `
        CREATE TABLE IF NOT EXISTS databases (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            type TEXT NOT NULL,
            connection_string TEXT NOT NULL
        )
    `
    _, err := metadataDB.Exec(query)
    if err != nil {
        log.Fatal(err)
    }
}

