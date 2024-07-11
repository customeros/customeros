package server

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events"
	orderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/order"
	"net"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/interceptors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/service"
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/comment"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	countrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/country"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	iepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_event"
	ispb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_session"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/issue"
	jobrolepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/job_role"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	logentrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/log_entry"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	orgplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/org_plan"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	phonenumpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (server *Server) NewEventProcessorGrpcServer() (func() error, *grpc.Server, error) {
	l, err := net.Listen(events.Tcp, server.Config.GRPC.Port)
	if err != nil {
		return nil, nil, errors.Wrap(err, "net.Listen")
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
			interceptors.CheckApiKeyInterceptor(server.Config.GRPC.ApiKey),
		),
		),
	)

	RegisterGrpcServices(grpcServer, server.Services)

	go func() {
		server.Log.Infof("%s gRPC Server is listening on port: {%s}", GetMicroserviceName(server.Config), server.Config.GRPC.Port)
		server.Log.Error(grpcServer.Serve(l))
	}()

	return l.Close, grpcServer, nil
}

func RegisterGrpcServices(grpcServer *grpc.Server, services *service.Services) {
	contactpb.RegisterContactGrpcServiceServer(grpcServer, services.ContactService)
	organizationpb.RegisterOrganizationGrpcServiceServer(grpcServer, services.OrganizationService)
	phonenumpb.RegisterPhoneNumberGrpcServiceServer(grpcServer, services.PhoneNumberService)
	emailpb.RegisterEmailGrpcServiceServer(grpcServer, services.EmailService)
	userpb.RegisterUserGrpcServiceServer(grpcServer, services.UserService)
	locationpb.RegisterLocationGrpcServiceServer(grpcServer, services.LocationService)
	jobrolepb.RegisterJobRoleGrpcServiceServer(grpcServer, services.JobRoleService)
	iepb.RegisterInteractionEventGrpcServiceServer(grpcServer, services.InteractionEventService)
	ispb.RegisterInteractionSessionGrpcServiceServer(grpcServer, services.InteractionSessionService)
	logentrypb.RegisterLogEntryGrpcServiceServer(grpcServer, services.LogEntryService)
	issuepb.RegisterIssueGrpcServiceServer(grpcServer, services.IssueService)
	commentpb.RegisterCommentGrpcServiceServer(grpcServer, services.CommentService)
	opportunitypb.RegisterOpportunityGrpcServiceServer(grpcServer, services.OpportunityService)
	contractpb.RegisterContractGrpcServiceServer(grpcServer, services.ContractService)
	servicelineitempb.RegisterServiceLineItemGrpcServiceServer(grpcServer, services.ServiceLineItemService)
	masterplanpb.RegisterMasterPlanGrpcServiceServer(grpcServer, services.MasterPlanService)
	invoicepb.RegisterInvoiceGrpcServiceServer(grpcServer, services.InvoiceService)
	countrypb.RegisterCountryGrpcServiceServer(grpcServer, services.CountryService)
	tenantpb.RegisterTenantGrpcServiceServer(grpcServer, services.TenantService)
	orgplanpb.RegisterOrganizationPlanGrpcServiceServer(grpcServer, services.OrganizationPlanService)
	reminderpb.RegisterReminderGrpcServiceServer(grpcServer, services.ReminderService)
	orderpb.RegisterOrderGrpcServiceServer(grpcServer, services.OrderService)
	eventstorepb.RegisterEventStoreGrpcServiceServer(grpcServer, services.EventStoreService)
}
