package workflows

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type WHWorkflowParam struct {
	TargetUrl                  string
	RequestBody                string
	AuthHeaderName             string
	AuthHeaderValue            string
	RetryPolicy                *temporal.RetryPolicy
	Notification               string
	NotificationProviderApiKey string
	NotifyFailure              bool
	NotifyAfterAttempts        int32
}

// Example retry policy
// retryPolicy := &temporal.RetryPolicy{
// 	InitialInterval:        time.Second,
// 	BackoffCoefficient:     2.0,
// 	MaximumInterval:        time.Second * 100, // 100 * InitialInterval
// 	MaximumAttempts:        3,                 // if set to 0 means Unlimited attempts; not inclusive eg. n < 3.
// 	NonRetryableErrorTypes: []string{},        // empty
// }

func WebhookWorkflow(ctx workflow.Context, param WHWorkflowParam) error {
	// Define the Activity Execution options
	retrypolicy := param.RetryPolicy
	if retrypolicy == nil {
		retrypolicy = &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Second * 100,
			MaximumAttempts:    3,
		}
	}
	// StartToCloseTimeout or ScheduleToCloseTimeout must be set
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy:         retrypolicy,
		TaskQueue:           WEBHOOK_CALLS_TASK_QUEUE,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	// Execute the Activity synchronously (wait for the result before proceeding)
	err := workflow.ExecuteActivity(
		ctx,
		activity.WebhookActivity,
		param.TargetUrl,
		param.AuthHeaderName,
		param.AuthHeaderValue,
		param.RequestBody,
		param.NotifyFailure,
		param.Notification,
		param.NotificationProviderApiKey,
		param.NotifyAfterAttempts,
	).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
