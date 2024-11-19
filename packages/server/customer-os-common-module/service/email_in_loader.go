package service

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/sirupsen/logrus"
)

func (l *emailService) LoadEmail(rawEmail *entity.RawEmail) (EmailMessageData, error) {
	email := EmailMessageData{}

	emailData, err := l.getRawEmailData(rawEmail)
	if err != nil {
		return email, err
	}

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
	email.Participants.InReplyTo = l.parseEmailAndName(emailData.InReplyTo)
	l.getAllEmails(&email.Participants)

	email.Headers = l.parseHeaders(emailData.Headers)

	return email, nil
}

func (l *emailService) getRawEmailData(rawEmail *entity.RawEmail) (EmailRawData, error) {
	rawEmailData := EmailRawData{}
	err := json.Unmarshal([]byte(rawEmail.Data), &rawEmailData)
	if err != nil {
		logrus.Errorf("Unmarshal Raw Email Data Failed: %v", err)
		return rawEmailData, err
	}

	return rawEmailData, nil
}

func (l *emailService) getAllEmails(contacts *EmailParticipants) {
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

func (l *emailService) parseEmailAndName(s string) EmailParticipant {
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

func (l *emailService) parseParticipants(s string) []EmailParticipant {
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

func (l *emailService) extractEmailFromBrackets(s string) string {
	s = strings.ToLower(s)
	if start := strings.LastIndex(s, "<"); start >= 0 {
		if end := strings.LastIndex(s, ">"); end > start {
			return s[start+1 : end]
		}
	}
	return s
}

func (l *emailService) parseHeaders(headers map[string]string) EmailHeaders {
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
	}

	return eh
}

func extractLines(input string) []string {
	lines := strings.Fields(input)
	return lines
}
