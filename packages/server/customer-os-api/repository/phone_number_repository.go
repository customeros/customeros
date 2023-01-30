package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type PhoneNumberRepository interface {
	MergePhoneNumberToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error)
	UpdatePhoneNumberByContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error)
	FindAllForContact(session neo4j.Session, tenant, contactId string) (any, error)
	SetOtherContactPhoneNumbersNonPrimaryInTx(tx neo4j.Transaction, tenant, contactId, e164 string) error
}

type phoneNumberRepository struct {
	driver *neo4j.Driver
}

func NewPhoneNumberRepository(driver *neo4j.Driver) PhoneNumberRepository {
	return &phoneNumberRepository{
		driver: driver,
	}
}

func (r *phoneNumberRepository) MergePhoneNumberToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[r:PHONE_ASSOCIATED_WITH]->(p:PhoneNumber {e164: $e164}) " +
		" ON CREATE SET p.label=$label, " +
		"				r.primary=$primary, " +
		"				p.id=randomUUID(), " +
		"				p.source=$source, " +
		"				p.sourceOfTruth=$sourceOfTruth, " +
		"				p.createdAt=datetime({timezone: 'UTC'}), " +
		"				p:%s " +
		" ON MATCH SET 	p.label=$label, " +
		"				r.primary=$primary, " +
		"				p.sourceOfTruth=$sourceOfTruth " +
		" RETURN p, r"

	queryResult, err := tx.Run(fmt.Sprintf(query, "PhoneNumber_"+tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"e164":          entity.E164,
			"label":         entity.Label,
			"primary":       entity.Primary,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}

func (r *phoneNumberRepository) UpdatePhoneNumberByContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.PhoneNumberEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				(c)-[r:PHONE_ASSOCIATED_WITH]->(p:PhoneNumber {id: $phoneId})
            SET p.e164=$e164,
				p.label=$label,
				r.primary=$primary,
				p.sourceOfTruth=$sourceOfTruth,
			RETURN p, r`,
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"emailId":       entity.Id,
			"email":         entity.E164,
			"label":         entity.Label,
			"primary":       entity.Primary,
			"sourceOfTruth": entity.SourceOfTruth,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}

func (r *phoneNumberRepository) FindAllForContact(session neo4j.Session, tenant, contactId string) (any, error) {
	return session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[r:PHONE_ASSOCIATED_WITH]->(p:PhoneNumber)
				RETURN p, r`,
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
			})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
}

func (r *phoneNumberRepository) SetOtherContactPhoneNumbersNonPrimaryInTx(tx neo4j.Transaction, tenant, contactId, e164 string) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				 (c)-[r:PHONE_ASSOCIATED_WITH]->(p:PhoneNumber)
			WHERE p.e164 <> $e164
            SET r.primary=false`,
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"e164":      e164,
		})
	return err
}
