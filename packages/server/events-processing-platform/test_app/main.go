package main

import (
	"context"
	eventstorepb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	"log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client/interceptor"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	commonpb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	organizationpb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
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

	//testRefreshRenewalSummary()
	//PleasePayInvoiceNotification()
	//testCreateInvoice()
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
