package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OpportunityRepository interface {
	CreateForOrganization(ctx context.Context, tenant, opportunityId string, evt event.OpportunityCreateEvent) error
	Update(ctx context.Context, tenant, opportunityId string, evt event.OpportunityUpdateEvent) error
	ReplaceOwner(ctx context.Context, tenant, opportunityId, userId string) error

	CreateRenewal(ctx context.Context, tenant, opportunityId string, evt event.OpportunityCreateRenewalEvent) error
	UpdateRenewal(ctx context.Context, tenant, opportunityId string, evt event.OpportunityUpdateRenewalEvent, setUpdatedByUserId bool) error
	UpdateNextCycleDate(ctx context.Context, tenant, opportunityId string, evt event.OpportunityUpdateNextCycleDateEvent) error

	GetOpenRenewalOpportunityForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetOpenRenewalOpportunitiesForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error)
	GetOpportunityById(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error)
	CloseWin(ctx context.Context, tenant, opportunityId string, data event.OpportunityCloseWinEvent) error
	CloseLoose(ctx context.Context, tenant, opportunityId string, data event.OpportunityCloseLooseEvent) error
}

type opportunityRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOpportunityRepository(driver *neo4j.DriverWithContext, database string) OpportunityRepository {
	return &opportunityRepository{
		driver:   driver,
		database: database,
	}
}

func (r *opportunityRepository) CreateForOrganization(ctx context.Context, tenant, opportunityId string, evt event.OpportunityCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.CreateForOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId), log.Object("event", evt))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$orgId})
							MERGE (t)<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity {id:$opportunityId})<-[:HAS_OPPORTUNITY]-(org)
							ON CREATE SET 
								op:Opportunity_%s,
								op.createdAt=$createdAt,
								op.updatedAt=$updatedAt,
								op.source=$source,
								op.sourceOfTruth=$sourceOfTruth,
								op.appSource=$appSource,
								op.name=$name,
								op.amount=$amount,
								op.internalType=$internalType,
								op.externalType=$externalType,
								op.internalStage=$internalStage,
								op.externalStage=$externalStage,
								op.estimatedClosedAt=$estimatedClosedAt,
								op.generalNotes=$generalNotes,
								op.nextSteps=$nextSteps
							WITH op, t
							OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$createdByUserId}) 
							WHERE $createdByUserId <> ""
							FOREACH (ignore IN CASE WHEN u IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (op)-[:CREATED_BY]->(u))
							`, tenant)
	params := map[string]any{
		"tenant":            tenant,
		"opportunityId":     opportunityId,
		"orgId":             evt.OrganizationId,
		"createdAt":         evt.CreatedAt,
		"updatedAt":         evt.UpdatedAt,
		"source":            helper.GetSource(evt.Source.Source),
		"sourceOfTruth":     helper.GetSourceOfTruth(evt.Source.Source),
		"appSource":         helper.GetAppSource(evt.Source.AppSource),
		"name":              evt.Name,
		"amount":            evt.Amount,
		"internalType":      evt.InternalType,
		"externalType":      evt.ExternalType,
		"internalStage":     evt.InternalStage,
		"externalStage":     evt.ExternalStage,
		"estimatedClosedAt": utils.TimePtrFirstNonNilNillableAsAny(evt.EstimatedClosedAt),
		"generalNotes":      evt.GeneralNotes,
		"nextSteps":         evt.NextSteps,
		"createdByUserId":   evt.CreatedByUserId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}

func (r *opportunityRepository) CreateRenewal(ctx context.Context, tenant, opportunityId string, evt event.OpportunityCreateRenewalEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.CreateRenewal")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId), log.Object("event", evt))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
							MERGE (c)-[:HAS_OPPORTUNITY]->(op:Opportunity {id:$opportunityId})
							ON CREATE SET 
								op:Opportunity_%s,
								op:RenewalOpportunity,
								op.createdAt=$createdAt,
								op.updatedAt=$updatedAt,
								op.source=$source,
								op.sourceOfTruth=$sourceOfTruth,
								op.appSource=$appSource,
								op.internalType=$internalType,
								op.internalStage=$internalStage,
								op.renewalLikelihood=$renewalLikelihood
							WITH op, c
							MERGE (c)-[:ACTIVE_RENEWAL]->(op)
							`, tenant)
	params := map[string]any{
		"tenant":            tenant,
		"opportunityId":     opportunityId,
		"contractId":        evt.ContractId,
		"createdAt":         evt.CreatedAt,
		"updatedAt":         evt.UpdatedAt,
		"source":            helper.GetSource(evt.Source.Source),
		"sourceOfTruth":     helper.GetSourceOfTruth(evt.Source.Source),
		"appSource":         helper.GetAppSource(evt.Source.AppSource),
		"internalType":      evt.InternalType,
		"internalStage":     evt.InternalStage,
		"renewalLikelihood": evt.RenewalLikelihood,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}

func (r *opportunityRepository) ReplaceOwner(ctx context.Context, tenant, opportunityId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.ReplaceOwner")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId), log.String("userId", userId))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}), (op:Opportunity {id:$opportunityId}) WHERE op:Opportunity_%s
			WITH op, t
			OPTIONAL MATCH (:User)-[rel:OWNS]->(op)
			DELETE rel
			WITH op, t
			MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
			WHERE (u.internal=false OR u.internal is null) AND (u.bot=false OR u.bot is null)
			MERGE (u)-[:OWNS]->(op)
			SET op.updatedAt=$now`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
		"userId":        userId,
		"now":           utils.Now(),
	}
	span.LogFields(log.String("query", query))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}

func (r *opportunityRepository) UpdateNextCycleDate(ctx context.Context, tenant, opportunityId string, evt event.OpportunityUpdateNextCycleDateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.UpdateNextCycleDate")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId), log.Object("event", evt))

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) 
							WHERE op:Opportunity_%s AND op.internalStage=$internalStage
							SET op.updatedAt=$updatedAt, op.renewedAt=$renewedAt`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
		"updatedAt":     evt.UpdatedAt,
		"internalStage": string(model.OpportunityInternalStageStringOpen),
		"renewedAt":     utils.TimePtrFirstNonNilNillableAsAny(evt.RenewedAt),
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}

func (r *opportunityRepository) Update(ctx context.Context, tenant, opportunityId string, evt event.OpportunityUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId), log.Object("event", evt))

	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
		"updatedAt":     evt.UpdatedAt,
		"sourceOfTruth": helper.GetSourceOfTruth(evt.Source),
		"overwrite":     helper.GetSourceOfTruth(evt.Source) == constants.SourceOpenline,
	}
	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) WHERE op:Opportunity_%s SET `, tenant)
	if evt.UpdateName() {
		cypher += ` op.name = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR op.name = '' THEN $name ELSE op.name END, `
		params["name"] = evt.Name
	}
	if evt.UpdateAmount() {
		cypher += ` op.amount = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $amount ELSE op.amount END, `
		params["amount"] = evt.Amount
	}
	if evt.UpdateMaxAmount() {
		cypher += ` op.maxAmount = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $maxAmount ELSE op.maxAmount END, `
		params["maxAmount"] = evt.MaxAmount
	}
	cypher += ` op.updatedAt = $updatedAt,
				op.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE op.sourceOfTruth END`
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}

func (r *opportunityRepository) GetOpportunityById(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.GetOpportunityById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId))

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$id}) WHERE op:Opportunity_%s RETURN op`, tenant)
	params := map[string]any{
		"id": opportunityId,
	}
	span.LogFields(log.String("query", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *opportunityRepository) GetOpenRenewalOpportunityForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.GetOpenRenewalOpportunityForContract")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
				MATCH (c)-[:ACTIVE_RENEWAL]->(op:Opportunity)
				WHERE op:RenewalOpportunity AND op.internalStage=$internalStage
				RETURN op`
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    contractId,
		"internalStage": string(model.OpportunityInternalStageStringOpen),
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		node := result.Record().Values[0].(dbtype.Node)
		return &node, nil
	}

	return nil, nil
}

func (r *opportunityRepository) GetOpenRenewalOpportunitiesForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.GetOpenRenewalOpportunitiesForOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})-[:HAS_CONTRACT]->(c:Contract)-[:ACTIVE_RENEWAL]->(op:Opportunity)
				WHERE op:RenewalOpportunity AND op.internalStage=$internalStage
				RETURN op`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"internalStage":  string(model.OpportunityInternalStageStringOpen),
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}

func (r *opportunityRepository) UpdateRenewal(ctx context.Context, tenant, opportunityId string, evt event.OpportunityUpdateRenewalEvent, setUpdatedByUserId bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.UpdateRenewal")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId), log.Object("event", evt))

	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
		"updatedAt":     evt.UpdatedAt,
		"sourceOfTruth": helper.GetSourceOfTruth(evt.Source),
		"overwrite":     helper.GetSourceOfTruth(evt.Source) == constants.SourceOpenline,
	}
	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) WHERE op:RenewalOpportunity AND op:Opportunity_%s 
				SET `, tenant)
	if setUpdatedByUserId {
		params["renewalUpdatedByUserId"] = evt.UpdatedByUserId
		cypher += ` op.renewalUpdatedByUserAt = $updatedAt, 
					op.renewalUpdatedByUserId = $renewalUpdatedByUserId, `
	}
	if evt.UpdateComments() {
		cypher += ` op.comments = $comments, `
		params["comments"] = evt.Comments
	}
	if evt.UpdateAmount() {
		cypher += ` op.amount = $amount, `
		params["amount"] = evt.Amount
	}
	if evt.UpdateRenewalLikelihood() {
		cypher += ` op.renewalLikelihood = $renewalLikelihood, `
		params["renewalLikelihood"] = evt.RenewalLikelihood
	}
	cypher += ` op.updatedAt = $updatedAt,
				op.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE op.sourceOfTruth END`

	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}

func (r *opportunityRepository) CloseWin(ctx context.Context, tenant, opportunityId string, data event.OpportunityCloseWinEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.CloseWin")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId), log.Object("data", data))

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) WHERE op:Opportunity_%s 
							SET 
								op.closedAt=$closedAt, 
								op.internalStage=$internalStage,
								op.updatedAt=$updatedAt
							WITH op
							OPTIONAL MATCH (op)<-[rel:ACTIVE_RENEWAL]-(c:Contract)
							DELETE rel`, tenant)
	params := map[string]any{
		"opportunityId": opportunityId,
		"updatedAt":     data.UpdatedAt,
		"closedAt":      data.ClosedAt,
		"internalStage": string(model.OpportunityInternalStageStringClosedWon),
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}

func (r *opportunityRepository) CloseLoose(ctx context.Context, tenant, opportunityId string, data event.OpportunityCloseLooseEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.CloseLoose")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId), log.Object("data", data))

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) WHERE op:Opportunity_%s 
							SET op.closedAt=$closedAt, 
								op.internalStage=$internalStage,
								op.updatedAt=$updatedAt
							WITH op
							OPTIONAL MATCH (op)<-[rel:ACTIVE_RENEWAL]-(c:Contract)
							DELETE rel`, tenant)
	params := map[string]any{
		"opportunityId": opportunityId,
		"updatedAt":     data.UpdatedAt,
		"closedAt":      data.ClosedAt,
		"internalStage": string(model.OpportunityInternalStageStringClosedLost),
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}
