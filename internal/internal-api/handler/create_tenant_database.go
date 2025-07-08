package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dxta-dev/app/internal/onboarding"
	"github.com/dxta-dev/app/internal/onboarding/workflow"
	"github.com/dxta-dev/app/internal/util"
	"go.temporal.io/sdk/client"
)

type CreateDatabaseRequestBody struct {
	DBName           string `json:"dbName"`
	OrganizationName string `json:"organizationName"`
}

type TemporalClient struct {
	temporalClient client.Client
	config         onboarding.Config
}

func NewTemporalHandler(temporalClient client.Client, config onboarding.Config) *TemporalClient {
	return &TemporalClient{
		temporalClient,
		config,
	}
}

func (tc *TemporalClient) CreateTenantDB(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := &CreateDatabaseRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		fmt.Printf("Issue while parsing body. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	if body.DBName == "" || body.OrganizationName == "" {
		fmt.Printf("Bad request body: %v", body)
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	authId := ctx.Value(util.AuthIdCtxKey).(string)

	_, err := workflow.ExecuteCreateTenantDBWorkflow(ctx, tc.temporalClient, workflow.ExecuteCreateTenantDBWorkflowParams{
		TemporalOnboardingQueueName: tc.config.TemporalOnboardingQueueName,
		AuthID:                      authId,
		DBName:                      body.DBName,
		OrganizationName:            body.OrganizationName,
	})

	if err != nil {
		fmt.Printf("Failed to execute CreateTenantDBWorkflow with error: %s", err.Error())
		util.JSONError(
			w,
			util.ErrorParam{Error: "Failed to execute CreateTenantDBWorkflow"},
			http.StatusInternalServerError,
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(struct{ Message string }{Message: "OK"}); err != nil {
		fmt.Printf("Issue while formatting response. Error: %s", err.Error())
		util.JSONError(
			w,
			util.ErrorParam{Error: "Internal Server Error"},
			http.StatusInternalServerError,
		)
		return
	}

}
