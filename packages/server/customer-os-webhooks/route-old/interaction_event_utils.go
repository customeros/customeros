package route

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func (h *EmailWebhookHandler) handleError(c *gin.Context, span opentracing.Span, status int, message string, err error) {
	if err != nil {
		tracing.TraceErr(span, err)
		h.logger.Errorf("(SyncInteractionEvent) %s: %s", message, err.Error())
	}
	c.JSON(status, gin.H{"error": message})
}

func (h *EmailWebhookHandler) setTenantContext(ctx context.Context, span opentracing.Span, tenant string) context.Context {
	ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})
	span.LogFields(tracingLog.Bool("tenant.found", true))
	span.LogFields(tracingLog.String("tenant.name", tenant))
	span.SetTag(tracing.SpanTagTenant, tenant)
	return ctx
}

func (h *EmailWebhookHandler) isExcludedEmail(ctx context.Context, tenant string, data *model.PostmarkEmailWebhookData) bool {
	exclusions := h.services.CommonServices.Cache.GetEmailExclusion(tenant)
	htmlData := strings.ReplaceAll(data.HtmlBody, "&amp;", "&")
	textData := strings.ReplaceAll(data.TextBody, "&amp;", "&")

	for _, exclusion := range exclusions {
		if h.matchesExclusion(exclusion, data.Subject, htmlData, textData) {
			return true
		}
	}
	return false
}

func (h *EmailWebhookHandler) matchesExclusion(exclusion dto.EmailExclusion, subject, htmlBody, textBody string) bool {
	if exclusion.ExcludeSubject != nil && strings.Contains(subject, *exclusion.ExcludeSubject) {
		return true
	}
	if exclusion.ExcludeBody != nil {
		return strings.Contains(htmlBody, *exclusion.ExcludeBody) || strings.Contains(textBody, *exclusion.ExcludeBody)
	}
	return false
}
