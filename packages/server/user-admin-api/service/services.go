package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	authServices "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/service"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"gorm.io/gorm"
)

type Services struct {
	CommonServices *commonService.Services
	AuthServices   *authServices.Services

	CustomerOsClient CustomerOsClient
}

func InitServices(cfg *config.Config, db *gorm.DB, driver *neo4j.DriverWithContext) *Services {

	return &Services{
		CommonServices:   commonService.InitServices(db, driver),
		AuthServices:     authServices.InitServices(nil, db),
		CustomerOsClient: NewCustomerOsClient(cfg),
	}
}
