package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	contactgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"time"
)

type ContactService interface {
	Create(ctx context.Context, contact *ContactCreateData) (string, error)
	Update(ctx context.Context, contactUpdateData *ContactUpdateData) (string, error)
	GetById(ctx context.Context, id string) (*entity.ContactEntity, error)
	GetFirstContactByEmail(ctx context.Context, email string) (*entity.ContactEntity, error)
	GetFirstContactByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.ContactEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	PermanentDelete(ctx context.Context, id string) (bool, error)
	Archive(ctx context.Context, contactId string) (bool, error)
	RestoreFromArchive(ctx context.Context, contactId string) (bool, error)
	GetContactsForJobRoles(ctx context.Context, jobRoleIds []string) (*entity.ContactEntities, error)
	GetContactsForOrganization(ctx context.Context, organizationId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	Merge(ctx context.Context, primaryContactId, mergedContactId string) error
	GetContactsForEmails(ctx context.Context, emailIds []string) (*entity.ContactEntities, error)
	GetContactsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*entity.ContactEntities, error)
	AddTag(ctx context.Context, contactId, tagId string) (*entity.ContactEntity, error)
	RemoveTag(ctx context.Context, contactId, tagId string) (*entity.ContactEntity, error)
	AddOrganization(ctx context.Context, contactId, organizationId, source, appSource string) (*entity.ContactEntity, error)
	RemoveOrganization(ctx context.Context, contactId, organizationId string) (*entity.ContactEntity, error)
	RemoveLocation(ctx context.Context, contactId string, locationId string) error
	CustomerContactCreate(ctx context.Context, entity *CustomerContactCreateData) (*model.CustomerContact, error)

	mapDbNodeToContactEntity(dbNode dbtype.Node) *entity.ContactEntity
}

type ContactCreateData struct {
	ContactEntity     *entity.ContactEntity
	EmailEntity       *entity.EmailEntity
	PhoneNumberEntity *entity.PhoneNumberEntity
	ExternalReference *entity.ExternalSystemEntity
	Source            entity.DataSource
	AppSource         string
}

type CustomerContactCreateData struct {
	ContactEntity *entity.ContactEntity
	EmailEntity   *entity.EmailEntity
}

type ContactUpdateData struct {
	ContactEntity *entity.ContactEntity
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
	upsertContactRequest := contactgrpc.UpsertContactGrpcRequest{
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
		Name:            contactDetails.ContactEntity.Name,
		Timezone:        contactDetails.ContactEntity.Timezone,
	}
	if contactDetails.ContactEntity.CreatedAt != nil {
		upsertContactRequest.CreatedAt = timestamppb.New(*contactDetails.ContactEntity.CreatedAt)
	}
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
	response, err := s.grpcClients.ContactClient.UpsertContact(ctx, &upsertContactRequest)
	for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
		user, findErr := s.GetById(ctx, response.Id)
		if user != nil && findErr == nil {
			break
		}
		time.Sleep(time.Duration(i*100) * time.Millisecond)
	}
	return response.Id, err
}

func (s *contactService) linkEmailByEvents(ctx context.Context, contactId, appSource string, emailEntity entity.EmailEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.linkEmailByEvents")
	defer span.Finish()

	emailId, err := s.services.EmailService.CreateEmailAddressByEvents(ctx, utils.StringFirstNonEmpty(emailEntity.RawEmail, emailEntity.Email), appSource)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to create email address for contact %s: %s", contactId, err.Error())
	}
	if emailId != "" {
		_, err = s.grpcClients.ContactClient.LinkEmailToContact(ctx, &contactgrpc.LinkEmailToContactGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			ContactId:      contactId,
			EmailId:        emailId,
			Primary:        emailEntity.Primary,
			Label:          emailEntity.Label,
			AppSource:      appSource,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to link email address %s with contact %s: %s", emailId, contactId, err.Error())
		}
	}
}

func (s *contactService) linkPhoneNumberByEvents(ctx context.Context, contactId, appSource string, phoneNumberEntity entity.PhoneNumberEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.linkPhoneNumberByEvents")
	defer span.Finish()

	phoneNumberId, err := s.services.PhoneNumberService.CreatePhoneNumberByEvents(ctx, utils.StringFirstNonEmpty(phoneNumberEntity.RawPhoneNumber, phoneNumberEntity.E164), appSource)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to create phone number for contact %s: %s", contactId, err.Error())
	}
	if phoneNumberId != "" {
		_, err = s.grpcClients.ContactClient.LinkPhoneNumberToContact(ctx, &contactgrpc.LinkPhoneNumberToContactGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			ContactId:      contactId,
			PhoneNumberId:  phoneNumberId,
			Primary:        phoneNumberEntity.Primary,
			Label:          phoneNumberEntity.Label,
			AppSource:      appSource,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to link phone number %s with contact %s: %s", phoneNumberId, contactId, err.Error())
		}
	}
}

func (s *contactService) Update(ctx context.Context, contactUpdateData *ContactUpdateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contactUpdateData", contactUpdateData))

	if contactUpdateData.ContactEntity == nil {
		err := fmt.Errorf("(ContactService.Update) contact entity is nil")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return "", err
	} else if contactUpdateData.ContactEntity.Id == "" {
		err := fmt.Errorf("(ContactService.Update) contact id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return "", err
	}

	currentContactEntity, err := s.GetById(ctx, contactUpdateData.ContactEntity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return "", err
	}

	contactDetails := *contactUpdateData.ContactEntity

	upsertContactRequest := contactgrpc.UpsertContactGrpcRequest{
		Tenant: common.GetTenantFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    string(entity.DataSourceOpenline),
			AppSource: utils.StringFirstNonEmpty(contactDetails.AppSource, constants.AppSourceCustomerOsApi),
		},
		LoggedInUserId:  common.GetUserIdFromContext(ctx),
		Id:              contactDetails.Id,
		Prefix:          contactDetails.Prefix,
		Name:            contactDetails.Name,
		FirstName:       contactDetails.FirstName,
		LastName:        contactDetails.LastName,
		Description:     contactDetails.Description,
		Timezone:        contactDetails.Timezone,
		ProfilePhotoUrl: currentContactEntity.ProfilePhotoUrl,
	}
	response, err := s.grpcClients.ContactClient.UpsertContact(ctx, &upsertContactRequest)
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

func (s *contactService) GetById(ctx context.Context, contactId string) (*entity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))

	if contactDbNode, err := s.repositories.ContactRepository.GetById(ctx, common.GetContext(ctx).Tenant, contactId); err != nil {
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Contact with id {%s} not found", contactId))
		return nil, wrappedErr
	} else {
		return s.mapDbNodeToContactEntity(*contactDbNode), nil
	}
}

func (s *contactService) GetFirstContactByEmail(ctx context.Context, email string) (*entity.ContactEntity, error) {
	dbNodes, err := s.repositories.ContactRepository.GetContactsForEmail(ctx, common.GetContext(ctx).Tenant, email)
	if err != nil || len(dbNodes) == 0 {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(*dbNodes[0]), nil
}

func (s *contactService) GetFirstContactByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.ContactEntity, error) {
	dbNodes, err := s.repositories.ContactRepository.GetContactsForPhoneNumber(ctx, common.GetContext(ctx).Tenant, phoneNumber)
	if err != nil || len(dbNodes) == 0 {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(*dbNodes[0]), nil
}

func (s *contactService) FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.ContactEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(entity.ContactEntity{}))
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

	contacts := make(entity.ContactEntities, 0, len(dbNodesWithTotalCount.Nodes))

	for _, v := range dbNodesWithTotalCount.Nodes {
		contacts = append(contacts, *s.mapDbNodeToContactEntity(*v))
	}
	paginatedResult.SetRows(&contacts)
	return &paginatedResult, nil
}

func (s *contactService) GetContactsForJobRoles(ctx context.Context, jobRoleIds []string) (*entity.ContactEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactsForJobRoles")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("jobRoleIds", jobRoleIds))

	contacts, err := s.repositories.ContactRepository.GetAllForJobRoles(ctx, common.GetTenantFromContext(ctx), jobRoleIds)
	if err != nil {
		return nil, err
	}
	contactEntities := make(entity.ContactEntities, 0, len(contacts))
	for _, v := range contacts {
		contactEntity := s.mapDbNodeToContactEntity(*v.Node)
		contactEntity.DataloaderKey = v.LinkedNodeId
		contactEntities = append(contactEntities, *contactEntity)
	}
	return &contactEntities, nil
}

func (s *contactService) GetContactsForOrganization(ctx context.Context, organizationId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactsForOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
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
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.ContactEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(entity.ContactEntity{}))
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

	contacts := make(entity.ContactEntities, 0, len(dbNodesWithTotalCount.Nodes))
	for _, v := range dbNodesWithTotalCount.Nodes {
		contacts = append(contacts, *s.mapDbNodeToContactEntity(*v))
	}
	paginatedResult.SetRows(&contacts)
	return &paginatedResult, nil
}

func (s *contactService) Merge(ctx context.Context, primaryContactId, mergedContactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
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
		err = s.repositories.ContactRepository.MergeContactPropertiesInTx(ctx, tx, tenant, primaryContactId, mergedContactId, entity.DataSourceOpenline)
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
		s.services.OrganizationService.UpdateLastTouchpointSyncByContactId(ctx, primaryContactId)
	}

	return err
}

func (s *contactService) AddTag(ctx context.Context, contactId, tagId string) (*entity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.AddTag")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("tagId", tagId))

	contactNodePtr, err := s.repositories.ContactRepository.AddTag(ctx, common.GetTenantFromContext(ctx), contactId, tagId)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(*contactNodePtr), nil
}

func (s *contactService) RemoveTag(ctx context.Context, contactId, tagId string) (*entity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.RemoveTag")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("tagId", tagId))

	contactNodePtr, err := s.repositories.ContactRepository.RemoveTag(ctx, common.GetTenantFromContext(ctx), contactId, tagId)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(*contactNodePtr), nil
}

func (s *contactService) AddOrganization(ctx context.Context, contactId, organizationId, source, appSource string) (*entity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.AddOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("organizationId", organizationId))

	contactNodePtr, err := s.repositories.ContactRepository.AddOrganization(ctx, common.GetTenantFromContext(ctx), contactId, organizationId, source, appSource)
	if err != nil {
		return nil, err
	}
	s.services.OrganizationService.UpdateLastTouchpointSync(ctx, organizationId)
	return s.mapDbNodeToContactEntity(*contactNodePtr), nil
}

func (s *contactService) RemoveOrganization(ctx context.Context, contactId, organizationId string) (*entity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.RemoveOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("organizationId", organizationId))

	contactNodePtr, err := s.repositories.ContactRepository.RemoveOrganization(ctx, common.GetTenantFromContext(ctx), contactId, organizationId)
	if err != nil {
		return nil, err
	}
	s.services.OrganizationService.UpdateLastTouchpointSync(ctx, organizationId)
	return s.mapDbNodeToContactEntity(*contactNodePtr), nil
}

func (s *contactService) GetContactsForEmails(ctx context.Context, emailIds []string) (*entity.ContactEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactsForEmails")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("emailIds", emailIds))

	contacts, err := s.repositories.ContactRepository.GetAllForEmails(ctx, common.GetTenantFromContext(ctx), emailIds)
	if err != nil {
		return nil, err
	}
	contactEntities := make(entity.ContactEntities, 0, len(contacts))
	for _, v := range contacts {
		contactEntity := s.mapDbNodeToContactEntity(*v.Node)
		contactEntity.DataloaderKey = v.LinkedNodeId
		contactEntities = append(contactEntities, *contactEntity)
	}
	return &contactEntities, nil
}

func (s *contactService) GetContactsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*entity.ContactEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactsForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("phoneNumberIds", phoneNumberIds))

	contacts, err := s.repositories.ContactRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), phoneNumberIds)
	if err != nil {
		return nil, err
	}
	contactEntities := make(entity.ContactEntities, 0, len(contacts))
	for _, v := range contacts {
		contactEntity := s.mapDbNodeToContactEntity(*v.Node)
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

	contactCreateRequest := &contactgrpc.UpsertContactGrpcRequest{
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
	if data.ContactEntity.CreatedAt != nil {
		contactCreateRequest.CreatedAt = timestamppb.New(*data.ContactEntity.CreatedAt)
	}

	contextWithTimeout, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()
	contactId, err := s.grpcClients.ContactClient.UpsertContact(contextWithTimeout, contactCreateRequest)
	if err != nil {
		s.log.Errorf("(%s) Failed to call method: {%v}", utils.GetFunctionName(), err.Error())
		return nil, err
	}
	result.ID = contactId.Id

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
		if data.ContactEntity.CreatedAt != nil {
			emailCreate.CreatedAt = timestamppb.New(*data.ContactEntity.CreatedAt)
		}
		emailId, err := s.grpcClients.EmailClient.UpsertEmail(contextWithTimeout, emailCreate)
		if err != nil {
			s.log.Errorf("(%s) Failed to call method: {%v}", utils.GetFunctionName(), err.Error())
			return nil, err
		}

		result.Email = &model.CustomerEmail{
			ID: emailId.Id,
		}
		_, err = s.grpcClients.ContactClient.LinkEmailToContact(contextWithTimeout, &contactgrpc.LinkEmailToContactGrpcRequest{
			Primary:   data.EmailEntity.Primary,
			Label:     data.EmailEntity.Label,
			ContactId: contactId.Id,
			EmailId:   emailId.Id,
			Tenant:    common.GetTenantFromContext(ctx),
			AppSource: data.ContactEntity.AppSource,
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

func (s *contactService) mapDbNodeToContactEntity(dbNode dbtype.Node) *entity.ContactEntity {
	props := utils.GetPropsFromNode(dbNode)
	contact := entity.ContactEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		FirstName:       utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:        utils.GetStringPropOrEmpty(props, "lastName"),
		Name:            utils.GetStringPropOrEmpty(props, "name"),
		Description:     utils.GetStringPropOrEmpty(props, "description"),
		Timezone:        utils.GetStringPropOrEmpty(props, "timezone"),
		ProfilePhotoUrl: utils.GetStringPropOrEmpty(props, "profilePhotoUrl"),
		Prefix:          utils.GetStringPropOrEmpty(props, "prefix"),
		CreatedAt:       utils.ToPtr(utils.GetTimePropOrEpochStart(props, "createdAt")),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:          entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &contact
}
