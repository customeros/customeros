package store

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	es "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"math"
)

const (
	count = math.MaxInt64
)

type aggregateStore struct {
	log        logger.Logger
	esdbClient *esdb.Client
}

func NewAggregateStore(log logger.Logger, esdbClient *esdb.Client) *aggregateStore {
	return &aggregateStore{log: log, esdbClient: esdbClient}
}

func (as *aggregateStore) Load(ctx context.Context, aggregate es.Aggregate) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AggregateStore.Load")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", aggregate.GetID()))

	stream, err := as.esdbClient.ReadStream(ctx, aggregate.GetID(), esdb.ReadStreamOptions{}, count)
	if err != nil {
		tracing.TraceErr(span, err)
		as.log.Errorf("(Load) esdbClient.ReadStream: {%+v}", err)
		return errors.Wrap(err, "esdbClient.ReadStream")
	}
	defer stream.Close()

	for {
		event, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			tracing.TraceErr(span, err)
			as.log.Errorf("(Load) esdbClient.ReadStream: {%+v}", err)
			return errors.Wrap(err, "stream.Recv")
		}

		esEvent := es.NewEventFromRecorded(event.Event)
		if err := aggregate.RaiseEvent(esEvent); err != nil {
			tracing.TraceErr(span, err)
			as.log.Errorf("(Load) aggregate.RaiseEvent: {%+v}", err)
			return errors.Wrap(err, "RaiseEvent")
		}
		as.log.Debugf("(Load) esEvent: {%s}", esEvent.String())
	}

	as.log.Debugf("(Load) aggregate: {%s}", aggregate.String())
	return nil
}

func (as *aggregateStore) Save(ctx context.Context, aggregate es.Aggregate) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AggregateStore.Save")
	defer span.Finish()
	span.LogFields(
		log.String("aggregateID", aggregate.GetID()),
		log.Int64("version", aggregate.GetVersion()),
		log.Object("aggregateType", aggregate.GetType()),
		log.Int("uncommittedEvents", len(aggregate.GetUncommittedEvents())),
		log.Int("appliedEvents", len(aggregate.GetAppliedEvents())))

	if len(aggregate.GetUncommittedEvents()) == 0 {
		as.log.Debugf("(Save) [no uncommittedEvents] len: {%d}", len(aggregate.GetUncommittedEvents()))
		return nil
	}

	eventsData := make([]esdb.EventData, 0, len(aggregate.GetUncommittedEvents()))
	for _, event := range aggregate.GetUncommittedEvents() {
		eventsData = append(eventsData, event.ToEventData())
	}

	// check if new aggregate, version is 0, or now events applied or version is not greater than uncommitted events
	var expectedRevision esdb.ExpectedRevision
	if aggregate.GetVersion() == 0 || int64(len(aggregate.GetUncommittedEvents()))-aggregate.GetVersion()-1 == 0 || (aggregate.IsWithAppliedEvents() && len(aggregate.GetAppliedEvents()) == 0) {
		expectedRevision = esdb.NoStream{}
		as.log.Debugf("(Save) expectedRevision: {%T}", expectedRevision)
		span.LogFields(log.String("expectedRevision", "NoStream"))

		appendStream, err := as.esdbClient.AppendToStream(
			ctx,
			aggregate.GetID(),
			esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
			eventsData...,
		)
		if err != nil {
			// skip tracing error if wrong expected version, retry will be done
			if !es.IsEventStoreErrorCodeWrongExpectedVersion(err) {
				tracing.TraceErr(span, err)
			}
			as.log.Errorf("(Save) esdbClient.AppendToStream: {%+v}", err)
			return errors.Wrap(err, "esdbClient.AppendToStream")
		}

		as.log.Debugf("(Save) stream: {%+v}", appendStream)
		return nil
	}

	readOps := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.End{}}
	stream, err := as.esdbClient.ReadStream(context.Background(), aggregate.GetID(), readOps, 1)
	if err != nil {
		tracing.TraceErr(span, err)
		as.log.Errorf("(Save) esdbClient.ReadStream: {%+v}", err)
		return errors.Wrap(err, "esdbClient.ReadStream")
	}
	defer stream.Close()

	lastEvent, err := stream.Recv()
	if err != nil {
		tracing.TraceErr(span, err)
		as.log.Errorf("(Save) stream.Recv: {%+v}", err)
		return errors.Wrap(err, "stream.Recv")
	}

	expectedRevision = esdb.Revision(lastEvent.OriginalEvent().EventNumber)
	span.LogFields(log.Object("expectedRevision", expectedRevision))
	as.log.Debugf("(Save) expectedRevision: {%T}", expectedRevision)

	appendStream, err := as.esdbClient.AppendToStream(
		ctx,
		aggregate.GetID(),
		esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
		eventsData...,
	)
	if err != nil {
		// skip tracing error if wrong expected version, retry will be done
		if !es.IsEventStoreErrorCodeWrongExpectedVersion(err) {
			tracing.TraceErr(span, err)
		}
		as.log.Errorf("(Save) esdbClient.AppendToStream: {%+v}", err)
		return err
	}

	as.log.Debugf("(Save) stream: {%+v}", appendStream)
	aggregate.ClearUncommittedEvents()
	return nil
}

func (as *aggregateStore) Exists(ctx context.Context, aggregateID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AggregateStore.Exists")
	defer span.Finish()
	span.SetTag("aggregateID", aggregateID)

	readStreamOptions := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.Revision(1)}

	stream, err := as.esdbClient.ReadStream(ctx, aggregateID, readStreamOptions, 1)
	if err != nil {
		return errors.Wrap(err, "esdbClient.ReadStream")
	}
	defer stream.Close()

	for {
		_, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			span.LogFields(log.String("exists", "false"))
			if es.IsEventStoreErrorCodeResourceNotFound(err) {
				as.log.Warnf("(AggregateStore.Exists) esdbClient.ReadStream: {%+v}", err)
				return es.ErrAggregateNotFound
			}
			tracing.TraceErr(span, err)
			as.log.Errorf("(AggregateStore.Exists) esdbClient.ReadStream: {%+v}", err)
			return errors.Wrap(err, "stream.Recv")
		}
	}

	span.LogFields(log.String("exists", "true"))
	return nil
}
