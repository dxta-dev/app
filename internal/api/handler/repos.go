package handler

import (
	"net/http"

	"github.com/dxta-dev/app/internal/api"
	"github.com/dxta-dev/app/internal/api/data"
	"github.com/labstack/echo/v4"
)

func ReposHandler(c echo.Context) error {

	ctx := c.Request().Context()

	reposDB, err := api.GetReposDB()
	if err != nil {
		return err
	}

	defer reposDB.Close()

	repos, err := data.GetRepos(ctx, reposDB)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, repos)
}
