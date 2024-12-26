package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEmailAndName(t *testing.T) {
	svc := &mailService{}
	tests := []struct {
		name     string
		input    string
		expected EmailParticipant
	}{
		{
			name:  "Full name with email",
			input: "John Smith <john.smith@example.com>",
			expected: EmailParticipant{
				FirstName: "john",
				LastName:  "smith",
				Email:     "john.smith@example.com",
			},
		},
		{
			name:  "First name only with email",
			input: "John <john@example.com>",
			expected: EmailParticipant{
				FirstName: "john",
				Email:     "john@example.com",
			},
		},
		{
			name:  "Bare email",
			input: "john@example.com",
			expected: EmailParticipant{
				Email: "john@example.com",
			},
		},
		{
			name:  "Three part name",
			input: "John Middle Smith <john.smith@example.com>",
			expected: EmailParticipant{
				FirstName: "john",
				LastName:  "middle smith",
				Email:     "john.smith@example.com",
			},
		},
		{
			name:  "Empty input",
			input: "",
			expected: EmailParticipant{
				Email: "",
			},
		},
		{
			name:  "Mixed case input",
			input: "John Smith <JOHN.SMITH@EXAMPLE.COM>",
			expected: EmailParticipant{
				FirstName: "john",
				LastName:  "smith",
				Email:     "john.smith@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.parseEmailAndName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseParticipants(t *testing.T) {
	svc := &mailService{}
	tests := []struct {
		name     string
		input    string
		expected []EmailParticipant
	}{
		{
			name:  "Single participant",
			input: "John Smith <john@example.com>",
			expected: []EmailParticipant{
				{
					FirstName: "john",
					LastName:  "smith",
					Email:     "john@example.com",
				},
			},
		},
		{
			name:  "Multiple participants",
			input: "John Smith <john@example.com>, Jane Doe <jane@example.com>",
			expected: []EmailParticipant{
				{
					FirstName: "john",
					LastName:  "smith",
					Email:     "john@example.com",
				},
				{
					FirstName: "jane",
					LastName:  "doe",
					Email:     "jane@example.com",
				},
			},
		},
		{
			name:  "Mixed formats",
			input: "john@example.com, Jane Doe <jane@example.com>",
			expected: []EmailParticipant{
				{
					Email: "john@example.com",
				},
				{
					FirstName: "jane",
					LastName:  "doe",
					Email:     "jane@example.com",
				},
			},
		},
		{
			name:     "Empty input",
			input:    "",
			expected: []EmailParticipant{{Email: "", FirstName: "", LastName: ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.parseParticipants(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHeaders(t *testing.T) {
	svc := &mailService{}
	tests := []struct {
		name     string
		headers  map[string]string
		expected EmailHeaders
	}{
		{
			name: "Auto-submitted header",
			headers: map[string]string{
				"Auto-Submitted": "yes",
			},
			expected: EmailHeaders{
				AutoSubmitted: true,
				RawHeaders: map[string]string{
					"Auto-Submitted": "yes",
				},
			},
		},
		{
			name: "Content-Description header",
			headers: map[string]string{
				"Content-Description": "test description",
			},
			expected: EmailHeaders{
				ContentDescription: "test description",
				RawHeaders: map[string]string{
					"Content-Description": "test description",
				},
			},
		},
		{
			name: "Delivery status content type",
			headers: map[string]string{
				"Content-Type": "message/delivery-status",
			},
			expected: EmailHeaders{
				DeliveryStatus: true,
				RawHeaders: map[string]string{
					"Content-Type": "message/delivery-status",
				},
			},
		},
		{
			name:     "Empty headers",
			headers:  map[string]string{},
			expected: EmailHeaders{RawHeaders: map[string]string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.parseHeaders(tt.headers)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Multiple lines",
			input:    "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "Single line",
			input:    "line1",
			expected: []string{"line1"},
		},
		{
			name:     "Empty input",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Multiple spaces",
			input:    "line1   line2\t\tline3",
			expected: []string{"line1", "line2", "line3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractLines(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractEmail(t *testing.T) {
	svc := &mailService{}
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple email",
			input:    "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "Email with angle brackets",
			input:    "<test@example.com>",
			expected: "test@example.com",
		},
		{
			name:     "Email with display name",
			input:    "Test User <test@example.com>",
			expected: "test@example.com",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Invalid email",
			input:    "not-an-email",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.extractEmail(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAllEmails(t *testing.T) {
	svc := &mailService{}
	tests := []struct {
		name     string
		input    EmailParticipants
		expected []string
	}{
		{
			name: "All unique emails",
			input: EmailParticipants{
				From: EmailParticipant{Email: "from@example.com"},
				To:   []EmailParticipant{{Email: "to@example.com"}},
				Cc:   []EmailParticipant{{Email: "cc@example.com"}},
				Bcc:  []EmailParticipant{{Email: "bcc@example.com"}},
			},
			expected: []string{"from@example.com", "to@example.com", "cc@example.com", "bcc@example.com"},
		},
		{
			name: "Duplicate emails",
			input: EmailParticipants{
				From: EmailParticipant{Email: "same@example.com"},
				To:   []EmailParticipant{{Email: "same@example.com"}},
				Cc:   []EmailParticipant{{Email: "same@example.com"}},
				Bcc:  []EmailParticipant{{Email: "different@example.com"}},
			},
			expected: []string{"same@example.com", "different@example.com"},
		},
		{
			name: "Empty emails should be ignored",
			input: EmailParticipants{
				From: EmailParticipant{Email: "from@example.com"},
				To:   []EmailParticipant{{Email: ""}, {Email: "to@example.com"}},
				Cc:   []EmailParticipant{{Email: ""}},
				Bcc:  []EmailParticipant{{Email: ""}},
			},
			expected: []string{"from@example.com", "to@example.com"},
		},
		{
			name:     "Empty participants",
			input:    EmailParticipants{},
			expected: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc.getAllEmails(&tt.input)
			assert.Equal(t, tt.expected, tt.input.AllEmails)
		})
	}
}
