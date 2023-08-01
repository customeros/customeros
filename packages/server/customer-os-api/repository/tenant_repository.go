package repository

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type TenantRepository interface {
	Merge(ctx context.Context, tenant entity.TenantEntity) (*dbtype.Node, error)
	GetForWorkspace(ctx context.Context, workspaceEntity entity.WorkspaceEntity) (*dbtype.Node, error)
	LinkWithWorkspace(ctx context.Context, tenant string, workspace entity.WorkspaceEntity) (bool, error)
	GetByName(ctx context.Context, tenant string) (*dbtype.Node, error)
}

type tenantRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTenantRepository(driver *neo4j.DriverWithContext) TenantRepository {
	return &tenantRepository{
		driver: driver,
	}
}

func (r *tenantRepository) LinkWithWorkspace(ctx context.Context, tenant string, workspace entity.WorkspaceEntity) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.LinkWithWorkspace")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	query := `
			MATCH (t:Tenant {name:$tenant})
			MATCH (w:Workspace {name:$name, provider:$provider})
			WHERE NOT ()-[:HAS_WORKSPACE]->(w)
			CREATE (t)-[:HAS_WORKSPACE]->(w)
			RETURN t`
	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":   tenant,
				"name":     workspace.Name,
				"provider": workspace.Provider,
			})
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		return false, err
	}
	convertedResult, isOk := result.([]*dbtype.Node)
	if !isOk {
		return false, errors.New("LinkWithWorkspace: cannot convert result")
	}
	if len(convertedResult) == 0 {
		return false, nil
	}
	return true, nil
}

func (r *tenantRepository) Merge(ctx context.Context, tenant entity.TenantEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MERGE (t:Tenant {name:$name}) " +
		" ON CREATE SET " +
		"  t.id=randomUUID(), " +
		"  t.createdAt=$now, " +
		"  t.updatedAt=$now, " +
		"  t.source=$source, " +
		"  t.appSource=$appSource " +
		" RETURN t"

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"name":      tenant.Name,
				"source":    tenant.Source,
				"appSource": tenant.AppSource,
				"now":       utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *tenantRepository) GetByName(ctx context.Context, tenant string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.GetByName")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$name}) " +
		" RETURN t"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"name": tenant,
			})
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})

	if err != nil {
		return nil, err
	}
	convertedResult, isOk := result.([]*dbtype.Node)
	if !isOk {
		return nil, errors.New("GetByName: cannot convert result")
	}
	if len(convertedResult) == 0 {
		return nil, nil
	}
	return convertedResult[0], err
}

func (r *tenantRepository) GetForWorkspace(ctx context.Context, workspaceEntity entity.WorkspaceEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.GetForWorkspace")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant)-[:HAS_WORKSPACE]->(w:Workspace)
			WHERE w.name=$name AND w.provider=$provider
			RETURN DISTINCT t LIMIT 1`,
			map[string]any{
				"name":     workspaceEntity.Name,
				"provider": workspaceEntity.Provider,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	convertedResult, isOk := result.([]*dbtype.Node)
	if !isOk {
		return nil, errors.New("GetForWorkspace: cannot convert result")
	}
	if len(convertedResult) == 0 {
		return nil, nil
	}
	return convertedResult[0], err
}
