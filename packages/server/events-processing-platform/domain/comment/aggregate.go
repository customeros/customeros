package comment

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/comment"
	"github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

const (
	CommentAggregateType eventstore.AggregateType = "comment"
)

type CommentAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Comment *Comment
}

func GetCommentObjectID(aggregateID string, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, CommentAggregateType)
}

func NewCommentAggregateWithTenantAndID(tenant, id string) *CommentAggregate {
	commentAggregate := CommentAggregate{}
	commentAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(CommentAggregateType, tenant, id)
	commentAggregate.SetWhen(commentAggregate.When)
	commentAggregate.Comment = &Comment{}
	commentAggregate.Tenant = tenant

	return &commentAggregate
}

func (a *CommentAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *commentpb.UpsertCommentGrpcRequest:
		return a.UpsertCommentGrpcRequest(ctx, r)
	default:
		return nil, nil
	}
}

func (a *CommentAggregate) UpsertCommentGrpcRequest(ctx context.Context, request *commentpb.UpsertCommentGrpcRequest) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CommentAggregate.UpsertCommentGrpcRequest")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	var err error
	var event eventstore.Event

	dataFields := CommentDataFields{
		Content:          request.Content,
		ContentType:      request.ContentType,
		AuthorUserId:     request.AuthorUserId,
		CommentedIssueId: request.CommentedIssueId,
	}

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)
	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), createdAtNotNil)

	if eventstore.IsAggregateNotFound(a) {
		event, err = NewCommentCreateEvent(a, dataFields, source, externalSystem, createdAtNotNil, updatedAtNotNil)
	} else {
		event, err = NewCommentUpdateEvent(a, dataFields.Content, dataFields.ContentType, source.Source, externalSystem, updatedAtNotNil)
	}

	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "CommentAggregate.UpsertCommentGrpcRequest failed to create event")
	}
	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.UserId,
		App:    source.AppSource,
	})

	return request.Id, a.Apply(event)
}

func (a *CommentAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case CommentCreateV1:
		return a.onCommentCreate(evt)
	case CommentUpdateV1:
		return a.onCommentUpdate(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *CommentAggregate) onCommentCreate(evt eventstore.Event) error {
	var eventData CommentCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Comment.ID = a.ID
	a.Comment.Tenant = a.Tenant
	a.Comment.Content = eventData.Content
	a.Comment.ContentType = eventData.ContentType
	a.Comment.AuthorUserId = eventData.AuthorUserId
	a.Comment.CommentedIssueId = eventData.CommentedIssueId
	a.Comment.Source = commonmodel.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.Source,
		AppSource:     eventData.AppSource,
	}
	a.Comment.CreatedAt = eventData.CreatedAt
	a.Comment.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.Comment.ExternalSystems = []commonmodel.ExternalSystem{eventData.ExternalSystem}
	}
	return nil
}

func (a *CommentAggregate) onCommentUpdate(evt eventstore.Event) error {
	var eventData CommentUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if eventData.Source == constants.SourceOpenline {
		a.Comment.Source.SourceOfTruth = eventData.Source
	}
	if eventData.Source != a.Comment.Source.SourceOfTruth && a.Comment.Source.SourceOfTruth == constants.SourceOpenline {
		if a.Comment.Content == "" {
			a.Comment.Content = eventData.Content
		}
		if a.Comment.ContentType == "" {
			a.Comment.ContentType = eventData.ContentType
		}
	} else {
		a.Comment.Content = eventData.Content
		a.Comment.ContentType = eventData.ContentType
	}
	a.Comment.UpdatedAt = eventData.UpdatedAt
	return nil
}
