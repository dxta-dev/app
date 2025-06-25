package onboarding

import (
	"errors"
	"log"
	"os"
)

type Config struct {
	TemporalHostPort            string
	TemporalOnboardingNamespace string
	TemporalOnboardingQueueName string
	UsersDSN                    string
}

func LoadConfig() (*Config, error) {
	var hostport, onboardingNamespace, onboardingTaskQueue, usersDSN string

	if hostport = os.Getenv("TEMPORAL_HOSTPORT"); hostport == "" {
		log.Println("TEMPORAL_HOSTPORT not set; using default")
	}

	if onboardingNamespace = os.Getenv("TEMPORAL_ONBOARDING_NAMESPACE"); onboardingNamespace == "" {
		return nil, errors.New("TEMPORAL_ONBOARDING_NAMESPACE is not defined")
	}

	if onboardingTaskQueue = os.Getenv("TEMPORAL_ONBOARDING_TASK_QUEUE"); onboardingTaskQueue == "" {
		return nil, errors.New("TEMPORAL_ONBOARDING_TASK_QUEUE is not defined")
	}

	if usersDSN = os.Getenv("USERS_DSN"); usersDSN == "" {
		return nil, errors.New("USERS_DSN is not defined")
	}

	return &Config{
		TemporalHostPort:            hostport,
		TemporalOnboardingNamespace: onboardingNamespace,
		TemporalOnboardingQueueName: onboardingTaskQueue,
		UsersDSN:                    usersDSN,
	}, nil
}
