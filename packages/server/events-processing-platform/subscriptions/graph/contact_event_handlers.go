package graph

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ContactEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func (h *ContactEventHandler) OnContactCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnContactCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.ContactCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)

	session := utils.NewNeo4jWriteSession(ctx, *h.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var err error
		err = h.repositories.ContactRepository.CreateContactInTx(ctx, tx, contactId, eventData)
		if err != nil {
			h.log.Errorf("Error while saving contact %s: %s", contactId, err.Error())
			return nil, err
		}
		if eventData.ExternalSystem.Available() {
			err = h.repositories.ExternalSystemRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, contactId, constants.NodeLabel_Contact, eventData.ExternalSystem)
			if err != nil {
				h.log.Errorf("Error while link contact %s with external system %s: %s", contactId, eventData.ExternalSystem.ExternalSystemId, err.Error())
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

func (h *ContactEventHandler) OnContactUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnContactUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.ContactUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.ContactRepository.UpdateContact(ctx, contactId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving contact %s: %s", contactId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		session := utils.NewNeo4jWriteSession(ctx, *h.repositories.Drivers.Neo4jDriver)
		defer session.Close(ctx)

		_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			//var err error
			if eventData.ExternalSystem.Available() {
				innerErr := h.repositories.ExternalSystemRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, contactId, constants.NodeLabel_Contact, eventData.ExternalSystem)
				if innerErr != nil {
					h.log.Errorf("Error while link contact %s with external system %s: %s", contactId, eventData.ExternalSystem.ExternalSystemId, err.Error())
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

	return err
}

func (e *ContactEventHandler) OnPhoneNumberLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnPhoneNumberLinkToContact")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.ContactLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := e.repositories.PhoneNumberRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.PhoneNumberId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *ContactEventHandler) OnEmailLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnEmailLinkToContact")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.ContactLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.EmailRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.EmailId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *ContactEventHandler) OnLocationLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnLocationLinkToContact")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.ContactLinkLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.LocationRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.LocationId, eventData.UpdatedAt)

	return err
}
