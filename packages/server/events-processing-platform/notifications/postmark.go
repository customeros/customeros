package notifications

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"strings"

	"github.com/Boostport/mjml-go"
	"github.com/mrz1836/postmark"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type PostmarkEmail struct {
	WorkflowId   string
	TemplateData map[string]string

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
}

type PostmarkProvider struct {
	Client *postmark.Client
	log    logger.Logger
}

func NewPostmarkProvider(log logger.Logger, serverToken string) *PostmarkProvider {
	return &PostmarkProvider{
		Client: postmark.NewClient(serverToken, ""),
		log:    log,
	}
}

func (np *PostmarkProvider) SendNotification(ctx context.Context, postmarkEmail PostmarkEmail, span *opentracing.Span) error {
	rawEmailTemplate, err := np.LoadEmailBody(postmarkEmail.WorkflowId)
	if err != nil {
		tracing.TraceErr(*span, err)
		return err
	}

	htmlEmailTemplate, err := np.FillTemplate(rawEmailTemplate, postmarkEmail.TemplateData)
	if err != nil {
		tracing.TraceErr(*span, err)
		return err
	}

	email := postmark.Email{
		From:       postmarkEmail.From,
		To:         postmarkEmail.To,
		Cc:         postmarkEmail.CC,
		Bcc:        postmarkEmail.BCC,
		Subject:    postmarkEmail.Subject,
		HTMLBody:   htmlEmailTemplate,
		TrackOpens: true,
	}

	if postmarkEmail.Attachments != nil {
		for _, attachment := range postmarkEmail.Attachments {
			email.Attachments = append(email.Attachments, postmark.Attachment{
				Name:        attachment.Filename,
				Content:     attachment.ContentEncoded,
				ContentType: attachment.ContentType,
			})
		}
	}

	_, err = np.Client.SendEmail(ctx, email)

	if err != nil {
		tracing.TraceErr(*span, err)
		np.log.Errorf("(PostmarkProvider.SendNotification) error: %s", err.Error())
		return err
	}

	return nil
}

func (np *PostmarkProvider) LoadEmailBody(workflowId string) (string, error) {
	var fileName string
	switch workflowId {
	case WorkflowInvoicePaid:
		fileName = "invoice.paid.mjml"
	case WorkflowInvoiceReady:
		fileName = "invoice.ready.mjml"
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

func (np *PostmarkProvider) FillTemplate(template string, replace map[string]string) (string, error) {
	filledTemplate := template
	for k, v := range replace {
		filledTemplate = strings.Replace(filledTemplate, k, v, -1)
	}

	html, err := mjml.ToHTML(context.Background(), filledTemplate)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		return "", fmt.Errorf("(PostmarkProvider.Template) error: %s", mjmlError.Message)
	}

	return html, err
}
