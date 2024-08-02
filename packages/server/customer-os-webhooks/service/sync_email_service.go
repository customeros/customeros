package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	repository2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"time"
)

type SyncEmailService interface {
	SyncEmail(ctx context.Context, email model.EmailData) (organizationSync SyncResult, interactionEventSync SyncResult, contactSync SyncResult, err error)
	GetEmailIdForEmail(ctx context.Context, tenant string, email string, personalEmailProviderList []postgresEntity.PersonalEmailProvider, source string) (string, error)
}

type syncEmailService struct {
	log                    logger.Logger
	repositories           *repository.Repositories
	grpcClients            *grpc_client.Clients
	services               *Services
	cache                  *caches.Cache
	maxWorkers             int
	personalEmailProviders []postgresEntity.PersonalEmailProvider
}

func NewSyncEmailService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services, cache *caches.Cache, personalEmailProvidersList []postgresEntity.PersonalEmailProvider) SyncEmailService {
	personalEmailProviders := cache.GetPersonalEmailProviders()
	var err error
	if len(personalEmailProviders) == 0 {
		personalEmailProvidersList, err = repositories.PostgresRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
		if err != nil {
			log.Errorf("error while getting personal email providers: %v", err)
		}
		personalEmailProviders = make([]string, 0)
		for _, personalEmailProvider := range personalEmailProvidersList {
			personalEmailProviders = append(personalEmailProviders, personalEmailProvider.ProviderDomain)
		}
		cache.SetPersonalEmailProviders(personalEmailProviders)
	}
	return &syncEmailService{
		log:                    log,
		repositories:           repositories,
		grpcClients:            grpcClients,
		services:               services,
		maxWorkers:             services.cfg.ConcurrencyConfig.InteractionEventSyncConcurrency,
		cache:                  cache,
		personalEmailProviders: personalEmailProvidersList,
	}
}

func (s syncEmailService) SyncEmail(ctx context.Context, emailData model.EmailData) (organizationSync SyncResult, interactionEventSync SyncResult, contactSync SyncResult, err error) {
	var orgSyncResult, interactionEventSyncResult, contactSyncResult, orgSync, contactsSync SyncResult

	var toOrgSync, toContactsSync SyncResult
	var ccOrgSync, ccContactsSync SyncResult
	var bccOrgSync, bccContactsSync SyncResult

	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncEmailService.SyncEmails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "emailData", emailData)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		reason := fmt.Sprintf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		s.log.Errorf(reason)
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, SyncResult{}, SyncResult{}, errors.ErrTenantNotValid
	}

	if strings.HasSuffix(emailData.Subject, "• lemwarmup") || strings.HasSuffix(emailData.Subject, "• lemwarm") {
		return SyncResult{Skipped: 1}, SyncResult{}, SyncResult{}, nil
	}

	interactionEventId, err := s.repositories.InteractionEventRepository.GetInteractionEventIdByExternalId(ctx, common.GetTenantFromContext(ctx), emailData.ExternalId, emailData.ExternalSystem)
	if err != nil {
		reason := fmt.Sprintf("failed to check if interaction event exists for external id %v for tenant %v :%v", emailData.ExternalId, common.GetTenantFromContext(ctx), err)
		s.log.Error(reason)
		return SyncResult{}, SyncResult{Failed: 1}, SyncResult{}, nil
	}

	if interactionEventId == "" {
		emailSentDate, err := utils.UnmarshalDateTime(emailData.CreatedAtStr)
		if err != nil {
			reason := fmt.Sprintf("failed convert email date to utc %s", err.Error())
			s.log.Error(reason)
		}

		from, to, cc, bcc, inReplyTo := extractEmailData(emailData)

		allEmailsString, err := buildEmailsListExcludingPersonalEmails(s.personalEmailProviders, "", emailData.SentBy.Address, to, cc, bcc)
		if err != nil {
			reason := fmt.Sprintf("failed to build emails list: %v", err)
			s.log.Error(reason)
		}

		if len(allEmailsString) == 0 {
			return SyncResult{}, SyncResult{}, SyncResult{}, nil
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
			reason := "more than 5 domains belongs to a workspace domain"
			s.log.Error(reason)
			return SyncResult{}, SyncResult{}, SyncResult{}, nil
		}

		channelData, err := buildEmailChannelData(emailData.Subject, inReplyTo)
		if err != nil {
			reason := fmt.Sprintf("failed to build emailData channel data for emailData with id %v: %v", emailData.Id, err)
			s.log.Error(reason)
			return SyncResult{Skipped: 1}, SyncResult{}, SyncResult{}, err
		}
		session := model.InteractionSessionData{
			BaseData: model.BaseData{
				CreatedAt:            emailSentDate,
				ExternalSystem:       emailData.ExternalSystem,
				ExternalId:           emailData.ThreadId,
				ExternalSourceEntity: emailData.ExternalSourceEntity,
				SyncId:               emailData.SyncId,
				AppSource:            emailData.AppSource,
			},
			Channel:     "EMAIL",
			ChannelData: *channelData,
			Identifier:  emailData.ThreadId,
			Status:      "ACTIVE",
			Type:        "THREAD",
			Name:        emailData.Subject,
		}

		integrationEvent := model.InteractionEventData{
			BaseData: model.BaseData{
				CreatedAt:            emailSentDate,
				ExternalSystem:       emailData.ExternalSystem,
				ExternalId:           emailData.ExternalId,
				ExternalIdSecond:     emailData.ExternalIdSecond,
				ExternalUrl:          emailData.ExternalUrl,
				ExternalSourceEntity: emailData.ExternalSourceEntity,
				SyncId:               emailData.SyncId,
				AppSource:            emailData.AppSource,
			},
			Content:        emailData.Content,
			ContentType:    emailData.ContentType,
			Channel:        "EMAIL",
			ChannelData:    *channelData,
			Identifier:     emailData.ExternalId,
			Hide:           emailData.Hide,
			SessionDetails: session,
		}
		var interactionEvents []model.InteractionEventData
		interactionEvents = append(interactionEvents, integrationEvent)

		interactionEventSyncResult, err = s.services.InteractionEventService.SyncInteractionEvents(ctx, interactionEvents)
		if err != nil {
			reason := fmt.Sprintf("failed merge interaction event for emailData id %v :%v", emailData.Id, err)
			s.log.Error(reason)
			return SyncResult{}, interactionEventSyncResult, SyncResult{}, nil
		}

		err = s.linkInteractionEventToSessionWithRetry(ctx, &emailData, interactionEventId, session.Identifier)
		if err != nil {
			reason := fmt.Sprintf("failed to associate interaction event to session for raw emailData id %v :%v", emailData.Id, err)
			s.log.Error(reason)
			return SyncResult{}, SyncResult{}, SyncResult{}, nil
		}
		var source string
		if emailData.ExternalSystem == "gmail" {
			source = "GMAIL"
		} else if emailData.ExternalSystem == "outlook" {
			source = "OUTLOOK"
		} else {
			err = fmt.Errorf("unknown emailData source: %s", emailData.ExternalSystem)
			s.log.Error(err.Error())
			return SyncResult{}, SyncResult{}, SyncResult{}, err
		}

		// Process the "from" email
		orgSyncResult, contactSyncResult, err = s.processEmail(ctx, emailData.SentBy.Name, from, emailData, s.personalEmailProviders, source, interactionEventId)
		if err != nil {
			reason := fmt.Sprintf("failed to process emailData for emailData id %v :%v", emailData.Id, err)
			s.log.Error(reason)
			return SyncResult{}, SyncResult{}, SyncResult{}, nil
		}

		fromEmailId, err := s.GetEmailIdForEmail(ctx, common.GetTenantFromContext(ctx), from, s.personalEmailProviders, source)

		if fromEmailId == "" {
			reason := fmt.Sprintf("unable to retrieve emailData id for tenant %s and emailData %s", common.GetTenantFromContext(ctx), from)
			s.log.Error(reason)
			return SyncResult{}, SyncResult{}, SyncResult{}, nil
		}

		// Loop for To
		for i, email := range to {
			name := emailData.SentTo[i].Name
			toOrgSync, toContactsSync, err = s.processEmail(ctx, name, email, emailData, s.personalEmailProviders, source, interactionEventId)
			if err != nil {
				reason := fmt.Sprintf("failed to process emailData for emailData id %v: %v", emailData.Id, err)
				s.log.Error(reason)
				return SyncResult{}, SyncResult{}, SyncResult{}, nil
			}
		}
		// Loop for Cc
		for i, email := range cc {
			name := emailData.Cc[i].Name
			ccOrgSync, ccContactsSync, err = s.processEmail(ctx, name, email, emailData, s.personalEmailProviders, source, interactionEventId)
			if err != nil {
				reason := fmt.Sprintf("failed to process emailData for emailData id %v: %v", emailData.Id, err)
				s.log.Error(reason)
				return SyncResult{}, SyncResult{}, SyncResult{}, nil
			}
		}
		// Loop for Bcc
		for i, email := range bcc {
			name := emailData.Bcc[i].Name
			bccOrgSync, bccContactsSync, err = s.processEmail(ctx, name, email, emailData, s.personalEmailProviders, source, interactionEventId)
			if err != nil {
				reason := fmt.Sprintf("failed to process emailData for emailData id %v: %v", emailData.Id, err)
				s.log.Error(reason)
				return SyncResult{}, SyncResult{}, SyncResult{}, nil
			}
		}

		// Merge results
		orgSync = mergeSyncResults(toOrgSync, ccOrgSync, bccOrgSync, orgSyncResult)
		contactsSync = mergeSyncResults(toContactsSync, ccContactsSync, bccContactsSync, contactSyncResult)

	} else {
		reason := fmt.Sprintf("interaction event already exists for raw emailData id %v", emailData.Id)
		s.log.Info(reason)
		return SyncResult{}, SyncResult{}, SyncResult{}, nil
	}
	return orgSync, interactionEventSyncResult, contactsSync, nil
}

func (s *syncEmailService) GetEmailIdForEmail(ctx context.Context, tenant string, email string, personalEmailProviderList []postgresEntity.PersonalEmailProvider, source string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncEmailService.GetEmailIdForEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	emailId, err := s.repositories.Neo4jRepositories.EmailReadRepository.GetEmailIdIfExists(ctx, tenant, email)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve email id for tenant: %v", err)
	}
	if emailId != "" {
		return emailId, nil
	}

	//if it's a personal email, we create just the email node in tenant
	domain := utils.ExtractDomain(email)
	for _, personalEmailProvider := range personalEmailProviderList {
		if strings.Contains(domain, personalEmailProvider.ProviderDomain) {
			id, err := utils.GenerateUUID()
			if err != nil {
				return "", fmt.Errorf("unable to generate email id: %v", err)
			}
			err = s.repositories.Neo4jRepositories.EmailWriteRepository.CreateEmail(ctx, tenant, id, repository2.EmailCreateFields{
				RawEmail:     email,
				SourceFields: neo4jmodel.Source{Source: source},
			})
			if err != nil {
				return "", fmt.Errorf("unable to create email: %v", err)
			}
			return id, nil
		}
	}
	return emailId, nil
}

func (s *syncEmailService) processEmail(ctx context.Context, name string, email string, emailData model.EmailData, personalEmailProviderList []postgresEntity.PersonalEmailProvider, source string, interactionEventId string) (SyncResult, SyncResult, error) {
	from, to, cc, bcc, _ := extractEmailData(emailData)

	emailId, err := s.GetEmailIdForEmail(ctx, common.GetTenantFromContext(ctx), email, personalEmailProviderList, source)
	if err != nil {
		reason := fmt.Sprintf("unable to retrieve emailData id for tenant: %v", err)
		s.log.Error(reason)
	}
	if emailId == "" {
		reason := fmt.Sprintf("unable to retrieve emailData id for tenant %s and emailData %s", common.GetTenantFromContext(ctx), email)
		s.log.Error(reason)
	}

	orgSyncResult, err := s.createOrganizationDataAndSync(ctx, name, email, emailData)
	if err != nil {
		reason := fmt.Sprintf("unable to sync org: %v", err)
		s.log.Error(reason)
	}

	contactSyncResult, err := s.createContactDataAndSync(ctx, name, email, emailData)
	if err != nil {
		reason := fmt.Sprintf("unable sync contact: %v", err)
		s.log.Error(reason)
	}

	// Set the timeout for waiting
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Define a channel for the timeout
	_ = time.After(timeout)

	select {
	case <-ctx.Done():
		reason := fmt.Sprintf("timeout waiting for operation to complete")
		s.log.Error(reason)
		return SyncResult{}, SyncResult{}, nil
	default:
		var eventType string
		if email == from {
			// Handle "from" case
			eventType = "FROM"
			err = s.repositories.Neo4jRepositories.InteractionEventWriteRepository.InteractionEventSentByEmail(ctx, common.GetTenantFromContext(ctx), interactionEventId, emailId)
		} else {
			// Handle other cases
			if contains(to, email) {
				eventType = "TO"
			} else if contains(cc, email) {
				eventType = "CC"
			} else if contains(bcc, email) {
				eventType = "BCC"
			}
			err = s.repositories.Neo4jRepositories.InteractionEventWriteRepository.InteractionEventSentToEmails(ctx, common.GetTenantFromContext(ctx), interactionEventId, eventType, []string{emailId})
		}
	}
	if err != nil {
		reason := fmt.Sprintf("unable to link emailData to interaction event: %v", err)
		s.log.Error(reason)
		return SyncResult{}, SyncResult{}, err
	}
	return orgSyncResult, contactSyncResult, nil
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
				AppSource:      emailData.AppSource,
				ExternalSystem: emailData.ExternalSystem,
			},
			Name:  name,
			Email: email,
		},
	}

	contactSyncResult, err := s.services.ContactService.SyncContacts(ctx, contactsData)
	return contactSyncResult, err
}

// Define a function to link interaction event to session with retry and timeout
func (s *syncEmailService) linkInteractionEventToSessionWithRetry(ctx context.Context, emailData *model.EmailData, interactionEventId, sessionId string) error {
	// Set the timeout for waiting on node persistence
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Define a function for retry with backoff
	retry := func(ctx context.Context) error {
		err := s.repositories.Neo4jRepositories.InteractionEventWriteRepository.LinkInteractionEventToSession(ctx, common.GetTenantFromContext(ctx), interactionEventId, sessionId)
		return err
	}

	// Use retry with exponential backoff until timeout
	err := retryWithExponentialBackoff(ctx, retry)
	if err != nil {
		reason := fmt.Sprintf("failed to associate interaction event to session for raw emailData id %v: %v", emailData.Id, err)
		s.log.Error(reason)
		return err
	}

	return nil
}

// retry with exponential backoff
func retryWithExponentialBackoff(ctx context.Context, retryFunc func(context.Context) error) error {
	initialDelay := 100 * time.Millisecond
	maxDelay := 2 * time.Second

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := retryFunc(ctx)
			if err == nil {
				return nil
			}

			// Exponential backoff
			delay := initialDelay
			initialDelay *= 2
			if initialDelay > maxDelay {
				initialDelay = maxDelay
			}

			time.Sleep(delay)
		}
	}
}

func mergeSyncResults(results ...SyncResult) SyncResult {
	mergedResult := SyncResult{}

	for _, result := range results {
		mergedResult.Skipped += result.Skipped
		mergedResult.Failed += result.Failed
		mergedResult.Completed += result.Completed
	}

	return mergedResult
}
