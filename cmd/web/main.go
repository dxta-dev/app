package main

import (
	"dxta-dev/app/internal/handlers"
	"dxta-dev/app/internal/middlewares"
	"dxta-dev/app/internal/utils"
	"fmt"
	"log"
	"os"

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
	e.Use(middlewares.HtmxMiddleware)

	e.GET("/*", handlers.PublicHandler())

	e.GET("/", app.Home)

	e.GET("/charts", app.Charts)


	g := e.Group("")

	g.Use(middlewares.ConfigMiddleware(config))
	g.Use(middlewares.TenantMiddleware)

	g.GET("/dashboard", app.Dashboard)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

}
