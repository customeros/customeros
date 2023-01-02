package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"reflect"
)

type ContactService interface {
	Create(ctx context.Context, contact *ContactCreateData) (*entity.ContactEntity, error)
	Update(ctx context.Context, contactUpdateData *ContactUpdateData) (*entity.ContactEntity, error)
	FindContactById(ctx context.Context, id string) (*entity.ContactEntity, error)
	FindContactByEmail(ctx context.Context, email string) (*entity.ContactEntity, error)
	FindContactByPhoneNumber(ctx context.Context, e164 string) (*entity.ContactEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	FindAllForContactGroup(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy, contactGroupId string) (*utils.Pagination, error)
	PermanentDelete(ctx context.Context, id string) (bool, error)
	SoftDelete(ctx context.Context, id string) (bool, error)
}

type ContactCreateData struct {
	ContactEntity     *entity.ContactEntity
	CustomFields      *entity.CustomFieldEntities
	FieldSets         *entity.FieldSetEntities
	EmailEntity       *entity.EmailEntity
	PhoneNumberEntity *entity.PhoneNumberEntity
	TemplateId        *string
	ContactTypeId     *string
	OwnerUserId       *string
	ExternalReference *entity.ExternalReferenceRelationship
}

type ContactUpdateData struct {
	ContactEntity *entity.ContactEntity
	ContactTypeId *string
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

func (s *contactService) getNeo4jDriver() neo4j.Driver {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *contactService) Create(ctx context.Context, newContact *ContactCreateData) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getNeo4jDriver())
	defer session.Close()

	contactDbNode, err := session.WriteTransaction(s.createContactInDBTxWork(ctx, newContact))
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(contactDbNode.(*dbtype.Node)), nil
}

func (s *contactService) createContactInDBTxWork(ctx context.Context, newContact *ContactCreateData) func(tx neo4j.Transaction) (any, error) {
	return func(tx neo4j.Transaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		contactDbNode, err := s.repositories.ContactRepository.Create(tx, tenant, *newContact.ContactEntity)
		if err != nil {
			return nil, err
		}
		var contactId = utils.GetPropsFromNode(*contactDbNode)["id"].(string)

		if newContact.ContactTypeId != nil {
			err := s.repositories.ContactRepository.LinkWithContactTypeInTx(tx, tenant, contactId, *newContact.ContactTypeId)
			if err != nil {
				return nil, err
			}
		}
		if newContact.TemplateId != nil {
			err := s.repositories.ContactRepository.LinkWithEntityTemplateInTx(tx, tenant, contactId, *newContact.TemplateId)
			if err != nil {
				return nil, err
			}
		}
		if newContact.ExternalReference != nil {
			err := s.repositories.ExternalSystemRepository.LinkContactWithExternalSystemInTx(tx, tenant, contactId, *newContact.ExternalReference)
			if err != nil {
				return nil, err
			}
		}
		if newContact.CustomFields != nil {
			for _, customField := range *newContact.CustomFields {
				dbNode, err := s.repositories.CustomFieldRepository.MergeCustomFieldToContactInTx(tx, tenant, contactId, &customField)
				if err != nil {
					return nil, err
				}
				if customField.TemplateId != nil {
					var fieldId = utils.GetPropsFromNode(*dbNode)["id"].(string)
					err := s.repositories.CustomFieldRepository.LinkWithCustomFieldTemplateForContactInTx(tx, fieldId, contactId, *customField.TemplateId)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if newContact.FieldSets != nil {
			for _, fieldSet := range *newContact.FieldSets {
				setDbNode, _, err := s.repositories.FieldSetRepository.MergeFieldSetToContactInTx(tx, tenant, contactId, fieldSet)
				if err != nil {
					return nil, err
				}
				var fieldSetId = utils.GetPropsFromNode(*setDbNode)["id"].(string)
				if fieldSet.TemplateId != nil {
					err := s.repositories.FieldSetRepository.LinkWithFieldSetTemplateInTx(tx, tenant, fieldSetId, *fieldSet.TemplateId)
					if err != nil {
						return nil, err
					}
				}
				if fieldSet.CustomFields != nil {
					for _, customField := range *fieldSet.CustomFields {
						fieldDbNode, err := s.repositories.CustomFieldRepository.MergeCustomFieldToFieldSetInTx(tx, tenant, contactId, fieldSetId, &customField)
						if err != nil {
							return nil, err
						}
						if customField.TemplateId != nil {
							var fieldId = utils.GetPropsFromNode(*fieldDbNode)["id"].(string)
							err := s.repositories.CustomFieldRepository.LinkWithCustomFieldTemplateForFieldSetInTx(tx, fieldId, fieldSetId, *customField.TemplateId)
							if err != nil {
								return nil, err
							}
						}
					}
				}
			}
		}
		if newContact.EmailEntity != nil {
			_, _, err := s.repositories.EmailRepository.MergeEmailToContactInTx(tx, tenant, contactId, *newContact.EmailEntity)
			if err != nil {
				return nil, err
			}
		}
		if newContact.PhoneNumberEntity != nil {
			_, _, err := s.repositories.PhoneNumberRepository.MergePhoneNumberToContactInTx(tx, tenant, contactId, *newContact.PhoneNumberEntity)
			if err != nil {
				return nil, err
			}
		}
		if newContact.OwnerUserId != nil {
			err := s.repositories.ContactRepository.SetOwner(tx, tenant, contactId, *newContact.OwnerUserId)
			if err != nil {
				return nil, err
			}
		}
		return contactDbNode, nil
	}
}

func (s *contactService) Update(ctx context.Context, contactUpdateData *ContactUpdateData) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getNeo4jDriver())
	defer session.Close()

	contactDbNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		contactId := contactUpdateData.ContactEntity.Id

		dbNode, err := s.repositories.ContactRepository.Update(tx, tenant, contactId, contactUpdateData.ContactEntity)
		if err != nil {
			return nil, err
		}
		err = s.repositories.ContactRepository.UnlinkFromContactTypesInTx(tx, tenant, contactId)
		if err != nil {
			return nil, err
		}
		if contactUpdateData.ContactTypeId != nil {
			err := s.repositories.ContactRepository.LinkWithContactTypeInTx(tx, tenant, contactId, *contactUpdateData.ContactTypeId)
			if err != nil {
				return nil, err
			}
		}

		err = s.repositories.ContactRepository.RemoveOwner(tx, tenant, contactId)
		if err != nil {
			return nil, err
		}
		if contactUpdateData.OwnerUserId != nil {
			err := s.repositories.ContactRepository.SetOwner(tx, tenant, contactId, *contactUpdateData.OwnerUserId)
			if err != nil {
				return nil, err
			}
		}

		return dbNode, nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactEntity(contactDbNode.(*dbtype.Node)), nil
}

func (s *contactService) PermanentDelete(ctx context.Context, contactId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getNeo4jDriver())
	defer session.Close()

	err := s.repositories.ContactRepository.Delete(session, common.GetContext(ctx).Tenant, contactId)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *contactService) SoftDelete(ctx context.Context, contactId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getNeo4jDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
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

func (s *contactService) FindContactById(ctx context.Context, id string) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jReadSession(s.getNeo4jDriver())
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) RETURN c`,
			map[string]interface{}{
				"id":     id,
				"tenant": common.GetContext(ctx).Tenant,
			})
		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactEntity(utils.NodePtr(queryResult.(dbtype.Node))), nil
}

func (s *contactService) FindContactByEmail(ctx context.Context, email string) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jReadSession(s.getNeo4jDriver())
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (:Email {email:$email})<-[:EMAILED_AT]-(c:Contact),
					(c)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
			RETURN c`,
			map[string]interface{}{
				"email":  email,
				"tenant": common.GetContext(ctx).Tenant,
			})
		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactEntity(utils.NodePtr(queryResult.(dbtype.Node))), nil
}

func (s *contactService) FindContactByPhoneNumber(ctx context.Context, e164 string) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jReadSession(s.getNeo4jDriver())
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (:PhoneNumber {e164:$e164})<-[:CALLED_AT]-(c:Contact),
					(c)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
			RETURN c`,
			map[string]interface{}{
				"e164":   e164,
				"tenant": common.GetContext(ctx).Tenant,
			})
		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactEntity(utils.NodePtr(queryResult.(dbtype.Node))), nil
}

func (s *contactService) FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(s.getNeo4jDriver())
	defer session.Close()

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
		session,
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
		contacts = append(contacts, *s.mapDbNodeToContactEntity(v))
	}
	paginatedResult.SetRows(&contacts)
	return &paginatedResult, nil
}

func (s *contactService) FindAllForContactGroup(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy, contactGroupId string) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(s.getNeo4jDriver())
	defer session.Close()

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
		session,
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
		contacts = append(contacts, *s.mapDbNodeToContactEntity(v))
	}
	paginatedResult.SetRows(&contacts)
	return &paginatedResult, nil
}

func (s *contactService) mapDbNodeToContactEntity(dbContactNode *dbtype.Node) *entity.ContactEntity {
	props := utils.GetPropsFromNode(*dbContactNode)
	contact := entity.ContactEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		FirstName: utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:  utils.GetStringPropOrEmpty(props, "lastName"),
		Label:     utils.GetStringPropOrEmpty(props, "label"),
		Title:     utils.GetStringPropOrEmpty(props, "title"),
		CreatedAt: utils.GetTimePropOrNil(props, "createdAt"),
		Readonly:  utils.GetBoolPropOrFalse(props, "readonly"),
	}
	return &contact
}
