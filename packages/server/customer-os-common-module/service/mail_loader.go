package service

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

func (l *mailService) LoadEmail(ctx context.Context, rawEmail *entity.RawEmail) (EmailMessageData, error) {
	span, ctx := l.initializeTracing(ctx, "MailService.LoadEmail")
	defer span.Finish()
	span.LogFields(tracingLog.Object("rawEmail", rawEmail))

	email := EmailMessageData{}

	emailData, err := l.getRawEmailData(rawEmail, span)
	if err != nil {
		err = fmt.Errorf("failed to get raw email data: %v", err)
		tracing.TraceErr(span, err)
		return email, err
	}

	email.Headers = l.parseHeaders(emailData.Headers)

	email.Identifiers.ProviderMessageId = emailData.ProviderMessageId
	email.Identifiers.MessageId = emailData.MessageId
	email.Identifiers.EmailThreadId = emailData.ThreadId
	email.Identifiers.References = extractLines(emailData.Reference)

	email.Content.SentDate = emailData.Sent
	email.Content.Subject = emailData.Subject
	email.Content.Html = emailData.Html
	email.Content.Text = emailData.Text

	email.Participants.From = l.parseEmailAndName(emailData.From)
	email.Participants.To = l.parseParticipants(emailData.To)
	email.Participants.Cc = l.parseParticipants(emailData.Cc)
	email.Participants.Bcc = l.parseParticipants(emailData.Bcc)
	email.Participants.ReplyTo = []EmailParticipant{l.parseEmailAndName(email.Headers.ReplyTo)}
	l.getAllEmails(&email.Participants)

	return email, nil
}

func (l *mailService) getRawEmailData(rawEmail *entity.RawEmail, span opentracing.Span) (EmailRawData, error) {
	rawEmailData := EmailRawData{}
	err := json.Unmarshal([]byte(rawEmail.Data), &rawEmailData)
	if err != nil {
		err = fmt.Errorf("Unmarshal Raw Email Data Failed: %v", err)
		tracing.TraceErr(span, err)
		return rawEmailData, err
	}

	return rawEmailData, nil
}

func (l *mailService) getAllEmails(contacts *EmailParticipants) {
	var all []string
	all = append(all, contacts.From.Email)

	for _, c := range contacts.To {
		if !slices.Contains(all, c.Email) && c.Email != "" {
			all = append(all, c.Email)
		}
	}

	for _, c := range contacts.Cc {
		if !slices.Contains(all, c.Email) && c.Email != "" {
			all = append(all, c.Email)
		}
	}

	for _, c := range contacts.Bcc {
		if !slices.Contains(all, c.Email) && c.Email != "" {
			all = append(all, c.Email)
		}
	}

	contacts.AllEmails = all
}

func (l *mailService) parseEmailAndName(s string) EmailParticipant {
	s = strings.ToLower(s)
	results := EmailParticipant{}

	// Handle bare email case
	if !strings.Contains(s, "<") {
		results.Email = strings.TrimSpace(s)
		return results
	}

	// Extract email
	results.Email = l.extractEmailFromBrackets(s)

	// Extract name part
	namePart := strings.TrimSpace(strings.Split(s, "<")[0])

	// Split name into parts
	names := strings.Fields(namePart)
	if len(names) > 0 {
		results.FirstName = names[0]
		if len(names) > 1 {
			results.LastName = strings.Join(names[1:], " ")
		}
	}

	return results
}

func (l *mailService) parseParticipants(s string) []EmailParticipant {
	s = strings.ToLower(s)
	participants := []EmailParticipant{}

	// Split on commas
	parts := strings.Split(s, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		participant := l.parseEmailAndName(part)
		participants = append(participants, participant)
	}

	return participants
}

func (l *mailService) extractEmailFromBrackets(s string) string {
	s = strings.ToLower(s)
	if start := strings.LastIndex(s, "<"); start >= 0 {
		if end := strings.LastIndex(s, ">"); end > start {
			return s[start+1 : end]
		}
	}
	return s
}

func (l *mailService) parseHeaders(headers map[string]string) EmailHeaders {
	eh := EmailHeaders{}

	// Store raw headers
	eh.RawHeaders = headers

	// Parse specific headers
	for header, value := range headers {
		if strings.EqualFold(header, "Auto-Submitted") {
			eh.AutoSubmitted = true
		}
		if strings.EqualFold(header, "Content-Description") {
			eh.ContentDescription = value
		}
		if strings.EqualFold(header, "Content-Type") &&
			strings.Contains(value, "delivery-status") {
			eh.DeliveryStatus = true
		}
		if strings.EqualFold(header, "List-Unsubscribe") {
			eh.ListUnsubscribe = true
		}
		if strings.EqualFold(header, "Precedence") {
			eh.Precedence = value
		}
		if strings.EqualFold(header, "Return-Path") {
			eh.ReturnPath = l.extractEmailFromBrackets(value)
		}
		if strings.EqualFold(header, "X-Autoreply") {
			eh.XAutoreply = value
		}
		if strings.EqualFold(header, "X-Autoresponse") {
			eh.XAutoresponse = value
		}
		if strings.EqualFold(header, "X-Loop") {
			eh.XLoop = true
		}
		if strings.EqualFold(header, "X-Failed-Recipients") {
			eh.XFailedRecepients = true
		}
		if strings.EqualFold(header, "Reply-To") {
			eh.ReplyTo = value
		}
		if strings.EqualFold(header, "Sender") {
			eh.ReturnPath = l.extractEmailFromBrackets(value)
		}
	}

	return eh
}

func extractLines(input string) []string {
	lines := strings.Fields(input)
	return lines
}
