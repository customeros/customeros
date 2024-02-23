package worker

import (
	"log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/activity"
	temporal_client "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/workflows"
	"go.temporal.io/sdk/worker"
)

func RunWebhookWorker(hostPort, namespace string) {
	// Initialize a Temporal Client
	// Specify the Namespace in the Client options
	temporalClient, err := temporal_client.TemporalClient(hostPort, namespace)
	if err != nil {
		log.Fatalln("Unable to create a Temporal Client", err)
	}
	if temporalClient == nil {
		log.Fatalln("Temporal Client is nil")
	}
	defer temporalClient.Close()
	// Create a new Worker
	yourWorker := worker.New(temporalClient, workflows.WEBHOOK_CALLS_TASK_QUEUE, worker.Options{
		// MaxConcurrentSessionExecutionSize: 10,
	})
	// Register Workflows
	yourWorker.RegisterWorkflow(workflows.WebhookWorkflow)
	// Register Activities
	yourWorker.RegisterActivity(activity.WebhookActivity)
	// Start the Worker Process
	err = yourWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start the Worker Process", err)
	}
}
