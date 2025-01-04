package main

import (
	"context"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	"log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client/interceptor"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"google.golang.org/grpc"
)

const grpcApiKey = "082c1193-a5a2-42fc-87fc-e960e692fffd"
const appSource = "test_app"

var tenant = "customerosai"
var orgId = "ceae019f-d1e3-49b3-87c5-35ebb68a5ff1"

type Clients struct {
	OrganizationClient organizationpb.OrganizationGrpcServiceClient
	InvoiceClient      invoicepb.InvoiceGrpcServiceClient
	EventStoreClient   eventstorepb.EventStoreGrpcServiceClient
}

var clients *Clients

func InitClients() {
	conn, _ := grpc.Dial("localhost:5001", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			interceptor.ApiKeyEnricher(grpcApiKey),
		))
	clients = &Clients{
		OrganizationClient: organizationpb.NewOrganizationGrpcServiceClient(conn),
		InvoiceClient:      invoicepb.NewInvoiceGrpcServiceClient(conn),
		EventStoreClient:   eventstorepb.NewEventStoreGrpcServiceClient(conn),
	}
}

func main() {
	InitClients()

	//testHideOrganization()
	//testAddCustomField()
	//testRemoveParentOrganization()
	//testContactLinkWithPhoneNumber()
	//testOrganizationLinkWithEmail()
	//testOrganizationLinkWithPhoneNumber()
	//testCreateComment()
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
