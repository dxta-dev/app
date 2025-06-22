package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dxta-dev/app/internal/api"
	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/markdown"
	"github.com/dxta-dev/app/internal/util"
)

func CodeChangeMarkdownHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apiState, err := api.NewAPIState(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weekParam := r.URL.Query().Get("weeks")

	weeksArray := util.GetWeeksArray(weekParam)
	weeksSorted := util.SortISOWeeks(weeksArray)

	query := data.BuildCodeChangeQuery(weeksSorted, apiState.TeamId)

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
		"Code Change Metrics",
		`The Code Change engineering metric quantifies the team’s weekly development activity by measuring the total number of lines of code added, modified, or deleted across our repositories.

* **Source**: Computed from commit diffs in our Git version-control system, excluding merge commits and auto-generated files.
* **Aggregation**: Grouped by ISO week (Monday–Sunday).
* **Purpose**:
   * Tracks engineering velocity and throughput over time.
   * Highlights spikes (e.g., major feature work or refactors) and troughs (e.g., stabilization periods, planning, or holidays).
   * Helps correlate process changes (code freezes, new tooling) with fluctuations in developer output.`,
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

func CodeChangeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apiState, err := api.NewAPIState(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weekParam := r.URL.Query().Get("weeks")

	weeksArray := util.GetWeeksArray(weekParam)
	weeksSorted := util.SortISOWeeks(weeksArray)

	query := data.BuildCodeChangeQuery(weeksSorted, apiState.TeamId)

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
