package events_platform

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	eventcompletionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_completion"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	jobrolepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/job_role"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"

	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
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
	contractpb.RegisterContractGrpcServiceServer(server, &MockContractService{})
	servicelineitempb.RegisterServiceLineItemGrpcServiceServer(server, &MockServiceLineItemService{})
	opportunitypb.RegisterOpportunityGrpcServiceServer(server, &MockOpportunityService{})
	invoicepb.RegisterInvoiceGrpcServiceServer(server, &MockInvoiceService{})
	tenantpb.RegisterTenantGrpcServiceServer(server, &MockTenantService{})
	eventstorepb.RegisterEventStoreGrpcServiceServer(server, &MockEventStoreService{})
	eventcompletionpb.RegisterEventCompletionGrpcServiceServer(server, &MockEventCompletionService{})

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
