package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/dxta-dev/app/internal/internal-api"
	"github.com/dxta-dev/app/internal/util"
	"github.com/go-playground/validator/v10"
)

type CreateTeamRequestBody struct {
	TeamName string `json:"teamName"`
}

type CreateTeamResponse struct {
	TeamId int64 `json:"team_id" validate:"required"`
}

type CreateTeamHandler struct {
	validate *validator.Validate
}

func NewCreateTeamHandler(validate *validator.Validate) *CreateTeamHandler {
	return &CreateTeamHandler{
		validate,
	}
}

func (cth CreateTeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := &CreateTeamRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		fmt.Printf("Issue while parsing body. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	err := cth.validate.Struct(body)

	if err != nil {
		fmt.Printf("Bad request body: %v", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	authId := ctx.Value(util.AuthIdCtxKey).(string)

	apiState, err := api.InternalApiState(authId, ctx)

	if err != nil {
		util.JSONError(w, util.ErrorParam{Error: "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	organizationId, err := apiState.DB.GetOrganizationIdByAuthId(authId, ctx)

	if err != nil {
		util.JSONError(w, util.ErrorParam{Error: "Bad request"}, http.StatusBadRequest)
		return
	}

	newTeamRes, err := apiState.DB.CreateTeam(body.TeamName, organizationId, ctx)

	if err != nil {
		util.JSONError(
			w,
			util.ErrorParam{Error: "Could not create new team"},
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(CreateTeamResponse{TeamId: newTeamRes.Id}); err != nil {
		fmt.Printf("Issue while formatting response. Error: %s", err.Error())
		util.JSONError(
			w,
			util.ErrorParam{Error: "Internal Server Error"},
			http.StatusInternalServerError,
		)
		return
	}
}
