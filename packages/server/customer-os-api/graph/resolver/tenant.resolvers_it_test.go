package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_GetTenantBillingProfiles(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	profileId := neo4jtest.CreateTenantBillingProfile(ctx, driver, tenantName, neo4jentity.TenantBillingProfileEntity{
		Phone:                  "123456789",
		LegalName:              "test",
		AddressLine1:           "address1",
		AddressLine2:           "address2",
		AddressLine3:           "address3",
		Locality:               "locality",
		Country:                "country",
		Region:                 "region",
		Zip:                    "zip",
		VatNumber:              "vatNumber",
		SendInvoicesFrom:       "sendInvoicesFrom",
		SendInvoicesBcc:        "sendInvoicesBcc",
		CanPayWithPigeon:       true,
		CanPayWithBankTransfer: true,
		Check:                  true,
	})

	rawResponse, err := c.RawPost(getQuery("tenant/get_tenant_billing_profiles"))
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantGraphqlResponse struct {
		TenantBillingProfiles []model.TenantBillingProfile
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantGraphqlResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantGraphqlResponse)

	require.Equal(t, 1, len(tenantGraphqlResponse.TenantBillingProfiles))
	tenantBillingProfile := tenantGraphqlResponse.TenantBillingProfiles[0]
	require.Equal(t, profileId, tenantBillingProfile.ID)
	require.Equal(t, "123456789", tenantBillingProfile.Phone)
	require.Equal(t, "test", tenantBillingProfile.LegalName)
	require.Equal(t, "address1", tenantBillingProfile.AddressLine1)
	require.Equal(t, "address2", tenantBillingProfile.AddressLine2)
	require.Equal(t, "address3", tenantBillingProfile.AddressLine3)
	require.Equal(t, "locality", tenantBillingProfile.Locality)
	require.Equal(t, "country", tenantBillingProfile.Country)
	require.Equal(t, "region", tenantBillingProfile.Region)
	require.Equal(t, "zip", tenantBillingProfile.Zip)
	require.Equal(t, "vatNumber", tenantBillingProfile.VatNumber)
	require.Equal(t, "sendInvoicesFrom", tenantBillingProfile.SendInvoicesFrom)
	require.Equal(t, "sendInvoicesBcc", tenantBillingProfile.SendInvoicesBcc)
	require.True(t, tenantBillingProfile.CanPayWithPigeon)
	require.True(t, tenantBillingProfile.CanPayWithBankTransfer)
	require.True(t, tenantBillingProfile.Check)
}

func TestQueryResolver_GetTenantBillingProfile(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	profileId := neo4jtest.CreateTenantBillingProfile(ctx, driver, tenantName, neo4jentity.TenantBillingProfileEntity{
		Phone:                  "123456789",
		LegalName:              "test",
		AddressLine1:           "address1",
		AddressLine2:           "address2",
		AddressLine3:           "address3",
		Locality:               "locality",
		Country:                "country",
		Region:                 "region",
		Zip:                    "zip",
		VatNumber:              "vatNumber",
		SendInvoicesFrom:       "sendInvoicesFrom",
		SendInvoicesBcc:        "sendInvoicesBcc",
		CanPayWithPigeon:       true,
		CanPayWithBankTransfer: true,
		Check:                  true,
	})

	rawResponse, err := c.RawPost(getQuery("tenant/get_tenant_billing_profile"), client.Var("id", profileId))
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantGraphqlResponse struct {
		TenantBillingProfile model.TenantBillingProfile
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantGraphqlResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantGraphqlResponse)

	tenantBillingProfile := tenantGraphqlResponse.TenantBillingProfile
	require.Equal(t, profileId, tenantBillingProfile.ID)
	require.Equal(t, "123456789", tenantBillingProfile.Phone)
	require.Equal(t, "test", tenantBillingProfile.LegalName)
	require.Equal(t, "address1", tenantBillingProfile.AddressLine1)
	require.Equal(t, "address2", tenantBillingProfile.AddressLine2)
	require.Equal(t, "address3", tenantBillingProfile.AddressLine3)
	require.Equal(t, "locality", tenantBillingProfile.Locality)
	require.Equal(t, "country", tenantBillingProfile.Country)
	require.Equal(t, "region", tenantBillingProfile.Region)
	require.Equal(t, "zip", tenantBillingProfile.Zip)
	require.Equal(t, "vatNumber", tenantBillingProfile.VatNumber)
	require.Equal(t, "sendInvoicesFrom", tenantBillingProfile.SendInvoicesFrom)
	require.Equal(t, "sendInvoicesBcc", tenantBillingProfile.SendInvoicesBcc)
	require.True(t, tenantBillingProfile.CanPayWithPigeon)
	require.True(t, tenantBillingProfile.CanPayWithBankTransfer)
	require.True(t, tenantBillingProfile.Check)
}

func TestQueryResolver_GetTenantSettings(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenantSettings(ctx, driver, tenantName, neo4jentity.TenantSettingsEntity{
		LogoRepositoryFileId: "logoRepositoryFileId",
		BaseCurrency:         neo4jenum.CurrencyUSD,
		InvoicingEnabled:     true,
	})

	rawResponse, err := c.RawPost(getQuery("tenant/get_tenant_settings"))
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantGraphqlResponse struct {
		TenantSettings model.TenantSettings
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantGraphqlResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantGraphqlResponse)

	tenantSettings := tenantGraphqlResponse.TenantSettings
	require.Equal(t, "logoRepositoryFileId", *tenantSettings.LogoRepositoryFileID)
	require.Equal(t, model.CurrencyUsd, *tenantSettings.BaseCurrency)
	require.Equal(t, true, tenantSettings.BillingEnabled)
}

//func TestMutationResolver_TenantHardDelete(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//
//	insertTenantDataWithNodeChecks(ctx, t, "tenant1", testUserId)
//	insertTenantDataWithNodeChecks(ctx, t, "tenant2", "1")
//
//	rawResponseTenant2, err := cCustomerOsPlatformOwner.RawPost(getQuery("tenant/hard_delete_tenant"),
//		client.Var("tenant", "tenant2"),
//		client.Var("confirmTenant", "tenant2"))
//	assertRawResponseSuccess(t, rawResponseTenant2, err)
//
//	response2 := map[string]interface{}{}
//
//	err = decode.Decode(rawResponseTenant2.Data.(map[string]any), &response2)
//	require.Nil(t, err)
//
//	tenantDataNodeChecks(ctx, t, "tenant1", 1)
//	tenantDataNodeChecks(ctx, t, "tenant2", 0)
//
//	rawResponseTenant1, err := cCustomerOsPlatformOwner.RawPost(getQuery("tenant/hard_delete_tenant"),
//		client.Var("tenant", "tenant1"),
//		client.Var("confirmTenant", "tenant1"))
//	assertRawResponseSuccess(t, rawResponseTenant2, err)
//
//	response1 := map[string]interface{}{}
//
//	err = decode.Decode(rawResponseTenant1.Data.(map[string]any), &response1)
//	require.Nil(t, err)
//
//	tenantDataNodeChecks(ctx, t, "tenant1", 0)
//	tenantDataNodeChecks(ctx, t, "tenant2", 0)
//
//	tenantNodesCount := neo4jtest.GetCountOfNodes(ctx, driver, "Tenant")
//	require.Equal(t, 0, tenantNodesCount)
//}

func insertTenantDataWithNodeChecks(ctx context.Context, t *testing.T, tenant, userId string) {
	neo4jtest.CreateTenant(ctx, driver, tenant)
	neo4jtest.CreateExternalSystem(ctx, driver, tenant, neo4jentity.ExternalSystemEntity{})
	neo4jtest.CreateWorkspace(ctx, driver, "testworkspace", "testprovider", tenant)
	neo4jtest.CreateUser(ctx, driver, tenant, neo4jentity.UserEntity{
		Id:    userId,
		Roles: []string{"PLATFORM_OWNER"},
	})
	neo4jtest.CreateTenantSettings(ctx, driver, tenant, neo4jentity.TenantSettingsEntity{})
	neo4jtest.CreateTenantBillingProfile(ctx, driver, tenant, neo4jentity.TenantBillingProfileEntity{})
	neo4jtest.CreateBillingProfile(ctx, driver, tenant, neo4jentity.BillingProfileEntity{})
	neo4jtest.CreateBankAccount(ctx, driver, tenant, neo4jentity.BankAccountEntity{})
	neo4jtest.CreateContact(ctx, driver, tenant, neo4jentity.ContactEntity{})
	neo4jtest.CreateSocial(ctx, driver, tenant, neo4jentity.SocialEntity{})
	neo4jtest.CreateEmail(ctx, driver, tenant, neo4jentity.EmailEntity{})
	neo4jtest.CreateLogEntry(ctx, driver, tenant, neo4jentity.LogEntryEntity{})
	neo4jtest.CreateTag(ctx, driver, tenant, neo4jentity.TagEntity{})

	neo4jt.CreateMeeting(ctx, driver, tenant, "", utils.Now())
	neo4jt.CreateAttachment(ctx, driver, tenant, neo4jentity.AttachmentEntity{})
	neo4jtest.CreatePhoneNumber(ctx, driver, tenant, neo4jentity.PhoneNumberEntity{})
	neo4jt.CreateIssue(ctx, driver, tenant, entity.IssueEntity{})
	neo4jtest.CreateInteractionEvent(ctx, driver, tenant, "1", "c", "", "EMAIL", utils.Now())
	neo4jtest.CreateInteractionSession(ctx, driver, tenant, "1", "c", "", "", "", utils.Now(), true)

	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenant, neo4jentity.OrganizationEntity{})
	neo4jtest.CreateReminder(ctx, driver, tenant, testUserId, organizationId, utils.Now(), neo4jentity.ReminderEntity{})
	neo4jt.CreateActionForOrganization(ctx, driver, tenant, organizationId, enum.ActionCreated, utils.Now())

	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenant, organizationId, neo4jentity.ContractEntity{}, neo4jentity.OpportunityEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenant, contractId, neo4jentity.ServiceLineItemEntity{})

	tenantDataNodeChecks(ctx, t, tenant, 1)
}

func tenantDataNodeChecks(ctx context.Context, t *testing.T, tenant string, numberOfNodes int) {
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		commonModel.NodeLabelExternalSystem + "_" + tenant:       numberOfNodes,
		commonModel.NodeLabelUser + "_" + tenant:                 numberOfNodes,
		commonModel.NodeLabelTenantBillingProfile + "_" + tenant: numberOfNodes,
		commonModel.NodeLabelBillingProfile + "_" + tenant:       numberOfNodes,
		commonModel.NodeLabelBankAccount + "_" + tenant:          numberOfNodes,
		commonModel.NodeLabelContact + "_" + tenant:              numberOfNodes,
		commonModel.NodeLabelSocial + "_" + tenant:               numberOfNodes,
		commonModel.NodeLabelEmail + "_" + tenant:                numberOfNodes,
		commonModel.NodeLabelLogEntry + "_" + tenant:             numberOfNodes,
		commonModel.NodeLabelOrganization + "_" + tenant:         numberOfNodes,
		commonModel.NodeLabelReminder + "_" + tenant:             numberOfNodes,
		commonModel.NodeLabelMeeting + "_" + tenant:              numberOfNodes,
		commonModel.NodeLabelAttachment + "_" + tenant:           numberOfNodes,
		commonModel.NodeLabelIssue + "_" + tenant:                numberOfNodes,
		commonModel.NodeLabelPhoneNumber + "_" + tenant:          numberOfNodes,
		commonModel.NodeLabelAction + "_" + tenant:               numberOfNodes,
		commonModel.NodeLabelTag + "_" + tenant:                  numberOfNodes,
		commonModel.NodeLabelContract + "_" + tenant:             numberOfNodes,
		commonModel.NodeLabelOpportunity + "_" + tenant:          numberOfNodes,
		commonModel.NodeLabelServiceLineItem + "_" + tenant:      numberOfNodes,
		commonModel.NodeLabelInteractionEvent + "_" + tenant:     numberOfNodes,
		commonModel.NodeLabelInteractionSession + "_" + tenant:   numberOfNodes,
	})
}
