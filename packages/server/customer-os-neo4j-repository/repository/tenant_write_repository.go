package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type TenantBillingProfileCreateFields struct {
	Id                     string       `json:"id"`
	CreatedAt              time.Time    `json:"createdAt"`
	SourceFields           model.Source `json:"sourceFields"`
	Phone                  string       `json:"phone"`
	LegalName              string       `json:"legalName"`
	AddressLine1           string       `json:"addressLine1"`
	AddressLine2           string       `json:"addressLine2"`
	AddressLine3           string       `json:"addressLine3"`
	Locality               string       `json:"locality"`
	Country                string       `json:"country"`
	Region                 string       `json:"region"`
	Zip                    string       `json:"zip"`
	VatNumber              string       `json:"vatNumber"`
	SendInvoicesFrom       string       `json:"sendInvoicesFrom"`
	SendInvoicesBcc        string       `json:"sendInvoicesBcc"`
	CanPayWithPigeon       bool         `json:"canPayWithPigeon"`
	CanPayWithBankTransfer bool         `json:"canPayWithBankTransfer"`
	Check                  bool         `json:"check"`
}

type TenantBillingProfileUpdateFields struct {
	Id                           string    `json:"id"`
	UpdatedAt                    time.Time `json:"updatedAt"`
	Phone                        string    `json:"phone"`
	LegalName                    string    `json:"legalName"`
	AddressLine1                 string    `json:"addressLine1"`
	AddressLine2                 string    `json:"addressLine2"`
	AddressLine3                 string    `json:"addressLine3"`
	Locality                     string    `json:"locality"`
	Country                      string    `json:"country"`
	Region                       string    `json:"region"`
	Zip                          string    `json:"zip"`
	VatNumber                    string    `json:"vatNumber"`
	SendInvoicesFrom             string    `json:"sendInvoicesFrom"`
	SendInvoicesBcc              string    `json:"sendInvoicesBcc"`
	CanPayWithPigeon             bool      `json:"canPayWithPigeon"`
	CanPayWithBankTransfer       bool      `json:"canPayWithBankTransfer"`
	Check                        bool      `json:"check"`
	UpdatePhone                  bool      `json:"updatePhone"`
	UpdateLegalName              bool      `json:"updateLegalName"`
	UpdateAddressLine1           bool      `json:"updateAddressLine1"`
	UpdateAddressLine2           bool      `json:"updateAddressLine2"`
	UpdateAddressLine3           bool      `json:"updateAddressLine3"`
	UpdateLocality               bool      `json:"updateLocality"`
	UpdateCountry                bool      `json:"updateCountry"`
	UpdateRegion                 bool      `json:"updateRegion"`
	UpdateZip                    bool      `json:"updateZip"`
	UpdateVatNumber              bool      `json:"updateVatNumber"`
	UpdateSendInvoicesFrom       bool      `json:"updateSendInvoicesFrom"`
	UpdateSendInvoicesBcc        bool      `json:"updateSendInvoicesBcc"`
	UpdateCanPayWithPigeon       bool      `json:"updateCanPayWithPigeon"`
	UpdateCanPayWithBankTransfer bool      `json:"updateCanPayWithBankTransfer"`
	UpdateCheck                  bool      `json:"updateCheck"`
}

type TenantSettingsFields struct {
	UpdatedAt                  time.Time     `json:"updatedAt"`
	LogoRepositoryFileId       string        `json:"logoRepositoryFileId"`
	BaseCurrency               enum.Currency `json:"baseCurrency"`
	InvoicingEnabled           bool          `json:"invoicingEnabled"`
	InvoicingPostpaid          bool          `json:"invoicingPostpaid"`
	UpdateLogoRepositoryFileId bool          `json:"updateLogoRepositoryFileId"`
	UpdateInvoicingEnabled     bool          `json:"updateInvoicingEnabled"`
	UpdateInvoicingPostpaid    bool          `json:"updateInvoicingPostpaid"`
	UpdateBaseCurrency         bool          `json:"updateBaseCurrency"`
}

type TenantWriteRepository interface {
	CreateTenantBillingProfile(ctx context.Context, tenant string, data TenantBillingProfileCreateFields) error
	UpdateTenantBillingProfile(ctx context.Context, tenant string, data TenantBillingProfileUpdateFields) error
	UpdateTenantSettings(ctx context.Context, tenant string, data TenantSettingsFields) error
}

type tenantWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewTenantWriteRepository(driver *neo4j.DriverWithContext, database string) TenantWriteRepository {
	return &tenantWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *tenantWriteRepository) CreateTenantBillingProfile(ctx context.Context, tenant string, data TenantBillingProfileCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.CreateTenantBillingProfile")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)-[:HAS_BILLING_PROFILE]->(tbp:TenantBillingProfile {id:$billingProfileId}) 
							ON CREATE SET 
								tbp:TenantBillingProfile_%s,
								tbp.createdAt=$createdAt,
								tbp.updatedAt=$updatedAt,
								tbp.source=$source,
								tbp.sourceOfTruth=$sourceOfTruth,
								tbp.appSource=$appSource,
								tbp.phone=$phone,
								tbp.legalName=$legalName,	
								tbp.addressLine1=$addressLine1,	
								tbp.addressLine2=$addressLine2,
								tbp.addressLine3=$addressLine3,
								tbp.locality=$locality,
								tbp.country=$country,
								tbp.region=$region,
								tbp.zip=$zip,
								tbp.vatNumber=$vatNumber,	
								tbp.sendInvoicesFrom=$sendInvoicesFrom,
								tbp.sendInvoicesBcc=$sendInvoicesBcc,
								tbp.canPayWithPigeon=$canPayWithPigeon,
								tbp.canPayWithBankTransfer=$canPayWithBankTransfer,
								tbp.check=$check
							`, tenant)
	params := map[string]any{
		"tenant":                 tenant,
		"billingProfileId":       data.Id,
		"createdAt":              data.CreatedAt,
		"updatedAt":              data.CreatedAt,
		"source":                 data.SourceFields.Source,
		"sourceOfTruth":          data.SourceFields.Source,
		"appSource":              data.SourceFields.AppSource,
		"phone":                  data.Phone,
		"legalName":              data.LegalName,
		"addressLine1":           data.AddressLine1,
		"addressLine2":           data.AddressLine2,
		"addressLine3":           data.AddressLine3,
		"locality":               data.Locality,
		"country":                data.Country,
		"region":                 data.Region,
		"zip":                    data.Zip,
		"vatNumber":              data.VatNumber,
		"sendInvoicesFrom":       data.SendInvoicesFrom,
		"sendInvoicesBcc":        data.SendInvoicesBcc,
		"canPayWithPigeon":       data.CanPayWithPigeon,
		"canPayWithBankTransfer": data.CanPayWithBankTransfer,
		"check":                  data.Check,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *tenantWriteRepository) UpdateTenantBillingProfile(ctx context.Context, tenant string, data TenantBillingProfileUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.UpdateTenantBillingProfile")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (:Tenant {name:$tenant})-[:HAS_BILLING_PROFILE]->(tbp:TenantBillingProfile {id:$billingProfileId}) 
							SET tbp.updatedAt=$updatedAt
							
							`
	params := map[string]any{
		"tenant":           tenant,
		"billingProfileId": data.Id,
		"updatedAt":        data.UpdatedAt,
	}
	if data.UpdatePhone {
		cypher += `,tbp.phone=$phone`
		params["phone"] = data.Phone
	}
	if data.UpdateLegalName {
		cypher += `,tbp.legalName=$legalName`
		params["legalName"] = data.LegalName
	}
	if data.UpdateAddressLine1 {
		cypher += `,tbp.addressLine1=$addressLine1`
		params["addressLine1"] = data.AddressLine1
	}
	if data.UpdateAddressLine2 {
		cypher += `,tbp.addressLine2=$addressLine2`
		params["addressLine2"] = data.AddressLine2
	}
	if data.UpdateAddressLine3 {
		cypher += `,tbp.addressLine3=$addressLine3`
		params["addressLine3"] = data.AddressLine3
	}
	if data.UpdateLocality {
		cypher += `,tbp.locality=$locality`
		params["locality"] = data.Locality
	}
	if data.UpdateCountry {
		cypher += `,tbp.country=$country`
		params["country"] = data.Country
	}
	if data.UpdateRegion {
		cypher += `,tbp.region=$region`
		params["region"] = data.Region
	}
	if data.UpdateZip {
		cypher += `,tbp.zip=$zip`
		params["zip"] = data.Zip
	}
	if data.UpdateVatNumber {
		cypher += `,tbp.vatNumber=$vatNumber`
		params["vatNumber"] = data.VatNumber
	}
	if data.UpdateSendInvoicesFrom {
		cypher += `,tbp.sendInvoicesFrom=$sendInvoicesFrom`
		params["sendInvoicesFrom"] = data.SendInvoicesFrom
	}
	if data.UpdateSendInvoicesBcc {
		cypher += `,tbp.sendInvoicesBcc=$sendInvoicesBcc`
		params["sendInvoicesBcc"] = data.SendInvoicesBcc
	}
	if data.UpdateCanPayWithPigeon {
		cypher += `,tbp.canPayWithPigeon=$canPayWithPigeon`
		params["canPayWithPigeon"] = data.CanPayWithPigeon
	}
	if data.UpdateCanPayWithBankTransfer {
		cypher += `,tbp.canPayWithBankTransfer=$canPayWithBankTransfer`
		params["canPayWithBankTransfer"] = data.CanPayWithBankTransfer
	}
	if data.UpdateCheck {
		cypher += `,tbp.check=$check`
		params["check"] = data.Check
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *tenantWriteRepository) UpdateTenantSettings(ctx context.Context, tenant string, data TenantSettingsFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.UpdateTenantBillingProfile")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (t:Tenant {name:$tenant})
				MERGE (t)-[:HAS_SETTINGS]->(ts:TenantSettings {tenant:$tenant})
				ON CREATE SET
					ts.id=randomUUID(),
					ts.createdAt=$now
				SET
					ts.updatedAt=$updatedAt`
	params := map[string]any{
		"tenant":    tenant,
		"updatedAt": data.UpdatedAt,
		"now":       utils.Now(),
	}
	if data.UpdateInvoicingEnabled {
		cypher += ", ts.invoicingEnabled=$invoicingEnabled"
		params["invoicingEnabled"] = data.InvoicingEnabled
	}
	if data.UpdateInvoicingPostpaid {
		cypher += ", ts.invoicingPostpaid=$invoicingPostpaid"
		params["invoicingPostpaid"] = data.InvoicingPostpaid
	}
	if data.UpdateBaseCurrency {
		cypher += ", ts.baseCurrency=$baseCurrency"
		params["baseCurrency"] = data.BaseCurrency.String()
	}
	if data.UpdateLogoRepositoryFileId {
		cypher += ", ts.logoRepositoryFileId=$logoRepositoryFileId"
		params["logoRepositoryFileId"] = data.LogoRepositoryFileId
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
