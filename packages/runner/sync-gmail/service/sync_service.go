package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/tracing"
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go/log"
	"net/mail"
	"strings"
	"time"
)

const AppSource = "sync-gmail"

type syncService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type SyncService interface {
	GetWhitelistedDomain(domain string, whitelistedDomains []commonEntity.WhitelistDomain) *commonEntity.WhitelistDomain
	GetEmailIdForEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId, email string, whitelistDomain *commonEntity.WhitelistDomain, personalEmailProviderList []commonEntity.PersonalEmailProvider, now time.Time, source string) (string, error)

	BuildEmailsListExcludingPersonalEmails(personalEmailProviderList []commonEntity.PersonalEmailProvider, usernameSource, from string, to []string, cc []string, bcc []string) ([]string, error)

	ConvertToUTC(datetimeStr string) (time.Time, error)
	IsValidEmailSyntax(email string) bool
}

func (s *syncService) GetWhitelistedDomain(domain string, whitelistedDomains []commonEntity.WhitelistDomain) *commonEntity.WhitelistDomain {
	for _, allowedOrganization := range whitelistedDomains {
		if strings.Contains(domain, allowedOrganization.Domain) {
			return &allowedOrganization
		}
	}
	return nil
}

func (s *syncService) BuildEmailsListExcludingPersonalEmails(personalEmailProviderList []commonEntity.PersonalEmailProvider, usernameSource, from string, to []string, cc []string, bcc []string) ([]string, error) {
	var allEmails []string

	if from != "" && from != usernameSource && !hasPersonalEmailProvider(personalEmailProviderList, extractDomain(from)) {
		allEmails = append(allEmails, from)
	}
	allEmails = append(allEmails, from)
	for _, email := range [][]string{to, cc, bcc} {
		for _, email := range email {
			if email != "" && email != usernameSource && !hasPersonalEmailProvider(personalEmailProviderList, extractDomain(email)) {
				allEmails = append(allEmails, email)
			}
		}
	}
	return allEmails, nil
}

func (s *syncService) ConvertToUTC(datetimeStr string) (time.Time, error) {
	var err error

	layouts := []string{
		"2006-01-02T15:04:05Z07:00",
		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
		"Mon, 02 Jan 2006 15:04:05 -0700 (MST)",

		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 -0700",

		"Mon, 2 Jan 2006 15:04:05 +0000 (GMT)",
		"Mon, 02 Jan 2006 15:04:05 +0000 (GMT)",

		"Thu, 2 Jun 2023 03:53:38 -0700 (PDT)",
		"Thu, 02 Jun 2023 03:53:38 -0700 (PDT)",

		"Wed, 2 Sep 2021 13:02:25 GMT",
		"Wed, 02 Sep 2021 13:02:25 GMT",

		"2 Jan 2006 15:04:05 -0700",
		"02 Jan 2006 15:04:05 -0700",
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

func (s *syncService) IsValidEmailSyntax(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "" // Invalid email format
	}
	split := strings.Split(parts[1], ".")
	if len(split) < 2 {
		return parts[1]
	}
	return strings.ToLower(split[len(split)-2] + "." + split[len(split)-1])
}

func hasPersonalEmailProvider(providers []commonEntity.PersonalEmailProvider, domain string) bool {
	for _, provider := range providers {
		if provider.ProviderDomain == domain {
			return true
		}
	}
	return false
}

func (s *syncService) GetEmailIdForEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId, email string, whitelistDomain *commonEntity.WhitelistDomain, personalEmailProviderList []commonEntity.PersonalEmailProvider, now time.Time, source string) (string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.getEmailIdForEmail")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("email", email))

	fromEmailId, err := s.repositories.EmailRepository.GetEmailId(ctx, tenant, email)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve email id for tenant: %v", err)
	}
	if fromEmailId != "" {
		return fromEmailId, nil
	}

	//if it's a personal email, we create just the email node in tenant
	domain := extractDomain(email)
	for _, personalEmailProvider := range personalEmailProviderList {
		if strings.Contains(domain, personalEmailProvider.ProviderDomain) {
			emailId, err := s.repositories.EmailRepository.CreateEmail(ctx, tx, tenant, email, source, AppSource)
			if err != nil {
				return "", fmt.Errorf("unable to create email: %v", err)
			}
			return emailId, nil
		}
	}

	var domainNode *neo4j.Node
	var organizationNode *neo4j.Node
	var organizationId string

	domainNode, err = s.repositories.DomainRepository.GetDomainInTx(ctx, tx, domain)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve domain for tenant: %v", err)
	}

	if domainNode == nil {
		domainNode, err = s.repositories.DomainRepository.CreateDomain(ctx, tx, domain, source, AppSource, now)
		if err != nil {
			return "", fmt.Errorf("unable to create domain: %v", err)
		}
	}

	organizationNode, err = s.repositories.OrganizationRepository.GetOrganizationWithDomain(ctx, tx, tenant, utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*domainNode), "domain"))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve organization for tenant: %v", err)
	}

	if organizationNode == nil {

		var organizationName string

		if whitelistDomain == nil || whitelistDomain.Name == "" {

			//TODO to insert into the allowed organization table with allowed = false t have it for the next time ????
			organizationName, err = s.services.OpenAiService.AskForOrganizationNameByDomain(tenant, interactionEventId, domain)
			if err != nil {
				return "", fmt.Errorf("unable to retrieve organization name for tenant: %v", err)
			}
			if organizationName == "" {
				return "", fmt.Errorf("unable to retrieve organization name for tenant: %v", err)
			}
		} else {
			organizationName = whitelistDomain.Name
		}

		hide := whitelistDomain == nil || !whitelistDomain.Allowed
		organizationNode, err = s.repositories.OrganizationRepository.CreateOrganization(ctx, tx, tenant, organizationName, source, "openline", AppSource, now, hide)
		if err != nil {
			return "", fmt.Errorf("unable to create organization for tenant: %v", err)
		}

		organizationId = utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id")
		domainName := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*domainNode), "domain")
		err = s.repositories.OrganizationRepository.LinkDomainToOrganization(ctx, tx, tenant, domainName, organizationId)
		if err != nil {
			return "", fmt.Errorf("unable to link domain to organization: %v", err)
		}

		_, err := s.repositories.ActionRepository.Create(ctx, tx, tenant, organizationId, entity.ORGANIZATION, entity.ActionCreated, source, AppSource)
		if err != nil {
			return "", fmt.Errorf("unable to create action: %v", err)
		}
	} else {
		organizationId = utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id")
	}

	firstName := ""
	lastname := ""

	//split email address by @ and take the first part to determine first name and last name
	emailParts := strings.Split(email, "@")
	if len(emailParts) > 0 {
		firstPart := emailParts[0]
		nameParts := strings.Split(firstPart, ".")
		if len(nameParts) > 0 {
			firstName = nameParts[0]
			if len(nameParts) > 1 {
				lastname = nameParts[1]
			}
		}
	}

	emailId, err := s.repositories.EmailRepository.CreateContactWithEmailLinkedToOrganization(ctx, tx, tenant, organizationId, email, firstName, lastname, source, AppSource)
	if err != nil {
		return "", fmt.Errorf("unable to create email linked to organization: %v", err)
	}

	return emailId, nil
}

func NewSyncService(cfg *config.Config, repositories *repository.Repositories, services *Services) SyncService {
	return &syncService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
	}
}
