package handler

import (
	"github.com/dxta-dev/app/internal/onboarding"
	"go.temporal.io/sdk/client"
)

type OnboardingTemporal struct {
	temporalClient client.Client
	config         onboarding.Config
}

func SetupOnboardingTemporal(temporalClient client.Client, config onboarding.Config) *OnboardingTemporal {
	return &OnboardingTemporal{
		temporalClient: temporalClient,
		config:         config,
	}
}
