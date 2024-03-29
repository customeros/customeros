package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
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

type OpportunityUpdateFields struct {
	UpdatedAt       time.Time `json:"updatedAt"`
	Source          string    `json:"source"`
	Name            string    `json:"name"`
	Amount          float64   `json:"amount"`
	MaxAmount       float64   `json:"maxAmount"`
	UpdateName      bool      `json:"updateName"`
	UpdateAmount    bool      `json:"updateAmount"`
	UpdateMaxAmount bool      `json:"updateMaxAmount"`
}

type RenewalOpportunityCreateFields struct {
	ContractId        string       `json:"contractId"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
	SourceFields      model.Source `json:"sourceFields"`
	InternalType      string       `json:"internalType"`
	InternalStage     string       `json:"internalStage"`
	RenewalLikelihood string       `json:"renewalLikelihood"`
	RenewalApproved   bool         `json:"renewalApproved"`
}

type RenewalOpportunityUpdateFields struct {
	UpdatedAt               time.Time `json:"updatedAt"`
	Source                  string    `json:"source"`
	UpdatedByUserId         string    `json:"updatedByUserId"`
	SetUpdatedByUserId      bool      `json:"setUpdatedByUserId"`
	Comments                string    `json:"comments"`
	Amount                  float64   `json:"amount"`
	RenewalLikelihood       string    `json:"renewalLikelihood"`
	RenewalApproved         bool      `json:"renewalApproved"`
	UpdateComments          bool      `json:"updateComments"`
	UpdateAmount            bool      `json:"updateAmount"`
	UpdateRenewalLikelihood bool      `json:"updateRenewalLikelihood"`
	UpdateRenewalApproved   bool      `json:"updateRenewalApproved"`
}

type OpportunityWriteRepository interface {
	CreateForOrganization(ctx context.Context, tenant, opportunityId string, data OpportunityCreateFields) error
	Update(ctx context.Context, tenant, opportunityId string, data OpportunityUpdateFields) error
	ReplaceOwner(ctx context.Context, tenant, opportunityId, userId string) error
	CreateRenewal(ctx context.Context, tenant, opportunityId string, data RenewalOpportunityCreateFields) error
	UpdateRenewal(ctx context.Context, tenant, opportunityId string, data RenewalOpportunityUpdateFields) error
	UpdateNextRenewalDate(ctx context.Context, tenant, opportunityId string, updatedAt time.Time, renewedAt *time.Time) error
	CloseWin(ctx context.Context, tenant, opportunityId string, updatedAt, closedAt time.Time) error
	CloseLoose(ctx context.Context, tenant, opportunityId string, updatedAt, closedAt time.Time) error
	MarkRenewalRequested(ctx context.Context, tenant, opportunityId string) error
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

func (r *opportunityWriteRepository) Update(ctx context.Context, tenant, opportunityId string, data OpportunityUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	tracing.LogObjectAsJson(span, "data", data)

	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
		"updatedAt":     data.UpdatedAt,
		"sourceOfTruth": data.Source,
		"overwrite":     data.Source == constants.SourceOpenline,
	}
	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) WHERE op:Opportunity_%s SET `, tenant)
	if data.UpdateName {
		cypher += ` op.name = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR op.name = '' THEN $name ELSE op.name END, `
		params["name"] = data.Name
	}
	if data.UpdateAmount {
		cypher += ` op.amount = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $amount ELSE op.amount END, `
		params["amount"] = data.Amount
	}
	if data.UpdateMaxAmount {
		cypher += ` op.maxAmount = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $maxAmount ELSE op.maxAmount END, `
		params["maxAmount"] = data.MaxAmount
	}
	cypher += ` op.updatedAt = $updatedAt,
				op.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE op.sourceOfTruth END`
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *opportunityWriteRepository) ReplaceOwner(ctx context.Context, tenant, opportunityId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.ReplaceOwner")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	span.LogFields(log.String("userId", userId))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}), (op:Opportunity {id:$opportunityId}) WHERE op:Opportunity_%s
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
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *opportunityWriteRepository) CreateRenewal(ctx context.Context, tenant, opportunityId string, data RenewalOpportunityCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.CreateRenewal")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	tracing.LogObjectAsJson(span, "data", data)

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
								op.renewalLikelihood=$renewalLikelihood,
								op.renewalApproved=$renewalApproved
							WITH op, c
							MERGE (c)-[:ACTIVE_RENEWAL]->(op)
							`, tenant)
	params := map[string]any{
		"tenant":            tenant,
		"opportunityId":     opportunityId,
		"contractId":        data.ContractId,
		"createdAt":         data.CreatedAt,
		"updatedAt":         data.UpdatedAt,
		"source":            data.SourceFields.Source,
		"sourceOfTruth":     data.SourceFields.Source,
		"appSource":         data.SourceFields.AppSource,
		"internalType":      data.InternalType,
		"internalStage":     data.InternalStage,
		"renewalLikelihood": data.RenewalLikelihood,
		"renewalApproved":   data.RenewalApproved,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *opportunityWriteRepository) UpdateRenewal(ctx context.Context, tenant, opportunityId string, data RenewalOpportunityUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.UpdateRenewal")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	tracing.LogObjectAsJson(span, "data", data)

	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
		"updatedAt":     data.UpdatedAt,
		"sourceOfTruth": data.Source,
		"overwrite":     data.Source == constants.SourceOpenline,
	}
	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) WHERE op:RenewalOpportunity AND op:Opportunity_%s 
				SET op.updatedAt = $updatedAt,
					op.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE op.sourceOfTruth END`, tenant)
	if data.SetUpdatedByUserId {
		params["renewalUpdatedByUserId"] = data.UpdatedByUserId
		cypher += `, op.renewalUpdatedByUserAt = $updatedAt, 
					op.renewalUpdatedByUserId = $renewalUpdatedByUserId `
	}
	if data.UpdateComments {
		cypher += `, op.comments = $comments `
		params["comments"] = data.Comments
	}
	if data.UpdateAmount {
		cypher += `, op.amount = $amount `
		params["amount"] = data.Amount
	}
	if data.UpdateRenewalLikelihood {
		cypher += `, op.renewalLikelihood = $renewalLikelihood `
		params["renewalLikelihood"] = data.RenewalLikelihood
	}
	if data.UpdateRenewalApproved {
		cypher += `, op.renewalApproved = $renewalApproved `
		params["renewalApproved"] = data.RenewalApproved
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *opportunityWriteRepository) UpdateNextRenewalDate(ctx context.Context, tenant, opportunityId string, updatedAt time.Time, renewedAt *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.UpdateNextRenewalDate")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) 
							WHERE op:RenewalOpportunity AND op:Opportunity_%s AND op.internalStage=$internalStage
							SET op.updatedAt=$updatedAt, 
								op.renewedAt=$renewedAt`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
		"updatedAt":     updatedAt,
		"internalStage": enum.OpportunityInternalStageOpen.String(),
		"renewedAt":     utils.ToNeo4jDateAsAny(renewedAt),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *opportunityWriteRepository) CloseWin(ctx context.Context, tenant, opportunityId string, updatedAt, closedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.CloseWin")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) 
							WHERE op:Opportunity_%s AND op.internalStage <> $internalStage
							SET 
								op.closedAt=$closedAt, 
								op.internalStage=$internalStage,
								op.updatedAt=$updatedAt
							WITH op
							OPTIONAL MATCH (op)<-[rel:ACTIVE_RENEWAL]-(c:Contract)
							DELETE rel`, tenant)
	params := map[string]any{
		"opportunityId": opportunityId,
		"updatedAt":     updatedAt,
		"closedAt":      closedAt,
		"internalStage": enum.OpportunityInternalStageClosedWon.String(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *opportunityWriteRepository) CloseLoose(ctx context.Context, tenant, opportunityId string, updatedAt, closedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.CloseLoose")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) 
							WHERE op:Opportunity_%s AND op.internalStage <> $internalStage
							SET op.closedAt=$closedAt, 
								op.internalStage=$internalStage,
								op.updatedAt=$updatedAt
							WITH op
							OPTIONAL MATCH (op)<-[rel:ACTIVE_RENEWAL]-(c:Contract)
							DELETE rel`, tenant)
	params := map[string]any{
		"opportunityId": opportunityId,
		"updatedAt":     updatedAt,
		"closedAt":      closedAt,
		"internalStage": enum.OpportunityInternalStageClosedLost,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *opportunityWriteRepository) MarkRenewalRequested(ctx context.Context, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.MarkRenewalRequested")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId})
				WHERE op:Opportunity_%s
				SET op.techRolloutRenewalRequestedAt=$now`, tenant)
	params := map[string]any{
		"opportunityId": opportunityId,
		"now":           utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
