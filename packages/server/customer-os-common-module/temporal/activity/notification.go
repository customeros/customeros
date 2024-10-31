package activity

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"

	"github.com/opentracing/opentracing-go"
)

func NotifyUserActivity(notification string, apiKey string) error {

	var n *commonService.NovuNotification

	// Unmarshal the JSON string into a map
	if err := json.Unmarshal([]byte(notification), &n); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}

	provider := commonService.NewNovuService(&commonService.Services{
		GlobalConfig: &config.GlobalConfig{
			NovuConfig: &config.NovuConfig{
				ApiKey: apiKey,
			},
		},
	})
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "CommonModule.Temporal.Activity.NotififyUserActivity")
	defer span.Finish()
	return provider.SendNotification(ctx, n)
}
