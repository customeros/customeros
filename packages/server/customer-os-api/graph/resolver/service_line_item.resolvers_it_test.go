package resolver

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_ServiceLineItemCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	serviceLineItemId := uuid.New().String()
	calledCreateServiceLineItem := false

	serviceLineItemServiceCallbacks := events_platform.MockServiceLineItemServiceCallbacks{
		CreateServiceLineItem: func(context context.Context, serviceLineItem *servicelineitempb.CreateServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
			require.Equal(t, tenantName, serviceLineItem.Tenant)
			require.Equal(t, contractId, serviceLineItem.ContractId)
			require.Equal(t, testUserId, serviceLineItem.LoggedInUserId)
			require.Equal(t, string(neo4jentity.DataSourceOpenline), serviceLineItem.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, serviceLineItem.SourceFields.AppSource)
			require.Equal(t, "Service Line Item 1", serviceLineItem.Name)
			require.Equal(t, commonpb.BilledType_MONTHLY_BILLED, serviceLineItem.Billed)
			require.Equal(t, int64(2), serviceLineItem.Quantity)
			require.Equal(t, 30.0, serviceLineItem.Price)
			require.Equal(t, 20.0, serviceLineItem.VatRate)

			calledCreateServiceLineItem = true
			neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
				ID:       serviceLineItemId,
				ParentID: serviceLineItemId,
			})
			return &servicelineitempb.ServiceLineItemIdGrpcResponse{
				Id: serviceLineItemId,
			}, nil
		},
	}
	events_platform.SetServiceLineItemCallbacks(&serviceLineItemServiceCallbacks)

	rawResponse := callGraphQL(t, "service_line_item/create_service_line_item", map[string]interface{}{
		"contractId": contractId,
	})

	var serviceLineItemStruct struct {
		ContractLineItem_Create model.ServiceLineItem
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &serviceLineItemStruct)
	require.Nil(t, err)
	serviceLineItem := serviceLineItemStruct.ContractLineItem_Create
	require.Equal(t, serviceLineItemId, serviceLineItem.Metadata.ID)
	require.Equal(t, serviceLineItemId, serviceLineItem.ParentID)
	require.True(t, calledCreateServiceLineItem)
}

func TestMutationResolver_ServiceLineItemUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	serviceLineItemIdParentId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{Name: "service"})
	//Using serviceLineItemIdParentId as the parent id
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{Name: "service", ParentID: serviceLineItemIdParentId})
	calledUpdateServiceLineItem := false
	serviceLineItemServiceCallbacks := events_platform.MockServiceLineItemServiceCallbacks{
		UpdateServiceLineItem: func(context context.Context, serviceLineItem *servicelineitempb.UpdateServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
			require.Equal(t, tenantName, serviceLineItem.Tenant)
			require.Equal(t, testUserId, serviceLineItem.LoggedInUserId)
			require.Equal(t, string(neo4jentity.DataSourceOpenline), serviceLineItem.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, serviceLineItem.SourceFields.AppSource)
			require.Equal(t, "Service Line Item 1", serviceLineItem.Name)
			require.Equal(t, commonpb.BilledType_MONTHLY_BILLED, serviceLineItem.Billed)
			require.Equal(t, int64(2), serviceLineItem.Quantity)
			require.Equal(t, float64(30), serviceLineItem.Price)
			require.Equal(t, "test comments", serviceLineItem.Comments)
			require.Equal(t, true, serviceLineItem.IsRetroactiveCorrection)
			require.Equal(t, 10.5, serviceLineItem.VatRate)
			calledUpdateServiceLineItem = true
			return &servicelineitempb.ServiceLineItemIdGrpcResponse{
				Id: serviceLineItemId,
			}, nil
		},
	}
	events_platform.SetServiceLineItemCallbacks(&serviceLineItemServiceCallbacks)

	rawResponse := callGraphQL(t, "service_line_item/update_service_line_item", map[string]interface{}{
		"serviceLineItemId": serviceLineItemId,
	})

	var serviceLineItemStruct struct {
		ContractLineItem_Update model.ServiceLineItem
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &serviceLineItemStruct)
	require.Nil(t, err)
	serviceLineItem := serviceLineItemStruct.ContractLineItem_Update
	require.Equal(t, serviceLineItemId, serviceLineItem.Metadata.ID)
	require.Equal(t, serviceLineItemIdParentId, serviceLineItem.ParentID)
	require.True(t, calledUpdateServiceLineItem)
}

func TestMutationResolver_ServiceLineItemDelete(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{})

	calledDeleteServiceLineItem := false
	serviceLineItemServiceCallbacks := events_platform.MockServiceLineItemServiceCallbacks{
		DeleteServiceLineItem: func(context context.Context, serviceLineItem *servicelineitempb.DeleteServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
			require.Equal(t, tenantName, serviceLineItem.Tenant)
			require.Equal(t, serviceLineItemId, serviceLineItem.Id)
			require.Equal(t, testUserId, serviceLineItem.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, constants.AppSourceCustomerOsApi)
			calledDeleteServiceLineItem = true
			return &servicelineitempb.ServiceLineItemIdGrpcResponse{
				Id: serviceLineItemId,
			}, nil
		},
	}
	events_platform.SetServiceLineItemCallbacks(&serviceLineItemServiceCallbacks)

	rawResponse := callGraphQL(t, "service_line_item/delete_service_line_item", map[string]interface{}{
		"serviceLineItemId": serviceLineItemId,
	})

	var response struct {
		ServiceLineItem_Delete model.DeleteResponse
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.True(t, response.ServiceLineItem_Delete.Accepted)
	require.False(t, response.ServiceLineItem_Delete.Completed)
	require.True(t, calledDeleteServiceLineItem)
}

func TestMutationResolver_ServiceLineItemClose(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{})

	calledCloseServiceLineItem := false
	serviceLineItemServiceCallbacks := events_platform.MockServiceLineItemServiceCallbacks{
		CloseServiceLineItem: func(context context.Context, serviceLineItem *servicelineitempb.CloseServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
			require.Equal(t, tenantName, serviceLineItem.Tenant)
			require.Equal(t, serviceLineItemId, serviceLineItem.Id)
			require.Equal(t, testUserId, serviceLineItem.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, constants.AppSourceCustomerOsApi)
			require.Nil(t, serviceLineItem.EndedAt)
			require.Nil(t, serviceLineItem.UpdatedAt)
			calledCloseServiceLineItem = true
			return &servicelineitempb.ServiceLineItemIdGrpcResponse{
				Id: serviceLineItemId,
			}, nil
		},
	}
	events_platform.SetServiceLineItemCallbacks(&serviceLineItemServiceCallbacks)

	rawResponse := callGraphQL(t, "service_line_item/close_service_line_item", map[string]interface{}{
		"serviceLineItemId": serviceLineItemId,
	})

	var response map[string]interface{}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.True(t, calledCloseServiceLineItem)
	require.Equal(t, serviceLineItemId, response["id"])
}
