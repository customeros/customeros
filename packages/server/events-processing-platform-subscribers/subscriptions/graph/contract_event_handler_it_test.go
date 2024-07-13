package graph

import (
	"context"
	"github.com/google/uuid"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContractEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	userIdCreator := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{})
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
			Name:                 "New Contract",
			ContractUrl:          "http://contract.url",
			OrganizationId:       orgId,
			CreatedByUserId:      userIdCreator,
			ServiceStartedAt:     &timeNow,
			SignedAt:             &timeNow,
			LengthInMonths:       int64(1),
			BillingCycleInMonths: 1,
			Currency:             neo4jenum.CurrencyUSD.String(),
			InvoicingStartDate:   &timeNow,
			AutoRenew:            true,
			Check:                true,
			DueDays:              30,
			Country:              "US",
			Approved:             true,
		},
		events.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
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
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, org.AppSource)
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
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, op.SourceFields.AppSource)
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
	contract := mapper.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "New Contract", contract.Name)
	require.Equal(t, "http://contract.url", contract.ContractUrl)
	require.Equal(t, neo4jenum.ContractStatusLive, contract.ContractStatus)
	require.Equal(t, int64(1), contract.LengthInMonths)
	require.True(t, timeNow.Equal(contract.CreatedAt.UTC()))
	test.AssertRecentTime(t, contract.UpdatedAt)
	require.True(t, utils.ToDate(timeNow).Equal(*contract.ServiceStartedAt))
	require.Equal(t, utils.ToDatePtr(&timeNow), contract.SignedAt)
	require.True(t, utils.ToDatePtr(&timeNow).Equal(*contract.InvoicingStartDate))
	require.Equal(t, neo4jenum.CurrencyUSD, contract.Currency)
	require.Equal(t, int64(1), contract.BillingCycleInMonths)
	require.Nil(t, contract.EndedAt)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), contract.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, contract.AppSource)
	require.True(t, contract.AutoRenew)
	require.True(t, contract.Check)
	require.Equal(t, int64(30), contract.DueDays)
	require.Equal(t, "US", contract.Country)
	require.True(t, contract.Approved)

	// Verify events platform was called
	require.False(t, calledEventsPlatformForOnboardingStatusChange)
	require.True(t, calledEventsPlatformToCreateRenewalOpportunity)
}

func TestContractEventHandler_OnUpdate_FrequencySet(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
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
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, op.SourceFields.AppSource)
			require.Nil(t, op.CreatedAt)
			require.Nil(t, op.UpdatedAt)
			calledEventsPlatformCreateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityServiceRefreshCallbacks)

	calledEventsPlatformForOnboardingStatusChange := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		UpdateOnboardingStatus: func(context context.Context, org *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, org.AppSource)
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
			LengthInMonths:   int64(1),
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		now,
		[]string{event.FieldMaskLengthInMonths, event.FieldMaskName, event.FieldMaskServiceStartedAt, event.FieldMaskSignedAt, event.FieldMaskEndedAt, event.FieldMaskContractURL})
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Contract_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := mapper.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "test contract updated", contract.Name)
	require.Equal(t, "http://contract.url/updated", contract.ContractUrl)
	require.Equal(t, neo4jenum.ContractStatusDraft, contract.ContractStatus)
	require.Equal(t, int64(1), contract.LengthInMonths)
	test.AssertRecentTime(t, contract.UpdatedAt)
	require.True(t, utils.ToDate(yesterday).Equal(*contract.ServiceStartedAt))
	require.Equal(t, utils.ToDate(daysAgo2), *contract.SignedAt)
	require.Equal(t, utils.ToDate(tomorrow), *contract.EndedAt)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), contract.SourceOfTruth)

	// Verify call to events platform
	require.True(t, calledEventsPlatformCreateRenewalOpportunity)
	require.False(t, calledEventsPlatformForOnboardingStatusChange)
}

func TestContractEventHandler_OnUpdate_FrequencyNotChanged(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		LengthInMonths: 1,
	})
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
	})

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
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, op.SourceFields.AppSource)
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
			Name:           "test contract updated",
			LengthInMonths: int64(1),
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		utils.Now(),
		[]string{})
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
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		LengthInMonths: 1,
	})
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
	})

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
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, op.SourceFields.AppSource)
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
			Name:           "test contract updated",
			LengthInMonths: int64(12),
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		utils.Now(),
		[]string{})
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
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		LengthInMonths: 1,
	})
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
	})

	// prepare grpc mock
	calledEventsPlatformToRefreshRenewalSummary, calledEventsPlatformToRefreshArr := false, false
	organizationServiceRefreshCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		RefreshRenewalSummary: func(context context.Context, org *organizationpb.RefreshRenewalSummaryGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, org.AppSource)
			calledEventsPlatformToRefreshRenewalSummary = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
		RefreshArr: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, org.AppSource)
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
			Name:           "test contract updated",
			LengthInMonths: int64(0),
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		utils.Now(),
		[]string{event.FieldMaskLengthInMonths, event.FieldMaskName})
	require.Nil(t, err)

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	// Verify
	neo4jtest.AssertRelationships(ctx, t, testDatabase.Driver, contractId, []string{"SUSPENDED_RENEWAL", "HAS_OPPORTUNITY"}, opportunityId)
	opportunityDbNode, _ := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Opportunity", opportunityId)
	require.Nil(t, err)
	require.NotNil(t, opportunityDbNode)
	opportunity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
	require.Equal(t, neo4jenum.OpportunityInternalStageSuspended, opportunity.InternalStage)

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)
	contract := mapper.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, int64(0), contract.LengthInMonths)

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
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		LengthInMonths:   1,
		ServiceStartedAt: &yesterday,
	})
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
	})

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
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, op.SourceFields.AppSource)
			require.Equal(t, constants.SourceOpenline, op.SourceFields.Source)
			calledEventsPlatformToUpdateOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, op.AppSource)
			require.NotNil(t, op.RenewedAt)
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare grpc mock for onboarding status update
	calledEventsPlatformForOnboardingStatusChange := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		UpdateOnboardingStatus: func(context context.Context, org *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, org.AppSource)
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
			LengthInMonths:   int64(1),
			ServiceStartedAt: &now,
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		utils.Now(),
		[]string{})
	require.Nil(t, err)

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	// Verify
	require.True(t, calledEventsPlatformToUpdateOpportunity)
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)
	require.False(t, calledEventsPlatformForOnboardingStatusChange)
}

func TestContractEventHandler_OnUpdate_EndDateSet(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Name:           "test contract",
		ContractUrl:    "http://contract.url",
		LengthInMonths: 1,
		Approved:       true,
	})
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewedAt:         utils.Ptr(utils.Now().AddDate(0, 0, 20)),
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Opportunity": 1})

	// prepare grpc mock
	calledEventsPlatformToUpdateRenewalOpportunity := false
	calledEventsPlatformToUpdateOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunity: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, op.SourceFields.AppSource)
			require.Equal(t, opportunitypb.RenewalLikelihood_ZERO_RENEWAL, op.RenewalLikelihood)
			require.Equal(t, int64(0), op.RenewalAdjustedRate)
			require.Equal(t, []opportunitypb.OpportunityMaskField{
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_ADJUSTED_RATE,
			}, op.FieldsMask)
			calledEventsPlatformToUpdateRenewalOpportunity = true
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
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, op.SourceFields.AppSource)
			require.Equal(t, constants.SourceOpenline, op.SourceFields.Source)
			calledEventsPlatformToUpdateOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare grpc mock for onboarding status update
	calledEventsPlatformForOnboardingStatusChange := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		UpdateOnboardingStatus: func(context context.Context, org *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, org.AppSource)
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
			LengthInMonths:   int64(1),
			AutoRenew:        false,
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		now,
		[]string{})
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Contract_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := mapper.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "test contract updated", contract.Name)
	require.Equal(t, "http://contract.url/updated", contract.ContractUrl)
	require.Equal(t, neo4jenum.ContractStatusLive, contract.ContractStatus)
	require.Equal(t, int64(1), contract.LengthInMonths)
	test.AssertRecentTime(t, contract.UpdatedAt)
	require.True(t, utils.ToDate(yesterday).Equal(*contract.ServiceStartedAt))
	require.Equal(t, utils.ToDate(daysAgo2), *contract.SignedAt)
	require.Equal(t, utils.ToDate(tomorrow), *contract.EndedAt)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), contract.SourceOfTruth)

	// Verify event platform was called
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunity)
	require.True(t, calledEventsPlatformToUpdateOpportunity)
	require.False(t, calledEventsPlatformForOnboardingStatusChange)
}

func TestContractEventHandler_OnUpdate_CurrentSourceOpenline_UpdateSourceNonOpenline(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	now := utils.Now()
	tomorrow := now.AddDate(0, 0, 1)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Name:             "test contract",
		ContractUrl:      "http://contract.url",
		LengthInMonths:   12,
		ContractStatus:   neo4jenum.ContractStatusEnded,
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
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			Name:             "test contract updated",
			ContractUrl:      "http://contract.url/updated",
			ServiceStartedAt: &yesterday,
			SignedAt:         &daysAgo2,
			EndedAt:          &tomorrow,
			LengthInMonths:   int64(1),
		},
		commonmodel.ExternalSystem{},
		"hubspot",
		now,
		[]string{})
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Contract_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := mapper.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "test contract", contract.Name)
	require.Equal(t, "http://contract.url", contract.ContractUrl)
	require.Equal(t, neo4jenum.ContractStatusEnded, contract.ContractStatus)
	require.Equal(t, int64(12), contract.LengthInMonths)
	test.AssertRecentTime(t, contract.UpdatedAt)
	require.Equal(t, utils.ToDate(now), *contract.ServiceStartedAt)
	require.Equal(t, utils.ToDate(now), *contract.SignedAt)
	require.Equal(t, utils.ToDate(now), *contract.EndedAt)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), contract.SourceOfTruth)
}

func TestContractEventHandler_OnRefreshStatus_Ended(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Name:           "test contract",
		ContractStatus: neo4jenum.ContractStatusDraft,
		EndedAt:        utils.Ptr(utils.Now().AddDate(0, 0, -1)),
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
	refreshStatusEvent, err := event.NewContractRefreshStatusEvent(contractAggregate)
	require.Nil(t, err)

	// EXECUTE
	err = contractEventHandler.OnRefreshStatus(context.Background(), refreshStatusEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"Contract": 1, "Contract_" + tenantName: 1,
		"Action": 1, "Action_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := mapper.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionContractStatusUpdated, action.Type)
	require.Equal(t, "test contract has ended", action.Content)
	require.Equal(t, `{"status":"ENDED","contract-name":"test contract","comment":"test contract has ended"}`, action.Metadata)
}

func TestContractEventHandler_OnRefreshStatus_Live(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Name:             "test contract",
		ContractStatus:   neo4jenum.ContractStatusDraft,
		Approved:         true,
		ServiceStartedAt: utils.Ptr(utils.Now().AddDate(0, 0, -1)),
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
	refreshStatusEvent, err := event.NewContractRefreshStatusEvent(contractAggregate)
	require.Nil(t, err)

	// prepare grpc mock for onboarding status update
	calledEventsPlatformForOnboardingStatusChange := false
	organizationServiceCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		UpdateOnboardingStatus: func(context context.Context, org *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, org.AppSource)
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
	err = contractEventHandler.OnRefreshStatus(context.Background(), refreshStatusEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"Contract": 1, "Contract_" + tenantName: 1,
		"Action": 1, "Action_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := mapper.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)

	// verify grpc was called
	require.False(t, calledEventsPlatformForOnboardingStatusChange)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionContractStatusUpdated, action.Type)
	require.Equal(t, "test contract is now live", action.Content)
	require.Equal(t, `{"status":"LIVE","contract-name":"test contract","comment":"test contract is now live"}`, action.Metadata)
}

func TestContractEventHandler_OnUpdate_SubsetOfFieldsSet(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1})

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	now := utils.Now()

	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	updateEvent, err := event.NewContractUpdateEvent(contractAggregate,
		model.ContractDataFields{
			CanPayWithCard:         true,
			CanPayWithBankTransfer: false,
			AutoRenew:              true,
			Check:                  true,
			DueDays:                60,
			InvoiceEmail:           "to@gmail.com",
			InvoiceEmailCC:         []string{"cc1@gmail.com", "cc2@gmail.com"},
			InvoiceEmailBCC:        []string{"bcc1@gmail.com", "bcc2@gmail.com"},
		},
		commonmodel.ExternalSystem{},
		constants.SourceOpenline,
		now,
		[]string{event.FieldMaskAutoRenew,
			event.FieldMaskCanPayWithCard,
			event.FieldMaskCheck,
			event.FieldMaskDueDays,
			event.FieldMaskCanPayWithBankTransfer,
			event.FieldMaskInvoiceEmail,
			event.FieldMaskInvoiceEmailCC,
			event.FieldMaskInvoiceEmailBCC})
	require.Nil(t, err, "failed to create event")

	// EXECUTE
	err = contractEventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err, "failed to execute event handler")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contract": 1, "Contract_" + tenantName: 1})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// verify contract
	contract := mapper.MapDbNodeToContractEntity(contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, true, contract.CanPayWithCard)
	require.Equal(t, false, contract.CanPayWithBankTransfer)
	require.Equal(t, true, contract.AutoRenew)
	require.Equal(t, true, contract.Check)
	require.Equal(t, int64(60), contract.DueDays)
	require.Equal(t, "to@gmail.com", contract.InvoiceEmail)
	require.Equal(t, []string{"cc1@gmail.com", "cc2@gmail.com"}, contract.InvoiceEmailCC)
	require.Equal(t, []string{"bcc1@gmail.com", "bcc2@gmail.com"}, contract.InvoiceEmailBCC)
}

func TestContractEventHandler_OnDeleteV1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId1 := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganization: 1, model2.NodeLabelOrganization + "_" + tenantName: 1,
		model2.NodeLabelContract: 2, model2.NodeLabelContract + "_" + tenantName: 2})

	// prepare grpc mock
	calledEventsPlatformToRefreshRenewalSummary, calledEventsPlatformToRefreshArr := false, false
	organizationServiceRefreshCallbacks := mocked_grpc.MockOrganizationServiceCallbacks{
		RefreshRenewalSummary: func(context context.Context, org *organizationpb.RefreshRenewalSummaryGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, org.AppSource)
			calledEventsPlatformToRefreshRenewalSummary = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: orgId,
			}, nil
		},
		RefreshArr: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, orgId, org.OrganizationId)
			require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, org.AppSource)
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
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId1)
	now := utils.Now()
	deleteEvent, err := event.NewContractDeleteEvent(contractAggregate, now)
	require.Nil(t, err)

	// EXECUTE
	err = contractEventHandler.OnDeleteV1(context.Background(), deleteEvent)
	require.Nil(t, err)

	// VERIFY
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model2.NodeLabelOrganization: 1, model2.NodeLabelOrganization + "_" + tenantName: 1,
		model2.NodeLabelContract: 1, model2.NodeLabelContract + "_" + tenantName: 1,
		model2.NodeLabelDeletedContract: 1, model2.NodeLabelDeletedContract + "_" + tenantName: 1,
	})

	contractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId2)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	deletedContractDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "DeletedContract_"+tenantName, contractId1)
	require.Nil(t, err)
	require.NotNil(t, deletedContractDbNode)

	// verify call to events platform
	require.True(t, calledEventsPlatformToRefreshRenewalSummary)
	require.True(t, calledEventsPlatformToRefreshArr)
}

func TestContractEventHandler_DeriveContractStatus_Ended(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	now := utils.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractEntity := neo4jentity.ContractEntity{
		Name:           "test contract",
		ContractStatus: neo4jenum.ContractStatusDraft,
		EndedAt:        &yesterday,
	}
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, contractEntity)
	contractEntity.Id = contractId
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewedAt:         &tomorrow,
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
			RenewalApproved:   false,
		},
	})

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}
	// EXECUTE
	status, err := contractEventHandler.deriveContractStatus(ctx, tenantName, contractEntity)
	require.Nil(t, err)
	require.Equal(t, neo4jenum.ContractStatusEnded.String(), status)
}

func TestContractEventHandler_DeriveContractStatus_Draft_NoServiceStartedAt(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractEntity := neo4jentity.ContractEntity{
		Name:             "test contract",
		ContractStatus:   neo4jenum.ContractStatusDraft,
		ServiceStartedAt: nil,
	}
	_ = neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, contractEntity)

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// EXECUTE
	status, err := contractEventHandler.deriveContractStatus(ctx, tenantName, contractEntity)
	require.Nil(t, err)
	require.Equal(t, neo4jenum.ContractStatusDraft.String(), status)
}

func TestContractEventHandler_DeriveContractStatus_Draft_FutureServiceStartedAt(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	tomorrow := utils.Now().AddDate(0, 0, 1)
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractEntity := neo4jentity.ContractEntity{
		Name:             "test contract",
		ContractStatus:   neo4jenum.ContractStatusDraft,
		ServiceStartedAt: &tomorrow,
	}
	_ = neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, contractEntity)

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// EXECUTE
	status, err := contractEventHandler.deriveContractStatus(ctx, tenantName, contractEntity)
	require.Nil(t, err)
	require.Equal(t, neo4jenum.ContractStatusDraft.String(), status)
}

func TestContractEventHandler_DeriveContractStatus_Live_AutoRenew_ActiveRenewalOpportunity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	now := utils.Now()
	tomorrow := now.AddDate(0, 0, 1)
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractEntity := neo4jentity.ContractEntity{
		Name:             "test contract",
		ContractStatus:   neo4jenum.ContractStatusLive,
		AutoRenew:        true,
		Approved:         true,
		ServiceStartedAt: &now,
	}
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, contractEntity)
	contractEntity.Id = contractId
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewedAt:         &tomorrow,
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
			RenewalApproved:   false,
		},
	})

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// EXECUTE
	status, err := contractEventHandler.deriveContractStatus(ctx, tenantName, contractEntity)
	require.Nil(t, err)
	require.Equal(t, neo4jenum.ContractStatusLive.String(), status)
}

func TestContractEventHandler_DeriveContractStatus_Live_NoAutoRenew_NoActiveRenewalOpportunity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	now := utils.Now()
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractEntity := neo4jentity.ContractEntity{
		Name:             "test contract",
		ContractStatus:   neo4jenum.ContractStatusLive,
		AutoRenew:        false,
		Approved:         true,
		ServiceStartedAt: &now,
	}
	_ = neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, contractEntity)

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// EXECUTE
	status, err := contractEventHandler.deriveContractStatus(ctx, tenantName, contractEntity)
	require.Nil(t, err)
	require.Equal(t, neo4jenum.ContractStatusLive.String(), status)
}

func TestContractEventHandler_DeriveContractStatus_OutOfContract_NoAutoRenew_ActiveRenewalOpportunity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	now := utils.Now()
	yesterday := now.AddDate(0, 0, -1)
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractEntity := neo4jentity.ContractEntity{
		Name:             "test contract",
		ContractStatus:   neo4jenum.ContractStatusLive,
		AutoRenew:        false,
		Approved:         true,
		ServiceStartedAt: &now,
	}
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, contractEntity)
	contractEntity.Id = contractId
	neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewedAt:       &yesterday,
			RenewalApproved: false,
		},
	})

	// prepare event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// EXECUTE
	status, err := contractEventHandler.deriveContractStatus(ctx, tenantName, contractEntity)
	require.Nil(t, err)
	require.Equal(t, neo4jenum.ContractStatusOutOfContract.String(), status)
}
