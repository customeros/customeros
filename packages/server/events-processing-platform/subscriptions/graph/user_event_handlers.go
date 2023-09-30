package graph

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphUserEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func (h *GraphUserEventHandler) OnUserCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphUserEventHandler.OnUserCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

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
		err = h.repositories.UserRepository.CreateUserInTx(ctx, tx, userId, eventData)
		if err != nil {
			h.log.Errorf("Error while saving user %s: %s", userId, err.Error())
			return nil, err
		}
		if eventData.ExternalSystem.Available() {
			err = h.repositories.ExternalSystemRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, userId, constants.NodeLabel_User, eventData.ExternalSystem)
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

func (h *GraphUserEventHandler) OnUserUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphUserEventHandler.OnUserUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.UserUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.UserRepository.UpdateUser(ctx, userId, eventData)
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
				innerErr := h.repositories.ExternalSystemRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, userId, constants.NodeLabel_User, eventData.ExternalSystem)
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

func (h *GraphUserEventHandler) OnJobRoleLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphUserEventHandler.OnJobRoleLinkedToUser")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.UserLinkJobRoleEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.JobRoleRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.JobRoleId, eventData.UpdatedAt)

	return err
}

func (h *GraphUserEventHandler) OnPhoneNumberLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphUserEventHandler.OnPhoneNumberLinkedToUser")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.UserLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.PhoneNumberRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.PhoneNumberId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *GraphUserEventHandler) OnEmailLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphUserEventHandler.OnEmailLinkedToUser")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.UserLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.EmailRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.EmailId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *GraphUserEventHandler) OnAddPlayer(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphUserEventHandler.OnAddPlayer")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

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
