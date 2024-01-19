package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestContractReadRepository_GetContractsToGenerateOnCycleInvoices(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	referenceDate := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &referenceDate,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelTenant:         1,
		neo4jutil.NodeLabelTenantSettings: 1,
		neo4jutil.NodeLabelOrganization:   1,
		neo4jutil.NodeLabelContract:       1,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOnCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateOnCycleInvoices_OrganizationIsHidden(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"
	referenceDate := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})

	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	organizationHidden := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{
		Hide: true,
	})

	contractId := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleQuarterlyBilling,
		InvoicingStartDate: &referenceDate,
	})
	neo4jtest.CreateContract(ctx, driver, tenant, organizationHidden, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &referenceDate,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelTenant:         1,
		neo4jutil.NodeLabelTenantSettings: 1,
		neo4jutil.NodeLabelOrganization:   2,
		neo4jutil.NodeLabelContract:       2,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOnCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateOnCycleInvoices_InvoicingNodeEnabled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"
	tenantNok := "tenant2"
	referenceDate := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})

	neo4jtest.CreateTenant(ctx, driver, tenantNok)
	neo4jtest.CreateTenantSettings(ctx, driver, tenantNok, entity.TenantSettingsEntity{
		InvoicingEnabled: false,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})

	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleAnnualBilling,
		InvoicingStartDate: &referenceDate,
	})

	organizationIdNok := neo4jtest.CreateOrganization(ctx, driver, tenantNok, entity.OrganizationEntity{})
	neo4jtest.CreateContract(ctx, driver, tenantNok, organizationIdNok, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &referenceDate,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOnCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateOnCycleInvoices_MissingCurrency(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"
	referenceDate := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		Currency:           neo4jenum.CurrencyAUD,
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &referenceDate,
	})
	neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &referenceDate,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOnCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateOnCycleInvoices_MissingBillingCycle(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant"
	referenceDate := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &referenceDate,
	})
	neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		InvoicingStartDate: &referenceDate,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOnCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateOnCycleInvoices_CheckByNextInvoiceDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant"
	referenceDate := utils.Now()
	tomorrow := referenceDate.Add(24 * time.Hour)
	yesterday := referenceDate.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractIdYesterday := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &referenceDate,
		NextInvoiceDate:    &yesterday,
	})
	contractIdToday := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &referenceDate,
		NextInvoiceDate:    &referenceDate,
	})
	contractIdTomorrow := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &referenceDate,
		NextInvoiceDate:    &tomorrow,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOnCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 2)
	props1 := utils.GetPropsFromNode(*result[0].Node)
	props2 := utils.GetPropsFromNode(*result[1].Node)
	contractId1 := utils.GetStringPropOrEmpty(props1, "id")
	contractId2 := utils.GetStringPropOrEmpty(props2, "id")
	require.ElementsMatch(t, []string{contractIdToday, contractIdYesterday}, []string{contractId1, contractId2})
	require.NotContains(t, []string{contractIdTomorrow}, []string{contractId1, contractId2})
}

func TestContractReadRepository_GetContractsToGenerateOnCycleInvoices_CheckByInvoicingStartDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant"
	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractIdYesterday := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &yesterday,
	})
	contractIdToday := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &today,
	})
	contractIdTomorrow := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &tomorrow,
	})
	contractIdNoInvoicingStartDate := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle: neo4jenum.BillingCycleMonthlyBilling,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOnCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 2)
	props1 := utils.GetPropsFromNode(*result[0].Node)
	props2 := utils.GetPropsFromNode(*result[1].Node)
	contractId1 := utils.GetStringPropOrEmpty(props1, "id")
	contractId2 := utils.GetStringPropOrEmpty(props2, "id")
	require.ElementsMatch(t, []string{contractIdToday, contractIdYesterday}, []string{contractId1, contractId2})
	require.NotContains(t, []string{contractIdTomorrow, contractIdNoInvoicingStartDate}, []string{contractId1, contractId2})
}

func TestContractReadRepository_GetContractsToGenerateOnCycleInvoices_CheckByContractEndDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant"
	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractIdYesterday := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &today,
		EndedAt:            &yesterday,
	})
	contractIdToday := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &today,
		EndedAt:            &today,
	})
	contractIdTomorrow := neo4jtest.CreateContract(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &today,
		EndedAt:            &tomorrow,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOnCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props1 := utils.GetPropsFromNode(*result[0].Node)
	contractId1 := utils.GetStringPropOrEmpty(props1, "id")
	require.ElementsMatch(t, []string{contractIdTomorrow}, []string{contractId1})
	require.NotContains(t, []string{contractIdToday, contractIdYesterday}, []string{contractId1})
}
