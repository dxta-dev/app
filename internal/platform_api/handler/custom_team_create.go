package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dxta-dev/app/internal/platform_api"
	"github.com/go-chi/jwtauth/v5"
)

type RequestBody struct {
	OrganizationId string `json:"organizationId"`
	TeamName       string `json:"teamName"`
}

func CustomTeamCreate(w http.ResponseWriter, r *http.Request) {
	body := &RequestBody{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	enableJWTAuth := os.Getenv("ENABLE_JWT_AUTH")
	organizationId := body.OrganizationId

	if enableJWTAuth == "true" {
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if claims["organizationId"] == "" || claims["organizationId"] != organizationId {
			http.Error(w, "Invalid organization id", http.StatusBadRequest)
			return
		}
	}

	apiState, err := platform_api.PlatformApiState(r, organizationId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	newTeamRes, err := apiState.DB.CreateCustomTeam(body.TeamName, organizationId, ctx)

	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": newTeamRes.Id}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
