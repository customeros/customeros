package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func InitUserSettingsRoutes(r *gin.Engine, ctx context.Context, commonRepositoryContainer *commonRepository.Repositories, services *service.Services) {
	r.GET("/user/settings/oauth/:playerIdentityId",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),

		func(c *gin.Context) {
			playerIdentityId := c.Param("playerIdentityId")
			userSettings, err := services.OAuthUserSettingsService.GetOAuthUserSettings(playerIdentityId)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, userSettings)
		})

	//r.POST("/user/settings",
	//	commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
	//	commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),
	//	func(c *gin.Context) {
	//
	//		// -- TODO vasi -  move this to the read email async service
	// // move this to the read email async service

	//		conf, err := google.ConfigFromJSON([]byte(secret), gmail.GmailReadonlyScope)
	//		if err != nil {
	//			println("get conf for read only", err.Error())
	//		}
	//
	//		token := oauth2.Token{
	//			AccessToken: "ya29.a0AfH6SM",
	//		}
	//
	//		client := conf.Client(oauth2.NoContext, &token)
	//
	//		var request model.UserSettings
	//
	//		gmailService, err := gmail.New(client)
	//		if err != nil {
	//			println("gmail service failed", err.Error())
	//		}
	//
	//		if err != nil {
	//			println("new service", err.Error())
	//		}
	//
	//		mails, err := gmailService.Users.Messages.List("vasi@openline.ai").Do()
	//
	//		if mails != nil {
	//			println("get emails", mails)
	//		}
	//
	//		if err := c.BindJSON(&request); err != nil {
	//			println(err.Error())
	//			c.AbortWithStatus(500) //todo
	//			return
	//		}
	//		request.TenantName = c.Keys["TenantName"].(string)
	//		services.OAuthUserSettingsService.Save(&request)
	//		c.JSON(200, request)
	//	})
}
