package emailing

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type emailingService struct {
	log      logger.Logger
	postgres *postgres_repository.Repositories
}

func NewEmailingService(log logger.Logger, postgres *postgres_repository.Repositories) interfaces.EmailingService {
	return &emailingService{
		log:      log,
		postgres: postgres,
	}
}

func (s emailingService) GenerateEmailSpyPixelUrl(ctx context.Context, tenant, publicUrl, uniqueMessageId, campaign, recipientId string, trackOpens bool) (url string, mid string, err error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailingService.GenerateEmailSpyPixelUrl")
	defer spans.Finish()
	spans.LogKV("publicUrl", publicUrl)
	spans.LogKV("uniqueMessageId", uniqueMessageId)
	spans.LogKV("campaign", campaign)
	spans.LogKV("recipientId", recipientId)
	spans.LogKV("trackOpens", trackOpens)

	mid = uniqueMessageId
	if mid == "" {
		mid = utils.GenerateRandomString(64)
	}

	emailLookup, err := s.postgres.EmailLookupRepository.Create(ctx, postgresentity.EmailLookup{
		Tenant:        tenant,
		MessageId:     mid,
		Campaign:      campaign,
		Type:          postgresentity.EmailLookupTypeSpyPixel,
		TrackOpens:    trackOpens,
		TrackerDomain: publicUrl,
		RecipientId:   recipientId,
	})
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error creating email lookup: %v", err)
		return "", "", err
	}

	return publicUrl + "/v1/s?c=" + emailLookup.ID, mid, nil
}

func (s emailingService) GenerateEmailLinkUrl(ctx context.Context, tenant, publicUrl, redirectUrl, uniqueMessageId, campaign, recipientId string, trackClicks bool) (url string, mid string, lid string, err error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailingService.GenerateEmailLinkUrl")
	defer spans.Finish()
	spans.LogKV("publicUrl", publicUrl)
	spans.LogKV("redirectUrl", redirectUrl)
	spans.LogKV("uniqueMessageId", uniqueMessageId)
	spans.LogKV("campaign", campaign)
	spans.LogKV("recipientId", recipientId)
	spans.LogKV("trackClicks", trackClicks)

	mid = uniqueMessageId
	if mid == "" {
		mid = utils.GenerateRandomString(64)
	}

	lid = utils.GenerateRandomString(64)

	emailLookup, err := s.postgres.EmailLookupRepository.Create(ctx, postgresentity.EmailLookup{
		Tenant:        tenant,
		MessageId:     mid,
		LinkId:        lid,
		RedirectUrl:   redirectUrl,
		Campaign:      campaign,
		Type:          postgresentity.EmailLookupTypeLink,
		RecipientId:   recipientId,
		TrackClicks:   trackClicks,
		TrackerDomain: publicUrl,
	})
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error creating email lookup: %v", err)
		return "", "", "", err
	}

	return publicUrl + "/v1/l?c=" + emailLookup.ID, mid, lid, nil
}

func (s emailingService) GenerateEmailUnsubscribeUrl(ctx context.Context, tenant, publicUrl, unsubscribeUrl, uniqueMessageId, campaign, recipientId string) (url string, mid string, err error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailingService.GenerateEmailUnsubscribeUrl")
	defer spans.Finish()
	spans.LogKV("publicUrl", publicUrl)
	spans.LogKV("unsubscribeUrl", unsubscribeUrl)
	spans.LogKV("uniqueMessageId", uniqueMessageId)
	spans.LogKV("campaign", campaign)
	spans.LogKV("recipientId", recipientId)

	mid = uniqueMessageId
	if mid == "" {
		mid = utils.GenerateRandomString(64)
	}

	emailLookup, err := s.postgres.EmailLookupRepository.Create(ctx, postgresentity.EmailLookup{
		Tenant:         tenant,
		MessageId:      mid,
		UnsubscribeUrl: unsubscribeUrl,
		Campaign:       campaign,
		Type:           postgresentity.EmailLookupTypeUnsubscribe,
		TrackerDomain:  publicUrl,
		RecipientId:    recipientId,
	})
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error creating email lookup: %v", err)
		return "", "", err
	}

	return publicUrl + "/v1/l?u=" + emailLookup.ID, mid, nil
}
