package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantAndOpportunityId struct {
	Tenant        string
	OpportunityId string
}

type OpportunityReadRepository interface {
	GetOpportunityById(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error)
	GetActiveRenewalOpportunityForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetActiveRenewalOpportunitiesForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error)
	GetRenewalOpportunitiesForClosingAsLost(ctx context.Context) ([]TenantAndOpportunityId, error)
}

type opportunityReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOpportunityReadRepository(driver *neo4j.DriverWithContext, database string) OpportunityReadRepository {
	return &opportunityReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *opportunityReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *opportunityReadRepository) GetOpportunityById(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityReadRepository.GetOpportunityById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$id}) WHERE op:Opportunity_%s RETURN op`, tenant)
	params := map[string]any{
		"id": opportunityId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *opportunityReadRepository) GetActiveRenewalOpportunityForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityReadRepository.GetActiveRenewalOpportunityForContract")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})-[:ACTIVE_RENEWAL]->(op:RenewalOpportunity)
				WHERE op.internalStage=$internalStage
				RETURN op`
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    contractId,
		"internalStage": enum.OpportunityInternalStageOpen.String(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if result.Next(ctx) {
		node := result.Record().Values[0].(dbtype.Node)
		span.LogFields(log.Bool("result.found", true))
		return &node, nil
	}

	span.LogFields(log.Bool("result.found", false))
	return nil, nil
}

func (r *opportunityReadRepository) GetActiveRenewalOpportunitiesForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityReadRepository.GetActiveRenewalOpportunitiesForOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("organizationId", organizationId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})-[:HAS_CONTRACT]->(c:Contract)-[:ACTIVE_RENEWAL]->(op:Opportunity)
				WHERE op:RenewalOpportunity AND op.internalStage=$internalStage
				RETURN op`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"internalStage":  enum.OpportunityInternalStageOpen.String(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), nil
}

func (r *opportunityReadRepository) GetRenewalOpportunitiesForClosingAsLost(ctx context.Context) ([]TenantAndOpportunityId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.GetRenewalOpportunitiesForClosingAsLost")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")

	cypher := `MATCH (t:Tenant)<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:ACTIVE_RENEWAL]->(op:RenewalOpportunity)
				WHERE 
					c.status = $endedStatus AND 
					c.endedAt < $now AND 
					op.internalStage = $internalStageOpen AND 
					(op.techRolloutRenewalRequestedAt IS NULL OR op.techRolloutRenewalRequestedAt + duration({hours: 1}) < $now)
				RETURN t.name, op.id LIMIT 100`
	params := map[string]any{
		"now":               utils.Now(),
		"endedStatus":       enum.ContractStatusEnded.String(),
		"internalStageOpen": enum.OpportunityInternalStageOpen.String(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)

	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantAndOpportunityId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOpportunityId{
				Tenant:        v.Values[0].(string),
				OpportunityId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}
