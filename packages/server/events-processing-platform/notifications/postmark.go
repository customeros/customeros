package notifications

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"

	"github.com/Boostport/mjml-go"
	"github.com/mrz1836/postmark"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
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
	CC      string
	BCC     string
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

func (np *PostmarkProvider) getPostmarkClient(tenant string) (*postmark.Client, error) {
	p := np.repository.CommonRepositories.PostmarkApiKeyRepository.GetPostmarkApiKey(tenant)
	if p.Error != nil {
		return nil, p.Error
	}

	if p.Result == nil {
		return nil, errors.New("postmark api key not found")
	}

	serverToken := p.Result.(entity.PostmarkApiKey).Key

	return postmark.NewClient(serverToken, ""), nil
}

func (np *PostmarkProvider) SendNotification(ctx context.Context, postmarkEmail PostmarkEmail, span opentracing.Span, tenant string) error {
	postmarkClient, err := np.getPostmarkClient(tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		np.log.Errorf("(PostmarkProvider.SendNotification) error: %s", err.Error())
		return err
	}

	htmlContent, err := np.LoadEmailContent(postmarkEmail.WorkflowId, "mjml", postmarkEmail.TemplateData)
	if err != nil {
		tracing.TraceErr(span, err)
		np.log.Errorf("(PostmarkProvider.SendNotification.LoadEmailContent) error: %s", err.Error())
		return err
	}
	htmlContent, err = np.ConvertMjmlToHtml(htmlContent)
	if err != nil {
		tracing.TraceErr(span, err)
		np.log.Errorf("(PostmarkProvider.SendNotification. ConvertMjmlToHtml) error: %s", err.Error())
		return err
	}

	textContent, err := np.LoadEmailContent(postmarkEmail.WorkflowId, "txt", postmarkEmail.TemplateData)
	if err != nil {
		tracing.TraceErr(span, err)
		np.log.Errorf("(PostmarkProvider.SendNotification.LoadEmailContent) error: %s", err.Error())
		return err
	}

	email := postmark.Email{
		From:       postmarkEmail.From,
		To:         postmarkEmail.To,
		Cc:         postmarkEmail.CC,
		Bcc:        postmarkEmail.BCC,
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
		tracing.TraceErr(span, err)
		np.log.Errorf("(PostmarkProvider.SendNotification) error: %s", err.Error())
		return err
	}

	return nil
}

func (np *PostmarkProvider) LoadEmailContent(workflowId, fileExtension string, templateData map[string]string) (string, error) {
	rawEmailTemplate, err := np.LoadEmailBody(workflowId, fileExtension)
	if err != nil {
		return "", err
	}

	emailTemplate := np.FillTemplate(rawEmailTemplate, templateData)
	if err != nil {
		return "", err
	}

	return emailTemplate, nil
}

func (np *PostmarkProvider) GetFileName(workflowId, fileExtension string) string {
	var fileName string
	switch workflowId {
	case WorkflowInvoicePaid:
		fileName = "invoice.paid." + fileExtension
	case WorkflowInvoiceReady:
		fileName = "invoice.ready." + fileExtension
	case WorkflowInvoiceVoided:
		fileName = "invoice.voided." + fileExtension
	}
	return fileName
}

func (np *PostmarkProvider) LoadEmailBody(workflowId, fileExtension string) (string, error) {
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

func (np *PostmarkProvider) ConvertMjmlToHtml(filledTemplate string) (string, error) {
	html, err := mjml.ToHTML(context.Background(), filledTemplate)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		return "", fmt.Errorf("(PostmarkProvider.Template) error: %s", mjmlError.Message)
	}

	return html, err
}
