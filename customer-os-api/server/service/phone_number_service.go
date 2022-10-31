package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type PhoneNumberService interface {
	FindAllForContact(ctx context.Context, obj *model.Contact) (*entity.PhoneNumberEntities, error)
	MergePhoneNumberToContact(ctx context.Context, id string, toEntity *entity.PhoneNumberEntity) (*entity.PhoneNumberEntity, error)
	Delete(ctx context.Context, contactId string, phoneNumber string) (bool, error)
}

type phoneNumberService struct {
	driver *neo4j.Driver
}

func NewPhoneNumberService(driver *neo4j.Driver) PhoneNumberService {
	return &phoneNumberService{
		driver: driver,
	}
}

func (s *phoneNumberService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.PhoneNumberEntities, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
				MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c:Contact {id:$id})-[:CALLED_AT]->(p:PhoneNumber) 
				RETURN p `,
			map[string]interface{}{
				"id":     contact.ID,
				"tenant": common.GetContext(ctx).Tenant})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	phoneNumberEntities := entity.PhoneNumberEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		phoneNumberEntity := s.mapDbNodeToPhoneNumberEntity(dbRecord.Values[0].(dbtype.Node))
		phoneNumberEntities = append(phoneNumberEntities, *phoneNumberEntity)
	}

	return &phoneNumberEntities, nil
}

func (s *phoneNumberService) MergePhoneNumberToContact(ctx context.Context, contactId string, entity *entity.PhoneNumberEntity) (*entity.PhoneNumberEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (p:PhoneNumber {number: $number})<-[:CALLED_AT]-(c)
            ON CREATE SET p.label=$label
            ON MATCH SET p.label=$label
			RETURN p`,
			map[string]interface{}{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"number":    entity.Number,
				"label":     entity.Label,
			})
		record, err := txResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToPhoneNumberEntity(queryResult.(dbtype.Node)), nil
}

func addPhoneNumberToContactInTx(ctx context.Context, contactId string, input entity.PhoneNumberEntity, tx neo4j.Transaction) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})
			CREATE (p:PhoneNumber {
				  number: $number,
				  label: $label
			})<-[:CALLED_AT]-(c)
			RETURN p`,
		map[string]interface{}{
			"contactId": contactId,
			"number":    input.Number,
			"label":     input.Label,
		})

	return err
}

func (s *phoneNumberService) Delete(ctx context.Context, contactId string, phoneNumber string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c:Contact {id:$id})-[:CALLED_AT]->(p:PhoneNumber {number:$number})
            DETACH DELETE p
			`,
			map[string]interface{}{
				"id":     contactId,
				"number": phoneNumber,
				"tenant": common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *phoneNumberService) mapDbNodeToPhoneNumberEntity(dbContactGroupNode dbtype.Node) *entity.PhoneNumberEntity {
	props := utils.GetPropsFromNode(dbContactGroupNode)
	result := entity.PhoneNumberEntity{
		Number: utils.GetStringPropOrEmpty(props, "number"),
		Label:  utils.GetStringPropOrEmpty(props, "label"),
	}
	return &result
}
