package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type AuthenticationWriteRepository interface {
	CreateAuthentication(ctx context.Context, tx neo4j.ManagedTransaction, data neo4j_entity.AuthenticationEntity) (string, error)
	CreateAuthenticationUser(ctx context.Context, tx neo4j.ManagedTransaction, authenticationId string, data neo4j_entity.AuthenticationUserEntity) (string, error)
	LinkAuthenticationWithAuthenticationUser(ctx context.Context, tx neo4j.ManagedTransaction, authId, authUserId string) error
	LinkAuthenticationUserWithTenant(ctx context.Context, tx *neo4j.ManagedTransaction, authUserId, tenant string) error
	UnlinkAuthenticationUserWithTenant(ctx context.Context, tx *neo4j.ManagedTransaction, authUserId, tenant string) error
	LinkAuthenticationUserWithUser(ctx context.Context, tx neo4j.ManagedTransaction, authUserId, userId string) error
	SetDefaultTenant(ctx context.Context, tx neo4j.ManagedTransaction, authUserId, defaultTenant string) error
	SetCurrentTenant(ctx context.Context, tx *neo4j.ManagedTransaction, authUserId, currentTenant string) error
}

type authenticationWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewAuthenticationWriteRepository(driver *neo4j.DriverWithContext, database string) AuthenticationWriteRepository {
	return &authenticationWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *authenticationWriteRepository) CreateAuthentication(c context.Context, tx neo4j.ManagedTransaction, data neo4j_entity.AuthenticationEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "AuthenticationWriteRepository.CreateAuthentication")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MERGE (a:Authentication {authId:$authId, provider:$provider})
				ON CREATE SET a.id=randomUUID(),
							  a.identityId=$identityId,
							  a.createdAt=$createdAt
				RETURN a.id`
	params := map[string]any{
		"authId":     data.AuthId,
		"provider":   data.Provider,
		"identityId": data.IdentityId,
		"createdAt":  utils.NowIfZero(data.CreatedAt),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
		tracing.TraceErr(span, err)
		return "", err
	} else {
		return utils.ExtractSingleRecordFirstValueAsString(ctx, queryResult, err)
	}
}

func (r *authenticationWriteRepository) CreateAuthenticationUser(ctx context.Context, tx neo4j.ManagedTransaction, authenticationId string, data neo4j_entity.AuthenticationUserEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationWriteRepository.CreateAuthenticationUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (a:Authentication {id:$authenticationId})
				MERGE (u:AuthenticationUser {id:randomUUID()})
				ON CREATE SET 
					u.createdAt=datetime(),
					u.firstName=$firstName,
					u.lastName=$lastName
				MERGE (a)-[:%s]->(u)
				RETURN u.id`, model.HAS.String())
	params := map[string]any{
		"authenticationId": authenticationId,
		"firstName":        data.FirstName,
		"lastName":         data.LastName,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
		tracing.TraceErr(span, err)
		return "", err
	} else {
		return utils.ExtractSingleRecordFirstValueAsString(ctx, queryResult, err)
	}
}

func (r *authenticationWriteRepository) LinkAuthenticationWithAuthenticationUser(ctx context.Context, tx neo4j.ManagedTransaction, authId, authUserId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationWriteRepository.LinkAuthenticationWithAuthenticationUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogKV("authUserId", authUserId)
	span.LogKV("authId", authId)

	cypher := fmt.Sprintf(`MATCH (a:Authentication {id:$authId}), (au:AuthenticationUser {id:$authUserId}) MERGE (a)-[:%s]->(au)`, model.HAS.String())
	params := map[string]any{
		"authId":     authId,
		"authUserId": authUserId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if _, err := tx.Run(ctx, cypher, params); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *authenticationWriteRepository) LinkAuthenticationUserWithTenant(ctx context.Context, tx *neo4j.ManagedTransaction, authUserId, tenant string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "AuthenticationWriteRepository.LinkAuthenticationUserWithTenant")
	defer spans.Finish()

	spans.LogKV("authUserId", authUserId)
	spans.LogKV("tenant", tenant)

	cypher := fmt.Sprintf(`MATCH (u:AuthenticationUser {id:$authUserId}), (t:Tenant {name:$tenant}) MERGE (u)-[:%s]->(t)`, model.HAS_WORKSPACE.String())
	params := map[string]any{
		"tenant":     tenant,
		"authUserId": authUserId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *authenticationWriteRepository) UnlinkAuthenticationUserWithTenant(ctx context.Context, tx *neo4j.ManagedTransaction, authUserId, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationWriteRepository.UnlinkAuthenticationUserWithTenant")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogKV("authUserId", authUserId)
	span.LogKV("tenant", tenant)

	cypher := `MATCH (u:AuthenticationUser {id:$authUserId})-[r:HAS_WORKSPACE]->(t:Tenant {name:$tenant}) delete r`
	params := map[string]any{
		"tenant":     tenant,
		"authUserId": authUserId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *authenticationWriteRepository) LinkAuthenticationUserWithUser(ctx context.Context, tx neo4j.ManagedTransaction, authUserId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationWriteRepository.LinkAuthenticationUserWithUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogKV("authUserId", authUserId)
	span.LogKV("userId", userId)

	cypher := fmt.Sprintf(`MATCH (a:AuthenticationUser {id:$authUserId}), (u:User {id:$userId}) MERGE (u)-[:%s]->(a)`, model.AUTHENTICATED_BY.String())
	params := map[string]any{
		"authUserId": authUserId,
		"userId":     userId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if _, err := tx.Run(ctx, cypher, params); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *authenticationWriteRepository) SetDefaultTenant(ctx context.Context, tx neo4j.ManagedTransaction, authUserId, defaultTenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationWriteRepository.SetDefaultTenant")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogKV("authUserId", authUserId)
	span.LogKV("defaultTenant", defaultTenant)

	cypher := `MATCH (a:AuthenticationUser {id:$authUserId}) set a.defaultTenant = $defaultTenant`
	params := map[string]any{
		"authUserId":    authUserId,
		"defaultTenant": defaultTenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if _, err := tx.Run(ctx, cypher, params); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *authenticationWriteRepository) SetCurrentTenant(ctx context.Context, tx *neo4j.ManagedTransaction, authUserId, currentTenant string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "AuthenticationWriteRepository.SetCurrentTenant")
	defer spans.Finish()
	spans.LogKV("authUserId", authUserId)
	spans.LogKV("currentTenant", currentTenant)

	cypher := `MATCH (a:AuthenticationUser {id:$authUserId}) set a.currentTenant = $currentTenant`
	params := map[string]any{
		"authUserId":    authUserId,
		"currentTenant": currentTenant,
	}
	spans.LogFields(log.String("cypher", cypher))
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		spans.TraceError(err)
	}

	return nil
}
