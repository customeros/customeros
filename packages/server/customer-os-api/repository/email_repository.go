package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type EmailRepository interface {
	MergeEmailToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error)
}

type emailRepository struct {
	driver *neo4j.Driver
}

func NewEmailRepository(driver *neo4j.Driver) EmailRepository {
	return &emailRepository{
		driver: driver,
	}
}

func (r *emailRepository) MergeEmailToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[r:EMAILED_AT]->(e:Email {email: $email}) " +
		" ON CREATE SET e.label=$label, r.primary=$primary, e.id=randomUUID(), e:%s " +
		" ON MATCH SET e.label=$label, r.primary=$primary " +
		" RETURN e, r"

	queryResult, err := tx.Run(fmt.Sprintf(query, "Email_"+tenant),
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"email":     entity.Email,
			"label":     entity.Label,
			"primary":   entity.Primary,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}
