package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dxta-dev/app/internal/handler"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/util"

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
		DebugMode:      DEBUG == "true",
	}

	app.GenerateNonce()

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.GzipWithConfig(echoMiddleware.GzipConfig{Level: 6}))
	e.Use(echoMiddleware.SecureWithConfig(echoMiddleware.SecureConfig{
		Skipper:               echoMiddleware.DefaultSkipper,
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'; img-src 'self' https://avatars.githubusercontent.com; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://unpkg.com; style-src 'self' 'unsafe-inline';",
	}))

	e.Use(middleware.HtmxMiddleware)
	g := e.Group("")

	g.Use(middleware.ConfigMiddleware(config))
	g.Use(middleware.TenantMiddleware)

	g.GET("/", app.GetCrawlInstancesInfo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

}
