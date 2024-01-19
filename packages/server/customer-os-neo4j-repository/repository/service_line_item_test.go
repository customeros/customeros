package repository

import (
	"context"
	"github.com/cucumber/godog"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"testing"
)

var contextData map[string]interface{}

func TestFeaturesCustomSLIsAreProperlyInserted(t *testing.T) {
	contextData = make(map[string]interface{})
	contextData["testingInstance"] = t
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenarioCustomSLIsAreProperlyInserted,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t,
			StopOnFailure: false,
			Tags:          "tag_custom_sLIs_are_properly_inserted",
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestFeaturesDefaultSLIsAreProperlyInserted(t *testing.T) {
	contextData = make(map[string]interface{})
	contextData["testingInstance"] = t
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenarioDefaultSLIsAreProperlyInserted,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t,
			StopOnFailure: false,
			Strict:        false,
			Tags:          "tag_default_sLIs_are_properly_inserted",
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenarioCustomSLIsAreProperlyInserted(sc *godog.ScenarioContext) {
	ctx := context.Background()
	sc.Step(`^a tenant was created$`, TenantWasInserted)
	sc.Step(`^the organization with the id ([^"]*) was created$`, OrganizationWasInserted)
	sc.Step(`^a contract with the id ([^"]*) was created for the organization having the id ([^"]*)$`, ContractWasInserted)
	sc.Step(`^the following SLIs are inserted in the database$`, func(table *godog.Table) (context.Context, error) {
		return CustomSlisWereInserted(ctx, table)
	})
	sc.Step(`^the SLIs should exist in the neo4j database in a consistent format$`, CustomSlisShouldExist)
	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		neo4jtest.CleanupAllData(ctx, driver)
		return ctx, nil
	})
}

func InitializeScenarioDefaultSLIsAreProperlyInserted(sc *godog.ScenarioContext) {
	sc.Step(`^a tenant was created$`, TenantWasInserted)
	sc.Step(`^the organization with the id ([^"]*) was created$`, OrganizationWasInserted)
	sc.Step(`^a contract with the id ([^"]*) was created for the organization having the id ([^"]*)$`, ContractWasInserted)
	sc.Step(`^(\d+) SLIs are inserted in the database for the contract ([^"]*)$`, SlisWereInserted)
	sc.Step(`^(\d+) should exist in the neo4j database$`, SlisShouldExist)
	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		neo4jtest.CleanupAllData(ctx, driver)
		return ctx, nil
	})
}
