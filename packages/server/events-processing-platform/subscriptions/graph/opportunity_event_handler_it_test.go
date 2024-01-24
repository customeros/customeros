package graph

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestOpportunityEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	userIdOwner := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	userIdCreator := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, "sf")
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "User": 2, "ExternalSystem": 1, "Opportunity": 0})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create an OpportunityCreateEvent
	opportunityId := uuid.New().String()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	timeNow := utils.Now()
	opportunityData := model.OpportunityDataFields{
		Name:              "New Opportunity",
		Amount:            10000,
		InternalType:      model.NBO,
		ExternalType:      "TypeA",
		InternalStage:     model.OPEN,
		ExternalStage:     "Stage1",
		EstimatedClosedAt: &timeNow,
		OwnerUserId:       userIdOwner,
		CreatedByUserId:   userIdCreator,
		GeneralNotes:      "Some general notes about the opportunity",
		NextSteps:         "Next steps to proceed with the opportunity",
		OrganizationId:    orgId,
	}
	createEvent, err := event.NewOpportunityCreateEvent(
		opportunityAggregate,
		opportunityData,
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		commonmodel.ExternalSystem{
			ExternalSystemId: "sf",
			ExternalId:       "ext-id-1",
		},
		timeNow,
		timeNow,
	)
	require.Nil(t, err, "failed to create opportunity create event")

	// EXECUTE
	err = opportunityEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute opportunity create event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization":   1,
		"User":           2,
		"ExternalSystem": 1,
		"Opportunity":    1, "Opportunity_" + tenantName: 1})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, userIdOwner, "OWNS", opportunityId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, opportunityId, "CREATED_BY", userIdCreator)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, orgId, "HAS_OPPORTUNITY", opportunityId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, opportunityId, "IS_LINKED_WITH", "sf")

	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Opportunity_"+tenantName, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), opportunity.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, opportunity.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), opportunity.SourceOfTruth)
	require.Equal(t, timeNow, opportunity.CreatedAt)
	test.AssertRecentTime(t, opportunity.UpdatedAt)
	require.Equal(t, timeNow, *opportunity.EstimatedClosedAt)
	require.Equal(t, opportunityData.Name, opportunity.Name)
	require.Equal(t, opportunityData.Amount, opportunity.Amount)
	require.Equal(t, string(opportunityData.InternalType.StringEnumValue()), opportunity.InternalType)
	require.Equal(t, opportunityData.ExternalType, opportunity.ExternalType)
	require.Equal(t, string(opportunityData.InternalStage.StringEnumValue()), opportunity.InternalStage)
	require.Equal(t, opportunityData.ExternalStage, opportunity.ExternalStage)
	require.Equal(t, opportunityData.GeneralNotes, opportunity.GeneralNotes)
	require.Equal(t, opportunityData.NextSteps, opportunity.NextSteps)
}

func TestOpportunityEventHandler_OnCreateRenewal(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "Contract": 1, "Opportunity": 0})

	// prepare grpc mock
	calledEventsPlatformToRefreshRenewalSummary := false
	organizationServiceRefreshCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		RefreshRenewalSummary: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, org.AppSource)
			calledEventsPlatformToRefreshRenewalSummary = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
	}
	mocked_grpc.SetOrganizationCallbacks(&organizationServiceRefreshCallbacks)

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create an OpportunityCreateEvent
	opportunityId := uuid.New().String()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	timeNow := utils.Now()
	createEvent, err := event.NewOpportunityCreateRenewalEvent(
		opportunityAggregate,
		contractId,
		neo4jenum.RenewalLikelihoodLow.String(),
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
	)
	require.Nil(t, err, "failed to create opportunity create renewal event")

	// EXECUTE
	err = opportunityEventHandler.OnCreateRenewal(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute opportunity create event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":    1,
		"Opportunity": 1, "Opportunity_" + tenantName: 1})
	neo4jtest.AssertRelationships(ctx, t, testDatabase.Driver, contractId, []string{"HAS_OPPORTUNITY", "ACTIVE_RENEWAL"}, opportunityId)

	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Opportunity_"+tenantName, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), opportunity.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, opportunity.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), opportunity.SourceOfTruth)
	require.Equal(t, timeNow, opportunity.CreatedAt)
	test.AssertRecentTime(t, opportunity.UpdatedAt)
	require.Nil(t, opportunity.EstimatedClosedAt)
	require.Equal(t, "", opportunity.Name)
	require.Equal(t, "", opportunity.ExternalType)
	require.Equal(t, "", opportunity.ExternalStage)
	require.Equal(t, neo4jenum.OpportunityInternalTypeRenewal.String(), opportunity.InternalType)
	require.Equal(t, neo4jenum.OpportunityInternalStageOpen.String(), opportunity.InternalStage)
	require.Equal(t, neo4jenum.RenewalLikelihoodLow.String(), opportunity.RenewalDetails.RenewalLikelihood)

	require.True(t, calledEventsPlatformToRefreshRenewalSummary)
}

func TestOpportunityEventHandler_OnUpdateNextCycleDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen.String(),
	})
	updatedAt := utils.Now()
	renewedAt := updatedAt.AddDate(0, 6, 0) // 6 months later

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create an OpportunityUpdateNextCycleDateEvent
	updateEvent, err := event.NewOpportunityUpdateNextCycleDateEvent(
		aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId),
		updatedAt,
		&renewedAt,
	)
	require.Nil(t, err)

	// Execute the event handler
	err = opportunityEventHandler.OnUpdateNextCycleDate(ctx, updateEvent)
	require.Nil(t, err)

	// Assert Neo4j Node
	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelOpportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// Validate that the opportunity next cycle date is updated in the repository
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, updatedAt, opportunity.UpdatedAt)
	require.Equal(t, renewedAt, *opportunity.RenewalDetails.RenewedAt)
}

func TestOpportunityEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Name:      "test opportunity",
		Amount:    10000,
		MaxAmount: 20000,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelOpportunity: 1})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	now := utils.Now()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateEvent(opportunityAggregate,
		model.OpportunityDataFields{
			Name:      "updated opportunity",
			Amount:    30000,
			MaxAmount: 40000,
		},
		constants.SourceOpenline,
		commonmodel.ExternalSystem{},
		now,
		nil)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelOpportunity: 1, neo4jutil.NodeLabelOpportunity + "_" + tenantName: 1})

	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelOpportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, "updated opportunity", opportunity.Name)
	require.Equal(t, float64(30000), opportunity.Amount)
	require.Equal(t, float64(40000), opportunity.MaxAmount)
}

func TestOpportunityEventHandler_OnUpdate_OnlyAmountIsChangedByFieldsMask(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Name:   "test opportunity",
		Amount: 10000,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelOpportunity: 1})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	now := utils.Now()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateEvent(opportunityAggregate,
		model.OpportunityDataFields{
			Name:   "updated opportunity",
			Amount: 20000,
		},
		constants.SourceOpenline,
		commonmodel.ExternalSystem{},
		now,
		[]string{"amount"})
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelOpportunity: 1, neo4jutil.NodeLabelOpportunity + "_" + tenantName: 1})

	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Opportunity_"+tenantName, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, "test opportunity", opportunity.Name)
	require.Equal(t, float64(20000), opportunity.Amount)
}

func TestOpportunityEventHandler_OnUpdateRenewal_AmountAndRenewalChangedByUser(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Name:         "test opportunity",
		Amount:       10000,
		InternalType: neo4jenum.OpportunityInternalTypeRenewal.String(),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: "HIGH",
		},
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		RenewalCycle: neo4jenum.RenewalCycleMonthlyRenewal,
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelOpportunity: 1})

	// prepare grpc mock
	calledEventsPlatformToRefreshRenewalSummary, calledEventsPlatformToRefreshArr := false, false
	organizationServiceRefreshCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		RefreshRenewalSummary: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, org.AppSource)
			calledEventsPlatformToRefreshRenewalSummary = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
		RefreshArr: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, org.AppSource)
			calledEventsPlatformToRefreshArr = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
	}
	mocked_grpc.SetOrganizationCallbacks(&organizationServiceRefreshCallbacks)

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	now := utils.Now()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateRenewalEvent(opportunityAggregate,
		"MEDIUM",
		"some comments",
		"user-123",
		"openline",
		float64(10),
		now,
		[]string{},
		"user-123")
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdateRenewal(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelOpportunity: 1, neo4jutil.NodeLabelOpportunity + "_" + tenantName: 1})

	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelOpportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, "MEDIUM", opportunity.RenewalDetails.RenewalLikelihood)
	require.Equal(t, "user-123", opportunity.RenewalDetails.RenewalUpdatedByUserId)
	require.Equal(t, now, *opportunity.RenewalDetails.RenewalUpdatedByUserAt)
	require.Equal(t, float64(10), opportunity.Amount)
	require.Equal(t, "some comments", opportunity.Comments)

	require.True(t, calledEventsPlatformToRefreshRenewalSummary)
	require.True(t, calledEventsPlatformToRefreshArr)
}

func TestOpportunityEventHandler_OnUpdateRenewal_OnlyCommentsChangedByUser_DoNotUpdatePreviousUpdatedByUser(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstoret.NewTestAggregateStore()

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Name:         "test opportunity",
		Amount:       10000,
		Comments:     "no comments",
		InternalType: string(neo4jenum.OpportunityInternalTypeRenewal),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood:      "HIGH",
			RenewalUpdatedByUserId: "orig-user",
		},
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelOpportunity: 1})

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	now := utils.Now()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateRenewalEvent(opportunityAggregate,
		"HIGH",
		"some comments",
		"user-123",
		"openline",
		float64(10000),
		now,
		[]string{},
		"user-123")
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdateRenewal(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelOpportunity: 1, neo4jutil.NodeLabelOpportunity + "_" + tenantName: 1})

	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelOpportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, "HIGH", opportunity.RenewalDetails.RenewalLikelihood)
	require.Equal(t, "orig-user", opportunity.RenewalDetails.RenewalUpdatedByUserId)
	require.Nil(t, opportunity.RenewalDetails.RenewalUpdatedByUserAt)
	require.Equal(t, float64(10000), opportunity.Amount)
	require.Equal(t, "some comments", opportunity.Comments)

	// Check no events were generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 0, len(eventsMap))
}

func TestOpportunityEventHandler_OnUpdateRenewal_LikelihoodChangedByUser_GenerateEventToRecalculateArr(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, neo4jentity.ContractEntity{
		RenewalCycle: neo4jenum.RenewalCycleMonthlyRenewal,
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		Amount:        10000,
		InternalType:  string(neo4jenum.OpportunityInternalTypeRenewal),
		InternalStage: string(neo4jenum.OpportunityInternalStageOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood:      "HIGH",
			RenewalUpdatedByUserId: "orig-user",
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelOpportunity: 1,
		neo4jutil.NodeLabelContract:    1})

	// prepare grpc client
	calledEventsPlatformToUpdateOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateOpportunity: func(context context.Context, op *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, float64(0), op.Amount)
			require.Equal(t, float64(0), op.MaxAmount)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT}, op.FieldsMask)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, constants.SourceOpenline, op.SourceFields.Source)
			calledEventsPlatformToUpdateOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	now := utils.Now()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateRenewalEvent(opportunityAggregate,
		"MEDIUM",
		"Updated likelihood",
		"user-123",
		"openline",
		float64(10000),
		now,
		[]string{},
		"user-123")
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdateRenewal(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelOpportunity: 1, neo4jutil.NodeLabelOpportunity + "_" + tenantName: 1})

	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelOpportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, "MEDIUM", opportunity.RenewalDetails.RenewalLikelihood)
	require.Equal(t, "user-123", opportunity.RenewalDetails.RenewalUpdatedByUserId)
	require.Equal(t, now, *opportunity.RenewalDetails.RenewalUpdatedByUserAt)
	require.Equal(t, float64(10000), opportunity.Amount)
	require.Equal(t, "Updated likelihood", opportunity.Comments)

	// Verify Action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, neo4jenum.ActionRenewalLikelihoodUpdated, action.Type)
	require.Equal(t, "Renewal likelihood set to Medium", action.Content)
	require.Equal(t, fmt.Sprintf(`{"likelihood":"%s","reason":"%s"}`, "MEDIUM", "Updated likelihood"), action.Metadata)
	// Check extra properties
	props := utils.GetPropsFromNode(*actionDbNode)
	require.Equal(t, "Updated likelihood", utils.GetStringPropOrEmpty(props, "comments"))

	// check that event was invoked
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestOpportunityEventHandler_OnCloseWin(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalStage: string(neo4jenum.OpportunityInternalStageOpen),
	})
	now := utils.Now()

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	closeOpportunityEvent, err := event.NewOpportunityCloseWinEvent(aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId), now, now)
	require.Nil(t, err)

	// Execute the event handler
	err = opportunityEventHandler.OnCloseWin(ctx, closeOpportunityEvent)
	require.Nil(t, err)

	// Assert Neo4j Node
	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelOpportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// Validate that the opportunity next cycle date is updated in the repository
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, now, *opportunity.ClosedAt)
	require.Equal(t, string(neo4jenum.OpportunityInternalStageClosedWon), opportunity.InternalStage)
}

func TestOpportunityEventHandler_OnCloseLoose(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalStage: string(neo4jenum.OpportunityInternalStageOpen),
	})
	now := utils.Now()

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	closeOpportunityEvent, err := event.NewOpportunityCloseLooseEvent(aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId), now, now)
	require.Nil(t, err)

	// Execute the event handler
	err = opportunityEventHandler.OnCloseLoose(ctx, closeOpportunityEvent)
	require.Nil(t, err)

	// Assert Neo4j Node
	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelOpportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// Validate that the opportunity next cycle date is updated in the repository
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, now, opportunity.UpdatedAt)
	require.Equal(t, now, *opportunity.ClosedAt)
	require.Equal(t, string(neo4jenum.OpportunityInternalStageClosedLost), opportunity.InternalStage)
}

func TestOpportunityEventHandler_OnUpdateRenewal_ChangeOwner(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	userIdOwner := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	userIdCreator := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	userIdOwnerNew := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, "sf")
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1, "User": 3, "ExternalSystem": 1, "Opportunity": 0})

	// prepare grpc mock
	calledEventsPlatformToRefreshRenewalSummary := false
	organizationServiceRefreshCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		RefreshRenewalSummary: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, org.AppSource)
			calledEventsPlatformToRefreshRenewalSummary = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
	}
	mocked_grpc.SetOrganizationCallbacks(&organizationServiceRefreshCallbacks)

	// Prepare the event handler
	opportunityEventHandler := &OpportunityEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create an OpportunityCreateEvent
	opportunityId := uuid.New().String()
	opportunityAggregate := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	timeNow := utils.Now()
	opportunityData := model.OpportunityDataFields{
		Name:              "New Opportunity",
		Amount:            10000,
		InternalType:      model.NBO,
		ExternalType:      "TypeA",
		InternalStage:     model.OPEN,
		ExternalStage:     "Stage1",
		EstimatedClosedAt: &timeNow,
		OwnerUserId:       userIdOwner,
		CreatedByUserId:   userIdCreator,
		GeneralNotes:      "Some general notes about the opportunity",
		NextSteps:         "Next steps to proceed with the opportunity",
		OrganizationId:    orgId,
	}
	createEvent, err := event.NewOpportunityCreateEvent(
		opportunityAggregate,
		opportunityData,
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		commonmodel.ExternalSystem{
			ExternalSystemId: "sf",
			ExternalId:       "ext-id-1",
		},
		timeNow,
		timeNow,
	)
	require.Nil(t, err, "failed to create opportunity create event")

	// EXECUTE
	err = opportunityEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute opportunity create event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization":   1,
		"User":           3,
		"ExternalSystem": 1,
		"Opportunity":    1, "Opportunity_" + tenantName: 1})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, userIdOwner, "OWNS", opportunityId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, opportunityId, "CREATED_BY", userIdCreator)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, orgId, "HAS_OPPORTUNITY", opportunityId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, opportunityId, "IS_LINKED_WITH", "sf")

	opportunityDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Opportunity_"+tenantName, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)

	// verify opportunity
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, opportunityId, opportunity.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), opportunity.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, opportunity.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), opportunity.SourceOfTruth)
	require.Equal(t, timeNow, opportunity.CreatedAt)
	test.AssertRecentTime(t, opportunity.UpdatedAt)
	require.Equal(t, timeNow, *opportunity.EstimatedClosedAt)
	require.Equal(t, opportunityData.Name, opportunity.Name)
	require.Equal(t, opportunityData.Amount, opportunity.Amount)
	require.Equal(t, string(opportunityData.InternalType.StringEnumValue()), opportunity.InternalType)
	require.Equal(t, opportunityData.ExternalType, opportunity.ExternalType)
	require.Equal(t, string(opportunityData.InternalStage.StringEnumValue()), opportunity.InternalStage)
	require.Equal(t, opportunityData.ExternalStage, opportunity.ExternalStage)
	require.Equal(t, opportunityData.GeneralNotes, opportunity.GeneralNotes)
	require.Equal(t, opportunityData.NextSteps, opportunity.NextSteps)

	now := utils.Now()
	opportunityAggregate1 := aggregate.NewOpportunityAggregateWithTenantAndID(tenantName, opportunityId)
	updateEvent, err := event.NewOpportunityUpdateRenewalEvent(opportunityAggregate1,
		"MEDIUM",
		"Updated likelihood",
		"user-123",
		"openline",
		float64(10000),
		now,
		[]string{},
		userIdOwnerNew) //Changing the owner here
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = opportunityEventHandler.OnUpdateRenewal(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelOpportunity: 1, neo4jutil.NodeLabelOpportunity + "_" + tenantName: 1})
	//checking if the owner changed
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, userIdOwnerNew, "OWNS", opportunityId)

	opportunityDbNode1, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelOpportunity, opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode1)

	// verify opportunity
	opportunity1 := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode1)
	require.Equal(t, opportunityId, opportunity1.Id)

	require.True(t, calledEventsPlatformToRefreshRenewalSummary)

}
