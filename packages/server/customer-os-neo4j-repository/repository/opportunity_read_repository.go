package neo4j_repository

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type TenantAndOpportunityId struct {
	Tenant        string
	OpportunityId string
}

type OpportunityReadRepository interface {
	GetOpportunityById(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) (*dbtype.Node, error)
	GetActiveRenewalOpportunityForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetActiveRenewalOpportunitiesForOrganization(ctx context.Context, tenant, organizationId string, includeDraftContracts bool) ([]*dbtype.Node, error)
	GetRenewalOpportunitiesForClosingAsLost(ctx context.Context, limit int) ([]TenantAndOpportunityId, error)
	GetPreviousClosedWonRenewalOpportunityForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error)
	GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	GetForTasks(ctx context.Context, tenant string, taskIds []string) ([]*utils.DbNodeAndId, error)
	GetPaginatedOpportunitiesLinkedToAnOrganization(ctx context.Context, tenant string, skip, limit int) (*utils.DbNodesWithTotalCount, error)
}

type opportunityReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func (r *opportunityReadRepository) GetPreviousClosedWonRenewalOpportunityForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OpportunityReadRepository.GetPreviousClosedWonRenewalOpportunityForContract")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})-[:HAS_OPPORTUNITY]->(op:RenewalOpportunity)
				WHERE op.internalStage=$internalStage AND op.renewedAt < $now
				RETURN op ORDER BY op.renewedAt DESC limit 1`
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    contractId,
		"now":           utils.Now(),
		"internalStage": enum.OpportunityInternalStageClosedWon.String(),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if result.Next(ctx) {
		node := result.Record().Values[0].(dbtype.Node)
		spans.LogKV("result.found", true)
		return &node, nil
	}

	spans.LogKV("result.found", false)
	return nil, nil
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

func (r *opportunityReadRepository) GetOpportunityById(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OpportunityReadRepository.GetOpportunityById")
	defer spans.Finish()

	spans.TagEntity(opportunityId)

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$id}) WHERE op:Opportunity_%s RETURN op`, tenant)
	params := map[string]any{
		"id": opportunityId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *opportunityReadRepository) GetActiveRenewalOpportunityForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OpportunityReadRepository.GetActiveRenewalOpportunityForContract")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})-[:ACTIVE_RENEWAL]->(op:RenewalOpportunity)
				WHERE op.internalStage=$internalStage
				RETURN op ORDER BY op.renewedAt DESC`
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    contractId,
		"internalStage": enum.OpportunityInternalStageOpen.String(),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if result.Next(ctx) {
		node := result.Record().Values[0].(dbtype.Node)
		spans.LogKV("result.found", true)
		return &node, nil
	}

	spans.LogKV("result.found", false)
	return nil, nil
}

func (r *opportunityReadRepository) GetActiveRenewalOpportunitiesForOrganization(ctx context.Context, tenant, organizationId string, includeDraftContracts bool) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OpportunityReadRepository.GetActiveRenewalOpportunitiesForOrganization")
	defer spans.Finish()

	spans.LogKV("organizationId", organizationId)
	spans.LogKV("includeDraftContracts", includeDraftContracts)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})-[:HAS_CONTRACT]->(c:Contract)-[:ACTIVE_RENEWAL]->(op:Opportunity)
				WHERE op:RenewalOpportunity AND op.internalStage=$internalStage AND (c.status<>$draftStatus OR $includeDraftContracts)
				RETURN op`
	params := map[string]any{
		"tenant":                tenant,
		"organizationId":        organizationId,
		"internalStage":         enum.OpportunityInternalStageOpen.String(),
		"draftStatus":           enum.ContractStatusDraft.String(),
		"includeDraftContracts": includeDraftContracts,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), nil
}

func (r *opportunityReadRepository) GetRenewalOpportunitiesForClosingAsLost(ctx context.Context, limit int) ([]TenantAndOpportunityId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContractRepository.GetRenewalOpportunitiesForClosingAsLost")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	cypher := `MATCH (t:Tenant)<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:ACTIVE_RENEWAL]->(op:RenewalOpportunity)
				WHERE 
					c.status = $endedStatus AND 
					c.endedAt < $now AND 
					op.internalStage = $internalStageOpen AND 
					(op.techRolloutRenewalRequestedAt IS NULL OR op.techRolloutRenewalRequestedAt + duration({hours: 1}) < $now)
				RETURN t.name, op.id LIMIT $limit`
	params := map[string]any{
		"now":               utils.Now(),
		"endedStatus":       enum.ContractStatusEnded.String(),
		"internalStageOpen": enum.OpportunityInternalStageOpen.String(),
		"limit":             limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *opportunityReadRepository) GetForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OpportunityRepository.GetForContracts")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:HAS_OPPORTUNITY]->(op:Opportunity)
			WHERE c.id IN $contractIds
			RETURN op, c.id ORDER BY op.createdAt DESC`
	params := map[string]any{
		"tenant":      tenant,
		"contractIds": contractIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *opportunityReadRepository) GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OpportunityRepository.GetForOrganizations")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_OPPORTUNITY]->(op:Opportunity)
			WHERE org.id IN $orgIds
			RETURN op, org.id ORDER BY op.createdAt DESC`
	params := map[string]any{
		"tenant": tenant,
		"orgIds": organizationIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *opportunityReadRepository) GetForTasks(ctx context.Context, tenant string, taskIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OpportunityRepository.GetForTasks")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:TASK_BELONGS_TO_TENANT]-(tsk:Task)-[:LINKED_TO]->(op:Opportunity)
			WHERE tsk.id IN $taskIds
			RETURN op, tsk.id ORDER BY op.createdAt DESC`
	params := map[string]any{
		"tenant":  tenant,
		"taskIds": taskIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndId)))
	return result.([]*utils.DbNodeAndId), err
}

func (r *opportunityReadRepository) GetPaginatedOpportunitiesLinkedToAnOrganization(ctx context.Context, tenant string, skip, limit int) (*utils.DbNodesWithTotalCount, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OpportunityRepository.GetPaginatedOpportunitiesLinkedToAnOrganization")
	defer spans.Finish()

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	countCypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {hide:false})-[:HAS_OPPORTUNITY]->(op:Opportunity) RETURN count(op) as count`
	countParams := map[string]any{
		"tenant": tenant,
	}

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {hide:false})-[:HAS_OPPORTUNITY]->(op:Opportunity) RETURN op ORDER BY op.createdAt SKIP $skip LIMIT $limit`
	params := map[string]any{
		"tenant": tenant,
		"skip":   skip,
		"limit":  limit,
	}
	spans.LogKV("countCypher", countCypher)
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)
	spans.LogObjectAsJson("countParams", countParams)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, countCypher, countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	spans.LogKV("result.count", len(dbNodesWithTotalCount.Nodes))
	return dbNodesWithTotalCount, nil
}
