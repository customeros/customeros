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

func (aggregateStore *aggregateStore) Load(ctx context.Context, aggregate es.Aggregate) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Load")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", aggregate.GetID()))

	stream, err := aggregateStore.esdbClient.ReadStream(ctx, aggregate.GetID(), esdb.ReadStreamOptions{}, count)
	if err != nil {
		tracing.TraceErr(span, err)
		aggregateStore.log.Errorf("(Load) esdbClient.ReadStream: {%+v}", err)
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
			aggregateStore.log.Errorf("(Load) esdbClient.ReadStream: {%+v}", err)
			return errors.Wrap(err, "stream.Recv")
		}

		esEvent := es.NewEventFromRecorded(event.Event)
		if err := aggregate.RaiseEvent(esEvent); err != nil {
			tracing.TraceErr(span, err)
			aggregateStore.log.Errorf("(Load) aggregate.RaiseEvent: {%+v}", err)
			return errors.Wrap(err, "RaiseEvent")
		}
		aggregateStore.log.Debugf("(Load) esEvent: {%s}", esEvent.String())
	}

	aggregateStore.log.Debugf("(Load) aggregate: {%s}", aggregate.String())
	return nil
}

func (aggregateStore *aggregateStore) Save(ctx context.Context, aggregate es.Aggregate) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Save")
	defer span.Finish()
	span.LogFields(log.String("aggregate", aggregate.String()))

	if len(aggregate.GetUncommittedEvents()) == 0 {
		aggregateStore.log.Debugf("(Save) [no uncommittedEvents] len: {%d}", len(aggregate.GetUncommittedEvents()))
		return nil
	}

	eventsData := make([]esdb.EventData, 0, len(aggregate.GetUncommittedEvents()))
	for _, event := range aggregate.GetUncommittedEvents() {
		eventsData = append(eventsData, event.ToEventData())
	}

	// check for aggregate.GetVersion() == 0 or len(aggregate.GetAppliedEvents()) == 0 means new aggregate
	var expectedRevision esdb.ExpectedRevision
	if aggregate.GetVersion() == 0 {
		expectedRevision = esdb.NoStream{}
		aggregateStore.log.Debugf("(Save) expectedRevision: {%T}", expectedRevision)

		appendStream, err := aggregateStore.esdbClient.AppendToStream(
			ctx,
			aggregate.GetID(),
			esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
			eventsData...,
		)
		if err != nil {
			tracing.TraceErr(span, err)
			aggregateStore.log.Errorf("(Save) esdbClient.AppendToStream: {%+v}", err)
			return errors.Wrap(err, "esdbClient.AppendToStream")
		}

		aggregateStore.log.Debugf("(Save) stream: {%+v}", appendStream)
		return nil
	}

	readOps := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.End{}}
	stream, err := aggregateStore.esdbClient.ReadStream(context.Background(), aggregate.GetID(), readOps, 1)
	if err != nil {
		tracing.TraceErr(span, err)
		aggregateStore.log.Errorf("(Save) esdbClient.ReadStream: {%+v}", err)
		return errors.Wrap(err, "esdbClient.ReadStream")
	}
	defer stream.Close()

	lastEvent, err := stream.Recv()
	if err != nil {
		tracing.TraceErr(span, err)
		aggregateStore.log.Errorf("(Save) stream.Recv: {%+v}", err)
		return errors.Wrap(err, "stream.Recv")
	}

	expectedRevision = esdb.Revision(lastEvent.OriginalEvent().EventNumber)
	aggregateStore.log.Debugf("(Save) expectedRevision: {%T}", expectedRevision)

	appendStream, err := aggregateStore.esdbClient.AppendToStream(
		ctx,
		aggregate.GetID(),
		esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
		eventsData...,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		aggregateStore.log.Errorf("(Save) esdbClient.AppendToStream: {%+v}", err)
		return errors.Wrap(err, "esdbClient.AppendToStream")
	}

	aggregateStore.log.Debugf("(Save) stream: {%+v}", appendStream)
	aggregate.ClearUncommittedEvents()
	return nil
}

func (aggregateStore *aggregateStore) Exists(ctx context.Context, streamID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Exists")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", streamID))

	readStreamOptions := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.Revision(1)}

	stream, err := aggregateStore.esdbClient.ReadStream(ctx, streamID, readStreamOptions, 1)
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
			tracing.TraceErr(span, err)
			if es.IsErrEsResourceNotFound(err) {
				aggregateStore.log.Warnf("(aggregateStore.Exists) esdbClient.ReadStream: {%+v}", err)
				return es.ErrAggregateNotFound
			}
			aggregateStore.log.Errorf("(aggregateStore.Exists) esdbClient.ReadStream: {%+v}", err)
			return errors.Wrap(err, "stream.Recv")
		}
	}

	return nil
}
