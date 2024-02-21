package service

import (
	"encoding/json"
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

func GetWhitelistedDomain(domain string, whitelistedDomains []commonEntity.WhitelistDomain) *commonEntity.WhitelistDomain {
	for _, allowedOrganization := range whitelistedDomains {
		if strings.Contains(domain, allowedOrganization.Domain) {
			return &allowedOrganization
		}
	}
	return nil
}

func extractEmailAddresses(input string) []string {
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
				if IsValidEmailSyntax(email) {
					emailAddresses = append(emailAddresses, email)
				}
			}

		} else if IsValidEmailSyntax(email) {
			emailAddresses = append(emailAddresses, email)
		}
	}

	if len(emailAddresses) > 0 {
		return emailAddresses
	}

	return []string{input}
}

func IsValidEmailSyntax(email string) bool {
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

func extractSingleEmailAddresses(input string) []string {
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
				if IsValidEmailSyntax(email) {
					emailAddresses = append(emailAddresses, email)
				}
			}

		} else if IsValidEmailSyntax(email) {
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

func buildEmailsListExcludingPersonalEmails(personalEmailProviderList []commonEntity.PersonalEmailProvider, usernameSource, from string, to []string, cc []string, bcc []string) ([]string, error) {
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

func extractEmailData(emailData model.EmailData) (string, []string, []string, []string, []string, []string) {
	// Extract "from" email
	from := extractEmailAddresses(emailData.SentBy)[0]

	// Extract other email addresses
	to := extractEmailAddresses(emailData.SentTo)
	cc := extractEmailAddresses(emailData.Cc)
	bcc := extractEmailAddresses(emailData.Bcc)

	// Extract references
	references := extractLines(emailData.Reference)

	// Extract in-reply-to
	inReplyTo := extractLines(emailData.InReplyTo)

	return from, to, cc, bcc, references, inReplyTo
}
