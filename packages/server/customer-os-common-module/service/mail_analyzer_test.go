package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessEmailCheck(t *testing.T) {
	svc := &emailService{}

	tests := []struct {
		name     string
		email    *EmailMessageData
		expected HeaderAnalysis
	}{
		{
			name: "Bounce email",
			email: &EmailMessageData{
				Headers: EmailHeaders{
					DeliveryStatus: true,
				},
				Content: EmailContent{
					Subject: "Delivery Status Notification (Failure)",
				},
				Participants: EmailParticipants{
					From: EmailParticipant{
						Email: "mailer-daemon@example.com",
					},
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBounce:     true,
			},
		},
		{
			name: "Auto-responder email",
			email: &EmailMessageData{
				Headers: EmailHeaders{
					AutoSubmitted: true,
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail:    false,
				IsAutoResponder: true,
			},
		},
		{
			name: "Bulk mail",
			email: &EmailMessageData{
				Headers: EmailHeaders{
					ListUnsubscribe: true,
				},
			},
			expected: HeaderAnalysis{
				ProcessEmail: false,
				IsBulkMail:   true,
			},
		},
		{
			name: "Normal email",
			email: &EmailMessageData{
				Headers: EmailHeaders{},
			},
			expected: HeaderAnalysis{
				ProcessEmail: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.ProcessEmailCheck(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAutoResponder(t *testing.T) {
	svc := &emailService{}

	tests := []struct {
		name     string
		headers  EmailHeaders
		expected bool
	}{
		{
			name: "X-Autoreply header",
			headers: EmailHeaders{
				XAutoreply: "yes",
			},
			expected: true,
		},
		{
			name: "Auto-Submitted header",
			headers: EmailHeaders{
				AutoSubmitted: true,
			},
			expected: true,
		},
		{
			name:     "Not auto-response",
			headers:  EmailHeaders{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.isAutoResponder(tt.headers)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsBounce(t *testing.T) {
	svc := &emailService{}

	tests := []struct {
		name     string
		headers  EmailHeaders
		subject  string
		from     string
		expected bool
	}{
		{
			name: "Delivery Status header",
			headers: EmailHeaders{
				DeliveryStatus: true,
			},
			expected: true,
		},
		{
			name:     "Bounce subject",
			headers:  EmailHeaders{},
			subject:  "Delivery Status Notification (Failure)",
			expected: true,
		},
		{
			name:     "Mailer daemon from",
			headers:  EmailHeaders{},
			from:     "mailer-daemon@example.com",
			expected: true,
		},
		{
			name:     "Not bounce",
			headers:  EmailHeaders{},
			subject:  "Hello",
			from:     "user@example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.isBounce(tt.headers, tt.subject, tt.from)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsBounceSubject(t *testing.T) {
	svc := &emailService{}

	tests := []struct {
		name     string
		subject  string
		expected bool
	}{
		{
			name:     "Delivery status notification",
			subject:  "Delivery Status Notification (Failure)",
			expected: true,
		},
		{
			name:     "Undeliverable",
			subject:  "Undeliverable: Your message",
			expected: true,
		},
		{
			name:     "Normal subject",
			subject:  "Hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.isBounceSubject(tt.subject)
			assert.Equal(t, tt.expected, result)
		})
	}
}
