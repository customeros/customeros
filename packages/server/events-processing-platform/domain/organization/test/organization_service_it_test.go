package test

import (
	"context"
	"github.com/google/uuid"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	organizationAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	organizationEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"testing"
	"time"
)

var testDatabase *test.TestDatabase
var dialFactory *grpc.TestDialFactoryImpl

func TestMain(m *testing.M) {
	myDatabase, shutdown := test.SetupTestDatabase()
	testDatabase = &myDatabase

	dialFactory = &grpc.TestDialFactoryImpl{}
	defer shutdown()

	os.Exit(m.Run())
}

func TestOrganizationsService_UpsertOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to emailEvents processing platform: %v", err)
	}
	organizationClient := organization_grpc_service.NewOrganizationGrpcServiceClient(grpcConnection)
	timeNow := time.Now().UTC()
	organizationId := uuid.New().String()
	tenant := "ziggy"
	response, err := organizationClient.UpsertOrganization(ctx, &organization_grpc_service.UpsertOrganizationGrpcRequest{
		Tenant:            tenant,
		Id:                organizationId,
		Name:              "Test Organization",
		Description:       "This is a organization description",
		Website:           "https://www.openline.ai",
		Employees:         int64(12),
		Market:            "B2B",
		Industry:          "Software",
		SubIndustry:       "sub-industry",
		IndustryGroup:     "industry-group",
		TargetAudience:    "target-audience",
		ValueProposition:  "value-proposition",
		LastFundingRound:  "Seed",
		LastFundingAmount: "1.000.000",
		IsPublic:          false,
		AppSource:         "unit-test",
		Source:            "N/A",
		SourceOfTruth:     "N/A",
		CreatedAt:         timestamppb.New(timeNow),
	})
	if err != nil {
		t.Errorf("Failed to create organization: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	aggregate := organizationAggregate.NewOrganizationAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[aggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, organizationEvents.OrganizationCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(organizationAggregate.OrganizationAggregateType)+"-"+tenant+"-"+organizationId, eventList[0].GetAggregateID())
	var eventData organizationEvents.OrganizationCreateEvent
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, "unit-test", eventData.AppSource)
	require.Equal(t, "N/A", eventData.Source)
	require.Equal(t, "N/A", eventData.SourceOfTruth)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "Test Organization", eventData.Name)
	require.Equal(t, "This is a organization description", eventData.Description)
	require.Equal(t, "https://www.openline.ai", eventData.Website)
	require.Equal(t, int64(12), eventData.Employees)
	require.Equal(t, "B2B", eventData.Market)
	require.Equal(t, "Software", eventData.Industry)
	require.Equal(t, "sub-industry", eventData.SubIndustry)
	require.Equal(t, "industry-group", eventData.IndustryGroup)
	require.Equal(t, "target-audience", eventData.TargetAudience)
	require.Equal(t, "value-proposition", eventData.ValueProposition)
	require.Equal(t, "Seed", eventData.LastFundingRound)
	require.Equal(t, "1.000.000", eventData.LastFundingAmount)
	require.Equal(t, false, eventData.IsPublic)
}

func TestOrganizationsService_LinkDomain(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to emailEvents processing platform: %v", err)
	}
	organizationClient := organization_grpc_service.NewOrganizationGrpcServiceClient(grpcConnection)
	organizationId := uuid.New().String()
	domain := "openline.ai"
	tenant := "ziggy"
	response, err := organizationClient.LinkDomainToOrganization(ctx, &organization_grpc_service.LinkDomainToOrganizationGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
		Domain:         domain,
	})
	if err != nil {
		t.Errorf("Failed to link domain: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	aggregate := organizationAggregate.NewOrganizationAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[aggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, organizationEvents.OrganizationLinkDomainV1, eventList[0].GetEventType())
	require.Equal(t, string(organizationAggregate.OrganizationAggregateType)+"-"+tenant+"-"+organizationId, eventList[0].GetAggregateID())
	var eventData organizationEvents.OrganizationLinkDomainEvent
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, domain, eventData.Domain)
}

func tearDownTestCase(ctx context.Context, database *test.TestDatabase) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jt.CleanupAllData(ctx, database.Driver)
	}
}
