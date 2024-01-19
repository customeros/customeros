package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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
	DomesticPaymentsBankName          string       `json:"domesticPaymentsBankName"`
	DomesticPaymentsAccountNumber     string       `json:"domesticPaymentsAccountNumber"`
	DomesticPaymentsSortCode          string       `json:"domesticPaymentsSortCode"`
	InternationalPaymentsSwiftBic     string       `json:"internationalPaymentsSwiftBic"`
	InternationalPaymentsBankName     string       `json:"internationalPaymentsBankName"`
	InternationalPaymentsBankAddress  string       `json:"internationalPaymentsBankAddress"`
	InternationalPaymentsInstructions string       `json:"internationalPaymentsInstructions"`
}

type TenantWriteRepository interface {
	CreateTenantBillingProfile(ctx context.Context, tenant string, data TenantBillingProfileCreateFields) error
}

type tenantWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
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
								tbp.domesticPaymentsBankName=$domesticPaymentsBankName,
								tbp.domesticPaymentsAccountNumber=$domesticPaymentsAccountNumber,
								tbp.domesticPaymentsSortCode=$domesticPaymentsSortCode,
								tbp.internationalPaymentsSwiftBic=$internationalPaymentsSwiftBic,
								tbp.internationalPaymentsBankName=$internationalPaymentsBankName,
								tbp.internationalPaymentsBankAddress=$internationalPaymentsBankAddress,
								tbp.internationalPaymentsInstructions=$internationalPaymentsInstructions
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
		"domesticPaymentsBankName":          data.DomesticPaymentsBankName,
		"domesticPaymentsAccountNumber":     data.DomesticPaymentsAccountNumber,
		"domesticPaymentsSortCode":          data.DomesticPaymentsSortCode,
		"internationalPaymentsSwiftBic":     data.InternationalPaymentsSwiftBic,
		"internationalPaymentsBankName":     data.InternationalPaymentsBankName,
		"internationalPaymentsBankAddress":  data.InternationalPaymentsBankAddress,
		"internationalPaymentsInstructions": data.InternationalPaymentsInstructions,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func NewTenantWriteRepository(driver *neo4j.DriverWithContext, database string) TenantWriteRepository {
	return &tenantWriteRepository{
		driver:   driver,
		database: database,
	}
}
