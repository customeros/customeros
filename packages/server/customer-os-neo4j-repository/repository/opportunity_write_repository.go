package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"time"
)

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
	Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string, data data_fields.OpportunityFields) error
	ReplaceOwner(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId, userId string) error
	RemoveOwner(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) error
	CreateRenewal(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string, data data_fields.OpportunityFields) (bool, error)
	UpdateRenewal(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string, data data_fields.OpportunityFields) error
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

func (r *opportunityWriteRepository) Save(ctx context.Context, txx *neo4j.ManagedTransaction, tenant, opportunityId string, data data_fields.OpportunityFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.Save")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, opportunityId)
	tracing.LogObjectAsJson(span, "data", data)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, txx, func(tx neo4j.ManagedTransaction) (any, error) {

		//create if not exists
		cypherCreate := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) MERGE(t)<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity:Opportunity_%s {id:$opportunityId})
			ON CREATE SET
				op.createdAt=datetime(),
				op.updatedAt=datetime(),
				op.stageUpdatedAt=datetime(),
				op.source=$source,
				op.appSource=$appSource
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
		}

		cypherUpdate := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity:Opportunity_%s {id:$opportunityId}) 
			SET op.updatedAt = datetime()`, tenant)
		if data.Name != nil {
			cypherUpdate += `, op.name = $name `
			paramsUpdate["name"] = *data.Name
		}
		if data.Amount != nil {
			cypherUpdate += `, op.amount = $amount `
			paramsUpdate["amount"] = *data.Amount
		}
		if data.MaxAmount != nil {
			cypherUpdate += `, op.maxAmount = $maxAmount `
			paramsUpdate["maxAmount"] = *data.MaxAmount
		}
		if data.ExternalStage != nil {
			cypherUpdate += `, op.externalStage = $externalStage `
			paramsUpdate["externalStage"] = *data.ExternalStage
		}
		if data.ExternalType != nil {
			cypherUpdate += `, op.externalType = $externalType `
			paramsUpdate["externalType"] = *data.ExternalType
		}
		if data.EstimatedClosedAt != nil {
			cypherUpdate += `, op.estimatedClosedAt = $estimatedClosedAt `
			paramsUpdate["estimatedClosedAt"] = utils.TimePtrAsAny(data.EstimatedClosedAt)
		}
		if data.InternalStage != nil {
			cypherUpdate += `, op.internalStage = $internalStage `
			paramsUpdate["internalStage"] = *data.InternalStage
		}
		if data.InternalType != nil {
			cypherUpdate += `, op.internalType = $internalType `
			paramsUpdate["internalType"] = data.InternalType.String()
		}
		if data.Currency != nil {
			cypherUpdate += `, op.currency = $currency `
			paramsUpdate["currency"] = data.Currency.String()
		}
		if data.NextSteps != nil {
			cypherUpdate += `, op.nextSteps = $nextSteps `
			paramsUpdate["nextSteps"] = *data.NextSteps
		}
		if data.LikelihoodRate != nil {
			cypherUpdate += `, op.likelihoodRate = $likelihoodRate `
			paramsUpdate["likelihoodRate"] = *data.LikelihoodRate
		}

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

func (r *opportunityWriteRepository) CreateRenewal(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string, data data_fields.OpportunityFields) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.CreateRenewal")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	tracing.LogObjectAsJson(span, "data", data)

	if !data.IsRenewal() {
		err := fmt.Errorf("opportunity is not a renewal opportunity")
		tracing.TraceErr(span, err)
		return false, err
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
							WHERE NOT (c)-[:ACTIVE_RENEWAL]->(:RenewalOpportunity)
							MERGE (c)-[:ACTIVE_RENEWAL]->(newOp:Opportunity {id:$opportunityId})
							ON CREATE SET 
								newOp:Opportunity_%s,
								newOp:RenewalOpportunity,
								newOp.createdAt=$createdAt,
								newOp.updatedAt=datetime(),
								newOp.source=$source,
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
		"contractId":          utils.IfNotNilString(data.ContractId),
		"createdAt":           utils.IfNotNilTimeWithDefault(data.CreatedAt, utils.Now()),
		"source":              utils.IfNotNilString(data.Source),
		"appSource":           utils.IfNotNilString(data.AppSource),
		"internalStage":       utils.IfNotNilString(data.InternalStage),
		"renewalApproved":     utils.IfNotNilBool(data.RenewalApproved),
		"renewalAdjustedRate": utils.IfNotNilInt64(data.RenewalAdjustedRate),
		"renewedAt":           utils.ToDateAsAny(data.RenewedAt),
	}
	if data.InternalType != nil {
		params["internalType"] = data.InternalType.String()
	} else {
		params["internalType"] = ""
	}
	if data.RenewalLikelihood != nil {
		params["renewalLikelihood"] = data.RenewalLikelihood.String()
	} else {
		params["renewalLikelihood"] = ""
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.Bool("result.created", false))
		return false, err
	}
	span.LogFields(log.Bool("result.created", result.(bool)))
	return result.(bool), nil
}

func (r *opportunityWriteRepository) UpdateRenewal(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string, data data_fields.OpportunityFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityWriteRepository.UpdateRenewal")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	tracing.LogObjectAsJson(span, "data", data)

	params := map[string]any{
		"tenant":        tenant,
		"opportunityId": opportunityId,
	}
	cypher := fmt.Sprintf(`MATCH (op:Opportunity {id:$opportunityId}) WHERE op:RenewalOpportunity AND op:Opportunity_%s 
				SET op.updatedAt = datetime()`, tenant)
	if data.Comments != nil {
		cypher += `, op.comments = $comments `
		params["comments"] = *data.Comments
	}
	if data.Amount != nil {
		cypher += `, op.amount = $amount `
		params["amount"] = *data.Amount
	}
	if data.RenewalLikelihood != nil {
		cypher += `, op.renewalLikelihood = $renewalLikelihood `
		params["renewalLikelihood"] = data.RenewalLikelihood.String()
	}
	if data.RenewalApproved != nil {
		cypher += `, op.renewalApproved = $renewalApproved `
		params["renewalApproved"] = *data.RenewalApproved
	}
	if data.RenewedAt != nil {
		cypher += `, op.renewedAt = $renewedAt `
		params["renewedAt"] = utils.ToDateAsAny(data.RenewedAt)
	}
	if data.RenewalAdjustedRate != nil {
		cypher += `, op.renewalAdjustedRate = $renewalAdjustedRate `
		params["renewalAdjustedRate"] = data.RenewalAdjustedRate
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
