package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	utils2 "github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type invoiceService struct {
	invoicepb.UnimplementedInvoiceGrpcServiceServer
	repositories   *repository.Repositories
	services       *Services
	log            logger.Logger
	aggregateStore eventstore.AggregateStore
}

func NewInvoiceService(repositories *repository.Repositories, services *Services, log logger.Logger, aggregateStore eventstore.AggregateStore) *invoiceService {
	return &invoiceService{
		repositories:   repositories,
		services:       services,
		log:            log,
		aggregateStore: aggregateStore,
	}
}

func (s *invoiceService) NextPreviewInvoiceForContract(ctx context.Context, request *invoicepb.NextPreviewInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.NextPreviewInvoiceForContract")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, "")
	tracing.LogObjectAsJson(span, "request", request)

	// check and fetch contract
	contractExists, err := s.checkContractExists(ctx, request.Tenant, request.ContractId)
	if err != nil {
		s.log.Error(err, "error checking contract existence")
		return nil, status.Errorf(codes.Internal, "error checking contract existence: %v", err)
	}
	if !contractExists {
		return nil, status.Errorf(codes.NotFound, "contract with ID %s not found", request.ContractId)
	}

	contract, err := s.getContract(ctx, request.Tenant, request.ContractId)
	if err != nil {
		s.log.Errorf("Error while getting contract %s: %s", request.ContractId, err.Error())
		tracing.TraceErr(span, err)
		return nil, err
	}

	//get last issued on cycle invoice for contract
	//lastIssuedOnCycleInvoiceForContract, err := s.GetLastIssuedOnCycleInvoiceForContract(ctx, request.Tenant, request.EntityId)
	//if err != nil {
	//	s.log.Errorf("Error while getting last issued on cycle invoice for contract %s: %s", request.EntityId, err.Error())
	//	tracing.TraceErr(span, err)
	//	return nil, err
	//}

	var invoicePeriodStart, invoicePeriodEnd time.Time
	if contract.NextInvoiceDate != nil {
		invoicePeriodStart = *contract.NextInvoiceDate
	} else if contract.InvoicingStartDate != nil {
		invoicePeriodStart = *contract.InvoicingStartDate
	} else {
		err = fmt.Errorf("contract has no next invoice date or invoicing start date")
		tracing.TraceErr(span, err)
		return nil, err
	}
	invoicePeriodEnd = s.calculateInvoiceCycleEnd(ctx, invoicePeriodStart, request.Tenant, *contract)

	now := utils.Now()
	invoiceId := uuid.New().String()

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, invoiceId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, &invoicepb.NewInvoiceForContractRequest{
		Tenant:         request.Tenant,
		LoggedInUserId: "",
		ContractId:     request.ContractId,
		CreatedAt:      utils.ConvertTimeToTimestampPtr(&now),
		SourceFields: &commonpb.SourceFields{
			AppSource: utils2.AppSourceEventProcessingPlatform,
		},
		InvoicePeriodStart:   utils.ConvertTimeToTimestampPtr(&invoicePeriodStart),
		InvoicePeriodEnd:     utils.ConvertTimeToTimestampPtr(&invoicePeriodEnd),
		Currency:             contract.Currency.String(),
		BillingCycleInMonths: contract.BillingCycleInMonths,
		Note:                 "",
		DryRun:               true,
		Preview:              true,
		OffCycle:             false,
		Postpaid:             s.getTenantInvoicingPostpaidFlag(ctx, request.Tenant),
	}); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(NewOnCycleInvoiceForContract) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: invoiceId}, nil
}

func (s *invoiceService) GetLastIssuedOnCycleInvoiceForContract(ctx context.Context, tenant, contractId string) (*neo4jentity.InvoiceEntity, error) {
	lastIssuedOnCycleInvoiceForContractNode, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetLastIssuedOnCycleInvoiceForContract(ctx, tenant, contractId)
	if err != nil {
		return nil, err
	}
	entity := neo4jmapper.MapDbNodeToInvoiceEntity(lastIssuedOnCycleInvoiceForContractNode)

	if *entity == (neo4jentity.InvoiceEntity{}) {
		return nil, nil
	} else {
		return entity, nil
	}
}

func (s *invoiceService) getContract(ctx context.Context, tenant, contractId string) (*neo4jentity.ContractEntity, error) {
	contractNode, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToContractEntity(contractNode), nil
}

func (s *invoiceService) NewInvoiceForContract(ctx context.Context, request *invoicepb.NewInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.NewInvoiceForContract")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.Currency == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("currency"))
	} else if request.InvoicePeriodStart == nil {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoicePeriodStart"))
	} else if request.InvoicePeriodEnd == nil {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoicePeriodEnd"))
	} else if request.ContractId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("contractId"))
	}

	// Check if the contract aggregate exists
	contractExists, err := s.checkContractExists(ctx, request.Tenant, request.ContractId)
	if err != nil {
		s.log.Error(err, "error checking contract existence")
		return nil, status.Errorf(codes.Internal, "error checking contract existence: %v", err)
	}
	if !contractExists {
		return nil, status.Errorf(codes.NotFound, "contract with ID %s not found", request.ContractId)
	}

	invoiceId := uuid.New().String()

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, invoiceId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(NewOnCycleInvoiceForContract) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: invoiceId}, nil
}

func (s *invoiceService) prepareInvoiceNumber(tenant string) string {
	maxAttempts := 20
	var invoiceNumber string
	for attempt := 1; attempt < maxAttempts+1; attempt++ {
		invoiceNumber = s.services.CommonServices.InvoiceService.GenerateNewRandomInvoiceNumber()
		invoiceNumberEntity := postgresentity.InvoiceNumberEntity{
			InvoiceNumber: invoiceNumber,
			Tenant:        tenant,
			Attempts:      attempt,
		}
		innerErr := s.repositories.InvoiceRepository.Reserve(invoiceNumberEntity)
		if innerErr == nil {
			break
		}
	}

	return invoiceNumber
}

func (s *invoiceService) FillInvoice(ctx context.Context, request *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.FillInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	extraParams := map[string]any{}

	if request.InvoiceNumber != "" {
		extraParams[invoice.PARAM_INVOICE_NUMBER] = request.InvoiceNumber
	} else {
		if !request.DryRun || request.Preview {
			extraParams[invoice.PARAM_INVOICE_NUMBER] = s.prepareInvoiceNumber(request.Tenant)
		} else {
			extraParams[invoice.PARAM_INVOICE_NUMBER] = s.services.CommonServices.InvoiceService.GenerateNewRandomInvoiceNumber()
		}
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request, extraParams); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(FillInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) GenerateInvoicePdf(ctx context.Context, request *invoicepb.GenerateInvoicePdfRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.GenerateInvoicePdf")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(GenerateInvoicePdf) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) PdfGeneratedInvoice(ctx context.Context, request *invoicepb.PdfGeneratedInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.PdfGeneratedInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(PdfGeneratedInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) PayInvoiceNotification(ctx context.Context, request *invoicepb.PayInvoiceNotificationRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.PayInvoiceNotification")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(PayInvoiceNotification) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) checkContractExists(ctx context.Context, tenant, contractId string) (bool, error) {
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	err := s.aggregateStore.Exists(ctx, contractAggregate.GetID())
	if err != nil {
		if errors.Is(err, eventstore.ErrAggregateNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil // The contract exists
}

func (s *invoiceService) RequestFillInvoice(ctx context.Context, request *invoicepb.RequestFillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.RequestFillInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RequestFillInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) PermanentlyDeleteInitializedInvoice(ctx context.Context, request *invoicepb.PermanentlyDeleteInitializedInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.PermanentlyDeleteInitializedInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(PermanentlyDeleteInitializedInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) VoidInvoice(ctx context.Context, request *invoicepb.VoidInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.VoidInvoice")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(VoidInvoice) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) RemindInvoiceNotification(ctx context.Context, request *invoicepb.RemindInvoiceNotificationRequest) (*invoicepb.InvoiceIdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "InvoiceService.RemindInvoiceNotification")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.InvoiceId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("invoiceId"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return invoice.NewInvoiceAggregateWithTenantAndID(request.Tenant, request.InvoiceId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, *eventstore.NewLoadAggregateOptionsWithRequired(), request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RemindInvoiceNotification) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &invoicepb.InvoiceIdResponse{Id: request.InvoiceId}, nil
}

func (s *invoiceService) calculateInvoiceCycleEnd(ctx context.Context, start time.Time, tenant string, contractEntity neo4jentity.ContractEntity) time.Time {
	nextStart := start.AddDate(0, int(contractEntity.BillingCycleInMonths), 0)
	if start.Day() == 1 {
		// if previous invoice was generated end of month, we need to substract extra 1 day
		previousCycleInvoiceDbNode, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetPreviousCycleInvoice(ctx, tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(nil, err)
		}
		if previousCycleInvoiceDbNode != nil {
			previousInvoice := neo4jmapper.MapDbNodeToInvoiceEntity(previousCycleInvoiceDbNode)
			if previousInvoice.PeriodStartDate.Day() != 1 {
				nextStart = nextStart.AddDate(0, -1, 0)
				nextStart = time.Date(nextStart.Year(), nextStart.Month(), previousInvoice.PeriodStartDate.Day(), 0, 0, 0, 0, nextStart.Location())
			}
		}
	}
	return nextStart.AddDate(0, 0, -1)
}

func (s *invoiceService) getTenantInvoicingPostpaidFlag(ctx context.Context, tenant string) bool {
	dbNode, _ := s.repositories.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
	tenantSettings := neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode)
	return tenantSettings.InvoicingPostpaid
}
