package opensrs

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
)

type openSRSService struct {
	log           logger.Logger
	openSrsConfig *config.OpenSRSConfig
	postgres      *postgres_repository.Repositories
	mailstack     interfaces.MailstackService
}

func NewOpenSRSService(log logger.Logger, openSrsConfig *config.OpenSRSConfig, postgres *postgres_repository.Repositories) interfaces.OpenSrsService {
	return &openSRSService{
		log:           log,
		openSrsConfig: openSrsConfig,
		postgres:      postgres,
	}
}

func (s *openSRSService) SetMailstackService(mailstack interfaces.MailstackService) {
	s.mailstack = mailstack
}

func (s *openSRSService) SendEmail(ctx context.Context, request *postgres_entity.EmailMessage) error {
	spans, _ := telemetry.StartServiceSpan(ctx, "OpenSrsService.Reply")
	defer spans.Finish()

	// Define the SMTP server details
	smtpHost := "mail.hostedemail.com"
	smtpPort := "587"

	tenant := common.GetTenantFromContext(ctx)

	mailboxRecord, err := s.mailstack.GetByMailbox(ctx, tenant, request.From)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if mailboxRecord == nil {
		err = errors.New("mailbox not found")
		spans.TraceError(err)
		return err
	}

	toEmail := []string{}
	ccEmail := []string{}
	bccEmail := []string{}

	for _, to := range request.To {
		toEmail = append(toEmail, to)
	}
	if request.Cc != nil {
		for _, cc := range request.Cc {
			ccEmail = append(ccEmail, cc)
		}
	}
	if request.Bcc != nil {
		for _, bcc := range request.Bcc {
			bccEmail = append(bccEmail, bcc)
		}
	}

	subject := request.Subject
	inReplyTo := request.ProviderInReplyTo
	references := request.ProviderReferences

	// Compose the email headers and body
	messageTemplate := `From: {{.FromEmail}}
To: {{.ToEmail}}{{if .CCEmail}}
Cc: {{.CCEmail}}{{- end}}
Subject: {{.Subject}}
Date: {{.Date}}
Message-ID: {{.MessageId}}
In-Reply-To: {{.InReplyTo}}
References: {{.References}}
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="{{.Boundary}}"

--{{.Boundary}}
Content-Type: text/plain; charset=US-ASCII; format=flowed

{{.PlainBody}}
--{{.Boundary}}
Content-Type: text/html; charset=UTF-8

{{.HTMLBody}}
--{{.Boundary}}--
`

	plainText, err := HTMLToPlainText(request.Content)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	data := struct {
		FromEmail  string
		ToEmail    string
		CCEmail    string
		BCCEmail   string
		Subject    string
		Date       string
		MessageId  string
		InReplyTo  string
		References string
		Boundary   string
		PlainBody  string
		HTMLBody   string
	}{
		ToEmail:    strings.Join(toEmail, ", "),
		CCEmail:    strings.Join(ccEmail, ", "),
		BCCEmail:   strings.Join(bccEmail, ", "),
		Subject:    subject,
		Date:       time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		MessageId:  generateMessageID(mailboxRecord.Email),
		InReplyTo:  inReplyTo,
		References: references,
		Boundary:   fmt.Sprintf("=_%x", time.Now().UnixNano()),
		PlainBody:  plainText,
		HTMLBody:   request.Content,
	}

	if request.FromName != "" {
		data.FromEmail = fmt.Sprintf("%s <%s>", request.FromName, request.From)
	} else {
		data.FromEmail = request.From
	}

	tmpl, err := template.New("email").Parse(messageTemplate)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	var msgBuffer bytes.Buffer
	if err := tmpl.Execute(&msgBuffer, data); err != nil {
		spans.TraceError(err)
		return err
	}

	msg := msgBuffer.String()

	// Combine all recipients: To, CC, BCC
	recipients := []string{}
	recipients = append(recipients, toEmail...)
	recipients = append(recipients, ccEmail...)
	recipients = append(recipients, bccEmail...)

	auth := smtp.PlainAuth("", mailboxRecord.Email, mailboxRecord.Password, smtpHost)

	// Send the email
	err = smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		auth,
		request.From,
		recipients,
		[]byte(msg),
	)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	request.ProviderMessageId = data.MessageId
	request.ProviderThreadId = data.MessageId
	request.ProviderInReplyTo = data.InReplyTo
	request.ProviderReferences = data.References

	return nil
}

func HTMLToPlainText(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	// Remove script and style elements
	doc.Find("script, style").Each(func(i int, el *goquery.Selection) {
		el.Remove()
	})

	// Get text content
	text := doc.Find("body").Text()

	// Trim spaces and replace multiple newlines with a single one
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n\n", "\n")

	return text, nil
}

func generateMessageID(fromEmail string) string {
	// Extract the mailbox part of the email address
	mailbox := fromEmail[:strings.IndexByte(fromEmail, '@')]

	// Generate a unique identifier using the mailbox and current timestamp
	uniqueID := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s.%d", mailbox, time.Now().UnixNano()))))

	// Construct the final Message-ID
	domain := fromEmail[strings.IndexByte(fromEmail, '@')+1:]
	messageID := fmt.Sprintf("<%s@%s>", uniqueID, domain)

	return messageID
}
