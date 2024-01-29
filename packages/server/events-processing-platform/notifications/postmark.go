package notifications

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Boostport/mjml-go"
	"github.com/mrz1836/postmark"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type PostmarkProvider struct {
	Client          *postmark.Client
	TemplatePath    string
	log             logger.Logger
	emailRawContent string
	emailContent    string
}

func newPostmarkProvider(log logger.Logger, serverToken, accountToken string) *PostmarkProvider {
	return &PostmarkProvider{
		Client: postmark.NewClient(serverToken, accountToken),
		log:    log,
	}
}

// SendNotification is a method that sends notification to the user, payload must contain the following:
// - from: email address
// - subject: email subject
// - attachment: base64 encoded PDF attachment
func (np *PostmarkProvider) SendNotification(ctx context.Context, u *NotifiableUser, payload map[string]interface{}, workflowId string) error {
	email := postmark.Email{
		From:       payload["from"].(string),
		To:         u.Email,
		Subject:    payload["subject"].(string),
		HTMLBody:   np.emailContent,
		TrackOpens: true,
	}

	encoded := payload["attachment"].(string)
	var attachmentName string
	switch workflowId {
	case WorkflowInvoicePaid:
		attachmentName = "paid_invoice.pdf"
	case WorkflowInvoiceReady:
		attachmentName = "ready_invoice.pdf"
	}

	attachment := []postmark.Attachment{
		{
			Name:        attachmentName,
			Content:     encoded,
			ContentType: "application/pdf",
		},
	}

	email.Attachments = attachment

	_, err := np.Client.SendEmail(ctx, email)

	if err != nil {
		np.log.Errorf("(PostmarkProvider.SendNotification) error: %s", err.Error())
		return err
	}

	return nil
}

func (np *PostmarkProvider) LoadEmailBody(ctx context.Context, workflowId string) error {
	var fileName string
	switch workflowId {
	case WorkflowInvoicePaid:
		fileName = "invoice.paid.mjml"
	case WorkflowInvoiceReady:
		fileName = "invoice.ready.mjml"
	}

	if _, err := os.Stat(np.TemplatePath); os.IsNotExist(err) {
		return fmt.Errorf("(PostmarkProvider.LoadEmailBody) error: %s", err.Error())
	}
	emailPath := fmt.Sprintf("%s/%s", np.TemplatePath, fileName)
	if _, err := os.Stat(emailPath); err != nil {
		return fmt.Errorf("(PostmarkProvider.LoadEmailBody) error: %s", err.Error())
	}

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("(PostmarkProvider.LoadEmailBody) Getwd: %s", err.Error())
	}

	// Build the full path to the template file
	templatePath := filepath.Join(currentDir, emailPath)
	if _, err := os.Stat(templatePath); err != nil {
		return fmt.Errorf("(NovuProvider.LoadEmailBody) error: %s", err.Error())
	}

	rawMjml, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("(PostmarkProvider.LoadEmailBody) error: %s", err.Error())
	}
	np.emailRawContent = string(rawMjml[:])
	return nil
}

func (np *PostmarkProvider) Template(ctx context.Context, replace map[string]string) (string, error) {
	mjmlf := np.emailRawContent
	for k, v := range replace {
		mjmlf = strings.Replace(mjmlf, k, v, -1)
	}
	np.emailRawContent = mjmlf

	html, err := mjml.ToHTML(context.Background(), mjmlf)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		return "", fmt.Errorf("(PostmarkProvider.Template) error: %s", mjmlError.Message)
	}
	np.emailContent = html
	return html, err
}

func (np *PostmarkProvider) GetRawContent() string {
	return np.emailRawContent
}
