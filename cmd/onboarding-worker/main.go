package main

import (
	"errors"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type config struct {
	TemporalHostPort  string
	TemporalNamespace string
	TemporalQueueName string
}

func loadConfig() (*config, error) {
	var hostport, namespace, taskQueue string

	if hostport = os.Getenv("TEMPORAL_HOSTPORT"); hostport == "" {
		log.Println("TEMPORAL_HOSTPORT not set; using default")
	}

	if namespace = os.Getenv("TEMPORAL_NAMESPACE"); namespace != "" {
		return nil, errors.New("TEMPORAL_NAMESPACE is not defined")
	}

	if taskQueue = os.Getenv("TEMPORAL_TASK_QUEUE"); taskQueue == "" {
		return nil, errors.New("TEMPORAL_TASK_QUEUE is not defined")
	}

	return &config{
		TemporalHostPort:  hostport,
		TemporalNamespace: namespace,
		TemporalQueueName: taskQueue,
	}, nil
}

func main() {
	cfg, err := loadConfig()
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

	w := worker.New(temporalClient, cfg.TemporalQueueName, worker.Options{})

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Worker failed to start", err)
	}
}
