package grpc_client

import (
	eventstorepb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	invoice_grpc_service "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	organization_grpc_service "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"google.golang.org/grpc"
)

type Clients struct {
	InvoiceClient      invoice_grpc_service.InvoiceGrpcServiceClient
	OrganizationClient organization_grpc_service.OrganizationGrpcServiceClient
	EventStoreClient   eventstorepb.EventStoreGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		OrganizationClient: organization_grpc_service.NewOrganizationGrpcServiceClient(conn),
		InvoiceClient:      invoice_grpc_service.NewInvoiceGrpcServiceClient(conn),
		EventStoreClient:   eventstorepb.NewEventStoreGrpcServiceClient(conn),
	}
	return &clients
}
