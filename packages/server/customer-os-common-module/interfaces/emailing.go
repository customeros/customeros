package interfaces

import "context"

type EmailingService interface {
	GenerateEmailSpyPixelUrl(ctx context.Context, tenant, publicUrl, uniqueMessageId, campaign, recipientId string, trackOpens bool) (url string, mid string, err error)
	GenerateEmailLinkUrl(ctx context.Context, tenant, publicUrl, redirectUrl, uniqueMessageId, campaign, recipientId string, trackClicks bool) (url string, mid string, lid string, err error)
	GenerateEmailUnsubscribeUrl(ctx context.Context, tenant, publicUrl, unsubscribeUrl, uniqueMessageId, campaign, recipientId string) (url string, mid string, err error)
}
