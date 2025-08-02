package handler

import (
	"net/http"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/markdown"
	api "github.com/dxta-dev/app/internal/oss-api"
	"github.com/dxta-dev/app/internal/util"
)

var ReviewHandler = OSSMetricHandler(data.BuildReviewQuery)

func ReviewMarkdownHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apiState, err := api.NewAPIState(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weekParam := r.URL.Query().Get("weeks")

	weeksArray := util.GetWeeksArray(weekParam)
	weeksSorted := util.SortISOWeeks(weeksArray)

	query := data.BuildReviewQuery(weeksSorted, apiState.TeamId)

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

	m, err := markdown.GetAggregatedValuesMarkdown(
		ctx,
		"Reviewed Metrics",
		``,
		result,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(m)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
