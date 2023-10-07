package grpc_client

import (
	contactgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	emailgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	jobrolegrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/job_role"
	locationgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	logentrygrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	orggrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	phonenumbergrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	usergrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"google.golang.org/grpc"
)

type Clients struct {
	ContactClient      contactgrpc.ContactGrpcServiceClient
	OrganizationClient orggrpc.OrganizationGrpcServiceClient
	PhoneNumberClient  phonenumbergrpc.PhoneNumberGrpcServiceClient
	EmailClient        emailgrpc.EmailGrpcServiceClient
	UserClient         usergrpc.UserGrpcServiceClient
	JobRoleClient      jobrolegrpc.JobRoleGrpcServiceClient
	LogEntryClient     logentrygrpc.LogEntryGrpcServiceClient
	LocationClient     locationgrpc.LocationGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		ContactClient:      contactgrpc.NewContactGrpcServiceClient(conn),
		OrganizationClient: orggrpc.NewOrganizationGrpcServiceClient(conn),
		PhoneNumberClient:  phonenumbergrpc.NewPhoneNumberGrpcServiceClient(conn),
		EmailClient:        emailgrpc.NewEmailGrpcServiceClient(conn),
		UserClient:         usergrpc.NewUserGrpcServiceClient(conn),
		JobRoleClient:      jobrolegrpc.NewJobRoleGrpcServiceClient(conn),
		LogEntryClient:     logentrygrpc.NewLogEntryGrpcServiceClient(conn),
		LocationClient:     locationgrpc.NewLocationGrpcServiceClient(conn),
	}
	return &clients
}
