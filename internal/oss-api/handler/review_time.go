package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dxta-dev/app/internal/data"
	api "github.com/dxta-dev/app/internal/oss-api"
	"github.com/dxta-dev/app/internal/util"
)

func ReviewTimeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apiState, err := api.NewAPIState(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	weekParam := r.URL.Query().Get("weeks")
	weeksArray := util.GetWeeksArray(weekParam)
	weeksSorted := util.SortISOWeeks(weeksArray)

	query := data.BuildReviewTimeQuery(weeksSorted, apiState.TeamId)
	result, err := apiState.DB.GetAggregatedStatistics(
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
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
