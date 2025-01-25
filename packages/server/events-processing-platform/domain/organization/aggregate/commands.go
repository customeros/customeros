package aggregate

import (
	"context"
	organizationpb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/customeros/customeros/packages/server/events/event/common"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/organization/command"
	localerror "github.com/customeros/customeros/packages/server/events-processing-platform/domain/organization/errors"
	organizationEvents "github.com/customeros/customeros/packages/server/events-processing-platform/domain/organization/events"
	"github.com/customeros/customeros/packages/server/events-processing-platform/tracing"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *OrganizationAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *organizationpb.CreateBillingProfileGrpcRequest:
		return a.CreateBillingProfile(ctx, r)
	case *organizationpb.UpdateBillingProfileGrpcRequest:
		return nil, a.UpdateBillingProfile(ctx, r)
	case *organizationpb.LinkEmailToBillingProfileGrpcRequest:
		return nil, a.LinkEmailToBillingProfile(ctx, r)
	case *organizationpb.UnlinkEmailFromBillingProfileGrpcRequest:
		return nil, a.UnlinkEmailFromBillingProfile(ctx, r)
	case *organizationpb.LinkLocationToBillingProfileGrpcRequest:
		return nil, a.LinkLocationToBillingProfile(ctx, r)
	case *organizationpb.UnlinkLocationFromBillingProfileGrpcRequest:
		return nil, a.UnlinkLocationFromBillingProfile(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *OrganizationTempAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationTempAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *organizationpb.RefreshRenewalSummaryGrpcRequest:
		return nil, a.refreshRenewalSummary(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *OrganizationAggregate) CreateBillingProfile(ctx context.Context, request *organizationpb.CreateBillingProfileGrpcRequest) (billingProfileId string, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.CreateBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), createdAtNotNil)
	sourceFields := common.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	billingProfileId = utils.NewUUIDIfEmpty(request.BillingProfileId)

	event, err := organizationEvents.NewBillingProfileCreateEvent(a, billingProfileId, request.LegalName, request.TaxId, sourceFields, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "NewBillingProfileCreateEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return billingProfileId, a.Apply(event)
}

func (a *OrganizationAggregate) UpdateBillingProfile(ctx context.Context, request *organizationpb.UpdateBillingProfileGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UpdateBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())
	var fieldsMask []string
	if utils.ContainsElement(request.FieldsMask, organizationpb.BillingProfileFieldMask_BILLING_PROFILE_PROPERTY_LEGAL_NAME) {
		fieldsMask = append(fieldsMask, organizationEvents.FieldMaskLegalName)
	}
	if utils.ContainsElement(request.FieldsMask, organizationpb.BillingProfileFieldMask_BILLING_PROFILE_PROPERTY_TAX_ID) {
		fieldsMask = append(fieldsMask, organizationEvents.FieldMaskTaxId)
	}

	updateEvent, err := organizationEvents.NewBillingProfileUpdateEvent(a, request.BillingProfileId, request.LegalName, request.TaxId, updatedAtNotNil, fieldsMask)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewBillingProfileUpdateEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&updateEvent, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *OrganizationAggregate) LinkEmailToBillingProfile(ctx context.Context, request *organizationpb.LinkEmailToBillingProfileGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.LinkEmailToBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := organizationEvents.NewLinkEmailToBillingProfileEvent(a, request.BillingProfileId, request.EmailId, request.Primary, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLinkEmailToBillingProfileEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) UnlinkEmailFromBillingProfile(ctx context.Context, request *organizationpb.UnlinkEmailFromBillingProfileGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UnlinkEmailFromBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := organizationEvents.NewUnlinkEmailFromBillingProfileEvent(a, request.BillingProfileId, request.EmailId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUnlinkEmailFromBillingProfileEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) LinkLocationToBillingProfile(ctx context.Context, request *organizationpb.LinkLocationToBillingProfileGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.LinkLocationToBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := organizationEvents.NewLinkLocationToBillingProfileEvent(a, request.BillingProfileId, request.LocationId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLinkLocationToBillingProfileEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) UnlinkLocationFromBillingProfile(ctx context.Context, request *organizationpb.UnlinkLocationFromBillingProfileGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UnlinkLocationFromBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := organizationEvents.NewUnlinkLocationFromBillingProfileEvent(a, request.BillingProfileId, request.LocationId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUnlinkLocationFromBillingProfileEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationTempAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationTempAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.RefreshArrCommand:
		return a.refreshArr(ctx, c)

	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *OrganizationAggregate) SetPhoneNumberNonPrimary(ctx context.Context, tenant, phoneNumberId, userId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.SetPhoneNumberNonPrimary")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("phoneNumberId", phoneNumberId), log.String("userId", userId))

	updatedAtNotNil := utils.Now()

	phoneNumber, ok := a.Organization.PhoneNumbers[phoneNumberId]
	if !ok {
		return localerror.ErrPhoneNumberNotFound
	}

	if phoneNumber.Primary {
		event, err := organizationEvents.NewOrganizationLinkPhoneNumberEvent(a, phoneNumberId, phoneNumber.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationLinkPhoneNumberEvent")
		}

		eventstore.EnrichEventWithMetadata(&event, &span, a.Tenant, userId)
		return a.Apply(event)
	}
	return nil
}

func (a *OrganizationTempAggregate) refreshArr(ctx context.Context, cmd *command.RefreshArrCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.refreshArr")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	event, err := organizationEvents.NewOrganizationRefreshArrEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRefreshArrEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationTempAggregate) refreshRenewalSummary(ctx context.Context, request *organizationpb.RefreshRenewalSummaryGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.refreshRenewalSummary")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	event, err := organizationEvents.NewOrganizationRefreshRenewalSummaryEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRefreshRenewalSummaryEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}
