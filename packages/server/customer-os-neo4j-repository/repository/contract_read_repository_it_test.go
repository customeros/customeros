package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestContractReadRepository_GetContractsForInvoicing_Hidden_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{
		Hide:       true,
		IsCustomer: true,
	})

	day1 := neo4jtest.FirstTimeOfMonth(2023, 6)
	neo4jtest.CreateContract(ctx, driver, tenantName, organizationId, entity.ContractEntity{
		NextInvoiceDate: utils.TimePtr(day1),
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelTenant:       1,
		neo4jutil.NodeLabelOrganization: 1,
		neo4jutil.NodeLabelContract:     1,
	})

	day1contractsForInvoicing, err := repositories.ContractReadRepository.GetContractsForInvoicing(ctx, day1)
	require.NoError(t, err)
	require.Len(t, day1contractsForInvoicing, 0)
}

func TestContractReadRepository_GetContractsForInvoicing_Prospect_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{})

	day1 := neo4jtest.FirstTimeOfMonth(2023, 6)
	neo4jtest.CreateContract(ctx, driver, tenantName, organizationId, entity.ContractEntity{
		NextInvoiceDate: utils.TimePtr(day1),
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelTenant:       1,
		neo4jutil.NodeLabelOrganization: 1,
		neo4jutil.NodeLabelContract:     1,
	})

	day1contractsForInvoicing, err := repositories.ContractReadRepository.GetContractsForInvoicing(ctx, day1)
	require.NoError(t, err)
	require.Len(t, day1contractsForInvoicing, 0)
}

func TestContractReadRepository_GetContractsForInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant1 := "tenant1"
	tenant2 := "tenant2"

	day1 := neo4jtest.FirstTimeOfMonth(2023, 6)
	day2 := day1.Add(24 * time.Hour)
	day10 := day1.Add(10 * 24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant1)
	neo4jtest.CreateTenant(ctx, driver, tenant2)

	organization1Id := neo4jtest.CreateOrganization(ctx, driver, tenant1, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1Id := neo4jtest.CreateContract(ctx, driver, tenant1, organization1Id, entity.ContractEntity{
		NextInvoiceDate: utils.TimePtr(day1),
	})
	contract2Id := neo4jtest.CreateContract(ctx, driver, tenant1, organization1Id, entity.ContractEntity{
		NextInvoiceDate: utils.TimePtr(day2),
	})

	organization2Id := neo4jtest.CreateOrganization(ctx, driver, tenant2, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract3Id := neo4jtest.CreateContract(ctx, driver, tenant2, organization2Id, entity.ContractEntity{
		NextInvoiceDate: utils.TimePtr(day10),
	})
	contract4Id := neo4jtest.CreateContract(ctx, driver, tenant2, organization2Id, entity.ContractEntity{
		NextInvoiceDate: utils.TimePtr(day10),
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelTenant:       2,
		neo4jutil.NodeLabelOrganization: 2,
		neo4jutil.NodeLabelContract:     4,
	})

	day1contractsForInvoicing, err := repositories.ContractReadRepository.GetContractsForInvoicing(ctx, day1)
	require.NoError(t, err)
	require.Len(t, day1contractsForInvoicing, 1)
	for _, dbNode := range day1contractsForInvoicing {
		tn, ctr := getRowData(dbNode)

		require.Equal(t, tenant1, tn)
		require.Equal(t, contract1Id, ctr)
	}

	day2contractsForInvoicing, err := repositories.ContractReadRepository.GetContractsForInvoicing(ctx, day2)
	require.NoError(t, err)
	require.Len(t, day2contractsForInvoicing, 2)
	for _, dbNode := range day1contractsForInvoicing {
		tn, ctr := getRowData(dbNode)

		require.Equal(t, tenant1, tn)
		require.Contains(t, []string{contract1Id, contract2Id}, ctr)
	}

	neo4jtest.MarkInvoicingStarted(ctx, driver, tenant1, contract1Id, day1.Add(12*time.Hour))

	day2contractsForInvoicing, err = repositories.ContractReadRepository.GetContractsForInvoicing(ctx, day2)
	require.NoError(t, err)
	require.Len(t, day2contractsForInvoicing, 1)
	for _, dbNode := range day2contractsForInvoicing {
		tn, ctr := getRowData(dbNode)

		require.Equal(t, tenant1, tn)
		require.Equal(t, contract2Id, ctr)
	}

	neo4jtest.MarkInvoicingStarted(ctx, driver, tenant1, contract2Id, day2.Add(12*time.Hour))

	day10contractsForInvoicing, err := repositories.ContractReadRepository.GetContractsForInvoicing(ctx, day10)
	require.NoError(t, err)
	require.Len(t, day10contractsForInvoicing, 2)
	for _, dbNode := range day10contractsForInvoicing {
		tn, ctr := getRowData(dbNode)

		require.Equal(t, tenant2, tn)
		require.Contains(t, []string{contract3Id, contract4Id}, ctr)
	}
}

func getRowData(props map[string]any) (string, string) {
	tenant := utils.GetStringPropOrEmpty(props, "tenant")
	organizationId := utils.GetStringPropOrEmpty(props, "contractId")

	return tenant, organizationId
}
