package handler

import (
	"encoding/json"
	"net/http"

	api "github.com/dxta-dev/app/internal/oss-api"
)

func TeamsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apiState, err := api.NewAPIState(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teams, err := apiState.DB.GetTeams(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(teams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

