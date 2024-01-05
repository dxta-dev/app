package main

import (
	"dxta-dev/app/internals/handlers"
	"dxta-dev/app/internals/middlewares"

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
	e.Use(middlewares.TenantMiddleware)
	e.Use(middlewares.HtmxMiddleware)

	e.Static("/", "public")

	e.GET("/", app.Home)

	e.GET("/database", app.Database)
	e.GET("/database/:week", app.Database)

	e.GET("/charts", app.Charts)

	e.GET("/swarm", app.Swarm)

	e.GET("/oss", app.OSSIndex)

	e.Logger.Fatal(e.Start(":3000"))
}
