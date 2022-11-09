package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
	"time"
)

type ContactService interface {
	Create(ctx context.Context, contact *ContactCreateData) (*entity.ContactEntity, error)

	Update(ctx context.Context, contact *entity.ContactEntity) (*entity.ContactEntity, error)

	FindContactById(ctx context.Context, id string) (*entity.ContactEntity, error)
	FindContactByEmail(ctx context.Context, email string) (*entity.ContactEntity, error)
	FindContactByPhoneNumber(ctx context.Context, e164 string) (*entity.ContactEntity, error)
	FindAll(ctx context.Context, page int, limit int) (*utils.Pagination, error)

	HardDelete(ctx context.Context, id string) (bool, error)
	SoftDelete(ctx context.Context, id string) (bool, error)
}

type ContactCreateData struct {
	ContactEntity     *entity.ContactEntity
	TextCustomFields  *entity.TextCustomFieldEntities
	EmailEntity       *entity.EmailEntity
	PhoneNumberEntity *entity.PhoneNumberEntity
	CompanyPosition   *entity.CompanyPositionEntity
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
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(createContactInDBTxWork(ctx, newContact))
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(queryResult.(dbtype.Node)), nil
}

func createContactInDBTxWork(ctx context.Context, newContact *ContactCreateData) func(tx neo4j.Transaction) (interface{}, error) {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			CREATE (c:Contact {
				  id: randomUUID(),
				  title: $title,
				  firstName: $firstName,
				  lastName: $lastName,
				  label: $label,
				  notes: $notes,
				  contactType: $contactType,
                  createdAt :datetime({timezone: 'UTC'})
			})-[:CONTACT_BELONGS_TO_TENANT]->(t)
			RETURN c`,
			map[string]interface{}{
				"tenant":      common.GetContext(ctx).Tenant,
				"firstName":   newContact.ContactEntity.FirstName,
				"lastName":    newContact.ContactEntity.LastName,
				"label":       newContact.ContactEntity.Label,
				"contactType": newContact.ContactEntity.ContactType,
				"title":       newContact.ContactEntity.Title,
				"notes":       newContact.ContactEntity.Notes,
			})

		record, err := result.Single()
		if err != nil {
			return nil, err
		}

		var contactId = utils.GetPropsFromNode(record.Values[0].(dbtype.Node))["id"].(string)
		if newContact.TextCustomFields != nil {
			for _, textCustomField := range *newContact.TextCustomFields {
				err := addTextCustomFieldToContactInTx(ctx, contactId, textCustomField, tx)
				if err != nil {
					return nil, err
				}
			}
		}
		if newContact.CompanyPosition != nil {
			err := setCompanyPositionRelationshipInTx(ctx, contactId, *newContact.CompanyPosition, tx)
			if err != nil {
				return nil, err
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
		return record.Values[0], nil
	}
}

func (s *contactService) Update(ctx context.Context, contact *entity.ContactEntity) (*entity.ContactEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			SET c.firstName=$firstName,
				c.lastName=$lastName,
				c.label=$label,
				c.contactType=$contactType,
				c.title=$title,
				c.notes=$notes
			RETURN c`,
			map[string]interface{}{
				"tenant":      common.GetContext(ctx).Tenant,
				"contactId":   contact.Id,
				"firstName":   contact.FirstName,
				"lastName":    contact.LastName,
				"label":       contact.Label,
				"contactType": contact.ContactType,
				"title":       contact.Title,
				"notes":       contact.Notes,
			})
		record, err := txResult.Single()
		if err != nil {
			return nil, err
		}
		return record, nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactEntity(queryResult.(*db.Record).Values[0].(dbtype.Node)), nil
}

func (s *contactService) HardDelete(ctx context.Context, contactId string) (bool, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (c)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField)
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
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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

	return s.mapDbNodeToContactEntity(queryResult.(dbtype.Node)), nil
}

func (s *contactService) FindContactByEmail(ctx context.Context, email string) (*entity.ContactEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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

	return s.mapDbNodeToContactEntity(queryResult.(dbtype.Node)), nil
}

func (s *contactService) FindContactByPhoneNumber(ctx context.Context, e164 string) (*entity.ContactEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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

	return s.mapDbNodeToContactEntity(queryResult.(dbtype.Node)), nil
}

func (s *contactService) FindAll(ctx context.Context, page int, limit int) (*utils.Pagination, error) {
	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	dataResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
				MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) RETURN count(c) as count`,
			map[string]interface{}{
				"tenant": common.GetContext(ctx).Tenant,
			})
		count, _ := result.Single()
		paginatedResult.SetTotalRows(count.Values[0].(int64))

		result, err = tx.Run(`
				MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) RETURN c SKIP $skip LIMIT $limit`,
			map[string]interface{}{
				"tenant": common.GetContext(ctx).Tenant,
				"skip":   paginatedResult.GetSkip(),
				"limit":  paginatedResult.GetLimit(),
			})
		data, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return data, nil
	})
	if err != nil {
		return nil, err
	}

	contacts := entity.ContactEntities{}

	for _, dbRecord := range dataResult.([]*db.Record) {
		contact := s.mapDbNodeToContactEntity(dbRecord.Values[0].(dbtype.Node))
		contacts = append(contacts, *contact)
	}
	paginatedResult.SetRows(&contacts)
	return &paginatedResult, nil
}

func (s *contactService) mapDbNodeToContactEntity(dbContactNode dbtype.Node) *entity.ContactEntity {
	props := utils.GetPropsFromNode(dbContactNode)
	contact := entity.ContactEntity{
		Id:          utils.GetStringPropOrEmpty(props, "id"),
		FirstName:   utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:    utils.GetStringPropOrEmpty(props, "lastName"),
		Label:       utils.GetStringPropOrEmpty(props, "label"),
		Title:       utils.GetStringPropOrEmpty(props, "title"),
		Notes:       utils.GetStringPropOrEmpty(props, "notes"),
		ContactType: utils.GetStringPropOrEmpty(props, "contactType"),
		CreatedAt:   props["createdAt"].(time.Time),
	}
	return &contact
}
