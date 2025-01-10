package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	pkgerrors "github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

type domains struct {
	personalEmailProviders []string
}

type OrganizationService interface {
	SyncOrganizations(ctx context.Context, organizations []model.OrganizationData) (SyncResult, error)
	GetIdForReferencedOrganization(ctx context.Context, tenant, externalSystem string, org model.ReferencedOrganization) (string, error)
}

type organizationService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
	cache        *caches.Cache
	maxWorkers   int
}

func NewOrganizationService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services, cache *caches.Cache) OrganizationService {
	return &organizationService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
		maxWorkers:   services.cfg.ConcurrencyConfig.OrganizationSyncConcurrency,
		cache:        cache,
	}
}

func (s *organizationService) SyncOrganizations(ctx context.Context, organizations []model.OrganizationData) (SyncResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.SyncOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	// pre-validate organization input before syncing
	for _, org := range organizations {
		if org.ExternalSystem == "" && org.Source == "" {
			tracing.TraceErr(span, errors.ErrMissingExternalSystem)
			return SyncResult{}, errors.ErrMissingExternalSystem
		}
		if org.ExternalSystem != "" {
			if !neo4jentity.IsValidDataSource(strings.ToLower(org.ExternalSystem)) {
				tracing.TraceErr(span, errors.ErrExternalSystemNotAccepted, log.String("externalSystem", org.ExternalSystem))
				return SyncResult{}, errors.ErrExternalSystemNotAccepted
			}
		}
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	workerLimit := make(chan struct{}, s.maxWorkers)

	syncMutex := &sync.Mutex{}
	statusesMutex := &sync.Mutex{}
	syncDate := utils.Now()
	var statuses []SyncStatus

	personalEmailProviders := s.cache.GetPersonalEmailProviders()
	if len(personalEmailProviders) == 0 {
		personalEmailProviderEntities, err := s.repositories.PostgresRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
		if err != nil {
			s.log.Errorf("error while getting personal email providers: %v", err)
		}
		personalEmailProviders = make([]string, 0)
		for _, personalEmailProvider := range personalEmailProviderEntities {
			personalEmailProviders = append(personalEmailProviders, personalEmailProvider.ProviderDomain)
		}
		s.cache.SetPersonalEmailProviders(personalEmailProviders)
	}

	controlDomains := &domains{
		personalEmailProviders: personalEmailProviders,
	}

	// Sync all organizations concurrently
	for _, organizationData := range organizations {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return SyncResult{}, ctx.Err()
		default:
		}

		// Acquire a worker slot
		workerLimit <- struct{}{}
		wg.Add(1)

		go func(organizationData model.OrganizationData) {
			defer wg.Done()
			defer func() {
				// Release the worker slot when done
				<-workerLimit
			}()

			result := s.syncOrganization(ctx, syncMutex, organizationData, syncDate, controlDomains)
			statusesMutex.Lock()
			statuses = append(statuses, result)
			statusesMutex.Unlock()
		}(organizationData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), organizations[0].ExternalSystem,
		organizations[0].AppSource, "organization", syncDate, statuses)

	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}

func (s *organizationService) syncOrganization(ctx context.Context, syncMutex *sync.Mutex, orgInput model.OrganizationData, syncDate time.Time, controlDomains *domains) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.syncOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagExternalSystem, orgInput.ExternalSystem)
	span.SetTag(tracing.SpanTagExternalId, orgInput.ExternalId)
	span.LogFields(log.Object("syncDate", syncDate))
	tracing.LogObjectAsJson(span, "orgInput", orgInput)

	tenant := common.GetTenantFromContext(ctx)
	appSource := utils.StringFirstNonEmpty(orgInput.AppSource, constants.AppSourceCustomerOsWebhooks)
	var failedSync = false
	var reason = ""
	orgInput.Normalize()

	// Check if organization sync should be skipped
	if orgInput.Skip {
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(orgInput.SkipReason)
	}

	// remove any domain for sub organizations, as they are not supported
	if orgInput.IsSubOrg() {
		orgInput.Domains = []string{}
	} else {
		// prepare domains for organization
		orgDomains := make([]string, 0)
		for _, domainInput := range orgInput.Domains {
			orgDomains = append(orgDomains, utils.ExtractDomain(domainInput))
		}
		primaryDomainFromWebsite, _ := s.services.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, orgInput.Website)
		if primaryDomainFromWebsite != "" {
			orgDomains = append(orgDomains, primaryDomainFromWebsite)
		}
		orgInput.Domains = orgDomains
		orgInput.NormalizeDomains()
	}

	// Merge external system neo4j node
	if orgInput.ExternalSystem != "" {
		err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, orgInput.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", orgInput.ExternalSystem, tenant, err.Error())
			s.log.Error(reason)
			span.LogFields(log.String("output", "failed"))
			return NewFailedSyncStatus(reason)
		}
	}

	// Remove personal email provider domains from organization domains
	nonPersonalEmailProviderDomains := make([]string, 0)
	for _, domain := range orgInput.Domains {
		if !controlDomains.isPersonalEmailProvider(domain) {
			nonPersonalEmailProviderDomains = append(nonPersonalEmailProviderDomains, domain)
		}
	}
	orgInput.Domains = nonPersonalEmailProviderDomains

	// Check if organization should be skipped due to missing domain
	if orgInput.DomainRequired && !orgInput.IsSubOrg() && !orgInput.HasDomains() {
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus("Missing domain while required")
	}

	// Use fallback name if applicable
	if orgInput.Name == "" && orgInput.FallbackName != "" && !orgInput.HasDomains() {
		orgInput.Name = orgInput.FallbackName
	}

	// Lock organization creation
	syncMutex.Lock()
	defer syncMutex.Unlock()
	// Check if organization already exists
	organizationId, err := s.repositories.OrganizationRepository.GetMatchedOrganizationId(ctx, tenant, orgInput.ExternalSystem, orgInput.ExternalId, orgInput.CustomerOsId, orgInput.Domains)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched organization with external reference %s for tenant %s :%s", orgInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}
	if !failedSync {
		matchingOrganizationExists := organizationId != ""
		span.LogFields(log.Bool("found matching organization", matchingOrganizationExists))

		if orgInput.UpdateOnly {
			if !matchingOrganizationExists {
				span.LogFields(log.String("output", "skipped"))
				return NewSkippedSyncStatus("Update only flag enabled and no matching organization found")
			}
		}

		organizationDataFields := data_fields.OrganizationFields{
			AppSource: utils.StringPtr(appSource),
			Source:    utils.StringPtr(utils.StringFirstNonEmpty(orgInput.ExternalSystem, orgInput.Source)),
		}
		if !matchingOrganizationExists {
			organizationDataFields.Name = utils.StringPtr(orgInput.Name)
			organizationDataFields.Description = utils.StringPtr(orgInput.Description)
			organizationDataFields.Website = utils.StringPtr(orgInput.Website)
			organizationDataFields.Industry = utils.StringPtr(orgInput.Industry)
			organizationDataFields.IsPublic = utils.BoolPtr(orgInput.IsPublic)
			organizationDataFields.Employees = utils.Int64Ptr(orgInput.Employees)
			organizationDataFields.Market = utils.StringPtr(orgInput.Market)
			organizationDataFields.ValueProposition = utils.StringPtr(orgInput.ValueProposition)
			organizationDataFields.LastFundingRound = utils.StringPtr(orgInput.LastFundingRound)
			organizationDataFields.LastFundingAmount = utils.StringPtr(orgInput.LastFundingAmount)
			organizationDataFields.Note = utils.StringPtr(orgInput.Note)
			organizationDataFields.ReferenceId = utils.StringPtr(orgInput.ReferenceId)
			organizationDataFields.LogoUrl = utils.StringPtr(orgInput.LogoUrl)
			organizationDataFields.YearFounded = orgInput.YearFounded
			organizationDataFields.Headquarters = utils.StringPtr(orgInput.Headquarters)
			organizationDataFields.EmployeeGrowthRate = utils.StringPtr(orgInput.EmployeeGrowthRate)
			organizationDataFields.LeadSource = utils.StringPtr(utils.StringFirstNonEmpty(orgInput.ExternalSystem, orgInput.Source))
			if orgInput.IsCustomer {
				organizationDataFields.Relationship = utils.ToPtr(neo4jenum.OrganizationRelationshipCustomer)
			} else {
				if !matchingOrganizationExists {
					organizationDataFields.Stage = utils.ToPtr(neo4jenum.Trial)
					organizationDataFields.Relationship = utils.ToPtr(neo4jenum.OrganizationRelationshipProspect)
				}
			}
		}

		if orgInput.ExternalSystem != "" {
			organizationDataFields.ExternalSystem = &neo4jmodel.ExternalSystem{
				ExternalSystemId: orgInput.ExternalSystem,
				ExternalId:       orgInput.ExternalId,
				ExternalUrl:      orgInput.ExternalUrl,
				ExternalIdSecond: orgInput.ExternalIdSecond,
				ExternalSource:   orgInput.ExternalSourceEntity,
				SyncDate:         &syncDate,
			}
		}
		var orgIdPtr *string
		if organizationId != "" {
			orgIdPtr = &organizationId
		}

		savedOrgId, err := s.services.CommonServices.OrganizationService.Save(ctx, nil, orgIdPtr, organizationDataFields)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, pkgerrors.Wrap(err, "failed to save organization"))
			reason = fmt.Sprintf("failed to save organization  with external reference %s for tenant %s :%s", orgInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
		organizationId = savedOrgId
		orgInput.Id = organizationId
		span.LogFields(log.String("organizationId", organizationId))
	}
	if !failedSync && orgInput.HasDomains() {
		for _, domain := range orgInput.Domains {
			//check if the domain is already linked to an organization. If the domain is already linked, skip the link operation
			domainInUse, err := s.repositories.OrganizationRepository.IsDomainUsedByOrganization(ctx, tenant, domain, organizationId)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("error while checking if domain is linked to organization: %v", err.Error())
				continue
			}
			if !domainInUse {
				_, err = s.services.CommonServices.OrganizationService.LinkWithDomain(ctx, nil, organizationId, domain)
				if err != nil {
					tracing.TraceErr(span, pkgerrors.Wrapf(err, "failed to link domain %s with organization %s", domain, organizationId))
				}
			}
		}
	}
	if !failedSync && orgInput.IsSubOrg() {
		parentOrganizationId, _ := s.GetIdForReferencedOrganization(ctx, tenant, orgInput.ExternalSystem, orgInput.ParentOrganization.Organization)
		if parentOrganizationId != "" {
			err = s.services.CommonServices.OrganizationService.AddParentOrganization(ctx, nil, parentOrganizationId, organizationId, orgInput.ParentOrganization.Type)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("Failed to link with parent for organization %s: %s", organizationId, err.Error())
				s.log.Error(reason)
			}
		}
	}
	if !failedSync {
		if orgInput.HasEmail() {
			_, err = s.services.CommonServices.EmailService.Merge(ctx, nil, tenant,
				commonservice.EmailFields{
					Email:     orgInput.Email,
					AppSource: orgInput.AppSource,
					Source:    neo4jentity.DecodeDataSource(orgInput.ExternalSystem),
					Primary:   true,
				},
				&commonservice.LinkWith{
					Type: commonmodel.ORGANIZATION,
					Id:   organizationId,
				})
			if err != nil {
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("Failed to create and link email address %s with organization %s: %s", orgInput.Email, organizationId, err.Error())
				failedSync = true
			}
		}

		if orgInput.HasPhoneNumbers() {
			for _, phoneNumberDtls := range orgInput.PhoneNumbers {
				// Create or update phone number
				phoneNumberId, err := s.services.CommonServices.PhoneNumberService.Merge(ctx, phoneNumberDtls.Number, neo4jentity.DecodeDataSource(orgInput.ExternalSystem))
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("Failed to create phone number %s for organization %s: %s", phoneNumberDtls.Number, organizationId, err.Error())
					s.log.Error(reason)
				}
				// Link phone number to contact
				if phoneNumberId != "" {
					err = s.services.CommonServices.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithOrganization(ctx, tenant, organizationId, phoneNumberId, phoneNumberDtls.Label, phoneNumberDtls.Primary)
					if err != nil {
						failedSync = true
						tracing.TraceErr(span, err, log.String("method", "LinkWithOrganization"))
						reason = fmt.Sprintf("Failed to link phone number %s with organization %s: %s", phoneNumberDtls.Number, organizationId, err.Error())
						s.log.Error(reason)
					}
				}
			}
		}

		syncLocation := false // skip location sync for now
		if orgInput.HasLocation() && syncLocation && !failedSync {
			// Create or update location
			locationId, err := s.repositories.LocationRepository.GetMatchedLocationIdForOrganizationBySource(ctx, organizationId, orgInput.ExternalSystem)
			if err != nil {
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("Failed to get matched location for organization %s: %s", organizationId, err.Error())
				failedSync = true
				s.log.Error(reason)
			}
			if locationId == "" {
				locationId, err = s.services.CommonServices.LocationService.Create(ctx, nil,
					data_fields.LocationFields{
						Source:   utils.StringPtr(orgInput.ExternalSystem),
						Name:     orgInput.LocationName,
						Country:  orgInput.Country,
						Region:   orgInput.Region,
						Locality: orgInput.Locality,
						Address:  orgInput.Address,
						Address2: orgInput.Address2,
						Zip:      orgInput.Zip,
					}, &commonservice.LinkWith{
						Type: commonmodel.ORGANIZATION,
						Id:   organizationId,
					})
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("Failed to create location for organization %s: %s", organizationId, err.Error())
					s.log.Error(reason)
				}
			}
		}

		if !orgInput.HasSocials() {
			for _, social := range orgInput.Socials {
				// Link social to organization
				_, err = s.services.CommonServices.SocialService.AddSocialToEntity(ctx,
					nil,
					commonservice.LinkWith{
						Id:   organizationId,
						Type: commonmodel.ORGANIZATION,
					}, neo4jentity.SocialEntity{
						Url:       social.URL,
						Source:    neo4jentity.DecodeDataSource(orgInput.ExternalSystem),
						AppSource: appSource,
					})
				if err != nil {
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("Failed to link social %s with organization %s: %s", social.URL, organizationId, err.Error())
					s.log.Error(reason)
				}
			}
		}
	}

	span.LogFields(log.Bool("failedSync", failedSync))
	if failedSync {
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}
	span.LogFields(log.String("output", "success"))
	return NewSuccessfulSyncStatus()
}

func (d domains) isPersonalEmailProvider(domain string) bool {
	for _, v := range d.personalEmailProviders {
		if strings.ToLower(domain) == strings.ToLower(v) {
			return true
		}
	}
	return false
}

func (s *organizationService) GetIdForReferencedOrganization(ctx context.Context, tenant, externalSystemId string, org model.ReferencedOrganization) (string, error) {
	if !org.Available() {
		return "", nil
	}

	if org.ReferencedById() {
		return s.repositories.OrganizationRepository.GetOrganizationIdById(ctx, tenant, org.Id)
	} else if org.ReferencedByExternalId() {
		return s.repositories.OrganizationRepository.GetOrganizationIdByExternalId(ctx, tenant, org.ExternalId, externalSystemId)
	} else if org.ReferencedByDomain() {
		return s.repositories.OrganizationRepository.GetOrganizationIdByDomain(ctx, tenant, org.Domain)
	}
	return "", nil
}
