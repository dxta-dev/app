package activity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dxta-dev/app/internal/otel"
)

func CountUsersActivity(ctx context.Context, dsn string) (int, error) {
	db, err := sql.Open(otel.GetDriverName(), dsn)
	if err != nil {
		return 0, fmt.Errorf("failed to open DB: %w", err)
	}
	defer db.Close()

	var count int
	row := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user")
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("query failed: %w", err)
	}
	return count, nil
}
