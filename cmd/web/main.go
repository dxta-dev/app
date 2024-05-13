package main

import (
	"github.com/dxta-dev/app/internal/handler"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/util"

	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var BUILDTIME string
var DEBUG string

var tracer = otel.Tracer("dxta-dev")

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

	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	app := &handler.App{
		HTMX:           htmx.New(),
		BuildTimestamp: strconv.FormatInt(t.Unix(), 10),
		DebugMode:      DEBUG == "true",
	}

	app.GenerateNonce()

	e := echo.New()
	e.Use(otelecho.Middleware("app"))
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

	e.GET("/timestamp", func(c echo.Context) error {
		return c.String(http.StatusOK, app.BuildTimestamp)
	})
	e.GET("/*", app.PublicHandler())

	g := e.Group("")

	g.Use(middleware.ConfigMiddleware(config))
	g.Use(middleware.TenantMiddleware)

	g.GET("/", app.DashboardPage)

	g.GET("/mr-info/:mrid", app.GetMergeRequestInfo)
	g.DELETE("/mr-info/:mrid", app.RemoveMergeRequestInfo)
	g.GET("/mr-stack/in-progress", app.GetMergeRequestInProgressStack)
	g.GET("/mr-stack/ready-to-merge", app.GetMergeRequestReadyToMergeStack)
	g.GET("/mr-stack/waiting-for-review", app.GetMergeRequestWaitingForReviewStack)
	g.GET("/mr-stack/merged", app.GetMergeRequestMergedStack)
	g.GET("/mr-stack/closed", app.GetMergeRequestClosedStack)

	g.GET("/mr/:mrid", app.GetMergeRequestDetails)

	g.GET("/metrics/quality", app.QualityMetricsPage)
	g.GET("/metrics/throughput", app.ThroughputMetricsPage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
