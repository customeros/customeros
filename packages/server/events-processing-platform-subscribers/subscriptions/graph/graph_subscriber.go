package graph

import (
	"context"
	"strings"
	"time"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	invoiceevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
)

type GraphSubscriber struct {
	log                      logger.Logger
	db                       *esdb.Client
	cfg                      *config.Config
	organizationEventHandler *OrganizationEventHandler
	invoiceEventHandler      *InvoiceEventHandler
	services                 *service.Services
}

func NewGraphSubscriber(log logger.Logger, db *esdb.Client, services *service.Services, grpcClients *grpc_client.Clients, cfg *config.Config, cache caches.Cache) *GraphSubscriber {
	return &GraphSubscriber{
		log:                      log,
		db:                       db,
		cfg:                      cfg,
		organizationEventHandler: NewOrganizationEventHandler(log, grpcClients, cache, services.CommonServices.Events, services.Neo4jRepositories, services.PostgresRepositories, services.CommonServices.CurrencyService),
		invoiceEventHandler:      NewInvoiceEventHandler(log, grpcClients, services.Neo4jRepositories),
	}
}

func (s *GraphSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.GraphSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.GraphSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.GraphSubscription.BufferSizeClient,
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

func (s *GraphSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *GraphSubscriber) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {
	for {
		event := stream.Recv()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if event.SubscriptionDropped != nil {
			span, _ := opentracing.StartSpanFromContext(ctx, "GraphSubscriber.ProcessEvents")
			defer span.Finish()
			wrappedErr := errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
			tracing.TraceErr(span, wrappedErr)
			s.log.Errorf(wrappedErr.Error())
			return wrappedErr
		}

		if event.EventAppeared != nil {
			s.log.EventAppeared(s.cfg.Subscriptions.GraphSubscription.GroupName, event.EventAppeared.Event, workerID)

			if event.EventAppeared.Event.Event == nil {
				span, _ := opentracing.StartSpanFromContext(ctx, "GraphSubscriber.ProcessEvents")
				defer span.Finish()
				err := errors.Wrap(errors.New("event.EventAppeared.Event.Event is nil"), "GraphSubscriber")
				tracing.TraceErr(span, err)
				s.log.Errorf(err.Error())
			} else {
				err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
				if err != nil {
					span, _ := opentracing.StartSpanFromContext(ctx, "GraphSubscriber.ProcessEvents")
					defer span.Finish()
					tracing.TraceErr(span, err)
					s.log.Errorf("(GraphSubscriber.when) err: {%v}", err)

					if err := stream.Nack(err.Error(), esdb.NackActionPark, event.EventAppeared.Event); err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("(stream.Nack) err: {%v}", err)
						return errors.Wrap(err, "stream.Nack")
					}
				}
			}

			err := stream.Ack(event.EventAppeared.Event)
			if err != nil {
				span, _ := opentracing.StartSpanFromContext(ctx, "GraphSubscriber.ProcessEvents")
				defer span.Finish()
				tracing.TraceErr(span, err)
				s.log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}
			s.log.Debugf("(ACK) event: {%+v}", eventstore.NewRecordedBaseEventFromRecorded(event.EventAppeared.Event.Event))
		}
	}
}

func (s *GraphSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	if strings.HasPrefix(evt.GetAggregateID(), constants.EsInternalStreamPrefix) {
		return nil
	}

	ctx, span := tracing.StartProjectionTracerSpan(ctx, "GraphSubscriber.When", evt)
	defer span.Finish()

	// set 25 sec context deadline
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	switch evt.GetEventType() {

	case orgevents.OrganizationPhoneNumberLinkV1:
		_ = s.organizationEventHandler.OnPhoneNumberLinkedToOrganization(ctx, evt)
		return nil
	case orgevents.OrganizationRefreshArrV1:
		_ = s.organizationEventHandler.OnRefreshArr(ctx, evt)
		return nil
	case orgevents.OrganizationRefreshRenewalSummaryV1:
		_ = s.organizationEventHandler.OnRefreshRenewalSummaryV1(ctx, evt)
		return nil
	case orgevents.OrganizationRefreshDerivedDataV1:
		_ = s.organizationEventHandler.OnRefreshDerivedDataV1(ctx, evt)
		return nil
	case orgevents.OrganizationUpsertCustomFieldV1:
		_ = s.organizationEventHandler.OnUpsertCustomField(ctx, evt)
		return nil
	case orgevents.OrganizationCreateBillingProfileV1:
		_ = s.organizationEventHandler.OnCreateBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationUpdateBillingProfileV1:
		_ = s.organizationEventHandler.OnUpdateBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationEmailLinkToBillingProfileV1:
		_ = s.organizationEventHandler.OnEmailLinkedToBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationEmailUnlinkFromBillingProfileV1:
		_ = s.organizationEventHandler.OnEmailUnlinkedFromBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationLocationLinkToBillingProfileV1:
		_ = s.organizationEventHandler.OnLocationLinkedToBillingProfile(ctx, evt)
		return nil
	case orgevents.OrganizationLocationUnlinkFromBillingProfileV1:
		_ = s.organizationEventHandler.OnLocationUnlinkedFromBillingProfile(ctx, evt)
		return nil

	case invoiceevents.InvoiceCreateForContractV1:
		_ = s.invoiceEventHandler.OnInvoiceCreateForContractV1(ctx, evt)
		return nil
	case invoiceevents.InvoiceFillV1:
		_ = s.invoiceEventHandler.OnInvoiceFillV1(ctx, evt)
		return nil
	case invoiceevents.InvoicePdfGeneratedV1:
		_ = s.invoiceEventHandler.OnInvoicePdfGenerated(ctx, evt)
		return nil
	case invoiceevents.InvoiceVoidV1:
		_ = s.invoiceEventHandler.OnInvoiceVoidV1(ctx, evt)
		return nil
	case invoiceevents.InvoiceDeleteV1:
		_ = s.invoiceEventHandler.OnInvoiceDeleteV1(ctx, evt)
		return nil

	default:
		return nil
	}
}
