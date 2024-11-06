package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ContractWriteRepository interface {
	CreateForOrganization(ctx context.Context, tenant, contractId string, data data_fields.ContractSaveFields) error
	UpdateContract(ctx context.Context, tenant, contractId string, data data_fields.ContractSaveFields) error
	UpdateStatus(ctx context.Context, tenant, contractId, status string) error
	SuspendActiveRenewalOpportunity(ctx context.Context, tenant, contractId string) error
	ActivateSuspendedRenewalOpportunity(ctx context.Context, tenant, contractId string) error
	ContractCausedOnboardingStatusChange(ctx context.Context, tenant, contractId string) error
	MarkStatusRenewalRequested(ctx context.Context, tenant, contractId string) error
	MarkRolloutRenewalRequested(ctx context.Context, tenant, contractId string) error
	MarkCycleInvoicingRequested(ctx context.Context, tenant, contractId string, invoicingStartedAt time.Time) error
	MarkOffCycleInvoicingRequested(ctx context.Context, tenant, contractId string, invoicingStartedAt time.Time) error
	MarkNextPreviewInvoicingRequested(ctx context.Context, tenant, contractId string, invoicingStartedAt time.Time) error
	SoftDelete(ctx context.Context, tenant, contractId string, deletedAt time.Time) error
	SetLtv(ctx context.Context, tenant, contractId string, ltv float64) error
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

func (r *contractWriteRepository) CreateForOrganization(ctx context.Context, tenant, contractId string, data data_fields.ContractSaveFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.CreateForOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$orgId})
							MERGE (t)<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})<-[:HAS_CONTRACT]-(org)
							ON CREATE SET 
								ct:Contract_%s,
								ct.createdAt=$createdAt,
								ct.updatedAt=datetime(),
								ct.source=$source,
								ct.appSource=$appSource,
								ct.name=$name,
								ct.contractUrl=$contractUrl,
								ct.status=$status,
								ct.signedAt=$signedAt,
								ct.serviceStartedAt=$serviceStartedAt,
								ct.currency=$currency,
								ct.billingCycleInMonths=$billingCycleInMonths,
								ct.invoicingStartDate=$invoicingStartDate,
								ct.invoicingEnabled=$invoicingEnabled,
								ct.payOnline=$payOnline,
								ct.payAutomatically=$payAutomatically,
								ct.canPayWithCard=$canPayWithCard,
								ct.canPayWithDirectDebit=$canPayWithDirectDebit,
								ct.canPayWithBankTransfer=$canPayWithBankTransfer,
								ct.autoRenew=$autoRenew,
								ct.check=$check,
								ct.country=$country,
								ct.dueDays=$dueDays,
								ct.lengthInMonths=$lengthInMonths,
								ct.approved=$approved,
								org.updatedAt=datetime()
							WITH ct, t
							OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$createdByUserId}) 
							WHERE $createdByUserId <> ""
							FOREACH (ignore IN CASE WHEN u IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (ct)-[:CREATED_BY]->(u))
							`, tenant)
	currency := neo4jenum.CurrencyUSD.String()
	if data.Currency != nil {
		currency = data.Currency.String()
	}
	params := map[string]any{
		"tenant":                 tenant,
		"contractId":             contractId,
		"orgId":                  data.GetOrganizationId(),
		"createdAt":              utils.IfNotNilTimeWithDefault(data.CreatedAt, utils.Now()),
		"source":                 data.Source,
		"appSource":              data.AppSource,
		"name":                   utils.IfNotNilString(data.Name),
		"contractUrl":            utils.IfNotNilString(data.ContractUrl),
		"status":                 utils.IfNotNilString(data.Status),
		"signedAt":               utils.ToDateAsAny(data.SignedAt),
		"serviceStartedAt":       utils.ToDateAsAny(data.ServiceStartedAt),
		"createdByUserId":        data.GetCreatedByUserId(),
		"currency":               currency,
		"billingCycleInMonths":   utils.IfNotNilInt64(data.BillingCycleInMonths),
		"invoicingStartDate":     utils.ToNeo4jDateAsAny(data.InvoicingStartDate),
		"invoicingEnabled":       utils.IfNotNilBool(data.InvoicingEnabled),
		"payOnline":              utils.IfNotNilBool(data.PayOnline),
		"payAutomatically":       utils.IfNotNilBool(data.PayAutomatically),
		"canPayWithCard":         utils.IfNotNilBool(data.CanPayWithCard),
		"canPayWithDirectDebit":  utils.IfNotNilBool(data.CanPayWithDirectDebit),
		"canPayWithBankTransfer": utils.IfNotNilBool(data.CanPayWithBankTransfer),
		"autoRenew":              utils.IfNotNilBool(data.AutoRenew),
		"check":                  utils.IfNotNilBool(data.Check),
		"dueDays":                utils.IfNotNilInt64(data.DueDays),
		"country":                utils.IfNotNilString(data.Country),
		"lengthInMonths":         utils.IfNotNilInt64(data.LengthInMonths),
		"approved":               utils.IfNotNilBool(data.Approved),
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) UpdateContract(ctx context.Context, tenant, contractId string, data data_fields.ContractSaveFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.UpdateContract")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
				SET 
				ct.updatedAt = datetime()
				`
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
	}
	if data.Name != nil {
		cypher += `, ct.name =  $name `
		params["name"] = *data.Name
	}
	if data.ContractUrl != nil {
		cypher += `, ct.contractUrl = $contractUrl`
		params["contractUrl"] = *data.ContractUrl
	}
	if data.Status != nil {
		cypher += `, ct.status = $status`
		params["status"] = *data.Status
	}
	if data.ServiceStartedAt != nil {
		cypher += `, ct.serviceStartedAt = $serviceStartedAt`
		params["serviceStartedAt"] = utils.ToDateAsAny(data.ServiceStartedAt)
	}
	if data.SignedAt != nil {
		cypher += `, ct.signedAt = $signedAt`
		params["signedAt"] = utils.ToDateAsAny(data.SignedAt)
	}
	if data.EndedAt != nil {
		cypher += `, ct.endedAt = $endedAt`
		params["endedAt"] = utils.ToDateAsAny(data.EndedAt)
	}
	if data.BillingCycleInMonths != nil {
		cypher += `, ct.billingCycleInMonths = $billingCycleInMonths`
		params["billingCycleInMonths"] = utils.IfNotNilInt64(data.BillingCycleInMonths)
	}
	if data.Currency != nil {
		cypher += `, ct.currency = $currency`
		params["currency"] = data.Currency.String()
	}
	if data.InvoicingStartDate != nil {
		cypher += `, ct.invoicingStartDate = $invoicingStartDate`
		params["invoicingStartDate"] = utils.ToNeo4jDateAsAny(data.InvoicingStartDate)
	}
	if data.AddressLine1 != nil {
		cypher += `, ct.addressLine1 = $addressLine1`
		params["addressLine1"] = *data.AddressLine1
	}
	if data.AddressLine2 != nil {
		cypher += `, ct.addressLine2 = $addressLine2`
		params["addressLine2"] = *data.AddressLine2
	}
	if data.Locality != nil {
		cypher += `, ct.locality = $locality`
		params["locality"] = *data.Locality
	}
	if data.Country != nil {
		cypher += `, ct.country = $country`
		params["country"] = *data.Country
	}
	if data.Region != nil {
		cypher += `, ct.region = $region`
		params["region"] = *data.Region
	}
	if data.Zip != nil {
		cypher += `, ct.zip = $zip`
		params["zip"] = *data.Zip
	}
	if data.OrganizationLegalName != nil {
		cypher += `, ct.organizationLegalName = $organizationLegalName`
		params["organizationLegalName"] = *data.OrganizationLegalName
	}
	if data.InvoiceEmail != nil {
		cypher += `, ct.invoiceEmail = $invoiceEmail`
		params["invoiceEmail"] = *data.InvoiceEmail
	}
	if data.InvoiceEmailCC != nil {
		cypher += `, ct.invoiceEmailCC = $invoiceEmailCC`
		params["invoiceEmailCC"] = *data.InvoiceEmailCC
	}
	if data.InvoiceEmailBCC != nil {
		cypher += `, ct.invoiceEmailBCC = $invoiceEmailBCC`
		params["invoiceEmailBCC"] = *data.InvoiceEmailBCC
	}
	if data.InvoiceNote != nil {
		cypher += `, ct.invoiceNote = $invoiceNote`
		params["invoiceNote"] = *data.InvoiceNote
	}
	if data.NextInvoiceDate != nil {
		cypher += `, ct.nextInvoiceDate=$nextInvoiceDate `
		params["nextInvoiceDate"] = utils.ToNeo4jDateAsAny(data.NextInvoiceDate)
	}
	if data.CanPayWithCard != nil {
		cypher += `, ct.canPayWithCard=$canPayWithCard `
		params["canPayWithCard"] = *data.CanPayWithCard
	}
	if data.CanPayWithDirectDebit != nil {
		cypher += `, ct.canPayWithDirectDebit=$canPayWithDirectDebit `
		params["canPayWithDirectDebit"] = *data.CanPayWithDirectDebit
	}
	if data.CanPayWithBankTransfer != nil {
		cypher += `, ct.canPayWithBankTransfer=$canPayWithBankTransfer `
		params["canPayWithBankTransfer"] = *data.CanPayWithBankTransfer
	}
	if data.InvoicingEnabled != nil {
		cypher += `, ct.invoicingEnabled=$invoicingEnabled `
		params["invoicingEnabled"] = *data.InvoicingEnabled
	}
	if data.PayOnline != nil {
		cypher += `, ct.payOnline=$payOnline `
		params["payOnline"] = *data.PayOnline
	}
	if data.PayAutomatically != nil {
		cypher += `, ct.payAutomatically=$payAutomatically `
		params["payAutomatically"] = *data.PayAutomatically
	}
	if data.AutoRenew != nil {
		cypher += `, ct.autoRenew=$autoRenew `
		params["autoRenew"] = *data.AutoRenew
	}
	if data.Check != nil {
		cypher += `, ct.check=$check `
		params["check"] = *data.Check
	}
	if data.DueDays != nil {
		cypher += `, ct.dueDays=$dueDays `
		params["dueDays"] = *data.DueDays
	}
	if data.LengthInMonths != nil {
		cypher += `, ct.lengthInMonths=$lengthInMonths `
		params["lengthInMonths"] = *data.LengthInMonths
	}
	if data.Approved != nil {
		cypher += `, ct.approved=$approved `
		params["approved"] = *data.Approved
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) UpdateStatus(ctx context.Context, tenant, contractId, status string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.UpdateStatus")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
				SET 
					ct.status=$status,
					ct.updatedAt=datetime()
							`
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
		"status":     status,
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})-[r:ACTIVE_RENEWAL]->(op:RenewalOpportunity)
				SET op.internalStage=$internalStageSuspended, 
					op.updatedAt=datetime()
				MERGE (ct)-[:SUSPENDED_RENEWAL]->(op)
				DELETE r`
	params := map[string]any{
		"tenant":                 tenant,
		"contractId":             contractId,
		"internalStageSuspended": "SUSPENDED",
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})-[r:SUSPENDED_RENEWAL]->(op:RenewalOpportunity)
				SET op.internalStage=$internalStage, 
					op.updatedAt=datetime()
				MERGE (ct)-[:ACTIVE_RENEWAL]->(op)
				DELETE r`
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    contractId,
		"internalStage": neo4jenum.OpportunityInternalStageOpen.String(),
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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

func (r *contractWriteRepository) MarkStatusRenewalRequested(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.MarkStatusRenewalRequested")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
				SET ct.techStatusRenewalRequestedAt=$now`
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
		"now":        utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) MarkRolloutRenewalRequested(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.MarkRolloutRenewalRequested")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
				SET ct.techRolloutRenewalRequestedAt=$now`
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
		"now":        utils.Now(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) MarkCycleInvoicingRequested(ctx context.Context, tenant, contractId string, invoicingStartedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.MarkCycleInvoicingRequested")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
				SET c.techInvoicingStartedAt=$invoicingStartedAt`
	params := map[string]any{
		"tenant":             tenant,
		"contractId":         contractId,
		"invoicingStartedAt": invoicingStartedAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) MarkOffCycleInvoicingRequested(ctx context.Context, tenant, contractId string, invoicingStartedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.MarkOffCycleInvoicingRequested")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
				SET c.techOffCycleInvoicingStartedAt=$invoicingStartedAt`
	params := map[string]any{
		"tenant":             tenant,
		"contractId":         contractId,
		"invoicingStartedAt": invoicingStartedAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) MarkNextPreviewInvoicingRequested(ctx context.Context, tenant, contractId string, nextPreviewInvoiceRequestedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.MarkNextPreviewInvoicingRequested")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
				SET c.techNextPreviewInvoiceRequestedAt=$nextPreviewInvoiceRequestedAt`
	params := map[string]any{
		"tenant":                        tenant,
		"contractId":                    contractId,
		"nextPreviewInvoiceRequestedAt": nextPreviewInvoiceRequestedAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) SoftDelete(ctx context.Context, tenant, contractId string, deletedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.SoftDelete")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
			SET ct.updatedAt=$deletedAt,
				ct:%s,
				ct:%s
			REMOVE 	ct:%s, 
					ct:%s`,
		model2.NodeLabelDeletedContract, model2.NodeLabelDeletedContract+"_"+tenant,
		model2.NodeLabelContract, model2.NodeLabelContract+"_"+tenant)
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
		"deletedAt":  deletedAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return tx.Run(ctx, cypher, params)
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *contractWriteRepository) SetLtv(ctx context.Context, tenant, contractId string, ltv float64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.SetLtv")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
				SET ct.ltv=$ltv`
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
		"ltv":        ltv,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
