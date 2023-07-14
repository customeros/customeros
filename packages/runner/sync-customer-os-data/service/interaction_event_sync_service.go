package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/sirupsen/logrus"
	"time"
)

type interactionEventSyncService struct {
	repositories *repository.Repositories
}

func NewDefaultInteractionEventSyncService(repositories *repository.Repositories) SyncService {
	return &interactionEventSyncService{
		repositories: repositories,
	}
}

func (s *interactionEventSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	for {
		interactionEvents := dataService.GetDataForSync(common.INTERACTION_EVENTS, batchSize, runId)
		if len(interactionEvents) == 0 {
			logrus.Debugf("no interaction found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d interaction events from %s for tenant %s", len(interactionEvents), dataService.SourceId(), tenant)

		for _, v := range interactionEvents {
			var failedSync = false
			var reason string

			interactionEventInput := v.(entity.InteractionEventData)
			interactionEventInput.Normalize()

			if interactionEventInput.Skip {
				if err := dataService.MarkProcessed(interactionEventInput.SyncId, runId, true, true, interactionEventInput.SkipReason); err != nil {
					failed++
					continue
				}
				skipped++
				continue
			}

			interactionEventId, err := s.repositories.InteractionEventRepository.GetMatchedInteractionEvent(ctx, tenant, interactionEventInput)
			if err != nil {
				failedSync = true
				reason = fmt.Sprintf("failed finding existing matched interaction event with external reference id %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
				logrus.Error(reason)
			}

			// Create new note id if not found
			if interactionEventId == "" {
				ieUuid, _ := uuid.NewRandom()
				interactionEventId = ieUuid.String()
			}
			interactionEventInput.Id = interactionEventId

			if !failedSync {
				err = s.repositories.InteractionEventRepository.MergeInteractionEvent(ctx, tenant, syncDate, interactionEventInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merge interaction event with external reference %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
					logrus.Error(reason)
				}
			}

			if !failedSync && interactionEventInput.IsPartOf() {
				err = s.repositories.InteractionEventRepository.LinkInteractionEventAsPartOfByExternalId(ctx, tenant, interactionEventInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed link interaction event as part of by external reference %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
					logrus.Error(reason)
				}
			}

			if !failedSync && interactionEventInput.HasSender() {
				err = s.repositories.InteractionEventRepository.LinkInteractionEventWithSenderByExternalId(ctx, tenant, interactionEventId, interactionEventInput.ExternalSystem, interactionEventInput.SentBy)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed link interaction event with sender by external reference %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
					logrus.Error(reason)
				}
			}

			if !failedSync && interactionEventInput.HasRecipients() {
				for _, recipient := range interactionEventInput.SentTo {
					err = s.repositories.InteractionEventRepository.LinkInteractionEventWithRecipientByExternalId(ctx, tenant, interactionEventId, interactionEventInput.ExternalSystem, recipient)
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed link interaction event with recipient by external reference %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
						logrus.Error(reason)
					}
				}
			}

			if failedSync == false {
				logrus.Debugf("successfully merged interaction event with id %v for tenant %v from %v", interactionEventId, tenant, dataService.SourceId())
			}
			if err := dataService.MarkProcessed(interactionEventInput.SyncId, runId, failedSync == false, false, reason); err != nil {
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
	return completed, failed, skipped
}
