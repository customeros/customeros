package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseEmailAndName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected EmailParticipant
	}{
		{
			name:  "Full name with email",
			input: "John Smith <john@example.com>",
			expected: EmailParticipant{
				Email:     "john@example.com",
				FirstName: "John",
				LastName:  "Smith",
			},
		},
		{
			name:  "Only email in brackets",
			input: "<test@example.com>",
			expected: EmailParticipant{
				Email: "test@example.com",
			},
		},
		{
			name:  "Bare email",
			input: "test@example.com",
			expected: EmailParticipant{
				Email: "test@example.com",
			},
		},
		{
			name:  "Three part name",
			input: "John van Smith <john@example.com>",
			expected: EmailParticipant{
				Email:     "john@example.com",
				FirstName: "John",
				LastName:  "van Smith",
			},
		},
		{
			name:  "Service name with email",
			input: "Google Workspace Alerts <google-workspace-alerts-noreply@google.com>",
			expected: EmailParticipant{
				Email:     "google-workspace-alerts-noreply@google.com",
				FirstName: "Google",
				LastName:  "Workspace Alerts",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseEmailAndName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseParticipants(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []EmailParticipant
	}{
		{
			name:  "Multiple participants with commas",
			input: "John Smith <john@example.com>, Jane Doe <jane@example.com>",
			expected: []EmailParticipant{
				{
					Email:     "john@example.com",
					FirstName: "John",
					LastName:  "Smith",
				},
				{
					Email:     "jane@example.com",
					FirstName: "Jane",
					LastName:  "Doe",
				},
			},
		},
		{
			name:  "Single participant",
			input: "John Smith <john@example.com>",
			expected: []EmailParticipant{
				{
					Email:     "john@example.com",
					FirstName: "John",
					LastName:  "Smith",
				},
			},
		},
		{
			name:  "Bare email",
			input: "test@example.com",
			expected: []EmailParticipant{
				{
					Email: "test@example.com",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseParticipants(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected EmailHeaders
	}{
		{
			name: "Full headers",
			headers: map[string]string{
				"Auto-Submitted":      "auto-replied",
				"Content-Description": "notification",
				"Content-Type":        "multipart/report; report-type=delivery-status",
				"List-Unsubscribe":    "<mailto:unsub@example.com>",
				"Precedence":          "bulk",
				"Return-Path":         "<bounce@example.com>",
				"X-Autoreply":         "yes",
				"X-Autoresponse":      "out of office",
				"X-Loop":              "yes",
				"X-Failed-Recipients": "failed@example.com",
			},
			expected: EmailHeaders{
				AutoSubmitted:      true,
				ContentDescription: "notification",
				DeliveryStatus:     true,
				ListUnsubscribe:    true,
				Precedence:         "bulk",
				ReturnPath:         "bounce@example.com",
				XAutoreply:         "yes",
				XAutoresponse:      "out of office",
				XLoop:              true,
				XFailedRecepients:  true,
				RawHeaders:         map[string]string{}, // Will be populated with input headers
			},
		},
		{
			name:    "Empty headers",
			headers: map[string]string{},
			expected: EmailHeaders{
				RawHeaders: map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHeaders(tt.headers)
			tt.expected.RawHeaders = tt.headers // Set expected raw headers
			assert.Equal(t, tt.expected, result)
		})
	}
}
