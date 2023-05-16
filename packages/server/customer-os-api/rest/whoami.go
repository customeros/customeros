package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
)

type UserElement struct {
	UserId    string `json:"userId"`
	FisrtName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Tenant    string `json:"tenant"`
	Default   bool   `json:"default"`
}

type WhoAmIResponse struct {
	Users []UserElement `json:"users"`
}

func WhoamiHandler(serviceContainer *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		identityId := c.GetHeader(commonService.IdentityIdHeader)
		if identityId == "" {
			c.JSON(400, gin.H{"error": "Missing identityId header"})
			return
		}

		users, err := serviceContainer.PersonService.GetUsersByIdentityId(c, identityId)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		response := WhoAmIResponse{}
		for _, user := range *users {
			response.Users = append(response.Users, UserElement{
				UserId:    user.Id,
				FisrtName: user.FirstName,
				LastName:  user.LastName,
				Tenant:    user.Tenant,
				Default:   user.DefaultForPerson,
			})

		}

		c.JSON(200, response)
	}
}
