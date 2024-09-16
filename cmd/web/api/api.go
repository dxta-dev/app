package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/util"
	"github.com/labstack/echo/v4"
)

type Api struct {
	HTMX           *htmx.HTMX
	BuildTimestamp string
	DebugMode      bool
	Nonce          string
	State          State
}

type State struct {
	Team *int64
}

func (api *Api) GenerateNonce() error {
	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		return err
	}

	encodedNonce := hex.EncodeToString(nonce)
	api.Nonce = encodedNonce
	return nil
}

func (api *Api) GetMRsMergedWithoutReviewHandler(c echo.Context) error {
	r := c.Request()
	fmt.Print("----------------------------------", r.Context())

	store := r.Context().Value(middleware.StoreContextKey).(*data.Store)
	team := api.State.Team

	teamMembers, err := store.GetTeamMembers(team)

	weeks := util.GetLastNWeeks(time.Now(), 12)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	mrCountByWeeks, averageMergedByXWeeks, err := store.GetMRsMergedWithoutReview(weeks, teamMembers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := map[string]interface{}{
		"mrCountByWeeks":        mrCountByWeeks,
		"averageMergedByXWeeks": averageMergedByXWeeks,
	}

	return c.JSON(http.StatusOK, response)
}
