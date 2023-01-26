package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"google.golang.org/grpc/metadata"
)

const UsernameHeader = "X-Openline-USERNAME"

func UserToTenantEnhancer(userToTenantRepository repository.UserToTenantRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		uh := c.GetHeader(UsernameHeader)
		if uh != "" {

			tenantResult := userToTenantRepository.FindTenantByUsername(uh)

			if tenantResult.Error != nil {
				c.AbortWithStatus(401)
				return
			}

			tenant := tenantResult.Result.(string)

			if len(tenant) == 0 {
				c.AbortWithStatus(401)
				return
			} else {
				if c.Keys == nil {
					c.Keys = map[string]any{}
				}
				c.Keys["TenantName"] = tenant
			}

			c.Next()
			// illegal request, terminate the current process
		} else {
			c.AbortWithStatus(401)
			return
		}

	}
}

func GetUsernameMetadataForGRPC(ctx context.Context) (*string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata")
	}

	kh := md.Get(UsernameHeader)
	if kh != nil && len(kh) == 1 {
		return &kh[0], nil
	}
	return nil, errors.New("no username header")
}

func GetTenantForUsernameForGRPC(ctx context.Context, userToTenantRepository repository.UserToTenantRepository) (*string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata")
	}

	kh := md.Get(UsernameHeader)
	if kh != nil && len(kh) == 1 {
		tenantResult := userToTenantRepository.FindTenantByUsername(kh[0])

		if tenantResult.Error != nil && tenantResult.Error.Error() != "record not found" {
			return nil, tenantResult.Error
		}

		tenantName := tenantResult.Result.(string)
		return &tenantName, nil
	}
	return nil, errors.New("no username header")
}
