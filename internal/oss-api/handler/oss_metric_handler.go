package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dxta-dev/app/internal/data"
	api "github.com/dxta-dev/app/internal/oss-api"
	"github.com/dxta-dev/app/internal/util"
)

type OSSMetricQueryBuilder[T data.QueryKeys] func(weeks []string, team *int64) data.Query[T]

func OSSMetricHandler[T data.Executable[T]](
	queryBuilder OSSMetricQueryBuilder[T],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		apiState, err := api.NewAPIState(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		weekParam := r.URL.Query().Get("weeks")

		weeksArray := util.GetWeeksArray(weekParam)
		weeksSorted := util.SortISOWeeks(weeksArray)

		query := queryBuilder(weeksSorted, apiState.TeamId)

		var key T
		result, err := key.Execute(
			ctx,
			&apiState.DB,
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
}
