package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type PhoneNumberEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewPhoneNumberEventHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *PhoneNumberEventHandler {
	return &PhoneNumberEventHandler{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (h *PhoneNumberEventHandler) OnPhoneNumberCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberEventHandler.OnPhoneNumberCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	data := neo4jrepository.PhoneNumberCreateFields{
		RawPhoneNumber: eventData.RawPhoneNumber,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.SourceFields.Source),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.SourceFields.SourceOfTruth),
			AppSource:     helper.GetAppSource(eventData.SourceFields.AppSource),
		},
		CreatedAt: eventData.CreatedAt,
	}
	err := h.repositories.Neo4jRepositories.PhoneNumberWriteRepository.CreatePhoneNumber(ctx, eventData.Tenant, phoneNumberId, data)

	return err
}

func (h *PhoneNumberEventHandler) OnPhoneNumberUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberEventHandler.OnPhoneNumberUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberUpdatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("rawPhoneNumber", eventData.RawPhoneNumber))

	phoneNumberDbNode, err := h.repositories.Neo4jRepositories.PhoneNumberReadRepository.GetById(ctx, eventData.Tenant, phoneNumberId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	phoneNumberBeforeUpdate := neo4jmapper.MapDbNodeToPhoneNumberEntity(phoneNumberDbNode)

	err = h.repositories.Neo4jRepositories.PhoneNumberWriteRepository.UpdatePhoneNumber(ctx, eventData.Tenant, phoneNumberId, eventData.RawPhoneNumber, eventData.Source)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// email address updated
	if phoneNumberBeforeUpdate.RawPhoneNumber != eventData.RawPhoneNumber && phoneNumberBeforeUpdate.E164 != eventData.RawPhoneNumber {
		err = h.repositories.Neo4jRepositories.PhoneNumberWriteRepository.CleanPhoneNumberValidation(ctx, eventData.Tenant, phoneNumberId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil
		}

		if len(eventData.RawPhoneNumber) > 0 {
			ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
			_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*phonenumberpb.PhoneNumberIdGrpcResponse](func() (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
				return h.grpcClients.PhoneNumberClient.RequestPhoneNumberValidation(ctx, &phonenumberpb.RequestPhoneNumberValidationGrpcRequest{
					Tenant:    eventData.Tenant,
					Id:        phoneNumberId,
					AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to request phone number validation for phoneNumberId: %s, tenant: %s, err: %s", phoneNumberId, eventData.Tenant, err.Error())
				return nil
			}
		}
	}

	return nil
}

func (e *PhoneNumberEventHandler) OnPhoneNumberValidated(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberEventHandler.OnPhoneNumberValidated")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberValidatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.PhoneNumberValidateFields{
		E164:          eventData.E164,
		CountryCodeA2: eventData.CountryCodeA2,
		ValidatedAt:   eventData.ValidatedAt,
		Source:        constants.SourceOpenline,
		AppSource:     "validation-api",
	}
	err := e.repositories.Neo4jRepositories.PhoneNumberWriteRepository.PhoneNumberValidated(ctx, eventData.Tenant, phoneNumberId, data)

	return err
}

func (h *PhoneNumberEventHandler) OnPhoneNumberValidationFailed(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberEventHandler.OnPhoneNumberValidationFailed")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberFailedValidationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.PhoneNumberWriteRepository.FailPhoneNumberValidation(ctx, eventData.Tenant, phoneNumberId, eventData.ValidationError)

	return err
}
