
package handlers

import (
	"context"
	"dxta-dev/app/internals/templates"
	"github.com/labstack/echo/v4"
)

func (a *App) OSSIndex(c echo.Context) error {
	components := templates.OSSIndex()
	return components.Render(context.Background(), c.Response().Writer)
}
