package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dxta-dev/app/internal/onboarding"
	"github.com/dxta-dev/app/internal/onboarding/workflow"
	"github.com/dxta-dev/app/internal/util"
	"go.temporal.io/sdk/client"
)

type UsersCountResponse struct {
	Count int `json:"count"`
}

type Users struct {
	temporalClient client.Client
	config         onboarding.Config
}

func NewUsers(temporalClient client.Client, config onboarding.Config) *Users {
	return &Users{
		temporalClient: temporalClient,
		config:         config,
	}
}

func (u *Users) UsersCount(w http.ResponseWriter, r *http.Request) {
	out, err := workflow.ExecuteCountUsersWorkflow(r.Context(), u.temporalClient, u.config)
	if err != nil {
		log.Fatal(errors.Unwrap(err))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(UsersCountResponse{Count: out}); err != nil {
		fmt.Printf("Issue while formatting response. Error: %s", err.Error())
		util.JSONError(
			w,
			util.ErrorParam{Error: "Internal Server Error"},
			http.StatusInternalServerError,
		)
		return
	}
}
