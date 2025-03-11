package handlers

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

// ListMailboxes returns all configured mailboxes
func ListMailboxes(imapService interfaces.IMAPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would need to be implemented - for now just use the status info
		status := imapService.Status()
		c.JSON(http.StatusOK, status)
	}
}

// AddMailbox adds a new mailbox configuration
func AddMailbox(imapService interfaces.IMAPService, mailboxRepository interfaces.MailboxRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "AddMailbox", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		var config models.Mailbox
		err := c.ShouldBindJSON(&config)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = mailboxRepository.SaveMailbox(ctx, config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = imapService.AddMailbox(ctx, &config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "mailbox added", "id": config.ID})
	}
}

// RemoveMailbox removes a mailbox configuration
func RemoveMailbox(imapService interfaces.IMAPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RemoveMailbox", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		id := c.Param("id")
		if err := imapService.RemoveMailbox(ctx, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "mailbox removed", "id": id})
	}
}
