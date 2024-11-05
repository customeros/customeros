package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

// Deprecated
type ContractCreateFields struct {
	OrganizationId         string             `json:"organizationId"`
	Name                   string             `json:"name"`
	ContractUrl            string             `json:"contractUrl"`
	CreatedByUserId        string             `json:"createdByUserId"`
	ServiceStartedAt       *time.Time         `json:"serviceStartedAt,omitempty"`
	SignedAt               *time.Time         `json:"signedAt,omitempty"`
	LengthInMonths         int64              `json:"lengthInMonths"`
	Status                 string             `json:"status"`
	CreatedAt              time.Time          `json:"createdAt"`
	SourceFields           model.SourceFields `json:"sourceFields"`
	BillingCycleInMonths   int64              `json:"billingCycleInMonths"`
	Currency               neo4jenum.Currency `json:"currency"`
	InvoicingStartDate     *time.Time         `json:"invoicingStartDate,omitempty"`
	InvoicingEnabled       bool               `json:"invoicingEnabled"`
	PayOnline              bool               `json:"payOnline"`
	PayAutomatically       bool               `json:"payAutomatically"`
	CanPayWithCard         bool               `json:"canPayWithCard"`
	CanPayWithDirectDebit  bool               `json:"canPayWithDirectDebit"`
	CanPayWithBankTransfer bool               `json:"canPayWithBankTransfer"`
	AutoRenew              bool               `json:"autoRenew"`
	Check                  bool               `json:"check"`
	DueDays                int64              `json:"dueDays"`
	Country                string             `json:"country"`
	Approved               bool               `json:"approved"`
}

// Deprecated
type ContractUpdateFields struct {
	Name                         string             `json:"name"`
	ContractUrl                  string             `json:"contractUrl"`
	Status                       string             `json:"status"`
	Source                       string             `json:"source"`
	LengthInMonths               int64              `json:"lengthInMonths"`
	ServiceStartedAt             *time.Time         `json:"serviceStartedAt"`
	SignedAt                     *time.Time         `json:"signedAt"`
	EndedAt                      *time.Time         `json:"endedAt"`
	BillingCycleInMonths         int64              `json:"billingCycleInMonths"`
	Currency                     neo4jenum.Currency `json:"currency"`
	InvoicingStartDate           *time.Time         `json:"invoicingStartDate,omitempty"`
	NextInvoiceDate              *time.Time         `json:"nextInvoiceDate,omitempty"`
	AddressLine1                 string             `json:"addressLine1"`
	AddressLine2                 string             `json:"addressLine2"`
	Locality                     string             `json:"locality"`
	Country                      string             `json:"country"`
	Region                       string             `json:"region"`
	Zip                          string             `json:"zip"`
	OrganizationLegalName        string             `json:"organizationLegalName"`
	InvoiceEmail                 string             `json:"invoiceEmail"`
	InvoiceEmailCC               []string           `json:"invoiceEmailCC"`
	InvoiceEmailBCC              []string           `json:"invoiceEmailBCC"`
	InvoiceNote                  string             `json:"invoiceNote"`
	InvoicingEnabled             bool               `json:"invoicingEnabled"`
	PayOnline                    bool               `json:"payOnline"`
	PayAutomatically             bool               `json:"payAutomatically"`
	CanPayWithCard               bool               `json:"canPayWithCard"`
	CanPayWithDirectDebit        bool               `json:"canPayWithDirectDebit"`
	CanPayWithBankTransfer       bool               `json:"canPayWithBankTransfer"`
	AutoRenew                    bool               `json:"autoRenew"`
	DueDays                      int64              `json:"dueDays"`
	Check                        bool               `json:"check"`
	Approved                     bool               `json:"approved"`
	UpdateName                   bool               `json:"updateName"`
	UpdateContractUrl            bool               `json:"updateContractUrl"`
	UpdateStatus                 bool               `json:"updateStatus"`
	UpdateServiceStartedAt       bool               `json:"updateServiceStartedAt"`
	UpdateSignedAt               bool               `json:"updateSignedAt"`
	UpdateEndedAt                bool               `json:"updateEndedAt"`
	UpdateBillingCycleInMonths   bool               `json:"updateBillingCycleInMonths"`
	UpdateCurrency               bool               `json:"updateCurrency"`
	UpdateInvoicingStartDate     bool               `json:"updateInvoicingStartDate"`
	UpdateNextInvoiceDate        bool               `json:"updateNextInvoiceDate"`
	UpdateAddressLine1           bool               `json:"updateAddressLine1"`
	UpdateAddressLine2           bool               `json:"updateAddressLine2"`
	UpdateLocality               bool               `json:"updateLocality"`
	UpdateCountry                bool               `json:"updateCountry"`
	UpdateRegion                 bool               `json:"updateRegion"`
	UpdateZip                    bool               `json:"updateZip"`
	UpdateOrganizationLegalName  bool               `json:"updateOrganizationLegalName"`
	UpdateInvoiceEmail           bool               `json:"updateInvoiceEmail"`
	UpdateInvoiceEmailCC         bool               `json:"UpdateInvoiceEmailCC"`
	UpdateInvoiceEmailBCC        bool               `json:"UpdateInvoiceEmailBCC"`
	UpdateInvoiceNote            bool               `json:"updateInvoiceNote"`
	UpdateCanPayWithCard         bool               `json:"updateCanPayWithCard"`
	UpdateCanPayWithDirectDebit  bool               `json:"updateCanPayWithDirectDebit"`
	UpdateCanPayWithBankTransfer bool               `json:"updateCanPayWithBankTransfer"`
	UpdateInvoicingEnabled       bool               `json:"updateInvoicingEnabled"`
	UpdatePayOnline              bool               `json:"updatePayOnline"`
	UpdatePayAutomatically       bool               `json:"updatePayAutomatically"`
	UpdateAutoRenew              bool               `json:"updateAutoRenew"`
	UpdateCheck                  bool               `json:"updateCheck"`
	UpdateDueDays                bool               `json:"updateDueDays"`
	UpdateLengthInMonths         bool               `json:"updateLengthInMonths"`
	UpdateApproved               bool               `json:"updateApproved"`
}

type ContractWriteRepository interface {
	// Deprecated
	CreateForOrganizationOld(ctx context.Context, tenant, contractId string, data ContractCreateFields) error
	CreateForOrganization(ctx context.Context, tenant, contractId string, data data_fields.ContractSaveFields) error
	// Deprecated
	UpdateContractOld(ctx context.Context, tenant, contractId string, data ContractUpdateFields) error
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

func (r *contractWriteRepository) CreateForOrganizationOld(ctx context.Context, tenant, contractId string, data ContractCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.CreateForOrganizationOld")
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
								ct.sourceOfTruth=$sourceOfTruth,
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
	params := map[string]any{
		"tenant":                 tenant,
		"contractId":             contractId,
		"orgId":                  data.OrganizationId,
		"createdAt":              data.CreatedAt,
		"source":                 data.SourceFields.Source,
		"sourceOfTruth":          data.SourceFields.Source,
		"appSource":              data.SourceFields.AppSource,
		"name":                   data.Name,
		"contractUrl":            data.ContractUrl,
		"status":                 data.Status,
		"signedAt":               utils.ToDateAsAny(data.SignedAt),
		"serviceStartedAt":       utils.ToDateAsAny(data.ServiceStartedAt),
		"createdByUserId":        data.CreatedByUserId,
		"currency":               data.Currency.String(),
		"billingCycleInMonths":   data.BillingCycleInMonths,
		"invoicingStartDate":     utils.ToNeo4jDateAsAny(data.InvoicingStartDate),
		"invoicingEnabled":       data.InvoicingEnabled,
		"payOnline":              data.PayOnline,
		"payAutomatically":       data.PayAutomatically,
		"canPayWithCard":         data.CanPayWithCard,
		"canPayWithDirectDebit":  data.CanPayWithDirectDebit,
		"canPayWithBankTransfer": data.CanPayWithBankTransfer,
		"autoRenew":              data.AutoRenew,
		"check":                  data.Check,
		"dueDays":                data.DueDays,
		"country":                data.Country,
		"lengthInMonths":         data.LengthInMonths,
		"approved":               data.Approved,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
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
								ct.sourceOfTruth=$sourceOfTruth,
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

func (r *contractWriteRepository) UpdateContractOld(ctx context.Context, tenant, contractId string, data ContractUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractWriteRepository.UpdateContractOld")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract {id:$contractId})
				SET 
				ct.updatedAt = datetime(),
				ct.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE ct.sourceOfTruth END
				`
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    contractId,
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
	if data.UpdateServiceStartedAt {
		cypher += `, ct.serviceStartedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $serviceStartedAt ELSE ct.serviceStartedAt END `
		params["serviceStartedAt"] = utils.ToDateAsAny(data.ServiceStartedAt)
	}
	if data.UpdateSignedAt {
		cypher += `, ct.signedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $signedAt ELSE ct.signedAt END `
		params["signedAt"] = utils.ToDateAsAny(data.SignedAt)
	}
	if data.UpdateEndedAt {
		cypher += `, ct.endedAt = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $endedAt ELSE ct.endedAt END `
		params["endedAt"] = utils.ToDateAsAny(data.EndedAt)
	}
	if data.UpdateBillingCycleInMonths {
		cypher += `, ct.billingCycleInMonths = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $billingCycleInMonths ELSE ct.billingCycleInMonths END `
		params["billingCycleInMonths"] = data.BillingCycleInMonths
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
	if data.UpdateRegion {
		cypher += `, ct.region = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $region ELSE ct.region END `
		params["region"] = data.Region

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
	if data.UpdateInvoiceEmailCC {
		cypher += `, ct.invoiceEmailCC = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $invoiceEmailCC ELSE ct.invoiceEmailCC END `
		params["invoiceEmailCC"] = data.InvoiceEmailCC
	}
	if data.UpdateInvoiceEmailBCC {
		cypher += `, ct.invoiceEmailBCC = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $invoiceEmailBCC ELSE ct.invoiceEmailBCC END `
		params["invoiceEmailBCC"] = data.InvoiceEmailBCC
	}
	if data.UpdateInvoiceNote {
		cypher += `, ct.invoiceNote = CASE WHEN ct.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $invoiceNote ELSE ct.invoiceNote END `
		params["invoiceNote"] = data.InvoiceNote
	}
	if data.UpdateNextInvoiceDate {
		cypher += `, ct.nextInvoiceDate=$nextInvoiceDate `
		params["nextInvoiceDate"] = utils.ToNeo4jDateAsAny(data.NextInvoiceDate)
	}
	if data.UpdateCanPayWithCard {
		cypher += `, ct.canPayWithCard=$canPayWithCard `
		params["canPayWithCard"] = data.CanPayWithCard
	}
	if data.UpdateCanPayWithDirectDebit {
		cypher += `, ct.canPayWithDirectDebit=$canPayWithDirectDebit `
		params["canPayWithDirectDebit"] = data.CanPayWithDirectDebit
	}
	if data.UpdateCanPayWithBankTransfer {
		cypher += `, ct.canPayWithBankTransfer=$canPayWithBankTransfer `
		params["canPayWithBankTransfer"] = data.CanPayWithBankTransfer
	}
	if data.UpdateInvoicingEnabled {
		cypher += `, ct.invoicingEnabled=$invoicingEnabled `
		params["invoicingEnabled"] = data.InvoicingEnabled
	}
	if data.UpdatePayOnline {
		cypher += `, ct.payOnline=$payOnline `
		params["payOnline"] = data.PayOnline
	}
	if data.UpdatePayAutomatically {
		cypher += `, ct.payAutomatically=$payAutomatically `
		params["payAutomatically"] = data.PayAutomatically
	}
	if data.UpdateAutoRenew {
		cypher += `, ct.autoRenew=$autoRenew `
		params["autoRenew"] = data.AutoRenew
	}
	if data.UpdateCheck {
		cypher += `, ct.check=$check `
		params["check"] = data.Check
	}
	if data.UpdateDueDays {
		cypher += `, ct.dueDays=$dueDays `
		params["dueDays"] = data.DueDays
	}
	if data.UpdateLengthInMonths {
		cypher += `, ct.lengthInMonths=$lengthInMonths `
		params["lengthInMonths"] = data.LengthInMonths
	}
	if data.UpdateApproved {
		cypher += `, ct.approved=$approved `
		params["approved"] = data.Approved
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
