package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInvoiceReadRepository_GetInvoicesForPayNotifications(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	referenceDate := utils.Now()
	yesterday := referenceDate.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: yesterday,
		Customer: entity.InvoiceCustomer{
			Email: "email",
		},
		DryRun: false,
		Status: enum.InvoiceStatusDue,
	})

	result, err := repositories.InvoiceReadRepository.GetInvoicesForPayNotifications(ctx, 60, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, invoiceId, props["id"])
}

func TestInvoiceReadRepository_GetInvoicesForPayNotifications_InvoiceIsDryRun(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	referenceDate := utils.Now()
	yesterday := referenceDate.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: yesterday,
		Customer: entity.InvoiceCustomer{
			Email: "email",
		},
		DryRun: false,
		Status: enum.InvoiceStatusDue,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: yesterday,
		Customer: entity.InvoiceCustomer{
			Email: "email",
		},
		DryRun: true,
		Status: enum.InvoiceStatusDue,
	})

	result, err := repositories.InvoiceReadRepository.GetInvoicesForPayNotifications(ctx, 60, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, invoiceId, props["id"])
}

func TestInvoiceReadRepository_GetInvoicesForPayNotifications_StatusIsDraft(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	referenceDate := utils.Now()
	yesterday := referenceDate.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: yesterday,
		Customer: entity.InvoiceCustomer{
			Email: "email",
		},
		DryRun: false,
		Status: enum.InvoiceStatusDue,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: yesterday,
		Customer: entity.InvoiceCustomer{
			Email: "email",
		},
		DryRun: false,
		Status: enum.InvoiceStatusDraft,
	})

	result, err := repositories.InvoiceReadRepository.GetInvoicesForPayNotifications(ctx, 60, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, invoiceId, props["id"])
}

func TestInvoiceReadRepository_GetInvoicesForPayNotifications_StatusIsPaid(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	referenceDate := utils.Now()
	yesterday := referenceDate.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: yesterday,
		Customer: entity.InvoiceCustomer{
			Email: "email",
		},
		DryRun: false,
		Status: enum.InvoiceStatusDue,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: yesterday,
		Customer: entity.InvoiceCustomer{
			Email: "email",
		},
		DryRun: false,
		Status: enum.InvoiceStatusPaid,
	})

	result, err := repositories.InvoiceReadRepository.GetInvoicesForPayNotifications(ctx, 60, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, invoiceId, props["id"])
}

func TestInvoiceReadRepository_GetInvoicesForPayNotifications_MissingCustomerEmail(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	referenceDate := utils.Now()
	yesterday := referenceDate.Add(-24 * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: yesterday,
		Customer: entity.InvoiceCustomer{
			Email: "email",
		},
		DryRun: false,
		Status: enum.InvoiceStatusDue,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: yesterday,
		DryRun:    false,
		Status:    enum.InvoiceStatusDue,
	})

	result, err := repositories.InvoiceReadRepository.GetInvoicesForPayNotifications(ctx, 60, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, invoiceId, props["id"])
}

func TestInvoiceReadRepository_GetInvoicesForPayNotifications_RecentlyUpdated(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tenant := "tenant1"

	referenceDate := utils.Now()
	yesterday := referenceDate.Add(-24 * time.Hour)
	minAgo10 := referenceDate.Add(-10 * time.Minute)

	neo4jtest.CreateTenant(ctx, driver, tenant)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, entity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenant, organizationId, entity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: yesterday,
		Customer: entity.InvoiceCustomer{
			Email: "email",
		},
		DryRun: false,
		Status: enum.InvoiceStatusDue,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenant, contractId, entity.InvoiceEntity{
		UpdatedAt: minAgo10,
		Customer: entity.InvoiceCustomer{
			Email: "email",
		},
		DryRun: false,
		Status: enum.InvoiceStatusDue,
	})

	result, err := repositories.InvoiceReadRepository.GetInvoicesForPayNotifications(ctx, 60, referenceDate)
	require.NoError(t, err)
	require.Len(t, result, 1)
	props := utils.GetPropsFromNode(*result[0].Node)
	require.Equal(t, invoiceId, props["id"])
}
