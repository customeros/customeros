package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/mocked_grpc"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceLineItemEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Id: contractId,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Quantity:   10,
			Price:      100.50,
			VatRate:    20.5,
			Name:       "Test service line item",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
		"",
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreateV1(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentID)
	require.Equal(t, neo4jenum.BilledTypeMonthly, serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, 100.50, serviceLineItem.Price)
	require.Equal(t, 20.5, serviceLineItem.VatRate)
	require.Equal(t, "Test service line item", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)
}

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
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
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
		commonmodel.Source{
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
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
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
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
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
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
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
	require.Equal(t, now, serviceLineItem.UpdatedAt)
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
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
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
		commonmodel.Source{
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
	require.Equal(t, float64(200.0), serviceLineItem.Price)

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
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
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
		commonmodel.Source{
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
	require.Equal(t, float64(200.0), serviceLineItem.Price)
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
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
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
		commonmodel.Source{
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
	require.Equal(t, float64(50.0), serviceLineItem.Price)
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
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
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
		commonmodel.Source{
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
	require.Equal(t, float64(50.0), serviceLineItem.Price)
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
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Quantity: 20,
		},
		commonmodel.Source{
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
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemUpdateEvent
	updatedAt := utils.Now()
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId),
		model.ServiceLineItemDataFields{
			Quantity: 350,
		},
		commonmodel.Source{
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

func TestServiceLineItemEventHandler_OnCreateRecurringMonthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Quantity:   10,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
		"",
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreateV1(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentID)
	require.Equal(t, neo4jenum.BilledTypeMonthly, serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemBilledTypeRecurringCreated, action.Type)
	require.Equal(t, fmt.Sprintf("logged-in user added a recurring service to Contract 1: Service 1 at 10 x 170.25/month starting with %s", timeNow.Format("2006-01-02")), action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"currency":"","comment":"billed type is MONTHLY for service Service 1","reasonForChange":"","startedAt":"2024-05-29T00:00:00Z","billedType":"MONTHLY","quantity":10,"previousPrice":0,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateRecurringAnnually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.AnnuallyBilled,
			Quantity:   10,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
		"",
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreateV1(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentID)
	require.Equal(t, neo4jenum.BilledTypeAnnually, serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemBilledTypeRecurringCreated, action.Type)
	require.Equal(t, "logged-in user added a recurring service to Contract 1: Service 1 at 10 x 170.25/year starting with "+timeNow.Format("2006-01-02"), action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"currency":"","comment":"billed type is ANNUALLY for service Service 1","reasonForChange":"","startedAt":"2024-05-29T00:00:00Z","billedType":"ANNUALLY","quantity":10,"previousPrice":0,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateRecurringQuarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.QuarterlyBilled,
			Quantity:   10,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
		"",
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreateV1(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentID)
	require.Equal(t, neo4jenum.BilledTypeQuarterly, serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemBilledTypeRecurringCreated, action.Type)
	require.Equal(t, "logged-in user added a recurring service to Contract 1: Service 1 at 10 x 170.25/quarter starting with "+timeNow.Format("2006-01-02"), action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"currency":"","comment":"billed type is QUARTERLY for service Service 1","reasonForChange":"","startedAt":"2024-05-29T00:00:00Z","billedType":"QUARTERLY","quantity":10,"previousPrice":0,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateOnce(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.OnceBilled,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
		"",
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreateV1(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentID)
	require.Equal(t, neo4jenum.BilledTypeOnce, serviceLineItem.Billed)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemBilledTypeOnceCreated, action.Type)
	require.Equal(t, "logged-in user added a one time service to Contract 1: Service 1 at 170.25 starting with "+timeNow.Format("2006-01-02"), action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"currency":"","comment":"billed type is ONCE for service Service 1","reasonForChange":"","startedAt":"2024-05-29T00:00:00Z","billedType":"ONCE","quantity":0,"previousPrice":0,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreatePerUse(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.UsageBilled,
			Price:      170.25,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
		"",
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreateV1(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentID)
	require.Equal(t, neo4jenum.BilledTypeUsage, serviceLineItem.Billed)
	require.Equal(t, float64(170.25), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemBilledTypeUsageCreated, action.Type)
	require.Equal(t, "logged-in user added a per use service to Contract 1: Service 1 at 170.2500 starting with "+timeNow.Format("2006-01-02"), action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":170.25,"currency":"","comment":"billed type is USAGE for service Service 1","reasonForChange":"","startedAt":"2024-05-29T00:00:00Z","billedType":"USAGE","quantity":0,"previousPrice":0,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateNewVersionForNonRetroactiveQuantityDecrease(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Id:   contractId,
		Name: "Contract 1",
	})

	serviceLineItemParentId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:     "Service Parent",
		Quantity: 400,
	})
	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}
	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Quantity:   10,
			Name:       "Test service line item",
			ContractId: contractId,
			ParentId:   serviceLineItemParentId,
			Comments:   "reason for what change?",
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
		serviceLineItemParentId,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreateV1(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 2, "ServiceLineItem_" + tenantName: 2,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, serviceLineItemParentId, serviceLineItem.ParentID)
	require.Equal(t, neo4jenum.BilledTypeMonthly, serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, "Test service line item", serviceLineItem.Name)
	require.Equal(t, "reason for what change?", serviceLineItem.Comments)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemQuantityUpdated, action.Type)
	require.Equal(t, "logged-in user decreased the quantity of Test service line item from 400 to 10 starting with "+timeNow.Format("2006-01-02"), action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Test service line item","price":0,"currency":"","comment":"quantity is 10 for service Test service line item","reasonForChange":"reason for what change?","startedAt":"2024-05-29T00:00:00Z","billedType":"MONTHLY","quantity":10,"previousPrice":0,"previousQuantity":400}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateNewVersionForNonRetroactivePriceIncrease(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Id:       contractId,
		Name:     "Contract 1",
		Currency: neo4jenum.CurrencyUSD,
	})

	serviceLineItemParentId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:  "Service Parent",
		Price: 1500.56,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Price:      850.75,
			Name:       "Test service line item",
			ContractId: contractId,
			ParentId:   serviceLineItemParentId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
		serviceLineItemParentId,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreateV1(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 2, "ServiceLineItem_" + tenantName: 2,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, serviceLineItemParentId, serviceLineItem.ParentID)
	require.Equal(t, neo4jenum.BilledTypeMonthly, serviceLineItem.Billed)
	require.Equal(t, float64(850.75), serviceLineItem.Price)
	require.Equal(t, "Test service line item", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user decreased the price for Test service line item from 1500.56/month to 850.75/month starting with "+timeNow.Format("2006-01-02"), action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Test service line item","price":850.75,"currency":"USD","comment":"price is 850.75 for service Test service line item","reasonForChange":"","startedAt":"2024-05-29T00:00:00Z","billedType":"MONTHLY","quantity":0,"previousPrice":1500.56,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnCreateNewVersionForNonRetroactivePriceIncreaseNoNameService(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Id:       contractId,
		Name:     "Contract 1",
		Currency: neo4jenum.CurrencyUSD,
	})

	serviceLineItemParentId := neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Price: 1500.56,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Price:      850.75,
			ContractId: contractId,
			ParentId:   serviceLineItemParentId,
			Comments:   "This is a reason for change",
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
		serviceLineItemParentId,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreateV1(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 2, "ServiceLineItem_" + tenantName: 2,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.Equal(t, serviceLineItemParentId, serviceLineItem.ParentID)
	require.Equal(t, neo4jenum.BilledTypeMonthly, serviceLineItem.Billed)
	require.Equal(t, float64(850.75), serviceLineItem.Price)
	require.Equal(t, "This is a reason for change", serviceLineItem.Comments)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user decreased the price for Unnamed service from 1500.56/month to 850.75/month starting with "+timeNow.Format("2006-01-02"), action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Unnamed service","price":850.75,"currency":"USD","comment":"price is 850.75 for service Unnamed service","reasonForChange":"This is a reason for change","startedAt":"2024-05-29T00:00:00Z","billedType":"MONTHLY","quantity":0,"previousPrice":1500.56,"previousQuantity":0}`, action.Metadata)
}

func TestServiceLineItemEventHandler_OnUpdatePriceNonRetroactiveForExistingSLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	serviceLineItemId1 := "service-line-item-id-1"
	serviceLineItemId2 := "service-line-item-id-2"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "logged-in",
		LastName:  "user",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenantName, orgId, neo4jentity.ContractEntity{
		Id:       contractId,
		Name:     "Contract 1",
		Currency: neo4jenum.CurrencyUSD,
	})

	neo4jtest.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		ID:       serviceLineItemId1,
		Name:     "service 1",
		Billed:   neo4jenum.BilledTypeMonthly,
		Price:    170.25,
		Quantity: 10,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId2)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:     model.MonthlyBilled,
			Quantity:   10,
			Price:      100,
			Name:       "Service 1",
			ContractId: contractId,
			ParentId:   serviceLineItemId1,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
		timeNow,
		timeNow,
		timeNow,
		nil,
		serviceLineItemId1,
	)
	require.Nil(t, err, "failed to create service line item create event")

	metadata := make(map[string]string)
	metadata["user-id"] = userId
	err = createEvent.SetMetadata(metadata)
	require.Nil(t, err)

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreateV1(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 2, "ServiceLineItem_" + tenantName: 2,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1,
	})
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId2)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId2)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId2, serviceLineItem.ID)
	require.Equal(t, serviceLineItemId1, serviceLineItem.ParentID)
	require.Equal(t, neo4jenum.BilledTypeMonthly, serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Quantity)
	require.Equal(t, float64(100), serviceLineItem.Price)
	require.Equal(t, "Service 1", serviceLineItem.Name)
	require.Equal(t, timeNow, serviceLineItem.CreatedAt)
	require.Equal(t, timeNow, serviceLineItem.UpdatedAt)
	require.Equal(t, utils.ToDate(timeNow), serviceLineItem.StartedAt)
	require.Nil(t, serviceLineItem.EndedAt)

	// verify action
	actionDbNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := neo4jmapper.MapDbNodeToActionEntity(actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, action.AppSource)
	require.Equal(t, neo4jenum.ActionServiceLineItemPriceUpdated, action.Type)
	require.Equal(t, "logged-in user decreased the price for Service 1 from 170.25/month to 100.00/month starting with "+timeNow.Format("2006-01-02"), action.Content)
	require.Equal(t, `{"user-name":"logged-in user","service-name":"Service 1","price":100,"currency":"USD","comment":"price is 100.00 for service Service 1","reasonForChange":"","startedAt":"2024-05-29T00:00:00Z","billedType":"MONTHLY","quantity":10,"previousPrice":170.25,"previousQuantity":0}`, action.Metadata)
}
