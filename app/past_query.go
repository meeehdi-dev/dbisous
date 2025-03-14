package app

import (
	"github.com/google/uuid"
)

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

func (a *App) GetPastQueries() ([]PastQuery, error) {
	rows, err := metadataDB.Query(`SELECT id, query, last_used FROM past_query ORDER BY last_used DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pastQueries := make([]PastQuery, 0)
	for rows.Next() {
		var pastQuery PastQuery
		err := rows.Scan(&pastQuery.ID, &pastQuery.Query, &pastQuery.LastUsed)
		if err != nil {
			return nil, err
		}
		pastQueries = append(pastQueries, pastQuery)
	}

	return pastQueries, nil
}

func InsertPastQuery(query string) error {
	queryId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	_, err = metadataDB.Exec(`INSERT INTO past_query (id, query) VALUES (?, ?) ON CONFLICT(query) DO UPDATE SET last_used = CURRENT_TIMESTAMP`, queryId.String(), query)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) DeletePastQuery(id string) error {
	_, err := metadataDB.Exec(`DELETE FROM past_query WHERE id = ?`, id)
	return err
}
