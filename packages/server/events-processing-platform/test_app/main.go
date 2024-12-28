package main

import (
	"context"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	"log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client/interceptor"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/comment"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	"google.golang.org/grpc"
)

const grpcApiKey = "082c1193-a5a2-42fc-87fc-e960e692fffd"
const appSource = "test_app"

var tenant = "customerosai"
var userId = "05f382ba-0fa9-4828-940c-efb4e2e6b84c"
var orgId = "ceae019f-d1e3-49b3-87c5-35ebb68a5ff1"
var contractId = "769d1fb8-50a1-44bc-aff0-0f4338bd8ff2"

type Clients struct {
	OrganizationClient    organizationpb.OrganizationGrpcServiceClient
	CommentClient         commentpb.CommentGrpcServiceClient
	ServiceLineItemClient servicelineitempb.ServiceLineItemGrpcServiceClient
	OpportunityClient     opportunitypb.OpportunityGrpcServiceClient
	InvoiceClient         invoicepb.InvoiceGrpcServiceClient
	EventStoreClient      eventstorepb.EventStoreGrpcServiceClient
}

var clients *Clients

func InitClients() {
	conn, _ := grpc.Dial("localhost:5001", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			interceptor.ApiKeyEnricher(grpcApiKey),
		))
	clients = &Clients{
		OrganizationClient:    organizationpb.NewOrganizationGrpcServiceClient(conn),
		CommentClient:         commentpb.NewCommentGrpcServiceClient(conn),
		OpportunityClient:     opportunitypb.NewOpportunityGrpcServiceClient(conn),
		ServiceLineItemClient: servicelineitempb.NewServiceLineItemGrpcServiceClient(conn),
		InvoiceClient:         invoicepb.NewInvoiceGrpcServiceClient(conn),
		EventStoreClient:      eventstorepb.NewEventStoreGrpcServiceClient(conn),
	}
}

func main() {
	InitClients()

	//testHideOrganization()
	//testAddCustomField()
	//testCreatePhoneNumber()
	//testAddParentOrganization()
	//testRemoveParentOrganization()
	//testContactLinkWithPhoneNumber()
	//testContactLinkWithLocation()
	//testOrganizationLinkWithEmail()
	//testOrganizationLinkWithPhoneNumber()
	//testOrganizationLinkWithLocation()
	//testCreateComment()
	//testAddContractService()
	//testCloseLooseOpportunity()
	//testCreateRenewalOpportunity()
	//testArchiveOpportunity()
	//testUpdateOnboardingStatus()
	//testUpdateOrgOwner()
	//testRefreshRenewalSummary()
	//testAddTenantBillingProfile()
	//PleasePayInvoiceNotification()
	//testCreateInvoice()
	//testCreateReminder()
	//testUpdateReminder()
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

func PleasePayInvoiceNotification() {
	_, err := clients.InvoiceClient.PayInvoiceNotification(context.Background(), &invoicepb.PayInvoiceNotificationRequest{
		Tenant:    tenant,
		InvoiceId: "e3af66b0-8e74-4aa7-941d-4b87518d7131",
	})

	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
	}
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
