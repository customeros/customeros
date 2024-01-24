package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ContractCreateFields struct {
	OrganizationId     string                 `json:"organizationId"`
	Name               string                 `json:"name"`
	ContractUrl        string                 `json:"contractUrl"`
	CreatedByUserId    string                 `json:"createdByUserId"`
	ServiceStartedAt   *time.Time             `json:"serviceStartedAt,omitempty"`
	SignedAt           *time.Time             `json:"signedAt,omitempty"`
	RenewalCycle       string                 `json:"renewalCycle"`
	RenewalPeriods     *int64                 `json:"renewalPeriods,omitempty"`
	Status             string                 `json:"status"`
	CreatedAt          time.Time              `json:"createdAt"`
	UpdatedAt          time.Time              `json:"updatedAt"`
	SourceFields       model.Source           `json:"sourceFields"`
	BillingCycle       neo4jenum.BillingCycle `json:"billingCycle"`
	Currency           neo4jenum.Currency     `json:"currency"`
	InvoicingStartDate *time.Time             `json:"invoicingStartDate,omitempty"`
}

type ContractUpdateFields struct {
	Name                        string                 `json:"name"`
	ContractUrl                 string                 `json:"contractUrl"`
	Status                      string                 `json:"status"`
	Source                      string                 `json:"source"`
	RenewalPeriods              *int64                 `json:"renewalPeriods"`
	RenewalCycle                string                 `json:"renewalCycle"`
	UpdatedAt                   time.Time              `json:"updatedAt"`
	ServiceStartedAt            *time.Time             `json:"serviceStartedAt"`
	SignedAt                    *time.Time             `json:"signedAt"`
	EndedAt                     *time.Time             `json:"endedAt"`
	BillingCycle                neo4jenum.BillingCycle `json:"billingCycle"`
	Currency                    neo4jenum.Currency     `json:"currency"`
	InvoicingStartDate          *time.Time             `json:"invoicingStartDate,omitempty"`
	AddressLine1                string                 `json:"addressLine1"`
	AddressLine2                string                 `json:"addressLine2"`
	Locality                    string                 `json:"locality"`
	Country                     string                 `json:"country"`
	Zip                         string                 `json:"zip"`
	OrganizationLegalName       string                 `json:"organizationLegalName"`
	InvoiceEmail                string                 `json:"invoiceEmail"`
	InvoiceNote                 string                 `json:"invoiceNote"`
	UpdateName                  bool                   `json:"updateName"`
	UpdateContractUrl           bool                   `json:"updateContractUrl"`
	UpdateStatus                bool                   `json:"updateStatus"`
	UpdateRenewalPeriods        bool                   `json:"updateRenewalPeriods"`
	UpdateRenewalCycle          bool                   `json:"updateRenewalCycle"`
	UpdateServiceStartedAt      bool                   `json:"updateServiceStartedAt"`
	UpdateSignedAt              bool                   `json:"updateSignedAt"`
	UpdateEndedAt               bool                   `json:"updateEndedAt"`
	UpdateBillingCycle          bool                   `json:"updateBillingCycle"`
	UpdateCurrency              bool                   `json:"updateCurrency"`
	UpdateInvoicingStartDate    bool                   `json:"updateInvoicingStartDate"`
	UpdateAddressLine1          bool                   `json:"updateAddressLine1"`
	UpdateAddressLine2          bool                   `json:"updateAddressLine2"`
	UpdateLocality              bool                   `json:"updateLocality"`
	UpdateCountry               bool                   `json:"updateCountry"`
	UpdateZip                   bool                   `json:"updateZip"`
	UpdateOrganizationLegalName bool                   `json:"updateOrganizationLegalName"`
	UpdateInvoiceEmail          bool                   `json:"updateInvoiceEmail"`
	UpdateInvoiceNote           bool                   `json:"updateInvoiceNote"`
}

type ContractWriteRepository interface {
	CreateForOrganization(ctx context.Context, tenant, contractId string, data ContractCreateFields) error
	UpdateAndReturn(ctx context.Context, tenant, contractId string, data ContractUpdateFields) (*dbtype.Node, error)
	UpdateStatus(ctx context.Context, tenant, contractId, status string, serviceStartedAt, endedAt *time.Time) error
	SuspendActiveRenewalOpportunity(ctx context.Context, tenant, contractId string) error
	ActivateSuspendedRenewalOpportunity(ctx context.Context, tenant, contractId string) error
	ContractCausedOnboardingStatusChange(ctx context.Context, tenant, contractId string) error
	MarkStatusRenewalRequested(ctx context.Context, tenant, contractId string) error
	MarkRolloutRenewalRequested(ctx context.Context, tenant, contractId string) error
	MarkInvoicingStarted(ctx context.Context, tenant, contractId string, invoicingStartedAt time.Time, nextInvoiceDate *time.Time) error
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
								ct.serviceStartedAt=$serviceStartedAt,
								ct.currency=$currency,
								ct.billingCycle=$billingCycle,
								ct.invoicingStartDate=$invoicingStartDate
							WITH ct, t
							OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$createdByUserId}) 
							WHERE $createdByUserId <> ""
							FOREACH (ignore IN CASE WHEN u IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (ct)-[:CREATED_BY]->(u))
							`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"contractId":         contractId,
		"orgId":              data.OrganizationId,
		"createdAt":          data.CreatedAt,
		"updatedAt":          data.UpdatedAt,
		"source":             data.SourceFields.Source,
		"sourceOfTruth":      data.SourceFields.Source,
		"appSource":          data.SourceFields.AppSource,
		"name":               data.Name,
		"contractUrl":        data.ContractUrl,
		"status":             data.Status,
		"renewalCycle":       data.RenewalCycle,
		"renewalPeriods":     data.RenewalPeriods,
		"signedAt":           utils.TimePtrFirstNonNilNillableAsAny(data.SignedAt),
		"serviceStartedAt":   utils.TimePtrFirstNonNilNillableAsAny(data.ServiceStartedAt),
		"createdByUserId":    data.CreatedByUserId,
		"currency":           data.Currency.String(),
		"billingCycle":       data.BillingCycle.String(),
		"invoicingStartDate": utils.ToNeo4jDateAsAny(data.InvoicingStartDate),
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
				ct.updatedAt = $updatedAt,
				ct.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE ct.sourceOfTruth END
				`
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    contractId,
		"updatedAt":     data.UpdatedAt,
		"sourceOfTruth": data.Source,
		"overwrite":     data.Source == constants.SourceOpenline,
	}
	if data.UpdateName {
		cypher += `, ct.name = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR ct.name is null OR ct.name = '' THEN $name ELSE ct.name END `
		params["name"] = data.Name
	}
	if data.UpdateContractUrl {
		cypher += `, ct.contractUrl = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR ct.contractUrl is null OR ct.contractUrl = '' THEN $contractUrl ELSE ct.contractUrl END `
		params["contractUrl"] = data.ContractUrl
	}
	if data.UpdateStatus {
		cypher += `, ct.status = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $status ELSE ct.status END `
		params["status"] = data.Status
	}
	if data.UpdateRenewalPeriods {
		cypher += `, ct.renewalPeriods = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $renewalPeriods ELSE ct.renewalPeriods END `
		params["renewalPeriods"] = data.RenewalPeriods
	}
	if data.UpdateRenewalCycle {
		cypher += `, ct.renewalCycle = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $renewalCycle ELSE ct.renewalCycle END `
		params["renewalCycle"] = data.RenewalCycle
	}
	if data.UpdateServiceStartedAt {
		cypher += `, ct.serviceStartedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $serviceStartedAt ELSE ct.serviceStartedAt END `
		params["serviceStartedAt"] = utils.TimePtrFirstNonNilNillableAsAny(data.ServiceStartedAt)
	}
	if data.UpdateSignedAt {
		cypher += `, ct.signedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $signedAt ELSE ct.signedAt END `
		params["signedAt"] = utils.TimePtrFirstNonNilNillableAsAny(data.SignedAt)
	}
	if data.UpdateEndedAt {
		cypher += `, ct.endedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $endedAt ELSE ct.endedAt END `
		params["endedAt"] = utils.TimePtrFirstNonNilNillableAsAny(data.EndedAt)
	}
	if data.UpdateBillingCycle {
		cypher += `, ct.billingCycle = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $billingCycle ELSE ct.billingCycle END `
		params["billingCycle"] = data.BillingCycle.String()
	}
	if data.UpdateCurrency {
		cypher += `, ct.currency = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $currency ELSE ct.currency END `
		params["currency"] = data.Currency.String()
	}
	if data.UpdateInvoicingStartDate {
		cypher += `, ct.invoicingStartDate = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $invoicingStartDate ELSE ct.invoicingStartDate END `
		params["invoicingStartDate"] = utils.ToNeo4jDateAsAny(data.InvoicingStartDate)
	}
	if data.UpdateAddressLine1 {
		cypher += `, ct.addressLine1 = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $addressLine1 ELSE ct.addressLine1 END `
		params["addressLine1"] = data.AddressLine1
	}
	if data.UpdateAddressLine2 {
		cypher += `, ct.addressLine2 = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $addressLine2 ELSE ct.addressLine2 END `
		params["addressLine2"] = data.AddressLine2
	}
	if data.UpdateLocality {
		cypher += `, ct.locality = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $locality ELSE ct.locality END `
		params["locality"] = data.Locality
	}
	if data.UpdateCountry {
		cypher += `, ct.country = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $country ELSE ct.country END `
		params["country"] = data.Country
	}
	if data.UpdateZip {
		cypher += `, ct.zip = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $zip ELSE ct.zip END `
		params["zip"] = data.Zip
	}
	if data.UpdateOrganizationLegalName {
		cypher += `, ct.organizationLegalName = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $organizationLegalName ELSE ct.organizationLegalName END `
		params["organizationLegalName"] = data.OrganizationLegalName
	}
	if data.UpdateInvoiceEmail {
		cypher += `, ct.invoiceEmail = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $invoiceEmail ELSE ct.invoiceEmail END `
		params["invoiceEmail"] = data.InvoiceEmail
	}
	if data.UpdateInvoiceNote {
		cypher += `, ct.invoiceNote = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $invoiceNote ELSE ct.invoiceNote END `
		params["invoiceNote"] = data.InvoiceNote
	}
	cypher += ` RETURN ct`

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
		"internalStage": neo4jenum.OpportunityInternalStageOpen.String(),
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

func (r *contractWriteRepository) MarkStatusRenewalRequested(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.MarkStatusRenewalRequested")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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

func (r *contractWriteRepository) MarkInvoicingStarted(ctx context.Context, tenant, contractId string, invoicingStartedAt time.Time, nextInvoiceDate *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.MarkInvoicingStarted")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
				SET c.techInvoicingStartedAt=$invoicingStartedAt,
					c.nextInvoiceDate=$nextInvoiceDate`
	params := map[string]any{
		"tenant":             tenant,
		"contractId":         contractId,
		"invoicingStartedAt": invoicingStartedAt,
		"nextInvoiceDate":    utils.ToNeo4jDateAsAny(nextInvoiceDate),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
