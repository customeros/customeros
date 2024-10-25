package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantReadRepository interface {
	GetAll(ctx context.Context) ([]*dbtype.Node, error)
	TenantExists(ctx context.Context, name string) (bool, error)
	GetTenantByName(ctx context.Context, tenant string) (*dbtype.Node, error)
	GetTenantByNameIgnoreCase(ctx context.Context, tenant string) (*dbtype.Node, error)
	GetTenantForWorkspaceProvider(ctx context.Context, workspaceName, workspaceProvider string) (*dbtype.Node, error)
	GetTenantForUserEmail(ctx context.Context, email string) (*dbtype.Node, error)
	GetTenantSettings(ctx context.Context, tenant string) (*dbtype.Node, error)
	GetTenantBillingProfiles(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetTenantBillingProfileById(ctx context.Context, tenant, id string) (*dbtype.Node, error)
}

type tenantReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
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

func (r *tenantReadRepository) GetAll(ctx context.Context) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantReadRepository.GetAll")
	defer span.Finish()

	cypher := `MATCH (t:Tenant) return t`
	params := map[string]any{}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), nil
}

func (r *tenantReadRepository) TenantExists(ctx context.Context, tenantName string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.TenantExists")
	defer span.Finish()
	span.LogFields(log.String("tenantName", tenantName))

	session := (*r.driver).NewSession(
		ctx,
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeRead,
			BoltLogger: neo4j.ConsoleBoltLogger(),
		},
	)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$name}) RETURN t.name`,
			map[string]interface{}{
				"name": tenantName,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return false, err
	}
	if len(records.([]*neo4j.Record)) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *tenantReadRepository) GetTenantByName(ctx context.Context, tenant string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantReadRepository.GetTenantByName")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

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

func (r *tenantReadRepository) GetTenantByNameIgnoreCase(ctx context.Context, tenant string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantReadRepository.GetTenantByNameIgnoreCase")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (t:Tenant) where lower(t.name) = lower($tenant)  RETURN t`
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
	tracing.TagComponentNeo4jRepository(span)

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
	tracing.TagComponentNeo4jRepository(span)
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

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

func (r *tenantReadRepository) GetTenantBillingProfiles(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantReadRepository.GetTenantBillingProfiles")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

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

func (r *tenantReadRepository) GetTenantBillingProfileById(ctx context.Context, tenant, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantReadRepository.GetTenantBillingProfileById")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})-[:HAS_BILLING_PROFILE]->(tbp:TenantBillingProfile {id:$id})
			RETURN tbp`
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})

	if err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.Bool("result.found", false))
		return nil, err
	}

	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}
