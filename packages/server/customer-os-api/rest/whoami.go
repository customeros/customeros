package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
)

type EmailElement struct {
	Email   string `json:"email"`
	Label   string `json:"label"`
	Primary bool   `json:"primary"`
}

type UserElement struct {
	UserId    string         `json:"userId"`
	FisrtName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Tenant    string         `json:"tenant"`
	Default   bool           `json:"default"`
	Emails    []EmailElement `json:"emails"`
}

type WhoAmIResponse struct {
	Users []UserElement `json:"users"`
}

func WhoamiHandler(serviceContainer *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		identityId := c.GetHeader(security.IdentityIdHeader)
		if identityId == "" {
			c.JSON(400, gin.H{"error": "Missing identityId header"})
			return
		}

		users, err := serviceContainer.PlayerService.GetUsersByIdentityId(c, identityId)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		response := WhoAmIResponse{}
		for _, user := range *users {
			newContext := common.WithCustomContext(c, &common.CustomContext{Tenant: user.Tenant, Roles: []string{model.RoleUser.String()}, IdentityId: identityId})
			emails, err := serviceContainer.EmailService.GetAllFor(newContext, commonModel.USER, user.Id)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			var emailList []EmailElement
			for _, email := range *emails {
				emailList = append(emailList, EmailElement{
					Email:   email.Email,
					Label:   email.Label,
					Primary: email.Primary,
				})
			}
			response.Users = append(response.Users, UserElement{
				UserId:    user.Id,
				FisrtName: user.FirstName,
				LastName:  user.LastName,
				Tenant:    user.Tenant,
				Default:   user.DefaultForPlayer,
				Emails:    emailList,
			})

		}

		c.JSON(200, response)
	}
}
