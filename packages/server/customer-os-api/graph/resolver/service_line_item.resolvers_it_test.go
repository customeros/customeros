package resolver

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/service_line_item"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_ServiceLineItemCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})
	serviceLineItemId := uuid.New().String()
	calledCreateServiceLineItem := false

	serviceLineItemServiceCallbacks := events_platform.MockServiceLineItemServiceCallbacks{
		CreateServiceLineItem: func(context context.Context, serviceLineItem *servicelineitempb.CreateServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
			require.Equal(t, tenantName, serviceLineItem.Tenant)
			require.Equal(t, contractId, serviceLineItem.ContractId)
			require.Equal(t, testUserId, serviceLineItem.LoggedInUserId)
			require.Equal(t, string(entity.DataSourceOpenline), serviceLineItem.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, serviceLineItem.SourceFields.AppSource)
			require.Equal(t, "Service Line Item 1", serviceLineItem.Name)
			require.Equal(t, servicelineitempb.BilledType_MONTHLY_BILLED, serviceLineItem.Billed)
			require.Equal(t, int64(2), serviceLineItem.Quantity)
			require.Equal(t, float32(30), serviceLineItem.Price)

			calledCreateServiceLineItem = true
			neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
				ID: serviceLineItemId,
			})
			return &servicelineitempb.ServiceLineItemIdGrpcResponse{
				Id: serviceLineItemId,
			}, nil
		},
	}
	events_platform.SetServiceLineItemCallbacks(&serviceLineItemServiceCallbacks)

	rawResponse := callGraphQL(t, "create_service_line_item", map[string]interface{}{
		"contractId": contractId,
	})

	var serviceLineItemStruct struct {
		ServiceLineItemCreate model.ServiceLineItem
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &serviceLineItemStruct)
	require.Nil(t, err)
	serviceLineItem := serviceLineItemStruct.ServiceLineItemCreate
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)
	require.True(t, calledCreateServiceLineItem)
}

func TestMutationResolver_ServiceLineItemUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})
	serviceLineItemId := neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{Name: "service"})
	calledUpdateServiceLineItem := false
	serviceLineItemServiceCallbacks := events_platform.MockServiceLineItemServiceCallbacks{
		UpdateServiceLineItem: func(context context.Context, serviceLineItem *servicelineitempb.UpdateServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
			require.Equal(t, tenantName, serviceLineItem.Tenant)
			require.Equal(t, testUserId, serviceLineItem.LoggedInUserId)
			require.Equal(t, string(entity.DataSourceOpenline), serviceLineItem.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, serviceLineItem.SourceFields.AppSource)
			require.Equal(t, "Service Line Item 1", serviceLineItem.Name)
			require.Equal(t, servicelineitempb.BilledType_MONTHLY_BILLED, serviceLineItem.Billed)
			require.Equal(t, int64(2), serviceLineItem.Quantity)
			require.Equal(t, float32(30), serviceLineItem.Price)
			require.Equal(t, "test comments", serviceLineItem.Comments)
			calledUpdateServiceLineItem = true
			return &servicelineitempb.ServiceLineItemIdGrpcResponse{
				Id: serviceLineItemId,
			}, nil
		},
	}
	events_platform.SetServiceLineItemCallbacks(&serviceLineItemServiceCallbacks)

	rawResponse := callGraphQL(t, "update_service_line_item", map[string]interface{}{
		"serviceLineItemId": serviceLineItemId,
	})

	var serviceLineItemStruct struct {
		ServiceLineItemUpdate model.ServiceLineItem
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &serviceLineItemStruct)
	require.Nil(t, err)
	serviceLineItem := serviceLineItemStruct.ServiceLineItemUpdate
	require.Equal(t, serviceLineItemId, serviceLineItem.ID)

	require.True(t, calledUpdateServiceLineItem)
}
