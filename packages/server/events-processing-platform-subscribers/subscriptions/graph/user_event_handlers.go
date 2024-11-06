package graph

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type UserEventHandler struct {
	log      logger.Logger
	services *service.Services
}

func NewUserEventHandler(log logger.Logger, services *service.Services) *UserEventHandler {
	return &UserEventHandler{
		log:      log,
		services: services,
	}
}

func (h *UserEventHandler) OnUserCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnUserCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)

	session := utils.NewNeo4jWriteSession(ctx, *h.services.CommonServices.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var err error

		err = h.services.CommonServices.Neo4jRepositories.UserWriteRepository.CreateUserInTx(ctx, tx, eventData.Tenant, neo4jentity.UserEntity{
			Id:              userId,
			Name:            eventData.Name,
			FirstName:       eventData.FirstName,
			LastName:        eventData.LastName,
			CreatedAt:       eventData.CreatedAt,
			UpdatedAt:       eventData.UpdatedAt,
			Internal:        eventData.Internal,
			Test:            eventData.Test,
			Bot:             eventData.Bot,
			ProfilePhotoUrl: eventData.ProfilePhotoUrl,
			Timezone:        eventData.Timezone,
			Source:          neo4jentity.DecodeDataSource(eventData.SourceFields.Source),
			SourceOfTruth:   neo4jentity.DecodeDataSource(eventData.SourceFields.SourceOfTruth),
			AppSource:       helper.GetAppSource(eventData.SourceFields.AppSource),
		})
		if err != nil {
			h.log.Errorf("Error while saving user %s: %s", userId, err.Error())
			return nil, err
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
			err = h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, userId, model.NodeLabelUser, externalSystemData)
			if err != nil {
				h.log.Errorf("Error while link user %s with external system %s: %s", userId, eventData.ExternalSystem.ExternalSystemId, err.Error())
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (h *UserEventHandler) OnUserUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnUserUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.UserUpdateFields{
		Name:            eventData.Name,
		Source:          helper.GetSource(eventData.Source),
		FirstName:       eventData.FirstName,
		LastName:        eventData.LastName,
		ProfilePhotoUrl: eventData.ProfilePhotoUrl,
		Timezone:        eventData.Timezone,
	}
	err := h.services.CommonServices.Neo4jRepositories.UserWriteRepository.UpdateUser(ctx, eventData.Tenant, userId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving user %s: %s", userId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		session := utils.NewNeo4jWriteSession(ctx, *h.services.CommonServices.Neo4jRepositories.Neo4jDriver)
		defer session.Close(ctx)

		_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			//var err error
			if eventData.ExternalSystem.Available() {
				externalSystemData := neo4jmodel.ExternalSystem{
					ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
					ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
					ExternalId:       eventData.ExternalSystem.ExternalId,
					ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
					ExternalSource:   eventData.ExternalSystem.ExternalSource,
					SyncDate:         eventData.ExternalSystem.SyncDate,
				}
				innerErr := h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, userId, model.NodeLabelUser, externalSystemData)
				if innerErr != nil {
					h.log.Errorf("Error while link user %s with external system %s: %s", userId, eventData.ExternalSystem.ExternalSystemId, err.Error())
					return nil, innerErr
				}
			}
			return nil, nil
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (h *UserEventHandler) OnJobRoleLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnJobRoleLinkedToUser")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserLinkJobRoleEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.JobRoleWriteRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.JobRoleId)

	return err
}

func (h *UserEventHandler) OnPhoneNumberLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnPhoneNumberLinkedToUser")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.PhoneNumberId, eventData.Label, eventData.Primary)

	return err
}

func (h *UserEventHandler) OnAddRole(c context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(c, "UserEventHandler.OnAddRole")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserAddRoleEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: eventData.Tenant})

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.UserWriteRepository.AddRole(ctx, userId, eventData.Role)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while adding role %s to user %s: %s", eventData.Role, userId, err.Error())
	}

	return err
}

func (h *UserEventHandler) OnRemoveRole(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnRemoveRole")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserRemoveRoleEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.UserWriteRepository.RemoveRole(ctx, eventData.Tenant, userId, eventData.Role)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while removing role %s from user %s: %s", eventData.Role, userId, err.Error())
	}

	return err
}
