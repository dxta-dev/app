package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/dxta-dev/app/internal/internal-api"
	"github.com/dxta-dev/app/internal/util"
)

type CreateMemberRequestBody struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
}

type CreateMemberResponse struct {
	MemberId int64 `json:"member_id"`
}

func CreateMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := &CreateMemberRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		fmt.Printf("Issue while parsing body. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	if body.Name == "" {
		fmt.Println("No member name in request body")
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
	}

	authId := ctx.Value(util.AuthIdCtxKey).(string)

	apiState, err := api.InternalApiState(authId, ctx)

	if err != nil {
		util.JSONError(w, util.ErrorParam{Error: "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	newMemberRes, err := apiState.DB.CreateMember(body.Name, body.Email, ctx)

	if err != nil {
		util.JSONError(
			w,
			util.ErrorParam{Error: "Could not create new member"},
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(CreateMemberResponse{MemberId: newMemberRes.Id}); err != nil {
		fmt.Printf("Issue while formatting response. Error: %s", err.Error())
		util.JSONError(
			w,
			util.ErrorParam{Error: "Internal Server Error"},
			http.StatusInternalServerError,
		)
		return
	}
}
