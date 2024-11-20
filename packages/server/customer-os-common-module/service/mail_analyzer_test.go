package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestProcessEmailCheck(t *testing.T) {
	tests := []struct {
		name     string
		email    *EmailMessageData
		expected HeaderAnalysis
	}{
		{
			name: "Bounce - Failed Recipients",
			email: &EmailMessageData{
				Headers: EmailHeaders{XFailedRecepients: true},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBounce:     true,
				SkipReason:   "Bounce: X-Failed-Recipients",
			},
		},
		{
			name: "Bounce - Delivery Report",
			email: &EmailMessageData{
				Headers: EmailHeaders{ContentDescription: "delivery report"},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBounce:     true,
				SkipReason:   "Bounce: Content-Description: Delivery Report",
			},
		},
		{
			name: "Autoresponder - X-Autoreply",
			email: &EmailMessageData{
				Headers: EmailHeaders{XAutoreply: "yes"},
			},
			expected: HeaderAnalysis{
				ProcessEmail:    false,
				IsAutoResponder: true,
				SkipReason:      "Autoresponder: X-Autoreply",
			},
		},
		{
			name: "Autoresponder - X-Loop",
			email: &EmailMessageData{
				Headers: EmailHeaders{XLoop: true},
			},
			expected: HeaderAnalysis{
				ProcessEmail:    false,
				IsAutoResponder: true,
				SkipReason:      "Autoresponder: X-Loop",
			},
		},
		{
			name: "Bulk - List Unsubscribe",
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
				SkipReason:   "Bulk: Unsubscribe",
			},
		},
		{
			name: "Bulk - Different Reply-To",
			email: &EmailMessageData{
				Participants: EmailParticipants{
					From:    EmailParticipant{Email: "sender@example.com"},
					ReplyTo: []EmailParticipant{{Email: "different@example.com"}},
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBulkMail:   true,
				SkipReason:   "Bulk: Reply-To != From",
			},
		},
		{
			name: "Normal Email",
			email: &EmailMessageData{
				Participants: EmailParticipants{
					From:    EmailParticipant{Email: "sender@example.com"},
					ReplyTo: []EmailParticipant{{Email: "sender@example.com"}},
				},
				Headers: EmailHeaders{
					ReturnPath: "sender@example.com",
					Sender:     "sender@example.com",
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
			result := svc.ProcessEmailCheck(context.Background(), tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractEmailAddresses(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple email",
			input:    "user@example.com",
			expected: []string{"user@example.com"},
		},
		{
			name:     "Email with display name",
			input:    "User Name <user@example.com>",
			expected: []string{"user@example.com"},
		},
		{
			name:     "Multiple emails",
			input:    "first@example.com, second@example.com",
			expected: []string{"first@example.com", "second@example.com"},
		},
		{
			name:     "Multiple emails with display names",
			input:    "First User <first@example.com>, Second User <second@example.com>",
			expected: []string{"first@example.com", "second@example.com"},
		},
		{
			name:     "Empty input",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "Invalid email",
			input:    "not-an-email",
			expected: []string{},
		},
	}

	svc := &mailService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.extractEmailAddresses(tt.input)
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
			name:     "Delivery Failure",
			subject:  "Delivery Failure Notice",
			expected: true,
		},
		{
			name:     "Returned Mail",
			subject:  "Returned mail: User unknown",
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
