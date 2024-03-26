package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

// HandleCommand processes commands and applies the resulting events to the aggregate.
// Deprecated
func (a *ContractAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.CreateContractCommand:
		return a.createContract(ctx, c)
	case *command.RolloutRenewalOpportunityOnExpirationCommand:
		return a.rolloutRenewalOpportunityOnExpiration(ctx, c)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *ContractAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *contractpb.UpdateContractGrpcRequest:
		return nil, a.updateContract(ctx, r)
	case *contractpb.SoftDeleteContractGrpcRequest:
		return nil, a.softDeleteContract(ctx, r)
	case *contractpb.RefreshContractStatusGrpcRequest:
		return nil, a.refreshContractStatus(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *ContractTempAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractTempAggregate.HandleRequest")
	defer span.Finish()

	switch request.(type) {
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *ContractAggregate) createContract(ctx context.Context, cmd *command.CreateContractCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.createContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	// Assuming you have a utility function to get the current time if the passed time is nil
	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	if cmd.DataFields.RenewalCycle != model.AnnuallyRenewal.String() {
		cmd.DataFields.RenewalPeriods = nil
	}
	if cmd.DataFields.RenewalPeriods != nil {
		if *cmd.DataFields.RenewalPeriods < 1 {
			cmd.DataFields.RenewalPeriods = utils.Int64Ptr(1)
		}
		if *cmd.DataFields.RenewalPeriods > 100 {
			cmd.DataFields.RenewalPeriods = utils.Int64Ptr(100)
		}
	}

	// Determine contract status based start and end dates
	createEvent, err := event.NewContractCreateEvent(a, cmd.DataFields, cmd.Source, cmd.ExternalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *ContractAggregate) updateContract(ctx context.Context, request *contractpb.UpdateContractGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.updateContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	source := utils.StringFirstNonEmpty(sourceFields.Source, a.Contract.Source.SourceOfTruth)

	signedAt := utils.TimestampProtoToTimePtr(request.SignedAt)
	if signedAt != nil && signedAt.Equal(time.Time{}) {
		signedAt = nil
	}
	endedAt := utils.TimestampProtoToTimePtr(request.EndedAt)
	if endedAt != nil && endedAt.Equal(time.Time{}) {
		endedAt = nil
	}
	serviceStartedAt := utils.TimestampProtoToTimePtr(request.ServiceStartedAt)
	if serviceStartedAt != nil && serviceStartedAt.Equal(time.Time{}) {
		serviceStartedAt = nil
	}
	invoicingStartDate := utils.TimestampProtoToTimePtr(request.InvoicingStartDate)
	if invoicingStartDate != nil && invoicingStartDate.Equal(time.Time{}) {
		invoicingStartDate = nil
	}
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())

	dataFields := model.ContractDataFields{
		Name:                   request.Name,
		ServiceStartedAt:       serviceStartedAt,
		SignedAt:               signedAt,
		EndedAt:                endedAt,
		InvoicingStartDate:     invoicingStartDate,
		RenewalCycle:           model.RenewalCycle(request.RenewalCycle).String(),
		ContractUrl:            request.ContractUrl,
		RenewalPeriods:         request.RenewalPeriods,
		Currency:               request.Currency,
		BillingCycle:           model.BillingCycle(request.BillingCycle).String(),
		AddressLine1:           request.AddressLine1,
		AddressLine2:           request.AddressLine2,
		Locality:               request.Locality,
		Country:                request.Country,
		Region:                 request.Region,
		Zip:                    request.Zip,
		OrganizationLegalName:  request.OrganizationLegalName,
		InvoiceEmail:           request.InvoiceEmail,
		InvoiceNote:            request.InvoiceNote,
		NextInvoiceDate:        utils.TimestampProtoToTimePtr(request.NextInvoiceDate),
		CanPayWithCard:         request.CanPayWithCard,
		CanPayWithDirectDebit:  request.CanPayWithDirectDebit,
		CanPayWithBankTransfer: request.CanPayWithBankTransfer,
		PayOnline:              request.PayOnline,
		PayAutomatically:       request.PayAutomatically,
		InvoicingEnabled:       request.InvoicingEnabled,
		AutoRenew:              request.AutoRenew,
		Check:                  request.Check,
		DueDays:                request.DueDays,
	}
	fieldsMask := extractFieldsMask(request.FieldsMask)

	// Validate the dates
	if isUpdated(event.FieldMaskEndedAt, fieldsMask) && dataFields.EndedAt != nil && (dataFields.SignedAt != nil && dataFields.EndedAt.Before(*dataFields.SignedAt) ||
		dataFields.ServiceStartedAt != nil && dataFields.EndedAt.Before(*dataFields.ServiceStartedAt)) {
		return errors.New(constants.FieldValidation + ": endedAt date must be after both signedAt and serviceStartedAt dates")
	}

	// Determine contract status based start and end dates
	if !isUpdated(event.FieldMaskServiceStartedAt, fieldsMask) {
		dataFields.ServiceStartedAt = a.Contract.ServiceStartedAt
	}
	if !isUpdated(event.FieldMaskEndedAt, fieldsMask) {
		dataFields.EndedAt = a.Contract.EndedAt
	}
	// set status field mask if at least any other field is set
	if len(fieldsMask) > 0 {
		fieldsMask = append(fieldsMask, event.FieldMaskStatus)
	}

	// Set renewal periods
	if !isUpdated(event.FieldMaskRenewalCycle, fieldsMask) {
		dataFields.RenewalCycle = a.Contract.RenewalCycle
	}

	if dataFields.RenewalCycle != model.AnnuallyRenewal.String() {
		dataFields.RenewalPeriods = nil
	}
	if dataFields.RenewalPeriods != nil {
		if *dataFields.RenewalPeriods < 1 {
			dataFields.RenewalPeriods = utils.Int64Ptr(1)
		}
		if *dataFields.RenewalPeriods > 100 {
			dataFields.RenewalPeriods = utils.Int64Ptr(100)
		}
	}

	updateEvent, err := event.NewContractUpdateEvent(
		a,
		dataFields,
		externalSystem,
		source,
		updatedAtNotNil,
		fieldsMask,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *ContractAggregate) refreshContractStatus(ctx context.Context, request *contractpb.RefreshContractStatusGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.refreshContractStatus")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	updateEvent, err := event.NewContractRefreshStatusEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractRefreshStatusEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(updateEvent)
}

func (a *ContractAggregate) rolloutRenewalOpportunityOnExpiration(ctx context.Context, cmd *command.RolloutRenewalOpportunityOnExpirationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.rolloutRenewalOpportunityOnExpiration")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	updateEvent, err := event.NewRolloutRenewalOpportunityEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewRolloutRenewalOpportunityEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(updateEvent)
}

func extractFieldsMask(requestFieldsMask []contractpb.ContractFieldMask) []string {
	var fieldsMask []string
	for _, requestFieldMask := range requestFieldsMask {
		switch requestFieldMask {
		case contractpb.ContractFieldMask_CONTRACT_FIELD_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskName)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_SERVICE_STARTED_AT:
			fieldsMask = append(fieldsMask, event.FieldMaskServiceStartedAt)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_SIGNED_AT:
			fieldsMask = append(fieldsMask, event.FieldMaskSignedAt)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_ENDED_AT:
			fieldsMask = append(fieldsMask, event.FieldMaskEndedAt)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_RENEWAL_CYCLE:
			fieldsMask = append(fieldsMask, event.FieldMaskRenewalCycle)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CONTRACT_URL:
			fieldsMask = append(fieldsMask, event.FieldMaskContractURL)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_RENEWAL_PERIODS:
			fieldsMask = append(fieldsMask, event.FieldMaskRenewalPeriods)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_BILLING_CYCLE:
			fieldsMask = append(fieldsMask, event.FieldMaskBillingCycle)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICING_START_DATE:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoicingStartDate)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CURRENCY:
			fieldsMask = append(fieldsMask, event.FieldMaskCurrency)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_ADDRESS_LINE_1:
			fieldsMask = append(fieldsMask, event.FieldMaskAddressLine1)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_ADDRESS_LINE_2:
			fieldsMask = append(fieldsMask, event.FieldMaskAddressLine2)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_LOCALITY:
			fieldsMask = append(fieldsMask, event.FieldMaskLocality)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_COUNTRY:
			fieldsMask = append(fieldsMask, event.FieldMaskCountry)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_REGION:
			fieldsMask = append(fieldsMask, event.FieldMaskRegion)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_ZIP:
			fieldsMask = append(fieldsMask, event.FieldMaskZip)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_ORGANIZATION_LEGAL_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskOrganizationLegalName)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_EMAIL:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoiceEmail)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICE_NOTE:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoiceNote)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_NEXT_INVOICE_DATE:
			fieldsMask = append(fieldsMask, event.FieldMaskNextInvoiceDate)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_INVOICING_ENABLED:
			fieldsMask = append(fieldsMask, event.FieldMaskInvoicingEnabled)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_CARD:
			fieldsMask = append(fieldsMask, event.FieldMaskCanPayWithCard)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_DIRECT_DEBIT:
			fieldsMask = append(fieldsMask, event.FieldMaskCanPayWithDirectDebit)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CAN_PAY_WITH_BANK_TRANSFER:
			fieldsMask = append(fieldsMask, event.FieldMaskCanPayWithBankTransfer)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_PAY_ONLINE:
			fieldsMask = append(fieldsMask, event.FieldMaskPayOnline)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_PAY_AUTOMATICALLY:
			fieldsMask = append(fieldsMask, event.FieldMaskPayAutomatically)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_AUTO_RENEW:
			fieldsMask = append(fieldsMask, event.FieldMaskAutoRenew)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_CHECK:
			fieldsMask = append(fieldsMask, event.FieldMaskCheck)
		case contractpb.ContractFieldMask_CONTRACT_FIELD_DUE_DAYS:
			fieldsMask = append(fieldsMask, event.FieldMaskDueDays)
		}
	}
	fieldsMask = utils.RemoveDuplicates(fieldsMask)
	return fieldsMask
}

func isUpdated(field string, fieldsMask []string) bool {
	return len(fieldsMask) == 0 || utils.Contains(fieldsMask, field)
}

func (a *ContractAggregate) softDeleteContract(ctx context.Context, r *contractpb.SoftDeleteContractGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.softDeleteContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", r)

	deleteEvent, err := event.NewContractDeleteEvent(a, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractDeleteEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&deleteEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(deleteEvent)
}
