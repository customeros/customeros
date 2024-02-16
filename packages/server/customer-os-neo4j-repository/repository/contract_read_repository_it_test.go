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

func TestContractReadRepository_GetContractsToGenerateCycleInvoices(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	referenceDate := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled:  true,
		InvoicingPostpaid: false,
		DefaultCurrency:   neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelTenant:         1,
		neo4jutil.NodeLabelTenantSettings: 1,
		neo4jutil.NodeLabelOrganization:   1,
		neo4jutil.NodeLabelContract:       1,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateCycleInvoices_OrganizationIsHidden(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"
	referenceDate := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled:  true,
		InvoicingPostpaid: true,
		DefaultCurrency:   neo4jenum.CurrencyUSD,
	})

	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	organizationHidden := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{
		Hide: true,
	})

	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleQuarterlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationHidden, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId2, entity.ServiceLineItemEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelTenant:         1,
		neo4jutil.NodeLabelTenantSettings: 1,
		neo4jutil.NodeLabelOrganization:   2,
		neo4jutil.NodeLabelContract:       2,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateCycleInvoices_InvoicingNodeEnabled(t *testing.T) {
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
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleAnnuallyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})

	organizationIdNok := neo4jtest.CreateOrganization(ctx, driver, tenantNok, entity.OrganizationEntity{})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, driver, tenantNok, organizationIdNok, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId2, entity.ServiceLineItemEntity{})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateCycleInvoices_MissingCurrency(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"
	referenceDate := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		Currency:              neo4jenum.CurrencyAUD,
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId2, entity.ServiceLineItemEntity{})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateCycleInvoices_MissingBillingCycle(t *testing.T) {
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
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId2, entity.ServiceLineItemEntity{})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateCycleInvoices_CheckByNextInvoiceDate(t *testing.T) {
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
	contractIdYesterday := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		NextInvoiceDate:       &yesterday,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractIdToday := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		NextInvoiceDate:       &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractIdTomorrow := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractIdYesterday, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractIdToday, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractIdTomorrow, entity.ServiceLineItemEntity{})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 2)
	props1 := utils.GetPropsFromNode(*result[0].Node)
	props2 := utils.GetPropsFromNode(*result[1].Node)
	contractId1 := utils.GetStringPropOrEmpty(props1, "id")
	contractId2 := utils.GetStringPropOrEmpty(props2, "id")
	require.ElementsMatch(t, []string{contractIdToday, contractIdYesterday}, []string{contractId1, contractId2})
	require.NotContains(t, []string{contractIdTomorrow}, []string{contractId1, contractId2})
}

func TestContractReadRepository_GetContractsToGenerateCycleInvoices_CheckByInvoicingStartDate(t *testing.T) {
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
	contractIdYesterday := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &yesterday,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractIdToday := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &today,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractIdTomorrow := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractIdNoInvoicingStartDate := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractIdYesterday, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractIdToday, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractIdTomorrow, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractIdNoInvoicingStartDate, entity.ServiceLineItemEntity{})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 2)
	props1 := utils.GetPropsFromNode(*result[0].Node)
	props2 := utils.GetPropsFromNode(*result[1].Node)
	contractId1 := utils.GetStringPropOrEmpty(props1, "id")
	contractId2 := utils.GetStringPropOrEmpty(props2, "id")
	require.ElementsMatch(t, []string{contractIdToday, contractIdYesterday}, []string{contractId1, contractId2})
	require.NotContains(t, []string{contractIdTomorrow, contractIdNoInvoicingStartDate}, []string{contractId1, contractId2})
}

func TestContractReadRepository_GetContractsToGenerateCycleInvoices_CheckByContractEndDate(t *testing.T) {
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
	contractIdYesterday := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &today,
		EndedAt:               &yesterday,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractIdToday := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &today,
		EndedAt:               &today,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractIdTomorrow := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &today,
		EndedAt:               &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractIdTomorrow, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractIdToday, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractIdYesterday, entity.ServiceLineItemEntity{})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props1 := utils.GetPropsFromNode(*result[0].Node)
	contractId1 := utils.GetStringPropOrEmpty(props1, "id")
	require.ElementsMatch(t, []string{contractIdTomorrow}, []string{contractId1})
	require.NotContains(t, []string{contractIdToday, contractIdYesterday}, []string{contractId1})
}

func TestContractReadRepository_GetContractsToGenerateCycleInvoices_MissingOrganizationLegalName(t *testing.T) {
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
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:       neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate: &referenceDate,
		InvoiceEmail:       "invoiceEmail",
		InvoicingEnabled:   true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId2, entity.ServiceLineItemEntity{})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateCycleInvoices_MissingInvoiceEmail(t *testing.T) {
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
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId2, entity.ServiceLineItemEntity{})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateCycleInvoices_MissingServiceLineItems(t *testing.T) {
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
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &referenceDate,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateCycleInvoices(ctx, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_InvoicingNotEnabled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: false,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_OrganizationIsHidden(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{
		Hide: true,
	})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_MissingCurrency(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_MissingOrganizationLegalName(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:     neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:  &tomorrow,
		InvoiceEmail:     "invoiceEmail",
		InvoicingEnabled: true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_MissingInvoiceEmail(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_NextInvoiceDateNotSet(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		InvoicingStartDate:    &today,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_ContractAlreadyEnded(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		EndedAt:               &today,
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_ServiceLineAlreadyInvoiced(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	sliId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{})
	invoiceLineId := neo4jtest.CreateInvoiceLine(ctx, driver, tenant, invoiceId, entity.InvoiceLineEntity{})
	neo4jtest.LinkNodes(ctx, driver, invoiceLineId, sliId, "INVOICED")

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_ServiceLineStartedAtLeast2DaysBeforeNextOnCycleInvoiceDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	yesterday := today.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &today,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_ServiceLineStartedSameDayAsInvoicingRun(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: today,
	})

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_LastServiceLineItemIsInvoiced(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)
	beforeYesterday := yesterday.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	invoicedSliId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: beforeYesterday,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{})
	invoiceLineId := neo4jtest.CreateInvoiceLine(ctx, driver, tenant, invoiceId, entity.InvoiceLineEntity{})
	neo4jtest.LinkNodes(ctx, driver, invoiceLineId, invoicedSliId, "INVOICED")

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestContractReadRepository_GetContractsToGenerateOffCycleInvoices_LastServiceLineItemIsNotInvoiced(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	today := utils.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)
	beforeYesterday := yesterday.Add(-24 * time.Hour)
	threeDaysAgo := beforeYesterday.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, entity.TenantSettingsEntity{
		InvoicingEnabled: true,
		DefaultCurrency:  neo4jenum.CurrencyUSD,
	})
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{
		BillingCycle:          neo4jenum.BillingCycleMonthlyBilling,
		NextInvoiceDate:       &tomorrow,
		OrganizationLegalName: "organizationLegalName",
		InvoiceEmail:          "invoiceEmail",
		InvoicingEnabled:      true,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: yesterday,
	})
	invoicedSliId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: beforeYesterday,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		StartedAt: threeDaysAgo,
	})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{})
	invoiceLineId := neo4jtest.CreateInvoiceLine(ctx, driver, tenant, invoiceId, entity.InvoiceLineEntity{})
	neo4jtest.LinkNodes(ctx, driver, invoiceLineId, invoicedSliId, "INVOICED")

	result, err := repositories.ContractReadRepository.GetContractsToGenerateOffCycleInvoices(ctx, today)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, contractId, props["id"])
}
