package server

import (
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/comment"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	iepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/issue"
	jobrolepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/job_role"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	logentrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	phonenumpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	service_line_item_pb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/service_line_item"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/interceptors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/service"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"

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

func (server *server) newEventProcessorGrpcServer() (func() error, *grpc.Server, error) {
	l, err := net.Listen(constants.Tcp, server.cfg.GRPC.Port)
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
			interceptors.CheckApiKeyInterceptor(server.cfg.GRPC.ApiKey),
		),
		),
	)
	RegisterGrpcServices(server, grpcServer)

	go func() {
		server.log.Infof("%s gRPC server is listening on port: {%s}", GetMicroserviceName(server.cfg), server.cfg.GRPC.Port)
		server.log.Error(grpcServer.Serve(l))
	}()

	return l.Close, grpcServer, nil
}

func RegisterGrpcServices(server *server, grpcServer *grpc.Server) {
	contactService := service.NewContactService(server.log, server.repositories, server.commandHandlers.Contact)
	contactpb.RegisterContactGrpcServiceServer(grpcServer, contactService)

	organizationService := service.NewOrganizationService(server.log, server.repositories, server.commandHandlers.Organization)
	organizationpb.RegisterOrganizationGrpcServiceServer(grpcServer, organizationService)

	phoneNumberService := service.NewPhoneNumberService(server.log, server.repositories, server.commandHandlers.PhoneNumber)
	phonenumpb.RegisterPhoneNumberGrpcServiceServer(grpcServer, phoneNumberService)

	emailService := service.NewEmailService(server.log, server.repositories, server.commandHandlers.Email)
	emailpb.RegisterEmailGrpcServiceServer(grpcServer, emailService)

	userService := service.NewUserService(server.log, server.commandHandlers.User)
	userpb.RegisterUserGrpcServiceServer(grpcServer, userService)

	locationService := service.NewLocationService(server.log, server.repositories, server.commandHandlers.Location)
	locationpb.RegisterLocationGrpcServiceServer(grpcServer, locationService)

	jobRoleService := service.NewJobRoleService(server.log, server.repositories, server.commandHandlers.JobRole)
	jobrolepb.RegisterJobRoleGrpcServiceServer(grpcServer, jobRoleService)

	interactionEventService := service.NewInteractionEventService(server.log, server.commandHandlers.InteractionEvent)
	iepb.RegisterInteractionEventGrpcServiceServer(grpcServer, interactionEventService)

	logEntryService := service.NewLogEntryService(server.log, server.commandHandlers.LogEntry)
	logentrypb.RegisterLogEntryGrpcServiceServer(grpcServer, logEntryService)

	issueService := service.NewIssueService(server.log, server.commandHandlers.Issue)
	issuepb.RegisterIssueGrpcServiceServer(grpcServer, issueService)

	commentService := service.NewCommentService(server.log, server.commandHandlers.Comment)
	commentpb.RegisterCommentGrpcServiceServer(grpcServer, commentService)

	opportunityService := service.NewOpportunityService(server.log, server.commandHandlers.Opportunity, server.aggregateStore)
	opportunitypb.RegisterOpportunityGrpcServiceServer(grpcServer, opportunityService)

	contractService := service.NewContractService(server.log, server.commandHandlers.Contract, server.aggregateStore)
	contractpb.RegisterContractGrpcServiceServer(grpcServer, contractService)

	serviceLineItemService := service.NewServiceLineItemService(server.log, server.commandHandlers.ServiceLineItem, server.aggregateStore)
	service_line_item_pb.RegisterServiceLineItemGrpcServiceServer(grpcServer, serviceLineItemService)
}
