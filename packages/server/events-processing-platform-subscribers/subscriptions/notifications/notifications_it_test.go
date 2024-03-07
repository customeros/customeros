package notifications

import (
	"context"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"os"
	"testing"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
)

const tenantName = "ziggy"

var testDatabase *test.TestDatabase
var testLogger = test.SetupTestLogger()

func TestMain(m *testing.M) {
	myDatabase, shutdown := test.SetupTestDatabase()
	testDatabase = &myDatabase

	defer shutdown()

	os.Exit(m.Run())
}

func tearDownTestCase(ctx context.Context, database *test.TestDatabase) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jtest.CleanupAllData(ctx, database.Driver)
	}
}
