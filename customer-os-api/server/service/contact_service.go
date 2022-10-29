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
	Create(ctx context.Context, contact *entity.ContactEntity) (*entity.ContactEntity, error)
	FindContactById(ctx context.Context, id string) (*entity.ContactEntity, error)
	FindAll(ctx context.Context) (*entity.ContactNodes, error)
}

type contactService struct {
	driver *neo4j.Driver
}

func NewContactService(driver *neo4j.Driver) ContactService {
	return &contactService{
		driver: driver,
	}
}

func (s *contactService) Create(ctx context.Context, newContact *entity.ContactEntity) (*entity.ContactEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(createContactInDBTxWork(ctx, newContact))
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactEntity(queryResult.(dbtype.Node)), nil
}

func createContactInDBTxWork(ctx context.Context, newContact *entity.ContactEntity) func(tx neo4j.Transaction) (interface{}, error) {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			CREATE (c:Contact {
				  id: randomUUID(),
				  firstName: $firstName,
				  lastName: $lastName,
				  label: $label,
				  companyName: $companyName,
				  contactType: $contactType,
                  createdAt :datetime({timezone: 'UTC'})
			})-[:BELONGS_TO]->(t)
			RETURN c`,
			map[string]interface{}{
				"tenant":      common.GetContext(ctx).Tenant,
				"firstName":   newContact.FirstName,
				"lastName":    newContact.LastName,
				"label":       newContact.Label,
				"contactType": newContact.ContactType,
				"companyName": newContact.CompanyName,
			})

		record, err := result.Single()
		if err != nil {
			return nil, err
		}

		for _, textCustomField := range newContact.TextCustomFields {
			_, err := addTextCustomFieldToContact(ctx, utils.GetPropsFromNode(record.Values[0].(dbtype.Node))["id"].(string), textCustomField, tx)
			if err != nil {
				return nil, err
			}
		}

		return record.Values[0], nil
	}
}

func (s *contactService) FindContactById(ctx context.Context, id string) (*entity.ContactEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (c:Contact {id:$id})--(:Tenant {name:$tenant}) RETURN c`,
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

func (s *contactService) FindAll(ctx context.Context) (*entity.ContactNodes, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`MATCH (c:Contact)--(:Tenant {name:$tenant}) RETURN c`, map[string]interface{}{
			"tenant": common.GetContext(ctx).Tenant,
		})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	contacts := entity.ContactNodes{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		contact := s.mapDbNodeToContactEntity(dbRecord.Values[0].(dbtype.Node))
		contacts = append(contacts, *contact)
	}

	return &contacts, nil
}

func (s *contactService) mapDbNodeToContactEntity(dbContactNode dbtype.Node) *entity.ContactEntity {
	props := utils.GetPropsFromNode(dbContactNode)
	contact := entity.ContactEntity{
		Id:          props["id"].(string),
		FirstName:   props["firstName"].(string),
		LastName:    props["lastName"].(string),
		Label:       props["label"].(string),
		CompanyName: props["companyName"].(string),
		ContactType: props["contactType"].(string),
		CreatedAt:   props["createdAt"].(time.Time),
	}
	return &contact
}
