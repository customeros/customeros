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
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type ContactService interface {
	Create(ctx context.Context, contact *ContactCreateData) (string, error)
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

	externalSystem := neo4jmodel.ExternalSystem{}
	if contactDetails.ExternalReference != nil && contactDetails.ExternalReference.ExternalSystemId != "" {
		externalSystem = neo4jmodel.ExternalSystem{
			ExternalSystemId: string(contactDetails.ExternalReference.ExternalSystemId),
			ExternalId:       contactDetails.ExternalReference.Relationship.ExternalId,
			ExternalUrl:      utils.IfNotNilString(contactDetails.ExternalReference.Relationship.ExternalUrl),
			ExternalSource:   utils.IfNotNilString(contactDetails.ExternalReference.Relationship.ExternalSource),
		}
		if contactDetails.ExternalReference.Relationship.SyncDate != nil {
			externalSystem.SyncDate = contactDetails.ExternalReference.Relationship.SyncDate
		}
	}

	contactId, err := s.services.CommonServices.ContactService.SaveContact(ctx, nil,
		neo4jrepository.ContactFields{
			SourceFields: neo4jmodel.SourceFields{
				Source:    string(contactDetails.Source),
				AppSource: utils.StringFirstNonEmpty(contactDetails.AppSource, constants.AppSourceCustomerOsApi),
			},
			FirstName:       contactDetails.ContactEntity.FirstName,
			LastName:        contactDetails.ContactEntity.LastName,
			Prefix:          contactDetails.ContactEntity.Prefix,
			Description:     contactDetails.ContactEntity.Description,
			ProfilePhotoUrl: contactDetails.ContactEntity.ProfilePhotoUrl,
			Username:        contactDetails.ContactEntity.Username,
			Name:            contactDetails.ContactEntity.Name,
			Timezone:        contactDetails.ContactEntity.Timezone,
		}, contactDetails.SocialUrl, externalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to create contact: %s", err.Error())
		return "", err
	}

	if contactDetails.EmailEntity != nil {
		_, err := s.services.CommonServices.EmailService.Merge(ctx, common.GetTenantFromContext(ctx),
			commonservice.EmailFields{
				Email:     strings.TrimSpace(utils.FirstNotEmptyString(contactDetails.EmailEntity.Email, contactDetails.EmailEntity.RawEmail)),
				Primary:   utils.IfNotNilBool(contactDetails.EmailEntity.Primary),
				Source:    neo4jentity.DataSourceOpenline,
				AppSource: constants.AppSourceCustomerOsApi,
			}, &commonservice.LinkWith{
				Type: commonModel.CONTACT,
				Id:   contactId,
			})
		if err != nil {
			tracing.TraceErr(span, err)
			return contactId, err
		}
	}

	if contactDetails.PhoneNumberEntity != nil {
		s.linkPhoneNumberByEvents(ctx, contactId, utils.StringFirstNonEmpty(contactDetails.PhoneNumberEntity.AppSource, contactDetails.AppSource), *contactDetails.PhoneNumberEntity)
	}

	span.LogFields(log.String("output - createdContactId", contactId))
	return contactId, nil
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

	utils.EventCompleted(ctx, common.GetTenantFromContext(ctx), commonModel.CONTACT.String(), contactId, s.grpcClients, utils.NewEventCompletedDetails().WithUpdate())
	utils.EventCompleted(ctx, common.GetTenantFromContext(ctx), commonModel.ORGANIZATION.String(), organizationId, s.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

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

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	contactId, err := s.services.CommonServices.ContactService.SaveContact(ctx, nil,
		neo4jrepository.ContactFields{
			FirstName:   data.ContactEntity.FirstName,
			LastName:    data.ContactEntity.LastName,
			Prefix:      data.ContactEntity.Prefix,
			Description: data.ContactEntity.Description,
			SourceFields: neo4jmodel.SourceFields{
				Source:    string(data.ContactEntity.Source),
				AppSource: data.ContactEntity.AppSource,
			},
		}, "", neo4jmodel.ExternalSystem{})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	result.ID = contactId

	if data.EmailEntity != nil {
		emailId, err := s.services.CommonServices.EmailService.Merge(ctx, common.GetTenantFromContext(ctx),
			commonservice.EmailFields{
				Email:     strings.TrimSpace(utils.FirstNotEmptyString(data.EmailEntity.Email, data.EmailEntity.RawEmail)),
				Primary:   utils.IfNotNilBool(data.EmailEntity.Primary),
				Source:    neo4jentity.DataSourceOpenline,
				AppSource: constants.AppSourceCustomerOsApi,
			}, &commonservice.LinkWith{
				Type: commonModel.CONTACT,
				Id:   contactId,
			})
		if err != nil {
			tracing.TraceErr(span, err)
			return result, err
		}
		result.Email = &model.CustomerEmail{
			ID: utils.IfNotNilString(emailId),
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
	//	EntityId:  contactId,
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
