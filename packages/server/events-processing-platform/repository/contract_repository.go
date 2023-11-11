package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
				ct.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE ct.sourceOfTruth END
							`
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
