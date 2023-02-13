package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type EmailRepository interface {
	MergeEmailToInTx(tx neo4j.Transaction, tenant string, entityType entity.EntityType, entityId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error)
	UpdateEmailByInTx(tx neo4j.Transaction, tenant string, entityType entity.EntityType, entityId string, entity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error)
	SetOtherEmailsNonPrimaryInTx(tx neo4j.Transaction, tenantId string, entityType entity.EntityType, entityId string, email string) error
	FindAllFor(session neo4j.Session, tenant string, entityType entity.EntityType, entityId string) (any, error)
	RemoveRelationship(entityType entity.EntityType, tenant, entityId, email string) error
	RemoveRelationshipById(entityType entity.EntityType, tenant, entityId, emailId string) error
	DeleteById(tenant, emailId string) error
}

type emailRepository struct {
	driver *neo4j.Driver
}

func NewEmailRepository(driver *neo4j.Driver) EmailRepository {
	return &emailRepository{
		driver: driver,
	}
}

func (r *emailRepository) MergeEmailToInTx(tx neo4j.Transaction, tenant string, entityType entity.EntityType, entityId string, emailEntity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	query := ""

	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
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
			"email":         emailEntity.Email,
			"label":         emailEntity.Label,
			"primary":       emailEntity.Primary,
			"source":        emailEntity.Source,
			"sourceOfTruth": emailEntity.SourceOfTruth,
			"appSource":     emailEntity.AppSource,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}

func (r *emailRepository) UpdateEmailByInTx(tx neo4j.Transaction, tenant string, entityType entity.EntityType, entityId string, emailEntity entity.EmailEntity) (*dbtype.Node, *dbtype.Relationship, error) {
	query := ""

	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
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
			"emailId":       emailEntity.Id,
			"label":         emailEntity.Label,
			"primary":       emailEntity.Primary,
			"sourceOfTruth": emailEntity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordNodeAndRelationship(queryResult, err)
}

func (r *emailRepository) FindAllFor(session neo4j.Session, tenant string, entityType entity.EntityType, entityId string) (any, error) {
	return session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		query := ""

		switch entityType {
		case entity.CONTACT:
			query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
		case entity.USER:
			query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
		case entity.ORGANIZATION:
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

func (r *emailRepository) SetOtherEmailsNonPrimaryInTx(tx neo4j.Transaction, tenantId string, entityType entity.EntityType, entityId string, email string) error {
	query := ""

	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
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

func (r *emailRepository) RemoveRelationship(entityType entity.EntityType, tenant, entityId, email string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := ""
	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(query+`MATCH (entity)-[rel:HAS]->(e:Email {email:$email})
            DELETE rel`,
			map[string]any{
				"entityId": entityId,
				"email":    email,
				"tenant":   tenant,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *emailRepository) RemoveRelationshipById(entityType entity.EntityType, tenant, entityId, emailId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := ""
	switch entityType {
	case entity.CONTACT:
		query = `MATCH (entity:Contact {id:$entityId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.USER:
		query = `MATCH (entity:User {id:$entityId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	case entity.ORGANIZATION:
		query = `MATCH (entity:Organization {id:$entityId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) `
	}

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(query+`MATCH (entity)-[rel:HAS]->(e:Email {id:$emailId})
            DELETE rel`,
			map[string]any{
				"entityId": entityId,
				"emailId":  emailId,
				"tenant":   tenant,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *emailRepository) DeleteById(tenant, emailId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`MATCH (e:Email {id:$emailId})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
            DETACH DELETE e`,
			map[string]any{
				"tenant":  tenant,
				"emailId": emailId,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}
