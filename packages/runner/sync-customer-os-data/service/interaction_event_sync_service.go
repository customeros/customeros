package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/sirupsen/logrus"
	"time"
)

type InteractionEventSyncService interface {
	SyncInteractionEvents(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int)
}

type interactionEventSyncService struct {
	repositories *repository.Repositories
}

func NewInteractionEventSyncService(repositories *repository.Repositories) InteractionEventSyncService {
	return &interactionEventSyncService{
		repositories: repositories,
	}
}

func (s *interactionEventSyncService) SyncInteractionEvents(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int) {
	completed, failed := 0, 0
	for {
		interactionEvents := dataService.GetInteractionEventsForSync(batchSize, runId)
		if len(interactionEvents) == 0 {
			logrus.Debugf("no interaction found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d interaction events from %s for tenant %s", len(interactionEvents), dataService.SourceId(), tenant)

		for _, interactionEvent := range interactionEvents {
			var failedSync = false

			interactionEventId, err := s.repositories.InteractionEventRepository.GetMatchedInteractionEvent(ctx, tenant, interactionEvent)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed finding existing matched interaction event with external reference id %v for tenant %v :%v", interactionEvent.ExternalId, tenant, err)
			}

			// Create new note id if not found
			if len(interactionEventId) == 0 {
				ieUuid, _ := uuid.NewRandom()
				interactionEventId = ieUuid.String()
			}
			interactionEvent.Id = interactionEventId

			if !failedSync {
				err = s.repositories.InteractionEventRepository.MergeInteractionEvent(ctx, tenant, syncDate, interactionEvent)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge interaction event with external reference %v for tenant %v :%v", interactionEvent.ExternalId, tenant, err)
				}
			}

			if !failedSync && interactionEvent.IsPartOf() {
				err = s.repositories.InteractionEventRepository.LinkInteractionEventAsPartOfByExternalId(ctx, tenant, interactionEvent)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link interaction event as part of by external reference %v for tenant %v :%v", interactionEvent.ExternalId, tenant, err)
				}
			}

			if !failedSync && interactionEvent.HasSender() {
				err = s.repositories.InteractionEventRepository.LinkInteractionEventWithSenderByExternalId(ctx, tenant, interactionEventId, interactionEvent.ExternalSystem, interactionEvent.SentBy)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link interaction event with sender by external reference %v for tenant %v :%v", interactionEvent.ExternalId, tenant, err)
				}
			}

			if !failedSync && interactionEvent.HasRecipients() {
				for _, recipient := range interactionEvent.SentTo {
					err = s.repositories.InteractionEventRepository.LinkInteractionEventWithRecipientByExternalId(ctx, tenant, interactionEventId, interactionEvent.ExternalSystem, recipient)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link interaction event with recipient by external reference %v for tenant %v :%v", interactionEvent.ExternalId, tenant, err)
					}
				}
			}

			if failedSync == false {
				logrus.Debugf("successfully merged interaction event with id %v for tenant %v from %v", interactionEventId, tenant, dataService.SourceId())
			}
			if err := dataService.MarkInteractionEventProcessed(interactionEvent.ExternalSyncId, runId, failedSync == false); err != nil {
				failed++
				continue
			}
			if failedSync == true {
				failed++
			} else {
				completed++
			}
		}
		if len(interactionEvents) < batchSize {
			break
		}
	}
	return completed, failed
}
