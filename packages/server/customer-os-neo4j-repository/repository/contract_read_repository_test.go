package repository

import (
	"context"
	"github.com/cucumber/godog"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository/tableMappers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository/types"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var contextData map[string]interface{}

func TestFeaturesContractRead(t *testing.T) {
	contextData = make(map[string]interface{})
	//contextData["testingInstance"] = t
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t,
			StopOnFailure: false,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	t := &testing.T{}
	contextData["testingInstance"] = t
	ctx := context.WithValue(context.Background(), "testingInstance", t)
	//ctx := context.Background()
	//sc.Step(`^a tenant was created$`, TenantWasInserted)
	sc.Step(`^(\d+) SLIs are inserted in the database$`, SlisWereInserted)
	sc.Step(`^(\d+) should exist in the neo4j database$`, SlisShouldExist)
	sc.Step(`^the following SLIs are inserted in the database$`, func(table *godog.Table) (context.Context, error) {
		//ctx, _ = CustomSlisWereInserted(ctx, table)
		//return ctx, nil
		return CustomSlisWereInserted(ctx, table)

	})
	sc.Step(`^the SLIs should exist in the neo4j database in a consistent format$`, CustomSlisShouldExist)
	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		//tearDownTestCase(ctx)
		neo4jtest.CleanupAllData(ctx, driver)
		return ctx, nil
	})
}

func CustomSlisWereInserted(ctx context.Context, table *godog.Table) (context.Context, error) {
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organization1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractId := neo4jtest.CreateContract(ctx, driver, tenantName, organization1Id, entity.ContractEntity{})

	sliArray := tableMappers.SliToTable(table)

	for i := 0; i < len(table.Rows)-1; i++ {
		sliArray[i].Id = neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, enum.GetBilledType(sliArray[i].BillingType), sliArray[i].Price, sliArray[i].Quantity, sliArray[i].StartedAt)
	}
	return context.WithValue(ctx, ctxKey{}, sliArray), nil
}

func CustomSlisShouldExist(ctx context.Context) {
	t := contextData["testingInstance"].(*testing.T)
	expectedSlis := ctx.Value(ctxKey{}).([]types.SLI)

	for i := 0; i < len(expectedSlis); i++ {
		actualSlis, err := neo4jtest.GetNodeById(ctx, driver, "ServiceLineItem", expectedSlis[i].Id)
		assert.Nil(t, err)
		require.NotNil(t, actualSlis)
		sliProps := utils.GetPropsFromNode(*actualSlis)

		require.Equal(t, expectedSlis[i].BillingType, utils.GetStringPropOrEmpty(sliProps, "billed"))
		require.Equal(t, expectedSlis[i].Quantity, utils.GetInt64PropOrZero(sliProps, "quantity"))
		require.Equal(t, expectedSlis[i].Price, utils.GetFloatPropOrZero(sliProps, "price"))
		require.Equal(t, expectedSlis[i].StartedAt, utils.GetTimePropOrZeroTime(sliProps, "startedAt"))
	}
}
func SlisWereInserted(ctx context.Context, inserted_slis int) (context.Context, error) {
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organization1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	currentYear := 2023

	contractId := neo4jtest.CreateContract(ctx, driver, tenantName, organization1Id, entity.ContractEntity{})

	sliStartedAt := neo4jtest.FirstTimeOfMonth(currentYear, 1)
	for i := 0; i < inserted_slis; i++ {
		neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, enum.BilledTypeMonthly, 12, 2, sliStartedAt)
	}
	return context.WithValue(ctx, ctxKey{}, inserted_slis), nil
}

func SlisShouldExist(ctx context.Context, actual_number_of_SlIs int) {
	t := contextData["testingInstance"].(*testing.T)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelOrganization:    1,
		neo4jutil.NodeLabelContract:        1,
		neo4jutil.NodeLabelServiceLineItem: actual_number_of_SlIs,
	})
}

func TenantWasInserted(ctx context.Context) {
	neo4jtest.CreateTenant(ctx, driver, tenantName)
}
