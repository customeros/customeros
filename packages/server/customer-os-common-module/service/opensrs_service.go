package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/DusanKasan/parsemail"
	"github.com/PuerkitoBio/goquery"
	mimemail "github.com/emersion/go-message/mail"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"net/smtp"
	"strings"
	"text/template"
	"time"
)

type openSRSService struct {
	services *Services
}

type OpenSrsService interface {
	Reply(ctx context.Context, tenant string, request dto.MailRequest) (*parsemail.Email, error)
}

func (s *openSRSService) Reply(ctx context.Context, tenant string, request dto.MailRequest) (*parsemail.Email, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpenSrsService.Reply")
	defer span.Finish()

	// Define the SMTP server details
	smtpHost := "mail.hostedemail.com"
	smtpPort := "587"

	mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, tenant, request.From)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	toEmail := []string{}
	ccEmail := []string{}
	bccEmail := []string{}

	//build internal object for transfer
	retMail := parsemail.Email{}
	retMail.HTMLBody = request.Content
	retMail.Subject = *request.Subject

	sentByName := "" //TODO
	fromAddress := []*mimemail.Address{{sentByName, request.From}}
	retMail.From = fromAddress
	var toAddress []*mimemail.Address
	var ccAddress []*mimemail.Address
	var bccAddress []*mimemail.Address
	for _, to := range request.To {
		toAddress = append(toAddress, &mimemail.Address{Address: to})
		retMail.To = toAddress
		toEmail = append(toEmail, to)
	}
	if request.Cc != nil {
		for _, cc := range request.Cc {
			ccAddress = append(ccAddress, &mimemail.Address{Address: cc})
			retMail.Cc = ccAddress
			ccEmail = append(ccEmail, cc)
		}
	}
	if request.Bcc != nil {
		for _, bcc := range request.Bcc {
			bccAddress = append(bccAddress, &mimemail.Address{Address: bcc})
			retMail.Bcc = bccAddress
			bccEmail = append(bccEmail, bcc)
		}
	}

	interactionEventNode, err := s.services.Neo4jRepositories.CommonReadRepository.GetById(ctx, tenant, *request.ReplyTo, commonModel.NodeLabelInteractionEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	interactionEvent := neo4jmapper.MapDbNodeToInteractionEventEntity(interactionEventNode)

	emailChannelData := dto.EmailChannelData{}
	err = json.Unmarshal([]byte(interactionEvent.ChannelData), &emailChannelData)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("unable to parse email channel data for %s", *request.ReplyTo)
	}

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
	subject := emailChannelData.Subject
	if len(emailChannelData.Subject) < 3 || emailChannelData.Subject[:3] != "Re:" {
		subject = "Re: " + emailChannelData.Subject
	}

	plainText, err := HTMLToPlainText(request.Content)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if emailChannelData.Reference == nil {
		emailChannelData.Reference = []string{}
	}

	emailChannelData.Reference = append(emailChannelData.Reference, emailChannelData.ProviderMessageId)

	references := strings.Join(emailChannelData.Reference, " ")

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
		FromEmail:  request.From,
		ToEmail:    strings.Join(toEmail, ", "),
		CCEmail:    strings.Join(ccEmail, ", "),
		BCCEmail:   strings.Join(bccEmail, ", "),
		Subject:    subject,
		Date:       time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		MessageId:  generateMessageID(mailbox.MailboxUsername),
		InReplyTo:  emailChannelData.ProviderMessageId,
		References: references,
		Boundary:   fmt.Sprintf("=_%x", time.Now().UnixNano()),
		PlainBody:  plainText,
		HTMLBody:   request.Content,
	}

	tmpl, err := template.New("email").Parse(messageTemplate)
	if err != nil {
		return nil, err
	}

	var msgBuffer bytes.Buffer
	if err := tmpl.Execute(&msgBuffer, data); err != nil {
		return nil, err
	}

	msg := msgBuffer.String()

	// Combine all recipients: To, CC, BCC
	recipients := []string{}
	recipients = append(recipients, toEmail...)
	recipients = append(recipients, ccEmail...)
	recipients = append(recipients, bccEmail...)

	auth := smtp.PlainAuth("", mailbox.MailboxUsername, mailbox.MailboxPassword, smtpHost)

	// Send the email
	err = smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		auth,
		request.From,
		recipients,
		[]byte(msg),
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	retMail.MessageID = data.MessageId
	retMail.Header = make(map[string][]string)
	retMail.Header["Message-Id"] = []string{data.MessageId}
	retMail.InReplyTo = []string{data.InReplyTo}
	retMail.References = emailChannelData.Reference

	return &retMail, nil
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

func NewOpenSRSService(services *Services) OpenSrsService {
	return &openSRSService{
		services: services,
	}
}
