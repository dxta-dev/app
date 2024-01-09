package main

import (
	"dxta-dev/app/internal/handlers"
	"dxta-dev/app/internal/middlewares"
	"fmt"

	"github.com/donseba/go-htmx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	app := &handlers.App{
		HTMX: htmx.New(),
	}

	if err := godotenv.Load(); err != nil {
		return
	}

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(middlewares.TenantMiddleware)
	e.Use(middlewares.HtmxMiddleware)

	e.GET("/*", handlers.PublicHandler())

	e.GET("/", app.Home)

	e.GET("/database", app.Database)
	e.GET("/database/:week", app.Database)

	e.GET("/charts", app.Charts)

	e.GET("/swarm", app.Swarm)

	e.GET("/oss", app.OSSIndex)

	err := middlewares.LoadTenants()
	if err != nil {
		fmt.Print("HERE", err)
	}
	e.Logger.Fatal(e.Start(":3000"))

}
