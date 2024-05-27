package notifications

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/Boostport/mjml-go"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mrz1836/postmark"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/notifications"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/opentracing/opentracing-go"
)

const (
	PostmarkMessageStreamInvoice = "invoices"
)

type PostmarkEmail struct {
	WorkflowId    string
	MessageStream string
	TemplateData  map[string]string

	From    string
	To      string
	CC      []string
	BCC     []string
	Subject string

	Attachments []PostmarkEmailAttachment
}

type PostmarkEmailAttachment struct {
	Filename       string
	ContentEncoded string
	ContentType    string
	ContentID      string
}

type PostmarkProvider struct {
	log        logger.Logger
	repository *repository.Repositories
}

func NewPostmarkProvider(log logger.Logger, repo *repository.Repositories) *PostmarkProvider {
	return &PostmarkProvider{
		log:        log,
		repository: repo,
	}
}

func (np *PostmarkProvider) getPostmarkClient(ctx context.Context, tenant string) (*postmark.Client, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkProvider.getPostmarkClient")
	defer span.Finish()

	p := np.repository.PostgresRepositories.PostmarkApiKeyRepository.GetPostmarkApiKey(tenant)
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

func (np *PostmarkProvider) SendNotification(ctx context.Context, postmarkEmail PostmarkEmail, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkProvider.SendNotification")
	defer span.Finish()

	postmarkClient, err := np.getPostmarkClient(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		np.log.Errorf("(PostmarkProvider.SendNotification) error: %s", err.Error())
		return err
	}

	htmlContent, err := np.LoadEmailContent(ctx, postmarkEmail.WorkflowId, "mjml", postmarkEmail.TemplateData)
	if err != nil {
		tracing.TraceErr(span, err)
		np.log.Errorf("(PostmarkProvider.SendNotification.LoadEmailContent) error: %s", err.Error())
		return err
	}
	htmlContent, err = np.ConvertMjmlToHtml(ctx, htmlContent)
	if err != nil {
		tracing.TraceErr(span, err)
		np.log.Errorf("(PostmarkProvider.SendNotification. ConvertMjmlToHtml) error: %s", err.Error())
		return err
	}

	textContent, err := np.LoadEmailContent(ctx, postmarkEmail.WorkflowId, "txt", postmarkEmail.TemplateData)
	if err != nil {
		tracing.TraceErr(span, err)
		np.log.Errorf("(PostmarkProvider.SendNotification.LoadEmailContent) error: %s", err.Error())
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
		np.log.Errorf("(PostmarkProvider.SendNotification) error: %s", err.Error())
		return err
	}

	return nil
}

func (np *PostmarkProvider) LoadEmailContent(ctx context.Context, workflowId, fileExtension string, templateData map[string]string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkProvider.LoadEmailContent")
	defer span.Finish()

	rawEmailTemplate, err := np.LoadEmailBody(ctx, workflowId, fileExtension)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	emailTemplate := np.FillTemplate(rawEmailTemplate, templateData)

	return emailTemplate, nil
}

func (np *PostmarkProvider) GetFileName(workflowId, fileExtension string) string {
	var fileName string
	switch workflowId {
	case notifications.WorkflowInvoicePaid:
		fileName = "invoice.paid." + fileExtension
	case notifications.WorkflowInvoicePaymentReceived:
		fileName = "invoice.payment.received." + fileExtension
	case notifications.WorkflowInvoiceReadyWithPaymentLink:
		fileName = "invoice.ready." + fileExtension
	case notifications.WorkflowInvoiceReadyNoPaymentLink:
		fileName = "invoice.ready.nolink." + fileExtension
	case notifications.WorkflowInvoiceVoided:
		fileName = "invoice.voided." + fileExtension
	}
	return fileName
}

func (np *PostmarkProvider) LoadEmailBody(ctx context.Context, workflowId, fileExtension string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkProvider.LoadEmailBody")
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

func (np *PostmarkProvider) FillTemplate(template string, replace map[string]string) string {
	filledTemplate := template
	for k, v := range replace {
		filledTemplate = strings.Replace(filledTemplate, k, v, -1)
	}

	return filledTemplate
}

func (np *PostmarkProvider) ConvertMjmlToHtml(ctx context.Context, filledTemplate string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostmarkProvider.ConvertMjmlToHtml")
	defer span.Finish()

	html, err := mjml.ToHTML(context.Background(), filledTemplate)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("(PostmarkProvider.Template) error: %s", mjmlError.Message)
	}

	return html, err
}
