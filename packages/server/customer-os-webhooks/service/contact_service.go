package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	pkgerrors "github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

type ContactService interface {
	SyncContacts(ctx context.Context, contacts []model.ContactData) (SyncResult, error)
	GetIdForReferencedContact(ctx context.Context, tenant, externalSystem string, contact model.ReferencedContact) (string, error)
}

type contactService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
	maxWorkers   int
}

func NewContactService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) ContactService {
	return &contactService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
		maxWorkers:   services.cfg.ConcurrencyConfig.ContactSyncConcurrency,
	}
}

func (s *contactService) SyncContacts(ctx context.Context, contacts []model.ContactData) (SyncResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.SyncContacts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("num of contacts", len(contacts)))

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	// pre-validate contact input before syncing
	for _, contact := range contacts {
		if contact.ExternalSystem == "" {
			tracing.TraceErr(span, errors.ErrMissingExternalSystem)
			return SyncResult{}, errors.ErrMissingExternalSystem
		}
		if !neo4jentity.IsValidDataSource(strings.ToLower(contact.ExternalSystem)) {
			tracing.TraceErr(span, errors.ErrExternalSystemNotAccepted, log.String("externalSystem", contact.ExternalSystem))
			return SyncResult{}, errors.ErrExternalSystemNotAccepted
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

	// Sync all contacts concurrently
	for _, contactData := range contacts {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return SyncResult{}, ctx.Err()
		default:
		}

		// Acquire a worker slot
		workerLimit <- struct{}{}
		wg.Add(1)

		go func(contactData model.ContactData) {
			defer wg.Done()
			defer func() {
				// Release the worker slot when done
				<-workerLimit
			}()

			result := s.syncContact(ctx, syncMutex, contactData, syncDate)
			statusesMutex.Lock()
			statuses = append(statuses, result)
			statusesMutex.Unlock()
		}(contactData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), contacts[0].ExternalSystem,
		contacts[0].AppSource, "contact", syncDate, statuses)

	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}

func (s *contactService) syncContact(ctx context.Context, syncMutex *sync.Mutex, contactInput model.ContactData, syncDate time.Time) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.syncContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagExternalSystem, contactInput.ExternalSystem)
	span.LogFields(log.Object("syncDate", syncDate))
	tracing.LogObjectAsJson(span, "contactInput", contactInput)

	tenant := common.GetTenantFromContext(ctx)
	var appSource = utils.StringFirstNonEmpty(contactInput.AppSource, constants.AppSourceCustomerOsWebhooks)
	var failedSync = false
	var reason = ""

	contactInput.Normalize()

	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, contactInput.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", contactInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}

	// Check if contact sync should be skipped
	if contactInput.Skip {
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(contactInput.SkipReason)
	}

	// Filter out un-existing organizations
	var identifiedOrganizations = make(map[string]model.ReferencedOrganization)
	for _, org := range contactInput.Organizations {
		orgId, _ := s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, contactInput.ExternalSystem, org)
		if orgId != "" {
			identifiedOrganizations[orgId] = org
		}
	}

	if contactInput.OrganizationRequired && len(identifiedOrganizations) == 0 {
		reason = fmt.Sprintf("organization(s) not found for contact %s for tenant %s", contactInput.ExternalId, tenant)
		s.log.Warn(reason)
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(reason)
	}

	if contactInput.Name == "" {
		contactInput.Name = strings.TrimSpace(fmt.Sprintf("%s %s", contactInput.FirstName, contactInput.LastName))
	}

	// Lock contact creation
	syncMutex.Lock()
	defer syncMutex.Unlock()
	// Check if contact already exists
	contactId, err := s.repositories.ContactRepository.GetMatchedContactId(ctx, tenant, contactInput.ExternalSystem, contactInput.ExternalId, contactInput.EmailsForUnicity())
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched contact with external reference %s for tenant %s :%s", contactInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}
	if !failedSync {
		matchingContactExists := contactId != ""
		span.LogFields(log.Bool("found matching contact", matchingContactExists))

		// Create new contact id if not found
		createContact := true
		if contactId != "" {
			createContact = false
		}
		if createContact {
			contactId, err = s.services.CommonServices.ContactService.SaveContact(ctx, nil,
				neo4jrepo.ContactFields{
					Name:            contactInput.Name,
					FirstName:       contactInput.FirstName,
					LastName:        contactInput.LastName,
					Description:     contactInput.Description,
					Timezone:        contactInput.Timezone,
					ProfilePhotoUrl: contactInput.ProfilePhotoUrl,
					CreatedAt:       utils.TimeOrNowFromPtr(contactInput.CreatedAt),
					SourceFields: neo4jmodel.SourceFields{
						Source:    contactInput.ExternalSystem,
						AppSource: appSource,
					},
				},
				"",
				neo4jmodel.ExternalSystem{
					ExternalSystemId: contactInput.ExternalSystem,
					ExternalId:       contactInput.ExternalId,
					ExternalUrl:      contactInput.ExternalUrl,
					ExternalIdSecond: contactInput.ExternalIdSecond,
					ExternalSource:   contactInput.ExternalSourceEntity,
					SyncDate:         &syncDate,
				})
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed creating contact id for external reference %s for tenant %s :%s", contactInput.ExternalId, tenant, err.Error())
				s.log.Error(reason)
			}
			contactInput.Id = contactId
		} else {
			// update contact
			// fetch current contact data
			contactEntity, err := s.services.CommonServices.ContactService.GetContactById(ctx, contactId)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed fetching contact with id %s for tenant %s :%s", contactId, tenant, err.Error())
				s.log.Error(reason)
			}
			contactFields := neo4jrepo.ContactFields{
				SourceFields: neo4jmodel.SourceFields{
					Source:    contactInput.ExternalSystem,
					AppSource: appSource,
				},
			}
			if contactEntity.Name == "" && contactInput.Name != "" {
				contactFields.Name = contactInput.Name
				contactFields.UpdateName = true
			}
			if contactEntity.FirstName == "" && contactInput.FirstName != "" {
				contactFields.FirstName = contactInput.FirstName
				contactFields.UpdateFirstName = true
			}
			if contactEntity.LastName == "" && contactInput.LastName != "" {
				contactFields.LastName = contactInput.LastName
				contactFields.UpdateLastName = true
			}
			if contactEntity.Description == "" && contactInput.Description != "" {
				contactFields.Description = contactInput.Description
				contactFields.UpdateDescription = true
			}
			if contactEntity.Timezone == "" && contactInput.Timezone != "" {
				contactFields.Timezone = contactInput.Timezone
				contactFields.UpdateTimezone = true
			}
			if contactEntity.ProfilePhotoUrl == "" && contactInput.ProfilePhotoUrl != "" {
				contactFields.ProfilePhotoUrl = contactInput.ProfilePhotoUrl
				contactFields.UpdateProfilePhotoUrl = true
			}
			_, err = s.services.CommonServices.ContactService.SaveContact(ctx, &contactId,
				contactFields,
				"",
				neo4jmodel.ExternalSystem{
					ExternalSystemId: contactInput.ExternalSystem,
					ExternalId:       contactInput.ExternalId,
					ExternalUrl:      contactInput.ExternalUrl,
					ExternalIdSecond: contactInput.ExternalIdSecond,
					ExternalSource:   contactInput.ExternalSourceEntity,
					SyncDate:         &syncDate,
				})
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed updating contact with external reference %s for tenant %s :%s", contactInput.ExternalId, tenant, err.Error())
				s.log.Error(reason)
			}
		}
		span.LogFields(log.String("contactId", contactInput.Id))
	}
	if !failedSync && contactInput.HasPrimaryEmail() {
		_, err = s.services.CommonServices.EmailService.Merge(ctx, tenant,
			commonservice.EmailFields{
				Email:     contactInput.Email,
				AppSource: appSource,
				Source:    neo4jentity.DecodeDataSource(contactInput.ExternalSystem),
				Primary:   true,
			},
			&commonservice.LinkWith{
				Type: commonmodel.CONTACT,
				Id:   contactId,
			})
		if err != nil {
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("Failed to create and link email address %s with contact %s: %s", contactInput.Email, contactId, err.Error())
			failedSync = true
		}
	}
	if !failedSync && contactInput.HasAdditionalEmails() {
		for _, email := range contactInput.AdditionalEmails {
			_, err = s.services.CommonServices.EmailService.Merge(ctx, tenant,
				commonservice.EmailFields{
					Email:     email,
					AppSource: appSource,
					Source:    neo4jentity.DecodeDataSource(contactInput.ExternalSystem),
					Primary:   false,
				},
				&commonservice.LinkWith{
					Type: commonmodel.CONTACT,
					Id:   contactId,
				})
			if err != nil {
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("Failed to create and link email address %s with contact %s: %s", contactInput.Email, contactId, err.Error())
				failedSync = true
			}
		}
	}
	if !failedSync && contactInput.HasSocials() {
		for _, social := range contactInput.Socials {
			_, err = s.services.CommonServices.SocialService.MergeSocialWithEntity(ctx, tenant, contactId, commonmodel.CONTACT,
				neo4jentity.SocialEntity{
					Url:       social.URL,
					Source:    neo4jentity.DecodeDataSource(contactInput.ExternalSystem),
					AppSource: appSource,
				})
			if err != nil {
				tracing.TraceErr(span, pkgerrors.Wrap(err, "failed to merge social with contact"))
				reason = fmt.Sprintf("Failed to merge social %s with contact %s: %s", social.URL, contactId, err.Error())
				s.log.Error(reason)
			}
		}
	}

	if !failedSync {
		for orgId, referencedOrganization := range identifiedOrganizations {
			// Link contact to organization
			_, err = CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
				return s.grpcClients.ContactClient.LinkWithOrganization(ctx, &contactpb.LinkWithOrganizationGrpcRequest{
					Tenant:         common.GetTenantFromContext(ctx),
					ContactId:      contactId,
					OrganizationId: orgId,
					JobTitle:       referencedOrganization.JobTitle,
					Description:    referencedOrganization.JobDescription,
					SourceFields: &commonpb.SourceFields{
						Source:    contactInput.ExternalSystem,
						AppSource: appSource,
					},
				})
			})
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err, log.String("grpcMethod", "LinkWithOrganization"))
				reason = fmt.Sprintf("Failed to link contact %s with organization %s: %s", contactId, orgId, err.Error())
				s.log.Error(reason)
			}
		}
	}

	if !failedSync {
		if contactInput.HasPhoneNumbers() {
			for _, phoneNumberDtls := range contactInput.PhoneNumbers {
				// Create or update phone number
				phoneNumberId, err := s.services.PhoneNumberService.CreatePhoneNumber(ctx, phoneNumberDtls.Number, contactInput.ExternalSystem, contactInput.AppSource)
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("Failed to create phone number %s for contact %s: %s", phoneNumberDtls.Number, contactId, err.Error())
					s.log.Error(reason)
				}
				// Link phone number to contact
				if phoneNumberId != "" {
					_, err = CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
						return s.grpcClients.ContactClient.LinkPhoneNumberToContact(ctx, &contactpb.LinkPhoneNumberToContactGrpcRequest{
							Tenant:        common.GetTenantFromContext(ctx),
							ContactId:     contactId,
							PhoneNumberId: phoneNumberId,
							Primary:       phoneNumberDtls.Primary,
							Label:         phoneNumberDtls.Label,
							AppSource:     appSource,
						})
					})
					if err != nil {
						failedSync = true
						tracing.TraceErr(span, err, log.String("grpcMethod", "LinkPhoneNumberToContact"))
						reason = fmt.Sprintf("Failed to link phone number %s with contact %s: %s", phoneNumberDtls.Number, contactId, err.Error())
						s.log.Error(reason)
					}
				}
			}
		}

		syncLocation := false // skip location sync for now
		if contactInput.HasLocation() && syncLocation {
			// Create or update location
			locationId, err := s.repositories.LocationRepository.GetMatchedLocationIdForContactBySource(ctx, contactId, contactInput.ExternalSystem)
			if err != nil {
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("Failed to get matched location for contact %s: %s", contactId, err.Error())
				failedSync = true
				s.log.Error(reason)
			}

			if !failedSync {
				locationId, err = s.services.LocationService.CreateLocation(ctx, locationId, contactInput.ExternalSystem, contactInput.AppSource,
					contactInput.LocationName, contactInput.Country, contactInput.Region, contactInput.Locality, contactInput.Street, contactInput.Address, "", contactInput.Zip, contactInput.PostalCode)
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("Failed to create location for contact %s: %s", contactId, err.Error())
					s.log.Error(reason)
				}
			}
			// Link location to contact
			if locationId != "" {
				_, err = CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
					return s.grpcClients.ContactClient.LinkLocationToContact(ctx, &contactpb.LinkLocationToContactGrpcRequest{
						Tenant:     common.GetTenantFromContext(ctx),
						ContactId:  contactId,
						LocationId: locationId,
						AppSource:  appSource,
					})
				})
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err, log.String("grpcMethod", "LinkLocationToContact"))
					reason = fmt.Sprintf("Failed to link location %s with contact %s: %s", locationId, contactId, err.Error())
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

func (s *contactService) GetIdForReferencedContact(ctx context.Context, tenant, externalSystemId string, contact model.ReferencedContact) (string, error) {
	if !contact.Available() {
		return "", nil
	}

	if contact.ReferencedById() {
		return s.repositories.ContactRepository.GetContactIdById(ctx, tenant, contact.Id)
	} else if contact.ReferencedByExternalId() {
		return s.repositories.ContactRepository.GetContactIdByExternalId(ctx, tenant, contact.ExternalId, externalSystemId)
	}
	return "", nil
}
