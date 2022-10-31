package service

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ContactService interface {
	Create(ctx context.Context, contact *entity.ContactNode) (*entity.ContactNode, error)
	FindContactById(ctx context.Context, id string) (*entity.ContactNode, error)
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

func (s *contactService) Create(ctx context.Context, newContact *entity.ContactNode) (*entity.ContactNode, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
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
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}
	contact := entity.ContactNode{}
	mapstructure.Decode(utils.GetPropsFromNode(queryResult.(dbtype.Node)), &contact)

	return &contact, nil
}

func (s *contactService) FindContactById(ctx context.Context, id string) (*entity.ContactNode, error) {
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

	contact := entity.ContactNode{}
	err = mapstructure.Decode(utils.GetPropsFromNode(queryResult.(dbtype.Node)), &contact)
	if err != nil {
		return nil, err
	}

	return &contact, nil
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
		contact := entity.ContactNode{}
		mapstructure.Decode(utils.GetPropsFromNode(dbRecord.Values[0].(dbtype.Node)), &contact)
		contacts = append(contacts, contact)
	}

	return &contacts, nil
}
