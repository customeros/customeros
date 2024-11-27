package service

import (
	"context"
	"fmt"

	"github.com/customeros/mailwatcher/blscan"
	"github.com/customeros/mailwatcher/domainage"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func (s *mailboxService) ReputationScore(ctx context.Context, domain string) int {
	span, ctx := s.initializeTracing(ctx, "MailboxService.ReputationScore")
	span.LogFields(
		log.String("domain", domain),
	)
	defer span.Finish()

	domainAgePenalty := s.domainAgePenalty(span, domain)
	blacklistPenaltyPct := s.blacklistPenaltyPercent(domain)

	score := (100 - domainAgePenalty) * (1 - blacklistPenaltyPct)

	// 7 day lookback on bounces

	// 7 day lookback on dmarc and spf

	return score
}

func (s *mailboxService) domainAgePenalty(span opentracing.Span, domain string) int {

	domainDates, err := domainage.GetDomainDates(domain)
	if err != nil {
		tracing.TraceErr(span, fmt.Errorf("Cannot determine domain dates: %v", err))
		return 0
	}
	domainAgeInDays := domainDates.CreationAge

	switch {
	case domainAgeInDays <= 1:
		return 75
	case domainAgeInDays <= 7:
		return 60
	case domainAgeInDays <= 10:
		return 50
	case domainAgeInDays <= 15:
		return 40
	case domainAgeInDays <= 30:
		return 30
	case domainAgeInDays <= 90:
		return 15
	default:
		return 0
	}
}

func (s *mailboxService) blacklistPenaltyPercent(domain string) int {

	blacklists := blscan.ScanBlacklists

	pct := (blacklists.MajorLists * .8) + (blacklists.MinorLists * .1) + (blacklists.SpamTrapLists * .25)
	if pct > 1 {
		return 100
	} else {
		return int(pct * 100)
	}

}
