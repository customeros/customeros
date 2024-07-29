package email_validation

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
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type emailEventHandler struct {
	log          logger.Logger
	cfg          *config.Config
	grpcClients  *grpc_client.Clients
	repositories *repository.Repositories
}

func NewEmailEventHandler(log logger.Logger, cfg *config.Config, grpcClients *grpc_client.Clients, repositories *repository.Repositories) *emailEventHandler {
	return &emailEventHandler{
		log:          log,
		cfg:          cfg,
		grpcClients:  grpcClients,
		repositories: repositories,
	}
}

type EmailValidate struct {
	Email string `json:"email" validate:"required,email"`
}

type EmailValidationResponseV1 struct {
	Error           string `json:"error"`
	IsReachable     string `json:"isReachable"`
	Email           string `json:"email"`
	AcceptsMail     bool   `json:"acceptsMail"`
	CanConnectSmtp  bool   `json:"canConnectSmtp"`
	HasFullInbox    bool   `json:"hasFullInbox"`
	IsCatchAll      bool   `json:"isCatchAll"`
	IsDeliverable   bool   `json:"isDeliverable"`
	IsDisabled      bool   `json:"isDisabled"`
	IsDisposable    bool   `json:"isDisposable"`
	IsRoleAccount   bool   `json:"isRoleAccount"`
	Address         string `json:"address"`
	Domain          string `json:"domain"`
	IsValidSyntax   bool   `json:"isValidSyntax"`
	Username        string `json:"username"`
	NormalizedEmail string `json:"normalizedEmail"`
}

func (h *emailEventHandler) OnEmailCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event.EmailCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	emailId := email.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	return h.validateEmail(ctx, eventData.Tenant, emailId, eventData.RawEmail)
}

func (h *emailEventHandler) OnEmailValidate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailValidate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event.EmailValidateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	emailId := email.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	emailDbNode, err := h.repositories.Neo4jRepositories.EmailReadRepository.GetById(ctx, eventData.Tenant, emailId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	emailEntity := neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)

	return h.validateEmail(ctx, eventData.Tenant, emailId, emailEntity.RawEmail)
}

func (h *emailEventHandler) validateEmail(ctx context.Context, tenant, emailId, emailToValidate string) error {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "EmailEventHandler.validateEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)

	emailValidate := EmailValidate{
		Email: strings.TrimSpace(emailToValidate),
	}

	preValidationErr := validator.GetValidator().Struct(emailValidate)
	if preValidationErr != nil {
		return h.sendEmailFailedValidationEvent(ctx, tenant, emailId, emailToValidate, preValidationErr.Error())
	}
	evJSON, err := json.Marshal(emailValidate)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, tenant, emailId, emailToValidate, err.Error())
	}
	requestBody := []byte(string(evJSON))
	req, err := http.NewRequest("POST", h.cfg.Services.ValidationApi+"/validateEmail", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, tenant, emailId, emailToValidate, err.Error())
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
		return h.sendEmailFailedValidationEvent(ctx, tenant, emailId, emailToValidate, err.Error())
	}
	defer response.Body.Close()
	var result EmailValidationResponseV1
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, tenant, emailId, emailToValidate, err.Error())
	}
	if result.IsReachable == "" {
		errMsg := utils.StringFirstNonEmpty(result.Error, "IsReachable flag not set. Email not passed validation.")
		return h.sendEmailFailedValidationEvent(ctx, tenant, emailId, emailToValidate, errMsg)
	}
	email := utils.StringFirstNonEmpty(result.Address, result.NormalizedEmail)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
		return h.grpcClients.EmailClient.PassEmailValidation(ctx, &emailpb.PassEmailValidationGrpcRequest{
			Tenant:         tenant,
			EmailId:        emailId,
			AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
			RawEmail:       emailToValidate,
			Email:          email,
			IsReachable:    result.IsReachable,
			ErrorMessage:   result.Error,
			Domain:         result.Domain,
			Username:       result.Username,
			AcceptsMail:    result.AcceptsMail,
			CanConnectSmtp: result.CanConnectSmtp,
			HasFullInbox:   result.HasFullInbox,
			IsCatchAll:     result.IsCatchAll,
			IsDisabled:     result.IsDisabled,
			IsValidSyntax:  result.IsValidSyntax,
			IsDeliverable:  result.IsDeliverable,
			IsDisposable:   result.IsDisposable,
			IsRoleAccount:  result.IsRoleAccount,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed sending passed email validation event for email %s for tenant %s: %s", emailId, tenant, err.Error())
	}
	return err
}

func (h *emailEventHandler) sendEmailFailedValidationEvent(ctx context.Context, tenant, emailId, rawEmail string, errorMessage string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.sendEmailFailedValidationEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("emailId", emailId), log.String("errorMessage", errorMessage))

	h.log.Errorf("Failed validating email %s for tenant %s: %s", emailId, tenant, errorMessage)
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
		return h.grpcClients.EmailClient.FailEmailValidation(ctx, &emailpb.FailEmailValidationGrpcRequest{
			Tenant:       tenant,
			EmailId:      emailId,
			RawEmail:     rawEmail,
			AppSource:    constants.AppSourceEventProcessingPlatformSubscribers,
			ErrorMessage: utils.StringFirstNonEmpty(errorMessage, "Error message not available"),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed sending failed email validation event for email %s for tenant %s: %s", emailId, tenant, err.Error())
	}
	return err
}
