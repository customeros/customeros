package resolver

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMutationResolver_ServiceLineItemCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		InvoicingEnabled:     true,
		BillingCycleInMonths: 1,
		InvoicingStartDate:   &now,
	})
	serviceLineItemId := uuid.New().String()

	// mock grpc
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

	//wait for generate preview invoice grpc call
	time.Sleep(4 * time.Second)

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
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	baseSliId := uuid.New().String()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		InvoicingEnabled:     true,
		BillingCycleInMonths: 1,
		InvoicingStartDate:   &now,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:     "service",
		ID:       baseSliId,
		ParentID: baseSliId,
		Billed:   neo4jenum.BilledTypeMonthly,
	})

	//Using serviceLineItemIdParentId as the parent id
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{Name: "service", ParentID: baseSliId})

	//mock grpc
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

	//wait for generate preview invoice grpc call
	time.Sleep(4 * time.Second)

	var serviceLineItemStruct struct {
		ContractLineItem_Update model.ServiceLineItem
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &serviceLineItemStruct)
	require.Nil(t, err)
	serviceLineItem := serviceLineItemStruct.ContractLineItem_Update
	require.Equal(t, serviceLineItemId, serviceLineItem.Metadata.ID)
	require.Equal(t, baseSliId, serviceLineItem.ParentID)
	require.True(t, calledUpdateServiceLineItem)
}

func TestMutationResolver_ServiceLineItemNewVersion(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		InvoicingEnabled:     true,
		BillingCycleInMonths: 1,
		InvoicingStartDate:   &now,
	})
	baseSliId := uuid.New().String()
	newSliId := uuid.New().String()
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:      "base service",
		Billed:    neo4jenum.BilledTypeMonthly,
		ID:        baseSliId,
		ParentID:  baseSliId,
		StartedAt: utils.Today(),
	})

	//mock grpc
	calledCreateServiceLineItem := false
	serviceLineItemServiceCallbacks := events_platform.MockServiceLineItemServiceCallbacks{
		CreateServiceLineItem: func(context context.Context, serviceLineItem *servicelineitempb.CreateServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
			require.Equal(t, tenantName, serviceLineItem.Tenant)
			require.Equal(t, testUserId, serviceLineItem.LoggedInUserId)
			require.Equal(t, string(neo4jentity.DataSourceOpenline), serviceLineItem.SourceFields.Source)
			require.Equal(t, constants.AppSourceCustomerOsApi, serviceLineItem.SourceFields.AppSource)
			require.Equal(t, "new version", serviceLineItem.Name)
			require.Equal(t, commonpb.BilledType_MONTHLY_BILLED, serviceLineItem.Billed)
			require.Equal(t, int64(2), serviceLineItem.Quantity)
			require.Equal(t, float64(1.1), serviceLineItem.Price)
			require.Equal(t, "some comments", serviceLineItem.Comments)
			require.Equal(t, float64(33), serviceLineItem.VatRate)
			require.Equal(t, baseSliId, serviceLineItem.ParentId)
			calledCreateServiceLineItem = true
			neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
				ID:       newSliId,
				ParentID: baseSliId,
			})
			return &servicelineitempb.ServiceLineItemIdGrpcResponse{
				Id: newSliId,
			}, nil
		},
	}
	events_platform.SetServiceLineItemCallbacks(&serviceLineItemServiceCallbacks)

	tomorrow := utils.Now().AddDate(0, 0, 1)

	rawResponse := callGraphQL(t, "service_line_item/new_version_service_line_item", map[string]interface{}{
		"serviceLineItemId": baseSliId,
		"serviceStarted":    tomorrow,
	})

	//wait for generate preview invoice grpc call
	time.Sleep(4 * time.Second)

	var serviceLineItemStruct struct {
		ContractLineItem_NewVersion model.ServiceLineItem
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &serviceLineItemStruct)
	require.Nil(t, err)
	serviceLineItem := serviceLineItemStruct.ContractLineItem_NewVersion
	require.Equal(t, newSliId, serviceLineItem.Metadata.ID)
	require.True(t, calledCreateServiceLineItem)
}

func TestMutationResolver_ServiceLineItemNewVersion_VersionAlreadyExists_NotAllowed(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		InvoicingEnabled:     true,
		BillingCycleInMonths: 1,
		InvoicingStartDate:   &now,
	})
	baseSliId := uuid.New().String()

	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Billed:    neo4jenum.BilledTypeMonthly,
		ID:        baseSliId,
		ParentID:  baseSliId,
		StartedAt: utils.Today(),
	})

	response := callGraphQLExpectError(t, "service_line_item/new_version_service_line_item", map[string]interface{}{
		"serviceLineItemId": baseSliId,
		"serviceStarted":    utils.Today(),
	})

	require.Equal(t, fmt.Sprintf("failed to create new contract line item version"), response.Message)
	require.Equal(t, "contractLineItem_NewVersion", response.Path[0])
}

func TestMutationResolver_ServiceLineItemNewVersion_OneTime_NotAllowed(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		InvoicingEnabled:     true,
		BillingCycleInMonths: 1,
		InvoicingStartDate:   &now,
	})
	baseSliId := uuid.New().String()

	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Billed:    neo4jenum.BilledTypeOnce,
		ID:        baseSliId,
		ParentID:  baseSliId,
		StartedAt: utils.Today(),
	})

	response := callGraphQLExpectError(t, "service_line_item/new_version_service_line_item", map[string]interface{}{
		"serviceLineItemId": baseSliId,
		"serviceStarted":    utils.Today().AddDate(0, 0, 1),
	})

	require.Equal(t, fmt.Sprintf("failed to create new contract line item version"), response.Message)
	require.Equal(t, "contractLineItem_NewVersion", response.Path[0])
}

func TestMutationResolver_ServiceLineItemNewVersion_ContractInvoiced_PastVersion_NotAllowed(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		InvoicingEnabled:     true,
		BillingCycleInMonths: 1,
		InvoicingStartDate:   &now,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		DryRun: false,
	})
	baseSliId := uuid.New().String()

	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Billed:    neo4jenum.BilledTypeMonthly,
		ID:        baseSliId,
		ParentID:  baseSliId,
		StartedAt: utils.Today(),
	})

	response := callGraphQLExpectError(t, "service_line_item/new_version_service_line_item", map[string]interface{}{
		"serviceLineItemId": baseSliId,
		"serviceStarted":    utils.Today().AddDate(0, 0, -1),
	})

	require.Equal(t, fmt.Sprintf("failed to create new contract line item version"), response.Message)
	require.Equal(t, "contractLineItem_NewVersion", response.Path[0])
}

func TestMutationResolver_ServiceLineItemDelete(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		InvoicingEnabled:     true,
		BillingCycleInMonths: 1,
		InvoicingStartDate:   &now,
		ContractStatus:       neo4jenum.ContractStatusDraft,
	})
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

	//wait for generate preview invoice grpc call
	time.Sleep(4 * time.Second)

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

	now := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		InvoicingEnabled:     true,
		BillingCycleInMonths: 1,
		InvoicingStartDate:   &now,
	})
	serviceLineItemId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{})

	//mock grpc
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

	//wait for generate preview invoice grpc call
	time.Sleep(4 * time.Second)

	var response map[string]interface{}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.Equal(t, serviceLineItemId, response["id"])
	require.True(t, calledCloseServiceLineItem)
}
