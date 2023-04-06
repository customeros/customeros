package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/util"
	"net/http"
	"time"
)

func addCallCredentialRoutes(rg *gin.RouterGroup, config *c.Config) {

	rg.GET("/call_credentials", func(c *gin.Context) {
		expiresTime := time.Now().Unix() + int64(config.WebRTC.TTL)
		timeLimitedUser := fmt.Sprintf("%d:%s", expiresTime, c.GetHeader("X-Openline-USERNAME"))
		password := util.GetSignature(timeLimitedUser, config.WebRTC.AuthSecret)
		c.JSON(http.StatusOK, gin.H{"username": timeLimitedUser, "password": password, "ttl": config.WebRTC.TTL})
	})
}
