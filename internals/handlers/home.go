package handlers

import (
	"context"
	"github.com/dxta-dev/app/internals/templates"
	"github.com/labstack/echo/v4"

)

func (a *App) Home(c echo.Context) error {
	components := templates.Home("DXTA")
	return components.Render(context.Background(), c.Response().Writer)
}
