package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

type RegistrationService interface {
	PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error
}

type registrationService struct {
	services *Services
}

func NewRegistrationService(services *Services) RegistrationService {
	return &registrationService{
		services: services,
	}
}

func (r registrationService) PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.GetAllForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("loggedInUserEmail", loggedInUserEmail)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	//tenant := common.GetTenantFromContext(ctx)

	// Step 1 - Create default user

	return nil
}
