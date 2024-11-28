package graph

import (
	"context"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/mocked_grpc"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
	"testing"
)

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
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
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
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
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
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
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
		log:      testLogger,
		services: testDatabase.Services,
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
		log:      testLogger,
		services: testDatabase.Services,
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
		log:      testLogger,
		services: testDatabase.Services,
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
		log:      testLogger,
		services: testDatabase.Services,
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
		log:      testLogger,
		services: testDatabase.Services,
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
		log:      testLogger,
		services: testDatabase.Services,
	}

	// EXECUTE
	status, err := contractEventHandler.deriveContractStatus(ctx, tenantName, contractEntity)
	require.Nil(t, err)
	require.Equal(t, neo4jenum.ContractStatusOutOfContract.String(), status)
}
