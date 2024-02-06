package handlers

import (
	"context"
	"dxta-dev/app/internal/templates"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) Home(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:   "Home",
		Boosted: h.HxBoosted,
		CacheBust: a.BuildTimestamp,
	}
	components := templates.Home(page, page.Title)
	return components.Render(context.Background(), c.Response().Writer)
}
