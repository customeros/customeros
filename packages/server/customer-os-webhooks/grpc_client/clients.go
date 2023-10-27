package grpc_client

import (
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	interactioneventpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/issue"
	jobrolepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/job_role"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	logentrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	phonenumpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"google.golang.org/grpc"
)

type Clients struct {
	ContactClient          contactpb.ContactGrpcServiceClient
	OrganizationClient     organizationpb.OrganizationGrpcServiceClient
	PhoneNumberClient      phonenumpb.PhoneNumberGrpcServiceClient
	EmailClient            emailpb.EmailGrpcServiceClient
	UserClient             userpb.UserGrpcServiceClient
	JobRoleClient          jobrolepb.JobRoleGrpcServiceClient
	LogEntryClient         logentrypb.LogEntryGrpcServiceClient
	LocationClient         locationpb.LocationGrpcServiceClient
	IssueClient            issuepb.IssueGrpcServiceClient
	InteractionEventClient interactioneventpb.InteractionEventGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		ContactClient:          contactpb.NewContactGrpcServiceClient(conn),
		OrganizationClient:     organizationpb.NewOrganizationGrpcServiceClient(conn),
		PhoneNumberClient:      phonenumpb.NewPhoneNumberGrpcServiceClient(conn),
		EmailClient:            emailpb.NewEmailGrpcServiceClient(conn),
		UserClient:             userpb.NewUserGrpcServiceClient(conn),
		JobRoleClient:          jobrolepb.NewJobRoleGrpcServiceClient(conn),
		LogEntryClient:         logentrypb.NewLogEntryGrpcServiceClient(conn),
		LocationClient:         locationpb.NewLocationGrpcServiceClient(conn),
		IssueClient:            issuepb.NewIssueGrpcServiceClient(conn),
		InteractionEventClient: interactioneventpb.NewInteractionEventGrpcServiceClient(conn),
	}
	return &clients
}
