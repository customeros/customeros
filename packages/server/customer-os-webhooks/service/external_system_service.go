package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ExternalSystemService interface {
	MergeExternalSystem(ctx context.Context, tenant, externalSystem string) error
	SyncExternalSystem(ctx context.Context, data model.ExternalSystemData) (SyncResult, error)
}

type externalSystemService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
	caches       *caches.Cache
}

func NewExternalSystemService(log logger.Logger, repositories *repository.Repositories, caches *caches.Cache, services *Services) ExternalSystemService {
	return &externalSystemService{
		log:          log,
		repositories: repositories,
		caches:       caches,
		services:     services,
	}
}

func (s *externalSystemService) MergeExternalSystem(ctx context.Context, tenant, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.MergeExternalSystem")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", externalSystem))

	if externalSystem == "" {
		return nil
	}

	if !s.caches.CheckExternalSystem(tenant, externalSystem) {
		err := s.services.CommonServices.ExternalSystemService.MergeExternalSystem(ctx, tenant, externalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		s.caches.AddExternalSystem(tenant, externalSystem)
	}
	return nil
}

func (s *externalSystemService) SyncExternalSystem(ctx context.Context, externalSystemInput model.ExternalSystemData) (SyncResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "externalSystemService.SyncExternalSystem")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	if externalSystemInput.ExternalSystem == "" {
		tracing.TraceErr(span, errors.ErrMissingExternalSystem)
		return SyncResult{}, errors.ErrMissingExternalSystem
	}

	syncDate := utils.Now()
	var statuses []SyncStatus
	var tenant = common.GetTenantFromContext(ctx)
	reason := ""
	failedSync := false

	externalSystemInput.Normalize()
	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, externalSystemInput.ExternalSystem)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err, log.String("externalSystem", externalSystemInput.ExternalSystem))
		reason := fmt.Sprintf("failed merging external system %s for tenant %s :%s", externalSystemInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("result", "failed"))
		statuses = append(statuses, NewFailedSyncStatus(reason))
	}
	if !failedSync {
		switch neo4jenum.DecodeExternalSystemId(externalSystemInput.ExternalSystem) {
		case neo4jenum.Stripe:
			err = s.repositories.Neo4jRepositories.ExternalSystemWriteRepository.SetProperty(ctx, tenant, externalSystemInput.ExternalSystem, neo4jentity.PropertyExternalSystemStripePaymentMethodTypes, externalSystemInput.PaymentMethodTypes)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err, log.String("externalSystem", externalSystemInput.ExternalSystem))
				reason = fmt.Sprintf("failed setting stripe payment method types for tenant %s :%s", tenant, err.Error())
				s.log.Error(reason)
			}
		}
	}

	if !failedSync {
		span.LogFields(log.String("result", "success"))
		statuses = append(statuses, NewSuccessfulSyncStatus())
	} else {
		span.LogFields(log.String("result", "failed"))
		statuses = append(statuses, NewFailedSyncStatus(reason))
	}

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), externalSystemInput.ExternalSystem, externalSystemInput.AppSource, "externalSystem", syncDate, statuses)
	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}
