package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
	"reflect"
	"time"
)

type ContactService interface {
	Create(ctx context.Context, contact *ContactCreateData) (*entity.ContactEntity, error)

	Update(ctx context.Context, contactUpdateData *ContactUpdateData) (*entity.ContactEntity, error)

	FindContactById(ctx context.Context, id string) (*entity.ContactEntity, error)
	FindContactByEmail(ctx context.Context, email string) (*entity.ContactEntity, error)
	FindContactByPhoneNumber(ctx context.Context, e164 string) (*entity.ContactEntity, error)
	FindAll(ctx context.Context, page, limit int, sortBy []*model.SortBy) (*utils.Pagination, error)

	HardDelete(ctx context.Context, id string) (bool, error)
	SoftDelete(ctx context.Context, id string) (bool, error)

	getDriver() neo4j.Driver
}

type ContactCreateData struct {
	ContactEntity     *entity.ContactEntity
	CustomFields      *entity.CustomFieldEntities
	EmailEntity       *entity.EmailEntity
	PhoneNumberEntity *entity.PhoneNumberEntity
	DefinitionId      *string
	ContactTypeId     *string
	OwnerUserId       *string
}

type ContactUpdateData struct {
	ContactEntity *entity.ContactEntity
	ContactTypeId *string
	OwnerUserId   *string
}

type contactService struct {
	repository *repository.RepositoryContainer
}

func NewContactService(repository *repository.RepositoryContainer) ContactService {
	return &contactService{
		repository: repository,
	}
}

func (s *contactService) getDriver() neo4j.Driver {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *contactService) Create(ctx context.Context, newContact *ContactCreateData) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
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
		contactDbNode, err := s.repository.ContactRepository.Create(tx, tenant, *newContact.ContactEntity)
		if err != nil {
			return nil, err
		}
		var contactId = utils.GetPropsFromNode(*contactDbNode)["id"].(string)

		if newContact.ContactTypeId != nil {
			err := s.repository.ContactRepository.LinkWithContactTypeInTx(tx, common.GetContext(ctx).Tenant, contactId, *newContact.ContactTypeId)
			if err != nil {
				return nil, err
			}
		}

		if newContact.DefinitionId != nil {
			err := s.repository.ContactRepository.LinkWithEntityDefinitionInTx(tx, common.GetContext(ctx).Tenant, contactId, *newContact.DefinitionId)
			if err != nil {
				return nil, err
			}
		}
		if newContact.CustomFields != nil {
			for _, customField := range *newContact.CustomFields {
				queryResult, err := s.repository.CustomFieldRepository.MergeCustomFieldToContactInTx(tx, common.GetContext(ctx).Tenant, contactId, &customField)
				if err != nil {
					return nil, err
				}
				var fieldId = utils.GetPropsFromNode(*queryResult)["id"].(string)
				if customField.DefinitionId != nil {
					err := s.repository.CustomFieldRepository.LinkWithCustomFieldDefinitionForContactInTx(tx, fieldId, contactId, *customField.DefinitionId)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if newContact.EmailEntity != nil {
			err := addEmailToContactInTx(ctx, contactId, *newContact.EmailEntity, tx)
			if err != nil {
				return nil, err
			}
		}
		if newContact.PhoneNumberEntity != nil {
			err := addPhoneNumberToContactInTx(ctx, contactId, *newContact.PhoneNumberEntity, tx)
			if err != nil {
				return nil, err
			}
		}
		if newContact.OwnerUserId != nil {
			err := s.repository.ContactRepository.SetOwner(tx, tenant, contactId, *newContact.OwnerUserId)
			if err != nil {
				return nil, err
			}
		}
		return contactDbNode, nil
	}
}

func (s *contactService) Update(ctx context.Context, contactUpdateData *ContactUpdateData) (*entity.ContactEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	contactDbNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		contactId := contactUpdateData.ContactEntity.Id
		queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			SET c.firstName=$firstName,
				c.lastName=$lastName,
				c.label=$label,
				c.title=$title,
				c.notes=$notes
			RETURN c`,
			map[string]interface{}{
				"tenant":    tenant,
				"contactId": contactId,
				"firstName": contactUpdateData.ContactEntity.FirstName,
				"lastName":  contactUpdateData.ContactEntity.LastName,
				"label":     contactUpdateData.ContactEntity.Label,
				"title":     contactUpdateData.ContactEntity.Title,
				"notes":     contactUpdateData.ContactEntity.Notes,
			})

		err = s.repository.ContactRepository.UnlinkFromContactTypesInTx(tx, tenant, contactId)
		if err != nil {
			return nil, err
		}
		if contactUpdateData.ContactTypeId != nil {
			err := s.repository.ContactRepository.LinkWithContactTypeInTx(tx, tenant, contactId, *contactUpdateData.ContactTypeId)
			if err != nil {
				return nil, err
			}
		}

		err = s.repository.ContactRepository.RemoveOwner(tx, tenant, contactId)
		if err != nil {
			return nil, err
		}
		if contactUpdateData.OwnerUserId != nil {
			err := s.repository.ContactRepository.SetOwner(tx, tenant, contactId, *contactUpdateData.OwnerUserId)
			if err != nil {
				return nil, err
			}
		}

		return utils.ExtractSingleRecordFirstValueAsNodePtr(queryResult, err)
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactEntity(contactDbNode.(*dbtype.Node)), nil
}

func (s *contactService) HardDelete(ctx context.Context, contactId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (c)-[:HAS_PROPERTY]->(f:CustomField)
			OPTIONAL MATCH (c)-[:CALLED_AT]->(p:PhoneNumber)
			OPTIONAL MATCH (c)-[:EMAILED_AT]->(e:Email)
            DETACH DELETE p, e, f, c
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

func (s *contactService) SoftDelete(ctx context.Context, contactId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
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
	session := utils.NewNeo4jReadSession(s.getDriver())
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
	session := utils.NewNeo4jReadSession(s.getDriver())
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
	session := utils.NewNeo4jReadSession(s.getDriver())
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

func (s *contactService) FindAll(ctx context.Context, page, limit int, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	sortings, err := s.prepareContactsSorting(sortBy)
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repository.ContactRepository.GetPaginatedContacts(
		session,
		common.GetContext(ctx).Tenant,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		sortings)
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

func (s *contactService) prepareContactsSorting(sortBy []*model.SortBy) (*utils.Sorts, error) {
	transformedSorting := new(utils.Sorts)
	if sortBy != nil {
		for _, v := range sortBy {
			err := transformedSorting.NewSortRule(v.By, v.Direction.String(), *v.CaseSensitive, reflect.TypeOf(entity.ContactEntity{}))
			if err != nil {
				return nil, err
			}
		}
	}
	return transformedSorting, nil
}

func (s *contactService) mapDbNodeToContactEntity(dbContactNode *dbtype.Node) *entity.ContactEntity {
	props := utils.GetPropsFromNode(*dbContactNode)
	contact := entity.ContactEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		FirstName: utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:  utils.GetStringPropOrEmpty(props, "lastName"),
		Label:     utils.GetStringPropOrEmpty(props, "label"),
		Title:     utils.GetStringPropOrEmpty(props, "title"),
		Notes:     utils.GetStringPropOrEmpty(props, "notes"),
		CreatedAt: props["createdAt"].(time.Time),
	}
	return &contact
}
