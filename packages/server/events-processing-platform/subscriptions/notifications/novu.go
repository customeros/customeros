package notifications

import (
	"context"

	novu "github.com/novuhq/go-novu/lib"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type EmailableUser struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	SubscriberID string `json:"subscriberId"` // must be unique uuid for user
}

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

func (np *NovuProvider) SendEmail(ctx context.Context, u *EmailableUser, payload map[string]interface{}, eventId string) error {
	to := map[string]interface{}{
		"lastName":     u.LastName,
		"firstName":    u.FirstName,
		"subscriberId": u.SubscriberID,
		"email":        u.Email,
	}

	data := novu.ITriggerPayloadOptions{To: to, Payload: payload}

	_, err := np.NovuClient.EventApi.Trigger(ctx, eventId, data)

	if err != nil {
		np.log.Errorf("(NotificationsSubscriber.NovuProvider.SendEmail) error: %s", err.Error())
		return err
	}

	return nil
}
