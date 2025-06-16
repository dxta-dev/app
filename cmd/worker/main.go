package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	hostport := os.Getenv("TEMPORAL_HOSTPORT")
	if hostport == "" {
		log.Println("TEMPORAL_HOSTPORT not set; using default")
	}

	temporalClient, err := client.Dial(client.Options{
		HostPort: hostport,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer temporalClient.Close()

	w := worker.New(temporalClient, "SOME_TASK_QUEUE", worker.Options{})

	// w.RegisterWorkflow(pipeline.SingleActivityWorkflow)
	// w.RegisterActivity(pipeline.SimpleActivity)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Worker failed to start", err)
	}
}
