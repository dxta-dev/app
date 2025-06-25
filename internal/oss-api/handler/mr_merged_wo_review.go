package handler

import (
	"encoding/json"
	"net/http"

	api "github.com/dxta-dev/app/internal/oss-api"
	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/util"
)

func MRsMergedWithoutReviewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apiState, err := api.NewAPIState(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	weekParam := r.URL.Query().Get("weeks")
	weeksArray := util.GetWeeksArray(weekParam)
	weeksSorted := util.SortISOWeeks(weeksArray)

	query := data.BuildMRsMergedWithoutReviewQuery(weeksSorted, apiState.TeamId)
	result, err := apiState.DB.GetAggregatedValues(
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

