package service

import (
	"context"
	"encoding/json"
	"fmt"
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

func CallEventsPlatformGRPCWithRetry[T any](fn func() (T, error)) (T, error) {
	var err error
	var response T
	for attempt := 1; attempt <= constants.MaxRetryGrpcCallWhenUnavailable; attempt++ {
		response, err = fn()
		if err == nil {
			break
		}
		if grpcError, ok := status.FromError(err); ok && (grpcError.Code() == codes.Unavailable || grpcError.Code() == codes.DeadlineExceeded) {
			time.Sleep(utils.BackOffExponentialDelay(attempt))
		} else {
			break
		}
	}
	return response, err
}

type EmailChannelData struct {
	Subject   string   `json:"Subject"`
	InReplyTo []string `json:"InReplyTo"`
	Reference []string `json:"Reference"`
}

func buildEmailChannelData(subject string, references, inReplyTo []string) (*string, error) {
	emailContent := EmailChannelData{
		Subject:   subject,
		InReplyTo: utils.EnsureEmailRfcIds(inReplyTo),
		Reference: utils.EnsureEmailRfcIds(references),
	}
	jsonContent, err := json.Marshal(emailContent)
	if err != nil {
		return nil, err
	}
	jsonContentString := string(jsonContent)

	return &jsonContentString, nil
}

func (s *syncEmailService) GetWhitelistedDomain(domain string, whitelistedDomains []commonEntity.WhitelistDomain) *commonEntity.WhitelistDomain {
	for _, allowedOrganization := range whitelistedDomains {
		if strings.Contains(domain, allowedOrganization.Domain) {
			return &allowedOrganization
		}
	}
	return nil
}

func (s *syncEmailService) createOrganizationDataAndSync(ctx context.Context, name string, domain string, emailData model.EmailData) (SyncResult, error) {
	domainSlice := []string{domain}
	organizationsData := []model.OrganizationData{
		{
			BaseData: model.BaseData{
				AppSource: emailData.AppSource,
				Source:    emailData.ExternalSystem,
			},
			Name:           name,
			Domains:        domainSlice,
			DomainRequired: true,
		},
	}

	orgSyncResult, err := s.services.OrganizationService.SyncOrganizations(ctx, organizationsData)
	return orgSyncResult, err
}

func (s *syncEmailService) createContactDataAndSync(ctx context.Context, name string, email string, emailData model.EmailData) (SyncResult, error) {
	contactsData := []model.ContactData{
		{
			BaseData: model.BaseData{
				AppSource: emailData.AppSource,
				Source:    emailData.ExternalSystem,
			},
			Name:  name,
			Email: email,
		},
	}

	orgSyncResult, err := s.services.ContactService.SyncContacts(ctx, contactsData)
	return orgSyncResult, err
}

func (s *syncEmailService) ConvertToUTC(datetimeStr string) (time.Time, error) {
	var err error

	layouts := []string{
		"2006-01-02T15:04:05Z07:00",

		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",

		"Mon, 2 Jan 2006 15:04:05 MST",

		"Mon, 2 Jan 2006 15:04:05 -0700",

		"Mon, 2 Jan 2006 15:04:05 +0000 (GMT)",

		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",

		"2 Jan 2006 15:04:05 -0700",
	}
	var parsedTime time.Time

	// Try parsing with each layout until successful
	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, datetimeStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse datetime string: %s", datetimeStr)
	}

	return parsedTime.UTC(), nil
}

func (s *syncEmailService) extractEmailAddresses(input string) []string {
	if input == "" {
		return []string{""}
	}
	// Regular expression pattern to match email addresses between <>
	emailPattern := `<(.*?)>`

	emails := make([]string, 0)
	emailAddresses := make([]string, 0)

	if strings.Contains(input, ",") {
		split := strings.Split(input, ",")

		for _, email := range split {
			email = strings.TrimSpace(email)
			email = strings.ToLower(email)
			emails = append(emails, email)
		}
	} else {
		emails = append(emails, input)
	}

	for _, email := range emails {
		email = strings.TrimSpace(email)
		email = strings.ToLower(email)
		if strings.Contains(email, "<") && strings.Contains(email, ">") {
			// Extract email addresses using the regular expression pattern
			re := regexp.MustCompile(emailPattern)
			matches := re.FindAllStringSubmatch(email, -1)

			// Create a map to store unique email addresses
			emailMap := make(map[string]bool)
			for _, match := range matches {
				email := match[1]
				emailMap[email] = true
			}

			// Convert the map keys to an array of email addresses
			for email := range emailMap {
				if s.IsValidEmailSyntax(email) {
					emailAddresses = append(emailAddresses, email)
				}
			}

		} else if s.IsValidEmailSyntax(email) {
			emailAddresses = append(emailAddresses, email)
		}
	}

	if len(emailAddresses) > 0 {
		return emailAddresses
	}

	return []string{input}
}

func (s *syncEmailService) IsValidEmailSyntax(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func hasPersonalEmailProvider(providers []commonEntity.PersonalEmailProvider, domain string) bool {
	for _, provider := range providers {
		if provider.ProviderDomain == domain {
			return true
		}
	}
	return false
}

func (s *syncEmailService) extractSingleEmailAddresses(input string) []string {
	if input == "" {
		return []string{""}
	}
	// Regular expression pattern to match email addresses between <>
	emailPattern := `<(.*?)>`

	emails := make([]string, 0)
	emailAddresses := make([]string, 0)

	if strings.Contains(input, ",") {
		split := strings.Split(input, ",")

		for _, email := range split {
			email = strings.TrimSpace(email)
			email = strings.ToLower(email)
			emails = append(emails, email)
		}
	} else {
		emails = append(emails, input)
	}

	for _, email := range emails {
		email = strings.TrimSpace(email)
		email = strings.ToLower(email)
		if strings.Contains(email, "<") && strings.Contains(email, ">") {
			// Extract email addresses using the regular expression pattern
			re := regexp.MustCompile(emailPattern)
			matches := re.FindAllStringSubmatch(email, -1)

			// Create a map to store unique email addresses
			emailMap := make(map[string]bool)
			for _, match := range matches {
				email := match[1]
				emailMap[email] = true
			}

			// Convert the map keys to an array of email addresses
			for email := range emailMap {
				if s.IsValidEmailSyntax(email) {
					emailAddresses = append(emailAddresses, email)
				}
			}

		} else if s.IsValidEmailSyntax(email) {
			emailAddresses = append(emailAddresses, email)
		}
	}

	if len(emailAddresses) > 0 {
		return emailAddresses
	}

	return []string{input}
}

func extractLines(input string) []string {
	lines := strings.Fields(input)
	return lines
}

func (s *syncEmailService) BuildEmailsListExcludingPersonalEmails(personalEmailProviderList []commonEntity.PersonalEmailProvider, usernameSource, from string, to []string, cc []string, bcc []string) ([]string, error) {
	var allEmails []string

	if from != "" && !hasPersonalEmailProvider(personalEmailProviderList, utils.ExtractDomain(from)) {
		allEmails = append(allEmails, from)
	}
	for _, email := range [][]string{to, cc, bcc} {
		for _, e := range email {
			if e != "" && !hasPersonalEmailProvider(personalEmailProviderList, utils.ExtractDomain(e)) {
				allEmails = append(allEmails, e)
			}
		}
	}
	return allEmails, nil
}

// Helper function to check if an element exists in a slice
func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
