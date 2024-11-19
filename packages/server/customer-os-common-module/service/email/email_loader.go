package service

import (
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetEmailMessageData(rawEmail string) (EmailMessageData, error) {
	email := EmailMessageData{}

	emailData, err := getRawEmailData(rawEmail)
	if err != nil {
		return email, err
	}

	email.Identifiers.ProviderMessageId = emailData.ProviderMessageId
	email.Identifiers.MessageId = emailData.MessageId
	email.Identifiers.EmailThreadId = emailData.ThreadId
	email.Identifiers.Reference = emailData.Reference

	email.Content.SentDate = emailData.Sent
	email.Content.Subject = emailData.Subject
	email.Content.Html = emailData.Html
	email.Content.Text = emailData.Text

	email.Participants.From = parseEmailAndName(emailData.From)
	email.Participants.To = parseParticipants(emailData.To)
	email.Participants.Cc = parseParticipants(emailData.Cc)
	email.Participants.Bcc = parseParticipants(emailData.Bcc)
	email.Participants.InReplyTo = parseEmailAndName(emailData.InReplyTo)

	email.Headers = parseHeaders(emailData.Headers)

	return email, nil
}

func getRawEmailData(rawEmail string) (EmailRawData, error) {
	rawEmailData := EmailRawData{}
	err := json.Unmarshal([]byte(rawEmail), &rawEmailData)
	if err != nil {
		logrus.Errorf("Unmarshal Raw Email Data Failed: %v", err)
		return rawEmailData, err
	}

	return rawEmailData, nil
}

func parseEmailAndName(s string) EmailParticipant {
	results := EmailParticipant{}

	// Handle bare email case
	if !strings.Contains(s, "<") {
		results.Email = strings.TrimSpace(s)
		return results
	}

	// Extract email
	results.Email = extractEmailFromBrackets(s)

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

func parseParticipants(s string) []EmailParticipant {
	participants := []EmailParticipant{}

	// Split on commas
	parts := strings.Split(s, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		participant := parseEmailAndName(part)
		participants = append(participants, participant)
	}

	return participants
}

func extractEmailFromBrackets(s string) string {
	if start := strings.LastIndex(s, "<"); start >= 0 {
		if end := strings.LastIndex(s, ">"); end > start {
			return s[start+1 : end]
		}
	}
	return s
}

func parseHeaders(headers map[string]string) EmailHeaders {
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
			eh.ReturnPath = extractEmailFromBrackets(value)
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
