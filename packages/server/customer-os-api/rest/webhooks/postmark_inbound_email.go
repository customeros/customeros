package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
)

const EXTERNAL_SYSTEM = "mailstack"

func PostmarkInboundEmail(ctx rest.HTTPContext) {
	ctx.ServiceContext, ctx.Span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "PostmarkInboundEmail", c.Request.Header)
	defer span.Finish()
	tracing.TagComponentRest(span)

	// Validate Postmark User-Agent
	if c.Request.UserAgent() == "" || !strings.EqualFold(c.Request.UserAgent(), "Postmark") {
		rest.SendError(c, span, http.StatusForbidden, rest.ErrForbidden)
		return
	}

	// Parse email data
	emailData, err := parseInboundEmail(httpContext)
	if err != nil {
		tracing.LogObjectAsJson(span, "body", c.Request.Body)
		tracing.TraceErr(span, err)
		rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
		return
	}

	// Process email asynchronously
	go func() {
		// Create a new context for the goroutine
		asyncCtx, asyncSpan := tracing.StartSpanFromContext(ctx, "AsyncEmailProcessing")
		defer asyncSpan.Finish()

		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				err := fmt.Errorf("panic recovered in email processing: %v\n%s", r, stack)
				tracing.TraceErr(asyncSpan, err)
			}
		}()

		if err := processInboundEmail(asyncCtx, &emailData); err != nil {
			tracing.TraceErr(asyncSpan, errors.Wrap(err, "failed to process inbound email from Postmark"))
		}
	}()
}

func parseInboundEmail(ctx rest.HTTPContext) (PostmarkInboundEmailData, error) {
	var emailData PostmarkInboundEmailData
	err := ctx.GinContext.BindJSON(&emailData)
	if err != nil {
		return emailData, err
	}

	return emailData, nil
}

func processInboundEmail(ctx rest.HTTPContext, emailData *PostmarkInboundEmailData) error {
	tenant, err := getTenant(ctx, emailData)
	if err != nil {
		return err
	}
	ctx.Tenant = tenant

	participants := emailData.AllParticipantEmails()
	username, err := getUsername(ctx, participants)
	if err != nil || username == "" {
		ctx.Span.LogFields(tracingLog.Bool("mailbox.found", false))
		tracing.TraceErr(ctx.Span, err)
		return err
	}

	messageId := emailData.GetHeaderValue("Message-Id")
	emailExistsInDb, err := ctx.Services.CommonServices.PostgresRepositories.RawEmailRepository.EmailExistsByMessageId(
		*ctx.ServiceContext, EXTERNAL_SYSTEM, ctx.Tenant, username, messageId,
	)
	if err != nil {
		tracing.TraceErr(ctx.Span, err)
		return fmt.Errorf("Unable to determine if email exists in db for messageId %s: %v", messageId, err)
	}

	if emailExistsInDb {
		return nil
	}

	dbEmailEntity := emailData.ToRawDbObject()
	jsonEmailEntity, err := json.Marshal(dbEmailEntity)
	if err != nil {
		return fmt.Errorf("Unable to produce JSON email object for db: %v", err)
	}

	dbErr := ctx.Services.CommonServices.PostgresRepositories.RawEmailRepository.Store(*ctx.ServiceContext, EXTERNAL_SYSTEM, ctx.Tenant, username, dbEmailEntity.ProviderMessageId, messageId, string(jsonEmailEntity), dbEmailEntity.Sent, entity.REAL_TIME)
	if dbErr != nil {
		ctx.Span.LogFields(tracingLog.Object("raw_email", jsonEmailEntity))
		tracing.TraceErr(ctx.Span, err)
		return fmt.Errorf("Error writing email to db: %v", err)
	}

	// Check to see if email is a reply to a flow.  If so, mark as complete.
	// This should be handled in the email processor common service, not here.
	// Same with goal achieved.
	// Move slack notifications to processing service
	// Handle attachments

	return nil
}

func getTenant(ctx rest.HTTPContext, emailData *PostmarkInboundEmailData) (string, error) {
	nameFromBcc := emailData.TenantFromBcc()

	n, err := ctx.Services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantByNameIgnoreCase(*ctx.ServiceContext, nameFromBcc)
	if err != nil {
		ctx.Span.LogFields(tracingLog.Bool("tenant.found", false))
		tracing.TraceErr(ctx.Span, err)
		return "", fmt.Errorf("Cannot identify tenant %s: %v", nameFromBcc, err)
	}

	if n == nil {
		ctx.Span.LogFields(tracingLog.Bool("tenant.found", false))
		return "", fmt.Errorf("No valid tenant %s: %v", nameFromBcc, err)
	}

	tenant := mapper.MapDbNodeToTenantEntity(n)
	ctx.Span.LogFields(tracingLog.Bool("tenant.found", true))
	ctx.Span.LogFields(tracingLog.String("tenant.name", tenant.Name))

	return tenant.Name, nil
}

func getUsername(ctx rest.HTTPContext, EmailParticipants []string) (string, error) {
	for _, p := range EmailParticipants {
		userByEmail, err := ctx.Services.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(
			*ctx.ServiceContext, ctx.Tenant, p)
		if err != nil {
			return "", fmt.Errorf("Error getting username from email %s: %v", p, err)
		}
		if userByEmail != nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("Unable to find user amongst email participants")
}
