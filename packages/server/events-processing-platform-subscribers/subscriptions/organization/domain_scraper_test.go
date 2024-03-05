package organization

import (
	"context"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"testing"

	"os"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
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
		neo4jtest.CleanupAllData(ctx, database.Driver)
	}
}

type mockAiModel struct{}

func (m *mockAiModel) Inference(ctx context.Context, text string) (string, error) {
	return `{"website": "https://www.customeros.ai", "companyName": "CustomerOS", "industry": "Enterprise SaaS", "industryGroup": "Software", "market": "Business to Business", "subIndustry": "CRM Software", "targetAudience": "SaaS companies", "valueProposition": "Grow with your best customers. See every experience everywhere. Create a success plan that delivers results."}`, nil
}

func TestWebScraping(t *testing.T) {
	ctx := context.Background()

	defer tearDownTestCase(ctx, testDatabase)(t)
	//_, driver := neo4jt.InitTestNeo4jDB()

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	organizationId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{Name: "org 1"})
	_ = neo4jtest.CreateLogEntryForOrganization(ctx, testDatabase.Driver, tenantName, organizationId, neo4jentity.LogEntryEntity{Content: "test content", StartedAt: utils.Now()})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization":  1,
		"LogEntry":      1,
		"TimelineEvent": 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"CREATED_BY": 0,
		"LOGGED":     1,
	})

	ds := NewDomainScraper(testLogger, &config.Config{}, testDatabase.Repositories, &mockAiModel{})
	scrapedContent, err := ds.Scrape("https://www.customeros.ai", tenantName, organizationId, true)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{"website": "https://www.customeros.ai", "companyName": "CustomerOS", "industry": "Enterprise SaaS", "industryGroup": "Software", "market": "Business to Business", "subIndustry": "CRM Software", "targetAudience": "SaaS companies", "valueProposition": "Grow with your best customers. See every experience everywhere. Create a success plan that delivers results."}
	require.Equal(t, expected["companyName"], scrapedContent.CompanyName)
	require.Equal(t, expected["website"], scrapedContent.Website)
	require.Equal(t, expected["industry"], scrapedContent.Industry)
}
