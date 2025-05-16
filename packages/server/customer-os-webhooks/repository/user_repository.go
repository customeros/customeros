package repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
)

type UserRepository interface {
	// Deprecated
	GetMatchedUserId(ctx context.Context, tenant, externalSystem, externalId, email string) (string, error)
	// Deprecated
	GetUserIdById(ctx context.Context, tenant, id string) (string, error)
	// Deprecated
	GetUserIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
	// Deprecated
	GetUserIdByExternalIdSecond(ctx context.Context, tenant, externalIdSecond, externalSystemId string) (string, error)
}

type userRepository struct {
	driver *neo4j.DriverWithContext
}

func NewUserRepository(driver *neo4j.DriverWithContext) UserRepository {
	return &userRepository{
		driver: driver,
	}
}

func (r *userRepository) GetMatchedUserId(ctx context.Context, tenant, externalSystem, externalId, email string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserRepository.GetMatchedUserId")
	defer spans.Finish()
	spans.LogKV("externalSystem", externalSystem, "externalId", externalId, "email", email)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u1:User)-[:IS_LINKED_WITH {externalId:$userExternalId}]->(e)
				OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u2:User)-[:HAS]->(m:Email)
					WHERE (m.rawEmail=$email OR m.email=$email) AND $email <> '' 
				with coalesce(u1, u2) as user
				where user is not null
				return user.id limit 1`
	spans.LogKV("query", query)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"userExternalId": externalId,
				"email":          email,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	userIDs := dbRecords.([]*db.Record)
	if len(userIDs) == 1 {
		return userIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *userRepository) GetUserIdById(ctx context.Context, tenant, id string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserRepository.GetUserIdById")
	defer spans.Finish()
	spans.LogKV("tenant", tenant, "id", id)

	query := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
				return u.id order by u.createdAt`
	spans.LogKV("query", query)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant": tenant,
			"userId": id,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *userRepository) GetUserIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserRepository.GetUserIdByExternalId")
	defer spans.Finish()
	spans.LogKV("tenant", tenant, "externalId", externalId, "externalSystemId", externalSystemId)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				return u.id order by u.createdAt`
	spans.LogKV("query", query)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":           tenant,
			"externalId":       externalId,
			"externalSystemId": externalSystemId,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *userRepository) GetUserIdByExternalIdSecond(ctx context.Context, tenant, externalIdSecond, externalSystemId string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserRepository.GetUserIdByExternalIdSecond")
	defer spans.Finish()
	spans.LogKV("tenant", tenant, "externalIdSecond", externalIdSecond, "externalSystemId", externalSystemId)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:IS_LINKED_WITH {externalIdSecond:$externalIdSecond}]->(e)
				return u.id order by u.createdAt`
	spans.LogKV("query", query)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":           tenant,
			"externalIdSecond": externalIdSecond,
			"externalSystemId": externalSystemId,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}
