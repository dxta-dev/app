package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dxta-dev/app/internal/api/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	instrruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func initTracer(ctx context.Context, res *sdkresource.Resource) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			res,
		),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func main() {

	isEndpointProvided := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" || os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") != ""

	if isEndpointProvided {
		res, err := sdkresource.New(
			context.Background(),
			sdkresource.WithAttributes(
				semconv.ServiceName("dxta-api"),
			),
		)
		if err != nil {
			log.Fatal(err)
		}
		tp, err := initTracer(context.Background(), res)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()

		err = instrruntime.Start(instrruntime.WithMinimumReadMemStatsInterval(60 * time.Second))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("%v", fmt.Errorf("missing OTEL exporter configuration. Provide one of (OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_EXPORTER_OTLP_TRACES_ENDPOINT) options"))
	}

	e := echo.New()

	isProd := os.Getenv("ENV") == "production"

	if isProd {
		e.Debug = false

		e.Server.ReadTimeout = 10 * time.Second
		e.Server.WriteTimeout = 10 * time.Second
		e.Server.IdleTimeout = 30 * time.Second
	}

	if isEndpointProvided {
		e.Use(otelecho.Middleware("dxta-app"))
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	if os.Getenv("API_SECRET") != "" {
		e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("API_SECRET"), nil
		}))
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hell")
	})

	e.GET("/repos", handler.ReposHandler)
	e.GET("/teams/:org/:repo", handler.TeamsHandler)
	e.GET("/code-change/:org/:repo", handler.CodeChangeHandler)
	e.GET("/coding-time/:org/:repo", handler.CodingTimeHandler)
	e.GET("/commits/:org/:repo", handler.CommitsHandler)
	e.GET("/cycle-time/:org/:repo", handler.CycleTimeHandler)
	e.GET("/detailed-cycle-time/:org/:repo", handler.DetailedCycleTimeHandler)
	e.GET("/deploy-freq/:org/:repo", handler.DeployFrequencyHandler)
	e.GET("/deploy-time/:org/:repo", handler.DeployTimeHandler)
	e.GET("/handover/:org/:repo", handler.HandoverHandler)
	e.GET("/merge-freq/:org/:repo", handler.MergeFrequencyHandler)
	e.GET("/mr-merged-wo-review/:org/:repo", handler.MRsMergedWithoutReviewHandler)
	e.GET("/mr-opened/:org/:repo", handler.MRsOpenedHandler)
	e.GET("/mr-pickup-time/:org/:repo", handler.MRPickupTimeHandler)
	e.GET("/mr-size/:org/:repo", handler.MRSizeHandler)
	e.GET("/review/:org/:repo", handler.ReviewHandler)
	e.GET("/review-depth/:org/:repo", handler.ReviewDepthHandler)
	e.GET("/review-time/:org/:repo", handler.ReviewTimeHandler)
	e.GET("/time-to-merge/:org/:repo", handler.TimeToMergeHandler)
	e.GET("/small-mrs/:org/:repo", handler.SmallMRsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
