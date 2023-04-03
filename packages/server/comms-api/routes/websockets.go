package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/ContactHub"
)

type Message struct {
	Message string `json:"message"`
}

func addWebSocketRoutes(rg *gin.RouterGroup, pingInterval int, hub *ContactHub.ContactHub) {

	rg.GET("/ws-participant/:participantId", func(c *gin.Context) {
		participantId := c.Param("participantId")
		if participantId == "" {
			c.JSON(400, gin.H{"msg": "participantId missing from path"})
			return
		}
		ContactHub.ServeContactWs(participantId, hub, c.Writer, c.Request, pingInterval)
	})
}
