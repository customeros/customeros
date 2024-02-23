# Usage example

### Webhook workflow
```go
func Example() error {
    hostPort := "localhost:7233" // pick this up from config
    namespace := "customer-os" // same here
    tClient := TemporalClient(hostPort, namespace)
    defer tClient.Close()
    // this is left to be implemented at this level so each execution can be configured
    workflowOptions := client.StartWorkflowOptions{
        ID:                       fmt.Sprintf("%s_%s", workflows.WEBHOOK_CALLS_TASK_QUEUE, uuid.New().String()),
        WorkflowExecutionTimeout: time.Hour * 24 * 365 * 10, // 10 years is a lot
        TaskQueue:                workflows.WEBHOOK_CALLS_TASK_QUEUE,
    }
    // this is left to be implemented at this level so each execution can be configured
    retryPolicy := &temporal.RetryPolicy{
        InitialInterval:        time.Second,
        BackoffCoefficient:     2.0,
        MaximumInterval:        time.Second * 100, // 100 * InitialInterval
        MaximumAttempts:        3,                 // if set to 0 means Unlimited attempts; not inclusive eg. n < 3.
        NonRetryableErrorTypes: []string{},        // empty
    }

    result, err := tClient.ExecuteWorkflow(
        context.Background(), 
        workflowOptions, 
        workflows.WebhookWorkflow, 
        workflows.WHWorkflowParam{
            TargetUrl:       "https://some.callback.url/events",
            RequestBody:     bytes.NewBuffer([]byte(`{"key":"value"}`)),
            AuthHeaderName:  "Authorization",
            AuthHeaderValue: "Bearer TOKEN",
            RetryPolicy: retryPolicy,
        })
    if err != nil {
        return err
    }
    fmt.Print(result)
    return nil
}
```