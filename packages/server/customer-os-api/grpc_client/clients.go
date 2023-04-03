package grpc_client

import (
	events_processing_contact "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/proto/contact"
	events_processing_phone_number "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/proto/phone_number"
	"google.golang.org/grpc"
)

type Clients struct {
	ContactClient     events_processing_contact.ContactGrpcServiceClient
	PhoneNumberClient events_processing_phone_number.PhoneNumberGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		ContactClient:     events_processing_contact.NewContactGrpcServiceClient(conn),
		PhoneNumberClient: events_processing_phone_number.NewPhoneNumberGrpcServiceClient(conn),
	}

	return &clients
}
