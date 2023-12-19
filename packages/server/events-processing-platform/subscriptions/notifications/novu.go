package notifications

import (
	"context"
	"os"
	"strings"

	novu "github.com/novuhq/go-novu/lib"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

// album represents data about a record album.
type EmailableUser struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Message      string `json:"message"`
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

func (np *NovuProvider) SendEmail(ctx context.Context, u *EmailableUser, eventId string) error {
	to := map[string]interface{}{
		"lastName":     u.LastName,
		"firstName":    u.FirstName,
		"subscriberId": u.SubscriberID,
		"email":        u.Email,
	}

	msg := u.Message
	var html string
	var button map[string]string
	switch eventId {
	case "test_flow":
		rawHtml, _ := os.ReadFile("./email_templates/email2.html") // FIXME: replace this html with an actual template
		msg = "Welcome to CustomerOS!"
		html = strings.Replace(string(rawHtml[:]), "{{fName}}", u.FirstName, -1)
		html = strings.Replace(html, "{{lName}}", u.LastName, -1)
		html = strings.Replace(html, "{{message}}", msg, -1)
		button = map[string]string{
			"text": "Click me",
			"url":  "https://customeros.ai",
		}
	case "user_update":
		// TODO: do something
		html = ""
		msg = ""
		button = map[string]string{
			"text": "Click me",
			"url":  "https://customeros.ai",
		}
	default:
		html = ""
		msg = ""
	}

	payload := map[string]interface{}{
		"message": msg,
		"organization": map[string]interface{}{
			"logo": "https://happycorp.com/logo.png", // able to add tenant logo here
		},
		"subscriber": map[string]interface{}{
			"firstName": u.FirstName,
			"lastName":  u.LastName,
			"email":     u.Email,
		},
		"button": button,
		"html":   string(html[:]),
	}

	data := novu.ITriggerPayloadOptions{To: to, Payload: payload}

	_, err := np.NovuClient.EventApi.Trigger(ctx, eventId, data)

	if err != nil {
		np.log.Errorf("(NotificationsSubscriber.NovuProvider.SendEmail) error: %s", err.Error())
		return err
	}

	return nil
}
