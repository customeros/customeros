package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
)

type EmailingService interface {
	GenerateEmailSpyPixelUrl(ctx context.Context, tenant, publicUrl, uniqueMessageId, campaign, recipientId string, trackOpens bool) (url string, mid string, err error)
	GenerateEmailLinkUrl(ctx context.Context, tenant, publicUrl, redirectUrl, uniqueMessageId, campaign, recipientId string, trackClicks bool) (url string, mid string, lid string, err error)
	GenerateEmailUnsubscribeUrl(ctx context.Context, tenant, publicUrl, unsubscribeUrl, uniqueMessageId, campaign, recipientId string) (url string, mid string, err error)
}

type emailingService struct {
	log      logger.Logger
	services *Services
}

func NewEmailingService(log logger.Logger, services *Services) EmailingService {
	return &emailingService{
		log:      log,
		services: services,
	}
}

func (s emailingService) GenerateEmailSpyPixelUrl(ctx context.Context, tenant, publicUrl, uniqueMessageId, campaign, recipientId string, trackOpens bool) (url string, mid string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailingService.GenerateEmailSpyPixelUrl")
	defer span.Finish()
	span.SetTag("tenant", tenant)
	span.LogFields(
		tracingLog.String("publicUrl", publicUrl),
		tracingLog.String("uniqueMessageId", uniqueMessageId),
		tracingLog.String("campaign", campaign),
		tracingLog.String("recipientId", recipientId),
		tracingLog.Bool("trackOpens", trackOpens))

	mid = uniqueMessageId
	if mid == "" {
		mid = utils.GenerateRandomString(64)
	}

	emailLookup, err := s.services.PostgresRepositories.EmailLookupRepository.Create(ctx, postgresentity.EmailLookup{
		Tenant:        tenant,
		MessageId:     mid,
		Campaign:      campaign,
		Type:          postgresentity.EmailLookupTypeSpyPixel,
		TrackOpens:    trackOpens,
		TrackerDomain: publicUrl,
		RecipientId:   recipientId,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error creating email lookup: %v", err)
		return "", "", err
	}

	return publicUrl + "/v1/s?c=" + emailLookup.ID, mid, nil
}

func (s emailingService) GenerateEmailLinkUrl(ctx context.Context, tenant, publicUrl, redirectUrl, uniqueMessageId, campaign, recipientId string, trackClicks bool) (url string, mid string, lid string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailingService.GenerateEmailLinkUrl")
	defer span.Finish()
	span.SetTag("tenant", tenant)
	span.LogFields(
		tracingLog.String("publicUrl", publicUrl),
		tracingLog.String("redirectUrl", redirectUrl),
		tracingLog.String("uniqueMessageId", uniqueMessageId),
		tracingLog.String("campaign", campaign),
		tracingLog.String("recipientId", recipientId),
		tracingLog.Bool("trackClicks", trackClicks))

	mid = uniqueMessageId
	if mid == "" {
		mid = utils.GenerateRandomString(64)
	}

	lid = utils.GenerateRandomString(64)

	emailLookup, err := s.services.PostgresRepositories.EmailLookupRepository.Create(ctx, postgresentity.EmailLookup{
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
		tracing.TraceErr(span, err)
		s.log.Errorf("Error creating email lookup: %v", err)
		return "", "", "", err
	}

	return publicUrl + "/v1/l?c=" + emailLookup.ID, mid, lid, nil
}

func (s emailingService) GenerateEmailUnsubscribeUrl(ctx context.Context, tenant, publicUrl, unsubscribeUrl, uniqueMessageId, campaign, recipientId string) (url string, mid string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailingService.GenerateEmailUnsubscribeUrl")
	defer span.Finish()
	span.SetTag("tenant", tenant)
	span.LogFields(
		tracingLog.String("publicUrl", publicUrl),
		tracingLog.String("unsubscribeUrl", unsubscribeUrl),
		tracingLog.String("uniqueMessageId", uniqueMessageId),
		tracingLog.String("campaign", campaign),
		tracingLog.String("recipientId", recipientId))

	mid = uniqueMessageId
	if mid == "" {
		mid = utils.GenerateRandomString(64)
	}

	emailLookup, err := s.services.PostgresRepositories.EmailLookupRepository.Create(ctx, postgresentity.EmailLookup{
		Tenant:         tenant,
		MessageId:      mid,
		UnsubscribeUrl: unsubscribeUrl,
		Campaign:       campaign,
		Type:           postgresentity.EmailLookupTypeUnsubscribe,
		TrackerDomain:  publicUrl,
		RecipientId:    recipientId,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error creating email lookup: %v", err)
		return "", "", err
	}

	return publicUrl + "/v1/l?u=" + emailLookup.ID, mid, nil
}
