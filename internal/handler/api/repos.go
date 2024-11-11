package api

import (
	"net/http"

	"github.com/dxta-dev/app/internal/data/api"
	"github.com/labstack/echo/v4"
)

func ReposHandler(c echo.Context) error {

	ctx := c.Request().Context()

	reposDB, err := GetReposDB()
	if err != nil {
		return err
	}

	defer reposDB.Close()

	repos, err := api.GetRepos(ctx, reposDB)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, repos)
}
