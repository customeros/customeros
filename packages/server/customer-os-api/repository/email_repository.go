package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"time"
)

type EntityType string

const (
	CONTACT      EntityType = "CONTACT"
	USER         EntityType = "USER"
	ORGANIZATION EntityType = "ORGANIZATION"
)

type EmailRepository interface {
	MergeEmailToInTx(tx neo4j.Transaction, tenant string, entityType EntityType, entityId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error)
	UpdateEmailByContactInTx(tx neo4j.Transaction, tenant, contactId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error)
	SetOtherContactEmailsNonPrimaryInTx(tx neo4j.Transaction, tenantId, contactId, email string) error
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

func (r *emailRepository) MergeEmailToInTx(tx neo4j.Transaction, tenant string, entityType EntityType, entityId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	query := ""

	switch entityType {
	case CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}

	query = query +
		" MERGE (entity)-[r:EMAILED_AT]->(e:Email {email: $email}) " +
		" ON CREATE SET e.label=$label, " +
		"				r.primary=$primary, " +
		"				e.id=randomUUID(), " +
		"				e.source=$source, " +
		"				e.sourceOfTruth=$sourceOfTruth, " +
		" 				e.appSource=$appSource, " +
		"				e.createdAt=$now, " +
		"				e.updatedAt=$now, " +
		"				e:%s " +
		" ON MATCH SET 	e.label=$label, " +
		"				r.primary=$primary, " +
		"				e.sourceOfTruth=$sourceOfTruth, " +
		"				e.updatedAt=$now " +
		" RETURN e, r"

	queryResult, err := tx.Run(fmt.Sprintf(query, "Email_"+tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"entityId":      entityId,
			"email":         entity.Email,
			"label":         entity.Label,
			"primary":       entity.Primary,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
			"appSource":     entity.AppSource,
			"now":           time.Now().UTC(),
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
				e.updatedAt=$now
			RETURN e, r`,
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"emailId":       entity.Id,
			"email":         entity.Email,
			"label":         entity.Label,
			"primary":       entity.Primary,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           time.Now().UTC(),
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

func (r *emailRepository) SetOtherContactEmailsNonPrimaryInTx(tx neo4j.Transaction, tenantId, contactId, email string) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				 (c)-[r:EMAILED_AT]->(e:Email)
			WHERE e.email <> $email
            SET r.primary=false, 
				e.updatedAt=$now`,
		map[string]interface{}{
			"tenant":    tenantId,
			"contactId": contactId,
			"email":     email,
			"now":       time.Now().UTC(),
		})
	return err
}
