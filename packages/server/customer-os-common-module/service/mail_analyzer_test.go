package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestProcessEmailCheck(t *testing.T) {
	tests := []struct {
		name     string
		tenant   string
		email    *EmailMessageData
		expected HeaderAnalysis
	}{
		{
			name:   "Bounce - Failed Recipients",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{XFailedRecepients: true},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBounce:     true,
				SkipReason:   "BOUNCE | X-FAILED-RECIPIENTS",
			},
		},
		{
			name:   "Bounce - Delivery Report",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{ContentDescription: "delivery report"},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBounce:     true,
				SkipReason:   "BOUNCE | CONTENT-DESCRIPTION: DELIVERY REPORT",
			},
		},
		{
			name:   "Bounce - Return Path Contains Mailer Daemon",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{ReturnPath: "mailer-daemon@example.com"},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBounce:     true,
				SkipReason:   "BOUNCE | RETURN-PATH CONTAINS BOUNCE KEYWORDS",
			},
		},
		{
			name:   "Bounce - From Contains Mailer Daemon",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Participants: EmailParticipants{
					From: EmailParticipant{Email: "mailer-daemon@example.com"},
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBounce:     true,
				SkipReason:   "BOUNCE | FROM CONTAINS BOUNCE KEYWORDS",
			},
		},
		{
			name:   "Bounce - Subject Contains Bounce Keywords",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Content: EmailContent{Subject: "Delivery Status Notification (Failure)"},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBounce:     true,
				SkipReason:   "BOUNCE | SUBJECT CONTAINS BOUNCE KEYWORDS",
			},
		},
		{
			name:   "Autoresponder - X-Autoreply",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{XAutoreply: "yes"},
			},
			expected: HeaderAnalysis{
				ProcessEmail:    false,
				IsAutoResponder: true,
				SkipReason:      "AUTORESPONDER | X-AUTOREPLY",
			},
		},
		{
			name:   "Autoresponder - X-Autoresponse",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{XAutoresponse: "yes"},
			},
			expected: HeaderAnalysis{
				ProcessEmail:    false,
				IsAutoResponder: true,
				SkipReason:      "AUTORESPONDER | X-AUTORESPONSE",
			},
		},
		{
			name:   "Autoresponder - X-Loop",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{XLoop: true},
			},
			expected: HeaderAnalysis{
				ProcessEmail:    false,
				IsAutoResponder: true,
				SkipReason:      "AUTORESPONDER | X-LOOP",
			},
		},
		{
			name:   "Autoresponder - Precedence Auto Reply",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{Precedence: "auto_reply"},
			},
			expected: HeaderAnalysis{
				ProcessEmail:    false,
				IsAutoResponder: true,
				SkipReason:      "AUTORESPONDER | PRECEDENCE: AUTO_REPLY",
			},
		},
		{
			name:   "Bulk - Different Reply-To",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Participants: EmailParticipants{
					From:    EmailParticipant{Email: "sender@example.com"},
					ReplyTo: []EmailParticipant{{Email: "different@example.com"}},
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBulkMail:   true,
				SkipReason:   "BULK | REPLY-TO != FROM",
			},
		},
		{
			name:   "Bulk - List Unsubscribe",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{ListUnsubscribe: true},
				Participants: EmailParticipants{
					From:    EmailParticipant{Email: "sender@example.com"},
					ReplyTo: []EmailParticipant{{Email: "sender@example.com"}},
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBulkMail:   true,
				SkipReason:   "BULK | UNSUBSCRIBE",
			},
		},
		{
			name:   "Bulk - Precedence Bulk",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{Precedence: "bulk"},
				Participants: EmailParticipants{
					From:    EmailParticipant{Email: "sender@example.com"},
					ReplyTo: []EmailParticipant{{Email: "sender@example.com"}},
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBulkMail:   true,
				SkipReason:   "BULK | PRECEDENCE: BULK",
			},
		},
		{
			name:   "Bulk - Empty Return Path",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{ReturnPath: ""},
				Participants: EmailParticipants{
					From:    EmailParticipant{Email: "sender@example.com"},
					ReplyTo: []EmailParticipant{{Email: "sender@example.com"}},
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBulkMail:   true,
				SkipReason:   "BULK | EMPTY RETURN-PATH",
			},
		},
		{
			name:   "Bulk - Return Path Different From From",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{ReturnPath: "different@example.com"},
				Participants: EmailParticipants{
					From:    EmailParticipant{Email: "sender@example.com"},
					ReplyTo: []EmailParticipant{{Email: "sender@example.com"}},
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBulkMail:   true,
				SkipReason:   "BULK | RETURN-PATH != FROM",
			},
		},
		{
			name:   "Bulk - Sender Different From From",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{
					Sender:     "different@example.com",
					ReturnPath: "sender@example.com",
				},
				Participants: EmailParticipants{
					From:    EmailParticipant{Email: "sender@example.com"},
					ReplyTo: []EmailParticipant{{Email: "sender@example.com"}},
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBulkMail:   true,
				SkipReason:   "BULK | SENDER != FROM",
			},
		},
		{
			name:   "Normal Email",
			tenant: "test-tenant",
			email: &EmailMessageData{
				Headers: EmailHeaders{
					ReturnPath: "sender@example.com",
					Sender:     "sender@example.com",
				},
				Participants: EmailParticipants{
					From:    EmailParticipant{Email: "sender@example.com"},
					ReplyTo: []EmailParticipant{{Email: "sender@example.com"}},
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: true,
			},
		},
	}

	svc := &mailService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.ProcessEmailCheck(context.Background(), tt.tenant, tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsBounceSubject(t *testing.T) {
	tests := []struct {
		name     string
		subject  string
		expected bool
	}{
		{
			name:     "Delivery Status Notification",
			subject:  "Delivery Status Notification (Failure)",
			expected: true,
		},
		{
			name:     "Undeliverable",
			subject:  "Undeliverable: Your message",
			expected: true,
		},
		{
			name:     "Undelivered",
			subject:  "Undelivered Mail Returned to Sender",
			expected: true,
		},
		{
			name:     "Delivery Failure",
			subject:  "Delivery Failure Notice",
			expected: true,
		},
		{
			name:     "Failure Notice",
			subject:  "Failure Notice: Unable to deliver",
			expected: true,
		},
		{
			name:     "Returned Mail",
			subject:  "Returned mail: User unknown",
			expected: true,
		},
		{
			name:     "Returned to Sender",
			subject:  "Mail Returned to Sender",
			expected: true,
		},
		{
			name:     "Case Insensitive",
			subject:  "DELIVERY STATUS NOTIFICATION",
			expected: true,
		},
		{
			name:     "Normal Subject",
			subject:  "Meeting Tomorrow",
			expected: false,
		},
		{
			name:     "Empty Subject",
			subject:  "",
			expected: false,
		},
	}

	svc := &mailService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.isBounceSubject(tt.subject)
			assert.Equal(t, tt.expected, result)
		})
	}
}
