package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	api "github.com/dxta-dev/app/internal/internal-api"
	"github.com/dxta-dev/app/internal/onboarding/workflows"
	"github.com/dxta-dev/app/internal/util"
	"github.com/go-chi/chi/v5"
)

type ProvisionGithubInstallationDataResponse struct {
	Message string `json:"message"`
}

func (t *OnboardingTemporal) ProvisionGithubInstallationData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	installationId, err := strconv.ParseInt(chi.URLParam(r, "installation_id"), 10, 64)

	if err != nil {
		fmt.Printf("Issue while parsing installation id URL param. Error: %s", err.Error())
		util.JSONError(w, util.ErrorParam{Error: "Bad Request"}, http.StatusBadRequest)
		return
	}

	authId := ctx.Value(util.AuthIdCtxKey).(string)

	tenantData, err := api.GetTenantDBUrlByAuthId(ctx, authId)

	if err != nil {
		util.JSONError(w, util.ErrorParam{Error: "Internal Server Error"}, http.StatusInternalServerError)
		return
	}

	workflows.ExecuteGithubInstallationDataProvision(
		ctx,
		t.temporalClient,
		workflows.Args{
			TemporalOnboardingQueueName: t.config.TemporalOnboardingNamespace,
			InstallationId:              installationId,
			AuthId:                      authId,
			DBUrl:                       tenantData.DBUrl,
		})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ProvisionGithubInstallationDataResponse{Message: "Success"}); err != nil {
		fmt.Printf("Issue while formatting response. Error: %s", err.Error())
		util.JSONError(
			w,
			util.ErrorParam{Error: "Internal Server Error"},
			http.StatusInternalServerError,
		)
		return
	}
}
