package worker

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/activity"
	temporal_client "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/workflows"
	"go.temporal.io/sdk/worker"
)

func RunWebhookWorker(hostPort, namespace string) error {
	// Initialize a Temporal Client
	// Specify the Namespace in the Client options
	temporalClient, err := temporal_client.TemporalClient(hostPort, namespace)
	if err != nil {
		return fmt.Errorf("unable to create a Temporal Client: %v", err)
	}
	if temporalClient == nil {
		return fmt.Errorf("temporal Client is nil")
	}
	defer temporalClient.Close()
	// Create a new Worker
	wrkr := worker.New(temporalClient, workflows.WEBHOOK_CALLS_TASK_QUEUE, worker.Options{
		// MaxConcurrentSessionExecutionSize: 10,
	})
	// Register Workflows
	wrkr.RegisterWorkflow(workflows.WebhookWorkflow)
	// Register Activities
	wrkr.RegisterActivity(activity.WebhookActivity)
	wrkr.RegisterActivity(activity.NotifyUserActivity)
	// Start the Worker Process
	err = wrkr.Run(worker.InterruptCh())
	if err != nil {
		return fmt.Errorf("unable to start the Worker Process: %v", err)
	}
	return nil
}
