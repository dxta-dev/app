package activities

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dxta-dev/app/internal/onboarding"
	"github.com/dxta-dev/app/internal/otel"
)

type UserActivites struct {
	onboardingConfig onboarding.Config
}

func NewUserActivites(onboardingConfig onboarding.Config) *UserActivites {
	return &UserActivites{
		onboardingConfig: onboardingConfig,
	}
}

func (a *UserActivites) CountUsers(ctx context.Context) (int, error) {
	db, err := sql.Open(otel.GetDriverName(), a.onboardingConfig.UsersDSN)
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
