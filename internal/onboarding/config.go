package onboarding

import (
	"errors"
	"log"
	"os"
)

type Config struct {
	TemporalHostPort  string
	TemporalNamespace string
	TemporalQueueName string
	UsersDSN          string
}

func LoadConfig() (*Config, error) {
	var hostport, namespace, taskQueue, usersDSN string

	if hostport = os.Getenv("TEMPORAL_HOSTPORT"); hostport == "" {
		log.Println("TEMPORAL_HOSTPORT not set; using default")
	}

	if namespace = os.Getenv("TEMPORAL_NAMESPACE"); namespace != "" {
		return nil, errors.New("TEMPORAL_NAMESPACE is not defined")
	}

	if taskQueue = os.Getenv("TEMPORAL_TASK_QUEUE"); taskQueue == "" {
		return nil, errors.New("TEMPORAL_TASK_QUEUE is not defined")
	}

	if usersDSN = os.Getenv("USERS_DSN"); usersDSN == "" {
		return nil, errors.New("USERS_DSN is not defined")
	}

	return &Config{
		TemporalHostPort:  hostport,
		TemporalNamespace: namespace,
		TemporalQueueName: taskQueue,
		UsersDSN:          usersDSN,
	}, nil
}
