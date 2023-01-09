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
	UpdateEmailByContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error)
	FindAllForContact(session neo4j.Session, tenant, contactId string) (any, error)
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
		" ON CREATE SET e.label=$label, r.primary=$primary, e.id=randomUUID(), e.source=$source, e.sourceOfTruth=$sourceOfTruth, e:%s " +
		" ON MATCH SET e.label=$label, r.primary=$primary, e.sourceOfTruth=$sourceOfTruth " +
		" RETURN e, r"

	queryResult, err := tx.Run(fmt.Sprintf(query, "Email_"+tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"email":         entity.Email,
			"label":         entity.Label,
			"primary":       entity.Primary,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}

func (r *emailRepository) UpdateEmailByContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	queryResult, err := tx.Run(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				(c)-[r:EMAILED_AT]->(e:Email {id:$emailId}) 
			SET e.email=$email,
				e.label=$label,
				r.primary=$primary,
				e.sourceOfTruth=$sourceOfTruth,
			RETURN e, r`,
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"emailId":       entity.Id,
			"email":         entity.Email,
			"label":         entity.Label,
			"primary":       entity.Primary,
			"sourceOfTruth": entity.SourceOfTruth,
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}

func (r *emailRepository) FindAllForContact(session neo4j.Session, tenant, contactId string) (any, error) {
	return session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[r:EMAILED_AT]->(e:Email) 				
				RETURN e, r`,
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
