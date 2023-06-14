package graph

import (
	"context"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"os"
	"testing"
)

const tenantName = "ziggy"

var testDatabase *neo4jt.TestDatabase

func TestMain(m *testing.M) {
	myDatabase, shutdown := neo4jt.SetupTestDatabase()
	testDatabase = &myDatabase

	defer shutdown()

	os.Exit(m.Run())
}

func tearDownTestCase(ctx context.Context, database *neo4jt.TestDatabase) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jt.CleanupAllData(ctx, database.Driver)
	}
}
