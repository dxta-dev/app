package main

import (
	"context"
	"log"

	"github.com/dxta-dev/app/internal/onboarding"
	"github.com/dxta-dev/app/internal/onboarding/activities"
	"github.com/dxta-dev/app/internal/onboarding/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	cfg, err := onboarding.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	githubConfig, err := activities.LoadGithubConfig()

	if err != nil {
		log.Fatalln("Failed to load github configuration:", err)
	}

	temporalClient, err := client.Dial(client.Options{
		HostPort:  cfg.TemporalHostPort,
		Namespace: cfg.TemporalOnboardingNamespace,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer temporalClient.Close()

	err = activities.InitAppClient()

	if err != nil {
		log.Fatalf("Unable to init app client: %v", err)
	}

	err = onboarding.RegisterNamespace(
		context.Background(),
		cfg.TemporalHostPort,
		cfg.TemporalOnboardingNamespace,
		30,
	)

	if err != nil {
		log.Fatalln("Failed to register Temporal namespace:", err)
	}

	w := worker.New(temporalClient, cfg.TemporalOnboardingQueueName, worker.Options{})

	userActivities := activities.NewUserActivites(
		*cfg,
	)
	githubActivities := activities.InitGHActivities(*githubConfig)
	dbActivities := activities.InitDBActivities()

	w.RegisterWorkflow(workflows.CountUsers)
	w.RegisterWorkflow(workflows.ProvisionGithubInstallationData)
	w.RegisterActivity(userActivities)
	w.RegisterActivity(githubActivities)
	w.RegisterActivity(dbActivities)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Worker failed to start", err)
	}
}
