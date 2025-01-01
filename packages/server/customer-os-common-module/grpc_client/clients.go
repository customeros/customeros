package grpc_client

import (
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	interactionsessionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_session"
	invoice_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	opportunity_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"google.golang.org/grpc"
)

type Clients struct {
	InvoiceClient            invoice_grpc_service.InvoiceGrpcServiceClient
	OpportunityClient        opportunity_grpc_service.OpportunityGrpcServiceClient
	OrganizationClient       organization_grpc_service.OrganizationGrpcServiceClient
	LocationClient           locationpb.LocationGrpcServiceClient
	InteractionSessionClient interactionsessionpb.InteractionSessionGrpcServiceClient
	EventStoreClient         eventstorepb.EventStoreGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		OrganizationClient:       organization_grpc_service.NewOrganizationGrpcServiceClient(conn),
		OpportunityClient:        opportunity_grpc_service.NewOpportunityGrpcServiceClient(conn),
		InvoiceClient:            invoice_grpc_service.NewInvoiceGrpcServiceClient(conn),
		LocationClient:           locationpb.NewLocationGrpcServiceClient(conn),
		InteractionSessionClient: interactionsessionpb.NewInteractionSessionGrpcServiceClient(conn),
		EventStoreClient:         eventstorepb.NewEventStoreGrpcServiceClient(conn),
	}
	return &clients
}
