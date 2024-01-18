package repository

import (
	"context"
	"github.com/cucumber/godog"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"testing"
)

var contextData map[string]interface{}

func TestFeatures(t *testing.T) {
	contextData = make(map[string]interface{})
	contextData["testingInstance"] = t
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^(\d+) SLIs were inserted in the database$`, SlisWereInserted)
	sc.Step(`^(\d+) should exist in the neo4j database$`, SlisShouldExist)
	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		//tearDownTestCase(ctx)
		neo4jtest.CleanupAllData(ctx, driver)
		return ctx, nil
	})
}
