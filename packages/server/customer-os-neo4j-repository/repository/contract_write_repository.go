package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ContractCreateFields struct {
	OrganizationId   string       `json:"organizationId"`
	Name             string       `json:"name"`
	ContractUrl      string       `json:"contractUrl"`
	CreatedByUserId  string       `json:"createdByUserId"`
	ServiceStartedAt *time.Time   `json:"serviceStartedAt,omitempty"`
	SignedAt         *time.Time   `json:"signedAt,omitempty"`
	RenewalCycle     string       `json:"renewalCycle"`
	RenewalPeriods   *int64       `json:"renewalPeriods,omitempty"`
	Status           string       `json:"status"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	SourceFields     model.Source `json:"sourceFields"`
}

type ContractUpdateFields struct {
	Name             string     `json:"name"`
	ContractUrl      string     `json:"contractUrl"`
	Status           string     `json:"status"`
	Source           string     `json:"source"`
	RenewalPeriods   *int64     `json:"renewalPeriods"`
	RenewalCycle     string     `json:"renewalCycle"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	ServiceStartedAt *time.Time `json:"serviceStartedAt"`
	SignedAt         *time.Time `json:"signedAt"`
	EndedAt          *time.Time `json:"endedAt"`
}

type ContractWriteRepository interface {
	CreateForOrganization(ctx context.Context, tenant, contractId string, data ContractCreateFields) error
	UpdateAndReturn(ctx context.Context, tenant, contractId string, data ContractUpdateFields) (*dbtype.Node, error)
	UpdateStatus(ctx context.Context, tenant, contractId, status string, serviceStartedAt, endedAt *time.Time) error
	SuspendActiveRenewalOpportunity(ctx context.Context, tenant, contractId string) error
	ActivateSuspendedRenewalOpportunity(ctx context.Context, tenant, contractId string) error
	ContractCausedOnboardingStatusChange(ctx context.Context, tenant, contractId string) error
}

type contractWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewContractWriteRepository(driver *neo4j.DriverWithContext, database string) ContractWriteRepository {
	return &contractWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *contractWriteRepository) CreateForOrganization(ctx context.Context, tenant, contractId string, data ContractCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.CreateForOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)
	tracing.LogObjectAsJson(span, "data", data)

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
								ct.renewalPeriods=$renewalPeriods,
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
		"orgId":            data.OrganizationId,
		"createdAt":        data.CreatedAt,
		"updatedAt":        data.UpdatedAt,
		"source":           data.SourceFields.Source,
		"sourceOfTruth":    data.SourceFields.Source,
		"appSource":        data.SourceFields.AppSource,
		"name":             data.Name,
		"contractUrl":      data.ContractUrl,
		"status":           data.Status,
		"renewalCycle":     data.RenewalCycle,
		"renewalPeriods":   data.RenewalPeriods,
		"signedAt":         utils.TimePtrFirstNonNilNillableAsAny(data.SignedAt),
		"serviceStartedAt": utils.TimePtrFirstNonNilNillableAsAny(data.ServiceStartedAt),
		"createdByUserId":  data.CreatedByUserId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) UpdateAndReturn(ctx context.Context, tenant, contractId string, data ContractUpdateFields) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.UpdateAndReturn")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
				SET 
				ct.name = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR ct.name is null OR ct.name = '' THEN $name ELSE ct.name END,	
				ct.contractUrl = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR ct.contractUrl is null OR ct.contractUrl = '' THEN $contractUrl ELSE ct.contractUrl END,	
				ct.signedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $signedAt ELSE ct.signedAt END,
				ct.endedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $endedAt ELSE ct.endedAt END,
				ct.serviceStartedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $serviceStartedAt ELSE ct.serviceStartedAt END,
				ct.status = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $status ELSE ct.status END,
				ct.renewalCycle = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $renewalCycle ELSE ct.renewalCycle END,
				ct.renewalPeriods = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $renewalPeriods ELSE ct.renewalPeriods END,
				ct.updatedAt = $updatedAt,
				ct.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE ct.sourceOfTruth END
				RETURN ct`
	params := map[string]any{
		"tenant":           tenant,
		"contractId":       contractId,
		"updatedAt":        data.UpdatedAt,
		"name":             data.Name,
		"contractUrl":      data.ContractUrl,
		"status":           data.Status,
		"renewalCycle":     data.RenewalCycle,
		"renewalPeriods":   data.RenewalPeriods,
		"signedAt":         utils.TimePtrFirstNonNilNillableAsAny(data.SignedAt),
		"serviceStartedAt": utils.TimePtrFirstNonNilNillableAsAny(data.ServiceStartedAt),
		"endedAt":          utils.TimePtrFirstNonNilNillableAsAny(data.EndedAt),
		"sourceOfTruth":    data.Source,
		"overwrite":        data.Source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
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
	return result.(*dbtype.Node), nil
}

func (r *contractWriteRepository) UpdateStatus(ctx context.Context, tenant, contractId, status string, serviceStartedAt, endedAt *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.UpdateStatus")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
				SET 
					ct.status=$status,
					ct.serviceStartedAt=$serviceStartedAt,
					ct.endedAt=$endedAt,
					ct.updatedAt=$updatedAt
							`
	params := map[string]any{
		"tenant":           tenant,
		"contractId":       contractId,
		"status":           status,
		"serviceStartedAt": utils.TimePtrFirstNonNilNillableAsAny(serviceStartedAt),
		"endedAt":          utils.TimePtrFirstNonNilNillableAsAny(endedAt),
		"updatedAt":        utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) SuspendActiveRenewalOpportunity(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.SuspendActiveRenewalOpportunity")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})-[r:ACTIVE_RENEWAL]->(op:RenewalOpportunity)
				SET op.internalStage=$internalStageSuspended, 
					op.updatedAt=$updatedAt
				MERGE (ct)-[:SUSPENDED_RENEWAL]->(op)
				DELETE r`
	params := map[string]any{
		"tenant":                 tenant,
		"contractId":             contractId,
		"internalStageSuspended": "SUSPENDED",
		"updatedAt":              utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) ActivateSuspendedRenewalOpportunity(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.ActivateSuspendedRenewalOpportunity")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})-[r:SUSPENDED_RENEWAL]->(op:RenewalOpportunity)
				SET op.internalStage=$internalStage, 
					op.updatedAt=$updatedAt
				MERGE (ct)-[:ACTIVE_RENEWAL]->(op)
				DELETE r`
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    contractId,
		"internalStage": neo4jentity.OpportunityInternalStageOpen.String(),
		"updatedAt":     utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) ContractCausedOnboardingStatusChange(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.ContractCausedOnboardingStatusChange")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
				SET ct.triggeredOnboardingStatusChange=true`
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
