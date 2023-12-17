package grpc_client

import (
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/comment"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	interactioneventpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/issue"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	logentrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	phonenumpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"google.golang.org/grpc"
)

type Clients struct {
	OrganizationClient     organizationpb.OrganizationGrpcServiceClient
	PhoneNumberClient      phonenumpb.PhoneNumberGrpcServiceClient
	EmailClient            emailpb.EmailGrpcServiceClient
	UserClient             userpb.UserGrpcServiceClient
	LogEntryClient         logentrypb.LogEntryGrpcServiceClient
	LocationClient         locationpb.LocationGrpcServiceClient
	IssueClient            issuepb.IssueGrpcServiceClient
	InteractionEventClient interactioneventpb.InteractionEventGrpcServiceClient
	CommentClient          commentpb.CommentGrpcServiceClient
	OpportunityClient      opportunitypb.OpportunityGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		OrganizationClient:     organizationpb.NewOrganizationGrpcServiceClient(conn),
		PhoneNumberClient:      phonenumpb.NewPhoneNumberGrpcServiceClient(conn),
		EmailClient:            emailpb.NewEmailGrpcServiceClient(conn),
		UserClient:             userpb.NewUserGrpcServiceClient(conn),
		LogEntryClient:         logentrypb.NewLogEntryGrpcServiceClient(conn),
		LocationClient:         locationpb.NewLocationGrpcServiceClient(conn),
		IssueClient:            issuepb.NewIssueGrpcServiceClient(conn),
		InteractionEventClient: interactioneventpb.NewInteractionEventGrpcServiceClient(conn),
		CommentClient:          commentpb.NewCommentGrpcServiceClient(conn),
		OpportunityClient:      opportunitypb.NewOpportunityGrpcServiceClient(conn),
	}
	return &clients
}
