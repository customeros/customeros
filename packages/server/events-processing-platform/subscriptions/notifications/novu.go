package notifications

import (
	"context"
	"os"
	"strings"

	novu "github.com/novuhq/go-novu/lib"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
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

	from := payload["actor"].(*entity.UserEntity)

	var html string
	switch eventId {
	case EventIdTestFlow:
		rawHtml, _ := os.ReadFile("./email_templates/ownership.single.mjml")
		html = strings.Replace(string(rawHtml[:]), "{{fName}}", u.FirstName, -1)
		html = strings.Replace(html, "{{lName}}", u.LastName, -1)
	case EventIdOrgOwnerUpdateEmail:
		rawMjml, _ := os.ReadFile("./email_templates/ownership.single.mjml")
		mjmlf := strings.Replace(string(rawMjml[:]), "{{userFirstName}}", u.FirstName, -1)
		mjmlf = strings.Replace(mjmlf, "{{actorFirstName}}", from.FirstName, -1)
		mjmlf = strings.Replace(mjmlf, "{{actorLastName}}", from.LastName, -1)
		mjmlf = strings.Replace(mjmlf, "{{orgName}}", payload["orgName"].(string), -1)
		html = mjmlf // TODO: convert mjml to html
	default:
		html = ""
	}

	novuPayload := map[string]interface{}{
		"html":    html,
		"subject": payload["subject"],
		"email":   u.Email,
	}

	data := novu.ITriggerPayloadOptions{To: to, Payload: novuPayload}

	_, err := np.NovuClient.EventApi.Trigger(ctx, eventId, data)

	if err != nil {
		np.log.Errorf("(NotificationsSubscriber.NovuProvider.SendEmail) error: %s", err.Error())
		return err
	}

	return nil
}
