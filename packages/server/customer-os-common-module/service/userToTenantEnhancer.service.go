package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j"
	"google.golang.org/grpc/metadata"
)

const UsernameHeader = "X-Openline-USERNAME"

func UserToTenantEnhancer(userRepository repository.UserRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		usernameHeader := c.GetHeader(UsernameHeader)
		if usernameHeader != "" {

			userId, tenant, err := userRepository.FindUserByEmail(usernameHeader)

			if err != nil {
				c.AbortWithStatus(401)
				return
			}

			if len(tenant) == 0 {
				c.AbortWithStatus(401)
				return
			} else {
				if c.Keys == nil {
					c.Keys = map[string]any{}
				}
				c.Keys["TenantName"] = tenant
				c.Keys["UserId"] = userId
				c.Keys["UserEmail"] = usernameHeader
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

func GetTenantForUsernameForGRPC(ctx context.Context, userRepository repository.UserRepository) (*string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata")
	}

	kh := md.Get(UsernameHeader)
	if kh != nil && len(kh) == 1 {
		_, tenant, err := userRepository.FindUserByEmail(kh[0])

		if err != nil {
			return nil, err
		}

		if len(tenant) == 0 {
			return nil, errors.New("no user found")
		}

		return &tenant, nil
	}
	return nil, errors.New("no username header")
}
