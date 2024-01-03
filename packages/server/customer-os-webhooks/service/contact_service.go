package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
		contactId = utils.NewUUIDIfEmpty(contactId)
		contactInput.Id = contactId
		span.LogFields(log.String("contactId", contactId))

		// Create or update contact
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = s.grpcClients.ContactClient.UpsertContact(ctx, &contactpb.UpsertContactGrpcRequest{
			Tenant:          tenant,
			Id:              contactId,
			Name:            contactInput.Name,
			FirstName:       contactInput.FirstName,
			LastName:        contactInput.LastName,
			Description:     contactInput.Description,
			Timezone:        contactInput.Timezone,
			ProfilePhotoUrl: contactInput.ProfilePhotoUrl,
			CreatedAt:       utils.ConvertTimeToTimestampPtr(contactInput.CreatedAt),
			UpdatedAt:       utils.ConvertTimeToTimestampPtr(contactInput.UpdatedAt),
			SourceFields: &commonpb.SourceFields{
				Source:    contactInput.ExternalSystem,
				AppSource: appSource,
			},
			ExternalSystemFields: &commonpb.ExternalSystemFields{
				ExternalSystemId: contactInput.ExternalSystem,
				ExternalId:       contactInput.ExternalId,
				ExternalUrl:      contactInput.ExternalUrl,
				ExternalIdSecond: contactInput.ExternalIdSecond,
				ExternalSource:   contactInput.ExternalSourceEntity,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
		})
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err, log.String("grpcMethod", "UpsertContact"))
			reason = fmt.Sprintf("failed sending event to upsert contact  with external reference %s for tenant %s :%s", contactInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
		// Wait for contact to be created in neo4j
		if !failedSync && !matchingContactExists {
			for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
				contact, findErr := s.repositories.ContactRepository.GetById(ctx, tenant, contactId)
				if contact != nil && findErr == nil {
					break
				}
				time.Sleep(time.Duration(i*constants.TimeoutIntervalMs) * time.Millisecond)
			}
		}
	}
	if !failedSync && contactInput.HasPrimaryEmail() {
		// Create or update email
		emailId, err := s.services.EmailService.CreateEmail(ctx, contactInput.Email, contactInput.ExternalSystem, appSource)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("Failed to create email address for contact %s: %s", contactId, err.Error())
			s.log.Error(reason)
		}
		// Link email to contact
		if !failedSync {
			_, err = s.grpcClients.ContactClient.LinkEmailToContact(ctx, &contactpb.LinkEmailToContactGrpcRequest{
				Tenant:    common.GetTenantFromContext(ctx),
				ContactId: contactId,
				EmailId:   emailId,
				Primary:   true,
				AppSource: appSource,
			})
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err, log.String("grpcMethod", "LinkEmailToContact"))
				reason = fmt.Sprintf("Failed to link email address %s with contact %s: %s", contactInput.Email, contactId, err.Error())
				s.log.Error(reason)
			}
		}
	}
	if !failedSync && contactInput.HasAdditionalEmails() {
		for _, email := range contactInput.AdditionalEmails {
			// Create or update email
			emailId, err := s.services.EmailService.CreateEmail(ctx, email, contactInput.ExternalSystem, appSource)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("Failed to create email address for contact %s: %s", contactId, err.Error())
				s.log.Error(reason)
			}
			// Link email to contact
			if !failedSync {
				_, err = s.grpcClients.ContactClient.LinkEmailToContact(ctx, &contactpb.LinkEmailToContactGrpcRequest{
					Tenant:    common.GetTenantFromContext(ctx),
					ContactId: contactId,
					EmailId:   emailId,
					Primary:   false,
					AppSource: appSource,
				})
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err, log.String("grpcMethod", "LinkEmailToContact"))
					reason = fmt.Sprintf("Failed to link email address %s with contact %s: %s", email, contactId, err.Error())
					s.log.Error(reason)
				}
			}
		}
	}

	if !failedSync {
		for orgId, referencedOrganization := range identifiedOrganizations {
			// Link contact to organization
			_, err = s.grpcClients.ContactClient.LinkWithOrganization(ctx, &contactpb.LinkWithOrganizationGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				ContactId:      contactId,
				OrganizationId: orgId,
				JobTitle:       referencedOrganization.JobTitle,
				SourceFields: &commonpb.SourceFields{
					Source:    contactInput.ExternalSystem,
					AppSource: appSource,
				},
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
					_, err = s.grpcClients.ContactClient.LinkPhoneNumberToContact(ctx, &contactpb.LinkPhoneNumberToContactGrpcRequest{
						Tenant:        common.GetTenantFromContext(ctx),
						ContactId:     contactId,
						PhoneNumberId: phoneNumberId,
						Primary:       phoneNumberDtls.Primary,
						Label:         phoneNumberDtls.Label,
						AppSource:     appSource,
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
		if contactInput.HasLocation() {
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
				_, err = s.grpcClients.ContactClient.LinkLocationToContact(ctx, &contactpb.LinkLocationToContactGrpcRequest{
					Tenant:     common.GetTenantFromContext(ctx),
					ContactId:  contactId,
					LocationId: locationId,
					AppSource:  appSource,
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

func (s *contactService) mapDbNodeToContactEntity(dbNode dbtype.Node) *entity.ContactEntity {
	props := utils.GetPropsFromNode(dbNode)
	output := entity.ContactEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		Name:            utils.GetStringPropOrEmpty(props, "name"),
		FirstName:       utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:        utils.GetStringPropOrEmpty(props, "lastName"),
		Description:     utils.GetStringPropOrEmpty(props, "description"),
		Timezone:        utils.GetStringPropOrEmpty(props, "timezone"),
		ProfilePhotoUrl: utils.GetStringPropOrEmpty(props, "profilePhotoUrl"),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:          neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &output
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
