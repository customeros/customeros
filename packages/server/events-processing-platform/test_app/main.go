package main

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client/interceptor"
	interaction_event_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"google.golang.org/grpc"
)

const grpcApiKey = "082c1193-a5a2-42fc-87fc-e960e692fffd"

type Clients struct {
	InteractionEventClient interaction_event_grpc_service.InteractionEventGrpcServiceClient
	OrganizationClient     organization_grpc_service.OrganizationGrpcServiceClient
}

var clients *Clients

func main() {
	InitClients()
	//testRequestGenerateSummaryRequest()
	//testRequestGenerateActionItemsRequest()
	//testCreateOrganization()
	//testUpdateOrganization()
	//testHideOrganization()
	testShowOrganization()
}

func InitClients() {
	conn, _ := grpc.Dial("localhost:5001", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			interceptor.ApiKeyEnricher(grpcApiKey),
		))
	clients = &Clients{
		InteractionEventClient: interaction_event_grpc_service.NewInteractionEventGrpcServiceClient(conn),
		OrganizationClient:     organization_grpc_service.NewOrganizationGrpcServiceClient(conn),
	}
}

func testRequestGenerateSummaryRequest() {
	tenant := "openline"
	interactionEventId := "555263fe-2e39-48f0-a8c2-c4c7a5ffb23d"

	result, _ := clients.InteractionEventClient.RequestGenerateSummary(context.TODO(), &interaction_event_grpc_service.RequestGenerateSummaryGrpcRequest{
		Tenant:             tenant,
		InteractionEventId: interactionEventId,
	})
	print(result)
}

func testRequestGenerateActionItemsRequest() {
	tenant := "openline"
	interactionEventId := "555263fe-2e39-48f0-a8c2-c4c7a5ffb23d"

	result, _ := clients.InteractionEventClient.RequestGenerateActionItems(context.TODO(), &interaction_event_grpc_service.RequestGenerateActionItensGrpcRequest{
		Tenant:             tenant,
		InteractionEventId: interactionEventId,
	})
	print(result)
}

func testCreateOrganization() {
	tenant := "openline"
	userId := "697563a8-171c-4950-a067-1aaaaf2de1d8"
	organizationId := "ccc"
	website := ""

	result, _ := clients.OrganizationClient.UpsertOrganization(context.TODO(), &organization_grpc_service.UpsertOrganizationGrpcRequest{
		Tenant:  tenant,
		Id:      organizationId,
		Website: website,
		UserId:  userId,
	})
	print(result)
}

func testUpdateOrganization() {
	tenant := "openline"
	organizationId := "ccc"
	website := ""

	result, _ := clients.OrganizationClient.UpsertOrganization(context.TODO(), &organization_grpc_service.UpsertOrganizationGrpcRequest{
		Tenant:  tenant,
		Id:      organizationId,
		Website: website,
	})
	print(result)
}

func testHideOrganization() {
	tenant := "openline"
	organizationId := "ccc"

	result, _ := clients.OrganizationClient.HideOrganization(context.TODO(), &organization_grpc_service.OrganizationIdGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
	})
	print(result)
}

func testShowOrganization() {
	tenant := "openline"
	organizationId := "ccc"

	result, _ := clients.OrganizationClient.ShowOrganization(context.TODO(), &organization_grpc_service.OrganizationIdGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organizationId,
	})
	print(result)
}
