package data

import (
	"database/sql"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type MergeRequestInfo struct {
	MergeRequestId int64
	Title          string
	WebUrl         string
	AuthorName     string
	RepositoryName string
	MrSize         int64
}

func (s *Store) GetMergeRequestInfo(mrId int64) (*MergeRequestInfo, error) {
	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	// TODO: diffs? + - lines
	// TODO: Author avatar, url to author profile
	// TODO: Canon mr id
	// TODO: mr description?
	// TODO: namespace name
	query := `
		SELECT
			metrics.mr_size,
			mr.title,
			mr.web_url,
			author.name,
			tf.name
		FROM transform_merge_requests AS mr
		JOIN transform_merge_request_metrics AS metrics
		ON metrics.merge_request = mr.id
		JOIN transform_merge_request_fact_users_junk AS uj
		ON metrics.users_junk = uj.id
		JOIN transform_forge_users AS author
		ON uj.author = author.id
		JOIN transform_repositories AS tf
		ON metrics.repository = tf.id
		WHERE mr.id = ?;
	`

	mergeRequestInfo := MergeRequestInfo{MergeRequestId: mrId}

	err = db.QueryRow(query, mrId).Scan(
		&mergeRequestInfo.MrSize,
		&mergeRequestInfo.Title,
		&mergeRequestInfo.WebUrl,
		&mergeRequestInfo.AuthorName,
		&mergeRequestInfo.RepositoryName,
	)

	return &mergeRequestInfo, nil
}
