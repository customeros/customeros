package session_manager

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
)

func (s *sessionManager) AnalyzeSession(ctx context.Context, msg *nats.Msg) error {
	return nil
}

func (s *sessionManager) determineLeadSource(ctx context.Context, href, referrer string) (LeadSource, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "sessionManager.determineLeadSource")
	defer spans.Finish()

	leadSource := LeadSource{}

	if referrer == "" {
		leadSource.Channel = enum.ChannelDirect
		return leadSource, nil
	}

	parsedReferrer, err := utils.ParseURL(referrer)
	if err != nil {
		spans.TraceError(err)
		return leadSource, nil
	}

	isSearch, searchEngine := isSearchEngine(parsedReferrer.Host)
	isSocial, socialPlatform := isSocialPlatform(parsedReferrer)
	isEmail := isEmail(parsedReferrer.Host)

	parsedHref, err := utils.ParseURL(href)
	if err != nil {
		spans.TraceError(err)
		return leadSource, nil
	}

	campaignMetadata := parseCampaignMetadata(parsedHref)

	isPaid := campaignMetadata != nil && campaignMetadata.IsPaid

	switch {
	case isSearch && isPaid:
		leadSource.AdPlatform = campaignMetadata.AdPlatform
		leadSource.Channel = enum.ChannelPaidSearch

	case isSearch && !isPaid:
		leadSource.SearchEngine = searchEngine
		leadSource.Channel = enum.ChannelOrganicSearch

	case isSocial && isPaid:
		leadSource.AdPlatform = campaignMetadata.AdPlatform
		leadSource.Channel = enum.ChannelPaidSocial

	case isSocial && !isPaid:
		leadSource.SocialPlatform = socialPlatform
		leadSource.Channel = enum.ChannelOrganicSocial

	case isEmail && !isPaid:
		leadSource.Channel = enum.ChannelEmail

	case isPaid:
		leadSource.Channel = enum.ChannelPaidSocial
		leadSource.AdPlatform = campaignMetadata.AdPlatform

	default:
		leadSource.ReferrerHost = parsedReferrer.Host
		leadSource.ReferrerPath = parsedReferrer.Path
		leadSource.Channel = enum.ChannelReferral
	}

	leadSource.Campaign = campaignMetadata

	return leadSource, nil
}
