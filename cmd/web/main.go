package main

import (
	"dxta-dev/app/internal/handlers"
	"dxta-dev/app/internal/middlewares"
	"dxta-dev/app/internal/utils"
	"log"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {

	config, err := utils.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := &handlers.App{
		HTMX: htmx.New(),
	}

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(middlewares.ConfigMiddleware(config))

	if config.IsMultiTenant {
		e.Use(middlewares.MultiTenantMiddleware)
	}

	e.Use(middlewares.HtmxMiddleware)

	e.GET("/*", handlers.PublicHandler())

	e.GET("/", app.Home)

	e.GET("/database", app.Database)
	e.GET("/database/:week", app.Database)

	e.GET("/charts", app.Charts)

	e.GET("/swarm", app.Swarm)

	e.GET("/oss", app.OSSIndex)

	e.Logger.Fatal(e.Start(":3000"))

}
