package repository

import (
	"context"
	"github.com/cucumber/godog"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
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
