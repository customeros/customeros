package activity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/novu"

	"github.com/opentracing/opentracing-go"
)

func NotifyUserActivity(notification string, apiKey string) error {

	var n *interfaces.NovuNotification

	// Unmarshal the JSON string into a map
	if err := json.Unmarshal([]byte(notification), &n); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}

	provider := novu.NewNovuService(apiKey)
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "CommonModule.Temporal.Activity.NotififyUserActivity")
	defer span.Finish()
	return provider.SendNotification(ctx, n)
}
