package graph

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type ContactEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewContactEventHandler(log logger.Logger, repositories *repository.Repositories) *ContactEventHandler {
	return &ContactEventHandler{
		log:          log,
		repositories: repositories,
	}
}

func (h *ContactEventHandler) OnContactCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnContactCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)

	session := utils.NewNeo4jWriteSession(ctx, *h.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var err error

		data := neo4jrepository.ContactCreateFields{
			FirstName:       eventData.FirstName,
			LastName:        eventData.LastName,
			Prefix:          eventData.Prefix,
			Description:     eventData.Description,
			Timezone:        eventData.Timezone,
			ProfilePhotoUrl: eventData.ProfilePhotoUrl,
			Name:            eventData.Name,
			CreatedAt:       eventData.CreatedAt,
			UpdatedAt:       eventData.UpdatedAt,
			SourceFields: neo4jmodel.Source{
				Source:        helper.GetSource(eventData.Source),
				SourceOfTruth: helper.GetSourceOfTruth(eventData.SourceOfTruth),
				AppSource:     helper.GetAppSource(eventData.AppSource),
			},
		}
		err = h.repositories.Neo4jRepositories.ContactWriteRepository.CreateContactInTx(ctx, tx, eventData.Tenant, contactId, data)
		if err != nil {
			h.log.Errorf("Error while saving contact %s: %s", contactId, err.Error())
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
			err = h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, contactId, neo4jutil.NodeLabelContact, externalSystemData)
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
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.ContactUpdateFields{
		FirstName:             eventData.FirstName,
		LastName:              eventData.LastName,
		Prefix:                eventData.Prefix,
		Description:           eventData.Description,
		Timezone:              eventData.Timezone,
		ProfilePhotoUrl:       eventData.ProfilePhotoUrl,
		Name:                  eventData.Name,
		UpdatedAt:             eventData.UpdatedAt,
		Source:                eventData.Source,
		UpdateFirstName:       eventData.UpdateFirstName(),
		UpdateLastName:        eventData.UpdateLastName(),
		UpdateName:            eventData.UpdateName(),
		UpdatePrefix:          eventData.UpdatePrefix(),
		UpdateDescription:     eventData.UpdateDescription(),
		UpdateTimezone:        eventData.UpdateTimezone(),
		UpdateProfilePhotoUrl: eventData.UpdateProfilePhotoUrl(),
	}
	err := h.repositories.Neo4jRepositories.ContactWriteRepository.UpdateContact(ctx, eventData.Tenant, contactId, data)
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
				externalSystemData := neo4jmodel.ExternalSystem{
					ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
					ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
					ExternalId:       eventData.ExternalSystem.ExternalId,
					ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
					ExternalSource:   eventData.ExternalSystem.ExternalSource,
					SyncDate:         eventData.ExternalSystem.SyncDate,
				}
				innerErr := h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, contactId, neo4jutil.NodeLabelContact, externalSystemData)
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

	return nil
}

func (e *ContactEventHandler) OnPhoneNumberLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnPhoneNumberLinkToContact")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := e.repositories.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.PhoneNumberId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *ContactEventHandler) OnEmailLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnEmailLinkToContact")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.EmailWriteRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.EmailId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *ContactEventHandler) OnLocationLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnLocationLinkToContact")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactLinkLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.LocationWriteRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.LocationId, eventData.UpdatedAt)

	return err
}

func (h *ContactEventHandler) OnContactLinkToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnContactLinkToOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContactLinkWithOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.JobRoleCreateFields{
		Description: eventData.Description,
		JobTitle:    eventData.JobTitle,
		Primary:     eventData.Primary,
		CreatedAt:   eventData.CreatedAt,
		UpdatedAt:   eventData.UpdatedAt,
		StartedAt:   eventData.StartedAt,
		EndedAt:     eventData.EndedAt,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.SourceFields.Source),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.SourceFields.SourceOfTruth),
			AppSource:     helper.GetAppSource(eventData.SourceFields.AppSource),
		},
	}
	err := h.repositories.Neo4jRepositories.JobRoleWriteRepository.LinkContactWithOrganization(ctx, eventData.Tenant, contactId, eventData.OrganizationId, data)

	return err
}

func (h *ContactEventHandler) OnSocialAddedToContactV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnSocialAddedToContactV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.AddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.SocialFields{
		SocialId:  eventData.SocialId,
		Url:       eventData.Url,
		CreatedAt: eventData.CreatedAt,
		UpdatedAt: eventData.CreatedAt,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source.Source),
			SourceOfTruth: helper.GetSource(eventData.Source.Source),
			AppSource:     helper.GetSource(eventData.Source.AppSource),
		},
	}
	err := h.repositories.Neo4jRepositories.SocialWriteRepository.MergeSocialFor(ctx, eventData.Tenant, contactId, neo4jutil.NodeLabelContact, data)

	return err
}
