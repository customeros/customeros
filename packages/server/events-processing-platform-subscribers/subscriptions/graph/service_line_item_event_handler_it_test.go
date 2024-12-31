package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/mocked_grpc"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceLineItemEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Billed: neo4jenum.BilledTypeMonthly,
	})
	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Name:     "Updated Service Line Item",
			Price:    200.0,
			Quantity: 20,
			VatRate:  20.5,
			Billed:   model.AnnuallyBilled,
		},
		common.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		updatedAt,
		nil,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdateV1(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
	})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, neo4jenum.BilledTypeAnnually, serviceLineItem.Billed)
	require.Equal(t, int64(20), serviceLineItem.Quantity)
	require.Equal(t, 200.0, serviceLineItem.Price)
	require.Equal(t, 20.5, serviceLineItem.VatRate)
	require.Equal(t, "Updated Service Line Item", serviceLineItem.Name)
}

func TestServiceLineItemEventHandler_OnDeleteUnnamed(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})

	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		LengthInMonths: 12,
	})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Billed: neo4jenum.BilledTypeMonthly,
	})
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
	})

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
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

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	// Create a ServiceLineItemDeleteEvent
	deleteEvent, err := event.NewServiceLineItemDeleteEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
	)
	require.Nil(t, err)

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = deleteEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnDeleteV1(ctx, deleteEvent)
	require.Nil(t, err)

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 0, "ServiceLineItem_" + tenantName: 0,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Verify events platform was called
	require.True(t, calledEventsPlatformToUpdateOpportunity)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemRemoved, action.Type)
	require.Equal(t, "logged-in user removed Unnamed service from Unnamed contract", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Unnamed service","price":0,"currency":"","comment":"service line item removed is Unnamed service from Unnamed contract by logged-in user","reasonForChange":"","billedType":"","quantity":0,"previousPrice":0,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnDelete(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Name:           "Contract 1",
		LengthInMonths: 12,
	})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:   "Service 1",
		Billed: neo4jenum.BilledTypeMonthly,
	})
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
	})

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
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

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	// Create a ServiceLineItemDeleteEvent
	deleteEvent, err := event.NewServiceLineItemDeleteEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
	)
	require.Nil(t, err)

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = deleteEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnDeleteV1(ctx, deleteEvent)
	require.Nil(t, err)

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 0, "ServiceLineItem_" + tenantName: 0,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Verify events platform was called
	require.True(t, calledEventsPlatformToUpdateOpportunity)

	// Verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemRemoved, action.Type)
	require.Equal(t, "logged-in user removed Service 1 from Contract 1", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":0,"currency":"","comment":"service line item removed is Service 1 from Contract 1 by logged-in user","reasonForChange":"","billedType":"","quantity":0,"previousPrice":0,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnClose(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		LengthInMonths: 12,
	})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Billed: neo4jenum.BilledTypeMonthly,
	})
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
	})

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
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

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	now := utils.Now()
	// Create a ServiceLineItemCloseEvent
	closeEvent, err := event.NewServiceLineItemCloseEvent(aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId), now, now, true)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnClose(ctx, closeEvent)
	require.Nil(t, err)

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
	})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	test.AssertRecentTime(t, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(now), *serviceLineItem.EndedAt)
	require.True(t, serviceLineItem.Canceled)

	// verify events platform was called
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestServiceLineItemEventHandler_OnUpdatePriceIncreaseRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:   "Service 1",
		Billed: neo4jenum.BilledTypeMonthly,
		Price:  150.0,
	})
	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Billed:   model.MonthlyBilled,
			Price:    200.0,
			Comments: "this is the reason for change",
		},
		common.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		updatedAt,
		nil,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdateV1(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, neo4jenum.BilledTypeMonthly, serviceLineItem.Billed)
	require.Equal(t, 200.0, serviceLineItem.Price)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively increased the price for Service 1 from 150.00/month to 200.00/month", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":200,"currency":"","comment":"price changed is 150.00 for service Service 1","reasonForChange":"this is the reason for change","billedType":"MONTHLY","quantity":0,"previousPrice":150,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdatePriceIncreasePerUseRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:   "Service 1",
		Billed: neo4jenum.BilledTypeUsage,
		Price:  150.0,
	})
	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Price:    200.0,
			Billed:   model.UsageBilled,
			Comments: "test reason for change",
		},
		common.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		updatedAt,
		nil,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdateV1(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, neo4jenum.BilledTypeUsage, serviceLineItem.Billed)
	require.Equal(t, 200.0, serviceLineItem.Price)
	require.Equal(t, "test reason for change", serviceLineItem.Comments)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively increased the price for Service 1 from 150.0000 to 200.0000", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":200,"currency":"","comment":"price changed is 150.00 for service Service 1","reasonForChange":"test reason for change","billedType":"USAGE","quantity":0,"previousPrice":150,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdatePriceDecreaseRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Currency: neo4jenum.CurrencyUSD,
	})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:   "Service 1",
		Billed: neo4jenum.BilledTypeAnnually,
		Price:  150.0,
	})
	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Name:     "Service 1",
			Price:    50.0,
			Billed:   model.AnnuallyBilled,
			Comments: "Reason for change is x",
		},
		common.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		updatedAt,
		nil,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdateV1(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, neo4jenum.BilledTypeAnnually, serviceLineItem.Billed)
	require.Equal(t, 50.0, serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively decreased the price for Service 1 from 150.00/year to 50.00/year", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":50,"currency":"USD","comment":"price changed is 150.00 for service Service 1","reasonForChange":"Reason for change is x","billedType":"ANNUALLY","quantity":0,"previousPrice":150,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdatePriceDecreaseOnceRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Currency: neo4jenum.CurrencyEUR,
	})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:   "Service 1",
		Billed: neo4jenum.BilledTypeOnce,
		Price:  150.0,
	})
	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Name:   "Service 1",
			Price:  50.0,
			Billed: model.OnceBilled,
		},
		common.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		updatedAt,
		nil,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdateV1(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, neo4jenum.BilledTypeOnce, serviceLineItem.Billed)
	require.Equal(t, 50.0, serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively decreased the price for Service 1 from 150.00 to 50.00", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":50,"currency":"EUR","comment":"price changed is 150.00 for service Service 1","reasonForChange":"","billedType":"ONCE","quantity":0,"previousPrice":150,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdateQuantityIncreaseRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:     "Service 1",
		Quantity: 15,
	})
	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Quantity: 20,
		},
		common.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		updatedAt,
		nil,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdateV1(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, int64(20), serviceLineItem.Quantity)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemQuantityUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively increased the quantity of Service 1 from 15 to 20", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":0,"currency":"","comment":"quantity changed is 15 for service Service 1","reasonForChange":"","billedType":"","quantity":20,"previousPrice":0,"previousQuantity":15}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdateQuantityDecreaseRetroactively_TimelineEvent(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare test data in Neo4j
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:     "Service 1",
		Quantity: 400,
	})
	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"Action": 0, "TimelineEvent": 0,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:         testLogger,
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Quantity: 350,
		},
		common.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		updatedAt,
		nil,
	)
	require.Nil(t, err, "failed to create service line item update event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = updateEvent.SetMetadata(metadata)
	require.Nil(t, err)
	// Execute the event handler
	err = serviceLineItemEventHandler.OnUpdateV1(ctx, updateEvent)
	require.Nil(t, err, "failed to execute service line item update event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, neo4jenum.BilledTypeNone, serviceLineItem.Billed)
	require.Equal(t, int64(350), serviceLineItem.Quantity)

	// verify actionat
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemQuantityUpdated, action.Type)
	require.Equal(t, "logged-in user retroactively decreased the quantity of Service 1 from 400 to 350", action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":0,"currency":"","comment":"quantity changed is 400 for service Service 1","reasonForChange":"","billedType":"","quantity":350,"previousPrice":0,"previousQuantity":400}`, action.Metadata)
}
