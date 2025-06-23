package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dxta-dev/app/internal/internal_api"
	"github.com/dxta-dev/app/internal/util"
)

type CreateTeamRequestBody struct {
	TeamName string `json:"teamName"`
}

type CreateTeamResponse struct {
	TeamId int64 `json:"team_id"`
}

func CreateTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := &CreateTeamRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		fmt.Printf("Issue while parsing body. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	organizationId := ctx.Value(util.OrganizationIdCtxKey).(int64)

	if organizationId == 0 || body.TeamName == "" {
		fmt.Printf("No organization id or team name provided. Organization id: %d Team name: %s", organizationId, body.TeamName)
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
	}

	apiState := ctx.Value(util.ApiStateCtxKey).(internal_api.State)

	newTeamRes, err := apiState.DB.CreateTeam(body.TeamName, organizationId, ctx)

	if err != nil {
		util.JSONError(w, util.ErrorParam{Error: "Could not create new team"}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(CreateTeamResponse{TeamId: newTeamRes.Id}); err != nil {
		fmt.Printf("Issue while formatting response. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Internal Server Error"}, http.StatusInternalServerError)
		return
	}
}
