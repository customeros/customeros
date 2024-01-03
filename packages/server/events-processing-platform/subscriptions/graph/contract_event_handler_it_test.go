package graph

import (
	"context"
	"github.com/google/uuid"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	opportunitymodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContractEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	userIdCreator := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{})
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, "sf")
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "User": 1, "ExternalSystem": 1, "Contract": 0})

	// Prepare the event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ContractCreateEvent
	contractId := uuid.New().String()
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	timeNow := utils.Now()
	createEvent, err := event.NewContractCreateEvent(
		contractAggregate,
		model.ContractDataFields{
			Name:             "New Contract",
			ContractUrl:      "http://contract.url",
			OrganizationId:   orgId,
			CreatedByUserId:  userIdCreator,
			ServiceStartedAt: &timeNow,
			SignedAt:         &timeNow,
			RenewalCycle:     model.MonthlyRenewal,
			Status:           model.Live,
		},
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
	require.Nil(t, err)

	// prepare grpc mock for onboarding status update
	calledEventsPlatformForOnboardingStatusChange := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		UpdateOnboardingStatus: func(context context.Context, org *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, org.AppSource)
			require.Equal(t, organizationpb.OnboardingStatus_ONBOARDING_STATUS_NOT_STARTED, org.OnboardingStatus)
			require.Equal(t, "", org.LoggedInUserId)
			require.Equal(t, "", org.Comments)
			require.Equal(t, contractId, org.CausedByContractId)
			calledEventsPlatformForOnboardingStatusChange = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
	}
	mocked_grpc.SetOrganizationCallbacks(&organizationServiceCallbacks)

	calledEventsPlatformToCreateRenewalOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		CreateRenewalOpportunity: func(context context.Context, op *opportunitypb.CreateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, "", op.LoggedInUserId)
			require.Equal(t, contractId, op.ContractId)
			require.Nil(t, op.CreatedAt)
			require.Nil(t, op.UpdatedAt)
			require.Equal(t, opportunitypb.RenewalLikelihood_HIGH_RENEWAL, op.RenewalLikelihood)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			calledEventsPlatformToCreateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: "some-opportunity-id",
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// Execute
	err = contractEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute contract create event handler")

	// Verify
	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1,
		"User":         1,
		"Contract":     1, "Contract_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "CREATED_BY", userIdCreator)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, orgId, "HAS_CONTRACT", contractId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "IS_LINKED_WITH", "sf")

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// Verify contract
	contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "New Contract", contract.Name)
	require.Equal(t, "http://contract.url", contract.ContractUrl)
	require.Equal(t, model.Live.String(), contract.Status)
	require.Equal(t, model.MonthlyRenewal.String(), contract.RenewalCycle)
	require.True(t, timeNow.Equal(contract.CreatedAt.UTC()))
	require.True(t, timeNow.Equal(contract.UpdatedAt.UTC()))
	require.True(t, timeNow.Equal(*contract.ServiceStartedAt))
	require.True(t, timeNow.Equal(*contract.SignedAt))
	require.Nil(t, contract.EndedAt)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), contract.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, contract.AppSource)

	// Verify events platform was called
	require.True(t, calledEventsPlatformForOnboardingStatusChange)
	require.True(t, calledEventsPlatformToCreateRenewalOpportunity)
}

func TestContractEventHandler_OnUpdate_FrequencySet(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Name:        "test contract",
		ContractUrl: "http://contract.url",
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1})

	// prepare grpc mock
	calledEventsPlatformCreateRenewalOpportunity := false
	opportunityServiceRefreshCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		CreateRenewalOpportunity: func(context context.Context, op *opportunitypb.CreateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, contractId, op.ContractId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Nil(t, op.CreatedAt)
			require.Nil(t, op.UpdatedAt)
			calledEventsPlatformCreateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityServiceRefreshCallbacks)

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	now := utils.Now()
	yesterday := now.AddDate(0, 0, -1)
	daysAgo2 := now.AddDate(0, 0, -2)
	tomorrow := now.AddDate(0, 0, 1)
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name:             "test contract updated",
			ContractUrl:      "http://contract.url/updated",
			ServiceStartedAt: &yesterday,
			SignedAt:         &daysAgo2,
			EndedAt:          &tomorrow,
			RenewalCycle:     model.MonthlyRenewal,
			Status:           model.Live,
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Contract_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "test contract updated", contract.Name)
	require.Equal(t, "http://contract.url/updated", contract.ContractUrl)
	require.Equal(t, model.Live.String(), contract.Status)
	require.Equal(t, model.MonthlyRenewal.String(), contract.RenewalCycle)
	require.True(t, now.Equal(contract.UpdatedAt))
	require.True(t, yesterday.Equal(*contract.ServiceStartedAt))
	require.True(t, daysAgo2.Equal(*contract.SignedAt))
	require.True(t, tomorrow.Equal(*contract.EndedAt))
	require.Equal(t, entity.DataSource(constants.SourceOpenline), contract.SourceOfTruth)

	// Verify call to events platform
	require.True(t, calledEventsPlatformCreateRenewalOpportunity)
}

func TestContractEventHandler_OnUpdate_FrequencyNotChanged(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle: string(model.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

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

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name:         "test contract updated",
			RenewalCycle: model.MonthlyRenewal,
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		utils.Now())
	require.Nil(t, err)

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	// Verify
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestContractEventHandler_OnUpdate_FrequencyChanged(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle: string(model.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

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

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name:         "test contract updated",
			RenewalCycle: model.AnnuallyRenewal,
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		utils.Now())
	require.Nil(t, err)

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	// Verify
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestContractEventHandler_OnUpdate_FrequencyRemoved(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, entity.ContractEntity{
		RenewalCycle: string(model.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

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

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name:         "test contract updated",
			RenewalCycle: model.None,
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		utils.Now())
	require.Nil(t, err)

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	// Verify
	neo4jtest.AssertRelationships(ctx, t, testDatabase.Driver, contractId, []string{"SUSPENDED_RENEWAL", "HAS_OPPORTUNITY"}, opportunityId)
	opportunityDbNode, _ := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Opportunity", opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, "SUSPENDED", opportunity.InternalStage)

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)
	contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, "", contract.RenewalCycle)

	// verify call to events platform
	require.True(t, calledEventsPlatformToRefreshRenewalSummary)
	require.True(t, calledEventsPlatformToRefreshArr)
}

func TestContractEventHandler_OnUpdate_ServiceStartDateChanged(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	now := utils.Now()
	yesterday := now.AddDate(0, 0, -1)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
		ServiceStartedAt: &yesterday,
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare grpc client
	calledEventsPlatformToUpdateOpportunity := false
	calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate := false
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
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.AppSource)
			require.NotNil(t, op.RenewedAt)
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name:             "test contract updated",
			RenewalCycle:     model.MonthlyRenewal,
			ServiceStartedAt: &now,
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		utils.Now())
	require.Nil(t, err)

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	// Verify
	require.True(t, calledEventsPlatformToUpdateOpportunity)
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)
}

func TestContractEventHandler_OnUpdate_EndDateSet(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Name:         "test contract",
		ContractUrl:  "http://contract.url",
		RenewalCycle: string(model.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewedAt:         utils.Ptr(utils.Now().AddDate(0, 0, 20)),
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Opportunity": 1})

	// prepare grpc mock
	calledEventsPlatformToUpdateRenewalOpportunity := false
	calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate := false
	calledEventsPlatformToUpdateOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunity: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, opportunitypb.RenewalLikelihood_ZERO_RENEWAL, op.RenewalLikelihood)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD}, op.FieldsMask)
			calledEventsPlatformToUpdateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.AppSource)
			require.NotNil(t, op.RenewedAt)
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
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

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	now := utils.Now()
	yesterday := now.AddDate(0, 0, -1)
	daysAgo2 := now.AddDate(0, 0, -2)
	tomorrow := now.AddDate(0, 0, 1)
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name:             "test contract updated",
			ContractUrl:      "http://contract.url/updated",
			ServiceStartedAt: &yesterday,
			SignedAt:         &daysAgo2,
			EndedAt:          &tomorrow,
			RenewalCycle:     model.MonthlyRenewal,
			Status:           model.Live,
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Contract_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "test contract updated", contract.Name)
	require.Equal(t, "http://contract.url/updated", contract.ContractUrl)
	require.Equal(t, model.Live.String(), contract.Status)
	require.Equal(t, model.MonthlyRenewal.String(), contract.RenewalCycle)
	require.True(t, now.Equal(contract.UpdatedAt))
	require.True(t, yesterday.Equal(*contract.ServiceStartedAt))
	require.True(t, daysAgo2.Equal(*contract.SignedAt))
	require.True(t, tomorrow.Equal(*contract.EndedAt))
	require.Equal(t, entity.DataSource(constants.SourceOpenline), contract.SourceOfTruth)

	// Verify event platform was called
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunity)
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestContractEventHandler_OnUpdate_CurrentSourceOpenline_UpdateSourceNonOpenline(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	now := utils.Now()
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Name:             "test contract",
		ContractUrl:      "http://contract.url",
		Status:           "DRAFT",
		RenewalCycle:     "ANNUALLY",
		ServiceStartedAt: &now,
		SignedAt:         &now,
		EndedAt:          &now,
		SourceOfTruth:    constants.SourceOpenline,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1})

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	yesterday := now.AddDate(0, 0, -1)
	daysAgo2 := now.AddDate(0, 0, -2)
	tomorrow := now.AddDate(0, 0, 1)
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name:             "test contract updated",
			ContractUrl:      "http://contract.url/updated",
			ServiceStartedAt: &yesterday,
			SignedAt:         &daysAgo2,
			EndedAt:          &tomorrow,
			RenewalCycle:     model.MonthlyRenewal,
			Status:           model.Live,
		},
		commonmodel.ExternalSystem{},
		"hubspot",
		now)
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Contract_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "test contract", contract.Name)
	require.Equal(t, "http://contract.url", contract.ContractUrl)
	require.Equal(t, model.Draft.String(), contract.Status)
	require.Equal(t, model.AnnuallyRenewal.String(), contract.RenewalCycle)
	require.True(t, now.Equal(contract.UpdatedAt))
	require.True(t, now.Equal(*contract.ServiceStartedAt))
	require.True(t, now.Equal(*contract.SignedAt))
	require.True(t, now.Equal(*contract.EndedAt))
	require.Equal(t, entity.DataSource(constants.SourceOpenline), contract.SourceOfTruth)
}

func TestContractEventHandler_OnUpdateStatus_Ended(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, entity.ContractEntity{
		Name:   "test contract",
		Status: string(model.ContractStatusStringDraft),
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"Contract": 1, "Contract_" + tenantName: 1, "Action": 0, "TimelineEvent": 0})

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	now := utils.Now()
	event, err := event.NewContractUpdateStatusEvent(contractAggregate, string(model.ContractStatusStringEnded), &now, nil)
	require.Nil(t, err)

	// prepare grpc mock for onboarding status update
	calledEventsPlatformForOnboardingStatusChange := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		UpdateOnboardingStatus: func(context context.Context, org *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, org.AppSource)
			require.Equal(t, organizationpb.OnboardingStatus_ONBOARDING_STATUS_NOT_STARTED, org.OnboardingStatus)
			require.Equal(t, "", org.LoggedInUserId)
			require.Equal(t, "", org.Comments)
			require.Equal(t, contractId, org.CausedByContractId)
			calledEventsPlatformForOnboardingStatusChange = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
	}
	mocked_grpc.SetOrganizationCallbacks(&organizationServiceCallbacks)

	// EXECUTE
	err = contractEventHandler.OnUpdateStatus(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"Contract": 1, "Contract_" + tenantName: 1,
		"Action": 1, "Action_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)

	// verify grpc was called
	require.True(t, calledEventsPlatformForOnboardingStatusChange)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, neo4jentity.ActionContractStatusUpdated, action.Type)
	require.Equal(t, "test contract has ended", action.Content)
	require.Equal(t, `{"status":"ENDED","contract-name":"test contract","comment":"test contract is now ENDED"}`, action.Metadata)
}

func TestContractEventHandler_OnUpdateStatus_Live(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, entity.ContractEntity{
		Name:   "test contract",
		Status: string(model.ContractStatusStringDraft),
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"Contract": 1, "Contract_" + tenantName: 1, "Action": 0, "TimelineEvent": 0})

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	now := utils.Now()
	event, err := event.NewContractUpdateStatusEvent(contractAggregate, string(model.ContractStatusStringLive), &now, nil)
	require.Nil(t, err)

	// prepare grpc mock for onboarding status update
	calledEventsPlatformForOnboardingStatusChange := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		UpdateOnboardingStatus: func(context context.Context, org *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, org.AppSource)
			require.Equal(t, organizationpb.OnboardingStatus_ONBOARDING_STATUS_NOT_STARTED, org.OnboardingStatus)
			require.Equal(t, "", org.LoggedInUserId)
			require.Equal(t, "", org.Comments)
			require.Equal(t, contractId, org.CausedByContractId)
			calledEventsPlatformForOnboardingStatusChange = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
	}
	mocked_grpc.SetOrganizationCallbacks(&organizationServiceCallbacks)

	// EXECUTE
	err = contractEventHandler.OnUpdateStatus(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"Contract": 1, "Contract_" + tenantName: 1,
		"Action": 1, "Action_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)

	// verify grpc was called
	require.True(t, calledEventsPlatformForOnboardingStatusChange)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, neo4jentity.ActionContractStatusUpdated, action.Type)
	require.Equal(t, "test contract is now live", action.Content)
	require.Equal(t, `{"status":"LIVE","contract-name":"test contract","comment":"test contract is now LIVE"}`, action.Metadata)
}
