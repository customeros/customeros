package service

import "strings"

type HeaderAnalysis struct {
	ProcessEmail    bool
	IsBounce        bool
	IsAutoResponder bool
	IsBulkMail      bool
}

// TODO parse SMTP status code from message/deliver-status
// and classify bounced email as hard or soft bounce

func (a *mailService) ProcessEmailCheck(email *EmailMessageData) HeaderAnalysis {
	analysis := HeaderAnalysis{
		ProcessEmail: true, // Default to processing
	}

	// Check bounce
	if a.isBounce(email.Headers, email.Content.Subject, email.Participants.From.Email) {
		analysis.IsBounce = true
		analysis.ProcessEmail = false
		return analysis
	}

	// Check auto-responder first
	if a.isAutoResponder(email.Headers) {
		analysis.IsAutoResponder = true
		analysis.ProcessEmail = false
		return analysis
	}

	// Check bulk mail
	if a.isBulkMail(email.Headers) {
		analysis.IsBulkMail = true
		analysis.ProcessEmail = false
	}

	return analysis
}

func (a *mailService) isAutoResponder(headers EmailHeaders) bool {
	return headers.XAutoreply != "" ||
		headers.XAutoresponse != "" ||
		headers.AutoSubmitted ||
		headers.XLoop ||
		strings.EqualFold(headers.Precedence, "auto_reply")
}

func (a *mailService) isBounce(headers EmailHeaders, subject, from string) bool {
	return headers.XFailedRecepients ||
		headers.DeliveryStatus ||
		strings.EqualFold(headers.ContentDescription, "delivery report") ||
		a.isReturnPathBounce(headers.ReturnPath) ||
		a.isReturnPathBounce(from) ||
		a.isBounceSubject(subject)
}

func (a *mailService) isBulkMail(headers EmailHeaders) bool {
	return headers.ListUnsubscribe ||
		strings.EqualFold(headers.Precedence, "bulk")
}

func (a *mailService) isReturnPathBounce(returnPath string) bool {
	return returnPath == "" ||
		strings.Contains(returnPath, "mailer-daemon") ||
		strings.Contains(returnPath, "postmaster")
}

func (a *mailService) isBounceSubject(subject string) bool {
	subject = strings.ToLower(subject)
	keywords := []string{
		"delivery status notification",
		"undeliverable",
		"undelivered",
		"delivery failure",
		"failure notice",
		"returned mail",
		"returned to sender",
	}
	for _, phrase := range keywords {
		if strings.Contains(subject, phrase) {
			return true
		}
	}

	return false
}
