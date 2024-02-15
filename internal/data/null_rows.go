package data

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"

	_ "github.com/libsql/libsql-client-go/libsql"
)

type NullRows struct {
	DateId 		int64
	UserId 		int64
	MergeRequestId 	int64
	RepositoryId 	int64
}

var nullRows *NullRows = nil


func (s *Store) GetNullRows() (*NullRows, error) {

	if nullRows != nil {
		fmt.Println("Returning cached null rows")
		return nullRows, nil
	}

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	nullRows = &NullRows{}

	err = db.QueryRow("SELECT dates_id, users_id, merge_requests_id, repository_id FROM transform_null_rows LIMIT 1;").Scan(
		&nullRows.DateId, &nullRows.UserId,
		&nullRows.MergeRequestId, &nullRows.RepositoryId,
	)

	if err != nil {
		return nil, err
	}

	return nullRows, nil
}
