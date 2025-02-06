package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantHasWorkspace struct {
	Tenant     string
	CreatedBy  string
	IsPersonal bool
}

type AuthenticationReadRepository interface {
	GetByAuthIdAndProvider(ctx context.Context, authId string, provider string) (*dbtype.Node, error)
	GetByAuthId(ctx context.Context, authId string) ([]*dbtype.Node, error)
	GetAuthUser(ctx context.Context, authId string) (*dbtype.Node, error)
	GetTenants(ctx context.Context, authUserId string) ([]*dbtype.Node, error)
	GetTenantsForImpersonation(ctx context.Context, authUserId string) ([]*TenantHasWorkspace, error)
}

type authenticationReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewAuthenticationReadRepository(driver *neo4j.DriverWithContext, database string) AuthenticationReadRepository {
	return &authenticationReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *authenticationReadRepository) GetByAuthIdAndProvider(ctx context.Context, authId string, provider string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationReadRepository.GetPlayerByAuthIdProvider")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf("MATCH (a:Authentication {authId:$authId, provider:$provider}) RETURN a")
	params := map[string]any{
		"authId":   authId,
		"provider": provider,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})

	if err != nil && err.Error() == "Result contains no more records" {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if dbRecord == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))
	return dbRecord.(*dbtype.Node), err
}

func (r *authenticationReadRepository) GetByAuthId(ctx context.Context, authId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationReadRepository.GetByAuthId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (a:Authentication {authId:$authId}) RETURN a`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]any{
				"authId": authId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting player by authId: %w", err)
	}

	return result.([]*dbtype.Node), nil
}

func (r *authenticationReadRepository) GetAuthUser(ctx context.Context, authId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationReadRepository.GetAuthUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (a:Authentication {id: $authId})-[:%s]->(u:AuthenticationUser) RETURN u`, model.HAS.String())

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]any{
				"authId": authId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting users for player: %w", err)
	}

	data := result.(*neo4j.Node)
	if data == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))

	return data, nil
}

func (r *authenticationReadRepository) GetTenants(ctx context.Context, authUserId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationReadRepository.GetTenants")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (u:AuthenticationUser {id:$authUserId})-[:%s]->(t:Tenant) RETURN t`, model.HAS_WORKSPACE.String())

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]any{
				"authUserId": authUserId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error getting users for player: %w", err)
	}

	data := result.([]*dbtype.Node)
	if data == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Int("result.found", len(data)))

	return data, nil
}

func (r *authenticationReadRepository) GetTenantsForImpersonation(ctx context.Context, authUserId string) ([]*TenantHasWorkspace, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationReadRepository.GetTenantsForImpersonation")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `
			MATCH (au:AuthenticationUser{id:$authUserId})-[:HAS_WORKSPACE]-(t:Tenant)
			OPTIONAL MATCH (t)-[r]-(w:Workspace)
			RETURN 
			t.name,
			t.createdBy,
			CASE WHEN w IS NULL THEN true ELSE false END AS isPersonal`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query),
			map[string]any{
				"authUserId": authUserId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("error getting users for player: %w", err)
	}

	var results []*TenantHasWorkspace
	if result != nil {
		for _, v := range result.([]*neo4j.Record) {
			tenant := v.Values[0].(string)
			createdBy := v.Values[1].(string)
			isPersonal := v.Values[2].(bool)

			results = append(results, &TenantHasWorkspace{
				Tenant:     tenant,
				CreatedBy:  createdBy,
				IsPersonal: isPersonal,
			})
		}
	}

	return results, nil
}
