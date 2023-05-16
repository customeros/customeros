package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/utils"
	"log"
	"net/http"
)

const APP_SOURCE = "user-admin-api"

func addRegistrationRoutes(rg *gin.RouterGroup, config *config.Config, cosClient service.CustomerOsClient) {
	rg.POST("/register", func(c *gin.Context) {
		log.Printf("registering user")
		apiKey := c.GetHeader("X-Openline-Api-Key")
		if apiKey != config.Service.ApiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"result": fmt.Sprintf("invalid api key"),
			})
			return
		}
		log.Printf("api key is valid")
		var req model.RegisterRequest
		if err := c.BindJSON(&req); err != nil {
			log.Printf("unable to parse json: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}
		log.Printf("parsed json: %v", req)
		if req.Properties.Workspace != nil {
			tenant, err := cosClient.GetTenantByWorkspace(&model.WorkspaceInput{
				Name:     *req.Properties.Workspace,
				Provider: req.Properties.Provider,
			})
			if err != nil {
				log.Printf("unable to get workspace: %v", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to get workspace: %v", err.Error()),
				})
				return
			}

			if tenant != nil {
				log.Printf("tenant found %s", *tenant)
				var appSource = APP_SOURCE
				id, err := cosClient.CreateUser(&model.UserInput{
					FirstName: req.Properties.FirstName,
					LastName:  req.Properties.LastName,
					Email: model.EmailInput{
						Email:     req.Properties.Email,
						Primary:   true,
						AppSource: &appSource,
					},
					Person: model.PersonInput{
						IdentityId: req.Properties.IdentityId,
						Email:      req.Properties.Email,
						Provider:   req.Properties.Provider,
						AppSource:  &appSource,
					},
				}, *tenant)
				if err != nil {
					log.Printf("unable to create user: %v", err.Error())
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": fmt.Sprintf("unable to create user: %v", err.Error()),
					})
					return
				}
				log.Printf("user created: %s", id)
			} else {
				var appSource = APP_SOURCE
				tenantStr := utils.Sanitize(*req.Properties.Workspace)
				log.Printf("tenant not found for workspace, creating new tenant %s", tenantStr)
				// Workspace is not mapped to a tenant create a new tenant and map it to the workspace
				id, failed := makeTenentAndUser(c, cosClient, tenantStr, appSource, req)
				if failed {
					return
				}
				log.Printf("user created: %s", id)

			}
		} else {
			// no workspace for this e-mail invent a tenant name
			var appSource = APP_SOURCE
			tenantStr := utils.GenerateName()
			log.Printf("user has no workspace, inventing tenant %s", tenantStr)

			id, failed := makeTenentAndUser(c, cosClient, tenantStr, appSource, req)
			if failed {
				return
			}
			log.Printf("user created: %s", id)
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func makeTenentAndUser(c *gin.Context, cosClient service.CustomerOsClient, tenantStr string, appSource string, req model.RegisterRequest) (string, bool) {
	newTenantStr, err := cosClient.MergeTenant(&model.TenantInput{
		Name:      tenantStr,
		AppSource: &appSource})
	if err != nil {
		log.Printf("unable to create tenant: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": fmt.Sprintf("unable to create tenant: %v", err.Error()),
		})
		return "", true
	}

	if req.Properties.Workspace != nil {
		mergeWorkspaceRes, err := cosClient.MergeTenantToWorkspace(&model.WorkspaceInput{
			Name:      *req.Properties.Workspace,
			Provider:  req.Properties.Provider,
			AppSource: &appSource,
		}, newTenantStr)

		if err != nil {
			log.Printf("unable to merge workspace: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to merge workspace: %v", err.Error()),
			})
			return "", true
		}
		if !mergeWorkspaceRes {
			log.Printf("unable to merge workspace")
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to merge workspace"),
			})
			return "", true
		}
	}

	id, err := cosClient.CreateUser(&model.UserInput{
		FirstName: req.Properties.FirstName,
		LastName:  req.Properties.LastName,
		Email: model.EmailInput{
			Email:     req.Properties.Email,
			Primary:   true,
			AppSource: &appSource,
		},
		Person: model.PersonInput{
			IdentityId: req.Properties.IdentityId,
			Email:      req.Properties.Email,
			Provider:   req.Properties.Provider,
			AppSource:  &appSource,
		},
	}, newTenantStr)
	if err != nil {
		log.Printf("unable to create user: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": fmt.Sprintf("unable to create user: %v", err.Error()),
		})
		return "", true
	}
	return id, false
}
