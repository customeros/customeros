package activity

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/aws_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/notifications"
)

func NotifyUserActivity(notification *notifications.NovuNotification, apiKey string) error {
	s3 := aws_client.NewS3Client(&aws.Config{Region: aws.String("eu-west-1")})
	log := logger.NewAppLogger(&logger.Config{LogLevel: "info"})
	provider := notifications.NewNovuNotificationProvider(log, apiKey, s3)
	return provider.SendNotification(context.Background(), notification)
}
