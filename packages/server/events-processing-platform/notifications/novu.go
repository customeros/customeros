package notifications

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Boostport/mjml-go"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	novu "github.com/novuhq/go-novu/lib"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
)

type NovuProvider struct {
	NovuClient *novu.APIClient
	log        logger.Logger
}

type NotifiableUser struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	SubscriberID string `json:"subscriberId"` // must be unique uuid for user
}

func newNovuProvider(log logger.Logger, apiKey string) *NovuProvider {
	return &NovuProvider{
		NovuClient: novu.NewAPIClient(apiKey, &novu.Config{}),
		log:        log,
	}
}

func (np *NovuProvider) SendNotification(ctx context.Context, notification *NovuNotification, span opentracing.Span) error {
	payload := notification.Payload
	u := notification.To
	workflowId := notification.WorkflowId

	rawEmailTemplate, err := np.LoadEmailBody(workflowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if rawEmailTemplate != "" {
		htmlEmailTemplate, err := np.FillTemplate(rawEmailTemplate, notification.TemplateData)
		if err != nil {
			tracing.TraceErr(span, err)
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
	}

	if fileName == "" {
		return "", nil
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String("eu-west-1")})
	if err != nil {
		return "", err
	}

	downloader := s3manager.NewDownloader(session)

	buffer := &aws.WriteAtBuffer{}
	_, err = downloader.Download(buffer,
		&s3.GetObjectInput{
			Bucket: aws.String("openline-production-mjml-templates"),
			Key:    aws.String(fileName),
		})
	if err != nil {
		return "", err
	}

	return string(buffer.Bytes()), nil
}

func (np *NovuProvider) FillTemplate(template string, replace map[string]string) (string, error) {
	_, ok := replace["{{userFirstName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing userFirstName")
	}
	_, ok = replace["{{actorFirstName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing actorFirstName")
	}
	_, ok = replace["{{actorLastName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing actorLastName")
	}
	_, ok = replace["{{orgName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing orgName")
	}
	_, ok = replace["{{orgLink}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing orgLink")
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

func containsKey[M ~map[K]V, K comparable, V any](m M, k K) bool {
	_, ok := m[k]
	return ok
}
