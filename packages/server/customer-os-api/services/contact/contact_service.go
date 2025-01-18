package api_contact

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
	api_filters "github.com/customeros/customeros/packages/server/customer-os-api/services/filters"
	api_sort "github.com/customeros/customeros/packages/server/customer-os-api/services/sort"
)

type contactService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	contact      interfaces.ContactService
	email        interfaces.EmailService
	phoneNumber  interfaces.PhoneNumberService
	org          cosapi_interfaces.OrganizationService
}

func NewContactService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, contact interfaces.ContactService, email interfaces.EmailService, phoneNumber interfaces.PhoneNumberService, org cosapi_interfaces.OrganizationService) cosapi_interfaces.ContactService {
	return &contactService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		contact:      contact,
		email:        email,
		phoneNumber:  phoneNumber,
		org:          org,
	}
}

func (s *contactService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *contactService) Create(ctx context.Context, contactDetails *cosapi_interfaces.ContactCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contactDetails", contactDetails))

	tenant := common.GetTenantFromContext(ctx)

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
	contactId := ""
	var err error

	if (neo4jentity.SocialEntity{Url: contactDetails.SocialUrl}).IsLinkedin() {
		contactId, err = s.contact.CreateContactByLinkedIn(ctx, nil, contactDetails.SocialUrl)
	} else {
		contactId, err = s.contact.Save(ctx, nil, nil,
			data_fields.ContactFields{
				Source:          utils.StringPtr(string(contactDetails.Source)),
				AppSource:       utils.StringPtr(utils.StringFirstNonEmpty(contactDetails.AppSource, constants.AppSourceCustomerOsApi)),
				FirstName:       utils.StringPtr(contactDetails.ContactEntity.FirstName),
				LastName:        utils.StringPtr(contactDetails.ContactEntity.LastName),
				Prefix:          utils.StringPtr(contactDetails.ContactEntity.Prefix),
				Description:     utils.StringPtr(contactDetails.ContactEntity.Description),
				ProfilePhotoUrl: utils.StringPtr(contactDetails.ContactEntity.ProfilePhotoUrl),
				Username:        utils.StringPtr(contactDetails.ContactEntity.Username),
				Name:            utils.StringPtr(contactDetails.ContactEntity.Name),
				Timezone:        utils.StringPtr(contactDetails.ContactEntity.Timezone),
			}, false)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to create contact: %s", err.Error())
			return "", err
		}
	}

	if contactDetails.EmailEntity != nil {
		_, err := s.email.Merge(ctx, nil, common.GetTenantFromContext(ctx),
			interfaces.EmailFields{
				Email:     strings.TrimSpace(utils.FirstNotEmptyString(contactDetails.EmailEntity.Email, contactDetails.EmailEntity.RawEmail)),
				Primary:   utils.IfNotNilBool(contactDetails.EmailEntity.Primary),
				Source:    neo4jentity.DataSourceOpenline,
				AppSource: constants.AppSourceCustomerOsApi,
			}, &common_srv.LinkWith{
				Type: commonModel.CONTACT,
				Id:   contactId,
			})
		if err != nil {
			tracing.TraceErr(span, err)
			return contactId, err
		}
	}

	if contactDetails.PhoneNumberEntity != nil {
		phoneNumberId, err := s.phoneNumber.Merge(ctx, contactDetails.PhoneNumberEntity.RawPhoneNumber, neo4jentity.DataSourceOpenline)
		if err != nil {
			tracing.TraceErr(span, err)
			return contactId, err
		}

		err = s.repositories.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithContact(ctx, tenant, contactId, phoneNumberId, contactDetails.PhoneNumberEntity.Label, contactDetails.PhoneNumberEntity.Primary)
		if err != nil {
			tracing.TraceErr(span, err)
			return contactId, err
		}
	}

	span.LogFields(log.String("output - createdContactId", contactId))
	return contactId, nil
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

func (s *contactService) GetFirstContactByPhoneNumber(ctx context.Context, phoneNumber string) (*neo4jentity.ContactEntity, error) {
	dbNodes, err := s.repositories.ContactRepository.GetContactsForPhoneNumber(ctx, common.GetContext(ctx).Tenant, phoneNumber)
	if err != nil || len(dbNodes) == 0 {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToContactEntity(dbNodes[0]), nil
}

func (s *contactService) FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*commonModel.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	paginatedResult := utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := api_sort.BuildSort(sortBy, reflect.TypeOf(neo4jentity.ContactEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := api_filters.BuildFilter(filter, reflect.TypeOf(neo4jentity.ContactEntity{}))
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

func (s *contactService) GetContactsForOrganization(ctx context.Context, organizationId string, page, limit int, filter *model.Filter, sortBy []*commonModel.SortBy) (*utils.Pagination, error) {
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

	paginatedResult := utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := api_sort.BuildSort(sortBy, reflect.TypeOf(neo4jentity.ContactEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := api_filters.BuildFilter(filter, reflect.TypeOf(neo4jentity.ContactEntity{}))
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
		s.org.UpdateLastTouchpointByContactId(ctx, primaryContactId)
	}

	return err
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

func (s *contactService) CustomerContactCreate(ctx context.Context, data *cosapi_interfaces.CustomerContactCreateData) (*model.CustomerContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.CustomerContactCreate")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	result := &model.CustomerContact{}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	contactId, err := s.contact.Save(ctx, nil, nil,
		data_fields.ContactFields{
			FirstName:   utils.StringPtr(data.ContactEntity.FirstName),
			LastName:    utils.StringPtr(data.ContactEntity.LastName),
			Prefix:      utils.StringPtr(data.ContactEntity.Prefix),
			Description: utils.StringPtr(data.ContactEntity.Description),
			Source:      utils.StringPtr(string(data.ContactEntity.Source)),
			AppSource:   utils.StringPtr(data.ContactEntity.AppSource),
		}, false)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	result.ID = contactId

	if data.EmailEntity != nil {
		emailId, err := s.email.Merge(ctx, nil, common.GetTenantFromContext(ctx),
			interfaces.EmailFields{
				Email:     strings.TrimSpace(utils.FirstNotEmptyString(data.EmailEntity.Email, data.EmailEntity.RawEmail)),
				Primary:   utils.IfNotNilBool(data.EmailEntity.Primary),
				Source:    neo4jentity.DataSourceOpenline,
				AppSource: constants.AppSourceCustomerOsApi,
			}, &common_srv.LinkWith{
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

	// TODO implement
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
