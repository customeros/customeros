package activity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/aws_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/notifications"
	"github.com/opentracing/opentracing-go"
)

func NotifyUserActivity(notification string, apiKey string) error {

	var n *notifications.NovuNotification

	// Unmarshal the JSON string into a map
	if err := json.Unmarshal([]byte(notification), &n); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}

	s3 := aws_client.NewS3Client(&aws.Config{Region: aws.String("eu-west-1")})
	log := logger.NewAppLogger(&logger.Config{LogLevel: "info"})
	provider := notifications.NewNovuNotificationProvider(log, apiKey, s3)
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "CommonModule.Temporal.Activity.NotififyUserActivity")
	defer span.Finish()
	return provider.SendNotification(ctx, n, span)
}
