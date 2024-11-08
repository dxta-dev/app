package api

import (
	"net/http"

	"github.com/dxta-dev/app/internal/data/api"
	"github.com/labstack/echo/v4"
)

type RequestBody struct {
	UserEmail string `json:"userEmail"`
	RepoUrl   string `json:"repoUrl"`
}

func WaitlistHandler(c echo.Context) error {
	ctx := c.Request().Context()

	metricsDB, err := GetMetricsDB()
	if err != nil {
		return err
	}

	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if reqBody.UserEmail == "" || reqBody.RepoUrl == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing userEmail or repoUrl"})
	}

	err = api.InsertWaitlistData(metricsDB, ctx, reqBody.UserEmail, reqBody.RepoUrl)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
