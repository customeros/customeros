package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicingcyclepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoicing_cycle"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type InvoicingCycleService interface {
	CreateInvoicingCycle(ctx context.Context, invoicingCycleType invoicingcyclepb.InvoicingDateType) (string, error)
	UpdateInvoicingCycle(ctx context.Context, id string, invoicingCycleType invoicingcyclepb.InvoicingDateType) error
	GetInvoicingCycle(ctx context.Context) (*neo4jentity.InvoicingCycleEntity, error)
}
type invoicingCycleService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewInvoicingCycleService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) InvoicingCycleService {
	return &invoicingCycleService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *invoicingCycleService) CreateInvoicingCycle(ctx context.Context, invoicingCycleType invoicingcyclepb.InvoicingDateType) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleService.CreateInvoicingCycle")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("invoicingCycleType", string(invoicingCycleType)))

	grpcRequest := invoicingcyclepb.CreateInvoicingCycleTypeRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Type:           invoicingCycleType,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*invoicingcyclepb.InvoicingCycleTypeResponse](func() (*invoicingcyclepb.InvoicingCycleTypeResponse, error) {
		return s.grpcClients.InvoicingCycleClient.CreateInvoicingCycleType(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, neo4jutil.NodeLabelInvoicingCycle, span)

	return response.Id, nil
}

func (s *invoicingCycleService) UpdateInvoicingCycle(ctx context.Context, id string, invoicingCycleType invoicingcyclepb.InvoicingDateType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleService.UpdateInvoicingCycle")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, id)
	span.LogFields(log.Object("invoicingCycleType", invoicingCycleType))

	exists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), id, neo4jutil.NodeLabelInvoicingCycle)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !exists {
		err = errors.New(fmt.Sprintf("Invoicing cycle with id {%s} not found", id))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := invoicingcyclepb.UpdateInvoicingCycleTypeRequest{
		Tenant:               common.GetTenantFromContext(ctx),
		InvoicingCycleTypeId: id,
		LoggedInUserId:       common.GetUserIdFromContext(ctx),
		Type:                 invoicingCycleType,
		SourceFields: &commonpb.SourceFields{
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*invoicingcyclepb.InvoicingCycleTypeResponse](func() (*invoicingcyclepb.InvoicingCycleTypeResponse, error) {
		return s.grpcClients.InvoicingCycleClient.UpdateInvoicingCycleType(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}
	return nil
}

func (s *invoicingCycleService) GetInvoicingCycle(ctx context.Context) (*neo4jentity.InvoicingCycleEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleService.GetInvoicingCycle")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if dbNode, err := s.repositories.Neo4jRepositories.InvoicingCycleReadRepository.GetInvoicingCycle(ctx, common.GetContext(ctx).Tenant); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Invoicing cycle not found"))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToInvoicingCycleEntity(dbNode), nil
	}
}
