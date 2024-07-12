package repository

import (
	"context"
	"github.com/cucumber/godog"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository/tableMappers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository/types"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func CustomSlisWereInserted(ctx context.Context, table *godog.Table) (context.Context, error) {
	sliArray := tableMappers.SliToTable(table)

	for i := 0; i < len(table.Rows)-1; i++ {
		sliArray[i].Id = test.InsertServiceLineItem(ctx, driver, tenantName, sliArray[i].ContractId, enum.DecodeBilledType(sliArray[i].BillingType), sliArray[i].Price, sliArray[i].Quantity, sliArray[i].StartedAt)
	}
	return context.WithValue(ctx, ctxKey{}, sliArray), nil
}
func SlisWereInserted(ctx context.Context, inserted_slis int, contractId string) (context.Context, error) {
	currentYear := 2023

	sliStartedAt := utils.FirstTimeOfMonth(currentYear, 1)
	for i := 0; i < inserted_slis; i++ {
		test.InsertServiceLineItem(ctx, driver, tenantName, contractId, enum.BilledTypeMonthly, 12, 2, sliStartedAt)
	}
	return context.WithValue(ctx, ctxKey{}, inserted_slis), nil
}

func SlisShouldExist(ctx context.Context, actual_number_of_SlIs int) {
	t := contextData["testingInstance"].(*testing.T)

	test.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		model.NodeLabelOrganization:    1,
		model.NodeLabelContract:        1,
		model.NodeLabelServiceLineItem: actual_number_of_SlIs,
	})
}

func CustomSlisShouldExist(ctx context.Context) {
	t := contextData["testingInstance"].(*testing.T)
	expectedSlis := ctx.Value(ctxKey{}).([]types.SLI)

	for i := 0; i < len(expectedSlis); i++ {
		actualSli, err := test.GetNodeById(ctx, driver, "ServiceLineItem", expectedSlis[i].Id)
		assert.Nil(t, err)
		require.NotNil(t, actualSli)
		actualSliProps := utils.GetPropsFromNode(*actualSli)

		require.Equal(t, expectedSlis[i].BillingType, utils.GetStringPropOrEmpty(actualSliProps, "billed"))
		require.Equal(t, expectedSlis[i].Quantity, utils.GetInt64PropOrZero(actualSliProps, "quantity"))
		require.Equal(t, expectedSlis[i].Price, utils.GetFloatPropOrZero(actualSliProps, "price"))
		require.Equal(t, expectedSlis[i].StartedAt, utils.GetTimePropOrZeroTime(actualSliProps, "startedAt"))

		test.AssertRelationship(ctx, t, driver, expectedSlis[i].ContractId, "HAS_SERVICE", utils.GetStringPropOrEmpty(actualSliProps, "id"))
	}
}
