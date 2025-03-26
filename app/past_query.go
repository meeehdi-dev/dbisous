package app

import (
	"database/sql"

	"github.com/google/uuid"
)

type PastQuery struct {
	ID       string `json:"id"`
	Query    string `json:"query"`
	LastUsed string `json:"last_used"`
}

func getPastQueries(db *sql.DB) ([]PastQuery, error) {
	rows, err := db.Query(`SELECT id, query, last_used FROM past_query ORDER BY last_used DESC`)
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

func insertPastQuery(db *sql.DB, query string) error {
	queryId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO past_query (id, query) VALUES (?, ?) ON CONFLICT(query) DO UPDATE SET last_used = CURRENT_TIMESTAMP`, queryId.String(), query)
	if err != nil {
		return err
	}

	return nil
}

func deletePastQuery(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM past_query WHERE id = ?`, id)
	return err
}
