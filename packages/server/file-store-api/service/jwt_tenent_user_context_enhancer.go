package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/model"
	"log"
	"net/http"
	"time"
)

type JWTTenantUserService struct {
	cfg *config.Config
}

const JWTHeader = "X-Openline-JWT"

func NewJWTTenantUserService(cfg *config.Config) *JWTTenantUserService {
	return &JWTTenantUserService{
		cfg: cfg,
	}
}

func (jtus *JWTTenantUserService) GetJWTTenantUserEnhancer() func(c *gin.Context) {
	return func(c *gin.Context) {
		jwtHeader := c.GetHeader(JWTHeader)
		if jwtHeader == "" {
			c.Next()
			return
		} else {
			parsedToken, err := jwt.ParseWithClaims(jwtHeader, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				algo, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])

				}
				if algo.Alg() != jwt.SigningMethodHS256.Alg() {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(jtus.cfg.Service.FileStoreAPIJwtSecret), nil
			})
			// Check if the JWT is valid and hasn't expired
			if !parsedToken.Valid {
				log.Println("Invalid token")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			claims, ok := parsedToken.Claims.(*jwt.StandardClaims)
			if !ok {
				log.Println("Invalid token claims")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if time.Now().Unix() > claims.ExpiresAt {
				log.Println("Token expired")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			userTenant := &model.UserTenant{}
			err = json.Unmarshal([]byte(claims.Audience), userTenant)
			if err != nil {
				log.Println("Invalid Audience")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Set("TenantName", userTenant.TenantName)
			c.Set("UserId", userTenant.UserId)
			c.Set("UserEmail", userTenant.UserEmail)
			// skip the rest of the middleware
			c.Handler()(c)
			c.Abort()
			return
		}

	}
}

func (jtus *JWTTenantUserService) MakeJWT(ctx *gin.Context) {
	expirationTime := time.Now().Add(1 * time.Minute)

	// Create the JWT claims
	userTenant := &model.UserTenant{
		TenantName: ctx.Keys["TenantName"].(string),
		UserEmail:  ctx.Keys["UserEmail"].(string),
		UserId:     ctx.Keys["UserId"].(string),
	}
	userTenantStr, err := json.Marshal(userTenant)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "error creating jwt user tenant: " + err.Error()})
		return
	}
	claims := jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Id:        uuid.New().String(),
		Audience:  string(userTenantStr),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(jtus.cfg.Service.FileStoreAPIJwtSecret)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "error signing jwt token: " + err.Error()})
	}

	ctx.JSON(200, gin.H{"status": "OK", "token": signedToken})
}
