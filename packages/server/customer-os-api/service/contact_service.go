package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"reflect"
)

type ContactService interface {
	Create(ctx context.Context, contact *ContactCreateData) (*entity.ContactEntity, error)
	Update(ctx context.Context, contactUpdateData *ContactUpdateData) (*entity.ContactEntity, error)
	GetContactById(ctx context.Context, id string) (*entity.ContactEntity, error)
	GetFirstContactByEmail(ctx context.Context, email string) (*entity.ContactEntity, error)
	GetFirstContactByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.ContactEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	FindAllForContactGroup(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy, contactGroupId string) (*utils.Pagination, error)
	GetAllForConversation(ctx context.Context, conversationId string) (*entity.ContactEntities, error)
	PermanentDelete(ctx context.Context, id string) (bool, error)
	SoftDelete(ctx context.Context, id string) (bool, error)
	GetContactForRole(ctx context.Context, roleId string) (*entity.ContactEntity, error)
	GetContactsForOrganization(ctx context.Context, organizationId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	Merge(ctx context.Context, primaryContactId, mergedContactId string) error

	AddTag(ctx context.Context, contactId string, tagId string) (*entity.ContactEntity, error)
	RemoveTag(ctx context.Context, contactId string, tagId string) (*entity.ContactEntity, error)

	mapDbNodeToContactEntity(dbNode dbtype.Node) *entity.ContactEntity
}

type ContactCreateData struct {
	ContactEntity     *entity.ContactEntity
	CustomFields      *entity.CustomFieldEntities
	FieldSets         *entity.FieldSetEntities
	EmailEntity       *entity.EmailEntity
	PhoneNumberEntity *entity.PhoneNumberEntity
	TemplateId        *string
	OwnerUserId       *string
	ExternalReference *entity.ExternalReferenceRelationship
	Source            entity.DataSource
	SourceOfTruth     entity.DataSource
}

type ContactUpdateData struct {
	ContactEntity *entity.ContactEntity
	OwnerUserId   *string
}

type contactService struct {
	repositories *repository.Repositories
}

func NewContactService(repositories *repository.Repositories) ContactService {
	return &contactService{
		repositories: repositories,
	}
}

func (s *contactService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *contactService) Create(ctx context.Context, newContact *ContactCreateData) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	contactDbNode, err := session.ExecuteWrite(ctx, s.createContactInDBTxWork(ctx, newContact))
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(*contactDbNode.(*dbtype.Node)), nil
}

func (s *contactService) createContactInDBTxWork(ctx context.Context, newContact *ContactCreateData) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		contactDbNode, err := s.repositories.ContactRepository.Create(ctx, tx, tenant, *newContact.ContactEntity, newContact.Source, newContact.SourceOfTruth)
		if err != nil {
			return nil, err
		}
		var contactId = utils.GetPropsFromNode(*contactDbNode)["id"].(string)

		if newContact.TemplateId != nil {
			err := s.repositories.ContactRepository.LinkWithEntityTemplateInTx(ctx, tx, tenant, contactId, *newContact.TemplateId)
			if err != nil {
				return nil, err
			}
		}
		if newContact.ExternalReference != nil {
			err := s.repositories.ExternalSystemRepository.LinkContactWithExternalSystemInTx(ctx, tx, tenant, contactId, *newContact.ExternalReference)
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
					err := s.repositories.CustomFieldRepository.LinkWithCustomFieldTemplateForContactInTx(ctx, tx, fieldId, contactId, *customField.TemplateId)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if newContact.FieldSets != nil {
			for _, fieldSet := range *newContact.FieldSets {
				setDbNode, err := s.repositories.FieldSetRepository.MergeFieldSetToContactInTx(ctx, tx, tenant, contactId, fieldSet)
				if err != nil {
					return nil, err
				}
				var fieldSetId = utils.GetPropsFromNode(*setDbNode)["id"].(string)
				if fieldSet.TemplateId != nil {
					err := s.repositories.FieldSetRepository.LinkWithFieldSetTemplateInTx(ctx, tx, tenant, fieldSetId, *fieldSet.TemplateId)
					if err != nil {
						return nil, err
					}
				}
				if fieldSet.CustomFields != nil {
					for _, customField := range *fieldSet.CustomFields {
						fieldDbNode, err := s.repositories.CustomFieldRepository.MergeCustomFieldToFieldSetInTx(ctx, tx, tenant, contactId, fieldSetId, customField)
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
			_, _, err := s.repositories.PhoneNumberRepository.MergePhoneNumberToContactInTx(ctx, tx, tenant, contactId, *newContact.PhoneNumberEntity)
			if err != nil {
				return nil, err
			}
		}
		if newContact.OwnerUserId != nil {
			err := s.repositories.ContactRepository.SetOwner(ctx, tx, tenant, contactId, *newContact.OwnerUserId)
			if err != nil {
				return nil, err
			}
		}
		return contactDbNode, nil
	}
}

func (s *contactService) Update(ctx context.Context, contactUpdateData *ContactUpdateData) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	contactDbNode, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		contactId := contactUpdateData.ContactEntity.Id

		dbNode, err := s.repositories.ContactRepository.Update(ctx, tx, tenant, contactId, contactUpdateData.ContactEntity)
		if err != nil {
			return nil, err
		}

		err = s.repositories.ContactRepository.RemoveOwner(ctx, tx, tenant, contactId)
		if err != nil {
			return nil, err
		}
		if contactUpdateData.OwnerUserId != nil {
			err := s.repositories.ContactRepository.SetOwner(ctx, tx, tenant, contactId, *contactUpdateData.OwnerUserId)
			if err != nil {
				return nil, err
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

func (s *contactService) SoftDelete(ctx context.Context, contactId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[r:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MERGE (c)-[:CONTACT_REMOVED_FROM_TENANT {removedAt:datetime({timezone: 'UTC'})}]->(t)
			SET c.removed=true
            DELETE r
			`,
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
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

	contacts := entity.ContactEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		contacts = append(contacts, *s.mapDbNodeToContactEntity(*v))
	}
	paginatedResult.SetRows(&contacts)
	return &paginatedResult, nil
}

func (s *contactService) FindAllForContactGroup(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy, contactGroupId string) (*utils.Pagination, error) {
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

	dbNodesWithTotalCount, err := s.repositories.ContactRepository.GetPaginatedContactsForContactGroup(
		ctx, session,
		common.GetContext(ctx).Tenant,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort,
		contactGroupId)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	contacts := entity.ContactEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		contacts = append(contacts, *s.mapDbNodeToContactEntity(*v))
	}
	paginatedResult.SetRows(&contacts)
	return &paginatedResult, nil
}

func (s *contactService) GetAllForConversation(ctx context.Context, conversationId string) (*entity.ContactEntities, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	dbNodes, err := s.repositories.ContactRepository.GetAllForConversation(ctx, session, common.GetContext(ctx).Tenant, conversationId)
	if err != nil {
		return nil, err
	}

	contactEntities := entity.ContactEntities{}
	for _, dbNode := range dbNodes {
		contactEntities = append(contactEntities, *s.mapDbNodeToContactEntity(*dbNode))
	}
	return &contactEntities, nil
}

func (s *contactService) GetContactForRole(ctx context.Context, roleId string) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	dbNode, err := s.repositories.ContactRepository.GetContactForRole(ctx, session, common.GetContext(ctx).Tenant, roleId)
	if dbNode == nil || err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(*dbNode), nil
}

func (s *contactService) GetContactsForOrganization(ctx context.Context, organizationId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
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

	contacts := entity.ContactEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		contacts = append(contacts, *s.mapDbNodeToContactEntity(*v))
	}
	paginatedResult.SetRows(&contacts)
	return &paginatedResult, nil
}

func (s *contactService) Merge(ctx context.Context, primaryContactId, mergedContactId string) error {
	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	_, err := s.GetContactById(ctx, primaryContactId)
	if err != nil {
		logrus.Errorf("Primary contact with id %s not found: %v", primaryContactId, err)
		return err
	}
	_, err = s.GetContactById(ctx, mergedContactId)
	if err != nil {
		logrus.Errorf("Contact to merge with id %s not found: %v", mergedContactId, err)
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
	return err
}

func (s *contactService) AddTag(ctx context.Context, contactId string, tagId string) (*entity.ContactEntity, error) {
	contactNodePtr, err := s.repositories.ContactRepository.AddTag(ctx, common.GetTenantFromContext(ctx), contactId, tagId)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(*contactNodePtr), nil
}

func (s *contactService) RemoveTag(ctx context.Context, contactId string, tagId string) (*entity.ContactEntity, error) {
	contactNodePtr, err := s.repositories.ContactRepository.RemoveTag(ctx, common.GetTenantFromContext(ctx), contactId, tagId)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(*contactNodePtr), nil
}

func (s *contactService) mapDbNodeToContactEntity(dbNode dbtype.Node) *entity.ContactEntity {
	props := utils.GetPropsFromNode(dbNode)
	contact := entity.ContactEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		FirstName:     utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:      utils.GetStringPropOrEmpty(props, "lastName"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		Title:         utils.GetStringPropOrEmpty(props, "title"),
		CreatedAt:     utils.ToPtr(utils.GetTimePropOrEpochStart(props, "createdAt")),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &contact
}
