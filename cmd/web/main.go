package main

import (
	"context"
	"strings"

	"dxta-dev/app/internals/handlers"
	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	app := &handlers.App{
		HTMX: htmx.New(),
	}

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(tenantMiddleware)
	e.Use(htmxMiddleware)

	e.Static("/", "public")

	e.GET("/", app.Home)

	e.GET("/database", app.Database)

	e.Logger.Fatal(e.Start(":3000"))
}

func tenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		host := c.Request().Host
		parts := strings.Split(host, ".")

		if len(parts) > 2 {
			subdomain := parts[0]

			if subdomain == "dxta" {
				c.Set("is_root", true)
			} else {
				c.Set("is_root", false)
			}

			c.Set("subdomain", subdomain)
		}

		return next(c)
	}
}

func htmxMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		hxh := htmx.HxRequestHeader{
			HxBoosted:               htmx.HxStrToBool(c.Request().Header.Get("HX-Boosted")),
			HxCurrentURL:            c.Request().Header.Get("HX-Current-URL"),
			HxHistoryRestoreRequest: htmx.HxStrToBool(c.Request().Header.Get("HX-History-Restore-Request")),
			HxPrompt:                c.Request().Header.Get("HX-Prompt"),
			HxRequest:               htmx.HxStrToBool(c.Request().Header.Get("HX-Request")),
			HxTarget:                c.Request().Header.Get("HX-Target"),
			HxTriggerName:           c.Request().Header.Get("HX-Trigger-Name"),
			HxTrigger:               c.Request().Header.Get("HX-Trigger"),
		}

		ctx = context.WithValue(ctx, htmx.ContextRequestHeader, hxh)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
