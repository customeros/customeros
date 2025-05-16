package service

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-webhooks/caches"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/errors"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/model"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/repository"
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
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemService.MergeExternalSystem")
	defer spans.Finish()
	spans.TagString(telemetry.SpanTagExternalSystem, externalSystem)
	spans.LogKV("externalSystem", externalSystem)

	if externalSystem == "" {
		return nil
	}

	if !s.caches.CheckExternalSystem(tenant, externalSystem) {
		err := s.services.CommonServices.ExternalSystemService.MergeExternalSystem(ctx, tenant, externalSystem)
		if err != nil {
			spans.TraceError(err)
			return err
		}
		s.caches.AddExternalSystem(tenant, externalSystem)
	}
	return nil
}

func (s *externalSystemService) SyncExternalSystem(ctx context.Context, externalSystemInput model.ExternalSystemData) (SyncResult, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "externalSystemService.SyncExternalSystem")
	defer spans.Finish()
	spans.TagString(telemetry.SpanTagExternalSystem, externalSystemInput.ExternalSystem)
	spans.LogKV("externalSystem", externalSystemInput.ExternalSystem)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		spans.TraceError(errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	if externalSystemInput.ExternalSystem == "" {
		spans.TraceError(errors.ErrMissingExternalSystem)
		return SyncResult{}, errors.ErrMissingExternalSystem
	}

	syncDate := utils.Now()
	var statuses []SyncStatus
	tenant := common.GetTenantFromContext(ctx)
	reason := ""
	failedSync := false

	externalSystemInput.Normalize()
	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, externalSystemInput.ExternalSystem)
	if err != nil {
		failedSync = true
		spans.TraceError(err)
		spans.LogKV("externalSystem", externalSystemInput.ExternalSystem)
		reason := fmt.Sprintf("failed merging external system %s for tenant %s :%s", externalSystemInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		spans.LogKV("result", "failed")
		statuses = append(statuses, NewFailedSyncStatus(reason))
	}
	if !failedSync {
		switch enum.DecodeSource(externalSystemInput.ExternalSystem) {
		case enum.SourceStripe:
			err = s.repositories.Neo4jRepositories.ExternalSystemWriteRepository.SetProperty(ctx, tenant, externalSystemInput.ExternalSystem, neo4jentity.PropertyExternalSystemStripePaymentMethodTypes, externalSystemInput.PaymentMethodTypes)
			if err != nil {
				failedSync = true
				spans.TraceError(err)
				spans.LogKV("externalSystem", externalSystemInput.ExternalSystem)
				reason = fmt.Sprintf("failed setting stripe payment method types for tenant %s :%s", tenant, err.Error())
				s.log.Error(reason)
			}
		}
	}

	if !failedSync {
		spans.LogKV("result", "success")
		statuses = append(statuses, NewSuccessfulSyncStatus())
	} else {
		spans.LogKV("result", "failed")
		statuses = append(statuses, NewFailedSyncStatus(reason))
	}

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), externalSystemInput.ExternalSystem, externalSystemInput.AppSource, "externalSystem", syncDate, statuses)
	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}
