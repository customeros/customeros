package projection

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
	"golang.org/x/net/context"
)

type graphProjection struct {
	log logger.Logger
	db  *esdb.Client
	cfg *config.Config
}

func NewContactProjection(db *esdb.Client, cfg *config.Config) *graphProjection {
	return &graphProjection{db: db, cfg: cfg}
}

type Worker func(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error

//func (o *graphProjection) Subscribe(ctx context.Context, prefixes []string, poolSize int, worker Worker) error {
//	o.log.Infof("(starting order subscription) prefixes: {%+v}", prefixes)
//
//	err := o.db.CreatePersistentSubscriptionAll(ctx, o.cfg.Subscriptions.MongoProjectionGroupName, esdb.PersistentAllSubscriptionOptions{
//		Filter: &esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: prefixes},
//	})
//	if err != nil {
//		if subscriptionError, ok := err.(*esdb.PersistentSubscriptionError); !ok || ok && (subscriptionError.Code != 6) {
//			o.log.Errorf("(CreatePersistentSubscriptionAll) err: {%v}", subscriptionError.Error())
//		}
//	}
//
//	stream, err := o.db.ConnectToPersistentSubscription(
//		ctx,
//		constants.EsAll,
//		o.cfg.Subscriptions.MongoProjectionGroupName,
//		esdb.ConnectToPersistentSubscriptionOptions{},
//	)
//	if err != nil {
//		return err
//	}
//	defer stream.Close()
//
//	g, ctx := errgroup.WithContext(ctx)
//	for i := 0; i <= poolSize; i++ {
//		g.Go(o.runWorker(ctx, worker, stream, i))
//	}
//	return g.Wait()
//}
//
//func (o *mongoProjection) runWorker(ctx context.Context, worker Worker, stream *esdb.PersistentSubscription, i int) func() error {
//	return func() error {
//		return worker(ctx, stream, i)
//	}
//}
//
//func (o *mongoProjection) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {
//
//	for {
//		event := stream.Recv()
//		select {
//		case <-ctx.Done():
//			return ctx.Err()
//		default:
//		}
//
//		if event.SubscriptionDropped != nil {
//			o.log.Errorf("(SubscriptionDropped) err: {%v}", event.SubscriptionDropped.Error)
//			return errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
//		}
//
//		if event.EventAppeared != nil {
//			o.log.ProjectionEvent(constants.MongoProjection, o.cfg.Subscriptions.MongoProjectionGroupName, event.EventAppeared, workerID)
//
//			err := o.When(ctx, es.NewEventFromRecorded(event.EventAppeared.Event))
//			if err != nil {
//				o.log.Errorf("(mongoProjection.when) err: {%v}", err)
//
//				if err := stream.Nack(err.Error(), esdb.Nack_Retry, event.EventAppeared); err != nil {
//					o.log.Errorf("(stream.Nack) err: {%v}", err)
//					return errors.Wrap(err, "stream.Nack")
//				}
//			}
//
//			err = stream.Ack(event.EventAppeared)
//			if err != nil {
//				o.log.Errorf("(stream.Ack) err: {%v}", err)
//				return errors.Wrap(err, "stream.Ack")
//			}
//			o.log.Infof("(ACK) event commit: {%v}", *event.EventAppeared.Commit)
//		}
//	}
//}
//
//func (o *mongoProjection) When(ctx context.Context, evt es.Event) error {
//	ctx, span := tracing.StartProjectionTracerSpan(ctx, "mongoProjection.When", evt)
//	defer span.Finish()
//	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))
//
//	switch evt.GetEventType() {
//
//	case v1.OrderCreated:
//		return o.onOrderCreate(ctx, evt)
//	case v1.OrderPaid:
//		return o.onOrderPaid(ctx, evt)
//	case v1.OrderSubmitted:
//		return o.onSubmit(ctx, evt)
//	case v1.ShoppingCartUpdated:
//		return o.onShoppingCartUpdate(ctx, evt)
//	case v1.OrderCanceled:
//		return o.onCancel(ctx, evt)
//	case v1.OrderCompleted:
//		return o.onCompleted(ctx, evt)
//	case v1.DeliveryAddressChanged:
//		return o.onDeliveryAddressChnaged(ctx, evt)
//
//	default:
//		o.log.Warnf("(mongoProjection) [When unknown EventType] eventType: {%s}", evt.EventType)
//		return es.ErrInvalidEventType
//	}
//}
