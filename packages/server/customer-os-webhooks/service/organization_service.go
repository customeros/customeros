package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	comentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	commongrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	orggrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"sync"
	"time"
)

const maxWorkersOrganizationSync = 4

type domains struct {
	whitelistDomains       []comentity.WhitelistDomain
	personalEmailProviders []comentity.PersonalEmailProvider
}

type OrganizationService interface {
	SyncOrganizations(ctx context.Context, organizations []model.OrganizationData) error
	GetIdForReferencedOrganization(ctx context.Context, tenant, externalSystem string, org model.ReferencedOrganization) (string, error)
	mapDbNodeToOrganizationEntity(dbNode dbtype.Node) *entity.OrganizationEntity
}

type organizationService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewOrganizationService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) OrganizationService {
	return &organizationService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *organizationService) SyncOrganizations(ctx context.Context, organizations []model.OrganizationData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.SyncOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return errors.ErrTenantNotValid
	}

	// pre-validate organization input before syncing
	for _, org := range organizations {
		if org.ExternalSystem == "" {
			return errors.ErrMissingExternalSystem
		}
		if !entity.IsValidDataSource(strings.ToLower(org.ExternalSystem)) {
			return errors.ErrExternalSystemNotAccepted
		}
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	workerLimit := make(chan struct{}, maxWorkersOrganizationSync)

	syncMutex := &sync.Mutex{}
	syncDate := utils.Now()
	var statuses []SyncStatus

	whitelistDomains, err := s.repositories.CommonRepositories.WhitelistDomainRepository.GetWhitelistDomains(common.GetTenantFromContext(ctx))
	if err != nil {
		s.log.Errorf("error while getting whitelist domains: %v", err)
		whitelistDomains = make([]comentity.WhitelistDomain, 0)
	}
	personalEmailProviders, err := s.repositories.CommonRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
	if err != nil {
		s.log.Errorf("error while getting personal email providers: %v", err)
		personalEmailProviders = make([]comentity.PersonalEmailProvider, 0)
	}
	controlDomains := &domains{
		whitelistDomains:       whitelistDomains,
		personalEmailProviders: personalEmailProviders,
	}

	// Sync all organizations concurrently
	for _, organizationData := range organizations {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Continue with Slack sync
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
			statuses = append(statuses, result)
		}(organizationData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), organizations[0].ExternalSystem,
		organizations[0].AppSource, "organization", syncDate, statuses)

	return nil
}

func (s *organizationService) syncOrganization(ctx context.Context, syncMutex *sync.Mutex, orgInput model.OrganizationData, syncDate time.Time, controlDomains *domains) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.syncOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", orgInput.ExternalSystem), log.Object("orgInput", orgInput))

	tenant := common.GetTenantFromContext(ctx)
	var failedSync = false
	var reason = ""
	orgInput.Normalize()

	// clear domains for sub organizations
	if orgInput.IsSubOrg() {
		orgInput.Domains = []string{}
	} else {
		// prepare domains for organization
		orgDomains := make([]string, 0)
		for _, domainInput := range orgInput.Domains {
			orgDomains = append(orgDomains, utils.ExtractDomain(domainInput))
		}
		orgDomains = append(orgDomains, utils.ExtractDomain(orgInput.Website))
		orgInput.Domains = orgDomains
		orgInput.NormalizeDomains()
	}

	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, orgInput.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", orgInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		return NewFailedSyncStatus(reason)
	}

	// Check if organization sync should be skipped
	if orgInput.Skip {
		span.LogFields(log.Bool("skippedSync", true))
		return NewSkippedSyncStatus(orgInput.SkipReason)
	}

	nonPersonalEmailProviderDomains := make([]string, 0)
	for _, domain := range orgInput.Domains {
		if !controlDomains.isPersonalEmailProvider(domain) {
			nonPersonalEmailProviderDomains = append(nonPersonalEmailProviderDomains, domain)
		}
	}
	orgInput.Domains = nonPersonalEmailProviderDomains
	if orgInput.DomainRequired && !orgInput.IsSubOrg() && !orgInput.HasDomains() {
		span.LogFields(log.Bool("skippedSync", true))
		return NewSkippedSyncStatus("Missing domain while required")
	}
	orgHasWhitelistedDomain := false
	for _, domain := range orgInput.Domains {
		if controlDomains.isWhitelistedDomain(domain) {
			orgHasWhitelistedDomain = true
		}
	}

	// Lock organization creation
	syncMutex.Lock()
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

		// Create new organization id if not found
		organizationId = utils.NewUUIDIfEmpty(organizationId)
		orgInput.Id = organizationId
		span.LogFields(log.String("organizationId", organizationId))

		// Create or update organization
		_, err = s.grpcClients.OrganizationClient.UpsertOrganization(ctx, &orggrpc.UpsertOrganizationGrpcRequest{
			Tenant:            tenant,
			Id:                organizationId,
			LoggedInUserId:    "",
			Name:              orgInput.Name,
			Description:       orgInput.Description,
			Website:           orgInput.Website,
			Industry:          orgInput.Industry,
			IsPublic:          orgInput.IsPublic,
			IsCustomer:        orgInput.IsCustomer,
			Employees:         orgInput.Employees,
			Market:            orgInput.Market,
			CreatedAt:         utils.ConvertTimeToTimestampPtr(orgInput.CreatedAt),
			UpdatedAt:         utils.ConvertTimeToTimestampPtr(orgInput.UpdatedAt),
			SubIndustry:       orgInput.SubIndustry,
			IndustryGroup:     orgInput.IndustryGroup,
			TargetAudience:    orgInput.TargetAudience,
			ValueProposition:  orgInput.ValueProposition,
			LastFundingRound:  orgInput.LastFundingRound,
			LastFundingAmount: orgInput.LastFundingAmount,
			Hide:              !(orgHasWhitelistedDomain || orgInput.Whitelisted),
			Note:              orgInput.Note,
			ReferenceId:       orgInput.ReferenceId,
			IgnoreEmptyFields: false,
			SourceFields: &commongrpc.SourceFields{
				Source:    orgInput.ExternalSystem,
				AppSource: utils.StringFirstNonEmpty(orgInput.AppSource, constants.AppSourceCustomerOsWebhooks),
			},
			ExternalSystemFields: &commongrpc.ExternalSystemFields{
				ExternalSystemId: orgInput.ExternalSystem,
				ExternalId:       orgInput.ExternalId,
				ExternalUrl:      orgInput.ExternalUrl,
				ExternalIdSecond: orgInput.ExternalIdSecond,
				ExternalSource:   orgInput.ExternalSourceEntity,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
		})
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed sending event to upsert organization  with external reference %s for tenant %s :%s", orgInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
		// Wait for organization to be created in neo4j
		if !failedSync && !matchingOrganizationExists {
			for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
				organization, findErr := s.repositories.OrganizationRepository.GetById(ctx, tenant, organizationId)
				if organization != nil && findErr == nil {
					break
				}
				time.Sleep(time.Duration(i*constants.TimeoutIntervalMs) * time.Millisecond)
			}
		}
	}
	if !failedSync && orgInput.HasDomains() {
		for _, domain := range orgInput.Domains {
			_, err = s.grpcClients.OrganizationClient.LinkDomainToOrganization(ctx, &orggrpc.LinkDomainToOrganizationGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				LoggedInUserId: "",
				OrganizationId: organizationId,
				Domain:         domain,
			})
			if err != nil {
				tracing.TraceErr(span, err)
			}
		}
	}
	syncMutex.Unlock()

	if !failedSync && orgInput.IsSubOrg() {
		parentOrganizationId, _ := s.GetIdForReferencedOrganization(ctx, tenant, orgInput.ExternalSystem, orgInput.ParentOrganization.Organization)
		if parentOrganizationId != "" {
			_, err = s.grpcClients.OrganizationClient.AddParentOrganization(ctx, &orggrpc.AddParentOrganizationGrpcRequest{
				Tenant:               common.GetTenantFromContext(ctx),
				OrganizationId:       organizationId,
				ParentOrganizationId: parentOrganizationId,
				Type:                 orgInput.ParentOrganization.Type,
			})
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
			// Create or update email
			emailId, err := s.services.EmailService.CreateEmail(ctx, orgInput.Email, orgInput.ExternalSystem, orgInput.AppSource)
			if err != nil {
				tracing.TraceErr(span, err)
				failedSync = true
				reason = fmt.Sprintf("Failed to create email address %s for organization %s: %s", orgInput.Email, organizationId, err.Error())
				s.log.Error(reason)
			}
			// Link email to organization
			if emailId != "" {
				_, err = s.grpcClients.OrganizationClient.LinkEmailToOrganization(ctx, &orggrpc.LinkEmailToOrganizationGrpcRequest{
					Tenant:         common.GetTenantFromContext(ctx),
					OrganizationId: organizationId,
					EmailId:        emailId,
				})
				if err != nil {
					tracing.TraceErr(span, err)
					failedSync = true
					reason = fmt.Sprintf("Failed to link email address %s with organization %s: %s", orgInput.Email, organizationId, err.Error())
					s.log.Error(reason)
				}
			}
		}

		if orgInput.HasPhoneNumbers() {
			for _, phoneNumberDtls := range orgInput.PhoneNumbers {
				// Create or update phone number
				phoneNumberId, err := s.services.PhoneNumberService.CreatePhoneNumber(ctx, phoneNumberDtls.Number, orgInput.ExternalSystem, orgInput.AppSource)
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("Failed to create phone number %s for organization %s: %s", phoneNumberDtls.Number, organizationId, err.Error())
					s.log.Error(reason)
				}
				// Link phone number to organization
				if phoneNumberId != "" {
					_, err = s.grpcClients.OrganizationClient.LinkPhoneNumberToOrganization(ctx, &orggrpc.LinkPhoneNumberToOrganizationGrpcRequest{
						Tenant:         common.GetTenantFromContext(ctx),
						OrganizationId: organizationId,
						PhoneNumberId:  phoneNumberId,
						Primary:        phoneNumberDtls.Primary,
						Label:          phoneNumberDtls.Label,
					})
					if err != nil {
						failedSync = true
						tracing.TraceErr(span, err)
						reason = fmt.Sprintf("Failed to link phone number %s for organization %s: %s", phoneNumberDtls.Number, organizationId, err.Error())
						s.log.Error(reason)
					}
				}
			}
		}

		if orgInput.HasLocation() {
			// Create or update location
			locationId, err := s.repositories.LocationRepository.GetMatchedLocationIdForOrganizationBySource(ctx, organizationId, orgInput.ExternalSystem)
			if err != nil {
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("Failed to get matched location for organization %s: %s", organizationId, err.Error())
				failedSync = true
				s.log.Error(reason)
			}
			if !failedSync {
				locationId, err = s.services.LocationService.CreateLocation(ctx, locationId, orgInput.ExternalSystem, orgInput.AppSource,
					orgInput.LocationName, orgInput.Country, orgInput.Region, orgInput.Locality, "", orgInput.Address, orgInput.Address2, orgInput.Zip, "")
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("Failed to create location for organization %s: %s", organizationId, err.Error())
					s.log.Error(reason)
				}
			}

			// Link location to organization
			if locationId != "" {
				_, err = s.grpcClients.OrganizationClient.LinkLocationToOrganization(ctx, &orggrpc.LinkLocationToOrganizationGrpcRequest{
					Tenant:         common.GetTenantFromContext(ctx),
					OrganizationId: organizationId,
					LocationId:     locationId,
				})
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("Failed to link location %s with organization %s: %s", locationId, organizationId, err.Error())
					s.log.Error(reason)
				}
			}
		}
	}

	span.LogFields(log.Bool("failedSync", failedSync))
	if failedSync {
		return NewFailedSyncStatus(reason)
	}
	return NewSuccessfulSyncStatus()
}

func (s *organizationService) mapDbNodeToOrganizationEntity(dbNode dbtype.Node) *entity.OrganizationEntity {
	props := utils.GetPropsFromNode(dbNode)
	output := entity.OrganizationEntity{
		ID:                utils.GetStringPropOrEmpty(props, "id"),
		CustomerOsId:      utils.GetStringPropOrEmpty(props, "customerOsId"),
		ReferenceId:       utils.GetStringPropOrEmpty(props, "referenceId"),
		Name:              utils.GetStringPropOrEmpty(props, "name"),
		Description:       utils.GetStringPropOrEmpty(props, "description"),
		Website:           utils.GetStringPropOrEmpty(props, "website"),
		Industry:          utils.GetStringPropOrEmpty(props, "industry"),
		IndustryGroup:     utils.GetStringPropOrEmpty(props, "industryGroup"),
		SubIndustry:       utils.GetStringPropOrEmpty(props, "subIndustry"),
		TargetAudience:    utils.GetStringPropOrEmpty(props, "targetAudience"),
		ValueProposition:  utils.GetStringPropOrEmpty(props, "valueProposition"),
		LastFundingRound:  utils.GetStringPropOrEmpty(props, "lastFundingRound"),
		LastFundingAmount: utils.GetStringPropOrEmpty(props, "lastFundingAmount"),
		Note:              utils.GetStringPropOrEmpty(props, "note"),
		IsPublic:          utils.GetBoolPropOrFalse(props, "isPublic"),
		IsCustomer:        utils.GetBoolPropOrFalse(props, "isCustomer"),
		Hide:              utils.GetBoolPropOrFalse(props, "hide"),
		Employees:         utils.GetInt64PropOrZero(props, "employees"),
		Market:            utils.GetStringPropOrEmpty(props, "market"),
		CreatedAt:         utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:         utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:            entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:     entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:         utils.GetStringPropOrEmpty(props, "appSource"),
		LastTouchpointAt:  utils.GetTimePropOrNil(props, "lastTouchpointAt"),
		LastTouchpointId:  utils.GetStringPropOrNil(props, "lastTouchpointId"),
		RenewalLikelihood: entity.RenewalLikelihood{
			RenewalLikelihood:         utils.GetStringPropOrEmpty(props, "renewalLikelihood"),
			PreviousRenewalLikelihood: utils.GetStringPropOrEmpty(props, "renewalLikelihoodPrevious"),
			Comment:                   utils.GetStringPropOrNil(props, "renewalLikelihoodComment"),
			UpdatedBy:                 utils.GetStringPropOrNil(props, "renewalLikelihoodUpdatedBy"),
			UpdatedAt:                 utils.GetTimePropOrNil(props, "renewalLikelihoodUpdatedAt"),
		},
		RenewalForecast: entity.RenewalForecast{
			Amount:          utils.GetFloatPropOrNil(props, "renewalForecastAmount"),
			PotentialAmount: utils.GetFloatPropOrNil(props, "renewalForecastPotentialAmount"),
			Comment:         utils.GetStringPropOrNil(props, "renewalForecastComment"),
			UpdatedById:     utils.GetStringPropOrNil(props, "renewalForecastUpdatedBy"),
			UpdatedAt:       utils.GetTimePropOrNil(props, "renewalForecastUpdatedAt"),
		},
		BillingDetails: entity.BillingDetails{
			Amount:            utils.GetFloatPropOrNil(props, "billingDetailsAmount"),
			Frequency:         utils.GetStringPropOrEmpty(props, "billingDetailsFrequency"),
			RenewalCycle:      utils.GetStringPropOrEmpty(props, "billingDetailsRenewalCycle"),
			RenewalCycleStart: utils.GetTimePropOrNil(props, "billingDetailsRenewalCycleStart"),
			RenewalCycleNext:  utils.GetTimePropOrNil(props, "billingDetailsRenewalCycleNext"),
		},
	}
	return &output
}

func (d domains) isPersonalEmailProvider(domain string) bool {
	for _, v := range d.personalEmailProviders {
		if strings.ToLower(domain) == strings.ToLower(v.ProviderDomain) {
			return true
		}
	}
	return false
}

func (d domains) isWhitelistedDomain(domain string) bool {
	if domain == "" {
		return false
	}
	for _, v := range d.whitelistDomains {
		if v.Domain != "*" && strings.ToLower(domain) == strings.ToLower(v.Domain) && v.Allowed {
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
