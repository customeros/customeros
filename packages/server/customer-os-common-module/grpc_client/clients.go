package grpc_client

import (
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/comment"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	contract_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	interactioneventpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_event"
	interactionsessionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_session"
	invoice_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/issue"
	job_role_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/job_role"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	opportunity_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	phone_number_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	service_line_item_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	"google.golang.org/grpc"
)

type Clients struct {
	ContactClient            contactpb.ContactGrpcServiceClient
	ContractClient           contract_grpc_service.ContractGrpcServiceClient
	InvoiceClient            invoice_grpc_service.InvoiceGrpcServiceClient
	JobRoleClient            job_role_grpc_service.JobRoleGrpcServiceClient
	OpportunityClient        opportunity_grpc_service.OpportunityGrpcServiceClient
	OrganizationClient       organization_grpc_service.OrganizationGrpcServiceClient
	PhoneNumberClient        phone_number_grpc_service.PhoneNumberGrpcServiceClient
	ServiceLineItemClient    service_line_item_grpc_service.ServiceLineItemGrpcServiceClient
	LocationClient           locationpb.LocationGrpcServiceClient
	IssueClient              issuepb.IssueGrpcServiceClient
	InteractionEventClient   interactioneventpb.InteractionEventGrpcServiceClient
	InteractionSessionClient interactionsessionpb.InteractionSessionGrpcServiceClient
	CommentClient            commentpb.CommentGrpcServiceClient
	EventStoreClient         eventstorepb.EventStoreGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		ContactClient:            contactpb.NewContactGrpcServiceClient(conn),
		OrganizationClient:       organization_grpc_service.NewOrganizationGrpcServiceClient(conn),
		PhoneNumberClient:        phone_number_grpc_service.NewPhoneNumberGrpcServiceClient(conn),
		JobRoleClient:            job_role_grpc_service.NewJobRoleGrpcServiceClient(conn),
		ContractClient:           contract_grpc_service.NewContractGrpcServiceClient(conn),
		ServiceLineItemClient:    service_line_item_grpc_service.NewServiceLineItemGrpcServiceClient(conn),
		OpportunityClient:        opportunity_grpc_service.NewOpportunityGrpcServiceClient(conn),
		InvoiceClient:            invoice_grpc_service.NewInvoiceGrpcServiceClient(conn),
		LocationClient:           locationpb.NewLocationGrpcServiceClient(conn),
		IssueClient:              issuepb.NewIssueGrpcServiceClient(conn),
		InteractionEventClient:   interactioneventpb.NewInteractionEventGrpcServiceClient(conn),
		InteractionSessionClient: interactionsessionpb.NewInteractionSessionGrpcServiceClient(conn),
		CommentClient:            commentpb.NewCommentGrpcServiceClient(conn),
		EventStoreClient:         eventstorepb.NewEventStoreGrpcServiceClient(conn),
	}
	return &clients
}
