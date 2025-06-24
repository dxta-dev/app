package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dxta-dev/app/internal/internal_api/handler"
	"github.com/dxta-dev/app/internal/onboarding"
	"github.com/dxta-dev/app/internal/onboarding/workflow"
	"github.com/dxta-dev/app/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"go.temporal.io/sdk/client"

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
				semconv.ServiceName("dxta-internal-api"),
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

	cfg, err := onboarding.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	temporalClient, err := client.Dial(client.Options{
		HostPort:  cfg.TemporalHostPort,
		Namespace: cfg.TemporalNamespace,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer temporalClient.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(middleware.Compress(5, "gzip"))

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
				return "1324"
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

	r.Route("/tenant", func(r chi.Router) {
		if os.Getenv("ENABLE_JWT_AUTH") == "true" {
			pubKey, _ := util.GetRawPublicKey()

			tokenAuth := util.CreateAuthVerifier(pubKey)

			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(util.Authenticator())
		}

		// TO-DO Add middleware if we don't authenticate with JWT
		// https://app.plane.so/crocoder/browse/DXTA-307/

		r.Post("/teams", handler.CreateTeam)
		r.Post("/teams/{team_id}/members/{member_id}", handler.AddMemberToTeam)
		r.Post("/members", handler.CreateMember)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`OK`))
	})

	r.Get("/users-count", func(w http.ResponseWriter, r *http.Request) {
		out, err := workflow.ExecuteCountUsersWorkflow(context.Background(), temporalClient, *cfg)
		if err != nil {
			log.Fatal(errors.Unwrap(err))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{ \"count\": " + fmt.Sprintf("%d", out) + "}"))
	})

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
