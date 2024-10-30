package graph

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	eventcompletionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_completion"
	commonEvents "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"regexp"
	"testing"
	"time"
)

const customerOsIdPattern = `^C-[A-HJ-NP-Z2-9]{3}-[A-HJ-NP-Z2-9]{3}$`

func TestGraphOrganizationEventHandler_OnOrganizationCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 0,
		"User":         1, "User_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0})

	orgId := uuid.New().String()

	callbacks := mocked_grpc.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	mocked_grpc.SetEventCompletionServiceCallbacks(&callbacks)

	// prepare event handler
	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	event, err := events.NewOrganizationCreateEvent(orgAggregate, &model.OrganizationFields{
		ID: orgId,
		OrganizationDataFields: model.OrganizationDataFields{
			Name:         "test org",
			Relationship: "CUSTOMER",
			LeadSource:   "website",
		},
	}, now, now)
	require.Nil(t, err)
	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = event.SetMetadata(metadata)
	require.Nil(t, err)

	// EXECUTE
	err = orgEventHandler.OnOrganizationCreate(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "User_" + tenantName: 1,
		"Organization": 1, "Organization_" + tenantName: 1,
		"Action": 1, "Action_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"ACTION_ON": 1,
		"OWNS":      1,
	})

	orgDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, orgDbNode)

	// verify organization
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, "test org", organization.Name)
	require.Equal(t, now, organization.CreatedAt)
	require.NotNil(t, organization.UpdatedAt)
	require.Equal(t, string(neo4jenum.OnboardingStatusNotApplicable), organization.OnboardingDetails.Status)
	require.Nil(t, organization.OnboardingDetails.SortingOrder)
	require.Equal(t, neo4jenum.Customer, organization.Relationship)
	require.Equal(t, "website", organization.LeadSource)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, now, action.CreatedAt)
	require.Equal(t, neo4jenum.ActionCreated, action.Type)
	require.Equal(t, "", action.Content)
	require.Equal(t, "", action.Metadata)
}

func TestGraphOrganizationEventHandler_OnOrganizationShow(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: "test org",
		Hide: true,
	})
	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	event, err := events.NewShowOrganizationEventEvent(orgAggregate)
	require.Nil(t, err)
	err = orgEventHandler.OnOrganizationShow(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1})
	neo4jtest.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"Organization", "Organization_" + tenantName, "Tenant"})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	organization := neo4jmapper.MapDbNodeToOrganizationEntity(dbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, false, organization.Hide)
	require.NotEqual(t, "", organization.CustomerOsId)
	require.True(t, regexp.MustCompile(customerOsIdPattern).MatchString(organization.CustomerOsId), "Valid CustomerOsId should match the format")
}

func TestGraphOrganizationEventHandler_OnSocialAddedToOrganization_New(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	socialId := uuid.New().String()
	socialUrl := "https://www.facebook.com/organization"
	now := utils.Now()
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: "test org",
	})
	neo4jtest.CreateSocial(ctx, testDatabase.Driver, tenantName, neo4jentity.SocialEntity{
		Url: socialUrl,
	})
	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	sourceFields := commonEvents.Source{
		Source:        constants.SourceOpenline,
		SourceOfTruth: constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
	}
	event, err := events.NewOrganizationAddSocialEvent(orgAggregate, socialId, socialUrl, "alias", "ext-1", int64(123), sourceFields, now)
	require.Nil(t, err)
	err = orgEventHandler.OnSocialAddedToOrganization(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1, "Social": 2, "Social_" + tenantName: 2})
	neo4jtest.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"Organization", "Organization_" + tenantName, "Tenant", "Social", "Social_" + tenantName})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Social_"+tenantName, socialId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	social := neo4jmapper.MapDbNodeToSocialEntity(dbNode)
	require.Equal(t, socialId, social.Id)
	require.Equal(t, socialUrl, social.Url)
	require.Equal(t, "alias", social.Alias)
	require.Equal(t, "ext-1", social.ExternalId)
	require.Equal(t, int64(123), social.FollowersCount)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), social.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, social.AppSource)
	require.Equal(t, now, social.CreatedAt)
	test.AssertRecentTime(t, social.UpdatedAt)
}

func TestGraphOrganizationEventHandler_OnSocialAddedToOrganization_SocialUrlAlreadyExistsForOrg(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	socialUrl := "https://www.facebook.com/organization"
	now := utils.Now()
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: "test org",
	})
	existingSocialId := neo4jtest.CreateSocial(ctx, testDatabase.Driver, tenantName, neo4jentity.SocialEntity{
		Url:            socialUrl,
		Alias:          "existing alias",
		ExternalId:     "ext-1",
		FollowersCount: int64(5),
	})
	neo4jt.LinkSocial(ctx, testDatabase.Driver, existingSocialId, orgId)

	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	// prepare grpc mock
	callbacks := mocked_grpc.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	mocked_grpc.SetEventCompletionServiceCallbacks(&callbacks)

	sourceFields := commonEvents.Source{
		Source:        constants.SourceOpenline,
		SourceOfTruth: constants.SourceOpenline,
		AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
	}
	event, err := events.NewOrganizationAddSocialEvent(orgAggregate, existingSocialId, socialUrl, "alias", "ext-1", int64(100), sourceFields, now)
	require.Nil(t, err)
	err = orgEventHandler.OnSocialAddedToOrganization(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1, "Social": 1, "Social_" + tenantName: 1})
	neo4jtest.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"Organization", "Organization_" + tenantName, "Tenant", "Social", "Social_" + tenantName})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Social_"+tenantName, existingSocialId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	social := neo4jmapper.MapDbNodeToSocialEntity(dbNode)
	require.Equal(t, existingSocialId, social.Id)
	require.Equal(t, socialUrl, social.Url)
	require.Equal(t, "alias", social.Alias)
	require.Equal(t, "ext-1", social.ExternalId)
	require.Equal(t, int64(100), social.FollowersCount)
}

func TestGraphOrganizationEventHandler_OnLocationLinkedToOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	organizationName := "test_org_name"
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: organizationName,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1})
	dbNodeAfterOrganizationCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterOrganizationCreate)
	propsAfterOrganizationCreate := utils.GetPropsFromNode(*dbNodeAfterOrganizationCreate)
	require.Equal(t, organizationName, utils.GetStringPropOrEmpty(propsAfterOrganizationCreate, "name"))

	locationName := "test_location_name"
	locationId := neo4jtest.CreateLocation(ctx, testDatabase.Driver, tenantName, neo4jentity.LocationEntity{
		Name: locationName,
	})

	dbNodeAfterLocationCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterLocationCreate)
	propsAfterLocationCreate := utils.GetPropsFromNode(*dbNodeAfterLocationCreate)
	require.Equal(t, locationName, utils.GetStringPropOrEmpty(propsAfterLocationCreate, "name"))

	// prepare grpc mock
	callbacks := mocked_grpc.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	mocked_grpc.SetEventCompletionServiceCallbacks(&callbacks)

	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	event, err := events.NewOrganizationLinkLocationEvent(orgAggregate, locationId, now)
	require.Nil(t, err)
	err = orgEventHandler.OnLocationLinkedToOrganization(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "ASSOCIATED_WITH"), "Incorrect number of ASSOCIATED_WITH relationships in Neo4j")
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, orgId, "ASSOCIATED_WITH", locationId)
}

func TestGraphOrganizationEventHandler_OnRefreshArr(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId1 := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId1, neo4jentity.OpportunityEntity{
		Amount:       float64(10),
		MaxAmount:    float64(20),
		InternalType: neo4jenum.OpportunityInternalTypeRenewal,
	})
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId1, neo4jentity.OpportunityEntity{
		Amount:       float64(100),
		MaxAmount:    float64(200),
		InternalType: neo4jenum.OpportunityInternalTypeRenewal,
	})
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId2, neo4jentity.OpportunityEntity{
		Amount:       float64(1000),
		MaxAmount:    float64(2000),
		InternalType: neo4jenum.OpportunityInternalTypeRenewal,
	})
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId2, neo4jentity.OpportunityEntity{
		Amount:       float64(10000),
		MaxAmount:    float64(20000),
		InternalType: neo4jenum.OpportunityInternalTypeNBO,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Contract": 2, "Opportunity": 4})

	// prepare grpc mock
	callbacks := mocked_grpc.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	mocked_grpc.SetEventCompletionServiceCallbacks(&callbacks)

	// prepare event handler
	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	event, err := events.NewOrganizationRefreshArrEvent(orgAggregate)
	require.Nil(t, err)

	// EXECUTE
	err = orgEventHandler.OnRefreshArr(context.Background(), event)
	require.Nil(t, err)

	orgDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Organization", orgId)
	require.Nil(t, err)
	require.NotNil(t, orgDbNode)

	// verify organization
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, float64(1110), *organization.RenewalSummary.ArrForecast)
	require.Equal(t, float64(2220), *organization.RenewalSummary.MaxArrForecast)

	// Check no events were generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 0, len(eventsMap))
}

func TestGraphOrganizationEventHandler_OnRefreshRenewalSummary(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	tomorrow := utils.Now().Add(time.Duration(24) * time.Hour)
	afterTomorrow := utils.Now().Add(time.Duration(48) * time.Hour)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId1 := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId1, neo4jentity.OpportunityEntity{
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewedAt:         &tomorrow,
			RenewalLikelihood: "HIGH",
		},
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
	})
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId2, neo4jentity.OpportunityEntity{
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewedAt:         &afterTomorrow,
			RenewalLikelihood: "LOW",
		},
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
	})
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId2, neo4jentity.OpportunityEntity{
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewedAt:         &afterTomorrow,
			RenewalLikelihood: "ZERO",
		},
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageClosedWon,
	})
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId2, neo4jentity.OpportunityEntity{
		InternalType: neo4jenum.OpportunityInternalTypeNBO,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Contract": 2, "Opportunity": 4})

	// prepare event handler
	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	event, err := events.NewOrganizationRefreshRenewalSummaryEvent(orgAggregate)
	require.Nil(t, err)

	// EXECUTE
	err = orgEventHandler.OnRefreshRenewalSummaryV1(context.Background(), event)
	require.Nil(t, err)

	orgDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Organization", orgId)
	require.Nil(t, err)
	require.NotNil(t, orgDbNode)

	// verify organization
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, int64(20), *organization.RenewalSummary.RenewalLikelihoodOrder)
	require.Equal(t, "LOW", organization.RenewalSummary.RenewalLikelihood)
	require.Equal(t, utils.ToDatePtr(&tomorrow), organization.RenewalSummary.NextRenewalAt)

	// Check no events were generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 0, len(eventsMap))
}

func TestGraphOrganizationEventHandler_OnUpdateOnboardingStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "Olivia",
		LastName:  "Rhye",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: "test org",
	})
	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	// prepare grpc mock
	callbacks := mocked_grpc.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	mocked_grpc.SetEventCompletionServiceCallbacks(&callbacks)

	now := utils.Now()
	event, err := events.NewUpdateOnboardingStatusEvent(orgAggregate, "DONE", "Some comments", userId, "", now)
	require.Nil(t, err)
	err = orgEventHandler.OnUpdateOnboardingStatus(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization":                1,
		"Organization_" + tenantName:  1,
		"Action":                      1,
		"Action_" + tenantName:        1,
		"TimelineEvent":               1,
		"TimelineEvent_" + tenantName: 1,
	})
	neo4jtest.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"Organization", "Organization_" + tenantName, "Tenant",
		"Action", "Action_" + tenantName, "TimelineEvent", "TimelineEvent_" + tenantName, "User", "User_" + tenantName})

	// Check organization
	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(dbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, "DONE", organization.OnboardingDetails.Status)
	require.Equal(t, int64(constants.OnboardingStatus_Order_Done), *organization.OnboardingDetails.SortingOrder)
	require.Equal(t, "Some comments", organization.OnboardingDetails.Comments)
	require.Equal(t, now, *organization.OnboardingDetails.UpdatedAt)

	// Check action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionOnboardingStatusChanged, action.Type)
	require.Equal(t, "Olivia Rhye changed the onboarding status to Done", action.Content)
	require.Equal(t, fmt.Sprintf(`{"status":"%s","comments":"%s","userId":"%s","contractId":"%s"}`, "DONE", "Some comments", userId, ""), action.Metadata)
}

func TestGraphOrganizationEventHandler_OnUpdateOnboardingStatus_CausedByContractChange(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: "test org",
	})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})

	// prepare grpc mock
	callbacks := mocked_grpc.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	mocked_grpc.SetEventCompletionServiceCallbacks(&callbacks)

	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	now := utils.Now()
	event, err := events.NewUpdateOnboardingStatusEvent(orgAggregate, "NOT_STARTED", "Some comments", "", contractId, now)
	require.Nil(t, err)
	// EXECUTE
	err = orgEventHandler.OnUpdateOnboardingStatus(context.Background(), event)
	require.Nil(t, err)

	// Verify nodes
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Organization_" + tenantName: 1,
		"Contract": 1, "Contract_" + tenantName: 1,
		"Action": 1, "Action_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jtest.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"Organization", "Organization_" + tenantName, "Tenant",
		"Action", "Action_" + tenantName, "TimelineEvent", "TimelineEvent_" + tenantName, "Contract", "Contract_" + tenantName})

	// Verify Organization
	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(dbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, "NOT_STARTED", organization.OnboardingDetails.Status)
	require.Equal(t, "Some comments", organization.OnboardingDetails.Comments)
	require.Equal(t, now, *organization.OnboardingDetails.UpdatedAt)

	// Verify Contract
	dbNode, err = neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	contract := neo4jmapper.MapDbNodeToContractEntity(dbNode)
	require.Equal(t, contractId, contract.Id)
	require.True(t, contract.TriggeredOnboardingStatusChange)

	// Verify Action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionOnboardingStatusChanged, action.Type)
	require.Equal(t, "The onboarding status was automatically set to Not started", action.Content)
	require.Equal(t, fmt.Sprintf(`{"status":"%s","comments":"%s","userId":"%s","contractId":"%s"}`, "NOT_STARTED", "Some comments", "", contractId), action.Metadata)
	// Check extra properties
	props := utils.GetPropsFromNode(*actionDbNode)
	require.Equal(t, "Some comments", utils.GetStringPropOrEmpty(props, "comments"))
	require.Equal(t, "NOT_STARTED", utils.GetStringPropOrEmpty(props, "status"))
}

func TestGraphOrganizationEventHandler_OnCreateBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})

	// prepare grpc mock
	callbacks := mocked_grpc.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	mocked_grpc.SetEventCompletionServiceCallbacks(&callbacks)

	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	now := utils.Now()
	billingProfileId := uuid.New().String()
	event, err := events.NewBillingProfileCreateEvent(orgAggregate, billingProfileId, "Billing profile", "Tax id",
		commonEvents.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		}, now, now)
	require.Nil(t, err)
	err = orgEventHandler.OnCreateBillingProfile(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganization:                      1,
		model2.NodeLabelOrganization + "_" + tenantName:   1,
		model2.NodeLabelBillingProfile:                    1,
		model2.NodeLabelBillingProfile + "_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, orgId, "HAS_BILLING_PROFILE", billingProfileId)

	// Check billing profile
	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelBillingProfile, billingProfileId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	billingProfile := neo4jmapper.MapDbNodeToBillingProfileEntity(dbNode)
	require.Equal(t, billingProfileId, billingProfile.Id)
	require.Equal(t, "Billing profile", billingProfile.LegalName)
	require.Equal(t, "Tax id", billingProfile.TaxId)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), billingProfile.Source)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), billingProfile.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, billingProfile.AppSource)
	require.Equal(t, now, billingProfile.CreatedAt)
	test.AssertRecentTime(t, billingProfile.UpdatedAt)
}

func TestGraphOrganizationEventHandler_OnUpdateBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	billingProfileId := neo4jtest.CreateBillingProfileForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.BillingProfileEntity{
		LegalName: "Billing profile",
		TaxId:     "Tax id",
	})

	// prepare grpc mock
	callbacks := mocked_grpc.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	mocked_grpc.SetEventCompletionServiceCallbacks(&callbacks)

	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	now := utils.Now()
	event, err := events.NewBillingProfileUpdateEvent(orgAggregate, billingProfileId, "Updated name", "Updated tax id",
		now, []string{events.FieldMaskLegalName, events.FieldMaskTaxId})
	require.Nil(t, err)
	err = orgEventHandler.OnUpdateBillingProfile(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganization:   1,
		model2.NodeLabelBillingProfile: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, orgId, "HAS_BILLING_PROFILE", billingProfileId)

	// Check billing profile
	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelBillingProfile, billingProfileId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	billingProfile := neo4jmapper.MapDbNodeToBillingProfileEntity(dbNode)
	require.Equal(t, billingProfileId, billingProfile.Id)
	require.Equal(t, "Updated name", billingProfile.LegalName)
	require.Equal(t, "Updated tax id", billingProfile.TaxId)
	test.AssertRecentTime(t, billingProfile.UpdatedAt)
}

func TestGraphOrganizationEventHandler_OnEmailLinkedToBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	billingProfileId := neo4jtest.CreateBillingProfileForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.BillingProfileEntity{})
	existingEmailId := neo4jtest.CreateEmail(ctx, testDatabase.Driver, tenantName, neo4jentity.EmailEntity{})
	neo4jtest.LinkNodes(ctx, testDatabase.Driver, billingProfileId, existingEmailId, "HAS", map[string]interface{}{"primary": true})
	newEmailId := neo4jtest.CreateEmail(ctx, testDatabase.Driver, tenantName, neo4jentity.EmailEntity{})

	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	now := utils.Now()
	event, _ := events.NewLinkEmailToBillingProfileEvent(orgAggregate, billingProfileId, newEmailId, true, now)
	err := orgEventHandler.OnEmailLinkedToBillingProfile(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganization:   1,
		model2.NodeLabelBillingProfile: 1,
		model2.NodeLabelEmail:          2,
	})
	neo4jtest.AssertRelationshipWithProperties(ctx, t, testDatabase.Driver, billingProfileId, "HAS", existingEmailId, map[string]interface{}{"primary": false})
	neo4jtest.AssertRelationshipWithProperties(ctx, t, testDatabase.Driver, billingProfileId, "HAS", newEmailId, map[string]interface{}{"primary": true})

	// Check billing profile
	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelBillingProfile, billingProfileId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	billingProfile := neo4jmapper.MapDbNodeToBillingProfileEntity(dbNode)
	test.AssertRecentTime(t, billingProfile.UpdatedAt)
}

func TestGraphOrganizationEventHandler_OnEmailUnlinkedFromBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	billingProfileId := neo4jtest.CreateBillingProfileForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.BillingProfileEntity{})
	existingEmailId := neo4jtest.CreateEmail(ctx, testDatabase.Driver, tenantName, neo4jentity.EmailEntity{})
	neo4jtest.LinkNodes(ctx, testDatabase.Driver, billingProfileId, existingEmailId, "HAS", map[string]interface{}{"primary": true})

	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	now := utils.Now()
	event, _ := events.NewUnlinkEmailFromBillingProfileEvent(orgAggregate, billingProfileId, existingEmailId, now)
	err := orgEventHandler.OnEmailUnlinkedFromBillingProfile(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganization:   1,
		model2.NodeLabelBillingProfile: 1,
		model2.NodeLabelEmail:          1,
	})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"HAS": 0,
	})

	// Check billing profile
	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelBillingProfile, billingProfileId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	billingProfile := neo4jmapper.MapDbNodeToBillingProfileEntity(dbNode)
	test.AssertRecentTime(t, billingProfile.UpdatedAt)
}

func TestGraphOrganizationEventHandler_OnLocationLinkedToBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	billingProfileId := neo4jtest.CreateBillingProfileForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.BillingProfileEntity{})
	locationId := neo4jtest.CreateLocation(ctx, testDatabase.Driver, tenantName, neo4jentity.LocationEntity{})

	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	now := utils.Now()
	event, _ := events.NewLinkLocationToBillingProfileEvent(orgAggregate, billingProfileId, locationId, now)
	err := orgEventHandler.OnLocationLinkedToBillingProfile(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganization:   1,
		model2.NodeLabelBillingProfile: 1,
		model2.NodeLabelLocation:       1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, billingProfileId, "HAS", locationId)

	// Check billing profile
	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelBillingProfile, billingProfileId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	billingProfile := neo4jmapper.MapDbNodeToBillingProfileEntity(dbNode)
	test.AssertRecentTime(t, billingProfile.UpdatedAt)
}

func TestGraphOrganizationEventHandler_OnLocationUnlinkedFromBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	billingProfileId := neo4jtest.CreateBillingProfileForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.BillingProfileEntity{})
	existingLocationId := neo4jtest.CreateLocation(ctx, testDatabase.Driver, tenantName, neo4jentity.LocationEntity{})
	neo4jtest.LinkNodes(ctx, testDatabase.Driver, billingProfileId, existingLocationId, "HAS")

	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	now := utils.Now()
	event, _ := events.NewUnlinkLocationFromBillingProfileEvent(orgAggregate, billingProfileId, existingLocationId, now)
	err := orgEventHandler.OnLocationUnlinkedFromBillingProfile(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganization:   1,
		model2.NodeLabelBillingProfile: 1,
		model2.NodeLabelLocation:       1,
	})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"HAS": 0,
	})

	// Check billing profile
	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model2.NodeLabelBillingProfile, billingProfileId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	billingProfile := neo4jmapper.MapDbNodeToBillingProfileEntity(dbNode)
	test.AssertRecentTime(t, billingProfile.UpdatedAt)
}

func TestGraphOrganizationEventHandler_OnDomainUnlinkedFromOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.LinkDomainToOrganization(ctx, testDatabase.Driver, orgId, "openline.ai")

	orgEventHandler := &OrganizationEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	event, _ := events.NewOrganizationUnlinkDomainEvent(orgAggregate, "openline.ai")
	err := orgEventHandler.OnDomainUnlinkedFromOrganization(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganization: 1,
		model2.NodeLabelDomain:       1,
	})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"HAS_DOMAIN": 0,
	})
}
