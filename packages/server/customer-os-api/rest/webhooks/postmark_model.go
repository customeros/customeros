package webhooks

import (
	"fmt"
	"strings"
	"time"

	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type PostmarkInboundEmailData struct {
	FromName      string `json:"FromName"`
	MessageStream string `json:"MessageStream"`
	From          string `json:"From"`
	FromFull      struct {
		Email       string `json:"Email"`
		Name        string `json:"Name"`
		MailboxHash string `json:"MailboxHash"`
	} `json:"FromFull"`
	To     string `json:"To"`
	ToFull []struct {
		Email       string `json:"Email"`
		Name        string `json:"Name"`
		MailboxHash string `json:"MailboxHash"`
	} `json:"ToFull"`
	Cc     string `json:"Cc"`
	CcFull []struct {
		Email       string `json:"Email"`
		Name        string `json:"Name"`
		MailboxHash string `json:"MailboxHash"`
	} `json:"CcFull"`
	Bcc     string `json:"Bcc"`
	BccFull []struct {
		Email       string `json:"Email"`
		Name        string `json:"Name"`
		MailboxHash string `json:"MailboxHash"`
	} `json:"BccFull"`
	OriginalRecipient string `json:"OriginalRecipient"`
	Subject           string `json:"Subject"`
	MessageID         string `json:"MessageID"`
	ReplyTo           string `json:"ReplyTo"`
	MailboxHash       string `json:"MailboxHash"`
	Date              string `json:"Date"`
	TextBody          string `json:"TextBody"`
	HtmlBody          string `json:"HtmlBody"`
	StrippedTextReply string `json:"StrippedTextReply"`
	Tag               string `json:"Tag"`
	Headers           []struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	} `json:"Headers"`
	Attachments []struct {
		Name          string `json:"Name"`
		Content       string `json:"Content"`
		ContentType   string `json:"ContentType"`
		ContentLength int    `json:"ContentLength"`
	} `json:"Attachments"`
}

func (p *PostmarkInboundEmailData) TenantFromBcc() string {
	var tenant string
	for _, address := range p.BccFull {
		validation := mailvalidate.SyntaxValidation(address.Email)
		if validation.IsValid {
			tenant = strings.Split(validation.Domain, ".")[0]
		}
	}
	return tenant
}

func (p *PostmarkInboundEmailData) HasAttachments() bool {
	return len(p.Attachments) > 0
}

func (p *PostmarkInboundEmailData) HTMLBody() string {
	return strings.ReplaceAll(p.HtmlBody, "&amp;", "&")
}

func (p *PostmarkInboundEmailData) TEXTBody() string {
	return strings.ReplaceAll(p.TextBody, "&amp;", "&")
}

func (p *PostmarkInboundEmailData) AllToEmails() []string {
	allEmails := make([]string, 0)
	for _, to := range p.ToFull {
		if to.Email == "" {
			continue
		}
		allEmails = append(allEmails, to.Email)
	}
	return allEmails
}

func (p *PostmarkInboundEmailData) AllCcEmails() []string {
	allEmails := make([]string, 0)
	for _, cc := range p.CcFull {
		if cc.Email == "" {
			continue
		}
		allEmails = append(allEmails, cc.Email)
	}
	return allEmails
}

func (p *PostmarkInboundEmailData) AllBccEmails() []string {
	allEmails := make([]string, 0)
	for _, bcc := range p.BccFull {
		if bcc.Email == "" {
			continue
		}
		allEmails = append(allEmails, bcc.Email)
	}
	return allEmails
}

func (p *PostmarkInboundEmailData) AllParticipantEmails() []string {
	to := p.AllToEmails()
	cc := p.AllCcEmails()
	bcc := p.AllBccEmails()
	result := make([]string, 0, len(to)+len(cc)+len(bcc)+1)
	result = append(result, to...)
	result = append(result, cc...)
	result = append(result, bcc...)
	result = append(result, p.FromFull.Email)

	result = utils.RemoveDuplicates(result)

	return result
}

func (p *PostmarkInboundEmailData) EmailSentTimestamp() time.Time {
	var ts *time.Time
	ts, _ = utils.UnmarshalDateTime(p.Date)
	return *ts
}

func (p *PostmarkInboundEmailData) GetHeaderValue(header string) string {
	for _, h := range p.Headers {
		if strings.EqualFold(h.Name, header) {
			return h.Value
		}
	}
	return ""
}

func (p *PostmarkInboundEmailData) GetHeaders() map[string]string {
	headers := make(map[string]string)
	for _, header := range p.Headers {
		headers[header.Name] = header.Value
	}
	return headers
}

func (p *PostmarkInboundEmailData) ToRawDbObject() entity.EmailRawData {
	var result entity.RawEmailData

	messageId := p.GetHeaderValue("Message-Id")
	result.ProviderMessageId = messageId
	result.MessageId = messageId
	result.Sent = p.EmailSentTimestamp()
	result.Subject = p.Subject
	result.From = p.FromFull.Email
	result.Html = p.HTMLBody()
	result.Text = p.TEXTBody()
	result.InReplyTo = p.GetHeaderValue("In-Reply-To")
	result.Reference = p.GetHeaderValue("References")
	result.Headers = p.GetHeaders()
	result.To = EmailWithBrackets(p.AllToEmails())
	result.Cc = EmailWithBrackets(p.AllCcEmails())
	result.Bcc = EmailWithBrackets(p.AllBccEmails())

	return result
}

func EmailWithBrackets(s []string) string {
	emails := make([]string, 0)
	for _, v := range s {
		emails = append(emails, fmt.Sprintf("<%s>", v))
	}

	return strings.Join(emails, ", ")
}
