package graph

import (
	"context"
	"fmt"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type CommentEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewCommentEventHandler(log logger.Logger, repositories *repository.Repositories) *CommentEventHandler {
	return &CommentEventHandler{
		log:          log,
		repositories: repositories,
	}
}

func (h *CommentEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.CommentCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	if eventData.CommentedIssueId != "" {
		issueExists, err := h.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, eventData.Tenant, eventData.CommentedIssueId, neo4jentity.NodeLabel_Issue)
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
	data := neo4jrepository.CommentCreateFields{
		Content:          eventData.Content,
		ContentType:      eventData.ContentType,
		CreatedAt:        eventData.CreatedAt,
		UpdatedAt:        eventData.UpdatedAt,
		AuthorUserId:     eventData.AuthorUserId,
		CommentedIssueId: eventData.CommentedIssueId,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.Source),
			AppSource:     helper.GetAppSource(eventData.AppSource),
		},
	}
	err := h.repositories.Neo4jRepositories.CommentWriteRepository.Create(ctx, eventData.Tenant, commentId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving comment %s: %s", commentId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, commentId, neo4jentity.NodeLabel_Comment, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link comment %s with external system %s: %s", commentId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	return nil
}

func (h *CommentEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.CommentUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	commentId := aggregate.GetCommentObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.CommentUpdateFields{
		Content:     eventData.Content,
		ContentType: eventData.ContentType,
		UpdatedAt:   eventData.UpdatedAt,
		Source:      helper.GetSource(eventData.Source),
	}
	err := h.repositories.Neo4jRepositories.CommentWriteRepository.Update(ctx, eventData.Tenant, commentId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving comment %s: %s", commentId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, commentId, neo4jentity.NodeLabel_Comment, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link comment %s with external system %s: %s", commentId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	return nil
}
