package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/otel"
	"github.com/go-chi/chi/v5"
)

type APIState struct {
	DB     data.DB
	Org    string
	Repo   string
	TeamId *int64
}

var tenantRepoCache sync.Map

func getCachedTenantRepo(ctx context.Context, org, repo string) (data.TenantRepo, error) {
	cacheKey := strings.ToLower(org + "/" + repo)

	if cached, ok := tenantRepoCache.Load(cacheKey); ok {
		return cached.(data.TenantRepo), nil
	}

	driverName := otel.GetDriverName()
	superURL := os.Getenv("SUPER_DATABASE_URL")
	devToken := os.Getenv("DXTA_DEV_GROUP_TOKEN")
	if superURL == "" || devToken == "" {
		return data.TenantRepo{}, fmt.Errorf("missing SUPER_DATABASE_URL or DXTA_DEV_GROUP_TOKEN")
	}

	reposDB, err := sql.Open(
		driverName,
		superURL+"?authToken="+devToken,
	)
	if err != nil {
		return data.TenantRepo{}, err
	}
	defer reposDB.Close()

	tenantRepo, err := data.GetTenantRepo(ctx, reposDB, org, repo)
	if err != nil {
		return data.TenantRepo{}, err
	}

	tenantRepoCache.Store(cacheKey, tenantRepo)
	return tenantRepo, nil
}

func NewAPIState(r *http.Request) (APIState, error) {
	ctx := r.Context()

	org := chi.URLParam(r, "org")
	repo := chi.URLParam(r, "repo")
	if org == "" || repo == "" {
		return APIState{}, fmt.Errorf("org and repo are required path parameters")
	}

	tenantRepo, err := getCachedTenantRepo(ctx, org, repo)
	if err != nil {
		return APIState{}, err
	}

	db, err := data.NewDB(ctx, tenantRepo)
	if err != nil {
		return APIState{}, err
	}

	var teamInt *int64
	teamParam := r.URL.Query().Get("team")
	if teamParam != "" {
		if t, err := strconv.ParseInt(teamParam, 10, 64); err == nil && t > 0 {
			teamInt = &t
		}
	}

	return APIState{
		DB:     db,
		Org:    tenantRepo.Organization,
		Repo:   tenantRepo.Repository,
		TeamId: teamInt,
	}, nil
}

func GetReposDB() (*sql.DB, error) {
	driverName := otel.GetDriverName()
	superURL := os.Getenv("SUPER_DATABASE_URL")
	devToken := os.Getenv("DXTA_DEV_GROUP_TOKEN")
	if superURL == "" || devToken == "" {
		return nil, fmt.Errorf("missing SUPER_DATABASE_URL or DXTA_DEV_GROUP_TOKEN")
	}
	reposDB, err := sql.Open(
		driverName,
		superURL+"?authToken="+devToken,
	)
	if err != nil {
		return nil, err
	}
	return reposDB, nil
}
