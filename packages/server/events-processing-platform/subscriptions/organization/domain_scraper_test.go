package organization

import (
	"context"
	"testing"

	"os"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
)

const tenantName = "openline"

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
		neo4jt.CleanupAllData(ctx, database.Driver)
	}
}

func TestWebScraping(t *testing.T) {
	ctx := context.TODO()
	// cfg, err := config.InitConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	defer tearDownTestCase(ctx, testDatabase)(t)
	_, driver := neo4jt.InitTestNeo4jDB()

	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{Name: "org 1"})
	_ = neo4jt.CreateLogEntryForOrg(ctx, driver, tenantName, organizationId, entity.LogEntryEntity{Content: "test content", StartedAt: utils.Now()})

	neo4jt.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Organization":  1,
		"LogEntry":      1,
		"TimelineEvent": 1})
	neo4jt.AssertNeo4jRelationCount(ctx, t, driver, map[string]int{
		"CREATED_BY": 0,
		"LOGGED":     1,
	})

	ds := NewDomainScraper(testLogger, &config.Config{}, testDatabase.Repositories)
	scrapedContent, err := ds.Scrape("https://www.customeros.ai", tenantName, organizationId)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{"website": "https://www.customeros.ai", "companyName": "CustomerOS", "industry": "Enterprise SaaS", "industryGroup": "Software", "market": "Business to Business", "subIndustry": "CRM Software", "targetAudience": "SaaS companies", "valueProposition": "Grow with your best customers. See every experience everywhere. Create a success plan that delivers results."}
	require.Equal(t, expected["companyName"], scrapedContent.CompanyName)
	require.Equal(t, expected["website"], scrapedContent.Website)
	require.Equal(t, expected["industry"], scrapedContent.Industry)
}
