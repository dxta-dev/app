package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type MergeRequestData struct {
	MergeRequestId int64
}

func (s *Store) GetMergeRequestDetails(MergeRequestId int64) ([]MergeRequestData, error) {
	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}
	fmt.Println("ajdijeviii", MergeRequestId)
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			ev.id
		FROM transform_merge_request_events as ev
		WHERE ev.merge_request=?;
	`
	rows, err := db.Query(query, MergeRequestId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var mergeRequestsData []MergeRequestData

	for rows.Next() {
		var mergeRequest MergeRequestData

		if err := rows.Scan(
			&mergeRequest.MergeRequestId,
		); err != nil {
			log.Fatal(err)
		}

		mergeRequestsData = append(mergeRequestsData, mergeRequest)
	}

	return mergeRequestsData, nil
}
