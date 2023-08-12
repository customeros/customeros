package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"sync"
	"time"
)

type interactionEventSyncService struct {
	repositories *repository.Repositories
	log          logger.Logger
}

func NewDefaultInteractionEventSyncService(repositories *repository.Repositories, log logger.Logger) SyncService {
	return &interactionEventSyncService{
		repositories: repositories,
		log:          log,
	}
}

func (s *interactionEventSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {

	completed, failed, skipped := 0, 0, 0

	for {

		events := dataService.GetDataForSync(ctx, common.INTERACTION_EVENTS, batchSize, runId)

		if len(events) == 0 {
			break
		}

		s.log.Infof("Syncing %d interaction events", len(events))

		var wg sync.WaitGroup
		wg.Add(len(events))

		results := make(chan result, len(events))
		done := make(chan struct{})

		for _, v := range events {
			v := v

			go func(event entity.InteractionEventData) {
				defer wg.Done()

				var comp, fail, skip int
				s.syncInteractionEvent(ctx, event, dataService, syncDate, tenant, runId, &comp, &fail, &skip)

				results <- result{comp, fail, skip}
			}(v.(entity.InteractionEventData))
		}
		// Wait for goroutines to finish
		go func() {
			wg.Wait()
			close(done)
		}()
		go func() {
			<-done
			close(results)
		}()

		for r := range results {
			completed += r.completed
			failed += r.failed
			skipped += r.skipped
		}

		if len(events) < batchSize {
			break
		}

	}

	return completed, failed, skipped
}

func (s *interactionEventSyncService) syncInteractionEvent(ctx context.Context, interactionEventInput entity.InteractionEventData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventSyncService.syncInteractionEvent")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	interactionEventInput.Normalize()

	if interactionEventInput.Skip {
		if err := dataService.MarkProcessed(ctx, interactionEventInput.SyncId, runId, true, true, interactionEventInput.SkipReason); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	interactionEventId, err := s.repositories.InteractionEventRepository.GetMatchedInteractionEvent(ctx, tenant, interactionEventInput)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched interaction event with external reference id %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
		s.log.Error(reason)
	}

	// Create new note id if not found
	if interactionEventId == "" {
		ieUuid, _ := uuid.NewRandom()
		interactionEventId = ieUuid.String()
	}
	interactionEventInput.Id = interactionEventId
	span.LogFields(log.String("interactionEventId", interactionEventId))

	if !failedSync {
		err = s.repositories.InteractionEventRepository.MergeInteractionEvent(ctx, tenant, syncDate, interactionEventInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge interaction event with external reference %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
	}

	if !failedSync && interactionEventInput.HasSession() {
		err = s.repositories.InteractionEventRepository.MergeInteractionSessionForEvent(ctx, tenant, interactionEventId, interactionEventInput.ExternalSystem, syncDate, interactionEventInput.PartOfSession)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge interaction session by external id %v for tenant %v :%v", interactionEventInput.PartOfSession.ExternalId, tenant, err)
			s.log.Error(reason)
		}
	}

	if !failedSync && interactionEventInput.IsPartOfByExternalId() {
		err = s.repositories.InteractionEventRepository.LinkInteractionEventAsPartOfByExternalId(ctx, tenant, interactionEventInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed link interaction event as part of by external reference %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
	}

	if !failedSync && interactionEventInput.HasSender() {
		err = s.repositories.InteractionEventRepository.LinkInteractionEventWithSenderByExternalId(ctx, tenant, interactionEventId, interactionEventInput.ExternalSystem, interactionEventInput.SentBy)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed link interaction event with sender by external reference %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
	}

	if !failedSync && interactionEventInput.HasRecipients() {
		for _, recipient := range interactionEventInput.SentTo {
			err = s.repositories.InteractionEventRepository.LinkInteractionEventWithRecipientByExternalId(ctx, tenant, interactionEventId, interactionEventInput.ExternalSystem, recipient)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link interaction event with recipient by external reference %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
				s.log.Error(reason)
			}
		}
	}

	if failedSync == false {
		s.log.Debugf("successfully merged interaction event with id %v for tenant %v from %v", interactionEventId, tenant, dataService.SourceId())
	}
	if err := dataService.MarkProcessed(ctx, interactionEventInput.SyncId, runId, failedSync == false, false, reason); err != nil {
		*failed++
		span.LogFields(log.Bool("failedSync", true))
		return
	}
	if failedSync == true {
		*failed++
	} else {
		*completed++
	}
	span.LogFields(log.Bool("failedSync", failedSync))
}
