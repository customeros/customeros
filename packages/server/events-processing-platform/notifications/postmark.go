package notifications

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
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

func (np *PostmarkProvider) SendNotification(ctx context.Context, u *NotifiableUser, payload map[string]interface{}, workflowId string) error {
	email := postmark.Email{
		From:       "no-reply@example.com",
		To:         u.Email,
		Subject:    payload["subject"].(string),
		HTMLBody:   np.emailContent,
		TrackOpens: true,
	}

	attachment, err := np.readAttachment(workflowId)
	if err != nil {
		return err
	}

	email.Attachments = attachment

	_, err = np.Client.SendEmail(ctx, email)

	if err != nil {
		np.log.Errorf("(PostmarkProvider.SendNotification) error: %s", err.Error())
		return err
	}

	return nil
}

func (np *PostmarkProvider) LoadEmailBody(ctx context.Context, workflowId string) error {
	switch workflowId {
	case WorkflowIdOrgOwnerUpdateEmail:
		if _, err := os.Stat(np.TemplatePath); os.IsNotExist(err) {
			return fmt.Errorf("(PostmarkProvider.LoadEmailBody) error: %s", err.Error())
		}
		emailPath := fmt.Sprintf("%s/ownership.single.mjml", np.TemplatePath)
		if _, err := os.Stat(emailPath); err != nil {
			return fmt.Errorf("(PostmarkProvider.LoadEmailBody) error: %s", err.Error())
		}

		rawMjml, err := os.ReadFile(emailPath)
		if err != nil {
			return fmt.Errorf("(PostmarkProvider.LoadEmailBody) error: %s", err.Error())
		}
		np.emailRawContent = string(rawMjml[:])
	}
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

func (np *PostmarkProvider) readAttachment(workflowId string) ([]postmark.Attachment, error) {
	var readFile []byte
	var err error
	switch workflowId {
	case WorkflowInvoicePaid:
		readFile, err = os.ReadFile("/Users/efirut/work/openline/openline-customer-os/packages/server/events-processing-platform/test_app/sample.pdf")
		if err != nil {
			return nil, err
		}
		encoded := base64.StdEncoding.EncodeToString(readFile)
		return []postmark.Attachment{
			{
				Name:        "dummy.pdf",
				Content:     encoded,
				ContentType: "application/pdf",
			},
		}, nil
	}

	// convert to base64

	return nil, nil
}

// client := postmark.NewClient("TODO USE HERE STUFF", "")

// email := postmark.Email{
// 	From:       "hello@customeros.ai",
// 	To:         "edi@openline.ai",
// 	Subject:    "This is a test email sent through Postmark from Go!",
// 	HTMLBody:   "<p>Ce face fetele? Se da cu barca?</p>",
// 	TrackOpens: true,
// }

// readFile, err := os.ReadFile("/Users/efirut/work/openline/openline-customer-os/packages/server/events-processing-platform/test_app/sample.pdf")
// if err != nil {
// 	panic(err)
// }

// //convert to base64
// encoded := base64.StdEncoding.EncodeToString(readFile)

// email.Attachments = []postmark.Attachment{
// 	{
// 		Name:        "dummy.pdf",
// 		Content:     encoded,
// 		ContentType: "application/pdf",
// 	},
// }

// _, err = client.SendEmail(context.Background(), email)
// if err != nil {
// 	panic(err)
// }
