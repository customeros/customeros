package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go/log"
	"net/mail"
	"strings"
	"time"
)

const AppSource = "sync-email"

type syncService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type SyncService interface {
	GetEmailIdForEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId, email string, now time.Time, source string) (string, error)

	BuildEmailsListExcludingPersonalEmails(usernameSource, from string, to []string, cc []string, bcc []string) ([]string, error)

	ConvertToUTC(datetimeStr string) (time.Time, error)
	IsValidEmailSyntax(email string) bool
}

func (s *syncService) BuildEmailsListExcludingPersonalEmails(usernameSource, from string, to []string, cc []string, bcc []string) ([]string, error) {
	var allEmails []string

	if from != "" && !hasPersonalEmailProvider(s.services.Cache.GetPersonalEmailProviders(), utils.ExtractDomain(from)) {
		allEmails = append(allEmails, from)
	}
	for _, email := range [][]string{to, cc, bcc} {
		for _, e := range email {
			if e != "" && !hasPersonalEmailProvider(s.services.Cache.GetPersonalEmailProviders(), utils.ExtractDomain(e)) {
				allEmails = append(allEmails, e)
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

func (s *syncService) IsValidEmailSyntax(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func hasPersonalEmailProvider(providers []string, domain string) bool {
	for _, provider := range providers {
		if provider == domain {
			return true
		}
	}
	return false
}

func (s *syncService) GetEmailIdForEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId, email string, now time.Time, source string) (string, error) {
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
	domain := utils.ExtractDomain(email)
	if domain == "" {
		return "", fmt.Errorf("unable to extract domain from email: %s", email)
	}
	for _, personalEmailProvider := range s.services.Cache.GetPersonalEmailProviders() {
		if strings.Contains(domain, personalEmailProvider) {
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

		organizationName := domain
		hide := false
		relationship := neo4jenum.Prospect.String()
		stage := neo4jenum.Lead.String()
		leadSource := ""

		if source == "gmail" {
			leadSource = "Gmail"
		} else if source == "outlook" {
			leadSource = "Outlook"
		} else if source == "mailstack" {
			leadSource = "Mailstack"
			stage = neo4jenum.Target.String()
		} else {
			leadSource = "Email"
		}

		organizationNode, err = s.repositories.OrganizationRepository.CreateOrganization(ctx, tx, tenant, organizationName, relationship, stage, leadSource, source, "openline", AppSource, now, hide)
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
