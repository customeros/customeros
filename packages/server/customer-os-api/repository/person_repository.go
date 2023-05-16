package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type PersonRepository interface {
	Merge(ctx context.Context, tx neo4j.ManagedTransaction, entity *entity.PersonEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tx neo4j.ManagedTransaction, entity *entity.PersonEntity) (*dbtype.Node, error)
	SetDefaultUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, personId, userId string, relation entity.PersonRelation) (*dbtype.Node, error)
	LinkWithUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, personId, userId, userTenant string, relation entity.PersonRelation) error
	UnlinkUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, personId, userId, userTenant string, relation entity.PersonRelation) error
	GetUsersForPerson(ctx context.Context, ids []string) ([]*utils.DbNodeWithRelationIdAndTenant, error)
	GetPersonByEmailProvider(ctx context.Context, email string, provider string) (*dbtype.Node, error)
	GetPersonByIdentityId(ctx context.Context, identityId string) (*dbtype.Node, error)
	GetPersonForUser(ctx context.Context, tenant string, userId string, relation entity.PersonRelation) (*dbtype.Node, error)
}

type personRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPersonRepository(driver *neo4j.DriverWithContext) PersonRepository {
	return &personRepository{
		driver: driver,
	}
}

func (r *personRepository) Merge(ctx context.Context, tx neo4j.ManagedTransaction, entity *entity.PersonEntity) (*dbtype.Node, error) {
	query := "MERGE (p:Person {email:$email, provider:$provider}) " +
		" ON CREATE SET p.id=RandomUUID(), " +
		"				p.identityId=$identityId, " +
		"				p.createdAt=$createdAt, " +
		"				p.updatedAt=$updatedAt, " +
		"				p.appSource=$appSource, " +
		"				p.source=$source, " +
		"				p.sourceOfTruth=$sourceOfTruth " +
		" RETURN p"

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"email":         entity.Email,
			"provider":      entity.Provider,
			"identityId":    entity.IdentityId,
			"createdAt":     entity.CreatedAt,
			"updatedAt":     entity.CreatedAt,
			"appSource":     entity.AppSource,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *personRepository) Update(ctx context.Context, tx neo4j.ManagedTransaction, entity *entity.PersonEntity) (*dbtype.Node, error) {
	query := "MATCH (p:Person {id: $id}) " +
		"  SET p.identityId=$identityId, " +
		"				p.updatedAt=$updatedAt, " +
		"				p.sourceOfTruth=$sourceOfTruth, " +
		"				p.appSource=$appSource " +
		" RETURN p"

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"id":            entity.Id,
			"identityId":    entity.IdentityId,
			"updatedAt":     entity.UpdatedAt,
			"sourceOfTruth": entity.SourceOfTruth,
			"appSource":     entity.AppSource,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *personRepository) SetDefaultUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, personId, userId string, relation entity.PersonRelation) (*dbtype.Node, error) {
	queryResult, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (p:Person {id:$personId})-[r:%s]->(u:User)
			SET r.default=
				CASE u.id
					WHEN $userId THEN true
					ELSE false
				END
			RETURN DISTINCT(p)`, relation),
		map[string]interface{}{
			"personId": personId,
			"userId":   userId,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *personRepository) LinkWithUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, personId, userId, userTenant string, relation entity.PersonRelation) error {
	queryResult, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (p:Person {id:$personId}), (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$userTenant})
			MERGE (p)-[r:%s]->(u)
			SET r.default= CASE
				WHEN NOT EXISTS((p)-[:%s {default: true}]->(:User)) THEN true
				ELSE false
			END
			RETURN p`, relation, relation),
		map[string]interface{}{
			"personId":   personId,
			"userId":     userId,
			"userTenant": userTenant,
		})
	if err != nil {
		return fmt.Errorf("error linking person with user: %w", err)
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *personRepository) UnlinkUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, personId, userId, userTenant string, relation entity.PersonRelation) error {
	query := fmt.Sprintf(`
							MATCH (p:Person {id:$personId}), (u:User_%s {id:$userId})
							MATCH (p)-[r:%s]->(u)
							DELETE r return p`, userTenant, relation)

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"personId": personId,
			"userId":   userId,
		})
	if err != nil {
		return fmt.Errorf("Error unlinking person with user: %w", err)
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *personRepository) GetUsersForPerson(ctx context.Context, ids []string) ([]*utils.DbNodeWithRelationIdAndTenant, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (p:Person)-[rel:%s]->(u:User)-[:USER_BELONGS_TO_TENANT]->(t:Tenant) WHERE p.id IN $ids RETURN u, rel, p.id, t.name`, entity.IDENTIFIES)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]any{
				"ids": ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationIdAndTenant(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting users for person: %w", err)
	}

	return result.([]*utils.DbNodeWithRelationIdAndTenant), nil
}

func (r *personRepository) GetPersonByEmailProvider(ctx context.Context, email string, provider string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (p:Person {email:$email, provider:$provider}) RETURN p`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]any{
				"email":    email,
				"provider": provider,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting person by email and provider: %w", err)
	}

	return result.(*dbtype.Node), nil

}

func (r *personRepository) GetPersonForUser(ctx context.Context, tenant string, userId string, relation entity.PersonRelation) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (p:Person)-[:%s]->(u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) RETURN p`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, relation),
			map[string]any{
				"userId": userId,
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting person for user: %w", err)
	}

	return result.(*dbtype.Node), nil

}

func (r *personRepository) GetPersonByIdentityId(ctx context.Context, identityId string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (p:Person {identityId:$identityId}) RETURN p`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]any{
				"identityId": identityId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting person by identityId: %w", err)
	}

	return result.(*dbtype.Node), nil

}
