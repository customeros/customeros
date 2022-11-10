package resolver

import (
	"context"
	"fmt"
	"github.com.openline-ai.customer-os-analytics-api/common"
	"github.com.openline-ai.customer-os-analytics-api/config"
	"github.com.openline-ai.customer-os-analytics-api/dataloader"
	"github.com.openline-ai.customer-os-analytics-api/graph/generated"
	"github.com.openline-ai.customer-os-analytics-api/graph/model"
	"github.com.openline-ai.customer-os-analytics-api/integration_tests"
	"github.com.openline-ai.customer-os-analytics-api/repository"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"testing"
)

var (
	dbContainer testcontainers.Container
	db          *config.StorageDB
	c           *client.Client
)

func TestMain(m *testing.M) {
	dbContainer, db = integration_tests.InitTestDB()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			log.Fatal("Error during container termination")
		}
	}(dbContainer, context.Background())

	prepareClient()

	os.Exit(m.Run())
}

func prepareClient() {
	repositoryContainer := repository.InitRepositories(db.GormDB)
	graphResolver := NewResolver(repositoryContainer)
	loader := dataloader.NewDataLoader(repositoryContainer)
	customCtx := &common.CustomContext{
		Tenant: "openline",
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphResolver}))
	dataloaderSrv := dataloader.Middleware(loader, srv)
	h := common.CreateContext(customCtx, dataloaderSrv)
	c = client.New(h)
}

func prepareTestDatabase(resourceFolder string) {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db.SqlDB),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(fmt.Sprintf("testdata/%s", resourceFolder)),
	)
	if err != nil {
		log.Fatal("Error creating test fixtures")
	}
	if err = fixtures.Load(); err != nil {
		log.Panicf("Error loading test fixtures: %v", err.Error())
	}
}

func getQuery(resourceFolder string) string {
	b, err := os.ReadFile(fmt.Sprintf("testdata/%s/query.txt", resourceFolder))
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

func TestQueryApplicationsOnly(t *testing.T) {
	var testResourcesFolder = "applicationsOnly"
	prepareTestDatabase(testResourcesFolder)
	query := getQuery(testResourcesFolder)

	var resp struct {
		Applications []model.Application
	}

	c.MustPost(query, &resp)
	require.Equal(t, 2, len(resp.Applications))
	require.Equal(t, "public_app_tracker", resp.Applications[0].TrackerName)
	require.Equal(t, "public_app_2_tracker", resp.Applications[1].TrackerName)
}

func TestQuerySessionsWithPageViews(t *testing.T) {
	var testResourcesFolder = "sessionsWithPageViews"
	prepareTestDatabase(testResourcesFolder)
	query := getQuery(testResourcesFolder)

	var resp struct {
		Application model.Application
	}

	c.MustPost(query, &resp)
	require.Equal(t, int64(2), resp.Application.Sessions.TotalElements)
	require.Equal(t, 1, resp.Application.Sessions.TotalPages)
	require.Equal(t, 2, len(resp.Application.Sessions.Content))
	require.Equal(t, "session-1", resp.Application.Sessions.Content[0].ID)
	require.Equal(t, "session-2", resp.Application.Sessions.Content[1].ID)
	require.Equal(t, 2, len(resp.Application.Sessions.Content[0].PageViews))
	require.Equal(t, 1, len(resp.Application.Sessions.Content[1].PageViews))
	require.Equal(t, "page-1", resp.Application.Sessions.Content[0].PageViews[0].ID)
	require.Equal(t, "page-2", resp.Application.Sessions.Content[0].PageViews[1].ID)
	require.Equal(t, "page-3", resp.Application.Sessions.Content[1].PageViews[0].ID)
	require.Equal(t, "title-1", resp.Application.Sessions.Content[0].PageViews[0].Title)
	require.Equal(t, "title-2", resp.Application.Sessions.Content[0].PageViews[1].Title)
	require.Equal(t, "title-3", resp.Application.Sessions.Content[1].PageViews[0].Title)
	require.Equal(t, "path-1", resp.Application.Sessions.Content[0].PageViews[0].Path)
	require.Equal(t, "path-2", resp.Application.Sessions.Content[0].PageViews[1].Path)
	require.Equal(t, "path-3", resp.Application.Sessions.Content[1].PageViews[0].Path)
	require.Equal(t, 1, resp.Application.Sessions.Content[0].PageViews[0].Order)
	require.Equal(t, 2, resp.Application.Sessions.Content[0].PageViews[1].Order)
	require.Equal(t, 1, resp.Application.Sessions.Content[1].PageViews[0].Order)
	require.Equal(t, 10, resp.Application.Sessions.Content[0].PageViews[0].EngagedTime)
	require.Equal(t, 20, resp.Application.Sessions.Content[0].PageViews[1].EngagedTime)
	require.Equal(t, 30, resp.Application.Sessions.Content[1].PageViews[0].EngagedTime)
}

func TestQuerySessionDetails(t *testing.T) {
	var testResourcesFolder = "sessionDetails"
	prepareTestDatabase(testResourcesFolder)
	query := getQuery(testResourcesFolder)

	var resp struct {
		Application model.Application
	}

	c.MustPost(query, &resp)
	require.Equal(t, 1, len(resp.Application.Sessions.Content))
	require.Equal(t, "session-1", resp.Application.Sessions.Content[0].ID)
	require.Equal(t, "Some Country", resp.Application.Sessions.Content[0].Country)
	require.Equal(t, "The City", resp.Application.Sessions.Content[0].City)
	require.Equal(t, "Region", resp.Application.Sessions.Content[0].Region)
	require.Equal(t, "Google", resp.Application.Sessions.Content[0].ReferrerSource)
	require.Equal(t, "utm campaign", resp.Application.Sessions.Content[0].UtmCampaign)
	require.Equal(t, "utm content", resp.Application.Sessions.Content[0].UtmContent)
	require.Equal(t, "utm medium", resp.Application.Sessions.Content[0].UtmMedium)
	require.Equal(t, "utm source", resp.Application.Sessions.Content[0].UtmSource)
	require.Equal(t, "utm network", resp.Application.Sessions.Content[0].UtmNetwork)
	require.Equal(t, "utm term", resp.Application.Sessions.Content[0].UtmTerm)
	require.Equal(t, "Apple", resp.Application.Sessions.Content[0].DeviceBrand)
	require.Equal(t, "Iphone", resp.Application.Sessions.Content[0].DeviceName)
	require.Equal(t, "Mobile", resp.Application.Sessions.Content[0].DeviceClass)
	require.Equal(t, "Chrome", resp.Application.Sessions.Content[0].AgentName)
	require.Equal(t, "106", resp.Application.Sessions.Content[0].AgentVersion)
	require.Equal(t, "Mac", resp.Application.Sessions.Content[0].OperatingSystem)
	require.Equal(t, "10", resp.Application.Sessions.Content[0].OsVersionMajor)
	require.Equal(t, "15", resp.Application.Sessions.Content[0].OsVersionMinor)
	require.Equal(t, 1000, resp.Application.Sessions.Content[0].EngagedTime)
}
