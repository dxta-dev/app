package main

import (
	"context"
	"log"

	"github.com/dxta-dev/app/internal/onboarding"
	"github.com/dxta-dev/app/internal/onboarding/activity"
	"github.com/dxta-dev/app/internal/onboarding/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	cfg, err := onboarding.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	githubConfig, err := onboarding.LoadGithubConfig()

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

	githubAppClient, err := onboarding.NewAppClient(*githubConfig)

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

	userActivities := activity.NewUserActivites(
		*cfg,
	)
	githubInstallationActivities := activity.NewGithubInstallationActivities(*githubAppClient)

	w.RegisterWorkflow(workflow.CountUsers)
	w.RegisterWorkflow(workflow.AfterGithubInstallationWorkflow)
	w.RegisterActivity(userActivities)
	w.RegisterActivity(githubInstallationActivities)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Worker failed to start", err)
	}
}
