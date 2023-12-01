package main

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client/interceptor"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/comment"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	iepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/issue"
	logentrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/service_line_item"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

const grpcApiKey = "082c1193-a5a2-42fc-87fc-e960e692fffd"
const appSource = "test_app"

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
	}
}

func main() {
	InitClients()
	//testRequestGenerateSummaryRequest()
	//testRequestGenerateActionItemsRequest()
	//testCreateOrganization()
	//testUpdateOrganization()
	//testHideOrganization()
	//testShowOrganization()
	//testCreateLogEntry()
	//testUpdateLogEntry()
	//testAddCustomField()
	//testCreateEmail()
	//testCreatePhoneNumber()
	//testAddParentOrganization()
	//testRemoveParentOrganization()
	//testCreateContact()
	//testUpdateContact()
	//testContactLinkWithEmail()
	//testContactLinkWithPhoneNumber()
	//testContactLinkWithLocation()
	//testOrganizationLinkWithEmail()
	//testOrganizationLinkWithPhoneNumber()
	//testOrganizationLinkWithLocation()
	//testContactLinkWithOrganization()
	//testCreateIssue()
	//testUpdateIssue()
	//testCreateComment()
	//testUserLinkWithEmail()
	//testCreateContract()
	//testUpdateContract()
	testAddContractService()
}

func testRequestGenerateSummaryRequest() {
	tenant := "openline"
	interactionEventId := "555263fe-2e39-48f0-a8c2-c4c7a5ffb23d"

	result, _ := clients.InteractionEventClient.RequestGenerateSummary(context.Background(), &iepb.RequestGenerateSummaryGrpcRequest{
		Tenant:             tenant,
		InteractionEventId: interactionEventId,
	})
	print(result)
}

func testRequestGenerateActionItemsRequest() {
	tenant := "openline"
	interactionEventId := "555263fe-2e39-48f0-a8c2-c4c7a5ffb23d"

	result, _ := clients.InteractionEventClient.RequestGenerateActionItems(context.Background(), &iepb.RequestGenerateActionItemsGrpcRequest{
		Tenant:             tenant,
		InteractionEventId: interactionEventId,
	})
	print(result)
}

func testCreateOrganization() {
	tenant := "openline"
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

func testUpdateOrganization() {
	tenant := "openline"
	organizationId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"
	website := "xtz.com"
	lastFoundingAmont := "1Million"
	partial := true

	result, _ := clients.OrganizationClient.UpsertOrganization(context.Background(), &organizationpb.UpsertOrganizationGrpcRequest{
		Tenant:            tenant,
		Id:                organizationId,
		Website:           website,
		LastFundingAmount: lastFoundingAmont,
		IgnoreEmptyFields: partial,
	})
	print(result)
}

func testHideOrganization() {
	tenant := "openline"
	organizationId := "ccc"

	result, _ := clients.OrganizationClient.HideOrganization(context.Background(), &organizationpb.OrganizationIdGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
	})
	print(result)
}

func testShowOrganization() {
	tenant := "openline"
	organizationId := "ccc"

	result, _ := clients.OrganizationClient.ShowOrganization(context.Background(), &organizationpb.OrganizationIdGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
	})
	print(result)
}

func testCreateLogEntry() {
	tenant := "openline"
	organizationId := "5e72b6fb-5f20-4973-9b96-52f4543a0df3"
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
	tenant := "openline"
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
	tenant := "openline"
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

func testCreateEmail() {
	tenant := "openline"
	userId := "697563a8-171c-4950-a067-1aaaaf2de1d8"
	rawEmail := "aa@test.com"

	result, _ := clients.EmailClient.UpsertEmail(context.Background(), &emailpb.UpsertEmailGrpcRequest{
		Tenant:         tenant,
		RawEmail:       rawEmail,
		LoggedInUserId: userId,
		SourceFields: &commonpb.SourceFields{
			AppSource: "test_app",
		},
	})
	print(result)
}

func testCreatePhoneNumber() {
	tenant := "openline"
	userId := "697563a8-171c-4950-a067-1aaaaf2de1d8"
	rawPhoneNumber := "+1234"

	result, _ := clients.PhoneNumberClient.UpsertPhoneNumber(context.Background(), &phonenumberpb.UpsertPhoneNumberGrpcRequest{
		Tenant:         tenant,
		PhoneNumber:    rawPhoneNumber,
		LoggedInUserId: userId,
	})
	print(result)
}

func testAddParentOrganization() {
	tenant := "openline"
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
	tenant := "openline"
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

func testCreateContact() {
	tenant := "openline"
	userId := "697563a8-171c-4950-a067-1aaaaf2de1d8"
	name := "hubspot contact 3"

	result, _ := clients.ContactClient.UpsertContact(context.Background(), &contactpb.UpsertContactGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: userId,
		Name:           name,
		ExternalSystemFields: &commonpb.ExternalSystemFields{
			ExternalSystemId: "hubspot",
			ExternalId:       "123",
		},
	})
	print(result)
}

func testUpdateContact() {
	tenant := "openline"
	contactId := "dd7bd45e-d6d3-405c-a7ba-cd4127479c20"
	name := "hubspot contact 4"

	result, _ := clients.ContactClient.UpsertContact(context.Background(), &contactpb.UpsertContactGrpcRequest{
		Tenant: tenant,
		Name:   name,
		Id:     contactId,
		ExternalSystemFields: &commonpb.ExternalSystemFields{
			ExternalSystemId: "hubspot",
			ExternalId:       "ABC",
		},
	})
	print(result)
}

func testContactLinkWithLocation() {
	tenant := "openline"
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
	tenant := "openline"
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

func testContactLinkWithEmail() {
	tenant := "openline"
	contactId := "dd7bd45e-d6d3-405c-a7ba-cd4127479c20"
	emailId := "548a69d2-90fe-439d-b5bb-ee7b68e17d34"

	result, _ := clients.ContactClient.LinkEmailToContact(context.Background(), &contactpb.LinkEmailToContactGrpcRequest{
		Tenant:    tenant,
		ContactId: contactId,
		EmailId:   emailId,
		Primary:   true,
		Label:     "work",
		AppSource: appSource,
	})
	print(result)
}

func testOrganizationLinkWithLocation() {
	tenant := "openline"
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
	tenant := "openline"
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

func testOrganizationLinkWithEmail() {
	tenant := "openline"
	orgId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"
	emailId := "548a69d2-90fe-439d-b5bb-ee7b68e17d34"

	result, _ := clients.OrganizationClient.LinkEmailToOrganization(context.Background(), &organizationpb.LinkEmailToOrganizationGrpcRequest{
		Tenant:         tenant,
		OrganizationId: orgId,
		EmailId:        emailId,
		Primary:        true,
		Label:          "work",
	})
	print(result)
}

func testContactLinkWithOrganization() {
	tenant := "openline"
	contactId := "dd7bd45e-d6d3-405c-a7ba-cd4127479c20"
	orgId := "cfaaf31f-ec3b-44d1-836e-4e50834632ae"
	jobRole := "CTO"

	result, _ := clients.ContactClient.LinkWithOrganization(context.Background(), &contactpb.LinkWithOrganizationGrpcRequest{
		Tenant:         tenant,
		ContactId:      contactId,
		OrganizationId: orgId,
		JobTitle:       jobRole,
		Primary:        true,
		Description:    "CEO of the company",
		StartedAt:      timestamppb.Now(),
	})
	print(result)
}

func testCreateIssue() {
	tenant := "openline"
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
	tenant := "openline"
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
	tenant := "openline"
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

func testUserLinkWithEmail() {
	tenant := "openline"
	userId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	emailId := "548a69d2-90fe-439d-b5bb-ee7b68e17d34"
	appSource := "integration.app"
	loggedInUserid := ""

	result, err := clients.UserClient.LinkEmailToUser(context.Background(), &userpb.LinkEmailToUserGrpcRequest{
		Tenant:         tenant,
		UserId:         userId,
		LoggedInUserId: loggedInUserid,
		EmailId:        emailId,
		Primary:        true,
		Label:          "work",
		AppSource:      appSource,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testCreateContract() {
	tenant := "openline"
	userId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	organizationId := "a00dc1f8-aec9-4107-bb77-48ef39f897bc"

	result, err := clients.ContractClient.CreateContract(context.Background(), &contractpb.CreateContractGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
		LoggedInUserId: userId,
		Name:           "Saturday contract 2",
	})
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
	log.Printf("Result: %v", result.Id)
}

func testUpdateContract() {
	tenant := "openline"
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
	tenant := "openline"
	userId := "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
	contractId := "c5486341-c7d8-47eb-b75a-4016b8e3d6d5"
	price := 0.004
	//quantity := 3
	//billed := servicelineitempb.BilledType_ANNUALLY_BILLED
	billed := servicelineitempb.BilledType_ONCE_BILLED

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
