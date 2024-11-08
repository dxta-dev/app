package api

import (
	"context"
	"database/sql"
	"fmt"
)

func InsertWaitlistData(db *sql.DB, ctx context.Context, userEmail string, repoUrl string) error {
	query := `
		INSERT INTO users_waitlist (user_email, repository_url)
		VALUES (?, ?) ON CONFLICT (user_email, repository_url) DO NOTHING
	`
	_, err := db.ExecContext(ctx, query, userEmail, repoUrl)
	if err != nil {
		return fmt.Errorf("failed to insert or replace waitlist data: %w", err)
	}
	return nil
}