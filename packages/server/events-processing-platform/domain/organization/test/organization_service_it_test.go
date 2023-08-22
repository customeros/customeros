package test

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	organizationAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	organizationEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
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
		t.Fatalf("Failed to connect to processing platform: %v", err)
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
		t.Fatalf("Failed to connect to processing platform: %v", err)
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

func TestOrganizationsService_UpdateRenewalLikelihood(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to processing platform: %v", err)
	}
	organizationClient := organization_grpc_service.NewOrganizationGrpcServiceClient(grpcConnection)

	organizationId := uuid.New().String()
	tenant := "openline"
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenant)
	neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenant, entity.OrganizationEntity{
		ID: organizationId,
	})

	response, err := organizationClient.UpdateOrganizationRenewalLikelihood(ctx, &organization_grpc_service.OrganizationRenewalLikelihoodRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
		Likelihood:     organization_grpc_service.Likelihood_HIGH,
		Comment:        utils.StringPtr("test comment"),
		UserId:         "user-123",
	})
	if err != nil {
		t.Errorf("Failed to update renewal likelihood: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	aggregate := organizationAggregate.NewOrganizationAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[aggregate.ID]
	require.Equal(t, 2, len(eventList))

	aggregateId := string(organizationAggregate.OrganizationAggregateType) + "-" + tenant + "-" + organizationId
	require.Equal(t, organizationEvents.OrganizationCreateV1, eventList[0].GetEventType())
	require.Equal(t, aggregateId, eventList[0].GetAggregateID())

	require.Equal(t, organizationEvents.OrganizationUpdateRenewalLikelihoodV1, eventList[1].GetEventType())
	require.Equal(t, aggregateId, eventList[1].GetAggregateID())
	var eventData organizationEvents.OrganizationUpdateRenewalLikelihoodEvent
	if err := eventList[1].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, tenant, eventData.Tenant)
	require.NotNil(t, eventData.Comment)
	require.Equal(t, "test comment", *eventData.Comment)
	require.Equal(t, "user-123", eventData.UpdatedBy)
	require.Equal(t, models.RenewalLikelihoodHIGH, eventData.RenewalLikelihood)
	test.AssertRecentTime(t, eventData.UpdatedAt)
}

func TestOrganizationsService_UpdateRenewalForecast(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to processing platform: %v", err)
	}
	organizationClient := organization_grpc_service.NewOrganizationGrpcServiceClient(grpcConnection)

	organizationId := uuid.New().String()
	tenant := "openline"
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenant)
	neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenant, entity.OrganizationEntity{
		ID: organizationId,
	})

	response, err := organizationClient.UpdateOrganizationRenewalForecast(ctx, &organization_grpc_service.OrganizationRenewalForecastRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
		Amount:         utils.Float64Ptr(100),
		Comment:        utils.StringPtr("test comment"),
		UserId:         "user-123",
	})
	if err != nil {
		t.Errorf("Failed to update renewal forecast: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	aggregate := organizationAggregate.NewOrganizationAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[aggregate.ID]
	require.Equal(t, 2, len(eventList))

	aggregateId := string(organizationAggregate.OrganizationAggregateType) + "-" + tenant + "-" + organizationId
	require.Equal(t, organizationEvents.OrganizationCreateV1, eventList[0].GetEventType())
	require.Equal(t, aggregateId, eventList[0].GetAggregateID())

	require.Equal(t, organizationEvents.OrganizationUpdateRenewalForecastV1, eventList[1].GetEventType())
	require.Equal(t, aggregateId, eventList[1].GetAggregateID())
	var eventData organizationEvents.OrganizationUpdateRenewalForecastEvent
	if err := eventList[1].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, tenant, eventData.Tenant)
	require.NotNil(t, eventData.Comment)
	require.Equal(t, "test comment", *eventData.Comment)
	require.Equal(t, "user-123", eventData.UpdatedBy)
	require.Equal(t, utils.Float64Ptr(100), eventData.Amount)
	require.Nil(t, eventData.PotentialAmount)
	test.AssertRecentTime(t, eventData.UpdatedAt)
}

func TestOrganizationsService_UpdateBillingDetails(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to processing platform: %v", err)
	}
	organizationClient := organization_grpc_service.NewOrganizationGrpcServiceClient(grpcConnection)

	organizationId := uuid.New().String()
	tenant := "openline"
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenant)
	neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenant, entity.OrganizationEntity{
		ID: organizationId,
	})
	now := utils.Now()

	response, err := organizationClient.UpdateOrganizationBillingDetails(ctx, &organization_grpc_service.OrganizationBillingDetailsRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
		Amount:         utils.Float64Ptr(100),
		UserId:         "user-123",
		Frequency:      utils.ToPtr(organization_grpc_service.Frequency_WEEKLY),
		RenewalCycle:   utils.ToPtr(organization_grpc_service.Frequency_MONTHLY),
		CycleStart:     utils.ConvertTimeToTimestampPtr(utils.TimePtr(now)),
	})
	if err != nil {
		t.Errorf("Failed to update billing details: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	aggregate := organizationAggregate.NewOrganizationAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[aggregate.ID]
	require.Equal(t, 2, len(eventList))

	aggregateId := string(organizationAggregate.OrganizationAggregateType) + "-" + tenant + "-" + organizationId
	require.Equal(t, organizationEvents.OrganizationCreateV1, eventList[0].GetEventType())
	require.Equal(t, aggregateId, eventList[0].GetAggregateID())

	require.Equal(t, organizationEvents.OrganizationUpdateBillingDetailsV1, eventList[1].GetEventType())
	require.Equal(t, aggregateId, eventList[1].GetAggregateID())
	var eventData organizationEvents.OrganizationUpdateBillingDetailsEvent
	if err := eventList[1].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, utils.Float64Ptr(100), eventData.Amount)
	require.Equal(t, "user-123", eventData.UpdatedBy)
	require.Equal(t, "WEEKLY", eventData.Frequency)
	require.Equal(t, "MONTHLY", eventData.RenewalCycle)
	require.Equal(t, now, *eventData.RenewalCycleStart)
	require.Nil(t, eventData.RenewalCycleNext)
}

func tearDownTestCase(ctx context.Context, database *test.TestDatabase) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jt.CleanupAllData(ctx, database.Driver)
	}
}
