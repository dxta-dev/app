package handler

import (
	"encoding/json"
	"net/http"

	api "github.com/dxta-dev/app/internal/oss-api"
	"github.com/dxta-dev/app/internal/oss-api/data"
)

func ReposHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reposDB, err := api.GetReposDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reposDB.Close()

	repos, err := data.GetRepos(ctx, reposDB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(repos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
