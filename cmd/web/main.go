package main

import (
	"dxta-dev/app/internal/handlers"
	"dxta-dev/app/internal/middlewares"
	"dxta-dev/app/internal/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

var BUILDTIME string
var DEBUG string

func main() {

	config, err := utils.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	t, err := time.Parse(time.RFC3339, BUILDTIME)

	if err != nil {
		t = time.Now()
	}

	if DEBUG == "true" {
		log.Printf("--------------------------------------------------")
		log.Printf("Debug mode is enabled")
		log.Printf("Build timestamp: %v", t.Unix())
		log.Printf("--------------------------------------------------")
	}

	app := &handlers.App{
		HTMX:           htmx.New(),
		BuildTimestamp: strconv.FormatInt(t.Unix(), 10),
		DebugMode: DEBUG == "true",
	}

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.GzipWithConfig(echoMiddleware.GzipConfig{Level: 6}))

	e.Use(middlewares.HtmxMiddleware)

	e.GET("/timestamp", func(c echo.Context) error {
		return c.String(http.StatusOK, app.BuildTimestamp)
	})
	e.GET("/*", app.PublicHandler())

	e.GET("/", app.Home)

	e.GET("/charts", app.Charts)

	g := e.Group("")

	g.Use(middlewares.ConfigMiddleware(config))
	g.Use(middlewares.TenantMiddleware)

	g.GET("/dashboard", app.Dashboard)
	g.GET("/merge-request/:mrid", app.GetMergeRequestInfo)
	g.DELETE("/merge-request/:mrid", app.RemoveMergeRequestInfo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

}
