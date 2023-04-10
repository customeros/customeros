package grpc_client

import (
	events_processing_contact "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/proto/contact"
	events_processing_email "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/proto/email"
	events_processing_organization "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/proto/organization"
	events_processing_phone_number "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/proto/phone_number"
	"google.golang.org/grpc"
)

type Clients struct {
	ContactClient      events_processing_contact.ContactGrpcServiceClient
	OrganizationClient events_processing_organization.OrganizationGrpcServiceClient
	PhoneNumberClient  events_processing_phone_number.PhoneNumberGrpcServiceClient
	EmailClient        events_processing_email.EmailGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		ContactClient:      events_processing_contact.NewContactGrpcServiceClient(conn),
		OrganizationClient: events_processing_organization.NewOrganizationGrpcServiceClient(conn),
		PhoneNumberClient:  events_processing_phone_number.NewPhoneNumberGrpcServiceClient(conn),
		EmailClient:        events_processing_email.NewEmailGrpcServiceClient(conn),
	}
	return &clients
}
