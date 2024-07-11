package grpc_client

import (
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/comment"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	contract_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	email_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	interactioneventpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_event"
	interactionsessionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_session"
	invoice_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/issue"
	job_role_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/job_role"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	log_entry_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/log_entry"
	master_plan_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	offeringpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/offering"
	opportunity_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	orderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/order"
	organization_plan_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/org_plan"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	phone_number_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	reminder_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	service_line_item_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	tenant_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"google.golang.org/grpc"
)

type Clients struct {
	ContactClient            contactpb.ContactGrpcServiceClient
	ContractClient           contract_grpc_service.ContractGrpcServiceClient
	EmailClient              email_grpc_service.EmailGrpcServiceClient
	InvoiceClient            invoice_grpc_service.InvoiceGrpcServiceClient
	JobRoleClient            job_role_grpc_service.JobRoleGrpcServiceClient
	LogEntryClient           log_entry_grpc_service.LogEntryGrpcServiceClient
	MasterPlanClient         master_plan_grpc_service.MasterPlanGrpcServiceClient
	OpportunityClient        opportunity_grpc_service.OpportunityGrpcServiceClient
	OrganizationClient       organization_grpc_service.OrganizationGrpcServiceClient
	OrganizationPlanClient   organization_plan_grpc_service.OrganizationPlanGrpcServiceClient
	PhoneNumberClient        phone_number_grpc_service.PhoneNumberGrpcServiceClient
	ReminderClient           reminder_grpc_service.ReminderGrpcServiceClient
	ServiceLineItemClient    service_line_item_grpc_service.ServiceLineItemGrpcServiceClient
	TenantClient             tenant_grpc_service.TenantGrpcServiceClient
	UserClient               userpb.UserGrpcServiceClient
	OfferingClient           offeringpb.OfferingGrpcServiceClient
	LocationClient           locationpb.LocationGrpcServiceClient
	IssueClient              issuepb.IssueGrpcServiceClient
	InteractionEventClient   interactioneventpb.InteractionEventGrpcServiceClient
	InteractionSessionClient interactionsessionpb.InteractionSessionGrpcServiceClient
	CommentClient            commentpb.CommentGrpcServiceClient
	OrderClient              orderpb.OrderGrpcServiceClient
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
		EmailClient:              email_grpc_service.NewEmailGrpcServiceClient(conn),
		UserClient:               userpb.NewUserGrpcServiceClient(conn),
		JobRoleClient:            job_role_grpc_service.NewJobRoleGrpcServiceClient(conn),
		LogEntryClient:           log_entry_grpc_service.NewLogEntryGrpcServiceClient(conn),
		ContractClient:           contract_grpc_service.NewContractGrpcServiceClient(conn),
		ServiceLineItemClient:    service_line_item_grpc_service.NewServiceLineItemGrpcServiceClient(conn),
		OpportunityClient:        opportunity_grpc_service.NewOpportunityGrpcServiceClient(conn),
		MasterPlanClient:         master_plan_grpc_service.NewMasterPlanGrpcServiceClient(conn),
		InvoiceClient:            invoice_grpc_service.NewInvoiceGrpcServiceClient(conn),
		OrganizationPlanClient:   organization_plan_grpc_service.NewOrganizationPlanGrpcServiceClient(conn),
		TenantClient:             tenant_grpc_service.NewTenantGrpcServiceClient(conn),
		ReminderClient:           reminder_grpc_service.NewReminderGrpcServiceClient(conn),
		OfferingClient:           offeringpb.NewOfferingGrpcServiceClient(conn),
		LocationClient:           locationpb.NewLocationGrpcServiceClient(conn),
		IssueClient:              issuepb.NewIssueGrpcServiceClient(conn),
		InteractionEventClient:   interactioneventpb.NewInteractionEventGrpcServiceClient(conn),
		InteractionSessionClient: interactionsessionpb.NewInteractionSessionGrpcServiceClient(conn),
		CommentClient:            commentpb.NewCommentGrpcServiceClient(conn),
		OrderClient:              orderpb.NewOrderGrpcServiceClient(conn),
		EventStoreClient:         eventstorepb.NewEventStoreGrpcServiceClient(conn),
	}
	return &clients
}
