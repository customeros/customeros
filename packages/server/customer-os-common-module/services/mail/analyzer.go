package mail

import (
	"strings"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

// TODO parse SMTP status code from message/deliver-status
// and classify bounced email as hard or soft bounce
func (a *mailService) ProcessEmailCheck(ctx context.Context, tenant string, emailData *interfaces.EmailMessageData) interfaces.HeaderAnalysis {
	span, ctx := a.initializeTracing(ctx, "MailService.ProcessEmailCheck")
	defer span.Finish()

	analysis := interfaces.HeaderAnalysis{
		ProcessEmail: true, // Default to processing
	}

	// Check valid from email address
	if emailData.Participants.From.Email != "" {
		if !mailvalidate.ValidateEmailSyntax(emailData.Participants.From.Email).IsValid {
			analysis.ProcessEmail = false
			analysis.SkipReason = "INVALID FROM EMAIL ADDRESS FORMAT"

			span.LogFields(log.Bool("process_email", analysis.ProcessEmail))
			span.LogFields(log.String("reason", analysis.SkipReason))

			return analysis
		}
	}

	// Check bounce
	bounce, reason := a.isBounce(emailData.Headers, emailData.Content.Subject, emailData.Participants.From.Email)
	if bounce {
		analysis.IsBounce = true
		analysis.ProcessEmail = false
		analysis.SkipReason = reason

		xFailedHeaderExists := len(emailData.Headers.XFailedRecepients) > 0
		switch {
		case !xFailedHeaderExists && emailData.Identifiers.ExternalSystem == "outlook":
			analysis.BouncedEmails = emailData.Participants.GetToEmailAddresses()
		case !xFailedHeaderExists && emailData.Identifiers.ExternalSystem == "mailstack":
			analysis.BouncedEmails = a.ExtractEmails(emailData.Content.Text)
		default:
			analysis.BouncedEmails = emailData.Headers.XFailedRecepients
		}

		span.LogFields(log.Bool("process_email", analysis.ProcessEmail))
		span.LogFields(log.String("reason", analysis.SkipReason))

		return analysis
	}

	// Check warming email
	if a.isWarmingEmail(tenant, emailData) {
		analysis.ProcessEmail = false
		analysis.SkipReason = "WARMING"

		span.LogFields(log.Bool("process_email", analysis.ProcessEmail))
		span.LogFields(log.String("reason", analysis.SkipReason))

		return analysis
	}

	// Check auto-responder first
	autoresp, reason := a.isAutoResponder(emailData.Headers)
	if autoresp {
		analysis.IsAutoResponder = true
		analysis.ProcessEmail = false
		analysis.SkipReason = reason

		span.LogFields(log.Bool("process_email", analysis.ProcessEmail))
		span.LogFields(log.String("reason", analysis.SkipReason))

		return analysis
	}

	// Check bulk mail
	bulk, reason := a.isBulkMail(emailData.Headers, emailData.Participants.From.Email, emailData.Participants.ReplyTo)
	if bulk {
		analysis.IsBulkMail = true
		analysis.ProcessEmail = false
		analysis.SkipReason = reason
	}

	span.LogFields(log.Bool("process_email", analysis.ProcessEmail))
	span.LogFields(log.String("reason", analysis.SkipReason))

	return analysis
}

func (a *mailService) isAutoResponder(headers interfaces.EmailHeaders) (bool, string) {
	switch {
	case headers.XAutoreply != "":
		return true, "AUTORESPONDER | X-AUTOREPLY"
	case headers.XAutoresponse != "":
		return true, "AUTORESPONDER | X-AUTORESPONSE"
	case headers.XLoop:
		return true, "AUTORESPONDER | X-LOOP"
	case strings.EqualFold(headers.Precedence, "auto_reply"):
		return true, "AUTORESPONDER | PRECEDENCE: AUTO_REPLY"
	default:
		return false, ""
	}
}

func (a *mailService) isBounce(headers interfaces.EmailHeaders, subject, from string) (bool, string) {
	switch {
	case len(headers.XFailedRecepients) > 0:
		return true, "BOUNCE | X-FAILED-RECIPIENTS"
	case strings.EqualFold(headers.ContentDescription, "delivery report"):
		return true, "BOUNCE | CONTENT-DESCRIPTION: DELIVERY REPORT"
	case a.isReturnPathBounce(headers.ReturnPath):
		return true, "BOUNCE | RETURN-PATH CONTAINS BOUNCE KEYWORDS"
	case a.isReturnPathBounce(from):
		return true, "BOUNCE | FROM CONTAINS BOUNCE KEYWORDS"
	case a.isBounceSubject(subject):
		return true, "BOUNCE | SUBJECT CONTAINS BOUNCE KEYWORDS"
	default:
		return false, ""
	}
}

func (a *mailService) isBulkMail(headers interfaces.EmailHeaders, from string, replyTo []interfaces.EmailParticipant) (bool, string) {
	matchReplyTo := false
	for _, replyToParticipant := range replyTo {
		if replyToParticipant.Email == from {
			matchReplyTo = true
			break
		}
	}

	if headers.ForwardedFor == "" {
		switch {
		case (headers.ReplyToExists && !matchReplyTo):
			return true, "BULK | REPLY-TO != FROM"
		case headers.ReturnPathExists && headers.ReturnPath == "":
			return true, "BULK | EMPTY RETURN-PATH"
		case headers.ReturnPathExists && strings.Index(headers.ReturnPath, from) == -1:
			return true, "BULK | RETURN-PATH != FROM"
		default:
		}
	}

	switch {
	case headers.ListUnsubscribe:
		return true, "BULK | UNSUBSCRIBE"
	case strings.EqualFold(headers.Precedence, "bulk"):
		return true, "BULK | PRECEDENCE: BULK"
	case headers.Sender != "" && headers.Sender != from:
		return true, "BULK | SENDER != FROM"
	default:
		return a.mailsherpaChecks(from)
	}
}

func (a *mailService) isReturnPathBounce(returnPath string) bool {
	return strings.Contains(strings.ToLower(returnPath), "mailer-daemon")
}

func (a *mailService) mailsherpaChecks(from string) (failedCheck bool, reason string) {
	if from == "" {
		return true, "NO FROM EMAIL"
	}
	syntaxValidation := mailvalidate.ValidateEmailSyntax(from)
	if syntaxValidation.IsRoleAccount {
		return true, "BULK | FROM ROLE ACCOUNT"
	}

	if syntaxValidation.IsSystemGenerated {
		return true, "BULK | SYSTEM GENERATED"
	}

	isPrimaryDomain, _ := domaincheck.PrimaryDomainCheck(syntaxValidation.Domain)

	if !isPrimaryDomain {
		return true, "BULK | FROM NON-PRIMARY DOMAIN"
	}

	return false, ""
}

func (a *mailService) isBounceSubject(subject string) bool {
	subject = strings.ToLower(subject)
	keywords := []string{
		"mail delivery failure",
		"undelivered mail returned to sender",
		"delivery status notification",
		"undeliverable",
		"undelivered",
		"delivery failure",
		"failure notice",
		"returned mail",
		"returned to sender",
	}
	for _, phrase := range keywords {
		if strings.Contains(strings.ToLower(subject), phrase) {
			return true
		}
	}

	return false
}

func (s *mailService) isWarmingEmail(tenant string, email *interfaces.EmailMessageData) bool {
	emailExclusion := s.cache.GetEmailExclusion(tenant)

	for _, exclusion := range emailExclusion {
		if exclusion.ExcludeSubject != nil {
			if strings.Contains(email.Content.Subject, *exclusion.ExcludeSubject) {
				return true
			}
		}
		if exclusion.ExcludeBody != nil {
			if strings.Contains(email.Content.Html, *exclusion.ExcludeBody) {
				return true
			}
			if strings.Contains(email.Content.Text, *exclusion.ExcludeBody) {
				return true
			}
		}
	}
	return false
}
