package postmark

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"strings"

	"github.com/Boostport/mjml-go"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/mrz1836/postmark"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

const (
	PostmarkMessageStreamMagicLink = string(commonenum.WorkspaceProviderMagicLink)
	PostmarkMessageStreamInvoice   = "invoices"
)

const (
	WorkflowMagicLinkSubject = "One click away from CustomerOS"
	WorkflowMagicLink        = string(commonenum.WorkspaceProviderMagicLink)

	WorkflowInvoicePaid                  = "invoice-paid"
	WorkflowInvoicePaymentReceived       = "invoice-payment-received"
	WorkflowInvoicePaymentPending        = "invoice-payment-pending"
	WorkflowInvoiceReadyWithPaymentLink  = "invoice-ready"
	WorkflowInvoiceReadyNoPaymentLink    = "invoice-ready-nolink"
	WorkflowInvoiceVoided                = "invoice-voided"
	WorkflowFailedWebhook                = "failed-webhook"
	WorkflowInvoiceRemindWithPaymentLink = "invoice-remind"
	WorkflowInvoiceRemindNoPaymentLink   = "invoice-remind-nolink"

	WorkflowInvoiceVoidedSubject          = "Voided Invoice %s"
	WorkflowInvoicePaidSubject            = "Paid Invoice %s from %s"
	WorkflowInvoicePaymentReceivedSubject = "Payment Received for Invoice %s from %s"
	WorkflowInvoiceReadySubject           = "New invoice %s"
	WorkflowInvoiceRemindSubject          = "Follow-Up: Overdue Invoice %s"
)

type postmarkService struct {
	postmarkConfig *config.PostmarkConfig
	postgres       *postgres_repository.Repositories
}

func NewPostmarkService(postmarkConfig *config.PostmarkConfig, postgres *postgres_repository.Repositories) interfaces.PostmarkService {
	return &postmarkService{
		postmarkConfig: postmarkConfig,
		postgres:       postgres,
	}
}

func (s *postmarkService) getPostmarkClient(ctx context.Context, tenant string) (*postmark.Client, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PostmarkService.getPostmarkClient")
	defer spans.Finish()

	p := s.postgres.PostmarkApiKeyRepository.GetPostmarkApiKey(ctx, tenant)
	if p.Error != nil {
		spans.TraceError(p.Error)
		return nil, p.Error
	}

	if p.Result == nil {
		err := errors.New("postmark api key not found")
		spans.TraceError(err)
		return nil, err
	}

	serverToken := p.Result.(*postgres_entity.PostmarkApiKey).Key

	return postmark.NewClient(serverToken, ""), nil
}

func (s *postmarkService) SendNotification(ctx context.Context, postmarkEmail interfaces.PostmarkEmail, tenant string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PostmarkService.SendNotification")
	defer spans.Finish()
	spans.LogObjectAsJson("postmarkEmail", postmarkEmail)

	if postmarkEmail.From == "" {
		err := errors.New("missing from email address")
		spans.TraceError(err)
		return err
	}

	postmarkClient, err := s.getPostmarkClient(ctx, tenant)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	htmlContent, err := s.LoadEmailContent(ctx, postmarkEmail.WorkflowId, "mjml", postmarkEmail.TemplateData)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	htmlContent, err = s.convertMjmlToHtml(ctx, htmlContent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	textContent, err := s.LoadEmailContent(ctx, postmarkEmail.WorkflowId, "txt", postmarkEmail.TemplateData)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	email := postmark.Email{
		From:       postmarkEmail.From,
		To:         postmarkEmail.To,
		Cc:         strings.Join(postmarkEmail.CC, ","),
		Bcc:        strings.Join(postmarkEmail.BCC, ","),
		Subject:    postmarkEmail.Subject,
		TextBody:   textContent,
		HTMLBody:   htmlContent,
		TrackOpens: true,
	}

	if postmarkEmail.MessageStream != "" {
		email.MessageStream = postmarkEmail.MessageStream
	}

	if postmarkEmail.Attachments != nil {
		for _, attachment := range postmarkEmail.Attachments {
			email.Attachments = append(email.Attachments, postmark.Attachment{
				Name:        attachment.Filename,
				Content:     attachment.ContentEncoded,
				ContentType: attachment.ContentType,
				ContentID:   attachment.ContentID,
			})
		}
	}

	_, err = postmarkClient.SendEmail(ctx, email)
	if err != nil {
		wrappedError := fmt.Errorf("(postmarkClient.SendEmail) error: %s", err.Error())
		spans.TraceError(wrappedError)
		return err
	}

	return nil
}

func (s *postmarkService) LoadEmailContent(ctx context.Context, workflowId, fileExtension string, templateData map[string]string) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PostmarkService.LoadEmailContent")
	defer spans.Finish()

	rawEmailTemplate, err := s.LoadEmailBody(ctx, workflowId, fileExtension)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	emailTemplate := s.FillTemplate(rawEmailTemplate, templateData)

	return emailTemplate, nil
}

func (s *postmarkService) getFileName(workflowId, fileExtension string) string {
	var fileName string
	switch workflowId {
	case WorkflowMagicLink:
		fileName = "magic-link." + fileExtension
	case WorkflowInvoicePaid:
		fileName = "invoice.paid." + fileExtension
	case WorkflowInvoicePaymentReceived:
		fileName = "invoice.payment.received." + fileExtension
	case WorkflowInvoicePaymentPending:
		fileName = "invoice.payment.pending." + fileExtension
	case WorkflowInvoiceReadyWithPaymentLink:
		fileName = "invoice.ready." + fileExtension
	case WorkflowInvoiceReadyNoPaymentLink:
		fileName = "invoice.ready.nolink." + fileExtension
	case WorkflowInvoiceVoided:
		fileName = "invoice.voided." + fileExtension
	case WorkflowInvoiceRemindWithPaymentLink:
		fileName = "invoice.remind." + fileExtension
	case WorkflowInvoiceRemindNoPaymentLink:
		fileName = "invoice.remind.nolink." + fileExtension
	}
	return fileName
}

func (s *postmarkService) LoadEmailBody(ctx context.Context, workflowId, fileExtension string) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PostmarkService.LoadEmailBody")
	defer spans.Finish()

	fileName := s.getFileName(workflowId, fileExtension)
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

func (s *postmarkService) FillTemplate(template string, replace map[string]string) string {
	filledTemplate := template
	for k, v := range replace {
		filledTemplate = strings.Replace(filledTemplate, k, v, -1)
	}

	return filledTemplate
}

func (s *postmarkService) convertMjmlToHtml(ctx context.Context, filledTemplate string) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PostmarkService.convertMjmlToHtml")
	defer spans.Finish()

	html, err := mjml.ToHTML(ctx, filledTemplate)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		spans.TraceError(err)
		return "", fmt.Errorf("(PostmarkService.Template) error: %s", mjmlError.Message)
	}

	return html, err
}

func (s *postmarkService) CreateServerIfNotExists(ctx context.Context) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PostmarkService.CreateServerIfNotExists")
	defer spans.Finish()

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate postmark is configured
	if s.postmarkConfig == nil {
		err := errors.New("postmark not configured")
		spans.TraceError(err)
		return err
	}
	if s.postmarkConfig.Url == "" {
		err := errors.New("postmark url not configured")
		spans.TraceError(err)
		return err
	}
	if s.postmarkConfig.AccountApiKey == "" {
		err := errors.New("postmark api key not configured")
		spans.TraceError(err)
		return err
	}

	postmarkServerName := strings.ToLower(tenant)
	inboundWebhookURL := s.postmarkConfig.DefaultInboundStreamWebhook
	inboundForwardingDomain := postmarkServerName + ".customeros.ai"
	spans.LogKV("postmarkServerName", postmarkServerName)
	spans.LogKV("inboundWebhookURL", inboundWebhookURL)
	spans.LogKV("inboundForwardingDomain", inboundForwardingDomain)

	apiClient := NewPostmarkAPIClient(s.postmarkConfig.Url, s.postmarkConfig.AccountApiKey)

	// check if server already exists
	existingServerDetails, err := apiClient.GetServerByName(ctx, postmarkServerName)
	if err != nil {
		spans.TraceError(err)
	}
	if existingServerDetails != nil {
		spans.LogObjectAsJson("existingServerDetails", existingServerDetails)
	}

	if existingServerDetails == nil {
		serverResponse, err := apiClient.CreateServer(ctx, postmarkServerName, inboundWebhookURL, inboundForwardingDomain)
		if err != nil {
			spans.TraceError(err)
			return err
		}
		spans.LogObjectAsJson("createdServerResponse", serverResponse)
	}

	return nil
}

func (np *postmarkService) DeleteServer(ctx context.Context, tenant string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PostmarkService.DeleteServer")
	defer spans.Finish()

	// TODO implement me

	return nil
}
