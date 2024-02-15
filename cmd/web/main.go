package main

import (
	"github.com/dxta-dev/app/internal/handler"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/util"

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

	config, err := util.GetConfig()
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

	app := &handler.App{
		HTMX:           htmx.New(),
		BuildTimestamp: strconv.FormatInt(t.Unix(), 10),
		DebugMode: DEBUG == "true",
	}

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.GzipWithConfig(echoMiddleware.GzipConfig{Level: 6}))

	e.Use(middleware.HtmxMiddleware)

	e.GET("/timestamp", func(c echo.Context) error {
		return c.String(http.StatusOK, app.BuildTimestamp)
	})
	e.GET("/*", app.PublicHandler())

	e.GET("/", app.Home)

	g := e.Group("")

	g.Use(middleware.ConfigMiddleware(config))
	g.Use(middleware.TenantMiddleware)

	g.GET("/dashboard", app.Dashboard)
	g.GET("/merge-request/:mrid", app.GetMergeRequestInfo)
	g.DELETE("/merge-request/:mrid", app.RemoveMergeRequestInfo)

	g.GET("/metrics", app.Metrics)
	g.GET("/metrics/quality", app.QualityMetricsPage)
	g.GET("/metrics/throughput", app.ThroughputMetricsPage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

}
