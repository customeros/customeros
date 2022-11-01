package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
	"time"
)

type ContactService interface {
	Create(ctx context.Context, contact *ContactCreateData) (*entity.ContactEntity, error)
	FindContactById(ctx context.Context, id string) (*entity.ContactEntity, error)
	FindAll(ctx context.Context, page int, limit int) (*utils.Pagination, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type ContactCreateData struct {
	ContactEntity     *entity.ContactEntity
	TextCustomFields  *entity.TextCustomFieldEntities
	EmailEntity       *entity.EmailEntity
	PhoneNumberEntity *entity.PhoneNumberEntity
}

type contactService struct {
	driver *neo4j.Driver
}

func NewContactService(driver *neo4j.Driver) ContactService {
	return &contactService{
		driver: driver,
	}
}

func (s *contactService) Create(ctx context.Context, newContact *ContactCreateData) (*entity.ContactEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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
				  company: $company,
				  companyTitle: $companyTitle,
				  notes: $notes,
				  contactType: $contactType,
                  createdAt :datetime({timezone: 'UTC'})
			})-[:CONTACT_BELONGS_TO_TENANT]->(t)
			RETURN c`,
			map[string]interface{}{
				"tenant":       common.GetContext(ctx).Tenant,
				"firstName":    newContact.ContactEntity.FirstName,
				"lastName":     newContact.ContactEntity.LastName,
				"label":        newContact.ContactEntity.Label,
				"contactType":  newContact.ContactEntity.ContactType,
				"company":      newContact.ContactEntity.Company,
				"title":        newContact.ContactEntity.Title,
				"companyTitle": newContact.ContactEntity.CompanyTitle,
				"notes":        newContact.ContactEntity.Notes,
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

func (s *contactService) Delete(ctx context.Context, contactId string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (c)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField)
			OPTIONAL MATCH (c)-[:CALLED_AT]->(p:PhoneNumber)
			OPTIONAL MATCH (c)-[:EMAILED_AT]->(e:Email)
            DETACH DELETE p, e, f, c
			`,
			map[string]interface{}{
				"id":     contactId,
				"tenant": common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *contactService) FindContactById(ctx context.Context, id string) (*entity.ContactEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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

func (s *contactService) FindAll(ctx context.Context, page int, limit int) (*utils.Pagination, error) {
	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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

	contacts := entity.ContactNodes{}

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
		Id:           utils.GetStringPropOrEmpty(props, "id"),
		FirstName:    utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:     utils.GetStringPropOrEmpty(props, "lastName"),
		Label:        utils.GetStringPropOrEmpty(props, "label"),
		Company:      utils.GetStringPropOrEmpty(props, "company"),
		Title:        utils.GetStringPropOrEmpty(props, "title"),
		CompanyTitle: utils.GetStringPropOrEmpty(props, "companyTitle"),
		Notes:        utils.GetStringPropOrEmpty(props, "notes"),
		ContactType:  utils.GetStringPropOrEmpty(props, "contactType"),
		CreatedAt:    props["createdAt"].(time.Time),
	}
	return &contact
}
