package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

type OrganizationSyncService interface {
	SyncOrganizations(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int)
}

type organizationSyncService struct {
	repositories *repository.Repositories
}

func NewOrganizationSyncService(repositories *repository.Repositories) OrganizationSyncService {
	return &organizationSyncService{
		repositories: repositories,
	}
}

func (s *organizationSyncService) SyncOrganizations(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int) {
	completed, failed := 0, 0
	for {
		organizations := dataService.GetOrganizationsForSync(batchSize, runId)
		if len(organizations) == 0 {
			logrus.Debugf("no organizations found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d organizations from %s for tenant %s", len(organizations), dataService.SourceId(), tenant)

		for _, v := range organizations {
			var failedSync = false

			organizationId, err := s.repositories.OrganizationRepository.MergeOrganization(ctx, tenant, syncDate, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merge organization with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			err = s.repositories.OrganizationRepository.MergeOrganizationDefaultPlace(ctx, tenant, organizationId, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merge organization' place with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			if len(v.OrganizationTypeName) > 0 {
				err = s.repositories.OrganizationRepository.MergeOrganizationType(ctx, tenant, organizationId, v.OrganizationTypeName)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge organization type for organization %v, tenant %v :%v", organizationId, tenant, err)
				}
			}

			logrus.Debugf("successfully merged organization with id %v for tenant %v from %v", organizationId, tenant, dataService.SourceId())
			if err := dataService.MarkOrganizationProcessed(v.ExternalId, runId, failedSync == false); err != nil {
				failed++
				continue
			}
			if failedSync == true {
				failed++
			} else {
				completed++
			}
		}
		if len(organizations) < batchSize {
			break
		}
	}
	return completed, failed
}
