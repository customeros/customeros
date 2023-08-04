package main

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client/interceptor"
	interaction_event_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	"google.golang.org/grpc"
)

const grpcApiKey = "082c1193-a5a2-42fc-87fc-e960e692fffd"

type Clients struct {
	InteractionEventClient interaction_event_grpc_service.InteractionEventGrpcServiceClient
}

var clients *Clients

func main() {
	InitClients()
	testRequestGenerateSummaryRequest()
}

func InitClients() {
	conn, _ := grpc.Dial("localhost:5001", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			interceptor.ApiKeyEnricher(grpcApiKey),
		))
	clients = &Clients{
		InteractionEventClient: interaction_event_grpc_service.NewInteractionEventGrpcServiceClient(conn),
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
