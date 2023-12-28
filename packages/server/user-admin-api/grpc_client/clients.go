package grpc_client

import (
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/issue"
	"google.golang.org/grpc"
)

type Clients struct {
	IssueClient issuepb.IssueGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		IssueClient: issuepb.NewIssueGrpcServiceClient(conn),
	}
	return &clients
}
