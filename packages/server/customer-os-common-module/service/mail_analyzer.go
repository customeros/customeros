package service

import (
	"regexp"
	"strings"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/customeros/mailsherpa/mailvalidate"
	"golang.org/x/net/context"
)

type HeaderAnalysis struct {
	ProcessEmail    bool
	IsBounce        bool
	IsAutoResponder bool
	IsBulkMail      bool
	SkipReason      string
}

// TODO parse SMTP status code from message/deliver-status
// and classify bounced email as hard or soft bounce

func (a *mailService) ProcessEmailCheck(ctx context.Context, email *EmailMessageData) HeaderAnalysis {
	span, ctx := a.initializeTracing(ctx, "MailService.ProcessEmailCheck")
	defer span.Finish()

	analysis := HeaderAnalysis{
		ProcessEmail: true, // Default to processing
	}

	// Check bounce
	bounce, reason := a.isBounce(email.Headers, email.Content.Subject, email.Participants.From.Email)
	if bounce {
		analysis.IsBounce = true
		analysis.ProcessEmail = false
		analysis.SkipReason = reason
		return analysis
	}

	// Check auto-responder first
	autoresp, reason := a.isAutoResponder(email.Headers)
	if autoresp {
		analysis.IsAutoResponder = true
		analysis.ProcessEmail = false
		analysis.SkipReason = reason
		return analysis
	}

	// Check bulk mail
	bulk, reason := a.isBulkMail(email.Headers, email.Participants.From.Email, email.Participants.ReplyTo)
	if bulk {
		analysis.IsBulkMail = true
		analysis.ProcessEmail = false
		analysis.SkipReason = reason
	}

	return analysis
}

func (a *mailService) isAutoResponder(headers EmailHeaders) (bool, string) {
	switch {
	case headers.XAutoreply != "":
		return true, "Autoresponder: X-Autoreply"
	case headers.XAutoresponse != "":
		return true, "Autoresponder: X-Autoresponse"
	case headers.XLoop:
		return true, "Autoresponder: X-Loop"
	case strings.EqualFold(headers.Precedence, "auto_reply"):
		return true, "Autoresponder: Precedence: auto_reply"
	default:
		return false, ""
	}
}

func (a *mailService) isBounce(headers EmailHeaders, subject, from string) (bool, string) {
	switch {
	case headers.XFailedRecepients:
		return true, "Bounce: X-Failed-Recipients"
	case strings.EqualFold(headers.ContentDescription, "delivery report"):
		return true, "Bounce: Content-Description: Delivery Report"
	case a.isReturnPathBounce(headers.ReturnPath):
		return true, "Bounce: Return-Path containts bounce keywords"
	case a.isReturnPathBounce(from):
		return true, "Bounce: From contains bounce keywords"
	case a.isBounceSubject(subject):
		return true, "Bounce: Subject contains bounce keywords"
	default:
		return false, ""
	}
}

func (a *mailService) isBulkMail(headers EmailHeaders, from string, replyTo []EmailParticipant) (bool, string) {
	matchReplyTo := false
	for _, replyToParticipant := range replyTo {
		if replyToParticipant.Email == from {
			matchReplyTo = true
			break
		}
	}

	switch {
	case !matchReplyTo:
		return true, "Bulk: Reply-To != From"
	case headers.ListUnsubscribe:
		return true, "Bulk: Unsubscribe"
	case strings.EqualFold(headers.Precedence, "bulk"):
		return true, "Bulk: Precidence: Bulk"
	case headers.ReturnPath == "":
		return true, "Bulk: Empty Return-Path"
	case headers.ReturnPath != from:
		return true, "Bulk: Return-Path != From"
	case (headers.Sender != "" && headers.Sender != from):
		return true, "Bulk: Sender != From"
	default:
		return a.mailsherpaChecks(from)
	}
}

func (a *mailService) isReturnPathBounce(returnPath string) bool {
	return strings.Contains(returnPath, "mailer-daemon")
}

func (a *mailService) mailsherpaChecks(from string) (failedCheck bool, reason string) {
	if from == "" {
		return true, "No From email"
	}
	syntaxValidation := mailvalidate.ValidateEmailSyntax(from)
	if syntaxValidation.IsRoleAccount {
		return true, "Bulk: From Role Account"
	}

	if syntaxValidation.IsSystemGenerated {
		return true, "Bulk: System generated"
	}

	primaryDomainCheck, _ := domaincheck.PrimaryDomainCheck(syntaxValidation.Domain)

	if !primaryDomainCheck {
		return true, "Bulk: From non-primary domain"
	}

	return false, ""
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

func (a *mailService) extractEmailAddresses(input string) []string {
	if input == "" {
		return []string{}
	}

	// Compile regex
	emailRegex := regexp.MustCompile(`<([^>]+)>|([^\s,<>]+@[^\s,<>]+)`)

	// Split, clean and normalize input
	emails := strings.Split(strings.ToLower(input), ",")

	// Use a map for deduplication
	uniqueEmails := make(map[string]struct{})

	for _, email := range emails {
		matches := emailRegex.FindAllStringSubmatch(strings.TrimSpace(email), -1)
		for _, match := range matches {
			// match[1] is from <...>, match[2] is raw email
			if email := match[1]; email != "" {
				if mailvalidate.ValidateEmailSyntax(email).IsValid {
					uniqueEmails[email] = struct{}{}
				}
			} else if email := match[2]; email != "" {
				if mailvalidate.ValidateEmailSyntax(email).IsValid {
					uniqueEmails[email] = struct{}{}
				}
			}
		}
	}

	if len(uniqueEmails) == 0 {
		return []string{input}
	}

	result := make([]string, 0, len(uniqueEmails))
	for email := range uniqueEmails {
		result = append(result, email)
	}
	return result
}
