package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	orgaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestOrganizationsService_UpsertOrganization_NewOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to processing platform: %v", err)
	}
	organizationClient := organizationpb.NewOrganizationGrpcServiceClient(grpcConnection)
	timeNow := utils.Now()
	organizationId := uuid.New().String()
	tenant := "ziggy"
	response, err := organizationClient.UpsertOrganization(ctx, &organizationpb.UpsertOrganizationGrpcRequest{
		Tenant:             tenant,
		Id:                 organizationId,
		Name:               "Test Organization",
		Description:        "This is a organization description",
		Website:            "https://www.openline.ai",
		Employees:          int64(12),
		Market:             "B2B",
		Industry:           "Software",
		SubIndustry:        "sub-industry",
		IndustryGroup:      "industry-group",
		TargetAudience:     "target-audience",
		ValueProposition:   "value-proposition",
		LastFundingRound:   "Seed",
		LastFundingAmount:  "1.000.000",
		ReferenceId:        "100/200",
		Note:               "Some important notes",
		IsPublic:           false,
		YearFounded:        utils.ToPtr(int64(2019)),
		Headquarters:       "San Francisco, CA",
		EmployeeGrowthRate: "10%",
		SlackChannelId:     "channel-id",
		LogoUrl:            "https://www.openline.ai/logo.png",
		IconUrl:            "https://www.openline.ai/icon.png",
		SourceFields: &commonpb.SourceFields{
			AppSource: "event-processing-platform",
			Source:    "N/A",
		},
		CreatedAt:    timestamppb.New(timeNow),
		Relationship: "PROSPECT",
		Stage:        "LEAD",
		LeadSource:   "Email",
	})
	if err != nil {
		t.Errorf("Failed to create organization: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	aggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[aggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, orgevents.OrganizationCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(orgaggregate.OrganizationAggregateType)+"-"+tenant+"-"+organizationId, eventList[0].GetAggregateID())
	var eventData orgevents.OrganizationCreateEvent
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, "event-processing-platform", eventData.AppSource)
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
	require.Equal(t, "100/200", eventData.ReferenceId)
	require.Equal(t, "Some important notes", eventData.Note)
	require.Equal(t, false, eventData.IsPublic)
	require.Equal(t, utils.ToPtr(int64(2019)), eventData.YearFounded)
	require.Equal(t, "San Francisco, CA", eventData.Headquarters)
	require.Equal(t, "10%", eventData.EmployeeGrowthRate)
	require.Equal(t, "channel-id", eventData.SlackChannelId)
	require.Equal(t, "https://www.openline.ai/logo.png", eventData.LogoUrl)
	require.Equal(t, "https://www.openline.ai/icon.png", eventData.IconUrl)
	require.Equal(t, "PROSPECT", eventData.Relationship)
	require.Equal(t, "LEAD", eventData.Stage)
	require.Equal(t, "Email", eventData.LeadSource)
}

func TestOrganizationsService_UnlinkDomain(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to processing platform: %v", err)
	}
	organizationClient := organizationpb.NewOrganizationGrpcServiceClient(grpcConnection)
	organizationId := uuid.New().String()
	domain := "openline.ai"
	tenant := "ziggy"
	response, err := organizationClient.UnlinkDomainFromOrganization(ctx, &organizationpb.UnLinkDomainFromOrganizationGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
		Domain:         domain,
	})
	if err != nil {
		t.Errorf("Failed to un link domain: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	aggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[aggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, orgevents.OrganizationUnlinkDomainV1, eventList[0].GetEventType())
	require.Equal(t, string(orgaggregate.OrganizationAggregateType)+"-"+tenant+"-"+organizationId, eventList[0].GetAggregateID())
	var eventData orgevents.OrganizationUnlinkDomainEvent
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, domain, eventData.Domain)
}

func TestOrganizationService_UpdateOnboardingStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	organizationId := uuid.New().String()

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	organizationClient := organizationpb.NewOrganizationGrpcServiceClient(grpcConnection)

	// Grpc call
	response, err := organizationClient.UpdateOnboardingStatus(ctx, &organizationpb.UpdateOnboardingStatusGrpcRequest{
		Tenant:           tenant,
		OrganizationId:   organizationId,
		LoggedInUserId:   "user-id-123",
		OnboardingStatus: organizationpb.OnboardingStatus_ONBOARDING_STATUS_DONE,
		Comments:         "Some comments",
		AppSource:        "event-processing-platform",
	})
	require.Nil(t, err)

	// Assert response
	require.NotNil(t, response)
	require.Equal(t, organizationId, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[organizationAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, orgevents.OrganizationUpdateOnboardingStatusV1, eventList[0].GetEventType())
	require.Equal(t, string(orgaggregate.OrganizationAggregateType)+"-"+tenant+"-"+organizationId, eventList[0].GetAggregateID())

	var eventData orgevents.UpdateOnboardingStatusEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	test.AssertRecentTime(t, eventData.UpdatedAt)
	require.Equal(t, "user-id-123", eventData.UpdatedByUserId)
	require.Equal(t, string(model.OnboardingStatusDone), eventData.Status)
	require.Equal(t, "Some comments", eventData.Comments)
}

func TestOrganizationService_CreateBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	organizationId := uuid.New().String()

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	organizationClient := organizationpb.NewOrganizationGrpcServiceClient(grpcConnection)

	// Grpc call
	response, err := organizationClient.CreateBillingProfile(ctx, &organizationpb.CreateBillingProfileGrpcRequest{
		Tenant:           tenant,
		OrganizationId:   organizationId,
		LoggedInUserId:   "user-id-123",
		BillingProfileId: "",
		SourceFields: &commonpb.SourceFields{
			AppSource: "event-processing-platform",
			Source:    "N/A",
		},
		LegalName: "Test Billing Profile",
		TaxId:     "123456789",
	})
	require.Nil(t, err)

	// Assert response
	require.NotNil(t, response)
	require.NotEmpty(t, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, organizationId)
	eventList := eventsMap[organizationAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, orgevents.OrganizationCreateBillingProfileV1, eventList[0].GetEventType())
	require.Equal(t, string(orgaggregate.OrganizationAggregateType)+"-"+tenant+"-"+organizationId, eventList[0].GetAggregateID())

	var eventData orgevents.BillingProfileCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "Test Billing Profile", eventData.LegalName)
	require.Equal(t, "123456789", eventData.TaxId)
	test.AssertRecentTime(t, eventData.CreatedAt)
	test.AssertRecentTime(t, eventData.UpdatedAt)
	require.Equal(t, "N/A", eventData.SourceFields.Source)
	require.Equal(t, "event-processing-platform", eventData.SourceFields.AppSource)
}

func TestOrganizationService_LinkEmailToBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	organizationId := uuid.New().String()

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	organizationClient := organizationpb.NewOrganizationGrpcServiceClient(grpcConnection)

	// Grpc call
	response, err := organizationClient.LinkEmailToBillingProfile(ctx, &organizationpb.LinkEmailToBillingProfileGrpcRequest{
		Tenant:           tenant,
		OrganizationId:   organizationId,
		LoggedInUserId:   "user-id-123",
		BillingProfileId: "profile-123",
		EmailId:          "email-123",
		Primary:          true,
	})
	require.Nil(t, err)

	// Assert response
	require.NotNil(t, response)
	require.NotEmpty(t, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, organizationId)
	eventList := eventsMap[organizationAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, orgevents.OrganizationEmailLinkToBillingProfileV1, eventList[0].GetEventType())
	require.Equal(t, string(orgaggregate.OrganizationAggregateType)+"-"+tenant+"-"+organizationId, eventList[0].GetAggregateID())

	var eventData orgevents.LinkEmailToBillingProfileEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	test.AssertRecentTime(t, eventData.UpdatedAt)
	require.Equal(t, true, eventData.Primary)
	require.Equal(t, "profile-123", eventData.BillingProfileId)
	require.Equal(t, "email-123", eventData.EmailId)
}

func TestOrganizationService_UnlinkEmailFromBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	organizationId := uuid.New().String()

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	organizationClient := organizationpb.NewOrganizationGrpcServiceClient(grpcConnection)

	// Grpc call
	response, err := organizationClient.UnlinkEmailFromBillingProfile(ctx, &organizationpb.UnlinkEmailFromBillingProfileGrpcRequest{
		Tenant:           tenant,
		OrganizationId:   organizationId,
		LoggedInUserId:   "user-id-123",
		BillingProfileId: "profile-123",
		EmailId:          "email-123",
	})
	require.Nil(t, err)

	// Assert response
	require.NotNil(t, response)
	require.NotEmpty(t, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, organizationId)
	eventList := eventsMap[organizationAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, orgevents.OrganizationEmailUnlinkFromBillingProfileV1, eventList[0].GetEventType())
	require.Equal(t, string(orgaggregate.OrganizationAggregateType)+"-"+tenant+"-"+organizationId, eventList[0].GetAggregateID())

	var eventData orgevents.UnlinkEmailFromBillingProfileEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	test.AssertRecentTime(t, eventData.UpdatedAt)
	require.Equal(t, "profile-123", eventData.BillingProfileId)
	require.Equal(t, "email-123", eventData.EmailId)
}

func TestOrganizationService_LinkLocationToBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	organizationId := uuid.New().String()

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	organizationClient := organizationpb.NewOrganizationGrpcServiceClient(grpcConnection)

	// Grpc call
	response, err := organizationClient.LinkLocationToBillingProfile(ctx, &organizationpb.LinkLocationToBillingProfileGrpcRequest{
		Tenant:           tenant,
		OrganizationId:   organizationId,
		LoggedInUserId:   "user-id-123",
		BillingProfileId: "profile-123",
		LocationId:       "location-123",
	})
	require.Nil(t, err)

	// Assert response
	require.NotNil(t, response)
	require.NotEmpty(t, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, organizationId)
	eventList := eventsMap[organizationAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, orgevents.OrganizationLocationLinkToBillingProfileV1, eventList[0].GetEventType())
	require.Equal(t, string(orgaggregate.OrganizationAggregateType)+"-"+tenant+"-"+organizationId, eventList[0].GetAggregateID())

	var eventData orgevents.LinkLocationToBillingProfileEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	test.AssertRecentTime(t, eventData.UpdatedAt)
	require.Equal(t, "profile-123", eventData.BillingProfileId)
	require.Equal(t, "location-123", eventData.LocationId)
}

func TestOrganizationService_UnlinkLocationFromBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	tenant := "ziggy"
	organizationId := uuid.New().String()

	aggregateStore := eventstore.NewTestAggregateStore()

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	organizationClient := organizationpb.NewOrganizationGrpcServiceClient(grpcConnection)

	// Grpc call
	response, err := organizationClient.UnlinkLocationFromBillingProfile(ctx, &organizationpb.UnlinkLocationFromBillingProfileGrpcRequest{
		Tenant:           tenant,
		OrganizationId:   organizationId,
		LoggedInUserId:   "user-id-123",
		BillingProfileId: "profile-123",
		LocationId:       "location-123",
	})
	require.Nil(t, err)

	// Assert response
	require.NotNil(t, response)
	require.NotEmpty(t, response.Id)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	organizationAggregate := orgaggregate.NewOrganizationAggregateWithTenantAndID(tenant, organizationId)
	eventList := eventsMap[organizationAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, orgevents.OrganizationLocationUnlinkFromBillingProfileV1, eventList[0].GetEventType())
	require.Equal(t, string(orgaggregate.OrganizationAggregateType)+"-"+tenant+"-"+organizationId, eventList[0].GetAggregateID())

	var eventData orgevents.UnlinkLocationFromBillingProfileEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err)

	require.Equal(t, tenant, eventData.Tenant)
	test.AssertRecentTime(t, eventData.UpdatedAt)
	require.Equal(t, "profile-123", eventData.BillingProfileId)
	require.Equal(t, "location-123", eventData.LocationId)
}
