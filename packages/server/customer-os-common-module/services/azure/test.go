package azure

import (
	"testing"
	"time"

	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/stretchr/testify/assert"
)

func Test_buildEmailsRequestURL(t *testing.T) {
	s := &azureService{}

	tests := []struct {
		name     string
		cursor   string
		expected string
	}{
		{
			name:     "with empty cursor",
			cursor:   "",
			expected: "https://graph.microsoft.com/v1.0/me/messages?%24top=100",
		},
		{
			name:     "with existing cursor",
			cursor:   "https://next-page-link",
			expected: "https://next-page-link",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.buildEmailsRequestURL(tt.cursor)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_convertToEmailRawData(t *testing.T) {
	input := MicrosoftRawEmailsResponse{
		Value: []MicrosoftRawEmailResponse{
			{
				Id:                "msg1",
				InternetMessageId: "message1@domain.com",
				SentDateTime:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Subject:           "Test Subject",
				From: struct {
					EmailAddress struct {
						Name    string `json:"name"`
						Address string `json:"address"`
					} `json:"emailAddress"`
				}{
					EmailAddress: struct {
						Name    string `json:"name"`
						Address string `json:"address"`
					}{
						Name:    "Sender Name",
						Address: "sender@test.com",
					},
				},
				ToRecipients: []struct {
					EmailAddress struct {
						Name    string `json:"name"`
						Address string `json:"address"`
					} `json:"emailAddress"`
				}{
					{
						EmailAddress: struct {
							Name    string `json:"name"`
							Address string `json:"address"`
						}{
							Name:    "Recipient Name",
							Address: "recipient@test.com",
						},
					},
				},
				Body: struct {
					ContentType string `json:"contentType"`
					Content     string `json:"content"`
				}{
					ContentType: "html",
					Content:     "<p>Test content</p>",
				},
				ConversationId: "thread1",
			},
		},
	}

	result := convertToEmailRawData(input)

	assert.Len(t, result, 1)
	email := result[0]
	assert.Equal(t, "msg1", email.ProviderMessageId)
	assert.Equal(t, "message1@domain.com", email.MessageId)
	assert.Equal(t, time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), email.Sent)
	assert.Equal(t, "Test Subject", email.Subject)
	assert.Equal(t, "sender@test.com", email.From)
	assert.Equal(t, "Recipient Name <recipient@test.com>", email.To)
	assert.Equal(t, "<p>Test content</p>", email.Html)
	assert.Equal(t, "thread1", email.ThreadId)
}

func Test_buildMailRequest(t *testing.T) {
	s := &azureService{}

	tests := []struct {
		name     string
		input    *postgresEntity.EmailMessage
		expected MailRequest
	}{
		{
			name: "with from name",
			input: &postgresEntity.EmailMessage{
				Subject:  "Test Subject",
				Content:  "<p>Test content</p>",
				From:     "sender@test.com",
				FromName: "Sender Name",
				To:       []string{"recipient@test.com"},
				Cc:       []string{"cc@test.com"},
				Bcc:      []string{"bcc@test.com"},
			},
			expected: MailRequest{
				Subject: "Test Subject",
				From: Recipient{
					EmailAddress: struct {
						Address string `json:"address"`
					}{
						Address: "Sender Name <sender@test.com>",
					},
				},
				Body: struct {
					ContentType string `json:"contentType"`
					Content     string `json:"content"`
				}{
					ContentType: "HTML",
					Content:     "<p>Test content</p>",
				},
				ToRecipients: []Recipient{
					{
						EmailAddress: struct {
							Address string `json:"address"`
						}{
							Address: "recipient@test.com",
						},
					},
				},
				CcRecipients: []Recipient{
					{
						EmailAddress: struct {
							Address string `json:"address"`
						}{
							Address: "cc@test.com",
						},
					},
				},
				BccRecipients: []Recipient{
					{
						EmailAddress: struct {
							Address string `json:"address"`
						}{
							Address: "bcc@test.com",
						},
					},
				},
			},
		},
		// {
		// 	name: "without from name",
		// 	input: &postgresEntity.EmailMessage{
		// 		Subject: "Test Subject",
		// 		Content: "<p>Test content</p>",
		// 		From:    "sender@test.com",
		// 		To:      []string{"recipient@test.com"},
		// 	},
		// 	expected: MailRequest{
		// 		Subject: "Test Subject",
		// 		From: Recipient{
		// 			EmailAddress: struct {
		// 				Address string `json:"address"`
		// 			}{
		// 				Address: "sender@test.com",
		// 			},
		// 		},
		// 		Body: struct {
		// 			ContentType string `json:"contentType"`
		// 			Content     string `json:"content"`
		// 		}{
		// 			ContentType: "HTML",
		// 			Content:     "<p>Test content</p>",
		// 		},
		// 		ToRecipients: []Recipient{
		// 			{
		// 				EmailAddress: struct {
		// 					Address string `json:"address"`
		// 				}{
		// 					Address: "recipient@test.com",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.buildMailRequest(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_buildRecipients(t *testing.T) {
	addresses := []string{"test1@example.com", "test2@example.com"}
	expected := []Recipient{
		{
			EmailAddress: struct {
				Address string `json:"address"`
			}{
				Address: "test1@example.com",
			},
		},
		{
			EmailAddress: struct {
				Address string `json:"address"`
			}{
				Address: "test2@example.com",
			},
		},
	}

	result := buildRecipients(addresses)
	assert.Equal(t, expected, result)
}

func Test_getBodyContent(t *testing.T) {
	body := struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	}{
		ContentType: "html",
		Content:     "<p>Test content</p>",
	}

	tests := []struct {
		name        string
		contentType string
		expected    string
	}{
		{
			name:        "matching content type",
			contentType: "html",
			expected:    "<p>Test content</p>",
		},
		{
			name:        "non-matching content type",
			contentType: "text",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBodyContent(body, tt.contentType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_concatenateEmailAddresses(t *testing.T) {
	recipients := []struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	}{
		{
			EmailAddress: struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			}{
				Name:    "Test User 1",
				Address: "test1@example.com",
			},
		},
		{
			EmailAddress: struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			}{
				Name:    "Test User 2",
				Address: "test2@example.com",
			},
		},
	}

	expected := "Test User 1 <test1@example.com>, Test User 2 <test2@example.com>"
	result := concatenateEmailAddresses(recipients)
	assert.Equal(t, expected, result)
}
