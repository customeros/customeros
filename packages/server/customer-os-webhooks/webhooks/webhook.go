package webhook

import (
	"github.com/gin-gonic/gin"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
)

type Registry interface {
	Register(handler WebhookHandler, middleware ...gin.HandlerFunc)
}

type WebhookHandler interface {
	Handle(c *gin.Context)
	GetPath() string
}

type BaseHandler struct {
	services *service.Services
	logger   logger.Logger
}

func (h *BaseHandler) handleError(c *gin.Context, status int, message string, err error) {
	if err != nil {
		h.logger.Errorf("webhook error: %s - %s", message, err.Error())
	}
	c.JSON(status, gin.H{"error": message})
}
