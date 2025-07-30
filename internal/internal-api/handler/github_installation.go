package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/dxta-dev/app/internal/internal-api"
	"github.com/dxta-dev/app/internal/onboarding/workflow"
	"github.com/dxta-dev/app/internal/util"
)

type GithubInstallationRequestBody struct {
	InstallationID int64 `json:"installationId"`
}

func (th *OnboardingHandler) GithubInstallation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := &GithubInstallationRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		fmt.Printf("Issue while parsing body. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	err := th.validate.Struct(body)

	if err != nil {
		fmt.Printf("Bad request body: %v", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	authId := ctx.Value(util.AuthIdCtxKey).(string)

	tenantData, err := api.GetTenantDBDataByAuthId(ctx, authId)

	if err != nil {
		fmt.Printf("Failed to retrieve tenant db data: %v", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Failed to retrieve tenant db data"}, http.StatusInternalServerError)
		return
	}

	_, err = workflow.ExecuteAfterGithubInstallationWorkflow(ctx, th.temporalClient, workflow.ExecuteAfterGithubInstallationParams{
		TemporalOnboardingQueueName: th.config.TemporalOnboardingQueueName,
		AuthID:                      authId,
		InstallationID:              body.InstallationID,
		DBURL:                       tenantData.DBUrl,
		DBDomainName:                tenantData.Domain,
	})

	if err != nil {
		fmt.Printf("Failed to execute AfterGithubInstallationWorkflow with error: %s", err.Error())
		util.JSONError(
			w,
			util.ErrorParam{Error: "Failed to execute AfterGithubInstallationWorkflow"},
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
