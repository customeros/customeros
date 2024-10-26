package main

import (
	"context"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	"log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client/interceptor"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/comment"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	iepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_event"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/issue"
	logentrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/log_entry"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcApiKey = "082c1193-a5a2-42fc-87fc-e960e692fffd"
const appSource = "test_app"

var tenant = "customerosai"
var userId = "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
var orgId = "ceae019f-d1e3-49b3-87c5-35ebb68a5ff1"
var contractId = "769d1fb8-50a1-44bc-aff0-0f4338bd8ff2"

type Clients struct {
	InteractionEventClient iepb.InteractionEventGrpcServiceClient
	OrganizationClient     organizationpb.OrganizationGrpcServiceClient
	ContactClient          contactpb.ContactGrpcServiceClient
	EmailClient            emailpb.EmailGrpcServiceClient
	PhoneNumberClient      phonenumberpb.PhoneNumberGrpcServiceClient
	LogEntryClient         logentrypb.LogEntryGrpcServiceClient
	IssueClient            issuepb.IssueGrpcServiceClient
	CommentClient          commentpb.CommentGrpcServiceClient
	UserClient             userpb.UserGrpcServiceClient
	ContractClient         contractpb.ContractGrpcServiceClient
	ServiceLineItemClient  servicelineitempb.ServiceLineItemGrpcServiceClient
	OpportunityClient      opportunitypb.OpportunityGrpcServiceClient
	TenantClient           tenantpb.TenantGrpcServiceClient
	InvoiceClient          invoicepb.InvoiceGrpcServiceClient
	EventStoreClient       eventstorepb.EventStoreGrpcServiceClient
}

var clients *Clients

func InitClients() {
	conn, _ := grpc.Dial("localhost:5001", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			interceptor.ApiKeyEnricher(grpcApiKey),
		))
	clients = &Clients{
		InteractionEventClient: iepb.NewInteractionEventGrpcServiceClient(conn),
		OrganizationClient:     organizationpb.NewOrganizationGrpcServiceClient(conn),
		ContactClient:          contactpb.NewContactGrpcServiceClient(conn),
		LogEntryClient:         logentrypb.NewLogEntryGrpcServiceClient(conn),
		EmailClient:            emailpb.NewEmailGrpcServiceClient(conn),
		PhoneNumberClient:      phonenumberpb.NewPhoneNumberGrpcServiceClient(conn),
		IssueClient:            issuepb.NewIssueGrpcServiceClient(conn),
		CommentClient:          commentpb.NewCommentGrpcServiceClient(conn),
		UserClient:             userpb.NewUserGrpcServiceClient(conn),
		ContractClient:         contractpb.NewContractGrpcServiceClient(conn),
		OpportunityClient:      opportunitypb.NewOpportunityGrpcServiceClient(conn),
		ServiceLineItemClient:  servicelineitempb.NewServiceLineItemGrpcServiceClient(conn),
		TenantClient:           tenantpb.NewTenantGrpcServiceClient(conn),
		InvoiceClient:          invoicepb.NewInvoiceGrpcServiceClient(conn),
		EventStoreClient:       eventstorepb.NewEventStoreGrpcServiceClient(conn),
	}
}

func main() {
	InitClients()

	//testRequestGenerateSummaryRequest()
	//testRequestGenerateActionItemsRequest()
	//testCreateOrganization()
	//testEnrichOrganization()
	//testEnrichContact()
	//testAddSocialToContact()
	//testUpdateWithUpsertOrganization()
	//testUpdateOrganization()
	//testHideOrganization()
	//testShowOrganization()
	//testCreateLogEntry()
	//testUpdateLogEntry()
	//testAddCustomField()
	//testCreatePhoneNumber()
	//testAddParentOrganization()
	//testRemoveParentOrganization()
	//testContactLinkWithPhoneNumber()
	//testContactLinkWithLocation()
	//testOrganizationLinkWithEmail()
	//testOrganizationLinkWithPhoneNumber()
	//testOrganizationLinkWithLocation()
	//testContactLinkWithOrganization()
	//testCreateIssue()
	//testUpdateIssue()
	//testCreateComment()
	//testCreateContract()
	//testUpdateContract()
	//testAddContractService()
	//testCloseLooseOpportunity()
	//testCreateRenewalOpportunity()
	//testArchiveOpportunity()
	//testUpdateOnboardingStatus()
	//testUpdateOrgOwner()
	//testRefreshLastTouchpoint()
	//testRefreshRenewalSummary()
	//testAddTenantBillingProfile()
	//PaidInvoiceNotification()
	//PleasePayInvoiceNotification()
	//testCreateInvoice()
	//testTenantSettingsUpdate()
	//testCreateReminder()
	//testUpdateReminder()
	//testAddBankAccount()
}

func testEnrichOrganization() {
	organizationId := "0081162c-a80c-428c-b6ba-ae274ad81c9f"
	website := "openline.ai"

	result, err := clients.OrganizationClient.EnrichOrganization(context.Background(), &organizationpb.EnrichOrganizationGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
		Url:            website,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testEnrichContact() {
	contactId := "b497a882-c3c2-46ad-85ac-0b451503cd16"

	result, err := clients.ContactClient.EnrichContact(context.Background(), &contactpb.EnrichContactGrpcRequest{
		Tenant:    tenant,
		ContactId: contactId,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testAddSocialToContact() {
	contactId := "b497a882-c3c2-46ad-85ac-0b451503cd16"
	social := "https://www.linkedin.com/in/mateocafe/"

	result, err := clients.ContactClient.AddSocial(context.Background(), &contactpb.ContactAddSocialGrpcRequest{
		Tenant:    tenant,
		ContactId: contactId,
		Url:       social,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testCreateInvoice() {
	today := utils.Now()
	in1Month := today.AddDate(0, 1, 0)
	contractId := "769d1fb8-50a1-44bc-aff0-0f4338bd8ff2"
	result, err := clients.InvoiceClient.NewInvoiceForContract(context.Background(), &invoicepb.NewInvoiceForContractRequest{
		Tenant:               tenant,
		ContractId:           contractId,
		Currency:             "USD",
		InvoicePeriodStart:   utils.ConvertTimeToTimestampPtr(&today),
		InvoicePeriodEnd:     utils.ConvertTimeToTimestampPtr(&in1Month),
		OffCycle:             false,
		BillingCycleInMonths: 1,
		DryRun:               true,
		SourceFields: &commonpb.SourceFields{
			AppSource: appSource,
		},
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	print(result.Id)
}

func PaidInvoiceNotification() {
	_, err := clients.InvoiceClient.UpdateInvoice(context.Background(), &invoicepb.UpdateInvoiceRequest{
		Tenant:    tenant,
		InvoiceId: "5b052bf0-1027-4425-ba1e-4aa940754423",
		Status:    invoicepb.InvoiceStatus_INVOICE_STATUS_PAID,
	})

	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
}

func PleasePayInvoiceNotification() {
	_, err := clients.InvoiceClient.PayInvoiceNotification(context.Background(), &invoicepb.PayInvoiceNotificationRequest{
		Tenant:    tenant,
		InvoiceId: "e3af66b0-8e74-4aa7-941d-4b87518d7131",
	})

	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
}

func testAddTenantBillingProfile() {
	result, err := clients.TenantClient.AddBillingProfile(context.Background(), &tenantpb.AddBillingProfileRequest{
		Tenant:          tenant,
		SendInvoicesBcc: "invoice@openline.ai",
		LegalName:       "Openline",
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testRequestGenerateSummaryRequest() {

	interactionEventId := "555263fe-2e39-48f0-a8c2-c4c7a5ffb23d"

	result, _ := clients.InteractionEventClient.RequestGenerateSummary(context.Background(), &iepb.RequestGenerateSummaryGrpcRequest{
		Tenant:             tenant,
		InteractionEventId: interactionEventId,
	})
	print(result)
}

func testRequestGenerateActionItemsRequest() {

	interactionEventId := "555263fe-2e39-48f0-a8c2-c4c7a5ffb23d"

	result, _ := clients.InteractionEventClient.RequestGenerateActionItems(context.Background(), &iepb.RequestGenerateActionItemsGrpcRequest{
		Tenant:             tenant,
		InteractionEventId: interactionEventId,
	})
	print(result)
}

func testCreateOrganization() {
	userId := "697563a8-171c-4950-a067-1aaaaf2de1d8"
	website := ""

	result, err := clients.OrganizationClient.UpsertOrganization(context.Background(), &organizationpb.UpsertOrganizationGrpcRequest{
		Tenant:         tenant,
		Website:        website,
		LoggedInUserId: userId,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result)
}

func testUpdateWithUpsertOrganization() {

	organizationId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"
	website := "xtz.com"
	lastFoundingAmont := "1Million"

	result, _ := clients.OrganizationClient.UpsertOrganization(context.Background(), &organizationpb.UpsertOrganizationGrpcRequest{
		Tenant:            tenant,
		Id:                organizationId,
		Website:           website,
		LastFundingAmount: lastFoundingAmont,
	})
	print(result)
}

func testUpdateOrganization() {

	organizationId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"
	name := "xtz.com"

	result, _ := clients.OrganizationClient.UpdateOrganization(context.Background(), &organizationpb.UpdateOrganizationGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
		Name:           name,
		FieldsMask:     []organizationpb.OrganizationMaskField{organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_NAME},
	})
	print(result)
}

func testShowOrganization() {

	organizationId := "ccc"

	result, _ := clients.OrganizationClient.ShowOrganization(context.Background(), &organizationpb.OrganizationIdGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
	})
	print(result)
}

func testCreateLogEntry() {

	organizationId := "2829263d-b489-4e92-b0ba-b1bca9ff4d04"
	userId := "development@openline.ai"
	authorId := "c61f8af2-0e46-4464-a5db-ded8e4fe242f"

	result, _ := clients.LogEntryClient.UpsertLogEntry(context.Background(), &logentrypb.UpsertLogEntryGrpcRequest{
		Tenant:               tenant,
		LoggedOrganizationId: utils.StringPtr(organizationId),
		SourceFields: &commonpb.SourceFields{
			AppSource: "test_app",
		},
		AuthorUserId: utils.StringPtr(authorId),
		Content:      "I spoke with client",
		ContentType:  "text/plain",
		UserId:       userId,
	})
	print(result)
}

func testUpdateLogEntry() {

	userId := "development@openline.ai"
	logEntryId := "ccffe134-4bcd-4fa0-955f-c79b9e1a985f"

	result, _ := clients.LogEntryClient.UpsertLogEntry(context.Background(), &logentrypb.UpsertLogEntryGrpcRequest{
		Tenant:      tenant,
		Id:          logEntryId,
		Content:     "new content",
		ContentType: "text/plain2",
		UserId:      userId,
		StartedAt:   timestamppb.New(utils.Now()),
	})
	print(result)
}

func testAddCustomField() {

	organizationId := "5e72b6fb-5f20-4973-9b96-52f4543a0df3"
	userId := "development@openline.ai"
	result, _ := clients.OrganizationClient.UpsertCustomFieldToOrganization(context.Background(), &organizationpb.CustomFieldForOrganizationGrpcRequest{
		Tenant:                tenant,
		OrganizationId:        organizationId,
		UserId:                userId,
		CustomFieldTemplateId: utils.StringPtr("c70cd2fb-1c31-46fd-851c-2e47ceba508f"),
		CustomFieldName:       "CF1",
		CustomFieldDataType:   organizationpb.CustomFieldDataType_TEXT,
		CustomFieldValue: &organizationpb.CustomFieldValue{
			StringValue: utils.StringPtr("super secret value"),
		},
	})
	print(result)
}

func testCreatePhoneNumber() {

	userId := "697563a8-171c-4950-a067-1aaaaf2de1d8"
	rawPhoneNumber := "+12345"

	result, _ := clients.PhoneNumberClient.UpsertPhoneNumber(context.Background(), &phonenumberpb.UpsertPhoneNumberGrpcRequest{
		Tenant:         tenant,
		PhoneNumber:    rawPhoneNumber,
		LoggedInUserId: userId,
	})
	print(result)
}

func testAddParentOrganization() {

	orgId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"
	parentOrgId := ""
	relType := "store"
	userId := "697563a8-171c-4950-a067-1aaaaf2de1d8"
	result, err := clients.OrganizationClient.AddParentOrganization(context.Background(), &organizationpb.AddParentOrganizationGrpcRequest{
		Tenant:               tenant,
		OrganizationId:       orgId,
		ParentOrganizationId: parentOrgId,
		Type:                 relType,
		AppSource:            appSource,
		LoggedInUserId:       userId,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result)
}

func testRemoveParentOrganization() {

	orgId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"
	parentOrgId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	result, err := clients.OrganizationClient.RemoveParentOrganization(context.Background(), &organizationpb.RemoveParentOrganizationGrpcRequest{
		Tenant:               tenant,
		OrganizationId:       orgId,
		ParentOrganizationId: parentOrgId,
	})
	if err != nil {
		print(err)
	}
	print(result)
}

func testContactLinkWithLocation() {

	contactId := "dd7bd45e-d6d3-405c-a7ba-cd4127479c20"
	locationId := "bafff70d-7e45-49e5-8732-6e2a362a3ee9"

	result, _ := clients.ContactClient.LinkLocationToContact(context.Background(), &contactpb.LinkLocationToContactGrpcRequest{
		Tenant:     tenant,
		ContactId:  contactId,
		LocationId: locationId,
		AppSource:  appSource,
	})
	print(result)
}

func testContactLinkWithPhoneNumber() {

	contactId := "dd7bd45e-d6d3-405c-a7ba-cd4127479c20"
	phoneNumberId := "c21c0352-14d8-474a-afcd-167daa99e321"

	result, _ := clients.ContactClient.LinkPhoneNumberToContact(context.Background(), &contactpb.LinkPhoneNumberToContactGrpcRequest{
		Tenant:        tenant,
		ContactId:     contactId,
		PhoneNumberId: phoneNumberId,
		Primary:       true,
		Label:         "work",
		AppSource:     appSource,
	})
	print(result)
}

func testOrganizationLinkWithLocation() {

	orgId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"
	locationId := "bafff70d-7e45-49e5-8732-6e2a362a3ee9"

	result, _ := clients.OrganizationClient.LinkLocationToOrganization(context.Background(), &organizationpb.LinkLocationToOrganizationGrpcRequest{
		Tenant:         tenant,
		OrganizationId: orgId,
		LocationId:     locationId,
	})
	print(result)
}

func testOrganizationLinkWithPhoneNumber() {

	orgId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"
	phoneNumberId := "c21c0352-14d8-474a-afcd-167daa99e321"

	result, _ := clients.OrganizationClient.LinkPhoneNumberToOrganization(context.Background(), &organizationpb.LinkPhoneNumberToOrganizationGrpcRequest{
		Tenant:         tenant,
		OrganizationId: orgId,
		PhoneNumberId:  phoneNumberId,
		Primary:        true,
		Label:          "work",
	})
	print(result)
}

func testContactLinkWithOrganization() {
	contactId := "2f7660a8-a40b-4f21-b81f-1b73f025f79c"
	orgId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"
	jobRole := "CTO"

	result, _ := clients.ContactClient.LinkWithOrganization(context.Background(), &contactpb.LinkWithOrganizationGrpcRequest{
		Tenant:         tenant,
		ContactId:      contactId,
		OrganizationId: orgId,
		JobTitle:       jobRole,
		Primary:        true,
		Description:    "CEO of the company",
		SourceFields: &commonpb.SourceFields{
			AppSource: "integration.app",
			Source:    "hubspot",
		},
	})
	print(result)
}

func testCreateIssue() {

	userId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	subject := "test issue"
	description := "nice issue"
	status := "open"
	priority := "high"
	orgId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"

	result, err := clients.IssueClient.UpsertIssue(context.Background(), &issuepb.UpsertIssueGrpcRequest{
		Tenant:                   tenant,
		Subject:                  subject,
		Description:              description,
		Status:                   status,
		Priority:                 priority,
		LoggedInUserId:           userId,
		ReportedByOrganizationId: utils.StringPtr(orgId),
		SourceFields: &commonpb.SourceFields{
			AppSource: appSource,
		},
		ExternalSystemFields: &commonpb.ExternalSystemFields{
			ExternalSystemId: "hubspot",
			ExternalId:       "123",
		},
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Created issue id: %v", result.Id)
}

func testUpdateIssue() {

	userId := "697563a8-171c-4950-a067-1aaaaf2de1d8"
	issueId := "ed17dbab-e79b-4e87-8914-2d93ed55324b"
	desription := "updated description"

	result, err := clients.IssueClient.UpsertIssue(context.Background(), &issuepb.UpsertIssueGrpcRequest{
		Tenant:         tenant,
		Id:             issueId,
		LoggedInUserId: userId,
		Description:    desription,
		SourceFields: &commonpb.SourceFields{
			AppSource: appSource,
		},
		ExternalSystemFields: &commonpb.ExternalSystemFields{
			ExternalSystemId: "hubspot",
			ExternalId:       "456",
		},
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	if issueId != result.Id {
		log.Fatalf("Result is not expected")
	}
	log.Printf("Updated issue id: %v", result.Id)
}

func testCreateComment() {

	userId := "0fe25c46-bdac-485d-a5d5-a4a0390778ad"
	content := "hellow world"
	contentType := "text/plain"
	issueId := "ed17dbab-e79b-4e87-8914-2d93ed55324b"

	result, err := clients.CommentClient.UpsertComment(context.Background(), &commentpb.UpsertCommentGrpcRequest{
		Tenant:           tenant,
		Content:          content,
		ContentType:      contentType,
		AuthorUserId:     utils.StringPtr(userId),
		CommentedIssueId: utils.StringPtr(issueId),
		SourceFields: &commonpb.SourceFields{
			AppSource: appSource,
		},
		ExternalSystemFields: &commonpb.ExternalSystemFields{
			ExternalSystemId: "hubspot",
			ExternalId:       "123",
		},
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Created comment id: %v", result.Id)
}

func testCreateContract() {
	userId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	organizationId := "0f4114c4-c010-4303-a42a-460cf66ac598"
	now := utils.Now()
	aYearAgo := now.AddDate(-1, 0, 0)

	result, err := clients.ContractClient.CreateContract(context.Background(), &contractpb.CreateContractGrpcRequest{
		Tenant:           tenant,
		OrganizationId:   organizationId,
		LoggedInUserId:   userId,
		LengthInMonths:   1,
		ServiceStartedAt: utils.ConvertTimeToTimestampPtr(&aYearAgo),
		Name:             "year ago contract",
		AutoRenew:        false,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testUpdateContract() {

	userId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	contractId := "c5486341-c7d8-47eb-b75a-4016b8e3d6d5"
	in10Days := utils.Now().AddDate(0, 0, 10)
	yesterday := utils.Now().AddDate(0, 0, -1)
	in1year := utils.Now().AddDate(1, 0, 0)

	result, err := clients.ContractClient.UpdateContract(context.Background(), &contractpb.UpdateContractGrpcRequest{
		Tenant:           tenant,
		LoggedInUserId:   userId,
		Id:               contractId,
		Name:             "Saturday contract 1",
		SignedAt:         utils.ConvertTimeToTimestampPtr(&yesterday),
		ServiceStartedAt: utils.ConvertTimeToTimestampPtr(&in10Days),
		EndedAt:          utils.ConvertTimeToTimestampPtr(&in1year),
		SourceFields: &commonpb.SourceFields{
			AppSource: "test_app",
		},
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testAddContractService() {
	userId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	contractId := "769d1fb8-50a1-44bc-aff0-0f4338bd8ff2"
	price := 0.004
	billed := commonpb.BilledType_ONCE_BILLED

	result, err := clients.ServiceLineItemClient.CreateServiceLineItem(context.Background(), &servicelineitempb.CreateServiceLineItemGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: userId,
		Name:           "Custom",
		ContractId:     contractId,
		Price:          price,
		//Quantity:       int64(quantity),
		Billed: billed,
		SourceFields: &commonpb.SourceFields{
			AppSource: "test_app",
		},
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testCloseLooseOpportunity() {

	userId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	opportunityId := "d8305351-8568-4d97-9fe9-c6cf701636d0"

	result, err := clients.OpportunityClient.CloseLooseOpportunity(context.Background(), &opportunitypb.CloseLooseOpportunityGrpcRequest{
		Tenant:         tenant,
		Id:             opportunityId,
		LoggedInUserId: userId,
		AppSource:      appSource,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testCreateRenewalOpportunity() {

	result, err := clients.OpportunityClient.CreateRenewalOpportunity(context.Background(), &opportunitypb.CreateRenewalOpportunityGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: userId,
		ContractId:     contractId,
		SourceFields: &commonpb.SourceFields{
			AppSource: appSource,
		},
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testUpdateOnboardingStatus() {

	userId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	orgId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"

	result, err := clients.OrganizationClient.UpdateOnboardingStatus(context.Background(), &organizationpb.UpdateOnboardingStatusGrpcRequest{
		Tenant:           tenant,
		OrganizationId:   orgId,
		LoggedInUserId:   userId,
		Comments:         "test comments",
		AppSource:        appSource,
		OnboardingStatus: organizationpb.OnboardingStatus_ONBOARDING_STATUS_DONE,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testUpdateOrgOwner() {

	userId := "f7634527-ccda-4cbb-80d8-cc4af9124ef5"
	actorId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	orgId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"

	result, err := clients.OrganizationClient.UpdateOrganizationOwner(context.Background(), &organizationpb.UpdateOrganizationOwnerGrpcRequest{
		Tenant:         tenant,
		OrganizationId: orgId,
		LoggedInUserId: actorId,
		OwnerUserId:    userId,
		AppSource:      appSource,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testRefreshLastTouchpoint() {
	result, err := clients.OrganizationClient.RefreshLastTouchpoint(context.Background(), &organizationpb.OrganizationIdGrpcRequest{
		Tenant:         tenant,
		OrganizationId: orgId,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testRefreshRenewalSummary() {
	result, err := clients.OrganizationClient.RefreshRenewalSummary(context.Background(), &organizationpb.RefreshRenewalSummaryGrpcRequest{
		Tenant:         tenant,
		OrganizationId: orgId,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testTenantSettingsUpdate() {
	_, err := clients.TenantClient.UpdateTenantSettings(context.Background(), &tenantpb.UpdateTenantSettingsRequest{
		Tenant:               tenant,
		LogoRepositoryFileId: "123-abc",
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
}

func testAddBankAccount() {
	_, err := clients.TenantClient.AddBankAccount(context.Background(), &tenantpb.AddBankAccountGrpcRequest{
		Tenant:        tenant,
		OtherDetails:  "Some eur details",
		BankName:      "Bank of Europe",
		AccountNumber: "ACC-456",
		RoutingNumber: "ROUT-456",
		Iban:          "IBAN-456",
		Bic:           "BIC-456",
		SortCode:      "SORT-456",
		Currency:      "USD",
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
}
