package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
)

type ContactService interface {
	Create(ctx context.Context, contact *ContactCreateData) (string, error)
	Update(ctx context.Context, input model.ContactUpdateInput) (string, error)
	GetById(ctx context.Context, id string) (*neo4jentity.ContactEntity, error)
	GetFirstContactByEmail(ctx context.Context, email string) (*neo4jentity.ContactEntity, error)
	GetFirstContactByPhoneNumber(ctx context.Context, phoneNumber string) (*neo4jentity.ContactEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	PermanentDelete(ctx context.Context, id string) (bool, error)
	Archive(ctx context.Context, contactId string) (bool, error)
	RestoreFromArchive(ctx context.Context, contactId string) (bool, error)
	GetContactsForJobRoles(ctx context.Context, jobRoleIds []string) (*neo4jentity.ContactEntities, error)
	GetContactsForOrganization(ctx context.Context, organizationId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	Merge(ctx context.Context, primaryContactId, mergedContactId string) error
	GetContactsForEmails(ctx context.Context, emailIds []string) (*neo4jentity.ContactEntities, error)
	GetContactsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*neo4jentity.ContactEntities, error)
	AddOrganization(ctx context.Context, contactId, organizationId, source, appSource string) (*neo4jentity.ContactEntity, error)
	RemoveOrganization(ctx context.Context, contactId, organizationId string) (*neo4jentity.ContactEntity, error)
	RemoveLocation(ctx context.Context, contactId string, locationId string) error
	CustomerContactCreate(ctx context.Context, entity *CustomerContactCreateData) (*model.CustomerContact, error)
	GetContactCountByOrganizations(ctx context.Context, ids []string) (map[string]int64, error)
}

type ContactCreateData struct {
	ContactEntity     *neo4jentity.ContactEntity
	EmailEntity       *neo4jentity.EmailEntity
	PhoneNumberEntity *neo4jentity.PhoneNumberEntity
	ExternalReference *neo4jentity.ExternalSystemEntity
	Source            neo4jentity.DataSource
	SocialUrl         string
	AppSource         string
}

type CustomerContactCreateData struct {
	ContactEntity *neo4jentity.ContactEntity
	EmailEntity   *neo4jentity.EmailEntity
}

type contactService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewContactService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) ContactService {
	return &contactService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *contactService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *contactService) Create(ctx context.Context, contactDetails *ContactCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contactDetails", contactDetails))

	if contactDetails.ContactEntity == nil {
		err := fmt.Errorf("contact entity is nil")
		tracing.TraceErr(span, err)
		return "", err
	}

	contactId, err := s.createContactWithEvents(ctx, contactDetails)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	if contactDetails.EmailEntity != nil {
		s.linkEmailByEvents(ctx, contactId, utils.StringFirstNonEmpty(contactDetails.EmailEntity.AppSource, contactDetails.AppSource), *contactDetails.EmailEntity)
	}

	if contactDetails.PhoneNumberEntity != nil {
		s.linkPhoneNumberByEvents(ctx, contactId, utils.StringFirstNonEmpty(contactDetails.PhoneNumberEntity.AppSource, contactDetails.AppSource), *contactDetails.PhoneNumberEntity)
	}

	span.LogFields(log.String("output - createdContactId", contactId))
	return contactId, nil
}

func (s *contactService) createContactWithEvents(ctx context.Context, contactDetails *ContactCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	upsertContactRequest := contactpb.UpsertContactGrpcRequest{
		Tenant: common.GetTenantFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    string(contactDetails.Source),
			AppSource: utils.StringFirstNonEmpty(contactDetails.AppSource, constants.AppSourceCustomerOsApi),
		},
		LoggedInUserId:  common.GetUserIdFromContext(ctx),
		FirstName:       contactDetails.ContactEntity.FirstName,
		LastName:        contactDetails.ContactEntity.LastName,
		Prefix:          contactDetails.ContactEntity.Prefix,
		Description:     contactDetails.ContactEntity.Description,
		ProfilePhotoUrl: contactDetails.ContactEntity.ProfilePhotoUrl,
		Username:        contactDetails.ContactEntity.Username,
		SocialUrl:       contactDetails.SocialUrl,
		Name:            contactDetails.ContactEntity.Name,
		Timezone:        contactDetails.ContactEntity.Timezone,
	}
	upsertContactRequest.CreatedAt = timestamppb.New(contactDetails.ContactEntity.CreatedAt)
	if contactDetails.ExternalReference != nil && contactDetails.ExternalReference.ExternalSystemId != "" {
		upsertContactRequest.ExternalSystemFields = &commonpb.ExternalSystemFields{
			ExternalSystemId: string(contactDetails.ExternalReference.ExternalSystemId),
			ExternalId:       contactDetails.ExternalReference.Relationship.ExternalId,
			ExternalUrl:      utils.IfNotNilString(contactDetails.ExternalReference.Relationship.ExternalUrl),
			ExternalSource:   utils.IfNotNilString(contactDetails.ExternalReference.Relationship.ExternalSource),
		}
		if contactDetails.ExternalReference.Relationship.SyncDate != nil {
			upsertContactRequest.ExternalSystemFields.SyncDate = timestamppb.New(*contactDetails.ExternalReference.Relationship.SyncDate)
		}
	}
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
		return s.grpcClients.ContactClient.UpsertContact(ctx, &upsertContactRequest)
	})

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, model2.NodeLabelContact, span)

	return response.Id, err
}

func (s *contactService) linkEmailByEvents(ctx context.Context, contactId, appSource string, emailEntity neo4jentity.EmailEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.linkEmailByEvents")
	defer span.Finish()

	emailId, err := s.services.EmailService.CreateEmailAddressViaEvents(ctx, utils.StringFirstNonEmpty(emailEntity.RawEmail, emailEntity.Email), appSource)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to create email address for contact %s: %s", contactId, err.Error())
	}
	if emailId != "" {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
			return s.grpcClients.ContactClient.LinkEmailToContact(ctx, &contactpb.LinkEmailToContactGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				ContactId:      contactId,
				EmailId:        emailId,
				Primary:        emailEntity.Primary,
				Label:          emailEntity.Label,
				AppSource:      appSource,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to link email address %s with contact %s: %s", emailId, contactId, err.Error())
		}
	}
}

func (s *contactService) linkPhoneNumberByEvents(ctx context.Context, contactId, appSource string, phoneNumberEntity neo4jentity.PhoneNumberEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.linkPhoneNumberByEvents")
	defer span.Finish()

	phoneNumberId, err := s.services.PhoneNumberService.CreatePhoneNumberViaEvents(ctx, utils.StringFirstNonEmpty(phoneNumberEntity.RawPhoneNumber, phoneNumberEntity.E164), appSource)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to create phone number for contact %s: %s", contactId, err.Error())
	}
	if phoneNumberId != "" {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
			return s.grpcClients.ContactClient.LinkPhoneNumberToContact(ctx, &contactpb.LinkPhoneNumberToContactGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				ContactId:      contactId,
				PhoneNumberId:  phoneNumberId,
				Primary:        phoneNumberEntity.Primary,
				Label:          phoneNumberEntity.Label,
				AppSource:      appSource,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to link phone number %s with contact %s: %s", phoneNumberId, contactId, err.Error())
		}
	}
}

func (s *contactService) Update(ctx context.Context, input model.ContactUpdateInput) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	if input.ID == "" {
		err := fmt.Errorf("(ContactService.Update) contact id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return "", err
	}

	currentContactEntity, err := s.GetById(ctx, input.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return "", err
	}

	upsertContactRequest := contactpb.UpsertContactGrpcRequest{
		Tenant: common.GetTenantFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    string(neo4jentity.DataSourceOpenline),
			AppSource: constants.AppSourceCustomerOsApi,
		},
		LoggedInUserId:  common.GetUserIdFromContext(ctx),
		Id:              input.ID,
		Prefix:          utils.IfNotNilString(input.Prefix),
		Name:            utils.IfNotNilString(input.Name),
		FirstName:       utils.IfNotNilString(input.FirstName),
		LastName:        utils.IfNotNilString(input.LastName),
		Description:     utils.IfNotNilString(input.Description),
		Timezone:        utils.IfNotNilString(input.Timezone),
		ProfilePhotoUrl: utils.StringFirstNonEmpty(utils.IfNotNilString(input.ProfilePhotoURL), currentContactEntity.ProfilePhotoUrl),
		Username:        utils.StringFirstNonEmpty(utils.IfNotNilString(input.Username), currentContactEntity.Username),
	}
	if input.Patch != nil && *input.Patch {
		var fieldsMask []contactpb.ContactFieldMask
		if input.FirstName != nil {
			fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_FIRST_NAME)
		}
		if input.LastName != nil {
			fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_LAST_NAME)
		}
		if input.Prefix != nil {
			fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_PREFIX)
		}
		if input.Name != nil {
			fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_NAME)
		}
		if input.Description != nil {
			fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_DESCRIPTION)
		}
		if input.Timezone != nil {
			fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_TIMEZONE)
		}
		if input.ProfilePhotoURL != nil {
			fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_PROFILE_PHOTO_URL)
		}
		if input.Username != nil {
			fieldsMask = append(fieldsMask, contactpb.ContactFieldMask_CONTACT_FIELD_USERNAME)
		}
		if len(fieldsMask) == 0 {
			span.LogFields(log.String("result", "No fields to update"))
			return input.ID, nil
		}
		upsertContactRequest.FieldsMask = fieldsMask
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
		return s.grpcClients.ContactClient.UpsertContact(ctx, &upsertContactRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error("Error from events processing: %s", err.Error())
		return "", err
	}
	return response.Id, nil
}

func (s *contactService) PermanentDelete(ctx context.Context, contactId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	err := s.repositories.ContactRepository.Delete(ctx, session, common.GetContext(ctx).Tenant, contactId)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *contactService) Archive(ctx context.Context, contactId string) (bool, error) {
	err := s.repositories.ContactRepository.Archive(ctx, common.GetTenantFromContext(ctx), contactId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *contactService) RestoreFromArchive(ctx context.Context, contactId string) (bool, error) {
	err := s.repositories.ContactRepository.RestoreFromArchive(ctx, common.GetTenantFromContext(ctx), contactId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *contactService) GetById(ctx context.Context, contactId string) (*neo4jentity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))

	if contactDbNode, err := s.repositories.ContactRepository.GetById(ctx, common.GetContext(ctx).Tenant, contactId); err != nil {
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Contact with id {%s} not found", contactId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToContactEntity(contactDbNode), nil
	}
}

func (s *contactService) GetFirstContactByEmail(ctx context.Context, email string) (*neo4jentity.ContactEntity, error) {
	dbNodes, err := s.repositories.Neo4jRepositories.ContactReadRepository.GetContactsWithEmail(ctx, common.GetContext(ctx).Tenant, email)
	if err != nil || len(dbNodes) == 0 {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToContactEntity(dbNodes[0]), nil
}

func (s *contactService) GetFirstContactByPhoneNumber(ctx context.Context, phoneNumber string) (*neo4jentity.ContactEntity, error) {
	dbNodes, err := s.repositories.ContactRepository.GetContactsForPhoneNumber(ctx, common.GetContext(ctx).Tenant, phoneNumber)
	if err != nil || len(dbNodes) == 0 {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToContactEntity(dbNodes[0]), nil
}

func (s *contactService) FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(neo4jentity.ContactEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(neo4jentity.ContactEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repositories.ContactRepository.GetPaginatedContacts(
		ctx, session,
		common.GetContext(ctx).Tenant,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	contacts := make(neo4jentity.ContactEntities, 0, len(dbNodesWithTotalCount.Nodes))

	for _, v := range dbNodesWithTotalCount.Nodes {
		contacts = append(contacts, *neo4jmapper.MapDbNodeToContactEntity(v))
	}
	paginatedResult.SetRows(&contacts)
	return &paginatedResult, nil
}

func (s *contactService) GetContactsForJobRoles(ctx context.Context, jobRoleIds []string) (*neo4jentity.ContactEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactsForJobRoles")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("jobRoleIds", jobRoleIds))

	contacts, err := s.repositories.ContactRepository.GetAllForJobRoles(ctx, common.GetTenantFromContext(ctx), jobRoleIds)
	if err != nil {
		return nil, err
	}
	contactEntities := make(neo4jentity.ContactEntities, 0, len(contacts))
	for _, v := range contacts {
		contactEntity := neo4jmapper.MapDbNodeToContactEntity(v.Node)
		contactEntity.DataloaderKey = v.LinkedNodeId
		contactEntities = append(contactEntities, *contactEntity)
	}
	return &contactEntities, nil
}

func (s *contactService) GetContactsForOrganization(ctx context.Context, organizationId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactsForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.Int("page", page), log.Int("limit", limit))
	if filter != nil {
		span.LogFields(log.Object("filter", filter))
	}
	if sortBy != nil {
		span.LogFields(log.Object("sortBy", sortBy))
	}

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(neo4jentity.ContactEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(neo4jentity.ContactEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repositories.ContactRepository.GetPaginatedContactsForOrganization(
		ctx, session,
		common.GetTenantFromContext(ctx),
		organizationId,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	contacts := make(neo4jentity.ContactEntities, 0, len(dbNodesWithTotalCount.Nodes))
	for _, v := range dbNodesWithTotalCount.Nodes {
		contacts = append(contacts, *neo4jmapper.MapDbNodeToContactEntity(v))
	}
	paginatedResult.SetRows(&contacts)
	return &paginatedResult, nil
}

func (s *contactService) Merge(ctx context.Context, primaryContactId, mergedContactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.Merge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("primaryContactId", primaryContactId), log.String("mergedContactId", mergedContactId))

	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	_, err := s.GetById(ctx, primaryContactId)
	if err != nil {
		s.log.Errorf("(%s) Primary contact with id {%s} not found: {%v}", utils.GetFunctionName(), primaryContactId, err.Error())
		return err
	}
	_, err = s.GetById(ctx, mergedContactId)
	if err != nil {
		s.log.Errorf("(%s) Contact to merge with id {%s} not found: {%v}", utils.GetFunctionName(), mergedContactId, err.Error())
		return err
	}

	tenant := common.GetContext(ctx).Tenant
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		err = s.repositories.ContactRepository.MergeContactPropertiesInTx(ctx, tx, tenant, primaryContactId, mergedContactId, neo4jentity.DataSourceOpenline)
		if err != nil {
			return nil, err
		}

		err = s.repositories.ContactRepository.MergeContactRelationsInTx(ctx, tx, tenant, primaryContactId, mergedContactId)
		if err != nil {
			return nil, err
		}

		err = s.repositories.ContactRepository.UpdateMergedContactLabelsInTx(ctx, tx, tenant, mergedContactId)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, primaryContactId)
	}

	return err
}

func (s *contactService) AddOrganization(ctx context.Context, contactId, organizationId, source, appSource string) (*neo4jentity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.AddOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("organizationId", organizationId))

	contactNodePtr, err := s.repositories.ContactRepository.AddOrganization(ctx, common.GetTenantFromContext(ctx), contactId, organizationId, source, appSource)
	if err != nil {
		return nil, err
	}
	s.services.OrganizationService.UpdateLastTouchpoint(ctx, organizationId)
	return neo4jmapper.MapDbNodeToContactEntity(contactNodePtr), nil
}

func (s *contactService) RemoveOrganization(ctx context.Context, contactId, organizationId string) (*neo4jentity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.RemoveOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("organizationId", organizationId))

	contactNodePtr, err := s.repositories.ContactRepository.RemoveOrganization(ctx, common.GetTenantFromContext(ctx), contactId, organizationId)
	if err != nil {
		return nil, err
	}
	s.services.OrganizationService.UpdateLastTouchpoint(ctx, organizationId)
	return neo4jmapper.MapDbNodeToContactEntity(contactNodePtr), nil
}

func (s *contactService) GetContactsForEmails(ctx context.Context, emailIds []string) (*neo4jentity.ContactEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactsForEmails")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("emailIds", emailIds))

	contacts, err := s.repositories.ContactRepository.GetAllForEmails(ctx, common.GetTenantFromContext(ctx), emailIds)
	if err != nil {
		return nil, err
	}
	contactEntities := make(neo4jentity.ContactEntities, 0, len(contacts))
	for _, v := range contacts {
		contactEntity := neo4jmapper.MapDbNodeToContactEntity(v.Node)
		contactEntity.DataloaderKey = v.LinkedNodeId
		contactEntities = append(contactEntities, *contactEntity)
	}
	return &contactEntities, nil
}

func (s *contactService) GetContactsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*neo4jentity.ContactEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactsForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("phoneNumberIds", phoneNumberIds))

	contacts, err := s.repositories.ContactRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), phoneNumberIds)
	if err != nil {
		return nil, err
	}
	contactEntities := make(neo4jentity.ContactEntities, 0, len(contacts))
	for _, v := range contacts {
		contactEntity := neo4jmapper.MapDbNodeToContactEntity(v.Node)
		contactEntity.DataloaderKey = v.LinkedNodeId
		contactEntities = append(contactEntities, *contactEntity)
	}
	return &contactEntities, nil
}

func (s *contactService) CustomerContactCreate(ctx context.Context, data *CustomerContactCreateData) (*model.CustomerContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.CustomerContactCreate")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	result := &model.CustomerContact{}

	contactCreateRequest := &contactpb.UpsertContactGrpcRequest{
		Tenant:      common.GetTenantFromContext(ctx),
		FirstName:   data.ContactEntity.FirstName,
		LastName:    data.ContactEntity.LastName,
		Prefix:      data.ContactEntity.Prefix,
		Description: data.ContactEntity.Description,
		SourceFields: &commonpb.SourceFields{
			Source:    string(data.ContactEntity.Source),
			AppSource: data.ContactEntity.AppSource,
		},
		LoggedInUserId: common.GetUserIdFromContext(ctx),
	}
	contactCreateRequest.CreatedAt = timestamppb.New(data.ContactEntity.CreatedAt)

	contextWithTimeout, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
		return s.grpcClients.ContactClient.UpsertContact(contextWithTimeout, contactCreateRequest)
	})
	if err != nil {
		s.log.Errorf("(%s) Failed to call method: {%v}", utils.GetFunctionName(), err.Error())
		return nil, err
	}
	result.ID = response.Id

	if data.EmailEntity != nil {
		emailCreate := &emailpb.UpsertEmailGrpcRequest{
			Tenant:   common.GetTenantFromContext(ctx),
			RawEmail: data.EmailEntity.RawEmail,
			SourceFields: &commonpb.SourceFields{
				Source:    string(data.EmailEntity.Source),
				AppSource: data.EmailEntity.AppSource,
			},
			LoggedInUserId: common.GetUserIdFromContext(ctx),
		}
		emailCreate.CreatedAt = timestamppb.New(data.ContactEntity.CreatedAt)
		emailId, err := utils.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
			return s.grpcClients.EmailClient.UpsertEmail(contextWithTimeout, emailCreate)
		})
		if err != nil {
			s.log.Errorf("(%s) Failed to call method: {%v}", utils.GetFunctionName(), err.Error())
			return nil, err
		}

		result.Email = &model.CustomerEmail{
			ID: emailId.Id,
		}
		_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
			return s.grpcClients.ContactClient.LinkEmailToContact(contextWithTimeout, &contactpb.LinkEmailToContactGrpcRequest{
				Primary:   data.EmailEntity.Primary,
				Label:     data.EmailEntity.Label,
				ContactId: response.Id,
				EmailId:   emailId.Id,
				Tenant:    common.GetTenantFromContext(ctx),
				AppSource: data.ContactEntity.AppSource,
			})
		})
		if err != nil {
			s.log.Errorf("(%s) Failed to call method: {%v}", utils.GetFunctionName(), err.Error())
			return nil, err
		}

	}
	return result, nil
}

func (s *contactService) RemoveLocation(ctx context.Context, contactId string, locationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.RemoveLocation")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("locationId", locationId))

	//TODO implement
	panic("implement me")
	//_, err := s.grpcClients.ContactClient.UnlinkLocationFromContact(context.Background(), &contact_grpc_service.UnlinkLocationFromContactGrpcRequest{
	//	Tenant:     common.GetTenantFromContext(ctx),
	//	ContactId:  contactId,
	//	LocationId: locationId,
	//})
}

func (s *contactService) GetContactCountByOrganizations(ctx context.Context, ids []string) (map[string]int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactCountByOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("organizationIds", ids))

	return s.repositories.Neo4jRepositories.ContactReadRepository.GetContactCountByOrganizations(ctx, common.GetTenantFromContext(ctx), ids)
}
