package graph

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	event "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type OrganizationPlanEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewOrganizationPlanEventHandler(log logger.Logger, repositories *repository.Repositories) *OrganizationPlanEventHandler {
	return &OrganizationPlanEventHandler{
		log:          log,
		repositories: repositories,
	}
}

func (h *OrganizationPlanEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.OrganizationPlanCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := eventData.OrganizationId
	masterPlanId := eventData.MasterPlanId
	span.SetTag(tracing.SpanTagEntityId, eventData.OrganizationPlanId)

	source := helper.GetSource(eventData.SourceFields.Source)
	appSource := helper.GetAppSource(eventData.SourceFields.AppSource)

	// Create empty org plan
	err := h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.Create(ctx, eventData.Tenant, eventData.OrganizationPlanId, eventData.Name, source, appSource, eventData.CreatedAt, entity.OrganizationPlanStatusDetails{Status: model.NotStarted.String(), UpdatedAt: eventData.CreatedAt, Comments: ""})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving organization plan %s: %s", eventData.OrganizationPlanId, err.Error())
		return err
	}

	// Link org plan to master plan
	if masterPlanId != "" {
		err = h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.LinkWithMasterPlan(ctx, eventData.Tenant, eventData.OrganizationPlanId, masterPlanId, eventData.CreatedAt)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error linking master plan %s: %s", eventData.OrganizationPlanId, err.Error())
			return err
		}
	}
	// Link org plan to org
	err = h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.LinkWithOrganization(ctx, eventData.Tenant, eventData.OrganizationPlanId, organizationId, eventData.CreatedAt)

	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error linking organization to plan %s: %s", eventData.OrganizationPlanId, err.Error())
		return err
	}
	return err
}

func (h *OrganizationPlanEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.OrganizationPlanUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	span.SetTag(tracing.SpanTagEntityId, eventData.OrganizationPlanId)

	data := neo4jrepository.OrganizationPlanUpdateFields{
		Name:    eventData.Name,
		Retired: eventData.Retired,
		StatusDetails: entity.OrganizationPlanStatusDetails{
			Status:    eventData.StatusDetails.Status,
			UpdatedAt: eventData.StatusDetails.UpdatedAt,
			Comments:  eventData.StatusDetails.Comments,
		},
		UpdatedAt:           eventData.UpdatedAt,
		UpdateName:          eventData.UpdateName(),
		UpdateRetired:       eventData.UpdateRetired(),
		UpdateStatusDetails: eventData.UpdateStatusDetails(),
	}
	err := h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.Update(ctx, eventData.Tenant, eventData.OrganizationPlanId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating organization plan %s: %s", eventData.OrganizationPlanId, err.Error())
		return err
	}
	return err
}

func (h *OrganizationPlanEventHandler) OnCreateMilestone(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanEventHandler.OnCreateMilestone")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.OrganizationPlanMilestoneCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	span.SetTag(tracing.SpanTagEntityId, eventData.MilestoneId)

	source := helper.GetSource(eventData.SourceFields.Source)
	appSource := helper.GetAppSource(eventData.SourceFields.AppSource)
	err := h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.CreateMilestone(ctx, eventData.Tenant, eventData.OrganizationPlanId, eventData.MilestoneId,
		eventData.Name, source, appSource, eventData.Order, convertItemsStrToObject(eventData.Items), eventData.Optional, eventData.CreatedAt, eventData.DueDate, entity.OrganizationPlanMilestoneStatusDetails{Status: model.MilestoneNotStarted.String(), UpdatedAt: eventData.CreatedAt, Comments: ""})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving organization plan milestone %s: %s", eventData.OrganizationPlanId, err.Error())
		return err
	}
	return err
}

func (h *OrganizationPlanEventHandler) OnUpdateMilestone(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanEventHandler.OnUpdateMilestone")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.OrganizationPlanMilestoneUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	span.SetTag(tracing.SpanTagEntityId, eventData.MilestoneId)

	data := neo4jrepository.OrganizationPlanMilestoneUpdateFields{
		UpdatedAt: eventData.UpdatedAt,
		Name:      eventData.Name,
		Order:     eventData.Order,
		DueDate:   eventData.DueDate,
		Items:     convertItemsModelToEntity(eventData.Items),
		Optional:  eventData.Optional,
		Retired:   eventData.Retired,
		StatusDetails: entity.OrganizationPlanMilestoneStatusDetails{
			Status:    eventData.StatusDetails.Status,
			UpdatedAt: eventData.StatusDetails.UpdatedAt,
			Comments:  eventData.StatusDetails.Comments,
		},
		UpdateName:          eventData.UpdateName(),
		UpdateOrder:         eventData.UpdateOrder(),
		UpdateDueDate:       eventData.UpdateDueDate(),
		UpdateItems:         eventData.UpdateItems(),
		UpdateOptional:      eventData.UpdateOptional(),
		UpdateRetired:       eventData.UpdateRetired(),
		UpdateStatusDetails: eventData.UpdateStatusDetails(),
	}
	err := h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.UpdateMilestone(ctx, eventData.Tenant, eventData.OrganizationPlanId, eventData.MilestoneId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating master plan milestone %s: %s", eventData.MilestoneId, err.Error())
		return err
	}
	return err
}

func (h *OrganizationPlanEventHandler) OnReorderMilestones(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanEventHandler.OnReorderMilestones")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.OrganizationPlanMilestoneReorderEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationPlanId := eventData.OrganizationPlanId

	span.SetTag(tracing.SpanTagEntityId, organizationPlanId)

	for i, milestoneId := range eventData.MilestoneIds {
		data := neo4jrepository.OrganizationPlanMilestoneUpdateFields{
			Order:       int64(i),
			UpdatedAt:   eventData.UpdatedAt,
			UpdateOrder: true,
		}
		err := h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.UpdateMilestone(ctx, eventData.Tenant, organizationPlanId, milestoneId, data)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while updating organization plan milestone order %s: %s", milestoneId, err.Error())
			return err
		}
	}
	return nil
}

func convertMasterPlanMilestonesToOrganizationPlanMilestones(masterPlanMilestonesNodes []*dbtype.Node, createdAt time.Time) []entity.OrganizationPlanMilestoneEntity {
	organizationPlanMilestones := make([]entity.OrganizationPlanMilestoneEntity, len(masterPlanMilestonesNodes))
	for i, masterPlanMilestoneNode := range masterPlanMilestonesNodes {
		mpMilestone := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneNode)
		organizationPlanMilestones[i] = entity.OrganizationPlanMilestoneEntity{
			Id:            uuid.New().String(),
			Name:          mpMilestone.Name,
			Order:         mpMilestone.Order,
			DueDate:       createdAt.Add(time.Duration(mpMilestone.DurationHours) * time.Hour),
			Items:         convertItemsStrToObject(mpMilestone.Items),
			Optional:      mpMilestone.Optional,
			CreatedAt:     createdAt,
			UpdatedAt:     createdAt,
			StatusDetails: entity.OrganizationPlanMilestoneStatusDetails{Status: model.MilestoneNotStarted.String(), UpdatedAt: createdAt, Comments: ""},
		}
	}
	return organizationPlanMilestones
}

func convertItemsStrToObject(items []string) []entity.OrganizationPlanMilestoneItem {
	milestoneItems := make([]entity.OrganizationPlanMilestoneItem, len(items))
	for i, item := range items {
		milestoneItems[i] = entity.OrganizationPlanMilestoneItem{
			Text:      item,
			UpdatedAt: time.Now().UTC(),
			Status:    model.TaskNotDone.String(),
		}
	}
	return milestoneItems
}

func convertItemsModelToEntity(items []model.OrganizationPlanMilestoneItem) []entity.OrganizationPlanMilestoneItem {
	milestoneItems := make([]entity.OrganizationPlanMilestoneItem, len(items))
	for i, item := range items {
		milestoneItems[i] = entity.OrganizationPlanMilestoneItem{
			Text:      item.Text,
			UpdatedAt: time.Now().UTC(),
			Status:    item.Status,
		}
	}
	return milestoneItems
}
