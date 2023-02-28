package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

type TicketSyncService interface {
	SyncTickets(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int)
}

type ticketSyncService struct {
	repositories *repository.Repositories
}

func NewTicketSyncService(repositories *repository.Repositories) TicketSyncService {
	return &ticketSyncService{
		repositories: repositories,
	}
}

func (s *ticketSyncService) SyncTickets(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int) {
	completed, failed := 0, 0
	for {
		tickets := dataService.GetTicketsForSync(batchSize, runId)
		if len(tickets) == 0 {
			logrus.Debugf("no tickets found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d tickets from %s for tenant %s", len(tickets), dataService.SourceId(), tenant)

		for _, v := range tickets {
			var failedSync = false

			ticketId, err := s.repositories.TicketRepository.GetMatchedTicketId(ctx, tenant, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed finding existing matched ticket with external reference id %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			// Create new ticket id if not found
			if len(ticketId) == 0 {
				ticketUuid, _ := uuid.NewRandom()
				ticketId = ticketUuid.String()
			}
			v.Id = ticketId

			if !failedSync {
				err = s.repositories.TicketRepository.MergeTicket(ctx, tenant, syncDate, v)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merging ticket with external reference id %v for tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			if v.HasCollaborators() && !failedSync {
				for _, userExternalId := range v.CollaboratorUserExternalIds {
					err = s.repositories.TicketRepository.LinkTicketWithCollaboratorUserByExternalId(ctx, tenant, ticketId, userExternalId, v.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link ticket %v with collaborator user for tenant %v :%v", ticketId, tenant, err)
					}
				}
			}

			if v.HasFollowers() && !failedSync {
				for _, userExternalId := range v.FollowerUserExternalIds {
					err = s.repositories.TicketRepository.LinkTicketWithFollowerUserByExternalId(ctx, tenant, ticketId, userExternalId, v.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link ticket %v with follower user for tenant %v :%v", ticketId, tenant, err)
					}
				}
			}

			if v.HasSubmitter() && !failedSync {
				err = s.repositories.TicketRepository.LinkTicketWithSubmitterUserOrContactByExternalId(ctx, tenant, ticketId, v.SubmitterExternalId, v.ExternalSystem)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link ticket %v with submitter for tenant %v :%v", ticketId, tenant, err)
				}
			}

			if v.HasRequester() && !failedSync {
				err = s.repositories.TicketRepository.LinkTicketWithRequesterUserOrContactByExternalId(ctx, tenant, ticketId, v.RequesterExternalId, v.ExternalSystem)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link ticket %v with requester for tenant %v :%v", ticketId, tenant, err)
				}
			}

			logrus.Debugf("successfully merged ticket with id %v for tenant %v from %v", ticketId, tenant, dataService.SourceId())
			if err := dataService.MarkTicketProcessed(v.ExternalSyncId, runId, failedSync == false); err != nil {
				failed++
				continue
			}
			if failedSync == true {
				failed++
			} else {
				completed++
			}
		}
		if len(tickets) < batchSize {
			break
		}
	}
	return completed, failed
}
