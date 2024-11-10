package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"time"
)

// Deprecated
type OpportunityCreateFields struct {
	OrganizationId    string             `json:"organizationId"`
	CreatedAt         time.Time          `json:"createdAt"`
	SourceFields      model.SourceFields `json:"sourceFields"`
	Name              string             `json:"name"`
	MaxAmount         float64            `json:"maxAmount"`
	InternalType      string             `json:"internalType"`
	ExternalType      string             `json:"externalType"`
	InternalStage     string             `json:"internalStage"`
	ExternalStage     string             `json:"externalStage"`
	EstimatedClosedAt *time.Time         `json:"estimatedClosedAt"`
	GeneralNotes      string             `json:"generalNotes"`
	NextSteps         string             `json:"nextSteps"`
	CreatedByUserId   string             `json:"createdByUserId"`
	Currency          enum.Currency      `json:"currency"`
	LikelihoodRate    int64              `json:"likelihoodRate"`
}

// Deprecated
type OpportunityUpdateFields struct {
	Source                  string        `json:"source"`
	Name                    string        `json:"name"`
	Amount                  float64       `json:"amount"`
	MaxAmount               float64       `json:"maxAmount"`
	ExternalStage           string        `json:"externalStage"`
	ExternalType            string        `json:"externalType"`
	EstimatedClosedAt       *time.Time    `json:"estimatedClosedAt"`
	InternalStage           string        `json:"internalStage"`
	Currency                enum.Currency `json:"currency"`
	NextSteps               string        `json:"nextSteps"`
	LikelihoodRate          int64         `json:"likelihoodRate"`
	UpdateName              bool          `json:"updateName"`
	UpdateAmount            bool          `json:"updateAmount"`
	UpdateMaxAmount         bool          `json:"updateMaxAmount"`
	UpdateExternalStage     bool          `json:"updateExternalStage"`
	UpdateExternalType      bool          `json:"updateExternalType"`
	UpdateEstimatedClosedAt bool          `json:"updateEstimatedClosedAt"`
	UpdateInternalStage     bool          `json:"updateInternalStage"`
	UpdateCurrency          bool          `json:"updateCurrency"`
	UpdateNextSteps         bool          `json:"updateNextSteps"`
	UpdateLikelihoodRate    bool          `json:"updateLikelihoodRate"`
}

type OpportunitySaveFields struct {
	AppSource         string        `json:"appSource"`
	Source            string        `json:"source"`
	Name              string        `json:"name"`
	Amount            float64       `json:"amount"`
	MaxAmount         float64       `json:"maxAmount"`
	ExternalStage     string        `json:"externalStage"`
	ExternalType      string        `json:"externalType"`
	EstimatedClosedAt *time.Time    `json:"estimatedClosedAt"`
	InternalStage     string        `json:"internalStage"`
	InternalType      string        `json:"internalType"`
	Currency          enum.Currency `json:"currency"`
	NextSteps         string        `json:"nextSteps"`
	LikelihoodRate    int64         `json:"likelihoodRate"`
	OwnerId           string        `json:"ownerId"`

	UpdateName              bool `json:"updateName"`
	UpdateAmount            bool `json:"updateAmount"`
	UpdateMaxAmount         bool `json:"updateMaxAmount"`
	UpdateExternalStage     bool `json:"updateExternalStage"`
	UpdateExternalType      bool `json:"updateExternalType"`
	UpdateEstimatedClosedAt bool `json:"updateEstimatedClosedAt"`
	UpdateInternalStage     bool `json:"updateInternalStage"`
	UpdateInternalType      bool `json:"updateInternalType"`
	UpdateCurrency          bool `json:"updateCurrency"`
	UpdateNextSteps         bool `json:"updateNextSteps"`
	UpdateLikelihoodRate    bool `json:"updateLikelihoodRate"`
	UpdateOwnerId           bool `json:"updateOwnerId"`
}

type RenewalOpportunityCreateFields struct {
	ContractId          string             `json:"contractId"`
	CreatedAt           time.Time          `json:"createdAt"`
	SourceFields        model.SourceFields `json:"sourceFields"`
	InternalType        string             `json:"internalType"`
	InternalStage       string             `json:"internalStage"`
	RenewalLikelihood   string             `json:"renewalLikelihood"`
	RenewalApproved     bool               `json:"renewalApproved"`
	RenewedAt           *time.Time         `json:"renewedAt"`
	RenewalAdjustedRate int64              `json:"renewalAdjustedRate"`
}

type RenewalOpportunityUpdateFields struct {
	UpdatedAt                 time.Time  `json:"updatedAt"`
	Source                    string     `json:"source"`
	UpdatedByUserId           string     `json:"updatedByUserId"`
	SetUpdatedByUserId        bool       `json:"setUpdatedByUserId"`
	Comments                  string     `json:"comments"`
	Amount                    float64    `json:"amount"`
	RenewalLikelihood         string     `json:"renewalLikelihood"`
	RenewalApproved           bool       `json:"renewalApproved"`
	RenewedAt                 *time.Time `json:"renewedAt"`
	RenewalAdjustedRate       int64      `json:"renewalAdjustedRate"`
	UpdateComments            bool       `json:"updateComments"`
	UpdateAmount              bool       `json:"updateAmount"`
	UpdateRenewalLikelihood   bool       `json:"updateRenewalLikelihood"`
	UpdateRenewalApproved     bool       `json:"updateRenewalApproved"`
	UpdateRenewedAt           bool       `json:"updateRenewedAt"`
	UpdateRenewalAdjustedRate bool       `json:"updateRenewalAdjustedRate"`
}

type OpportunityWriteRepository interface {
	//Deprecated
	CreateForOrganization(ctx context.Context, tenant, opportunityId string, data OpportunityCreateFields) error
	//Deprecated
	Update(ctx context.Context, tenant, opportunityId string, data OpportunityUpdateFields) error

	Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string, data OpportunitySaveFields) error
	ReplaceOwner(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId, userId string) error
	RemoveOwner(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) error
	CreateRenewal(ctx context.Context, tenant, opportunityId string, data RenewalOpportunityCreateFields) (bool, error)
	UpdateRenewal(ctx context.Context, tenant, opportunityId string, data RenewalOpportunityUpdateFields) error
	UpdateNextRenewalDate(ctx context.Context, tenant, opportunityId string, renewedAt *time.Time) error
	CloseWon(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string, closedAt time.Time) error
	CloseLost(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string, closedAt time.Time) error
	MarkRenewalRequested(ctx context.Context, tenant, opportunityId string) error
	Archive(ctx context.Context, tenant, opportunityId string) error
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.CreateForOrganizationOld")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$orgId})
							MERGE (t)<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity {id:$opportunityId})<-[:HAS_OPPORTUNITY]-(org)
							ON CREATE SET 
								op:Opportunity_%s,
								op.createdAt=$createdAt,
								op.updatedAt=datetime(),
								op.stageUpdatedAt=datetime(),
								op.source=$source,
								op.sourceOfTruth=$sourceOfTruth,
								op.appSource=$appSource,
								op.name=$name,
								op.maxAmount=$maxAmount,
								op.internalType=$internalType,
								op.externalType=$externalType,
								op.internalStage=$internalStage,
								op.externalStage=$externalStage,
								op.estimatedClosedAt=$estimatedClosedAt,
								op.generalNotes=$generalNotes,
								op.nextSteps=$nextSteps,
								op.currency=$currency,
								op.likelihoodRate=$likelihoodRate,
								org.updatedAt=datetime()
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
		"source":            data.SourceFields.Source,
		"sourceOfTruth":     data.SourceFields.Source,
		"appSource":         data.SourceFields.AppSource,
		"name":              data.Name,
		"maxAmount":         data.MaxAmount,
		"internalType":      data.InternalType,
		"externalType":      data.ExternalType,
		"internalStage":     data.InternalStage,
		"externalStage":     data.ExternalStage,
		"estimatedClosedAt": utils.TimePtrAsAny(data.EstimatedClosedAt),
		"generalNotes":      data.GeneralNotes,
		"nextSteps":         data.NextSteps,
		"createdByUserId":   data.CreatedByUserId,
		"currency":          data.Currency.String(),
		"likelihoodRate":    data.LikelihoodRate,
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	tracing.LogObjectAsJson(span, "data", data)

	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
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
	if data.UpdateExternalType {
		cypher += ` op.externalType = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $externalType ELSE op.externalType END, `
		params["externalType"] = data.ExternalType
	}
	if data.UpdateExternalStage {
		cypher += ` op.externalStage = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $externalStage ELSE op.externalStage END, `
		params["externalStage"] = data.ExternalStage
	}
	if data.UpdateEstimatedClosedAt {
		cypher += ` op.estimatedClosedAt = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $estimatedClosedAt ELSE op.estimatedClosedAt END, `
		params["estimatedClosedAt"] = utils.TimePtrAsAny(data.EstimatedClosedAt)
	}
	if data.UpdateInternalStage {
		cypher += ` op.internalStage = $internalStage, `
		params["internalStage"] = data.InternalStage
	}
	if data.UpdateCurrency {
		cypher += ` op.currency = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $currency ELSE op.currency END, `
		params["currency"] = data.Currency.String()
	}
	if data.UpdateNextSteps {
		cypher += ` op.nextSteps = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $nextSteps ELSE op.nextSteps END, `
		params["nextSteps"] = data.NextSteps
	}
	if data.UpdateLikelihoodRate {
		cypher += ` op.likelihoodRate = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $likelihoodRate ELSE op.likelihoodRate END, `
		params["likelihoodRate"] = data.LikelihoodRate
	}
	cypher += ` op.updatedAt = datetime(),
				op.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE op.sourceOfTruth END`
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *opportunityWriteRepository) Save(ctx context.Context, txx *neo4j.ManagedTransaction, tenant, opportunityId string, data OpportunitySaveFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.Save")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	tracing.LogObjectAsJson(span, "data", data)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, txx, func(tx neo4j.ManagedTransaction) (any, error) {

		//create if not exists
		cypherCreate := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) MERGE(t)<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity:Opportunity_%s {id:$opportunityId})
			ON CREATE SET
				op.createdAt=datetime(),
				op.updatedAt=datetime(),
				op.stageUpdatedAt=datetime(),
				op.source=$source,
				op.appSource=$appSource,
			`, tenant)
		paramsCreate := map[string]any{
			"tenant":        tenant,
			"opportunityId": opportunityId,
			"appSource":     data.AppSource,
			"source":        data.Source,
		}

		span.LogFields(log.String("cypherCreate", cypherCreate))
		tracing.LogObjectAsJson(span, "paramsCreate", paramsCreate)

		_, err := tx.Run(ctx, cypherCreate, paramsCreate)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		paramsUpdate := map[string]any{
			"tenant":        tenant,
			"opportunityId": opportunityId,
			"sourceOfTruth": data.Source,
			"overwrite":     data.Source == constants.SourceOpenline,
		}

		cypherUpdate := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity:Opportunity_%s {id:$opportunityId}) SET `, tenant)
		if data.UpdateName {
			cypherUpdate += ` op.name = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR op.name = '' THEN $name ELSE op.name END, `
			paramsUpdate["name"] = data.Name
		}
		if data.UpdateAmount {
			cypherUpdate += ` op.amount = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $amount ELSE op.amount END, `
			paramsUpdate["amount"] = data.Amount
		}
		if data.UpdateMaxAmount {
			cypherUpdate += ` op.maxAmount = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $maxAmount ELSE op.maxAmount END, `
			paramsUpdate["maxAmount"] = data.MaxAmount
		}
		if data.UpdateExternalType {
			cypherUpdate += ` op.externalType = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $externalType ELSE op.externalType END, `
			paramsUpdate["externalType"] = data.ExternalType
		}
		if data.UpdateExternalStage {
			cypherUpdate += ` op.externalStage = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $externalStage ELSE op.externalStage END, `
			paramsUpdate["externalStage"] = data.ExternalStage
		}
		if data.UpdateEstimatedClosedAt {
			cypherUpdate += ` op.estimatedClosedAt = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $estimatedClosedAt ELSE op.estimatedClosedAt END, `
			paramsUpdate["estimatedClosedAt"] = utils.TimePtrAsAny(data.EstimatedClosedAt)
		}
		if data.UpdateInternalStage {
			cypherUpdate += ` op.internalStage = $internalStage, `
			paramsUpdate["internalStage"] = data.InternalStage
		}
		if data.UpdateInternalType {
			cypherUpdate += ` op.internalType = $internalType, `
			paramsUpdate["internalType"] = data.InternalType
		}
		if data.UpdateCurrency {
			cypherUpdate += ` op.currency = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $currency ELSE op.currency END, `
			paramsUpdate["currency"] = data.Currency.String()
		}
		if data.UpdateNextSteps {
			cypherUpdate += ` op.nextSteps = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $nextSteps ELSE op.nextSteps END, `
			paramsUpdate["nextSteps"] = data.NextSteps
		}
		if data.UpdateLikelihoodRate {
			cypherUpdate += ` op.likelihoodRate = CASE WHEN op.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $likelihoodRate ELSE op.likelihoodRate END, `
			paramsUpdate["likelihoodRate"] = data.LikelihoodRate
		}
		cypherUpdate += ` op.updatedAt = datetime(),
				op.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE op.sourceOfTruth END`

		span.LogFields(log.String("cypherUpdate", cypherUpdate))
		tracing.LogObjectAsJson(span, "paramsUpdate", paramsUpdate)

		_, err = tx.Run(ctx, cypherUpdate, paramsUpdate)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (r *opportunityWriteRepository) ReplaceOwner(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.ReplaceOwner")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	span.LogFields(log.String("userId", userId))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}), (op:Opportunity:Opportunity_%s {id:$opportunityId})
			WITH op, t
			OPTIONAL MATCH (:User_%s)-[rel:OWNS]->(op)
			DELETE rel
			WITH op, t
			MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User_%s {id:$userId})
			WHERE (u.internal=false OR u.internal is null) AND (u.bot=false OR u.bot is null)
			MERGE (u)-[:OWNS]->(op)
			SET op.updatedAt=datetime()`, tenant, tenant, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
		"userId":        userId,
		"now":           utils.Now(),
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *opportunityWriteRepository) RemoveOwner(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.RemoveOwner")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity:Opportunity_%s {id:$opportunityId})<-[rel:OWNS]-(:User_%s)
				SET op.updatedAt=datetime()
				DELETE rel`, tenant, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *opportunityWriteRepository) CreateRenewal(ctx context.Context, tenant, opportunityId string, data RenewalOpportunityCreateFields) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.CreateRenewal")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
							WHERE NOT (c)-[:ACTIVE_RENEWAL]->(:RenewalOpportunity)
							MERGE (c)-[:ACTIVE_RENEWAL]->(newOp:Opportunity {id:$opportunityId})
							ON CREATE SET 
								newOp:Opportunity_%s,
								newOp:RenewalOpportunity,
								newOp.createdAt=$createdAt,
								newOp.updatedAt=datetime(),
								newOp.source=$source,
								newOp.sourceOfTruth=$sourceOfTruth,
								newOp.appSource=$appSource,
								newOp.internalType=$internalType,
								newOp.internalStage=$internalStage,
								newOp.renewalLikelihood=$renewalLikelihood,
								newOp.renewalApproved=$renewalApproved,
								newOp.renewedAt=$renewedAt,
								newOp.renewalAdjustedRate=$renewalAdjustedRate
							WITH c, newOp
								MERGE (c)-[:HAS_OPPORTUNITY]->(newOp)
							RETURN count(newOp) > 0 AS created`, tenant)
	params := map[string]any{
		"tenant":              tenant,
		"opportunityId":       opportunityId,
		"contractId":          data.ContractId,
		"createdAt":           data.CreatedAt,
		"source":              data.SourceFields.Source,
		"sourceOfTruth":       data.SourceFields.Source,
		"appSource":           data.SourceFields.AppSource,
		"internalType":        data.InternalType,
		"internalStage":       data.InternalStage,
		"renewalLikelihood":   data.RenewalLikelihood,
		"renewalApproved":     data.RenewalApproved,
		"renewalAdjustedRate": data.RenewalAdjustedRate,
		"renewedAt":           utils.ToDateAsAny(data.RenewedAt),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.Bool("result.created", false))
		return false, err
	}
	span.LogFields(log.Bool("result.created", result.(bool)))
	return result.(bool), nil
}

func (r *opportunityWriteRepository) UpdateRenewal(ctx context.Context, tenant, opportunityId string, data RenewalOpportunityUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.UpdateRenewal")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
				SET op.updatedAt = datetime(),
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
	if data.UpdateRenewedAt {
		cypher += `, op.renewedAt = $renewedAt `
		params["renewedAt"] = utils.ToDateAsAny(data.RenewedAt)
	}
	if data.UpdateRenewalAdjustedRate {
		cypher += `, op.renewalAdjustedRate = $renewalAdjustedRate `
		params["renewalAdjustedRate"] = data.RenewalAdjustedRate
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *opportunityWriteRepository) UpdateNextRenewalDate(ctx context.Context, tenant, opportunityId string, renewedAt *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.UpdateNextRenewalDate")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) 
							WHERE op:RenewalOpportunity AND op:Opportunity_%s AND op.internalStage=$internalStage
							SET op.updatedAt=datetime(), 
								op.renewedAt=$renewedAt`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
		"internalStage": enum.OpportunityInternalStageOpen.String(),
		"renewedAt":     utils.ToDateAsAny(renewedAt),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *opportunityWriteRepository) CloseWon(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string, closedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.CloseWon")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) 
							WHERE op:Opportunity_%s AND op.internalStage <> $internalStage
							SET 
								op.closedAt=$closedAt, 
								op.internalStage=$internalStage,
								op.updatedAt=datetime(),
								op.stageUpdatedAt=datetime()
							WITH op
							OPTIONAL MATCH (op)<-[rel:ACTIVE_RENEWAL]-(c:Contract)
							DELETE rel`, tenant)
	params := map[string]any{
		"opportunityId": opportunityId,
		"closedAt":      closedAt,
		"internalStage": enum.OpportunityInternalStageClosedWon.String(),
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *opportunityWriteRepository) CloseLost(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string, closedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.CloseLost")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) 
							WHERE op:Opportunity_%s AND op.internalStage <> $internalStage
							SET op.closedAt=$closedAt, 
								op.internalStage=$internalStage,
								op.updatedAt=datetime(),
								op.stageUpdatedAt=datetime()
							WITH op
							OPTIONAL MATCH (op)<-[rel:ACTIVE_RENEWAL]-(c:Contract)
							DELETE rel`, tenant)
	params := map[string]any{
		"opportunityId": opportunityId,
		"closedAt":      closedAt,
		"internalStage": enum.OpportunityInternalStageClosedLost,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *opportunityWriteRepository) MarkRenewalRequested(ctx context.Context, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.MarkRenewalRequested")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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

func (r *opportunityWriteRepository) Archive(ctx context.Context, tenant, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.Archive")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity {id:$opportunityId}) 
							WHERE op:Opportunity_%s
							SET op.updatedAt=datetime(),
								op:ArchivedOpportunity,
								op:ArchivedOpportunity_%s
							REMOVE op:Opportunity, op:Opportunity_%s
							`, tenant, tenant, tenant)
	params := map[string]any{
		"opportunityId": opportunityId,
		"tenant":        tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
