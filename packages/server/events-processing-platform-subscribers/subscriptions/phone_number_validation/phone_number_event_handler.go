package phone_number_validation

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commonTracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
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
	"net/http"
	"strings"
)

type phoneNumberEventHandler struct {
	repositories *repository.Repositories
	log          logger.Logger
	cfg          *config.Config
	grpcClients  *grpc_client.Clients
}

func NewPhoneNumberEventHandler(repositories *repository.Repositories, log logger.Logger, cfg *config.Config, grpcClients *grpc_client.Clients) *phoneNumberEventHandler {
	return &phoneNumberEventHandler{
		repositories: repositories,
		log:          log,
		cfg:          cfg,
		grpcClients:  grpcClients,
	}
}

type PhoneNumberValidateRequest struct {
	PhoneNumber   string `json:"phoneNumber" validate:"required"`
	CountryCodeA2 string `json:"country"`
}

type PhoneNumberValidationResponseV1 struct {
	E164      string `json:"e164"`
	Error     string `json:"error"`
	Valid     bool   `json:"valid"`
	CountryA2 string `json:"countryA2"`
}

func (h *phoneNumberEventHandler) OnPhoneNumberCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "phoneNumberEventHandler.OnPhoneNumberCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.PhoneNumberCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tenant := eventData.Tenant
	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)
	span.SetTag(tracing.SpanTagTenant, tenant)

	phoneNumberNodeAvailable := subscriptions.WaitCheckNodeExistsInNeo4j(ctx, h.repositories.Neo4jRepositories, eventData.Tenant, phoneNumberId, neo4jutil.NodeLabelPhoneNumber)
	if !phoneNumberNodeAvailable {
		err := errors.Errorf("%s node %s not available in neo4j", neo4jutil.NodeLabelPhoneNumber, phoneNumberId)
		tracing.TraceErr(span, err)
		return err
	}

	rawPhoneNumber := eventData.RawPhoneNumber

	return h.validatePhoneNumber(ctx, tenant, phoneNumberId, rawPhoneNumber)
}

func (h *phoneNumberEventHandler) OnPhoneNumberValidate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "phoneNumberEventHandler.OnPhoneNumberValidate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.PhoneNumberValidateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)

	phoneNumberDbNode, err := h.repositories.Neo4jRepositories.PhoneNumberReadRepository.GetById(ctx, eventData.Tenant, phoneNumberId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	phoneNumberEntity := neo4jmapper.MapDbNodeToPhoneNumberEntity(phoneNumberDbNode)

	return h.validatePhoneNumber(ctx, eventData.Tenant, phoneNumberId, phoneNumberEntity.RawPhoneNumber)
}

func (h *phoneNumberEventHandler) validatePhoneNumber(ctx context.Context, tenant, phoneNumberId, rawPhoneNumber string) error {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "EmailEventHandler.validateEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId), log.String("rawPhoneNumber", rawPhoneNumber))

	countryCodeA2, err := h.repositories.Neo4jRepositories.PhoneNumberReadRepository.GetCountryCodeA2ForPhoneNumber(ctx, tenant, phoneNumberId)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, err.Error())
	}

	phoneNumberValidate := PhoneNumberValidateRequest{
		PhoneNumber:   strings.TrimSpace(rawPhoneNumber),
		CountryCodeA2: countryCodeA2,
	}

	preValidationErr := validator.GetValidator().Struct(phoneNumberValidate)
	if preValidationErr != nil {
		tracing.TraceErr(span, preValidationErr)
		return h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, preValidationErr.Error())
	}
	evJSON, err := json.Marshal(phoneNumberValidate)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, err.Error())
	}
	requestBody := []byte(string(evJSON))
	req, err := http.NewRequest("POST", h.cfg.Services.ValidationApi+"/validatePhoneNumber", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, err.Error())
	}
	// Inject span context into the HTTP request
	req = commonTracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, h.cfg.Services.ValidationApiKey)
	req.Header.Set(security.TenantHeader, tenant)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, err.Error())
	}
	defer response.Body.Close()
	var result PhoneNumberValidationResponseV1
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, err.Error())
	}
	if !result.Valid {
		return h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, result.Error)
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*phonenumberpb.PhoneNumberIdGrpcResponse](func() (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
		return h.grpcClients.PhoneNumberClient.PassPhoneNumberValidation(ctx, &phonenumberpb.PassPhoneNumberValidationGrpcRequest{
			Tenant:        tenant,
			PhoneNumberId: phoneNumberId,
			PhoneNumber:   rawPhoneNumber,
			E164:          result.E164,
			CountryCodeA2: result.CountryA2,
			AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed sending passed phone number validation event for phone number %s for tenant %s: %s", phoneNumberId, tenant, err.Error())
	}
	return err
}

func (h *phoneNumberEventHandler) sendPhoneNumberFailedValidationEvent(ctx context.Context, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, errorMessage string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.sendEmailFailedValidationEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId), log.String("rawPhoneNumber", rawPhoneNumber), log.String("errorMessage", errorMessage))

	h.log.Errorf("Failed validating phone number %s for tenant %s: %s", phoneNumberId, tenant, errorMessage)
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*phonenumberpb.PhoneNumberIdGrpcResponse](func() (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
		return h.grpcClients.PhoneNumberClient.FailPhoneNumberValidation(ctx, &phonenumberpb.FailPhoneNumberValidationGrpcRequest{
			Tenant:        tenant,
			PhoneNumberId: phoneNumberId,
			PhoneNumber:   rawPhoneNumber,
			CountryCodeA2: countryCodeA2,
			AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
			ErrorMessage:  utils.StringFirstNonEmpty(errorMessage, "Error message not available"),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed sending failed phone number validation event for phone number %s for tenant %s: %s", phoneNumberId, tenant, err.Error())
	}
	return err
}
