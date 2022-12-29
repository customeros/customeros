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
		" MERGE (c)-[r:CALLED_AT]->(p:PhoneNumber {e164: $e164}) " +
		" ON CREATE SET p.label=$label, r.primary=$primary, p.id=randomUUID(), p:%s " +
		" ON MATCH SET p.label=$label, r.primary=$primary " +
		" RETURN p, r"

	queryResult, err := tx.Run(fmt.Sprintf(query, "PhoneNumber_"+tenant),
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"e164":      entity.E164,
			"label":     entity.Label,
			"primary":   entity.Primary,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}
