package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dxta-dev/app/internal/onboarding/workflow"
	"github.com/dxta-dev/app/internal/util"
)

type GithubInstallationRequestBody struct {
	InstallationID int64  `json:"installationId"`
	DBURL          string `json:"dbUrl"`
	DBDomainName   string `json:"dbDomainName"`
}

func (th *TemporalHandler) GithubInstallation(w http.ResponseWriter, r *http.Request) {
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

	_, err = workflow.ExecuteAfterGithubInstallationWorkflow(ctx, th.temporalClient, workflow.ExecuteAfterGithubInstallationParams{
		TemporalOnboardingQueueName: th.config.TemporalOnboardingQueueName,
		AuthID:                      authId,
		InstallationID:              body.InstallationID,
		DBURL:                       body.DBURL,
		DBDomainName:                body.DBDomainName,
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
