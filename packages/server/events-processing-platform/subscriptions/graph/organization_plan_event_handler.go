package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	orgModel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
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
	err := h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.Create(ctx, eventData.Tenant, masterPlanId, eventData.OrganizationPlanId, eventData.Name, source, appSource, eventData.CreatedAt, entity.OrganizationPlanStatusDetails{Status: model.NotStarted.String(), UpdatedAt: eventData.CreatedAt, Comments: ""})
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

	// if plan status changed, propagate to organization
	if data.UpdateStatusDetails {
		err = h.propagateStatusToOrg(ctx, eventData.Tenant, eventData.OrganizationPlanId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while propagating status to organization for organization plan %s: %s", eventData.OrganizationPlanId, err.Error())
			return err
		}
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
	err := h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.CreateMilestone(
		ctx,
		eventData.Tenant,
		eventData.OrganizationPlanId,
		eventData.MilestoneId,
		eventData.Name,
		source,
		appSource,
		eventData.Order,
		convertItemsStrToObject(eventData.Items),
		eventData.Optional,
		eventData.Adhoc,
		eventData.CreatedAt,
		eventData.DueDate,
		entity.OrganizationPlanMilestoneStatusDetails{
			Status:    model.MilestoneNotStarted.String(),
			UpdatedAt: eventData.CreatedAt,
			Comments:  "",
		})
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

	dueDate, err := h.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetMilestoneDueDate(ctx, eventData.Tenant, eventData.MilestoneId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while retrieving milestone due date %s: %s", eventData.MilestoneId, err.Error())
		return err
	}

	if eventData.UpdateDueDate() {
		dueDate = eventData.DueDate
	}

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
		Adhoc:               eventData.Adhoc,
		UpdateName:          eventData.UpdateName(),
		UpdateOrder:         eventData.UpdateOrder(),
		UpdateDueDate:       eventData.UpdateDueDate(),
		UpdateItems:         eventData.UpdateItems(),
		UpdateOptional:      eventData.UpdateOptional(),
		UpdateRetired:       eventData.UpdateRetired(),
		UpdateStatusDetails: eventData.UpdateStatusDetails(),
		UpdateAdhoc:         eventData.UpdateAdhoc(),
	}

	// check if milestone status should update
	if eventData.UpdateItems() {
		if checkAllItemStatusesAreDoneOrSkipped(eventData.Items) {
			if eventData.UpdatedAt.After(dueDate) {
				data.StatusDetails.Status = model.MilestoneDoneLate.String()
			} else {
				data.StatusDetails.Status = model.MilestoneDone.String()
			}
		} else {
			if eventData.UpdatedAt.After(dueDate) {
				data.StatusDetails.Status = model.MilestoneStartedLate.String()
			} else {
				data.StatusDetails.Status = model.MilestoneStarted.String()
			}
		}
		data.StatusDetails.UpdatedAt = eventData.UpdatedAt
		data.StatusDetails.Comments = eventData.StatusDetails.Comments
		data.UpdateStatusDetails = true
	}

	err = h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.UpdateMilestone(ctx, eventData.Tenant, eventData.OrganizationPlanId, eventData.MilestoneId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating master plan milestone %s: %s", eventData.MilestoneId, err.Error())
		return err
	}
	// if milestone status changed, propagate to organization plan.
	if data.UpdateStatusDetails {
		h.propagateStatusUpdatesFromMilestone(ctx, data, eventData.Tenant, eventData.OrganizationPlanId)
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

func (h *OrganizationPlanEventHandler) propagateStatusUpdatesFromMilestone(ctx context.Context, data neo4jrepository.OrganizationPlanMilestoneUpdateFields, tenant, opid string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanEventHandler.propagateStatusUpdatesFromMilestone")
	defer span.Finish()
	span.SetTag(tracing.SpanTagEntityId, opid)

	opmNode, err := h.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetMilestonesForOrganizationPlan(ctx, tenant, opid)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while retrieving organization plan %s: %s", opid, err.Error())
		return err
	}

	opNode, err := h.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetOrganizationPlanById(ctx, tenant, opid)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while retrieving organization plan %s: %s", opid, err.Error())
		return err
	}

	op := neo4jmapper.MapDbNodeToOrganizationPlanEntity(opNode)

	opMilestones := convertMilestonesToOrganizationPlanMilestones(opmNode)

	opdata := neo4jrepository.OrganizationPlanUpdateFields{
		Name:    op.Name,
		Retired: op.Retired,
		StatusDetails: entity.OrganizationPlanStatusDetails{
			Status:    op.StatusDetails.Status,
			UpdatedAt: op.StatusDetails.UpdatedAt,
			Comments:  op.StatusDetails.Comments,
		},
		UpdatedAt:           op.UpdatedAt,
		UpdateName:          false,
		UpdateRetired:       false,
		UpdateStatusDetails: false,
	}

	// check if all milestones are done
	allMilestonesDone := true
	late := false
	started := false
	for _, milestone := range opMilestones {
		if milestone.StatusDetails.Status != model.MilestoneDone.String() && milestone.StatusDetails.Status != model.MilestoneDoneLate.String() {
			allMilestonesDone = false
		}
		if milestone.StatusDetails.Status == model.MilestoneDoneLate.String() || milestone.StatusDetails.Status == model.MilestoneStartedLate.String() || milestone.StatusDetails.Status == model.MilestoneNotStartedLate.String() {
			late = true
		}
		if milestone.StatusDetails.Status == model.MilestoneStarted.String() || milestone.StatusDetails.Status == model.MilestoneStartedLate.String() {
			started = true
		}
	}

	// update organization plan status
	if allMilestonesDone {
		if late {
			opdata.StatusDetails.Status = model.DoneLate.String()
		} else {
			opdata.StatusDetails.Status = model.Done.String()
		}
	} else {
		if !started {
			if late {
				opdata.StatusDetails.Status = model.NotStartedLate.String()
			} else {
				opdata.StatusDetails.Status = model.NotStarted.String()
			}
		} else {
			if late {
				opdata.StatusDetails.Status = model.Late.String()
			} else {
				opdata.StatusDetails.Status = model.OnTrack.String()
			}
		}
	}

	if op.StatusDetails.Status != opdata.StatusDetails.Status {
		opdata.StatusDetails.UpdatedAt = time.Now().UTC()
		opdata.UpdateStatusDetails = true
		err = h.repositories.Neo4jRepositories.OrganizationPlanWriteRepository.Update(ctx, tenant, opid, opdata)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while updating organization plan %s: %s", opid, err.Error())
			return err
		}

		// if plan status changed, propagate to organization
		err = h.propagateStatusToOrg(ctx, tenant, opid)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while propagating status to organization for organization plan %s: %s", opid, err.Error())
			return err
		}
	}

	return nil
}

func (h *OrganizationPlanEventHandler) propagateStatusToOrg(ctx context.Context, tenant, opid string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanEventHandler.propagateStatusToOrg")
	defer span.Finish()
	span.SetTag(tracing.SpanTagEntityId, opid)

	// propagate to organization
	orgNode, err := h.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetOrganizationFromOrganizationPlan(ctx, tenant, opid)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while retrieving organization for organization plan %s: %s", opid, err.Error())
		return err
	}
	org := neo4jmapper.MapDbNodeToOrganizationEntity(orgNode)

	// get all organization plans
	opNodes, err := h.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetOrganizationPlansForOrganization(ctx, tenant, org.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while retrieving organization plans for organization %s: %s", org.ID, err.Error())
		return err
	}

	// check if all organization plans are done
	allPlansDone := true
	late := false
	started := false
	for _, opNode := range opNodes {
		op := neo4jmapper.MapDbNodeToOrganizationPlanEntity(opNode)
		if op.StatusDetails.Status != model.Done.String() && op.StatusDetails.Status != model.DoneLate.String() {
			allPlansDone = false
		}
		if op.StatusDetails.Status == model.DoneLate.String() || op.StatusDetails.Status == model.NotStartedLate.String() || op.StatusDetails.Status == model.Late.String() {
			late = true
		}
		if op.StatusDetails.Status != model.NotStarted.String() || op.StatusDetails.Status != model.NotStartedLate.String() {
			started = true
		}
	}

	var statusStr string
	updatedAtNow := time.Now().UTC()

	// update organization status
	if allPlansDone {
		statusStr = orgModel.Done.String()
	} else if late {
		statusStr = orgModel.Late.String()
	} else if !started {
		statusStr = orgModel.NotStarted.String()
	} else {
		statusStr = orgModel.OnTrack.String()
	}

	// if status changed, write to Org DB Node and save change action
	if org.OnboardingDetails.Status != statusStr {
		err = h.repositories.Neo4jRepositories.OrganizationWriteRepository.UpdateOnboardingStatus(ctx, tenant, org.ID, statusStr, org.OnboardingDetails.Comments, getOrderForOnboardingStatus(statusStr), updatedAtNow)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to update onboarding status for organization %s: %s", org.ID, err.Error())
			return err
		}

		err = h.saveOnboardingStatusChangeAction(ctx, tenant, org.ID, statusStr, org.OnboardingDetails.Comments, span, updatedAtNow)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to save onboarding status change action for organization %s: %s", org.ID, err.Error())
		}
	}

	return nil
}

func (h *OrganizationPlanEventHandler) saveOnboardingStatusChangeAction(ctx context.Context, tenant, organizationId, status, comments string, span opentracing.Span, updatedAt time.Time) error {
	metadata, _ := utils.ToJson(ActionOnboardingStatusMetadata{
		Status:     status,
		Comments:   comments,
		UserId:     "",
		ContractId: "",
	})

	message := fmt.Sprintf("The onboarding status was automatically set to %s", onboardingStatusReadableStringForActionMessage(status))

	extraActionProperties := map[string]interface{}{
		"status":   status,
		"comments": comments,
	}
	_, err := h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, tenant, organizationId, neo4jenum.ORGANIZATION, neo4jenum.ActionOnboardingStatusChanged, message, metadata, updatedAt, extraActionProperties)
	return err
}

/////////////////////////////////////////////////// helper functions ///////////////////////////////////////////////////

func checkAllItemStatusesAreDoneOrSkipped(items []model.OrganizationPlanMilestoneItem) bool {
	for _, item := range items {
		if item.Status == model.TaskNotDone.String() || item.Status == model.TaskNotDoneLate.String() {
			return false
		}
	}
	return true
}

func convertMilestonesToOrganizationPlanMilestones(opMilestonesNodes []*dbtype.Node) []entity.OrganizationPlanMilestoneEntity {
	organizationPlanMilestones := make([]entity.OrganizationPlanMilestoneEntity, len(opMilestonesNodes))
	for i, opMilestoneNode := range opMilestonesNodes {
		milestone := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(opMilestoneNode)
		organizationPlanMilestones[i] = *milestone
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
