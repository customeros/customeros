package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type EntityType string

const (
	CONTACT      EntityType = "CONTACT"
	USER         EntityType = "USER"
	ORGANIZATION EntityType = "ORGANIZATION"
)

type EmailRepository interface {
	MergeEmailToInTx(tx neo4j.Transaction, tenant string, entityType EntityType, entityId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error)
	UpdateEmailByInTx(tx neo4j.Transaction, tenant string, entityType EntityType, entityId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error)
	SetOtherEmailsNonPrimaryInTx(tx neo4j.Transaction, tenantId string, entityType EntityType, entityId string, email string) error
	FindAllFor(session neo4j.Session, tenant string, entityType EntityType, entityId string) (any, error)
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
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	}

	query = query +
		" MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {email: $email}) " +
		" ON CREATE SET e.id=randomUUID(), " +
		"				e.source=$source, " +
		"				e.sourceOfTruth=$sourceOfTruth, " +
		" 				e.appSource=$appSource, " +
		"				e.createdAt=$now, " +
		"				e.updatedAt=$now, " +
		"				e:%s " +
		" WITH e, entity " +
		" MERGE (entity)-[rel:HAS]->(e) " +
		" SET 	rel.label=$label, " +
		"		rel.primary=$primary, " +
		"		e.sourceOfTruth=$sourceOfTruth," +
		"		e.updatedAt=$now " +
		" RETURN e, rel"

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
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}

func (r *emailRepository) UpdateEmailByInTx(tx neo4j.Transaction, tenant string, entityType EntityType, entityId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	query := ""

	switch entityType {
	case CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	}

	queryResult, err := tx.Run(query+`, (entity)-[rel:HAS]->(e:Email {id:$emailId}) 
			SET rel.label=$label,
				rel.primary=$primary,
				e.sourceOfTruth=$sourceOfTruth,
				e.updatedAt=$now
			RETURN e, rel`,
		map[string]interface{}{
			"tenant":        tenant,
			"entityId":      entityId,
			"emailId":       entity.Id,
			"label":         entity.Label,
			"primary":       entity.Primary,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}

func (r *emailRepository) FindAllFor(session neo4j.Session, tenant string, entityType EntityType, entityId string) (any, error) {
	return session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		query := ""

		switch entityType {
		case CONTACT:
			query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
		case USER:
			query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
		case ORGANIZATION:
			query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
		}

		result, err := tx.Run(query+`, (entity)-[rel:HAS]->(e:Email) 				
				RETURN e, rel`,
			map[string]interface{}{
				"entityId": entityId,
				"tenant":   tenant,
			})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
}

func (r *emailRepository) SetOtherEmailsNonPrimaryInTx(tx neo4j.Transaction, tenantId string, entityType EntityType, entityId string, email string) error {
	query := ""

	switch entityType {
	case CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}

	_, err := tx.Run(query+`, (entity)-[rel:HAS]->(e:Email)
			WHERE e.email <> $email
            SET rel.primary=false, 
				e.updatedAt=$now`,
		map[string]interface{}{
			"tenant":   tenantId,
			"entityId": entityId,
			"email":    email,
			"now":      utils.Now(),
		})
	return err
}
