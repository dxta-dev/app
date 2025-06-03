package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/dxta-dev/app/internal/api/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	return tp, nil
}

func main() {
	isEndpointProvided := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" ||
		os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") != ""

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

		if err := instrruntime.Start(instrruntime.WithMinimumReadMemStatsInterval(60 * time.Second)); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("%v", fmt.Errorf(
			"missing OTEL exporter configuration. Provide one of (OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_EXPORTER_OTLP_TRACES_ENDPOINT)"))
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(middleware.Compress(5, "gzip"))

	apiSecret := os.Getenv("API_SECRET")
	if apiSecret != "" {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				authHeader := req.Header.Get("Authorization")
				if !strings.HasPrefix(authHeader, "Bearer ") {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
					return
				}

				token := strings.TrimPrefix(authHeader, "Bearer ")
				if token != apiSecret {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
					return
				}
				next.ServeHTTP(w, req)
			})
		})
	}

	if isEndpointProvided {
		r.Use(func(next http.Handler) http.Handler {
			return otelhttp.NewHandler(next, "dxta-app")
		})
	}

	isProd := os.Getenv("ENV") == "production"

	srv := &http.Server{
		Addr: ":" + func() string {
			p := os.Getenv("PORT")
			if p == "" {
				return "1323"
			}
			return p
		}(),
		Handler: r,
	}
	if isProd {
		srv.ReadTimeout = 10 * time.Second
		srv.WriteTimeout = 10 * time.Second
		srv.IdleTimeout = 30 * time.Second
	}

	r.Get("/repos", handler.ReposHandler)
	r.Get("/teams/{org}/{repo}", handler.TeamsHandler)
	r.Get("/code-change/{org}/{repo}", handler.CodeChangeHandler)
	r.Get("/coding-time/{org}/{repo}", handler.CodingTimeHandler)
	r.Get("/commits/{org}/{repo}", handler.CommitsHandler)
	r.Get("/cycle-time/{org}/{repo}", handler.CycleTimeHandler)
	r.Get("/detailed-cycle-time/{org}/{repo}", handler.DetailedCycleTimeHandler)
	r.Get("/deploy-freq/{org}/{repo}", handler.DeployFrequencyHandler)
	r.Get("/deploy-time/{org}/{repo}", handler.DeployTimeHandler)
	r.Get("/handover/{org}/{repo}", handler.HandoverHandler)
	r.Get("/merge-freq/{org}/{repo}", handler.MergeFrequencyHandler)
	r.Get("/mr-merged-wo-review/{org}/{repo}", handler.MRsMergedWithoutReviewHandler)
	r.Get("/mr-opened/{org}/{repo}", handler.MRsOpenedHandler)
	r.Get("/mr-pickup-time/{org}/{repo}", handler.MRPickupTimeHandler)
	r.Get("/mr-size/{org}/{repo}", handler.MRSizeHandler)
	r.Get("/review/{org}/{repo}", handler.ReviewHandler)
	r.Get("/review-depth/{org}/{repo}", handler.ReviewDepthHandler)
	r.Get("/review-time/{org}/{repo}", handler.ReviewTimeHandler)
	r.Get("/time-to-merge/{org}/{repo}", handler.TimeToMergeHandler)
	r.Get("/small-mrs/{org}/{repo}", handler.SmallMRsHandler)

	go func() {
		log.Printf("Listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down server gracefully…")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shut down: %v", err)
	}
	log.Println("Server stopped")
}
