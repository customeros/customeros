package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"time"
)

type OpportunityCreateFields struct {
	OrganizationId    string       `json:"organizationId"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
	SourceFields      model.Source `json:"sourceFields"`
	Name              string       `json:"name"`
	Amount            float64      `json:"amount"`
	InternalType      string       `json:"internalType"`
	ExternalType      string       `json:"externalType"`
	InternalStage     string       `json:"internalStage"`
	ExternalStage     string       `json:"externalStage"`
	EstimatedClosedAt *time.Time   `json:"estimatedClosedAt"`
	GeneralNotes      string       `json:"generalNotes"`
	NextSteps         string       `json:"nextSteps"`
	CreatedByUserId   string       `json:"createdByUserId"`
}

type OpportunityWriteRepository interface {
	CreateForOrganization(ctx context.Context, tenant, opportunityId string, data OpportunityCreateFields) error
}

type opportunityWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOpportunityWriteRepository(driver *neo4j.DriverWithContext, database string) OpportunityWriteRepository {
	return &opportunityWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *opportunityWriteRepository) CreateForOrganization(ctx context.Context, tenant, opportunityId string, data OpportunityCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.CreateForOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	span.LogFields(log.String("opportunityId", opportunityId))
	tracing.LogObjectAsJson(span, "data", data)

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
		"orgId":             data.OrganizationId,
		"createdAt":         data.CreatedAt,
		"updatedAt":         data.UpdatedAt,
		"source":            data.SourceFields.Source,
		"sourceOfTruth":     data.SourceFields.Source,
		"appSource":         data.SourceFields.AppSource,
		"name":              data.Name,
		"amount":            data.Amount,
		"internalType":      data.InternalType,
		"externalType":      data.ExternalType,
		"internalStage":     data.InternalStage,
		"externalStage":     data.ExternalStage,
		"estimatedClosedAt": utils.TimePtrFirstNonNilNillableAsAny(data.EstimatedClosedAt),
		"generalNotes":      data.GeneralNotes,
		"nextSteps":         data.NextSteps,
		"createdByUserId":   data.CreatedByUserId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
