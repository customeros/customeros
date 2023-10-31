package events_paltform

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	jobrolepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/job_role"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
)

type TestDialFactoryImpl struct {
	eventsProcessingPlatformConn *grpc.ClientConn
}

func (dfi TestDialFactoryImpl) Close(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}

func (dfi TestDialFactoryImpl) GetEventsProcessingPlatformConn() (*grpc.ClientConn, error) {
	if dfi.eventsProcessingPlatformConn != nil {
		return dfi.eventsProcessingPlatformConn, nil
	}
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	contactpb.RegisterContactGrpcServiceServer(server, &MockContactService{})
	emailpb.RegisterEmailGrpcServiceServer(server, &MockEmailService{})
	phonenumberpb.RegisterPhoneNumberGrpcServiceServer(server, &MockPhoneNumberService{})
	jobrolepb.RegisterJobRoleGrpcServiceServer(server, &MockJobRoleService{})
	userpb.RegisterUserGrpcServiceServer(server, &MockUserService{})
	organizationpb.RegisterOrganizationGrpcServiceServer(server, &MockOrganizationService{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return grpc.DialContext(context.Background(), "", grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, target string) (net.Conn, error) {
			return listener.Dial()
		}))
}

func NewTestDialFactory() grpc_client.DialFactory {
	dfi := new(TestDialFactoryImpl)
	return *dfi
}
