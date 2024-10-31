package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"strings"

	"github.com/Boostport/mjml-go"
	novu "github.com/novuhq/go-novu/lib"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

const (
	//TODO rework all of this
	WorkflowIdOrgOwnerUpdateEmail                  = "org-owner-update-email"
	WorkflowIdOrgOwnerUpdateAppNotification        = "org-owner-update-in-app-notification"
	WorkflowIdOrgOwnerUpdateEmailSubject           = "%s %s added you as an owner"
	WorkflowIdOrgOwnerUpdateAppNotificationSubject = "%s %s added you as an owner"

	WorkflowFailedWebhookSubject = "[Action Required] Webhook %s is offline"

	WorkflowReminderNotificationSubject = "Reminder, %s"
	WorkflowReminderNotificationEmail   = "reminder-notification-email"
	WorkflowReminderInAppNotification   = "reminder-in-app-notification"

	WorkflowId_FlowParticipantGoalAchievedEmail = "flow-participant-goal-achieved-email"
)

var REQUIRED_TEMPLATE_VALUES = map[string][]string{
	WorkflowIdOrgOwnerUpdateEmail: {
		"{{userFirstName}}",
		"{{actorFirstName}}",
		"{{actorLastName}}",
		"{{orgName}}",
		"{{orgLink}}",
	},
	WorkflowFailedWebhook: {
		"{{userFirstName}}",
		"{{webhookName}}",
		"{{webhookUrl}}",
	},
	WorkflowReminderNotificationEmail: {
		"{{reminderContent}}",
		"{{reminderCreatedAt}}",
		"{{orgName}}",
		"{{orgLink}}",
	},
}

type NotifiableUser struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	SubscriberID string `json:"subscriberId"` // must be unique uuid for user
}

type NovuNotification struct {
	WorkflowId   string
	TemplateData map[string]string

	To      *NotifiableUser
	Subject string
	Payload map[string]interface{}
}

type NovuService interface {
	SendNotification(ctx context.Context, notification *NovuNotification) error
}

type novuService struct {
	services   *Services
	NovuClient *novu.APIClient
}

func NewNovuService(services *Services) NovuService {
	apiKey := ""

	if services.GlobalConfig.NovuConfig != nil {
		apiKey = services.GlobalConfig.NovuConfig.ApiKey
	}

	return &novuService{
		NovuClient: novu.NewAPIClient(apiKey, &novu.Config{}),
	}
}

func (np *novuService) SendNotification(ctx context.Context, notification *NovuNotification) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NovuService.SendNotification")
	defer span.Finish()

	payload := notification.Payload

	if payload == nil {
		payload = make(map[string]interface{})
		payload["html"] = ""
	}

	payload["subject"] = notification.Subject

	u := notification.To
	workflowId := notification.WorkflowId

	rawEmailTemplate, err := np.LoadEmailBody(ctx, workflowId, "mjml")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if rawEmailTemplate != "" {
		htmlEmailTemplate, err := np.FillTemplate(workflowId, rawEmailTemplate, notification.TemplateData)
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
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (np *novuService) LoadEmailBody(ctx context.Context, workflowId, fileExtension string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NovuService.LoadEmailBody")
	defer span.Finish()

	fileName := np.GetFileName(workflowId, fileExtension)
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

func (np *novuService) GetFileName(workflowId, fileExtension string) string {
	var fileName string
	switch workflowId {
	case WorkflowId_FlowParticipantGoalAchievedEmail:
		fileName = "novu-flow-participant-goal-achieved-email." + fileExtension
	case WorkflowIdOrgOwnerUpdateEmail:
		fileName = "ownership.single." + fileExtension
	case WorkflowFailedWebhook:
		fileName = "webhook.failed." + fileExtension
	case WorkflowReminderNotificationEmail:
		fileName = "reminder." + fileExtension
	}
	return fileName
}

func (np *novuService) FillTemplate(workflowId, template string, replace map[string]string) (string, error) {
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
		return "", fmt.Errorf("(novuService.FillTemplate) error: %s", mjmlError.Message)
	}
	return html, err
}

func checkRequiredTemplateVars(replace map[string]string, requiredVars []string) error {
	for _, rv := range requiredVars {
		if _, ok := replace[rv]; !ok {
			return fmt.Errorf("(novuService.FillTemplate) error: missing %s", rv)
		}
	}

	return nil
}

func containsKey[M ~map[K]V, K comparable, V any](m M, k K) bool {
	_, ok := m[k]
	return ok
}
