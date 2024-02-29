package notifications

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Boostport/mjml-go"
	novu "github.com/novuhq/go-novu/lib"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/aws_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

type NovuProvider struct {
	NovuClient *novu.APIClient
	log        logger.Logger
	s3Client   aws_client.S3ClientI
}

func newNovuProvider(log logger.Logger, apiKey string, s3 aws_client.S3ClientI) *NovuProvider {
	return &NovuProvider{
		NovuClient: novu.NewAPIClient(apiKey, &novu.Config{}),
		log:        log,
		s3Client:   s3,
	}
}

func (np *NovuProvider) SendNotification(ctx context.Context, notification *NovuNotification) error {
	payload := notification.Payload
	u := notification.To
	workflowId := notification.WorkflowId

	rawEmailTemplate, err := np.LoadEmailBody(workflowId)
	if err != nil {
		return err
	}

	if rawEmailTemplate != "" {
		htmlEmailTemplate, err := np.FillTemplate(workflowId, rawEmailTemplate, notification.TemplateData)
		if err != nil {
			return err
		}
		payload["html"] = htmlEmailTemplate
	}
	to := map[string]interface{}{
		"lastName":     u.LastName,
		"firstName":    u.FirstName,
		"subscriberId": u.SubscriberID,
		"email":        u.Email,
	}
	data := novu.ITriggerPayloadOptions{To: to, Payload: payload}
	if containsKey(payload, "overrides") {
		overrides := payload["overrides"].(map[string]interface{})
		if len(overrides) > 0 {
			if containsKey(overrides, "email") {
				data.Overrides = overrides
			}
		}
	}

	_, err = np.NovuClient.EventApi.Trigger(ctx, workflowId, data)

	if err != nil {
		np.log.Errorf("(NotificationsSubscriber.NovuProvider.SendNotification) error: %s", err.Error())
		return err
	}

	return nil
}

func (np *NovuProvider) LoadEmailBody(workflowId string) (string, error) {
	var fileName string
	switch workflowId {
	case WorkflowIdOrgOwnerUpdateEmail:
		fileName = "ownership.single.mjml"
	case WorkflowFailedWebhook:
		fileName = "webhook.failed.mjml"
	}

	if fileName == "" {
		return "", nil
	}

	np.s3Client.ChangeRegion("eu-west-1")
	return np.s3Client.Download("openline-production-mjml-templates", fileName)
}

func (np *NovuProvider) FillTemplate(workflowId, template string, replace map[string]string) (string, error) {
	requiredVars := REQUIRED_TEMPLATE_VALUES[workflowId]
	err := checkRequiredTemplateVars(replace, requiredVars)
	if err != nil {
		return "", err
	}
	mjmlf := template
	for k, v := range replace {
		mjmlf = strings.Replace(mjmlf, k, v, -1)
	}

	html, err := mjml.ToHTML(context.Background(), mjmlf)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		return "", fmt.Errorf("(NovuProvider.FillTemplate) error: %s", mjmlError.Message)
	}
	return html, err
}

func checkRequiredTemplateVars(replace map[string]string, requiredVars []string) error {
	for _, rv := range requiredVars {
		if _, ok := replace[rv]; !ok {
			return fmt.Errorf("(NovuProvider.FillTemplate) error: missing %s", rv)
		}
	}

	return nil
}

func containsKey[M ~map[K]V, K comparable, V any](m M, k K) bool {
	_, ok := m[k]
	return ok
}
