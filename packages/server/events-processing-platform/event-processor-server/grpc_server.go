package server

import (
	contact_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	email_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	location_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	phone_number_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	user_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	contact_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/service"
	email_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/service"
	location_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/service"
	organization_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/service"
	phone_number_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/service"
	user_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/service"
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
		),
		),
	)

	RegisterGrpcServices(server, grpcServer)

	//userService := organization_service.NewOrganizationService(server.log, server.repositories, server.commands.OrganizationCommands)
	//organization_grpc_service.RegisterOrganizationGrpcServiceServer(grpcServer, organizationService)

	go func() {
		server.log.Infof("%s gRPC server is listening on port: {%s}", GetMicroserviceName(server.cfg), server.cfg.GRPC.Port)
		server.log.Error(grpcServer.Serve(l))
	}()

	return l.Close, grpcServer, nil
}

func RegisterGrpcServices(server *server, grpcServer *grpc.Server) {
	contactService := contact_service.NewContactService(server.log, server.repositories, server.commands.ContactCommands)
	contact_grpc_service.RegisterContactGrpcServiceServer(grpcServer, contactService)

	organizationService := organization_service.NewOrganizationService(server.log, server.repositories, server.commands.OrganizationCommands)
	organization_grpc_service.RegisterOrganizationGrpcServiceServer(grpcServer, organizationService)

	phoneNumberService := phone_number_service.NewPhoneNumberService(server.log, server.repositories, server.commands.PhoneNumberCommands)
	phone_number_grpc_service.RegisterPhoneNumberGrpcServiceServer(grpcServer, phoneNumberService)

	emailService := email_service.NewEmailService(server.log, server.repositories, server.commands.EmailCommands)
	email_grpc_service.RegisterEmailGrpcServiceServer(grpcServer, emailService)

	userService := user_service.NewUserService(server.log, server.repositories, server.commands.UserCommands)
	user_grpc_service.RegisterUserGrpcServiceServer(grpcServer, userService)

	locationService := location_service.NewLocationService(server.log, server.repositories, server.commands.LocationCommands)
	location_grpc_service.RegisterLocationGrpcServiceServer(grpcServer, locationService)
}
