package service

import (
	"context"
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
	contact_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	email_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
)

type ContactService interface {
	Create(ctx context.Context, contact *ContactCreateData) (*entity.ContactEntity, error)
	Update(ctx context.Context, contactUpdateData *ContactUpdateData) (*entity.ContactEntity, error)
	GetContactById(ctx context.Context, id string) (*entity.ContactEntity, error)
	GetFirstContactByEmail(ctx context.Context, email string) (*entity.ContactEntity, error)
	GetFirstContactByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.ContactEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	GetAllForConversation(ctx context.Context, conversationId string) (*entity.ContactEntities, error)
	PermanentDelete(ctx context.Context, id string) (bool, error)
	Archive(ctx context.Context, contactId string) (bool, error)
	RestoreFromArchive(ctx context.Context, contactId string) (bool, error)
	GetContactForRole(ctx context.Context, roleId string) (*entity.ContactEntity, error)
	GetContactsForOrganization(ctx context.Context, organizationId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	Merge(ctx context.Context, primaryContactId, mergedContactId string) error
	GetContactsForEmails(ctx context.Context, emailIds []string) (*entity.ContactEntities, error)
	GetContactsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*entity.ContactEntities, error)
	AddTag(ctx context.Context, contactId, tagId string) (*entity.ContactEntity, error)
	RemoveTag(ctx context.Context, contactId, tagId string) (*entity.ContactEntity, error)
	AddOrganization(ctx context.Context, contactId, organizationId, source, appSource string) (*entity.ContactEntity, error)
	RemoveOrganization(ctx context.Context, contactId, organizationId string) (*entity.ContactEntity, error)
	RemoveLocation(ctx context.Context, contactId string, locationId string) error

	mapDbNodeToContactEntity(dbNode dbtype.Node) *entity.ContactEntity

	// Deprecated
	UpsertPhoneNumberRelationInEventStore(ctx context.Context, size int) (int, int, error)
	// Deprecated
	UpsertEmailRelationInEventStore(ctx context.Context, size int) (int, int, error)
	CustomerContactCreate(ctx context.Context, entity *CustomerContactCreateData) (*model.CustomerContact, error)
}

type ContactCreateData struct {
	ContactEntity     *entity.ContactEntity
	CustomFields      *entity.CustomFieldEntities
	FieldSets         *entity.FieldSetEntities
	EmailEntity       *entity.EmailEntity
	PhoneNumberEntity *entity.PhoneNumberEntity
	TemplateId        *string
	OwnerUserId       *string
	ExternalReference *entity.ExternalSystemEntity
	Source            entity.DataSource
	SourceOfTruth     entity.DataSource
}

type CustomerContactCreateData struct {
	ContactEntity *entity.ContactEntity
	EmailEntity   *entity.EmailEntity
}

type ContactUpdateData struct {
	ContactEntity *entity.ContactEntity
	OwnerUserId   *string
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

func (s *contactService) Create(ctx context.Context, newContact *ContactCreateData) (*entity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("newContact", newContact))

	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	contactDbNode, err := session.ExecuteWrite(ctx, s.createContactInDBTxWork(ctx, newContact))
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(*contactDbNode.(*dbtype.Node)), nil
}

func (s *contactService) createContactInDBTxWork(ctx context.Context, newContact *ContactCreateData) func(tx neo4j.ManagedTransaction) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.createContactInDBTxWork")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("newContact", newContact))

	return func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		contactDbNode, err := s.repositories.ContactRepository.Create(ctx, tx, tenant, *newContact.ContactEntity)
		if err != nil {
			return nil, err
		}
		var contactId = utils.GetPropsFromNode(*contactDbNode)["id"].(string)
		entityType := &model.CustomFieldEntityType{
			ID:         contactId,
			EntityType: model.EntityTypeContact,
		}
		if newContact.TemplateId != nil {
			err := s.repositories.ContactRepository.LinkWithEntityTemplateInTx(ctx, tx, tenant, entityType, *newContact.TemplateId)
			if err != nil {
				return nil, err
			}
		}
		if newContact.ExternalReference != nil {
			err := s.repositories.ExternalSystemRepository.LinkNodeWithExternalSystemInTx(ctx, tx, tenant, contactId, entity.ExternalNodeContact, *newContact.ExternalReference)
			if err != nil {
				return nil, err
			}
		}
		if newContact.CustomFields != nil {
			for _, customField := range *newContact.CustomFields {
				dbNode, err := s.repositories.CustomFieldRepository.MergeCustomFieldToContactInTx(ctx, tx, tenant, contactId, customField)
				if err != nil {
					return nil, err
				}
				if customField.TemplateId != nil {
					var fieldId = utils.GetPropsFromNode(*dbNode)["id"].(string)
					err := s.repositories.CustomFieldRepository.LinkWithCustomFieldTemplateInTx(ctx, tx, fieldId, entityType, *customField.TemplateId)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if newContact.FieldSets != nil {
			for _, fieldSet := range *newContact.FieldSets {
				setDbNode, err := s.repositories.FieldSetRepository.MergeFieldSetInTx(ctx, tx, tenant, entityType, fieldSet)
				if err != nil {
					return nil, err
				}
				var fieldSetId = utils.GetPropsFromNode(*setDbNode)["id"].(string)
				if fieldSet.TemplateId != nil {
					err := s.repositories.FieldSetRepository.LinkWithFieldSetTemplateInTx(ctx, tx, tenant, fieldSetId, *fieldSet.TemplateId, entityType.EntityType)
					if err != nil {
						return nil, err
					}
				}
				if fieldSet.CustomFields != nil {
					for _, customField := range *fieldSet.CustomFields {
						fieldDbNode, err := s.repositories.CustomFieldRepository.MergeCustomFieldToFieldSetInTx(ctx, tx, tenant, entityType, fieldSetId, customField)
						if err != nil {
							return nil, err
						}
						if customField.TemplateId != nil {
							var fieldId = utils.GetPropsFromNode(*fieldDbNode)["id"].(string)
							err := s.repositories.CustomFieldRepository.LinkWithCustomFieldTemplateForFieldSetInTx(ctx, tx, fieldId, fieldSetId, *customField.TemplateId)
							if err != nil {
								return nil, err
							}
						}
					}
				}
			}
		}
		if newContact.EmailEntity != nil {
			_, _, err := s.repositories.EmailRepository.MergeEmailToInTx(ctx, tx, tenant, entity.CONTACT, contactId, *newContact.EmailEntity)
			if err != nil {
				return nil, err
			}
		}
		if newContact.PhoneNumberEntity != nil {
			_, _, err := s.repositories.PhoneNumberRepository.MergePhoneNumberToInTx(ctx, tx, tenant, entity.CONTACT, contactId, *newContact.PhoneNumberEntity)
			if err != nil {
				return nil, err
			}
		}

		var ownerUserId = common.GetUserIdFromContext(ctx)
		if newContact.OwnerUserId != nil && *newContact.OwnerUserId != "" {
			ownerUserId = *newContact.OwnerUserId
		}
		if len(ownerUserId) > 0 {
			err := s.repositories.ContactRepository.SetOwner(ctx, tx, tenant, contactId, ownerUserId)
			if err != nil {
				return nil, err
			}
		}
		return contactDbNode, nil
	}
}

func (s *contactService) Update(ctx context.Context, contactUpdateData *ContactUpdateData) (*entity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contactUpdateData", contactUpdateData))

	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	contactDbNode, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		contactId := contactUpdateData.ContactEntity.Id

		dbNode, err := s.repositories.ContactRepository.Update(ctx, tx, tenant, contactId, contactUpdateData.ContactEntity)
		if err != nil {
			return nil, err
		}

		if contactUpdateData.OwnerUserId != nil {
			err = s.repositories.ContactRepository.RemoveOwner(ctx, tx, tenant, contactId)
			if err != nil {
				return nil, err
			}

			if len(*contactUpdateData.OwnerUserId) > 0 {
				err := s.repositories.ContactRepository.SetOwner(ctx, tx, tenant, contactId, *contactUpdateData.OwnerUserId)
				if err != nil {
					return nil, err
				}
			}
		}

		return dbNode, nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactEntity(*contactDbNode.(*dbtype.Node)), nil
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

func (s *contactService) GetContactById(ctx context.Context, id string) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) RETURN c`,
			map[string]interface{}{
				"id":     id,
				"tenant": common.GetContext(ctx).Tenant,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactEntity(queryResult.(dbtype.Node)), nil
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

func (s *contactService) GetAllForConversation(ctx context.Context, conversationId string) (*entity.ContactEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetAllForConversation")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("conversationId", conversationId))

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	dbNodes, err := s.repositories.ContactRepository.GetAllForConversation(ctx, session, common.GetContext(ctx).Tenant, conversationId)
	if err != nil {
		return nil, err
	}

	contactEntities := make(entity.ContactEntities, 0, len(dbNodes))
	for _, dbNode := range dbNodes {
		contactEntities = append(contactEntities, *s.mapDbNodeToContactEntity(*dbNode))
	}
	return &contactEntities, nil
}

func (s *contactService) GetContactForRole(ctx context.Context, roleId string) (*entity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactForRole")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("roleId", roleId))

	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	dbNode, err := s.repositories.ContactRepository.GetContactForRole(ctx, session, common.GetContext(ctx).Tenant, roleId)
	if dbNode == nil || err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(*dbNode), nil
}

func (s *contactService) GetContactsForOrganization(ctx context.Context, organizationId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactsForOrganization")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("primaryContactId", primaryContactId), log.String("mergedContactId", mergedContactId))

	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	_, err := s.GetContactById(ctx, primaryContactId)
	if err != nil {
		s.log.Errorf("(%s) Primary contact with id {%s} not found: {%v}", utils.GetFunctionName(), primaryContactId, err.Error())
		return err
	}
	_, err = s.GetContactById(ctx, mergedContactId)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.AddTag")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.RemoveTag")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.AddOrganization")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.RemoveOrganization")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactsForEmails")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactsForPhoneNumbers")
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

func (s *contactService) UpsertPhoneNumberRelationInEventStore(ctx context.Context, size int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	outputErr := error(nil)
	for size > 0 {
		batchSize := constants.Neo4jBatchSize
		if size < constants.Neo4jBatchSize {
			batchSize = size
		}
		records, err := s.repositories.ContactRepository.GetAllContactPhoneNumberRelationships(ctx, batchSize)
		if err != nil {
			return 0, 0, err
		}
		for _, v := range records {
			_, err := s.grpcClients.ContactClient.LinkPhoneNumberToContact(context.Background(), &contact_grpc_service.LinkPhoneNumberToContactGrpcRequest{
				Primary:       utils.GetBoolPropOrFalse(v.Values[0].(neo4j.Relationship).Props, "primary"),
				Label:         utils.GetStringPropOrEmpty(v.Values[0].(neo4j.Relationship).Props, "label"),
				ContactId:     v.Values[1].(string),
				PhoneNumberId: v.Values[2].(string),
				Tenant:        v.Values[3].(string),
			})
			if err != nil {
				failedRecords++
				if outputErr != nil {
					outputErr = err
				}
				s.log.Errorf("(%s) Failed to call method: {%v}", utils.GetFunctionName(), err.Error())
			} else {
				processedRecords++
			}
		}

		size -= batchSize
	}

	return processedRecords, failedRecords, outputErr
}

func (s *contactService) UpsertEmailRelationInEventStore(ctx context.Context, size int) (int, int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.UpsertEmailRelationInEventStore")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	processedRecords := 0
	failedRecords := 0
	outputErr := error(nil)
	for size > 0 {
		batchSize := constants.Neo4jBatchSize
		if size < constants.Neo4jBatchSize {
			batchSize = size
		}
		records, err := s.repositories.ContactRepository.GetAllContactEmailRelationships(ctx, batchSize)
		if err != nil {
			return 0, 0, err
		}
		for _, v := range records {
			_, err := s.grpcClients.ContactClient.LinkEmailToContact(context.Background(), &contact_grpc_service.LinkEmailToContactGrpcRequest{
				Primary:   utils.GetBoolPropOrFalse(v.Values[0].(neo4j.Relationship).Props, "primary"),
				Label:     utils.GetStringPropOrEmpty(v.Values[0].(neo4j.Relationship).Props, "label"),
				ContactId: v.Values[1].(string),
				EmailId:   v.Values[2].(string),
				Tenant:    v.Values[3].(string),
			})
			if err != nil {
				failedRecords++
				if outputErr != nil {
					outputErr = err
				}
				s.log.Errorf("(%s) Failed to call method: {%v}", utils.GetFunctionName(), err.Error())
			} else {
				processedRecords++
			}
		}

		size -= batchSize
	}

	return processedRecords, failedRecords, outputErr
}

func (s *contactService) CustomerContactCreate(ctx context.Context, data *CustomerContactCreateData) (*model.CustomerContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.CustomerContactCreate")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	result := &model.CustomerContact{}

	contactCreate := &contact_grpc_service.CreateContactGrpcRequest{
		Tenant:        common.GetTenantFromContext(ctx),
		FirstName:     data.ContactEntity.FirstName,
		LastName:      data.ContactEntity.LastName,
		Prefix:        data.ContactEntity.Prefix,
		Description:   data.ContactEntity.Description,
		Source:        string(data.ContactEntity.Source),
		SourceOfTruth: string(data.ContactEntity.SourceOfTruth),
		AppSource:     data.ContactEntity.AppSource,
	}
	if data.ContactEntity.CreatedAt != nil {
		contactCreate.CreatedAt = timestamppb.New(*data.ContactEntity.CreatedAt)
	}

	contextWithTimeout, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()
	contactId, err := s.grpcClients.ContactClient.CreateContact(contextWithTimeout, contactCreate)
	if err != nil {
		s.log.Errorf("(%s) Failed to call method: {%v}", utils.GetFunctionName(), err.Error())
		return nil, err
	}
	result.ID = contactId.Id

	if data.EmailEntity != nil {
		emailCreate := &email_grpc_service.UpsertEmailGrpcRequest{
			Tenant:        common.GetTenantFromContext(ctx),
			RawEmail:      data.EmailEntity.RawEmail,
			AppSource:     data.EmailEntity.AppSource,
			Source:        string(data.EmailEntity.Source),
			SourceOfTruth: string(data.EmailEntity.SourceOfTruth),
			Id:            "",
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
		_, err = s.grpcClients.ContactClient.LinkEmailToContact(contextWithTimeout, &contact_grpc_service.LinkEmailToContactGrpcRequest{
			Primary:   data.EmailEntity.Primary,
			Label:     data.EmailEntity.Label,
			ContactId: contactId.Id,
			EmailId:   emailId.Id,
			Tenant:    common.GetTenantFromContext(ctx),
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

	_, err := s.grpcClients.ContactClient.UnlinkLocationFromContact(context.Background(), &contact_grpc_service.UnlinkLocationFromContactGrpcRequest{
		Tenant:     common.GetTenantFromContext(ctx),
		ContactId:  contactId,
		LocationId: locationId,
	})

	return err
}

func (s *contactService) mapDbNodeToContactEntity(dbNode dbtype.Node) *entity.ContactEntity {
	props := utils.GetPropsFromNode(dbNode)
	contact := entity.ContactEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		FirstName:     utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:      utils.GetStringPropOrEmpty(props, "lastName"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		Description:   utils.GetStringPropOrEmpty(props, "description"),
		Timezone:      utils.GetStringPropOrEmpty(props, "timezone"),
		Prefix:        utils.GetStringPropOrEmpty(props, "prefix"),
		CreatedAt:     utils.ToPtr(utils.GetTimePropOrEpochStart(props, "createdAt")),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &contact
}
