package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dxta-dev/app/internal/api"
	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/util"
)

func DetailedCycleTimeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apiState, err := api.NewAPIState(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	weekParam := r.URL.Query().Get("weeks")
	weeksArray := util.GetWeeksArray(weekParam)
	weeksSorted := util.SortISOWeeks(weeksArray)

	query := data.BuildDetailedCycleTimeQuery(weeksSorted, apiState.TeamId)
	cycleTimes, err := apiState.DB.GetDetailedCycleTime(
		ctx,
		query,
		apiState.Org,
		apiState.Repo,
		weeksSorted,
		apiState.TeamId,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cycleTimes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

