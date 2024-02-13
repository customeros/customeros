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
	Id                                string       `json:"id"`
	CreatedAt                         time.Time    `json:"createdAt"`
	SourceFields                      model.Source `json:"sourceFields"`
	Email                             string       `json:"email"`
	Phone                             string       `json:"phone"`
	LegalName                         string       `json:"legalName"`
	AddressLine1                      string       `json:"addressLine1"`
	AddressLine2                      string       `json:"addressLine2"`
	AddressLine3                      string       `json:"addressLine3"`
	Locality                          string       `json:"locality"`
	Country                           string       `json:"country"`
	Zip                               string       `json:"zip"`
	DomesticPaymentsBankInfo          string       `json:"domesticPaymentsBankInfo"`
	DomesticPaymentsBankName          string       `json:"domesticPaymentsBankName"`
	DomesticPaymentsAccountNumber     string       `json:"domesticPaymentsAccountNumber"`
	DomesticPaymentsSortCode          string       `json:"domesticPaymentsSortCode"`
	InternationalPaymentsBankInfo     string       `json:"internationalPaymentsBankInfo"`
	InternationalPaymentsSwiftBic     string       `json:"internationalPaymentsSwiftBic"`
	InternationalPaymentsBankName     string       `json:"internationalPaymentsBankName"`
	InternationalPaymentsBankAddress  string       `json:"internationalPaymentsBankAddress"`
	InternationalPaymentsInstructions string       `json:"internationalPaymentsInstructions"`
	VatNumber                         string       `json:"vatNumber"`
	SendInvoicesFrom                  string       `json:"sendInvoicesFrom"`
	CanPayWithCard                    bool         `json:"canPayWithCard"`
	CanPayWithDirectDebitSEPA         bool         `json:"canPayWithDirectDebitSEPA"`
	CanPayWithDirectDebitACH          bool         `json:"canPayWithDirectDebitACH"`
	CanPayWithDirectDebitBacs         bool         `json:"canPayWithDirectDebitBacs"`
	CanPayWithPigeon                  bool         `json:"canPayWithPigeon"`
}

type TenantBillingProfileUpdateFields struct {
	Id                                  string    `json:"id"`
	UpdatedAt                           time.Time `json:"updatedAt"`
	Email                               string    `json:"email"`
	Phone                               string    `json:"phone"`
	LegalName                           string    `json:"legalName"`
	AddressLine1                        string    `json:"addressLine1"`
	AddressLine2                        string    `json:"addressLine2"`
	AddressLine3                        string    `json:"addressLine3"`
	Locality                            string    `json:"locality"`
	Country                             string    `json:"country"`
	Zip                                 string    `json:"zip"`
	DomesticPaymentsBankInfo            string    `json:"domesticPaymentsBankInfo"`
	InternationalPaymentsBankInfo       string    `json:"internationalPaymentsBankInfo"`
	VatNumber                           string    `json:"vatNumber"`
	SendInvoicesFrom                    string    `json:"sendInvoicesFrom"`
	CanPayWithCard                      bool      `json:"canPayWithCard"`
	CanPayWithDirectDebitSEPA           bool      `json:"canPayWithDirectDebitSEPA"`
	CanPayWithDirectDebitACH            bool      `json:"canPayWithDirectDebitACH"`
	CanPayWithDirectDebitBacs           bool      `json:"canPayWithDirectDebitBacs"`
	CanPayWithPigeon                    bool      `json:"canPayWithPigeon"`
	UpdateEmail                         bool      `json:"updateEmail"`
	UpdatePhone                         bool      `json:"updatePhone"`
	UpdateLegalName                     bool      `json:"updateLegalName"`
	UpdateAddressLine1                  bool      `json:"updateAddressLine1"`
	UpdateAddressLine2                  bool      `json:"updateAddressLine2"`
	UpdateAddressLine3                  bool      `json:"updateAddressLine3"`
	UpdateLocality                      bool      `json:"updateLocality"`
	UpdateCountry                       bool      `json:"updateCountry"`
	UpdateZip                           bool      `json:"updateZip"`
	UpdateDomesticPaymentsBankInfo      bool      `json:"updateDomesticPaymentsBankInfo"`
	UpdateInternationalPaymentsBankInfo bool      `json:"updateInternationalPaymentsBankInfo"`
	UpdateVatNumber                     bool      `json:"updateVatNumber"`
	UpdateSendInvoicesFrom              bool      `json:"updateSendInvoicesFrom"`
	UpdateCanPayWithCard                bool      `json:"updateCanPayWithCard"`
	UpdateCanPayWithDirectDebitSEPA     bool      `json:"updateCanPayWithDirectDebitSEPA"`
	UpdateCanPayWithDirectDebitACH      bool      `json:"updateCanPayWithDirectDebitACH"`
	UpdateCanPayWithDirectDebitBacs     bool      `json:"updateCanPayWithDirectDebitBacs"`
	UpdateCanPayWithPigeon              bool      `json:"updateCanPayWithPigeon"`
}

type TenantSettingsFields struct {
	UpdatedAt               time.Time     `json:"updatedAt"`
	LogoUrl                 string        `json:"logoUrl"`
	DefaultCurrency         enum.Currency `json:"defaultCurrency"`
	InvoicingEnabled        bool          `json:"invoicingEnabled"`
	InvoicingPostpaid       bool          `json:"invoicingPostpaid"`
	UpdateLogoUrl           bool          `json:"updateLogoUrl"`
	UpdateInvoicingEnabled  bool          `json:"updateInvoicingEnabled"`
	UpdateInvoicingPostpaid bool          `json:"updateInvoicingPostpaid"`
	UpdateDefaultCurrency   bool          `json:"updateDefaultCurrency"`
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
								tbp.email=$email,
								tbp.phone=$phone,
								tbp.legalName=$legalName,	
								tbp.addressLine1=$addressLine1,	
								tbp.addressLine2=$addressLine2,
								tbp.addressLine3=$addressLine3,
								tbp.locality=$locality,
								tbp.country=$country,
								tbp.zip=$zip,
								tbp.domesticPaymentsBankInfo=$domesticPaymentsBankInfo,
								tbp.domesticPaymentsBankName=$domesticPaymentsBankName,
								tbp.domesticPaymentsAccountNumber=$domesticPaymentsAccountNumber,
								tbp.domesticPaymentsSortCode=$domesticPaymentsSortCode,
								tbp.internationalPaymentsBankInfo=$internationalPaymentsBankInfo,
								tbp.internationalPaymentsSwiftBic=$internationalPaymentsSwiftBic,
								tbp.internationalPaymentsBankName=$internationalPaymentsBankName,
								tbp.internationalPaymentsBankAddress=$internationalPaymentsBankAddress,
								tbp.internationalPaymentsInstructions=$internationalPaymentsInstructions,
								tbp.vatNumber=$vatNumber,	
								tbp.sendInvoicesFrom=$sendInvoicesFrom,
								tbp.canPayWithCard=$canPayWithCard,
								tbp.canPayWithDirectDebitSEPA=$canPayWithDirectDebitSEPA,
								tbp.canPayWithDirectDebitACH=$canPayWithDirectDebitACH,	
								tbp.canPayWithDirectDebitBacs=$canPayWithDirectDebitBacs,
								tbp.canPayWithPigeon=$canPayWithPigeon
							`, tenant)
	params := map[string]any{
		"tenant":                            tenant,
		"billingProfileId":                  data.Id,
		"createdAt":                         data.CreatedAt,
		"updatedAt":                         data.CreatedAt,
		"source":                            data.SourceFields.Source,
		"sourceOfTruth":                     data.SourceFields.Source,
		"appSource":                         data.SourceFields.AppSource,
		"email":                             data.Email,
		"phone":                             data.Phone,
		"legalName":                         data.LegalName,
		"addressLine1":                      data.AddressLine1,
		"addressLine2":                      data.AddressLine2,
		"addressLine3":                      data.AddressLine3,
		"locality":                          data.Locality,
		"country":                           data.Country,
		"zip":                               data.Zip,
		"domesticPaymentsBankInfo":          data.DomesticPaymentsBankInfo,
		"domesticPaymentsBankName":          data.DomesticPaymentsBankName,
		"domesticPaymentsAccountNumber":     data.DomesticPaymentsAccountNumber,
		"domesticPaymentsSortCode":          data.DomesticPaymentsSortCode,
		"internationalPaymentsBankInfo":     data.InternationalPaymentsBankInfo,
		"internationalPaymentsSwiftBic":     data.InternationalPaymentsSwiftBic,
		"internationalPaymentsBankName":     data.InternationalPaymentsBankName,
		"internationalPaymentsBankAddress":  data.InternationalPaymentsBankAddress,
		"internationalPaymentsInstructions": data.InternationalPaymentsInstructions,
		"vatNumber":                         data.VatNumber,
		"sendInvoicesFrom":                  data.SendInvoicesFrom,
		"canPayWithCard":                    data.CanPayWithCard,
		"canPayWithDirectDebitSEPA":         data.CanPayWithDirectDebitSEPA,
		"canPayWithDirectDebitACH":          data.CanPayWithDirectDebitACH,
		"canPayWithDirectDebitBacs":         data.CanPayWithDirectDebitBacs,
		"canPayWithPigeon":                  data.CanPayWithPigeon,
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
	if data.UpdateEmail {
		cypher += `,tbp.email=$email`
		params["email"] = data.Email
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
	if data.UpdateZip {
		cypher += `,tbp.zip=$zip`
		params["zip"] = data.Zip
	}
	if data.UpdateDomesticPaymentsBankInfo {
		cypher += `,tbp.domesticPaymentsBankInfo=$domesticPaymentsBankInfo`
		params["domesticPaymentsBankInfo"] = data.DomesticPaymentsBankInfo
	}
	if data.UpdateInternationalPaymentsBankInfo {
		cypher += `,tbp.internationalPaymentsBankInfo=$internationalPaymentsBankInfo`
		params["internationalPaymentsBankInfo"] = data.InternationalPaymentsBankInfo
	}
	if data.UpdateVatNumber {
		cypher += `,tbp.vatNumber=$vatNumber`
		params["vatNumber"] = data.VatNumber
	}
	if data.UpdateSendInvoicesFrom {
		cypher += `,tbp.sendInvoicesFrom=$sendInvoicesFrom`
		params["sendInvoicesFrom"] = data.SendInvoicesFrom
	}
	if data.UpdateCanPayWithCard {
		cypher += `,tbp.canPayWithCard=$canPayWithCard`
		params["canPayWithCard"] = data.CanPayWithCard
	}
	if data.UpdateCanPayWithDirectDebitSEPA {
		cypher += `,tbp.canPayWithDirectDebitSEPA=$canPayWithDirectDebitSEPA`
		params["canPayWithDirectDebitSEPA"] = data.CanPayWithDirectDebitSEPA
	}
	if data.UpdateCanPayWithDirectDebitACH {
		cypher += `,tbp.canPayWithDirectDebitACH=$canPayWithDirectDebitACH`
		params["canPayWithDirectDebitACH"] = data.CanPayWithDirectDebitACH
	}
	if data.UpdateCanPayWithDirectDebitBacs {
		cypher += `,tbp.canPayWithDirectDebitBacs=$canPayWithDirectDebitBacs`
		params["canPayWithDirectDebitBacs"] = data.CanPayWithDirectDebitBacs
	}
	if data.UpdateCanPayWithPigeon {
		cypher += `,tbp.canPayWithPigeon=$canPayWithPigeon`
		params["canPayWithPigeon"] = data.CanPayWithPigeon
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
	if data.UpdateDefaultCurrency {
		cypher += ", ts.defaultCurrency=$defaultCurrency"
		params["defaultCurrency"] = data.DefaultCurrency.String()
	}
	if data.UpdateLogoUrl {
		cypher += ", ts.logoUrl=$logoUrl"
		params["logoUrl"] = data.LogoUrl
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
