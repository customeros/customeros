package invoice

import (
	"context"
	"strings"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	commonServices "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/invoice"
	"github.com/customeros/customeros/packages/server/events-processing-platform/tracing"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/config"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/constants"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/logger"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/service"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/subscriptions"
)

type InvoiceSubscriber struct {
	log                 logger.Logger
	db                  *esdb.Client
	cfg                 *config.Config
	grpcClients         *grpc_client.Clients
	invoiceEventHandler *InvoiceEventHandler
	CommonServices      *commonServices.CommonServices
}

func NewInvoiceSubscriber(log logger.Logger, db *esdb.Client, cfg *config.Config, service *service.Services, grpcClients *grpc_client.Clients) *InvoiceSubscriber {
	return &InvoiceSubscriber{
		log:         log,
		db:          db,
		cfg:         cfg,
		grpcClients: grpcClients,
		invoiceEventHandler: NewInvoiceEventHandler(
			log,
			*cfg,
			grpcClients,
			service.Neo4jRepositories,
			service.PostgresRepositories,
			service.CommonServices.InvoiceService,
			service.CommonServices.FileService,
			service.CommonServices.PostmarkService,
		),
	}
}

func (s *InvoiceSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.InvoiceSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.InvoiceSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.InvoiceSubscription.BufferSizeClient,
			},
		)
		if err != nil {
			return err
		}
		defer sub.Close()

		group.Go(s.runWorker(ctx, worker, sub, i))
	}
	return group.Wait()
}

func (consumer *InvoiceSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *InvoiceSubscriber) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {
	for {
		event := stream.Recv()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if event.SubscriptionDropped != nil {
			s.log.Errorf("(SubscriptionDropped) err: {%v}", event.SubscriptionDropped.Error)
			return errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
		}

		if event.EventAppeared != nil {
			s.log.EventAppeared(s.cfg.Subscriptions.InvoiceSubscription.GroupName, event.EventAppeared.Event, workerID)

			if event.EventAppeared.Event.Event == nil {
				s.log.Errorf("(InvoiceSubscriber) event.EventAppeared.Event.Event is nil")
			} else {
				err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
				if err != nil {
					s.log.Errorf("(InvoiceSubscription.when) err: {%v}", err)

					if err := stream.Nack(err.Error(), esdb.NackActionPark, event.EventAppeared.Event); err != nil {
						s.log.Errorf("(stream.Nack) err: {%v}", err)
						return errors.Wrap(err, "stream.Nack")
					}
				}
			}

			err := stream.Ack(event.EventAppeared.Event)
			if err != nil {
				s.log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}
			s.log.Debugf("(ACK) event: {%+v}", eventstore.NewRecordedBaseEventFromRecorded(event.EventAppeared.Event.Event))
		}
	}
}

func (s *InvoiceSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	if strings.HasPrefix(evt.GetAggregateID(), constants.EsInternalStreamPrefix) {
		return nil
	}

	acceptedEventTypes := []string{
		invoice.InvoiceFillRequestedV1,
		invoice.InvoicePdfRequestedV1,
		invoice.InvoicePdfGeneratedV1,
		invoice.InvoiceVoidV1,
		invoice.InvoicePayNotificationV1,
		invoice.InvoiceRemindNotificationV1,
	}

	if !utils.Contains(acceptedEventTypes, evt.GetEventType()) {
		return nil
	}

	ctx, span := tracing.StartProjectionTracerSpan(ctx, "InvoiceSubscriber.When", evt)
	defer span.Finish()

	if s.cfg.Subscriptions.InvoiceSubscription.IgnoreEvents {
		return nil
	}

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: eventstore.GetTenantFromAggregate(evt.GetAggregateID(), invoice.InvoiceAggregateType),
	})

	switch evt.GetEventType() {
	case invoice.InvoiceFillRequestedV1:
		return s.invoiceEventHandler.onInvoiceFillRequestedV1(ctx, evt)
	case invoice.InvoicePdfRequestedV1:
		return s.invoiceEventHandler.generateInvoicePDFV1(ctx, evt)
	case invoice.InvoicePdfGeneratedV1:
		return s.invoiceEventHandler.onInvoicePdfGeneratedV1(ctx, evt)
	case invoice.InvoiceVoidV1:
		return s.invoiceEventHandler.onInvoiceVoidV1(ctx, evt)
	case invoice.InvoicePaidV1:
		return s.invoiceEventHandler.onInvoicePaidV1(ctx, evt)
	case invoice.InvoicePayNotificationV1:
		return s.invoiceEventHandler.onInvoicePayNotificationV1(ctx, evt)
	case invoice.InvoiceRemindNotificationV1:
		return s.invoiceEventHandler.onInvoiceRemindNotificationV1(ctx, evt)
	default:
		return nil
	}
}
