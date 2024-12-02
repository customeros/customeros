package service

import (
	"testing"

	"bou.ke/monkey"
	"github.com/customeros/mailwatcher/blscan"
	"github.com/customeros/mailwatcher/domainage"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/stretchr/testify/assert"
)

func Test_MailboxService_DomainAgePenalty(t *testing.T) {
	tests := []struct {
		name            string
		domain          string
		creationAge     int
		expectedPenalty int
	}{
		{
			name:            "New domain (1 day)",
			domain:          "example.com",
			creationAge:     1,
			expectedPenalty: 75,
		},
		{
			name:            "Week old domain (7 days)",
			domain:          "example.com",
			creationAge:     7,
			expectedPenalty: 60,
		},
		{
			name:            "10 days old domain",
			domain:          "example.com",
			creationAge:     10,
			expectedPenalty: 50,
		},
		{
			name:            "15 days old domain",
			domain:          "example.com",
			creationAge:     15,
			expectedPenalty: 40,
		},
		{
			name:            "30 days old domain",
			domain:          "example.com",
			creationAge:     30,
			expectedPenalty: 30,
		},
		{
			name:            "90 days old domain",
			domain:          "example.com",
			creationAge:     90,
			expectedPenalty: 15,
		},
		{
			name:            "Older domain (>90 days)",
			domain:          "example.com",
			creationAge:     120,
			expectedPenalty: 0,
		},
	}

	mockSpan := &MockSpan{}
	s := &mailboxService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Patch the GetDomainDates function for this test
			patch := monkey.Patch(domainage.GetDomainDates, func(domain string) (domainage.DomainDates, error) {
				return domainage.DomainDates{
					CreationAge: tt.creationAge,
				}, nil
			})
			defer patch.Unpatch()

			result := s.domainAgePenalty(mockSpan, tt.domain)
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
			// Patch the ScanBlacklists function for this test
			patch := monkey.Patch(blscan.ScanBlacklists, func(domain string, scanType string) blscan.BlacklistResults {
				return blscan.BlacklistResults{
					MajorLists:    tt.majorLists,
					MinorLists:    tt.minorLists,
					SpamTrapLists: tt.spamTrapLists,
				}
			})
			defer patch.Unpatch()

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
