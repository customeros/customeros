package grpc_client

import (
	contact_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	email_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	location_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	log_entry_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	phone_number_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	user_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"google.golang.org/grpc"
)

type Clients struct {
	ContactClient      contact_grpc_service.ContactGrpcServiceClient
	OrganizationClient organization_grpc_service.OrganizationGrpcServiceClient
	PhoneNumberClient  phone_number_grpc_service.PhoneNumberGrpcServiceClient
	LocationClient     location_grpc_service.LocationGrpcServiceClient
	EmailClient        email_grpc_service.EmailGrpcServiceClient
	UserClient         user_grpc_service.UserGrpcServiceClient
	LogEntryClient     log_entry_grpc_service.LogEntryGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		ContactClient:      contact_grpc_service.NewContactGrpcServiceClient(conn),
		OrganizationClient: organization_grpc_service.NewOrganizationGrpcServiceClient(conn),
		PhoneNumberClient:  phone_number_grpc_service.NewPhoneNumberGrpcServiceClient(conn),
		EmailClient:        email_grpc_service.NewEmailGrpcServiceClient(conn),
		UserClient:         user_grpc_service.NewUserGrpcServiceClient(conn),
		LocationClient:     location_grpc_service.NewLocationGrpcServiceClient(conn),
		LogEntryClient:     log_entry_grpc_service.NewLogEntryGrpcServiceClient(conn),
	}
	return &clients
}
