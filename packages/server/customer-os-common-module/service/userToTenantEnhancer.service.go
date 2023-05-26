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

const (
	KEY_TENANT_NAME = "TenantName"
	KEY_USER_ID     = "UserId"
	KEY_USER_EMAIL  = "UserEmail"
	KEY_USER_ROLES  = "UserRoles"
	KEY_IDENTITY_ID = "IdentityId"
)

const UsernameHeader = "X-Openline-USERNAME"
const TenantHeader = "X-Openline-TENANT"
const IdentityIdHeader = "X-Openline-IDENTITY-ID"

func TenantUserContextEnhancer(ctx context.Context, headerAllowance HeaderAllowance, cr *cr.Repositories) func(c *gin.Context) {
	return func(c *gin.Context) {
		tenantHeader := c.GetHeader(TenantHeader)
		usernameHeader := c.GetHeader(UsernameHeader)
		var (
			tenantExists bool
			userId       string
			tenantName   string
			roles        []string
			err          error
		)
		c.Set(KEY_IDENTITY_ID, c.GetHeader(IdentityIdHeader))

		switch headerAllowance {
		case TENANT:
			tenantExists, err = checkTenantHeader(c, tenantHeader, cr, ctx)
			if err != nil {
				return
			}
			c.Set(KEY_TENANT_NAME, tenantHeader)

		case USERNAME:
			userId, tenantName, roles, err = checkUsernameHeader(c, usernameHeader, cr, ctx)
			if err != nil {
				return
			}
			c.Set(KEY_TENANT_NAME, tenantName)
			c.Set(KEY_USER_ID, userId)
			c.Set(KEY_USER_EMAIL, usernameHeader)
			c.Set(KEY_USER_ROLES, roles)

		case USERNAME_OR_TENANT:
			if tenantHeader == "" && usernameHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message": "User or Tenant must be specified"}},
				})
				c.Abort()
				return
			}

			if tenantHeader != "" {
				tenantExists, err = checkTenantHeader(c, tenantHeader, cr, ctx)
				if err != nil {
					return
				}
				if tenantExists {
					c.Set(KEY_TENANT_NAME, tenantHeader)
				}
			}
			if usernameHeader != "" {
				userId, tenantName, roles, err = checkUsernameHeader(c, usernameHeader, cr, ctx)
				if err != nil {
					return
				}
				c.Set(KEY_TENANT_NAME, tenantName)
				c.Set(KEY_USER_ID, userId)
				c.Set(KEY_USER_EMAIL, usernameHeader)
				c.Set(KEY_USER_ROLES, roles)
			}
			c.Next()
			return
		default:
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []gin.H{{"message": "unknown header Allowance"}},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func checkTenantHeader(c *gin.Context, tenantHeader string, cr *cr.Repositories, ctx context.Context) (bool, error) {
	if tenantHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": "missing tenant header"}},
		})
		c.Abort()
		return false, fmt.Errorf("missing tenant header")
	}

	exists, err := cr.TenantRepository.TenantExists(ctx, tenantHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": fmt.Sprintf("failed to check tenant existence: %v", err)}},
		})
		c.Abort()
		return false, fmt.Errorf("failed to check tenant existence: %v", err)
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": "tenant does not exist"}},
		})
		c.Abort()
		return false, fmt.Errorf("tenant does not exist")
	}

	return true, nil
}

func checkUsernameHeader(c *gin.Context, usernameHeader string, cr *cr.Repositories, ctx context.Context) (string, string, []string, error) {
	if usernameHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": "missing username header"}},
		})
		c.Abort()
		return "", "", []string{}, fmt.Errorf("missing username header")
	}

	userId, tenantName, roles, err := cr.UserRepository.FindUserByEmail(ctx, usernameHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": fmt.Sprintf("failed to find user: %v", err)}},
		})
		c.Abort()
		return "", "", []string{}, fmt.Errorf("failed to find user: %v", err)
	}
	if tenantName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": "user has no associated tenant"}},
		})
		c.Abort()
		return "", "", []string{}, fmt.Errorf("user has no associated tenant")
	}

	return userId, tenantName, roles, nil
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

func GetTenantForUsernameForGRPC(ctx context.Context, userRepository repository.UserRepository) (*string, []string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, []string{}, errors.New("no metadata")
	}

	kh := md.Get(UsernameHeader)
	if kh != nil && len(kh) == 1 {
		_, tenant, roles, err := userRepository.FindUserByEmail(ctx, kh[0])

		if err != nil {
			return nil, []string{}, err
		}

		if len(tenant) == 0 {
			return nil, []string{}, errors.New("no user found")
		}

		return &tenant, roles, nil
	}
	return nil, []string{}, errors.New("no username header")
}
