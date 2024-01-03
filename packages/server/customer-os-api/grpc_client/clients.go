package grpc_client

import (
	contact_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	contract_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	email_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	job_role_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/job_role"
	log_entry_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/log_entry"
	master_plan_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	opportunity_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	phone_number_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	service_line_item_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	user_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"google.golang.org/grpc"
)

type Clients struct {
	ContactClient         contact_grpc_service.ContactGrpcServiceClient
	OrganizationClient    organization_grpc_service.OrganizationGrpcServiceClient
	PhoneNumberClient     phone_number_grpc_service.PhoneNumberGrpcServiceClient
	EmailClient           email_grpc_service.EmailGrpcServiceClient
	UserClient            user_grpc_service.UserGrpcServiceClient
	JobRoleClient         job_role_grpc_service.JobRoleGrpcServiceClient
	LogEntryClient        log_entry_grpc_service.LogEntryGrpcServiceClient
	ContractClient        contract_grpc_service.ContractGrpcServiceClient
	ServiceLineItemClient service_line_item_grpc_service.ServiceLineItemGrpcServiceClient
	OpportunityClient     opportunity_grpc_service.OpportunityGrpcServiceClient
	MasterPlanClient      master_plan_grpc_service.MasterPlanGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		ContactClient:         contact_grpc_service.NewContactGrpcServiceClient(conn),
		OrganizationClient:    organization_grpc_service.NewOrganizationGrpcServiceClient(conn),
		PhoneNumberClient:     phone_number_grpc_service.NewPhoneNumberGrpcServiceClient(conn),
		EmailClient:           email_grpc_service.NewEmailGrpcServiceClient(conn),
		UserClient:            user_grpc_service.NewUserGrpcServiceClient(conn),
		JobRoleClient:         job_role_grpc_service.NewJobRoleGrpcServiceClient(conn),
		LogEntryClient:        log_entry_grpc_service.NewLogEntryGrpcServiceClient(conn),
		ContractClient:        contract_grpc_service.NewContractGrpcServiceClient(conn),
		ServiceLineItemClient: service_line_item_grpc_service.NewServiceLineItemGrpcServiceClient(conn),
		OpportunityClient:     opportunity_grpc_service.NewOpportunityGrpcServiceClient(conn),
		MasterPlanClient:      master_plan_grpc_service.NewMasterPlanGrpcServiceClient(conn),
	}
	return &clients
}
