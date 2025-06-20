package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dxta-dev/app/internal-api/api"
	"github.com/dxta-dev/app/internal-api/util"
	"github.com/go-chi/chi/v5"
)

type AddMemberToTeamResponse struct {
	Message string `json:"message"`
}

func AddMemberToTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	organizationId := chi.URLParam(r, "organization_id")

	teamId, err := strconv.ParseInt(chi.URLParam(r, "team_id"), 10, 64)
	if err != nil {
		fmt.Printf("Issue while parsing team id URL param. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	memberId, err := strconv.ParseInt(chi.URLParam(r, "member_id"), 10, 64)
	if err != nil {
		fmt.Printf("Issue while parsing member id URL param. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	apiState, err := api.InternalApiState(r, organizationId)

	if err != nil {
		util.JSONError(w, util.ErrorParam{Error: "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	err = apiState.DB.AddMemberToTeam(teamId, memberId, ctx)

	if err != nil {
		util.JSONError(w, util.ErrorParam{Error: "Could not add member to team"}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(AddMemberToTeamResponse{Message: "success"}); err != nil {
		fmt.Printf("Issue while formatting response. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Internal Server Error"}, http.StatusInternalServerError)
		return
	}
}
