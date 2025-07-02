package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dxta-dev/app/internal/onboarding/workflows"
	"github.com/dxta-dev/app/internal/util"
)

type UsersCountResponse struct {
	Count int `json:"count"`
}

func (t OnboardingTemporal) UsersCount(w http.ResponseWriter, r *http.Request) {
	out, err := workflows.ExecuteCountUsersWorkflow(r.Context(), t.temporalClient, t.config)
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
