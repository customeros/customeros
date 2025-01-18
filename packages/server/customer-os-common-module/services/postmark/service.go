package postmark

import (
	"context"
	"fmt"
	"strings"

	"github.com/Boostport/mjml-go"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/mrz1836/postmark"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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
	postgres       *repository.Repositories
}

func NewPostmarkService(postmarkConfig *config.PostmarkConfig, postgres *repository.Repositories) interfaces.PostmarkService {
	return &postmarkService{
		postmarkConfig: postmarkConfig,
		postgres:       postgres,
	}
}

func (s *postmarkService) getPostmarkClient(ctx context.Context, tenant string) (*postmark.Client, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkService.getPostmarkClient")
	defer span.Finish()

	p := s.postgres.PostmarkApiKeyRepository.GetPostmarkApiKey(ctx, tenant)
	if p.Error != nil {
		tracing.TraceErr(span, p.Error)
		return nil, p.Error
	}

	if p.Result == nil {
		err := errors.New("postmark api key not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	serverToken := p.Result.(*entity.PostmarkApiKey).Key

	return postmark.NewClient(serverToken, ""), nil
}

func (s *postmarkService) SendNotification(ctx context.Context, postmarkEmail interfaces.PostmarkEmail, tenant string) error {
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
	case WorkflowMagicLink:
		fileName = "magic-link." + fileExtension
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

func (s *postmarkService) CreateServerIfNotExists(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkService.CreateServerIfNotExists")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate postmark is configured
	if s.postmarkConfig == nil || s.postmarkConfig.Url == "" || s.postmarkConfig.AccountApiKey == "" {
		err := errors.New("postmark url not configured")
		tracing.TraceErr(span, err)
		return err
	}

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// check postmark url and api key
	if s.postmarkConfig.Url == "" {
		err := errors.New("postmark url not configured")
		tracing.TraceErr(span, err)
		return err
	}
	if s.postmarkConfig.AccountApiKey == "" {
		err := errors.New("postmark api key not configured")
		tracing.TraceErr(span, err)
		return err
	}

	postmarkServerName := strings.ToLower(tenant)
	inboundWebhookURL := s.postmarkConfig.DefaultInboundStreamWebhook
	inboundForwardingDomain := postmarkServerName + ".customeros.ai"
	span.LogKV("postmarkServerName", postmarkServerName)
	span.LogKV("inboundWebhookURL", inboundWebhookURL)
	span.LogKV("inboundForwardingDomain", inboundForwardingDomain)

	apiClient := NewPostmarkAPIClient(s.postmarkConfig.Url, s.postmarkConfig.AccountApiKey)

	// check if server already exists
	existingServerDetails, err := apiClient.GetServerByName(ctx, postmarkServerName)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	if existingServerDetails != nil {
		tracing.LogObjectAsJson(span, "existingServerDetails", existingServerDetails)
	}

	if existingServerDetails == nil {
		serverResponse, err := apiClient.CreateServer(ctx, postmarkServerName, inboundWebhookURL, inboundForwardingDomain)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		tracing.LogObjectAsJson(span, "createdServerResponse", serverResponse)
	}

	return nil
}

func (np *postmarkService) DeleteServer(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkService.DeleteServer")
	defer span.Finish()

	// TODO implement me

	return nil
}
