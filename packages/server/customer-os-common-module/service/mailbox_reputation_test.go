package service

import (
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/stretchr/testify/assert"
)

func Test_MailboxService_DomainAgePenalty(t *testing.T) {
	tests := []struct {
		name            string
		domainAgeData   string
		expectedPenalty int
	}{
		{
			name:            "New domain (1 day)",
			domainAgeData:   "1",
			expectedPenalty: 75,
		},
		{
			name:            "Week old domain (7 days)",
			domainAgeData:   "7",
			expectedPenalty: 60,
		},
		{
			name:            "10 days old domain",
			domainAgeData:   "10",
			expectedPenalty: 50,
		},
		{
			name:            "15 days old domain",
			domainAgeData:   "15",
			expectedPenalty: 40,
		},
		{
			name:            "30 days old domain",
			domainAgeData:   "30",
			expectedPenalty: 30,
		},
		{
			name:            "90 days old domain",
			domainAgeData:   "90",
			expectedPenalty: 15,
		},
		{
			name:            "Older domain (>90 days)",
			domainAgeData:   "120",
			expectedPenalty: 0,
		},
	}

	s := &mailboxService{}
	mockSpan := &MockSpan{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.domainAgePenalty(mockSpan, tt.domainAgeData)
			assert.Equal(t, tt.expectedPenalty, result)
		})
	}
}

func Test_MailboxService_BlacklistPenaltyPercent(t *testing.T) {
	tests := []struct {
		name            string
		domain          string
		majorLists      int
		minorLists      int
		spamTrapLists   int
		expectedPenalty int
	}{
		{
			name:            "No blacklists",
			domain:          "example.com",
			majorLists:      0,
			minorLists:      0,
			spamTrapLists:   0,
			expectedPenalty: 0,
		},
		{
			name:            "Only major lists",
			domain:          "example.com",
			majorLists:      1,
			minorLists:      0,
			spamTrapLists:   0,
			expectedPenalty: 80,
		},
		{
			name:            "Only minor lists",
			domain:          "example.com",
			majorLists:      0,
			minorLists:      1,
			spamTrapLists:   0,
			expectedPenalty: 10,
		},
		{
			name:            "Only spam trap lists",
			domain:          "example.com",
			majorLists:      0,
			minorLists:      0,
			spamTrapLists:   1,
			expectedPenalty: 20,
		},
		{
			name:            "Mixed lists below 100%",
			domain:          "example.com",
			majorLists:      1,
			minorLists:      1,
			spamTrapLists:   0,
			expectedPenalty: 90,
		},
		{
			name:            "Mixed lists above 100%",
			domain:          "example.com",
			majorLists:      2,
			minorLists:      1,
			spamTrapLists:   1,
			expectedPenalty: 100,
		},
	}

	s := &mailboxService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.blacklistPenaltyPercent(tt.domain)
			assert.Equal(t, tt.expectedPenalty, result)
		})
	}
}

// Simple mock span implementation for testing
type MockSpan struct {
	opentracing.Span
}

func (m *MockSpan) LogFields(fields ...log.Field) {}

func (m *MockSpan) Finish() {}
