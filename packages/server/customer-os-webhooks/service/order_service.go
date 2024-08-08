package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	orderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/order"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"sync"
	"time"
)

type OrderService interface {
	SyncOrders(ctx context.Context, orders []model.OrderData) (SyncResult, error)
}

type orderService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
	maxWorkers   int
}

func NewOrderService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) OrderService {
	return &orderService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
		maxWorkers:   services.cfg.ConcurrencyConfig.OrderSyncConcurrency,
	}
}

func (s *orderService) SyncOrders(ctx context.Context, orders []model.OrderData) (SyncResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderService.SyncOrders")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	// pre-validate order input before syncing
	for _, order := range orders {
		if order.ExternalSystem == "" {
			tracing.TraceErr(span, errors.ErrMissingExternalSystem)
			return SyncResult{}, errors.ErrMissingExternalSystem
		}
		if !neo4jentity.IsValidDataSource(strings.ToLower(order.ExternalSystem)) {
			tracing.TraceErr(span, errors.ErrExternalSystemNotAccepted, log.String("externalSystem", order.ExternalSystem))
			return SyncResult{}, errors.ErrExternalSystemNotAccepted
		}
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	workerLimit := make(chan struct{}, s.maxWorkers)

	syncMutex := &sync.Mutex{}
	statusesMutex := &sync.Mutex{}
	syncDate := utils.Now()
	var statuses []SyncStatus

	// Sync all orders
	for _, orderData := range orders {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return SyncResult{}, ctx.Err()
		default:
		}

		// Acquire a worker slot
		workerLimit <- struct{}{}
		wg.Add(1)

		go func(syncOrder model.OrderData) {
			defer wg.Done()
			defer func() {
				// Release the worker slot when done
				<-workerLimit
			}()

			result := s.syncOrder(ctx, syncMutex, syncOrder, syncDate)
			statusesMutex.Lock()
			statuses = append(statuses, result)
			statusesMutex.Unlock()
		}(orderData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), orders[0].ExternalSystem,
		orders[0].AppSource, "order", syncDate, statuses)

	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}

func (s *orderService) syncOrder(ctx context.Context, syncMutex *sync.Mutex, orderInput model.OrderData, syncDate time.Time) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderService.syncOrder")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagExternalSystem, orderInput.ExternalSystem)
	span.LogFields(log.Object("syncDate", syncDate))
	tracing.LogObjectAsJson(span, "orderInput", orderInput)

	var tenant = common.GetTenantFromContext(ctx)
	var failedSync = false
	var reason = ""
	orderInput.Normalize()

	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, orderInput.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err, log.String("externalSystem", orderInput.ExternalSystem))
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", orderInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}

	// Check if order sync should be skipped
	if orderInput.Skip {
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(orderInput.SkipReason)
	}

	// Lock order creation
	syncMutex.Lock()
	defer syncMutex.Unlock()
	// Check if order already exists
	orderId, err := s.repositories.Neo4jRepositories.OrderReadRepository.GetMatchedOrderId(ctx, tenant, orderInput.ExternalSystem, orderInput.ExternalId)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched order with external reference %s for tenant %s :%s", orderInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}

	// Check if organization exists
	organizationId, _, err := s.services.FinderService.FindReferencedEntityId(ctx, orderInput.ExternalSystem, &orderInput.OrderedByOrganization)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding organization with external reference %s for tenant %s :%s", orderInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}
	if organizationId == "" {
		failedSync = true
		reason = fmt.Sprintf("organization with external reference %s not found for tenant %s", orderInput.ExternalId, tenant)
		s.log.Error(reason)
	}

	if !failedSync {
		matchingOrderFound := orderId != ""
		span.LogFields(log.Bool("found matching order", matchingOrderFound))
		span.LogFields(log.String("orderId", orderId))

		var confirmedAt *timestamppb.Timestamp
		if orderInput.ConfirmedAt != nil {
			confirmedAt = timestamppb.New(*orderInput.ConfirmedAt)
		}
		var paidAt *timestamppb.Timestamp
		if orderInput.PaidAt != nil {
			paidAt = timestamppb.New(*orderInput.PaidAt)
		}
		var fulfilledAt *timestamppb.Timestamp
		if orderInput.FulfilledAt != nil {
			fulfilledAt = timestamppb.New(*orderInput.FulfilledAt)
		}
		var cancelledAt *timestamppb.Timestamp
		if orderInput.CanceledAt != nil {
			cancelledAt = timestamppb.New(*orderInput.CanceledAt)
		}

		request := orderpb.UpsertOrderGrpcRequest{
			Id:             orderId,
			Tenant:         tenant,
			OrganizationId: organizationId,
			ConfirmedAt:    confirmedAt,
			PaidAt:         paidAt,
			FulfilledAt:    fulfilledAt,
			CanceledAt:     cancelledAt,
			CreatedAt:      timestamppb.New(utils.TimePtrAsAny(orderInput.CreatedAt, utils.NowPtr()).(time.Time)),
			UpdatedAt:      timestamppb.New(utils.TimePtrAsAny(orderInput.UpdatedAt, utils.NowPtr()).(time.Time)),
			SourceFields: &commonpb.SourceFields{
				Source:    orderInput.ExternalSystem,
				AppSource: utils.StringFirstNonEmpty(orderInput.AppSource, constants.AppSourceCustomerOsWebhooks),
			},
			ExternalSystemFields: &commonpb.ExternalSystemFields{
				ExternalSystemId: orderInput.ExternalSystem,
				ExternalId:       orderInput.ExternalId,
				ExternalSource:   orderInput.ExternalSourceEntity,
				ExternalUrl:      orderInput.ExternalUrl,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
		}

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		response, err := CallEventsPlatformGRPCWithRetry[*orderpb.OrderIdGrpcResponse](func() (*orderpb.OrderIdGrpcResponse, error) {
			return s.grpcClients.OrderClient.UpsertOrder(ctx, &request)
		})
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err, log.String("grpcMethod", "UpsertOrder"))
			reason = fmt.Sprintf("failed sending event to upsert order with external reference %s for tenant %s :%s", orderInput.ExternalId, tenant, err.Error())
			s.log.Error(reason)
		} else {
			orderId = response.GetId()
		}
		// Wait for order to be created in neo4j
		if !failedSync && !matchingOrderFound {
			for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
				order, forErr := s.repositories.Neo4jRepositories.OrderReadRepository.GetById(ctx, tenant, orderId)
				if order != nil && forErr == nil {
					break
				}
				time.Sleep(utils.BackOffExponentialDelay(i))
			}
		}
	}

	span.LogFields(log.Bool("failedSync", failedSync))
	if failedSync {
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}
	span.LogFields(log.String("output", "success"))
	return NewSuccessfulSyncStatus()
}
