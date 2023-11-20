package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ContractRepository interface {
	CreateForOrganization(ctx context.Context, tenant, contractId string, evt event.ContractCreateEvent) error
	Update(ctx context.Context, tenant, contractId string, evt event.ContractUpdateEvent) error
	GetContractById(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetContractByServiceLineItemId(ctx context.Context, tenant string, serviceLineItemId string) (*dbtype.Node, error)
	GetContractByOpportunityId(ctx context.Context, tenant string, opportunityId string) (*dbtype.Node, error)
}

type contractRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewContractRepository(driver *neo4j.DriverWithContext, database string) ContractRepository {
	return &contractRepository{
		driver:   driver,
		database: database,
	}
}

func (r *contractRepository) CreateForOrganization(ctx context.Context, tenant, contractId string, evt event.ContractCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.CreateForOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("contractId", contractId), log.Object("event", evt))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$orgId})
							MERGE (t)<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})<-[:HAS_CONTRACT]-(org)
							ON CREATE SET 
								ct:Contract_%s,
								ct.createdAt=$createdAt,
								ct.updatedAt=$updatedAt,
								ct.source=$source,
								ct.sourceOfTruth=$sourceOfTruth,
								ct.appSource=$appSource,
								ct.name=$name,
								ct.contractUrl=$contractUrl,
								ct.status=$status,
								ct.renewalCycle=$renewalCycle,
								ct.signedAt=$signedAt,
								ct.serviceStartedAt=$serviceStartedAt
							WITH ct, t
							OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$createdByUserId}) 
							WHERE $createdByUserId <> ""
							FOREACH (ignore IN CASE WHEN u IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (ct)-[:CREATED_BY]->(u))
							`, tenant)
	params := map[string]any{
		"tenant":           tenant,
		"contractId":       contractId,
		"orgId":            evt.OrganizationId,
		"createdAt":        evt.CreatedAt,
		"updatedAt":        evt.UpdatedAt,
		"source":           helper.GetSource(evt.Source.Source),
		"sourceOfTruth":    helper.GetSourceOfTruth(evt.Source.Source),
		"appSource":        helper.GetAppSource(evt.Source.AppSource),
		"name":             evt.Name,
		"contractUrl":      evt.ContractUrl,
		"status":           evt.Status,
		"renewalCycle":     evt.RenewalCycle,
		"signedAt":         utils.TimePtrFirstNonNilNillableAsAny(evt.SignedAt),
		"serviceStartedAt": utils.TimePtrFirstNonNilNillableAsAny(evt.ServiceStartedAt),
		"createdByUserId":  evt.CreatedByUserId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}

func (r *contractRepository) Update(ctx context.Context, tenant, contractId string, evt event.ContractUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("contractId", contractId), log.Object("event", evt))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
				SET 
				ct.name = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR ct.name is null OR ct.name = '' THEN $name ELSE ct.name END,	
				ct.contractUrl = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR ct.contractUrl is null OR ct.contractUrl = '' THEN $contractUrl ELSE ct.contractUrl END,	
				ct.signedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $signedAt ELSE ct.signedAt END,
				ct.endedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $endedAt ELSE ct.endedAt END,
				ct.serviceStartedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $serviceStartedAt ELSE ct.serviceStartedAt END,
				ct.status = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $status ELSE ct.status END,
				ct.renewalCycle = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $renewalCycle ELSE ct.renewalCycle END,
				ct.updatedAt = $updatedAt,
				ct.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE ct.sourceOfTruth END`
	params := map[string]any{
		"tenant":           tenant,
		"contractId":       contractId,
		"updatedAt":        evt.UpdatedAt,
		"name":             evt.Name,
		"contractUrl":      evt.ContractUrl,
		"status":           evt.Status,
		"renewalCycle":     evt.RenewalCycle,
		"signedAt":         utils.TimePtrFirstNonNilNillableAsAny(evt.SignedAt),
		"serviceStartedAt": utils.TimePtrFirstNonNilNillableAsAny(evt.ServiceStartedAt),
		"endedAt":          utils.TimePtrFirstNonNilNillableAsAny(evt.EndedAt),
		"sourceOfTruth":    helper.GetSourceOfTruth(evt.Source),
		"overwrite":        helper.GetSourceOfTruth(evt.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
}

func (r *contractRepository) GetContractById(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.GetContractById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$id}) RETURN c`
	params := map[string]any{
		"tenant": tenant,
		"id":     contractId,
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

func (r *contractRepository) GetContractByServiceLineItemId(ctx context.Context, tenant string, serviceLineItemId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.GetContractByServiceLineItemId")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	cypher := fmt.Sprintf(`MATCH (sli:ServiceLineItem {id:$id})<-[:HAS_SERVICE]-(c:Contract:Contract_%s) RETURN c limit 1`, tenant)
	params := map[string]any{
		"id": serviceLineItemId,
	}
	span.LogFields(log.String("query", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		return nil, nil
	} else {
		return records[0], nil
	}
}

func (r *contractRepository) GetContractByOpportunityId(ctx context.Context, tenant string, opportunityId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.GetContractByServiceLineItemId")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId))

	cypher := fmt.Sprintf(`MATCH (:Opportunity {id:$id})<-[:HAS_OPPORTUNITY]-(c:Contract:Contract_%s) RETURN c limit 1`, tenant)
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
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		return nil, nil
	} else {
		return records[0], nil
	}
}
