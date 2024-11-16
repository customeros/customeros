package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type BankAccountWriteRepository interface {
	CreateBankAccount(ctx context.Context, tenant string, data data_fields.BankAccountFields) error
	UpdateBankAccount(ctx context.Context, tenant string, data data_fields.BankAccountFields) error
	DeleteBankAccount(ctx context.Context, tenant, id string) error
}

type bankAccountWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewBankAccountWriteRepository(driver *neo4j.DriverWithContext, database string) BankAccountWriteRepository {
	return &bankAccountWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *bankAccountWriteRepository) CreateBankAccount(ctx context.Context, tenant string, data data_fields.BankAccountFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountWriteRepository.CreateBankAccount")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "data", data)

	if data.ID == "" {
		err := fmt.Errorf("missing bank account id")
		tracing.TraceErr(span, err)
		return err
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)-[:HAS_BANK_ACCOUNT]->(ba:BankAccount {id:$bankAccountId}) 
							ON CREATE SET 
								ba:BankAccount_%s,
								ba.createdAt=datetime(),
								ba.updatedAt=datetime(),
								ba.source=$source,
								ba.appSource=$appSource,
								ba.bankName=$bankName,
								ba.bankTransferEnabled=$bankTransferEnabled,
								ba.allowInternational=$allowInternational,
								ba.currency=$currency,
								ba.iban=$iban,
								ba.bic=$bic,
								ba.sortCode=$sortCode,
								ba.accountNumber=$accountNumber,
								ba.routingNumber=$routingNumber,
								ba.otherDetails=$otherDetails
							`, tenant)
	params := map[string]any{
		"tenant":              tenant,
		"bankAccountId":       data.ID,
		"source":              utils.IfNotNilString(data.Source),
		"appSource":           utils.IfNotNilString(data.AppSource),
		"bankName":            utils.IfNotNilString(data.BankName),
		"bankTransferEnabled": utils.IfNotNilBool(data.BankTransferEnabled),
		"allowInternational":  utils.IfNotNilBool(data.AllowInternational),
		"currency":            utils.IfNotNilString(data.Currency),
		"iban":                utils.IfNotNilString(data.IBAN),
		"bic":                 utils.IfNotNilString(data.BIC),
		"sortCode":            utils.IfNotNilString(data.SortCode),
		"accountNumber":       utils.IfNotNilString(data.AccountNumber),
		"routingNumber":       utils.IfNotNilString(data.RoutingNumber),
		"otherDetails":        utils.IfNotNilString(data.OtherDetails),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *bankAccountWriteRepository) UpdateBankAccount(ctx context.Context, tenant string, data data_fields.BankAccountFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountWriteRepository.UpdateBankAccount")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "data", data)

	if data.ID == "" {
		err := fmt.Errorf("missing bank account id")
		tracing.TraceErr(span, err)
		return err
	}

	cypher := `MATCH (:Tenant {name:$tenant})-[:HAS_BANK_ACCOUNT]->(ba:BankAccount {id:$bankAccountId}) 
							SET ba.updatedAt=datetime()
							`
	params := map[string]any{
		"tenant":        tenant,
		"bankAccountId": data.ID,
	}
	if data.BankName != nil {
		cypher += `,ba.bankName=$bankName`
		params["bankName"] = *data.BankName
	}
	if data.BankTransferEnabled != nil {
		cypher += `,ba.bankTransferEnabled=$bankTransferEnabled`
		params["bankTransferEnabled"] = *data.BankTransferEnabled
	}
	if data.AllowInternational != nil {
		cypher += `,ba.allowInternational=$allowInternational`
		params["allowInternational"] = *data.AllowInternational
	}
	if data.Currency != nil {
		cypher += `,ba.currency=$currency`
		params["currency"] = *data.Currency
	}
	if data.IBAN != nil {
		cypher += `,ba.iban=$iban`
		params["iban"] = *data.IBAN
	}
	if data.BIC != nil {
		cypher += `,ba.bic=$bic`
		params["bic"] = *data.BIC
	}
	if data.SortCode != nil {
		cypher += `,ba.sortCode=$sortCode`
		params["sortCode"] = *data.SortCode
	}
	if data.AccountNumber != nil {
		cypher += `,ba.accountNumber=$accountNumber`
		params["accountNumber"] = *data.AccountNumber
	}
	if data.RoutingNumber != nil {
		cypher += `,ba.routingNumber=$routingNumber`
		params["routingNumber"] = *data.RoutingNumber
	}
	if data.OtherDetails != nil {
		cypher += `,ba.otherDetails=$otherDetails`
		params["otherDetails"] = *data.OtherDetails
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *bankAccountWriteRepository) DeleteBankAccount(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountWriteRepository.DeleteBankAccount")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})-[r:HAS_BANK_ACCOUNT]->(ba:BankAccount {id:$bankAccountId}) 
							DELETE r, ba`
	params := map[string]any{
		"tenant":        tenant,
		"bankAccountId": id,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
