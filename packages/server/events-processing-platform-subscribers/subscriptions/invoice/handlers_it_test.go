package invoice

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	"os"
	"testing"
)

const tenantName = "ziggy"

var testDatabase *test.TestDatabase
var testLogger = test.SetupTestLogger()
var testMockedGrpcClient = test.SetupMockedTestGrpcClient()

var ctx = common.WithCustomContext(context.Background(), &common.CustomContext{Tenant: tenantName})

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
