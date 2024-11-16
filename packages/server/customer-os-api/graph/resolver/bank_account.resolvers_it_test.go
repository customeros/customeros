package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_BankAccounts(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	today := utils.Now()
	yesterday := today.AddDate(0, 0, -1)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	account1 := neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{
		CreatedAt:           today,
		UpdatedAt:           today,
		BankName:            "bankName1",
		BankTransferEnabled: true,
		AllowInternational:  true,
		Currency:            neo4jenum.CurrencyEUR,
		Iban:                "iban1",
		Bic:                 "bic1",
		SortCode:            "sortCode1",
		AccountNumber:       "accountNumber1",
		RoutingNumber:       "routingNumber1",
		OtherDetails:        "otherDetails1",
	})
	account2 := neo4jtest.CreateBankAccount(ctx, driver, tenantName, neo4jentity.BankAccountEntity{
		CreatedAt:           yesterday,
		UpdatedAt:           yesterday,
		BankName:            "bankName2",
		BankTransferEnabled: false,
		AllowInternational:  false,
		Currency:            neo4jenum.CurrencyUSD,
		Iban:                "iban2",
		Bic:                 "bic2",
		SortCode:            "sortCode2",
		AccountNumber:       "accountNumber2",
		RoutingNumber:       "routingNumber2",
		OtherDetails:        "otherDetails2",
	})

	rawResponse, err := c.RawPost(getQuery("bank_account/get_bank_accounts"))
	assertRawResponseSuccess(t, rawResponse, err)

	var graphqlResponse struct {
		BankAccounts []model.BankAccount
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
	require.Nil(t, err)
	require.NotNil(t, graphqlResponse)

	require.Equal(t, 2, len(graphqlResponse.BankAccounts))
	ba1 := graphqlResponse.BankAccounts[0]
	require.Equal(t, account2, ba1.Metadata.ID)
	require.Equal(t, yesterday, ba1.Metadata.Created)
	require.Equal(t, "bankName2", *ba1.BankName)
	require.Equal(t, false, ba1.BankTransferEnabled)
	require.Equal(t, false, ba1.AllowInternational)
	require.Equal(t, model.CurrencyUsd, *ba1.Currency)
	require.Equal(t, "iban2", *ba1.Iban)
	require.Equal(t, "bic2", *ba1.Bic)
	require.Equal(t, "sortCode2", *ba1.SortCode)
	require.Equal(t, "accountNumber2", *ba1.AccountNumber)
	require.Equal(t, "routingNumber2", *ba1.RoutingNumber)
	require.Equal(t, "otherDetails2", *ba1.OtherDetails)

	ba2 := graphqlResponse.BankAccounts[1]
	require.Equal(t, account1, ba2.Metadata.ID)
	require.Equal(t, today, ba2.Metadata.Created)
	require.Equal(t, "bankName1", *ba2.BankName)
	require.Equal(t, true, ba2.BankTransferEnabled)
	require.Equal(t, true, ba2.AllowInternational)
	require.Equal(t, model.CurrencyEur, *ba2.Currency)
	require.Equal(t, "iban1", *ba2.Iban)
	require.Equal(t, "bic1", *ba2.Bic)
	require.Equal(t, "sortCode1", *ba2.SortCode)
	require.Equal(t, "accountNumber1", *ba2.AccountNumber)
	require.Equal(t, "routingNumber1", *ba2.RoutingNumber)
	require.Equal(t, "otherDetails1", *ba2.OtherDetails)
}
