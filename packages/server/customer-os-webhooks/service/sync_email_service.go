package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	repository2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

type SyncEmailService interface {
	SyncEmail(ctx context.Context, email model.EmailData) (SyncResult, error)
	ConvertToUTC(datetimeStr string) (time.Time, error)
	IsValidEmailSyntax(email string) bool
	BuildEmailsListExcludingPersonalEmails(personalEmailProviderList []commonEntity.PersonalEmailProvider, usernameSource, from string, to []string, cc []string, bcc []string) ([]string, error)
	GetEmailIdForEmail(ctx context.Context, tenant string, interactionEventId, email string, whitelistDomain *commonEntity.WhitelistDomain, personalEmailProviderList []commonEntity.PersonalEmailProvider, now time.Time, source string) (string, error)
	GetWhitelistedDomain(domain string, whitelistedDomains []commonEntity.WhitelistDomain) *commonEntity.WhitelistDomain
}

type syncEmailService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
	maxWorkers   int
}

func NewSyncEmailService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) SyncEmailService {
	return &syncEmailService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
		maxWorkers:   services.cfg.ConcurrencyConfig.InteractionEventSyncConcurrency,
	}
}

func (s syncEmailService) SyncEmail(ctx context.Context, emailData model.EmailData) (SyncResult, error) {
	var name string
	var orgSyncResult SyncResult

	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncEmailService.SyncEmails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "emailData", emailData)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	if strings.HasSuffix(emailData.Subject, "• lemwarmup") || strings.HasSuffix(emailData.Subject, "• lemwarm") {
		return SyncResult{Skipped: 1}, nil
	}

	interactionEventId, err := s.repositories.InteractionEventRepository.GetInteractionEventIdByExternalId(ctx, common.GetTenantFromContext(ctx), emailData.ExternalId, emailData.Id)
	if err != nil {
		logrus.Errorf("failed to check if interaction event exists for external id %v for tenant %v :%v", emailData.Id, common.GetTenantFromContext(ctx), err)
		return SyncResult{Failed: 1}, nil
	}

	if interactionEventId == "" {

		now := time.Now().UTC()

		emailSentDate, err := s.ConvertToUTC(emailData.CreatedAtStr)
		if err != nil {
			logrus.Errorf("failed to convert emailData sent date to UTC for emailData with id %v :%v", emailData.Id, err)
			//return entity.ERROR, nil, err
		}

		from := s.extractEmailAddresses(emailData.SentBy)[0]
		to := s.extractEmailAddresses(emailData.SentTo)
		cc := s.extractEmailAddresses(emailData.Cc)
		bcc := s.extractEmailAddresses(emailData.Bcc)
		references := extractLines(emailData.Reference)
		inReplyTo := extractLines(emailData.InReplyTo)

		personalEmailProviderList, err := s.services.CommonServices.CommonRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
		if err != nil {
			logrus.Errorf("failed to get personal emailData provider list: %v", err)
			//return
		}

		whitelistDomainList, err := s.services.CommonServices.CommonRepositories.WhitelistDomainRepository.GetWhitelistDomains(common.GetTenantFromContext(ctx))
		if err != nil {
			logrus.Errorf("failed to get personal emailData provider list: %v", err)
			//return
		}

		allEmailsString, err := s.BuildEmailsListExcludingPersonalEmails(personalEmailProviderList, "", emailData.SentBy, to, cc, bcc)
		if err != nil {
			logrus.Errorf("failed to build emails list: %v", err)
			//return entity.ERROR, nil, err
		}

		if len(allEmailsString) == 0 {
			//reason := "no emails address belongs to a workspace domain"
			//return entity.SKIPPED, &reason, nil
		}

		// Create a map to store the domain counts
		domainCount := make(map[string]int)

		// Iterate through the email addresses
		for _, email := range allEmailsString {
			domain := utils.ExtractDomain(email)
			if domain != "" {
				domainCount[domain]++
			}
		}

		if len(domainCount) > 5 {
			//reason := "more than 5 domains belongs to a workspace domain"
			//return entity.SKIPPED, &reason, nil
		}

		channelData, err := buildEmailChannelData(emailData.Subject, references, inReplyTo)
		if err != nil {
			logrus.Errorf("failed to build emailData channel data for emailData with id %v: %v", emailData.Id, err)
			//return entity.ERROR, nil, err
		}

		sessionId, err := s.services.InteractionSessionService.MergeInteractionSession(ctx, common.GetTenantFromContext(ctx), emailData.ExternalSystem, emailData.SessionDetails, now)

		if err != nil {
			logrus.Errorf("failed merge interaction session for emailData id %v :%v", emailData.Id, err)
			//return entity.ERROR, nil, err
		}
		integrationEvent1 := model.InteractionEventData{
			BaseData:        model.BaseData{CreatedAt: &emailSentDate},
			Content:         "",
			ContentType:     "",
			Channel:         "",
			ChannelData:     *channelData,
			Identifier:      "",
			EventType:       "",
			Hide:            false,
			BelongsTo:       model.BelongsTo{},
			SentBy:          model.InteractionEventParticipant{},
			SentTo:          nil,
			ContactRequired: false,
			ParentRequired:  false,
			SessionDetails:  model.InteractionSessionData{},
		}
		var interactionEvents []model.InteractionEventData
		interactionEvents = append(interactionEvents, integrationEvent1)
		//TODO append emailData to interaction event slice

		_, _ = s.services.InteractionEventService.SyncInteractionEvents(ctx, interactionEvents)
		if err != nil {
			logrus.Errorf("failed merge interaction event for emailData id %v :%v", emailData.Id, err)
			//return entity.ERROR, nil, err
		}

		err = s.repositories.Neo4jRepositories.InteractionEventWriteRepository.LinkInteractionEventToSession(ctx, common.GetTenantFromContext(ctx), interactionEventId, sessionId)
		if err != nil {
			logrus.Errorf("failed to associate interaction event to session for raw emailData id %v :%v", emailData.Id, err)
			//return entity.ERROR, nil, err
		}

		emailidList := []string{}

		var source string
		//TODO check the source outlook or gmail
		//from
		//check if domain exists for tenant by emailData. if so, link the emailData to the user otherwise create a contact and link the emailData to the contact
		fromEmailId, err := s.GetEmailIdForEmail(ctx, common.GetTenantFromContext(ctx), interactionEventId, from, s.GetWhitelistedDomain(utils.ExtractDomain(from), whitelistDomainList), personalEmailProviderList, now, source)

		//here we trigger the sync organization
		orgSyncResult, err = s.createOrganizationDataAndSync(ctx, name, from, emailData)
		if err != nil {
			logrus.Errorf("unable sync org: %v", err)
			//return entity.ERROR, nil, err
		}
		_, err = s.createContactDataAndSync(ctx, name, from, emailData)
		if err != nil {
			logrus.Errorf("unable sync contact: %v", err)
			//return entity.ERROR, nil, err
		}
		if fromEmailId == "" {
			logrus.Errorf("unable to retrieve emailData id for tenant %s and emailData %s", common.GetTenantFromContext(ctx), from)
			//return entity.ERROR, nil, err
		}

		err = s.repositories.Neo4jRepositories.InteractionEventWriteRepository.InteractionEventSentByEmail(ctx, common.GetTenantFromContext(ctx), interactionEventId, fromEmailId)
		if err != nil {
			logrus.Errorf("unable to link emailData to interaction event: %v", err)
			//return entity.ERROR, nil, err
		}
		emailidList = append(emailidList, fromEmailId)

		//to
		for _, toEmail := range to {
			toEmailId, err := s.GetEmailIdForEmail(ctx, common.GetTenantFromContext(ctx), interactionEventId, toEmail, s.GetWhitelistedDomain(utils.ExtractDomain(toEmail), whitelistDomainList), personalEmailProviderList, now, source)
			if err != nil {
				logrus.Errorf("unable to retrieve emailData id for tenant: %v", err)
				//return entity.ERROR, nil, err
			}
			if toEmailId == "" {
				logrus.Errorf("unable to retrieve emailData id for tenant %s and emailData %s", common.GetTenantFromContext(ctx), toEmail)
				//return entity.ERROR, nil, err
			}
			orgSyncResult, err = s.createOrganizationDataAndSync(ctx, name, toEmail, emailData)
			if err != nil {
				logrus.Errorf("unable to sync org: %v", err)
				//return entity.ERROR, nil, err
			}
			_, err = s.createContactDataAndSync(ctx, name, toEmail, emailData)
			if err != nil {
				logrus.Errorf("unable sync contact: %v", err)
				//return entity.ERROR, nil, err
			}

			err = s.repositories.Neo4jRepositories.InteractionEventWriteRepository.InteractionEventSentToEmails(ctx, common.GetTenantFromContext(ctx), interactionEventId, "TO", []string{toEmailId})
			if err != nil {
				logrus.Errorf("unable to link emailData to interaction event: %v", err)
				//return entity.ERROR, nil, err
			}
			emailidList = append(emailidList, toEmailId)
		}

		//cc
		for _, ccEmail := range cc {
			ccEmailId, err := s.GetEmailIdForEmail(ctx, common.GetTenantFromContext(ctx), interactionEventId, ccEmail, s.GetWhitelistedDomain(utils.ExtractDomain(ccEmail), whitelistDomainList), personalEmailProviderList, now, source)
			if err != nil {
				logrus.Errorf("unable to retrieve emailData id for tenant: %v", err)
				//return entity.ERROR, nil, err
			}
			if ccEmailId == "" {
				logrus.Errorf("unable to retrieve emailData id for tenant %s and emailData %s", common.GetTenantFromContext(ctx), ccEmail)
				//return entity.ERROR, nil, err
			}
			orgSyncResult, err = s.createOrganizationDataAndSync(ctx, name, ccEmail, emailData)
			if err != nil {
				logrus.Errorf("unable to sync org: %v", err)
				//return entity.ERROR, nil, err
			}
			_, err = s.createContactDataAndSync(ctx, name, ccEmail, emailData)
			if err != nil {
				logrus.Errorf("unable sync contact: %v", err)
				//return entity.ERROR, nil, err
			}

			err = s.repositories.Neo4jRepositories.InteractionEventWriteRepository.InteractionEventSentToEmails(ctx, common.GetTenantFromContext(ctx), interactionEventId, "CC", []string{ccEmailId})
			if err != nil {
				logrus.Errorf("unable to link emailData to interaction event: %v", err)
				///return entity.ERROR, nil, err
			}
			emailidList = append(emailidList, ccEmailId)
		}

		//bcc
		for _, bccEmail := range bcc {

			bccEmailId, err := s.GetEmailIdForEmail(ctx, common.GetTenantFromContext(ctx), interactionEventId, bccEmail, s.GetWhitelistedDomain(utils.ExtractDomain(bccEmail), whitelistDomainList), personalEmailProviderList, now, source)
			if err != nil {
				logrus.Errorf("unable to retrieve emailData id for tenant: %v", err)
				//return entity.ERROR, nil, err
			}
			if bccEmailId == "" {
				logrus.Errorf("unable to retrieve emailData id for tenant %s and emailData %s", common.GetTenantFromContext(ctx), bccEmail)
				//return entity.ERROR, nil, err
			}
			orgSyncResult, err = s.createOrganizationDataAndSync(ctx, name, bccEmail, emailData)
			if err != nil {
				logrus.Errorf("unable to sync org: %v", err)
				//return entity.ERROR, nil, err
			}
			_, err = s.createContactDataAndSync(ctx, name, bccEmail, emailData)
			if err != nil {
				logrus.Errorf("unable sync contact: %v", err)
				//return entity.ERROR, nil, err
			}

			err = s.repositories.Neo4jRepositories.InteractionEventWriteRepository.InteractionEventSentToEmails(ctx, common.GetTenantFromContext(ctx), interactionEventId, "BCC", []string{bccEmailId})
			if err != nil {
				logrus.Errorf("unable to link emailData to interaction event: %v", err)
				//return entity.ERROR, nil, err
			}

			emailidList = append(emailidList, bccEmailId)
		}

	} else {
		logrus.Infof("interaction event already exists for raw emailData id %v", emailData.Id)
		//reason := "interaction event already exists"
		//return entity.SKIPPED, &reason
	}

	return orgSyncResult, nil
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

func (s *syncEmailService) GetEmailIdForEmail(ctx context.Context, tenant string, interactionEventId, email string, whitelistDomain *commonEntity.WhitelistDomain, personalEmailProviderList []commonEntity.PersonalEmailProvider, now time.Time, source string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncEmailService.GetEmailIdForEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("email", email))

	fromEmailId, err := s.repositories.Neo4jRepositories.EmailReadRepository.GetEmailIdIfExists(ctx, tenant, email)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve email id for tenant: %v", err)
	}
	if fromEmailId != "" {
		return fromEmailId, nil
	}

	//if it's a personal email, we create just the email node in tenant
	domain := utils.ExtractDomain(email)
	for _, personalEmailProvider := range personalEmailProviderList {
		if strings.Contains(domain, personalEmailProvider.ProviderDomain) {
			err = s.repositories.Neo4jRepositories.EmailWriteRepository.CreateEmail(ctx, tenant, email, repository2.EmailCreateFields{
				RawEmail:     email,
				SourceFields: neo4jmodel.Source{Source: source},
			})
			if err != nil {
				return "", fmt.Errorf("unable to create email: %v", err)
			}
			return email, nil
		}
	}

	var domainNode *neo4j.Node
	var emailId string

	domainNode, err = s.repositories.Neo4jRepositories.DomainReadRepository.GetDomain(ctx, domain, common.GetTenantFromContext(ctx))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve domain for tenant: %v", err)
	}

	if domainNode == nil {
		//TODO check for the AppSource
		err = s.repositories.Neo4jRepositories.DomainWriteRepository.CreateDomain(ctx, domain, source, "AppSource", now)
		if err != nil {
			return "", fmt.Errorf("unable to create domain: %v", err)
		}
	}
	//TODO check here if needed
	//orgSyncResult, err := s.services.ContactService.SyncContacts(ctx, organizationsData)
	//
	//emailId, err := s.repositories.CreateContactWithEmailLinkedToOrganization(ctx, tx, tenant, organizationId, email, firstName, lastname, source, AppSource)
	//if err != nil {
	//	return "", fmt.Errorf("unable to create email linked to organization: %v", err)
	//}

	return emailId, nil
}
