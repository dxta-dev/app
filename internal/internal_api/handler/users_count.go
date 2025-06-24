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

func UsersCount(w http.ResponseWriter, r *http.Request) {
	cfg, err := onboarding.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	temporalClient, err := client.Dial(client.Options{
		HostPort:  cfg.TemporalHostPort,
		Namespace: cfg.TemporalOnboardingNamespace,
	})
	if err != nil {
		fmt.Printf("Unable to create Temporal client. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Internal Server Error"}, http.StatusInternalServerError)
	}

	defer temporalClient.Close()

	out, err := workflow.ExecuteCountUsersWorkflow(r.Context(), temporalClient, *cfg)
	if err != nil {
		log.Fatal(errors.Unwrap(err))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(UsersCountResponse{Count: out}); err != nil {
		fmt.Printf("Issue while formatting response. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Internal Server Error"}, http.StatusInternalServerError)
		return
	}
}
