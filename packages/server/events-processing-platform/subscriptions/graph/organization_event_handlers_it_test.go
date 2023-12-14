package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	opportunitymodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

const customerOsIdPattern = `^C-[A-HJ-NP-Z2-9]{3}-[A-HJ-NP-Z2-9]{3}$`

func TestGraphOrganizationEventHandler_OnOrganizationCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 0,
		"User":         1, "User_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0})

	orgId := uuid.New().String()

	// prepare grpc mock
	lastTouchpointInvoked := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		RefreshLastTouchpoint: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, org.AppSource)
			lastTouchpointInvoked = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
	}
	mocked_grpc.SetOrganizationCallbacks(&organizationServiceCallbacks)

	// prepare event handler
	orgEventHandler := &OrganizationEventHandler{
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	event, err := events.NewOrganizationCreateEvent(orgAggregate, &model.OrganizationFields{
		ID: orgId,
		OrganizationDataFields: model.OrganizationDataFields{
			Name: "test org",
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

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 1, "User_" + tenantName: 1,
		"Organization": 1, "Organization_" + tenantName: 1,
		"Action": 1, "Action_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})
	neo4jt.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"ACTION_ON": 1,
		"OWNS":      1,
	})

	orgDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, orgDbNode)

	// verify organization
	organization := graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, "test org", organization.Name)
	require.Equal(t, now, organization.CreatedAt)
	require.NotNil(t, organization.UpdatedAt)
	require.Equal(t, string(entity.OnboardingStatusNotApplicable), organization.OnboardingDetails.Status)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, now, action.CreatedAt)
	require.Equal(t, entity.ActionCreated, action.Type)
	require.Equal(t, "", action.Content)
	require.Equal(t, "", action.Metadata)

	// Check refresh last touch point
	require.Truef(t, lastTouchpointInvoked, "RefreshLastTouchpoint was not invoked")
}

func TestGraphOrganizationEventHandler_OnOrganizationHide(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
		Hide: false,
	})
	orgEventHandler := &OrganizationEventHandler{
		repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	event, err := events.NewHideOrganizationEventEvent(orgAggregate)
	require.Nil(t, err)
	err = orgEventHandler.OnOrganizationHide(context.Background(), event)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1})
	neo4jt.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"Organization", "Organization_" + tenantName, "Tenant"})

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	organization := graph_db.MapDbNodeToOrganizationEntity(*dbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, true, organization.Hide)
}

func TestGraphOrganizationEventHandler_OnOrganizationShow(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
		Hide: true,
	})
	orgEventHandler := &OrganizationEventHandler{
		repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	event, err := events.NewShowOrganizationEventEvent(orgAggregate)
	require.Nil(t, err)
	err = orgEventHandler.OnOrganizationShow(context.Background(), event)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1})
	neo4jt.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"Organization", "Organization_" + tenantName, "Tenant"})

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	organization := graph_db.MapDbNodeToOrganizationEntity(*dbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, false, organization.Hide)
	require.NotEqual(t, "", organization.CustomerOsId)
	require.True(t, regexp.MustCompile(customerOsIdPattern).MatchString(organization.CustomerOsId), "Valid CustomerOsId should match the format")
}

func TestGraphOrganizationEventHandler_OnSocialAddedToOrganization_New(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	socialId := uuid.New().String()
	socialUrl := "https://www.facebook.com/organization"
	platformName := "facebook"
	now := utils.Now()
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
	})
	neo4jt.CreateSocial(ctx, testDatabase.Driver, tenantName, entity.SocialEntity{
		Url: socialUrl,
	})
	orgEventHandler := &OrganizationEventHandler{
		repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	event, err := events.NewOrganizationAddSocialEvent(orgAggregate, socialId, platformName, socialUrl, constants.SourceOpenline, constants.SourceOpenline, constants.AppSourceEventProcessingPlatform, now, now)
	require.Nil(t, err)
	err = orgEventHandler.OnSocialAddedToOrganization(context.Background(), event)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1, "Social": 2, "Social_" + tenantName: 2})
	neo4jt.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"Organization", "Organization_" + tenantName, "Tenant", "Social", "Social_" + tenantName})

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Social_"+tenantName, socialId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	social := graph_db.MapDbNodeToSocialEntity(*dbNode)
	require.Equal(t, socialId, social.Id)
	require.Equal(t, socialUrl, social.Url)
	require.Equal(t, platformName, social.PlatformName)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), social.SourceFields.Source)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), social.SourceFields.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, social.SourceFields.AppSource)
	require.Equal(t, now, social.CreatedAt)
	require.Equal(t, now, social.UpdatedAt)
}

func TestGraphOrganizationEventHandler_OnSocialAddedToOrganization_SocialUrlAlreadyExistsForOrg_NoChanges(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	socialId := uuid.New().String()
	socialUrl := "https://www.facebook.com/organization"
	platformName := "facebook"
	now := utils.Now()
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
	})
	existingSocialId := neo4jt.CreateSocial(ctx, testDatabase.Driver, tenantName, entity.SocialEntity{
		Url:          socialUrl,
		PlatformName: platformName,
	})
	neo4jt.LinkSocial(ctx, testDatabase.Driver, existingSocialId, orgId)

	orgEventHandler := &OrganizationEventHandler{
		repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)

	event, err := events.NewOrganizationAddSocialEvent(orgAggregate, socialId, "other platform name", socialUrl, constants.SourceOpenline, constants.SourceOpenline, constants.AppSourceEventProcessingPlatform, now, now)
	require.Nil(t, err)
	err = orgEventHandler.OnSocialAddedToOrganization(context.Background(), event)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1, "Social": 1, "Social_" + tenantName: 1})
	neo4jt.AssertNeo4jLabels(ctx, t, testDatabase.Driver, []string{"Organization", "Organization_" + tenantName, "Tenant", "Social", "Social_" + tenantName})

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Social_"+tenantName, existingSocialId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	social := graph_db.MapDbNodeToSocialEntity(*dbNode)
	require.Equal(t, existingSocialId, social.Id)
	require.Equal(t, socialUrl, social.Url)
	require.Equal(t, platformName, social.PlatformName)
}

func TestGraphOrganizationEventHandler_OnLocationLinkedToOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)

	organizationName := "test_org_name"
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: organizationName,
	})

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1})
	dbNodeAfterOrganizationCreate, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterOrganizationCreate)
	propsAfterOrganizationCreate := utils.GetPropsFromNode(*dbNodeAfterOrganizationCreate)
	require.Equal(t, organizationName, utils.GetStringPropOrEmpty(propsAfterOrganizationCreate, "name"))

	locationName := "test_location_name"
	locationId := neo4jt.CreateLocation(ctx, testDatabase.Driver, tenantName, entity.LocationEntity{
		Name: locationName,
	})

	dbNodeAfterLocationCreate, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterLocationCreate)
	propsAfterLocationCreate := utils.GetPropsFromNode(*dbNodeAfterLocationCreate)
	require.Equal(t, locationName, utils.GetStringPropOrEmpty(propsAfterLocationCreate, "name"))

	orgEventHandler := &OrganizationEventHandler{
		repositories: testDatabase.Repositories,
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	event, err := events.NewOrganizationLinkLocationEvent(orgAggregate, locationId, now)
	require.Nil(t, err)
	err = orgEventHandler.OnLocationLinkedToOrganization(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "ASSOCIATED_WITH"), "Incorrect number of ASSOCIATED_WITH relationships in Neo4j")
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, orgId, "ASSOCIATED_WITH", locationId)
}

func TestGraphOrganizationEventHandler_OnRefreshArr(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	contractId1 := neo4jt.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, entity.ContractEntity{})
	contractId2 := neo4jt.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, entity.ContractEntity{})
	opportunityIdRenewal1_1 := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Amount:       float64(10),
		MaxAmount:    float64(20),
		InternalType: string(opportunitymodel.OpportunityInternalTypeStringRenewal),
	})
	opportunityIdRenewal2_1 := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Amount:       float64(100),
		MaxAmount:    float64(200),
		InternalType: string(opportunitymodel.OpportunityInternalTypeStringRenewal),
	})
	opportunityIdRenewal2_2 := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Amount:       float64(1000),
		MaxAmount:    float64(2000),
		InternalType: string(opportunitymodel.OpportunityInternalTypeStringRenewal),
	})
	opportunityIdNbo2_3 := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Amount:       float64(10000),
		MaxAmount:    float64(20000),
		InternalType: string(opportunitymodel.OpportunityInternalTypeStringNBO),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId1, opportunityIdRenewal1_1, true)
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId2, opportunityIdRenewal2_1, true)
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId2, opportunityIdRenewal2_2, true)
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId2, opportunityIdNbo2_3, false)
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Contract": 2, "Opportunity": 4})

	// prepare event handler
	orgEventHandler := &OrganizationEventHandler{
		repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	event, err := events.NewOrganizationRefreshArrEvent(orgAggregate)
	require.Nil(t, err)

	// EXECUTE
	err = orgEventHandler.OnRefreshArr(context.Background(), event)
	require.Nil(t, err)

	orgDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization", orgId)
	require.Nil(t, err)
	require.NotNil(t, orgDbNode)

	// verify organization
	organization := graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)
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
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	contractId1 := neo4jt.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, entity.ContractEntity{})
	contractId2 := neo4jt.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, entity.ContractEntity{})
	opportunityIdRenewal1_1 := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		RenewalDetails: entity.RenewalDetails{
			RenewedAt:         &tomorrow,
			RenewalLikelihood: "HIGH",
		},
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	opportunityIdRenewal2_1 := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		RenewalDetails: entity.RenewalDetails{
			RenewedAt:         &afterTomorrow,
			RenewalLikelihood: "LOW",
		},
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	opportunityIdRenewal2_2 := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		RenewalDetails: entity.RenewalDetails{
			RenewedAt:         &afterTomorrow,
			RenewalLikelihood: "ZERO",
		},
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringClosedWon),
	})
	opportunityIdNbo2_3 := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType: string(opportunitymodel.OpportunityInternalTypeStringNBO),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId1, opportunityIdRenewal1_1, true)
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId2, opportunityIdRenewal2_1, true)
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId2, opportunityIdRenewal2_2, true)
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId2, opportunityIdNbo2_3, false)
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Contract": 2, "Opportunity": 4})

	// prepare event handler
	orgEventHandler := &OrganizationEventHandler{
		repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewCommandHandlers(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	event, err := events.NewOrganizationRefreshRenewalSummaryEvent(orgAggregate)
	require.Nil(t, err)

	// EXECUTE
	err = orgEventHandler.OnRefreshRenewalSummary(context.Background(), event)
	require.Nil(t, err)

	orgDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization", orgId)
	require.Nil(t, err)
	require.NotNil(t, orgDbNode)

	// verify organization
	organization := graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, int64(20), *organization.RenewalSummary.RenewalLikelihoodOrder)
	require.Equal(t, "LOW", organization.RenewalSummary.RenewalLikelihood)
	require.Equal(t, tomorrow, *organization.RenewalSummary.NextRenewalAt)

	// Check no events were generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 0, len(eventsMap))
}
