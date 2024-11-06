package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
	"time"
)

func TestQueryResolver_Contract_WithServiceLineItems(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	tomorrow := now.Add(time.Duration(24) * time.Hour)
	yesterday := now.Add(time.Duration(-24) * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		AddressLine1:           "address line 1",
		AddressLine2:           "address line 2",
		Zip:                    "zip",
		Locality:               "locality",
		Country:                "country",
		Region:                 "region",
		OrganizationLegalName:  "organization legal name",
		InvoiceEmail:           "invoice email",
		InvoiceNote:            "invoice note",
		BillingCycleInMonths:   1,
		InvoicingStartDate:     &now,
		NextInvoiceDate:        &tomorrow,
		InvoicingEnabled:       true,
		CanPayWithCard:         true,
		CanPayWithDirectDebit:  true,
		CanPayWithBankTransfer: true,
		PayOnline:              true,
		PayAutomatically:       true,
		AutoRenew:              true,
		Check:                  true,
		Approved:               true,
		DueDays:                int64(7),
	})

	serviceLineItemId1 := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:      "service line item 1",
		CreatedAt: yesterday,
		UpdatedAt: yesterday,
		Billed:    neo4jenum.BilledTypeAnnually,
		Price:     13,
		Quantity:  2,
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: "test1",
		VatRate:   0.1,
	})
	serviceLineItemId2 := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:      "service line item 2",
		CreatedAt: now,
		UpdatedAt: now,
		Billed:    neo4jenum.BilledTypeUsage,
		Price:     255,
		Quantity:  23,
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: "test2",
		VatRate:   0.2,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Organization":    1,
		"Contract":        1,
		"ServiceLineItem": 2,
	})
	neo4jtest.AssertRelationship(ctx, t, driver, contractId, "HAS_SERVICE", serviceLineItemId1)
	neo4jtest.AssertRelationship(ctx, t, driver, contractId, "HAS_SERVICE", serviceLineItemId2)

	rawResponse := callGraphQL(t, "contract/get_contract_with_service_line_items",
		map[string]interface{}{"contractId": contractId})

	var contractStruct struct {
		Contract model.Contract
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contractStruct)
	require.Nil(t, err)

	contract := contractStruct.Contract
	require.NotNil(t, contract)
	require.Equal(t, contractId, contract.Metadata.ID)
	require.True(t, contract.BillingEnabled)
	require.True(t, contract.AutoRenew)
	require.True(t, contract.Approved)

	billingDetails := contract.BillingDetails
	require.Equal(t, "address line 1", *billingDetails.AddressLine1)
	require.Equal(t, "address line 2", *billingDetails.AddressLine2)
	require.Equal(t, "zip", *billingDetails.PostalCode)
	require.Equal(t, "locality", *billingDetails.Locality)
	require.Equal(t, "country", *billingDetails.Country)
	require.Equal(t, "region", *billingDetails.Region)
	require.Equal(t, "organization legal name", *billingDetails.OrganizationLegalName)
	require.Equal(t, "invoice email", *billingDetails.BillingEmail)
	require.Equal(t, "invoice note", *billingDetails.InvoiceNote)
	require.Equal(t, model.ContractBillingCycleMonthlyBilling, *billingDetails.BillingCycle)
	require.True(t, *billingDetails.CanPayWithCard)
	require.True(t, *billingDetails.CanPayWithDirectDebit)
	require.True(t, *billingDetails.CanPayWithBankTransfer)
	require.True(t, *billingDetails.Check)
	require.True(t, *billingDetails.PayOnline)
	require.True(t, *billingDetails.PayAutomatically)
	require.Equal(t, utils.ToDate(now), *billingDetails.InvoicingStarted)
	require.Equal(t, utils.ToDate(tomorrow), *billingDetails.NextInvoicing)
	require.Equal(t, int64(7), *billingDetails.DueDays)

	require.Equal(t, 2, len(contract.ContractLineItems))

	firstContractLineItem := contract.ContractLineItems[0]
	require.Equal(t, serviceLineItemId1, firstContractLineItem.Metadata.ID)
	require.Equal(t, "service line item 1", firstContractLineItem.Description)
	require.Equal(t, yesterday, firstContractLineItem.Metadata.Created)
	require.Equal(t, yesterday, firstContractLineItem.Metadata.LastUpdated)
	require.Equal(t, model.BilledTypeAnnually, firstContractLineItem.BillingCycle)
	require.Equal(t, float64(13), firstContractLineItem.Price)
	require.Equal(t, int64(2), firstContractLineItem.Quantity)
	require.Equal(t, model.DataSourceOpenline, firstContractLineItem.Metadata.Source)
	require.Equal(t, "test1", firstContractLineItem.Metadata.AppSource)
	require.Equal(t, 0.1, firstContractLineItem.Tax.TaxRate)

	secondContractLineItem := contract.ContractLineItems[1]
	require.Equal(t, serviceLineItemId2, secondContractLineItem.Metadata.ID)
	require.Equal(t, "service line item 2", secondContractLineItem.Description)
	require.Equal(t, now, secondContractLineItem.Metadata.Created)
	require.Equal(t, now, secondContractLineItem.Metadata.LastUpdated)
	require.Equal(t, model.BilledTypeUsage, secondContractLineItem.BillingCycle)
	require.Equal(t, float64(255), secondContractLineItem.Price)
	require.Equal(t, int64(23), secondContractLineItem.Quantity)
	require.Equal(t, model.DataSourceOpenline, secondContractLineItem.Metadata.Source)
	require.Equal(t, "test2", secondContractLineItem.Metadata.AppSource)
	require.Equal(t, 0.2, secondContractLineItem.Tax.TaxRate)
}

func TestQueryResolver_Contract_WithInvoices(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	tomorrow := now.Add(time.Duration(24) * time.Hour)
	yesterday := now.Add(time.Duration(-24) * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		AddressLine1:           "address line 1",
		AddressLine2:           "address line 2",
		Zip:                    "zip",
		Locality:               "locality",
		Country:                "country",
		Region:                 "region",
		OrganizationLegalName:  "organization legal name",
		InvoiceEmail:           "invoice email",
		InvoiceNote:            "invoice note",
		BillingCycleInMonths:   1,
		InvoicingStartDate:     &now,
		NextInvoiceDate:        &tomorrow,
		InvoicingEnabled:       true,
		CanPayWithCard:         true,
		CanPayWithDirectDebit:  true,
		CanPayWithBankTransfer: true,
		PayOnline:              true,
		PayAutomatically:       true,
		AutoRenew:              true,
		Check:                  true,
		DueDays:                int64(7),
	})

	serviceLineItemId1 := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:      "service line item 1",
		CreatedAt: yesterday,
		UpdatedAt: yesterday,
		Billed:    neo4jenum.BilledTypeAnnually,
		Price:     13,
		Quantity:  2,
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: "test1",
		VatRate:   0.1,
	})
	serviceLineItemId2 := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Name:      "service line item 2",
		CreatedAt: now,
		UpdatedAt: now,
		Billed:    neo4jenum.BilledTypeUsage,
		Price:     255,
		Quantity:  23,
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: "test2",
		VatRate:   0.2,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Organization":    1,
		"Contract":        1,
		"ServiceLineItem": 2,
	})
	neo4jtest.AssertRelationship(ctx, t, driver, contractId, "HAS_SERVICE", serviceLineItemId1)
	neo4jtest.AssertRelationship(ctx, t, driver, contractId, "HAS_SERVICE", serviceLineItemId2)

	rawResponse := callGraphQL(t, "contract/get_contract_with_service_line_items",
		map[string]interface{}{"contractId": contractId})

	var contractStruct struct {
		Contract model.Contract
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contractStruct)
	require.Nil(t, err)

	contract := contractStruct.Contract
	require.NotNil(t, contract)
	require.Equal(t, contractId, contract.Metadata.ID)
	require.True(t, contract.BillingEnabled)
	require.True(t, contract.AutoRenew)

	billingDetails := contract.BillingDetails
	require.Equal(t, "address line 1", *billingDetails.AddressLine1)
	require.Equal(t, "address line 2", *billingDetails.AddressLine2)
	require.Equal(t, "zip", *billingDetails.PostalCode)
	require.Equal(t, "locality", *billingDetails.Locality)
	require.Equal(t, "country", *billingDetails.Country)
	require.Equal(t, "region", *billingDetails.Region)
	require.Equal(t, "organization legal name", *billingDetails.OrganizationLegalName)
	require.Equal(t, "invoice email", *billingDetails.BillingEmail)
	require.Equal(t, "invoice note", *billingDetails.InvoiceNote)
	require.Equal(t, model.ContractBillingCycleMonthlyBilling, *billingDetails.BillingCycle)
	require.True(t, *billingDetails.CanPayWithCard)
	require.True(t, *billingDetails.CanPayWithDirectDebit)
	require.True(t, *billingDetails.CanPayWithBankTransfer)
	require.True(t, *billingDetails.Check)
	require.True(t, *billingDetails.PayOnline)
	require.True(t, *billingDetails.PayAutomatically)
	require.Equal(t, utils.ToDate(now), *billingDetails.InvoicingStarted)
	require.Equal(t, utils.ToDate(tomorrow), *billingDetails.NextInvoicing)
	require.Equal(t, int64(7), *billingDetails.DueDays)

	require.Equal(t, 2, len(contract.ContractLineItems))

	firstContractLineItem := contract.ContractLineItems[0]
	require.Equal(t, serviceLineItemId1, firstContractLineItem.Metadata.ID)
	require.Equal(t, "service line item 1", firstContractLineItem.Description)
	require.Equal(t, yesterday, firstContractLineItem.Metadata.Created)
	require.Equal(t, yesterday, firstContractLineItem.Metadata.LastUpdated)
	require.Equal(t, model.BilledTypeAnnually, firstContractLineItem.BillingCycle)
	require.Equal(t, float64(13), firstContractLineItem.Price)
	require.Equal(t, int64(2), firstContractLineItem.Quantity)
	require.Equal(t, model.DataSourceOpenline, firstContractLineItem.Metadata.Source)
	require.Equal(t, "test1", firstContractLineItem.Metadata.AppSource)
	require.Equal(t, 0.1, firstContractLineItem.Tax.TaxRate)

	secondContractLineItem := contract.ContractLineItems[1]
	require.Equal(t, serviceLineItemId2, secondContractLineItem.Metadata.ID)
	require.Equal(t, "service line item 2", secondContractLineItem.Description)
	require.Equal(t, now, secondContractLineItem.Metadata.Created)
	require.Equal(t, now, secondContractLineItem.Metadata.LastUpdated)
	require.Equal(t, model.BilledTypeUsage, secondContractLineItem.BillingCycle)
	require.Equal(t, float64(255), secondContractLineItem.Price)
	require.Equal(t, int64(23), secondContractLineItem.Quantity)
	require.Equal(t, model.DataSourceOpenline, secondContractLineItem.Metadata.Source)
	require.Equal(t, "test2", secondContractLineItem.Metadata.AppSource)
	require.Equal(t, 0.2, secondContractLineItem.Tax.TaxRate)
}

func TestQueryResolver_Contract_WithOpportunities(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	yesterday := now.AddDate(0, 0, -1)
	twoDaysAgo := now.AddDate(0, 0, -2)
	threeDaysAgo := now.AddDate(0, 0, -3)
	fourDaysAgo := now.AddDate(0, 0, -4)
	tomorrow := now.AddDate(0, 0, 1)
	afterTomorrow := now.AddDate(0, 0, 2)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: yesterday,
		DryRun:     true,
		Status:     neo4jenum.InvoiceStatusDue,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: twoDaysAgo,
		DryRun:     false,
		Status:     neo4jenum.InvoiceStatusInitialized,
	})
	invoiceIdPaid := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: threeDaysAgo,
		DryRun:     false,
		Status:     neo4jenum.InvoiceStatusPaid,
	})
	invoiceIdVoid := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: fourDaysAgo,
		DryRun:     false,
		Status:     neo4jenum.InvoiceStatusVoid,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: now,
		DryRun:     true,
		Preview:    false,
		Status:     neo4jenum.InvoiceStatusDue,
	})
	neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: now,
		DryRun:     true,
		Preview:    true,
		Status:     neo4jenum.InvoiceStatusInitialized,
	})
	invoiceIdScheduled1 := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: tomorrow,
		DryRun:     true,
		Preview:    true,
		Status:     neo4jenum.InvoiceStatusScheduled,
	})
	invoiceIdScheduled2 := neo4jtest.CreateInvoiceForContract(ctx, driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		IssuedDate: afterTomorrow,
		DryRun:     true,
		Preview:    true,
		Status:     neo4jenum.InvoiceStatusScheduled,
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		model2.NodeLabelOrganization: 1,
		model2.NodeLabelContract:     1,
		model2.NodeLabelInvoice:      8,
	})

	rawResponse := callGraphQL(t, "contract/get_contract_with_invoices",
		map[string]interface{}{"contractId": contractId})

	var contractStruct struct {
		Contract model.Contract
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contractStruct)
	require.Nil(t, err)

	contract := contractStruct.Contract
	require.NotNil(t, contract)

	require.Equal(t, 2, len(contract.Invoices))
	require.Equal(t, invoiceIdPaid, contract.Invoices[0].Metadata.ID)
	require.Equal(t, invoiceIdVoid, contract.Invoices[1].Metadata.ID)

	require.Equal(t, 2, len(contract.UpcomingInvoices))
	require.Equal(t, invoiceIdScheduled1, contract.UpcomingInvoices[0].Metadata.ID)
	require.Equal(t, invoiceIdScheduled2, contract.UpcomingInvoices[1].Metadata.ID)
}

func TestMutationResolver_ContractDelete(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})

	calledDeleteContractEvent := false

	contractCallbacks := events_platform.MockContractServiceCallbacks{
		SoftDeleteContract: func(context context.Context, contract *contractpb.SoftDeleteContractGrpcRequest) (*emptypb.Empty, error) {
			require.Equal(t, tenantName, contract.Tenant)
			require.Equal(t, contractId, contract.Id)
			require.Equal(t, testUserId, contract.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, contract.AppSource)
			calledDeleteContractEvent = true
			return &emptypb.Empty{}, nil
		},
	}
	events_platform.SetContractCallbacks(&contractCallbacks)

	rawResponse := callGraphQL(t, "contract/delete_contract", map[string]interface{}{
		"contractId": contractId,
	})

	var response struct {
		Contract_Delete model.DeleteResponse
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.True(t, response.Contract_Delete.Accepted)
	require.False(t, response.Contract_Delete.Completed)
	require.True(t, calledDeleteContractEvent)
}

func TestMutationResolver_AddAttachmentToContract(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})

	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		MimeType: "text/plain",
		FileName: "readme.txt",
	})

	rawResponse, err := c.RawPost(getQuery("contract/contract_add_attachment"),
		client.Var("contractId", contractId),
		client.Var("attachmentId", attachmentId))
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contract"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "INCLUDES"))

	var meeting struct {
		Contract_AddAttachment model.Contract
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &meeting)
	require.Nil(t, err)

	require.NotNil(t, meeting.Contract_AddAttachment.ID)
	require.Len(t, meeting.Contract_AddAttachment.Attachments, 1)
	require.Equal(t, meeting.Contract_AddAttachment.Attachments[0].ID, attachmentId)
}

func TestMutationResolver_RemoveAttachmentFromContract(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})

	attachmentId1 := neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		MimeType: "text/plain",
		FileName: "readme1.txt",
	})

	attachmentId2 := neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		MimeType: "text/plain",
		FileName: "readme2.txt",
	})

	addAttachment1Response, err := c.RawPost(getQuery("contract/contract_add_attachment"),
		client.Var("contractId", contractId),
		client.Var("attachmentId", attachmentId1))
	assertRawResponseSuccess(t, addAttachment1Response, err)

	addAttachment2Response, err := c.RawPost(getQuery("contract/contract_add_attachment"),
		client.Var("contractId", contractId),
		client.Var("attachmentId", attachmentId2))
	assertRawResponseSuccess(t, addAttachment2Response, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contract"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "INCLUDES"))

	removeAttachmentResponse, err := c.RawPost(getQuery("contract/contract_remove_attachment"),
		client.Var("contractId", contractId),
		client.Var("attachmentId", attachmentId2))
	assertRawResponseSuccess(t, removeAttachmentResponse, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contract"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "INCLUDES"))

	var meeting struct {
		Contract_RemoveAttachment model.Contract
	}

	err = decode.Decode(removeAttachmentResponse.Data.(map[string]any), &meeting)
	require.Nil(t, err)

	require.NotNil(t, meeting.Contract_RemoveAttachment.ID)
	require.Len(t, meeting.Contract_RemoveAttachment.Attachments, 1)
	require.Equal(t, meeting.Contract_RemoveAttachment.Attachments[0].ID, attachmentId1)
}

func TestMutationResolver_ContractRenew_NoActiveRenewalOpportunity_CreateRenewalOpportunity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		LengthInMonths: 1,
	})

	calledCreatedRenewalOpportunityGrpc := false

	opportunityCallbacks := events_platform.MockOpportunityServiceCallbacks{
		CreateRenewalOpportunity: func(context context.Context, opportunity *opportunitypb.CreateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, opportunity.Tenant)
			require.Equal(t, testUserId, opportunity.LoggedInUserId)
			require.Equal(t, contractId, opportunity.ContractId)
			require.Equal(t, constants.AppSourceCustomerOsApi, opportunity.SourceFields.AppSource)
			require.Equal(t, string(neo4jentity.DataSourceOpenline), opportunity.SourceFields.Source)
			calledCreatedRenewalOpportunityGrpc = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: uuid.New().String(),
			}, nil
		},
	}
	events_platform.SetOpportunityCallbacks(&opportunityCallbacks)

	rawResponse := callGraphQL(t, "contract/renew_contract", map[string]interface{}{
		"contractId": contractId,
	})

	var response struct {
		Contract_Renew model.Contract
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.Equal(t, contractId, response.Contract_Renew.Metadata.ID)
	require.True(t, calledCreatedRenewalOpportunityGrpc)
}

func TestMutationResolver_ContractRenew_ActiveRenewalNotExpired_ApproveRenewalOpportunity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	tomorrow := utils.Now().Add(time.Duration(24) * time.Hour)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		LengthInMonths: 1,
	})
	opportunityId := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewedAt:       &tomorrow,
			RenewalApproved: false,
		},
	})

	calledUpdateOpportunityGrpc := false

	opportunityCallbacks := events_platform.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunity: func(context context.Context, opportunity *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, opportunity.Tenant)
			require.Equal(t, testUserId, opportunity.LoggedInUserId)
			require.Equal(t, opportunityId, opportunity.Id)
			require.True(t, opportunity.RenewalApproved)
			require.Equal(t, constants.AppSourceCustomerOsApi, opportunity.SourceFields.AppSource)
			require.Equal(t, string(neo4jentity.DataSourceOpenline), opportunity.SourceFields.Source)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEW_APPROVED}, opportunity.FieldsMask)
			calledUpdateOpportunityGrpc = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: uuid.New().String(),
			}, nil
		},
	}
	events_platform.SetOpportunityCallbacks(&opportunityCallbacks)

	rawResponse := callGraphQL(t, "contract/renew_contract", map[string]interface{}{
		"contractId": contractId,
	})

	var response struct {
		Contract_Renew model.Contract
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.Equal(t, contractId, response.Contract_Renew.Metadata.ID)
	require.True(t, calledUpdateOpportunityGrpc)
}

func TestMutationResolver_ContractRenew_ActiveRenewalExpired_RolloutRenewalOpportunity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	yesterday := utils.Now().Add(time.Duration(-24) * time.Hour)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		LengthInMonths: 1,
	})
	neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewedAt:       &yesterday,
			RenewalApproved: false,
		},
	})

	calledRolloutRenewalOpportunity := false

	contractCallbacks := events_platform.MockContractServiceCallbacks{
		RolloutRenewalOpportunityOnExpiration: func(context context.Context, contract *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
			require.Equal(t, tenantName, contract.Tenant)
			require.Equal(t, testUserId, contract.LoggedInUserId)
			require.Equal(t, contractId, contract.Id)
			require.Equal(t, constants.AppSourceCustomerOsApi, contract.AppSource)
			calledRolloutRenewalOpportunity = true
			return &contractpb.ContractIdGrpcResponse{
				Id: contractId,
			}, nil
		},
	}
	events_platform.SetContractCallbacks(&contractCallbacks)

	rawResponse := callGraphQL(t, "contract/renew_contract", map[string]interface{}{
		"contractId": contractId,
	})

	var response struct {
		Contract_Renew model.Contract
	}

	require.Nil(t, rawResponse.Errors)
	err := decode.Decode(rawResponse.Data.(map[string]any), &response)
	require.Nil(t, err)
	require.Equal(t, contractId, response.Contract_Renew.Metadata.ID)
	require.True(t, calledRolloutRenewalOpportunity)
}

func TestQueryResolver_Contracts(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	org1 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	org2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	contract1 := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, org1, neo4jentity.ContractEntity{})
	contract2 := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, org1, neo4jentity.ContractEntity{})
	contract3 := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, org2, neo4jentity.ContractEntity{})

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, model2.NodeLabelOrganization))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, model2.NodeLabelContract))

	rawResponse, err := c.RawPost(getQuery("contract/get_contracts"),
		client.Var("page", 1),
		client.Var("limit", 4),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var contractsPageStruct struct {
		Contracts model.ContractPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contractsPageStruct)
	require.Nil(t, err)
	require.NotNil(t, contractsPageStruct)
	pagedOrganizations := contractsPageStruct.Contracts
	require.Equal(t, 1, pagedOrganizations.TotalPages)
	require.Equal(t, int64(3), pagedOrganizations.TotalElements)
	require.ElementsMatch(t, []string{contract1, contract2, contract3}, []string{pagedOrganizations.Content[0].Metadata.ID, pagedOrganizations.Content[1].Metadata.ID, pagedOrganizations.Content[2].Metadata.ID})
}
