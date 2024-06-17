package repository

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
)

type TenantRepository interface {
	Merge(ctx context.Context, tenant neo4jentity.TenantEntity) (*dbtype.Node, error)
	LinkWithWorkspace(ctx context.Context, tenant string, workspace entity.WorkspaceEntity) (bool, error)
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

func (r *tenantRepository) Merge(ctx context.Context, tenant neo4jentity.TenantEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := `MERGE (t:Tenant {name:$name}) 
		 ON CREATE SET 
		  t.id=randomUUID(), 
		  t.createdAt=$now, 
		  t.updatedAt=datetime(), 
		  t.source=$source, 
		  t.appSource=$appSource
		WITH t
		MERGE (t)-[:HAS_SETTINGS]->(ts:TenantSettings {tenant:$name})
		ON CREATE SET
			ts.id=randomUUID(),
		  	ts.createdAt=$now,
			ts.updatedAt=datetime(),
			ts.invoicingEnabled=$invoicingEnabled,
			ts.invoicingPostpaid=$invoicingPostpaid,
			ts.opportunityStages=$opportunityStages,
			ts.baseCurrency=$currency
		 RETURN t`
	params := map[string]any{
		"name":              tenant.Name,
		"source":            tenant.Source,
		"appSource":         tenant.AppSource,
		"now":               utils.Now(),
		"invoicingEnabled":  false,
		"invoicingPostpaid": false,
		"opportunityStages": []string{"Identified", "Qualified", "Committed"},
		"currency":          neo4jenum.CurrencyUSD.String(),
	}

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
