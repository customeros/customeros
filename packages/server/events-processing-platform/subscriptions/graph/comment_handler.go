package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphCommentEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func (h *GraphCommentEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphCommentEventHandler.OnCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event.CommentCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	if eventData.CommentedIssueId != "" {
		issueExists, err := h.repositories.IssueRepository.ExistsById(ctx, eventData.Tenant, eventData.CommentedIssueId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while checking if issue %s exists: %s", eventData.CommentedIssueId, err.Error())
			return err
		}
		if !issueExists {
			err := errors.New(fmt.Sprintf("commented issue %s does not exist", eventData.CommentedIssueId))
			tracing.TraceErr(span, err)
			h.log.Errorf("Issue %s does not exist", eventData.CommentedIssueId)
			return err
		}
	}

	commentId := aggregate.GetCommentObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.CommentRepository.Create(ctx, eventData.Tenant, commentId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving comment %s: %s", commentId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, commentId, constants.NodeLabel_Comment, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link comment %s with external system %s: %s", commentId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	return nil
}

func (h *GraphCommentEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphCommentEventHandler.OnCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event.CommentUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	commentId := aggregate.GetCommentObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.CommentRepository.Update(ctx, eventData.Tenant, commentId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving comment %s: %s", commentId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, commentId, constants.NodeLabel_Comment, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link comment %s with external system %s: %s", commentId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	return nil
}
