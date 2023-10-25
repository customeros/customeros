package server

import (
	contactgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	emailgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	iegrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	issuegrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/issue"
	jobrolegrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/job_role"
	locationgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	logentrygrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	orggrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	phonenumbergrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	usergrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
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
	contactService := service.NewContactService(server.log, server.repositories, server.commands.ContactCommands)
	contactgrpc.RegisterContactGrpcServiceServer(grpcServer, contactService)

	organizationService := service.NewOrganizationService(server.log, server.repositories, server.commands.OrganizationCommands)
	orggrpc.RegisterOrganizationGrpcServiceServer(grpcServer, organizationService)

	phoneNumberService := service.NewPhoneNumberService(server.log, server.repositories, server.commands.PhoneNumberCommands)
	phonenumbergrpc.RegisterPhoneNumberGrpcServiceServer(grpcServer, phoneNumberService)

	emailService := service.NewEmailService(server.log, server.repositories, server.commands.EmailCommands)
	emailgrpc.RegisterEmailGrpcServiceServer(grpcServer, emailService)

	userService := service.NewUserService(server.log, server.commands.UserCommands)
	usergrpc.RegisterUserGrpcServiceServer(grpcServer, userService)

	locationService := service.NewLocationService(server.log, server.repositories, server.commands.LocationCommands)
	locationgrpc.RegisterLocationGrpcServiceServer(grpcServer, locationService)

	jobRoleService := service.NewJobRoleService(server.log, server.repositories, server.commands.JobRoleCommands)
	jobrolegrpc.RegisterJobRoleGrpcServiceServer(grpcServer, jobRoleService)

	interactionEventService := service.NewInteractionEventService(server.log, server.commands.InteractionEventCommands)
	iegrpc.RegisterInteractionEventGrpcServiceServer(grpcServer, interactionEventService)

	logEntryService := service.NewLogEntryService(server.log, server.commands.LogEntryCommands)
	logentrygrpc.RegisterLogEntryGrpcServiceServer(grpcServer, logEntryService)

	issueService := service.NewIssueService(server.log, server.commands.IssueCommands)
	issuegrpc.RegisterIssueGrpcServiceServer(grpcServer, issueService)
}
