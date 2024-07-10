package graph

import (
	"context"
	"fmt"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
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

	var eventData comment.CommentCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	if eventData.CommentedIssueId != "" {
		issueExists, err := h.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, eventData.Tenant, eventData.CommentedIssueId, neo4jutil.NodeLabelIssue)
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

	commentId := comment.GetCommentObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.CommentCreateFields{
		Content:          eventData.Content,
		ContentType:      eventData.ContentType,
		CreatedAt:        eventData.CreatedAt,
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
		externalSystemData := neo4jmodel.ExternalSystem{
			ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
			ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
			ExternalId:       eventData.ExternalSystem.ExternalId,
			ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
			ExternalSource:   eventData.ExternalSystem.ExternalSource,
			SyncDate:         eventData.ExternalSystem.SyncDate,
		}
		err = h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, commentId, neo4jutil.NodeLabelComment, externalSystemData)
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

	var eventData comment.CommentUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	commentId := comment.GetCommentObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.CommentUpdateFields{
		Content:     eventData.Content,
		ContentType: eventData.ContentType,
		Source:      helper.GetSource(eventData.Source),
	}
	err := h.repositories.Neo4jRepositories.CommentWriteRepository.Update(ctx, eventData.Tenant, commentId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving comment %s: %s", commentId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		externalSystemData := neo4jmodel.ExternalSystem{
			ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
			ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
			ExternalId:       eventData.ExternalSystem.ExternalId,
			ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
			ExternalSource:   eventData.ExternalSystem.ExternalSource,
			SyncDate:         eventData.ExternalSystem.SyncDate,
		}
		err = h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, commentId, neo4jutil.NodeLabelComment, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link comment %s with external system %s: %s", commentId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	return nil
}
