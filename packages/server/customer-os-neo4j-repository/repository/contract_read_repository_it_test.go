package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestContractReadRepository_GetContractsForInvoicing_1_SLI_Monthly_V1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organization1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	currentYear := 2023

	contractId := neo4jtest.CreateContract(ctx, driver, tenantName, organization1Id, entity.ContractEntity{})

	sli1StartedAt := neo4jtest.FirstTimeOfMonth(currentYear, 1)
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, enum.BilledTypeMonthly, 12, 2, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelOrganization:    1,
		neo4jutil.NodeLabelContract:        1,
		neo4jutil.NodeLabelServiceLineItem: 1,
	})

	//firstDayOfJanuary := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)
	//lastDayOfJanuary := neo4jtest.LastTimeOfMonth(currentYear, 1)
	//
	////no invoice on the 1st of january
	//assertNoOrganization(t, ctx, firstDayOfJanuary)
	//assertNoOrganization(t, ctx, firstDayOfJanuary.Add(time.Hour*24).Add(-time.Nanosecond))
	//
	////1 invoice on first nanosecond on second day of january
	//assertInvoiceForOrganization(t, ctx, firstDayOfJanuary.Add(time.Hour*24), tenantName, organization1Id)
	//
	////1 invoice on last nanosecond on second day of january
	//assertInvoiceForOrganization(t, ctx, firstDayOfJanuary.Add(time.Hour*48).Add(-time.Nanosecond), tenantName, organization1Id)
	//assertInvoiceForOrganization(t, ctx, firstDayOfJanuary.Add(time.Hour*120), tenantName, organization1Id)
	//assertInvoiceForOrganization(t, ctx, lastDayOfJanuary, tenantName, organization1Id)
	//
	////todo mark SLI as included in an invoice in january
	////check that the invoiced is not created anymore in january
	//
	//assertNoOrganization(t, ctx, firstDayOfJanuary.Add(time.Hour*24))
	//assertNoOrganization(t, ctx, firstDayOfJanuary.Add(time.Hour*48))
	//assertNoOrganization(t, ctx, lastDayOfJanuary)
	//
	////1 invoice on the 1st of every month starting with february
	//for month := time.February; month <= time.December; month++ {
	//	firstDayOfMonth := time.Date(currentYear, month, 1, 0, 0, 0, 0, time.UTC)
	//	assertInvoiceForOrganization(t, ctx, firstDayOfMonth, tenantName, organization1Id)
	//}
	//
	////no invoices on any other day
	//for month := time.January; month <= time.December; month++ {
	//	for day := 2; day <= daysInMonth(currentYear, month); day++ {
	//		currentDate := time.Date(currentYear, month, day, 0, 0, 0, 0, time.UTC)
	//
	//		assertNoOrganization(t, ctx, currentDate)
	//	}
	//}
}

func TestContractReadRepository_GetContractsForInvoicing_1_SLI_Monthly_V2(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organization1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractId := neo4jtest.CreateContract(ctx, driver, tenantName, organization1Id, entity.ContractEntity{})

	sli1StartedAt := neo4jtest.FirstTimeOfMonth(2023, 1)
	sli1EndedAt := neo4jtest.MiddleTimeOfMonth(2023, 3)

	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, enum.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1EndedAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, enum.BilledTypeMonthly, 4, 1, enum.BilledTypeMonthly, 1, 1, sli1EndedAt, sli1Id)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelOrganization:    1,
		neo4jutil.NodeLabelContract:        1,
		neo4jutil.NodeLabelServiceLineItem: 2,
	})

	//1 invoice every month on the 1st
	//1 invoice in middle of March for proration of 4$/31 days * 15 days = 1.935483870967742 - 1$/31 days * 15 days = 0.4838709677419355
	//
	//currentYear := 2023
	//for month := time.January; month <= time.December; month++ {
	//	firstDayOfMonth := time.Date(currentYear, month, 1, 0, 0, 0, 0, time.UTC)
	//
	//	//first nanosecond of the month on day 1
	//	assertInvoiceForOrganization(t, ctx, firstDayOfMonth, tenantName, organization1Id)
	//	//last nanosecond of the month on day 1
	//	assertInvoiceForOrganization(t, ctx, firstDayOfMonth.Add(time.Hour*24).Add(-time.Nanosecond), tenantName, organization1Id)
	//
	//	for day := 2; day <= daysInMonth(currentYear, month); day++ {
	//		if month == time.March && day == 15 {
	//			assertInvoiceForOrganization(t, ctx, time.Date(currentYear, month, day, 0, 0, 0, 0, time.UTC), tenantName, organization1Id)
	//			continue
	//		}
	//
	//		currentDate := time.Date(currentYear, month, day, 0, 0, 0, 0, time.UTC)
	//
	//		fmt.Sprint(currentDate)
	//		assertNoOrganization(t, ctx, currentDate)
	//	}
	//}
}

func assertNoContract(t *testing.T, ctx context.Context, invoiceDate time.Time) {
	contractsForInvoicing, err := repositories.ContractReadRepository.GetContractsForInvoicing(ctx, invoiceDate)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, 0, len(contractsForInvoicing))
}

func assertInvoiceForContract(t *testing.T, ctx context.Context, invoiceDate time.Time, tenant, organizationId string) {
	contractsForInvoicing, err := repositories.ContractReadRepository.GetContractsForInvoicing(ctx, invoiceDate)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, 1, len(contractsForInvoicing))

	for _, dbNode := range contractsForInvoicing {
		tn, org := getRowData(dbNode)

		require.Equal(t, tenant, tn)
		require.Equal(t, organizationId, org)
	}
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func getRowData(props map[string]any) (string, string) {
	tenant := utils.GetStringPropOrEmpty(props, "tenant")
	organizationId := utils.GetStringPropOrEmpty(props, "organizationId")

	return tenant, organizationId
}
