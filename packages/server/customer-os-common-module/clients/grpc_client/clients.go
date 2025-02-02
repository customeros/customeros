package grpc_client

import (
	invoice_grpc_service "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"google.golang.org/grpc"
)

type Clients struct {
	InvoiceClient invoice_grpc_service.InvoiceGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		InvoiceClient: invoice_grpc_service.NewInvoiceGrpcServiceClient(conn),
	}
	return &clients
}
