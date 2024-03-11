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

type BankAccountCreateFields struct {
	Id                  string        `json:"id"`
	CreatedAt           time.Time     `json:"createdAt"`
	SourceFields        model.Source  `json:"sourceFields"`
	BankName            string        `json:"bankName"`
	BankTransferEnabled bool          `json:"bankTransferEnabled"`
	Currency            enum.Currency `json:"currency"`
	Iban                string        `json:"iban"`
	Bic                 string        `json:"bic"`
	SortCode            string        `json:"sortCode"`
	AccountNumber       string        `json:"accountNumber"`
	RoutingNumber       string        `json:"routingNumber"`
}

type BankAccountUpdateFields struct {
	Id                        string        `json:"id"`
	UpdatedAt                 time.Time     `json:"updatedAt"`
	BankName                  string        `json:"bankName"`
	BankTransferEnabled       bool          `json:"bankTransferEnabled"`
	Currency                  enum.Currency `json:"currency"`
	Iban                      string        `json:"iban"`
	Bic                       string        `json:"bic"`
	SortCode                  string        `json:"sortCode"`
	AccountNumber             string        `json:"accountNumber"`
	RoutingNumber             string        `json:"routingNumber"`
	UpdateBankName            bool          `json:"updateBankName"`
	UpdateBankTransferEnabled bool          `json:"updateBankTransferEnabled"`
	UpdateCurrency            bool          `json:"updateCurrency"`
	UpdateIban                bool          `json:"updateIban"`
	UpdateBic                 bool          `json:"updateBic"`
	UpdateSortCode            bool          `json:"updateSortCode"`
	UpdateAccountNumber       bool          `json:"updateAccountNumber"`
	UpdateRoutingNumber       bool          `json:"updateRoutingNumber"`
}

type BankAccountWriteRepository interface {
	CreateBankAccount(ctx context.Context, tenant string, data BankAccountCreateFields) error
	UpdateBankAccount(ctx context.Context, tenant string, data BankAccountUpdateFields) error
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

func (r *bankAccountWriteRepository) CreateBankAccount(ctx context.Context, tenant string, data BankAccountCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountWriteRepository.CreateBankAccount")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)-[:HAS_BANK_ACCOUNT]->(ba:BankAccount {id:$bankAccountId}) 
							ON CREATE SET 
								ba:BankAccount_%s,
								ba.createdAt=$createdAt,
								ba.updatedAt=$updatedAt,
								ba.source=$source,
								ba.sourceOfTruth=$sourceOfTruth,
								ba.appSource=$appSource,
								ba.bankName=$bankName,
								ba.bankTransferEnabled=$bankTransferEnabled,
								ba.currency=$currency,
								ba.iban=$iban,
								ba.bic=$bic,
								ba.sortCode=$sortCode,
								ba.accountNumber=$accountNumber,
								ba.routingNumber=$routingNumber
							`, tenant)
	params := map[string]any{
		"tenant":              tenant,
		"bankAccountId":       data.Id,
		"createdAt":           data.CreatedAt,
		"updatedAt":           data.CreatedAt,
		"source":              data.SourceFields.Source,
		"sourceOfTruth":       data.SourceFields.Source,
		"appSource":           data.SourceFields.AppSource,
		"bankName":            data.BankName,
		"bankTransferEnabled": data.BankTransferEnabled,
		"currency":            data.Currency.String(),
		"iban":                data.Iban,
		"bic":                 data.Bic,
		"sortCode":            data.SortCode,
		"accountNumber":       data.AccountNumber,
		"routingNumber":       data.RoutingNumber,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *bankAccountWriteRepository) UpdateBankAccount(ctx context.Context, tenant string, data BankAccountUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountWriteRepository.UpdateBankAccount")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (:Tenant {name:$tenant})-[:HAS_BANK_ACCOUNT]->(ba:BankAccount {id:$bankAccountId}) 
							SET ba.updatedAt=$updatedAt
							`
	params := map[string]any{
		"tenant":        tenant,
		"bankAccountId": data.Id,
		"updatedAt":     data.UpdatedAt,
	}
	if data.UpdateBankName {
		cypher += `,ba.bankName=$bankName`
		params["bankName"] = data.BankName
	}
	if data.UpdateBankTransferEnabled {
		cypher += `,ba.bankTransferEnabled=$bankTransferEnabled`
		params["bankTransferEnabled"] = data.BankTransferEnabled
	}
	if data.UpdateCurrency {
		cypher += `,ba.currency=$currency`
		params["currency"] = data.Currency.String()
	}
	if data.UpdateIban {
		cypher += `,ba.iban=$iban`
		params["iban"] = data.Iban
	}
	if data.UpdateBic {
		cypher += `,ba.bic=$bic`
		params["bic"] = data.Bic
	}
	if data.UpdateSortCode {
		cypher += `,ba.sortCode=$sortCode`
		params["sortCode"] = data.SortCode
	}
	if data.UpdateAccountNumber {
		cypher += `,ba.accountNumber=$accountNumber`
		params["accountNumber"] = data.AccountNumber
	}
	if data.UpdateRoutingNumber {
		cypher += `,ba.routingNumber=$routingNumber`
		params["routingNumber"] = data.RoutingNumber
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

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
