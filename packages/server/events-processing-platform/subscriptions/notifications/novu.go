package notifications

import (
	"context"

	novu "github.com/novuhq/go-novu/lib"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type NovuProvider struct {
	NovuClient *novu.APIClient
	log        logger.Logger
}

func NewNovuProvider(log logger.Logger, apiKey string) *NovuProvider {
	return &NovuProvider{
		NovuClient: novu.NewAPIClient(apiKey, &novu.Config{}),
		log:        log,
	}
}

func (np *NovuProvider) SendNotification(ctx context.Context, u *NotifiableUser, payload, overrides map[string]interface{}, workflowId string) error {
	to := map[string]interface{}{
		"lastName":     u.LastName,
		"firstName":    u.FirstName,
		"subscriberId": u.SubscriberID,
		"email":        u.Email,
	}
	data := novu.ITriggerPayloadOptions{To: to, Payload: payload}
	if len(overrides) > 0 {
		if containsKey(overrides, "email") {
			data.Overrides = overrides
		}
	}

	_, err := np.NovuClient.EventApi.Trigger(ctx, workflowId, data)

	if err != nil {
		np.log.Errorf("(NotificationsSubscriber.NovuProvider.SendNotification) error: %s", err.Error())
		return err
	}

	return nil
}

func containsKey[M ~map[K]V, K comparable, V any](m M, k K) bool {
	_, ok := m[k]
	return ok
}
