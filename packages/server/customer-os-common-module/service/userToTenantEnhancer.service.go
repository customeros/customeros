package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	cr "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type HeaderAllowance string

const (
	USERNAME           HeaderAllowance = "USERNAME"
	TENANT             HeaderAllowance = "TENANT"
	USERNAME_OR_TENANT HeaderAllowance = "USERNAME_OR_TENANT"
)

const UsernameHeader = "X-Openline-USERNAME"
const TenantHeader = "X-Openline-TENANT"

func TenantUserContextEnhancer(ctx context.Context, headerAllowance HeaderAllowance, cr *cr.Repositories) func(c *gin.Context) {
	return func(c *gin.Context) {
		tenantHeader := c.GetHeader(TenantHeader)
		usernameHeader := c.GetHeader(UsernameHeader)
		var (
			tenantExists bool
			userId       string
			tenantName   string
			err          error
		)

		switch headerAllowance {
		case TENANT:
			tenantExists, err = checkTenantHeader(c, tenantHeader, cr, ctx)
			if err != nil {
				return
			}
			c.Set("TenantName", tenantHeader)

		case USERNAME:
			userId, tenantName, err = checkUsernameHeader(c, usernameHeader, cr, ctx)
			if err != nil {
				return
			}
			c.Set("TenantName", tenantName)
			c.Set("UserId", userId)
			c.Set("UserEmail", usernameHeader)

		case USERNAME_OR_TENANT:
			if tenantHeader == "" && usernameHeader == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if tenantHeader != "" {
				tenantExists, err = checkTenantHeader(c, tenantHeader, cr, ctx)
				if err != nil {
					return
				}
				if tenantExists {
					c.Set("TenantName", tenantHeader)
				}
			}
			if usernameHeader != "" {
				userId, tenantName, err = checkUsernameHeader(c, usernameHeader, cr, ctx)
				if err != nil {
					return
				}
				c.Set("TenantName", tenantName)
				c.Set("UserId", userId)
				c.Set("UserEmail", usernameHeader)
			}
			c.Next()
			return
		default:
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func checkTenantHeader(c *gin.Context, tenantHeader string, cr *cr.Repositories, ctx context.Context) (bool, error) {
	if tenantHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return false, fmt.Errorf("missing tenant header")
	}

	exists, err := cr.TenantRepository.TenantExists(ctx, tenantHeader)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return false, fmt.Errorf("failed to check tenant existence: %v", err)
	}
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return false, fmt.Errorf("tenant does not exist")
	}

	return true, nil
}

func checkUsernameHeader(c *gin.Context, usernameHeader string, cr *cr.Repositories, ctx context.Context) (string, string, error) {
	if usernameHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return "", "", fmt.Errorf("missing username header")
	}

	userId, tenantName, err := cr.UserRepository.FindUserByEmail(ctx, usernameHeader)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return "", "", fmt.Errorf("failed to find user: %v", err)
	}
	if tenantName == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return "", "", fmt.Errorf("user has no associated tenant")
	}

	return userId, tenantName, nil
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
		_, tenant, err := userRepository.FindUserByEmail(ctx, kh[0])

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
