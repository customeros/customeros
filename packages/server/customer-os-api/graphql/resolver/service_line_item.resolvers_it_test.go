package resolver

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	"github.com/customeros/customeros/packages/server/customer-os-api/utils/decode"
)

func TestMutationResolver_ServiceLineItemUpdate_NoChanges(t *testing.T) {
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
		Quantity: 2,
		Price:    30,
		Comments: "test comments",
		VatRate:  10.5,
	})

	rawResponse := callGraphQL(t, "service_line_item/update_service_line_item_parameterized", map[string]interface{}{
		"serviceLineItemId":       baseSliId,
		"description":             "service",
		"price":                   int64(30),
		"quantity":                float64(2),
		"comments":                "test comments",
		"isRetroactiveCorrection": false,
		"taxRate":                 10.5,
	})

	var serviceLineItemStruct struct {
		ContractLineItem_Update model.ServiceLineItem
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &serviceLineItemStruct)
	require.Nil(t, err)
	serviceLineItem := serviceLineItemStruct.ContractLineItem_Update
	require.Equal(t, baseSliId, serviceLineItem.Metadata.ID)
	require.Equal(t, baseSliId, serviceLineItem.ParentID)
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
