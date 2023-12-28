package graph

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type UserEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewUserEventHandler(log logger.Logger, repositories *repository.Repositories) *UserEventHandler {
	return &UserEventHandler{
		log:          log,
		repositories: repositories,
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

	session := utils.NewNeo4jWriteSession(ctx, *h.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var err error
		userCreateData := neo4jrepository.UserCreateFields{
			Name:      eventData.Name,
			FirstName: eventData.FirstName,
			LastName:  eventData.LastName,
			SourceFields: neo4jmodel.Source{
				Source:        helper.GetSource(eventData.SourceFields.Source),
				SourceOfTruth: helper.GetSourceOfTruth(eventData.SourceFields.SourceOfTruth),
				AppSource:     helper.GetAppSource(eventData.SourceFields.AppSource),
			},
			CreatedAt:       eventData.CreatedAt,
			UpdatedAt:       eventData.UpdatedAt,
			Internal:        eventData.Internal,
			Bot:             eventData.Bot,
			ProfilePhotoUrl: eventData.ProfilePhotoUrl,
			Timezone:        eventData.Timezone,
		}
		err = h.repositories.Neo4jRepositories.UserWriteRepository.CreateUserInTx(ctx, tx, eventData.Tenant, userId, userCreateData)
		if err != nil {
			h.log.Errorf("Error while saving user %s: %s", userId, err.Error())
			return nil, err
		}
		if eventData.ExternalSystem.Available() {
			err = h.repositories.ExternalSystemRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, userId, neo4jentity.NodeLabel_User, eventData.ExternalSystem)
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
		UpdatedAt:       eventData.UpdatedAt,
		Internal:        eventData.Internal,
		Bot:             eventData.Bot,
		ProfilePhotoUrl: eventData.ProfilePhotoUrl,
		Timezone:        eventData.Timezone,
	}
	err := h.repositories.Neo4jRepositories.UserWriteRepository.UpdateUser(ctx, eventData.Tenant, userId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving user %s: %s", userId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		session := utils.NewNeo4jWriteSession(ctx, *h.repositories.Drivers.Neo4jDriver)
		defer session.Close(ctx)

		_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			//var err error
			if eventData.ExternalSystem.Available() {
				innerErr := h.repositories.ExternalSystemRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, userId, neo4jentity.NodeLabel_User, eventData.ExternalSystem)
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
	err := h.repositories.JobRoleRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.JobRoleId, eventData.UpdatedAt)

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
	err := h.repositories.PhoneNumberRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.PhoneNumberId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *UserEventHandler) OnEmailLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnEmailLinkedToUser")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.EmailRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.EmailId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *UserEventHandler) OnAddPlayer(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnAddPlayer")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserAddPlayerInfoEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.PlayerRepository.Merge(ctx, eventData.Tenant, userId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while adding player %s to user %s: %s", eventData.AuthId, userId, err.Error())
	}

	return err
}

func (h *UserEventHandler) OnAddRole(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEventHandler.OnAddRole")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UserAddRoleEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.UserWriteRepository.AddRole(ctx, eventData.Tenant, userId, eventData.Role, eventData.At)
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
	err := h.repositories.Neo4jRepositories.UserWriteRepository.RemoveRole(ctx, eventData.Tenant, userId, eventData.Role, eventData.At)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while removing role %s from user %s: %s", eventData.Role, userId, err.Error())
	}

	return err
}
