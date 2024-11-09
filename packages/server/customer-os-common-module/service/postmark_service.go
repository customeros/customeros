package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"
	"strings"

	"github.com/Boostport/mjml-go"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mrz1836/postmark"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
)

const (
	PostmarkMessageStreamInvoice = "invoices"
)

const (
	WorkflowIdTestFlow                   = "test-workflow"
	WorkflowInvoicePaid                  = "invoice-paid"
	WorkflowInvoicePaymentReceived       = "invoice-payment-received"
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

type PostmarkEmail struct {
	WorkflowId    string            `json:"workflowId"`
	MessageStream string            `json:"messageStream"`
	TemplateData  map[string]string `json:"templateData"`
	From          string            `json:"from"`
	To            string            `json:"to"`
	CC            []string          `json:"cc"`
	BCC           []string          `json:"bcc"`
	Subject       string            `json:"subject"`
	Attachments   []PostmarkEmailAttachment
}

type PostmarkEmailAttachment struct {
	Filename       string
	ContentEncoded string
	ContentType    string
	ContentID      string
}

type PostmarkService interface {
	SendNotification(ctx context.Context, postmarkEmail PostmarkEmail, tenant string) error
	CreateServer(ctx context.Context) error
	DeleteServer(ctx context.Context, tenant string) error
}

type postmarkService struct {
	services *Services
}

func NewPostmarkService(services *Services) PostmarkService {
	return &postmarkService{
		services: services,
	}
}

func (s *postmarkService) getPostmarkClient(ctx context.Context, tenant string) (*postmark.Client, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkService.getPostmarkClient")
	defer span.Finish()

	p := s.services.PostgresRepositories.PostmarkApiKeyRepository.GetPostmarkApiKey(ctx, tenant)
	if p.Error != nil {
		tracing.TraceErr(span, p.Error)
		return nil, p.Error
	}

	if p.Result == nil {
		err := errors.New("postmark api key not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	serverToken := p.Result.(*postgresEntity.PostmarkApiKey).Key

	return postmark.NewClient(serverToken, ""), nil
}

func (s *postmarkService) SendNotification(ctx context.Context, postmarkEmail PostmarkEmail, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkService.SendNotification")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "postmarkEmail", postmarkEmail)

	if postmarkEmail.From == "" {
		err := errors.New("missing from email address")
		tracing.TraceErr(span, err)
		return err
	}

	postmarkClient, err := s.getPostmarkClient(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	htmlContent, err := s.LoadEmailContent(ctx, postmarkEmail.WorkflowId, "mjml", postmarkEmail.TemplateData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	htmlContent, err = s.convertMjmlToHtml(ctx, htmlContent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	textContent, err := s.LoadEmailContent(ctx, postmarkEmail.WorkflowId, "txt", postmarkEmail.TemplateData)
	if err != nil {
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, wrappedError)
		return err
	}

	return nil
}

func (s *postmarkService) LoadEmailContent(ctx context.Context, workflowId, fileExtension string, templateData map[string]string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkService.LoadEmailContent")
	defer span.Finish()

	rawEmailTemplate, err := s.LoadEmailBody(ctx, workflowId, fileExtension)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	emailTemplate := s.FillTemplate(rawEmailTemplate, templateData)

	return emailTemplate, nil
}

func (s *postmarkService) getFileName(workflowId, fileExtension string) string {
	var fileName string
	switch workflowId {
	case WorkflowInvoicePaid:
		fileName = "invoice.paid." + fileExtension
	case WorkflowInvoicePaymentReceived:
		fileName = "invoice.payment.received." + fileExtension
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkService.LoadEmailBody")
	defer span.Finish()

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkService.convertMjmlToHtml")
	defer span.Finish()

	html, err := mjml.ToHTML(ctx, filledTemplate)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("(PostmarkService.Template) error: %s", mjmlError.Message)
	}

	return html, err
}

func (s *postmarkService) CreateServer(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkService.CreateServer")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// check postmark url and api key
	if s.services.GlobalConfig.ExternalServices.PostmarkConfig.Url == "" {
		err := errors.New("postmark url not configured")
		tracing.TraceErr(span, err)
		return err
	}
	if s.services.GlobalConfig.ExternalServices.PostmarkConfig.AccountApiKey == "" {
		err := errors.New("postmark api key not configured")
		tracing.TraceErr(span, err)
		return err
	}

	postmarkServerName := strings.ToLower(tenant)
	inboundWebhookURL := s.services.GlobalConfig.ExternalServices.PostmarkConfig.DefaultInboundStreamWebhook
	inboundForwardingDomain := postmarkServerName + ".customeros.ai"
	apiClient := NewPostmarkAPIClient(s.services.GlobalConfig.ExternalServices.PostmarkConfig.Url, s.services.GlobalConfig.ExternalServices.PostmarkConfig.AccountApiKey)

	serverResponse, err := apiClient.CreateServer(ctx, postmarkServerName, inboundWebhookURL, inboundForwardingDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tracing.LogObjectAsJson(span, "serverResponse", serverResponse)

	return nil
}

func (np *postmarkService) DeleteServer(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkService.DeleteServer")
	defer span.Finish()

	// TODO implement me

	return nil
}
