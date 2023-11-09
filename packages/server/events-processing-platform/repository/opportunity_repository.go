package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OpportunityRepository interface {
	CreateForOrganization(ctx context.Context, tenant, opportunityId string, evt event.OpportunityCreateEvent) error
	ReplaceOwner(ctx context.Context, tenant, opportunityId, userId string) error
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

func (r *opportunityRepository) ReplaceOwner(ctx context.Context, tenant, opportunityId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.ReplaceOwner")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId), log.String("userId", userId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity {id:$opportunityId})
			OPTIONAL MATCH (:User)-[rel:OWNS]->(op)
			DELETE rel
			WITH op, t
			MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
			WHERE (u.internal=false OR u.internal is null) AND (u.bot=false OR u.bot is null)
			MERGE (u)-[:OWNS]->(op)
			SET op.updatedAt=$now`

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
		"userId":        userId,
		"now":           utils.Now(),
	})
}
