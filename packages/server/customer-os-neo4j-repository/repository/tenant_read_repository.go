package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantReadRepository interface {
	GetTenantByName(ctx context.Context, tenant string) (*dbtype.Node, error)
	GetTenantForWorkspaceProvider(ctx context.Context, workspaceName, workspaceProvider string) (*dbtype.Node, error)
	GetTenantForUserEmail(ctx context.Context, email string) (*dbtype.Node, error)
	GetTenantSettings(ctx context.Context, tenant string) (*dbtype.Node, error)
	GetTenantBillingProfiles(ctx context.Context, tenant string) ([]*dbtype.Node, error)
}

type tenantReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func (r *tenantReadRepository) GetTenantBillingProfiles(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantReadRepository.GetTenantBillingProfiles")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})-[:HAS_BILLING_PROFILE]->(tbp:TenantBillingProfile)
			RETURN tbp ORDER BY tbp.createdAt ASC`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})

	if err != nil {
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	if result == nil {
		return nil, nil
	}
	return result.([]*dbtype.Node), nil
}

func NewTenantReadRepository(driver *neo4j.DriverWithContext, database string) TenantReadRepository {
	return &tenantReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *tenantReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *tenantReadRepository) GetTenantByName(ctx context.Context, tenant string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantReadRepository.GetTenantByName")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (t:Tenant {name:$tenant}) RETURN t`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
	})

	if err != nil {
		return nil, err
	}

	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}

func (r *tenantReadRepository) GetTenantForWorkspaceProvider(ctx context.Context, workspaceName, workspaceProvider string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantReadRepository.GetTenantForWorkspaceProvider")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")

	cypher := `MATCH (t:Tenant)-[:HAS_WORKSPACE]->(w:Workspace)
			WHERE w.name=$name AND w.provider=$provider
			RETURN DISTINCT t`
	params := map[string]any{
		"name":     workspaceName,
		"provider": workspaceProvider,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})

	if err != nil {
		return nil, err
	}

	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}

func (r *tenantReadRepository) GetTenantForUserEmail(ctx context.Context, email string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantReadRepository.GetTenantForUserEmail")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.String("email", email))

	cypher := `MATCH (t:Tenant)<-[:USER_BELONGS_TO_TENANT]-(:User)-[:HAS]->(e:Email)
		WHERE e.email=$email OR e.rawEmail=$email
		RETURN DISTINCT t order by t.createdAt ASC LIMIT 1`
	params := map[string]any{
		"email": email,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}

	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}

func (r *tenantReadRepository) GetTenantSettings(ctx context.Context, tenant string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantReadRepository.GetTenantSettings")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})-[:HAS_SETTINGS]->(ts:TenantSettings)
			RETURN ts`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
	})

	if err != nil {
		return nil, err
	}

	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}
