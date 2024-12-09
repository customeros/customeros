package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
)

func (s *trackingService) notifyOnSlack(c context.Context, r *entity.Tracking) error {
	span, ctx := opentracing.StartSpanFromContext(c, "TrackingService.notifyOnSlack")
	defer span.Finish()

	record, err := s.getAndUpdateTrackingRecord(ctx, span, r.ID)
	if err != nil {
		return err
	}

	if record.Notified || record.OrganizationId == nil {
		return nil
	}

	globalOrg, err := s.getEnrichedOrganizationDetails(ctx, span, record)
	if err != nil {
		return err
	}

	if shouldSkipNotification(ctx, span, record) {
		return nil
	}

	slackBlock := s.buildSlackNotification(record, globalOrg)
	return s.sendNotifications(ctx, span, record, slackBlock)
}

func (s *trackingService) getAndUpdateTrackingRecord(ctx context.Context, span opentracing.Span, id string) (*entity.Tracking, error) {
	record, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get tracking record"))
		return nil, err
	}

	if record.Tenant != "" {
		tracing.TagTenant(span, record.Tenant)
	}

	err = s.services.CommonServices.PostgresRepositories.TrackingRepository.IncrementNotificationTry(ctx, record.ID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to increment notification try"))
	}

	return record, nil
}

func (s *trackingService) getEnrichedOrganizationDetails(ctx context.Context, span opentracing.Span, record *entity.Tracking) (*entity.GlobalOrganization, error) {
	enrichDetails, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, record.IP)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	_, primaryDomain := domaincheck.PrimaryDomainCheck(*enrichDetails.CompanyDomain)
	return s.services.CommonServices.PostgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
}

func (s *trackingService) shouldSkipNotification(ctx context.Context, span opentracing.Span, record *entity.Tracking) bool {
	// Skip if domain matches workspace
	if record.OrganizationDomain != nil && *record.OrganizationDomain != "" {
		if s.isWorkspaceDomain(ctx, span, record) {
			return true
		}

		// Skip if recently notified
		notificationSent, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.WasNotifiedRecently(ctx, *record.OrganizationDomain, 24)
		if err != nil {
			tracing.TraceErr(span, err)
			return true
		}
		return notificationSent
	}
	return false
}

func (s *trackingService) isWorkspaceDomain(ctx context.Context, span opentracing.Span, record *entity.Tracking) bool {
	workspaceNodeList, err := s.services.CommonServices.Neo4jRepositories.WorkspaceReadRepository.GetAllForTenant(ctx, record.Tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get workspace nodes"))
		return true
	}

	for _, workspaceNode := range workspaceNodeList {
		props := utils.GetPropsFromNode(*workspaceNode)
		if utils.GetStringPropOrEmpty(props, "name") == *record.OrganizationDomain {
			err := s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsNotified(ctx, record.ID)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			span.LogFields(log.String("skip", "workspace is the same as organization domain"))
			return true
		}
	}
	return false
}

func (s *trackingService) buildSlackNotification(record *entity.Tracking, globalOrg *entity.GlobalOrganization) string {
	referrer := record.Referrer
	if referrer == "" {
		referrer = "Direct"
	}

	slackBlock := `[{"type":"header",...}]` // Your existing slack block template

	replacements := map[string]string{
		"{placeholder_organization_name}":     globalOrg.Name,
		"{placeholder_location}":              organizationLocation,
		"{placeholder_website}":               globalOrg.Website,
		"{placeholder_linkedin}":              globalOrg.LinkedInUrl,
		"{placeholder_referrer}":              referrer,
		"{placeholder_view_organization_url}": fmt.Sprintf("https://app.customeros.ai/organization/%s?tab=about", *record.OrganizationId),
	}

	for placeholder, value := range replacements {
		slackBlock = strings.Replace(slackBlock, placeholder, value, -1)
	}

	return slackBlock
}

func (s *trackingService) sendNotifications(ctx context.Context, span opentracing.Span, record *entity.Tracking, slackBlock string) error {
	slackChannels, err := s.services.CommonServices.PostgresRepositories.SlackChannelNotificationRepository.GetSlackChannels(ctx, record.Tenant, "REVEAL-AI")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, slackChannel := range slackChannels {
		if slackChannel.CreatedAt.After(record.CreatedAt) {
			err := s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsNotified(ctx, record.ID)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			continue
		}

		if err := s.sendSlackMessage(ctx, slackChannel.Tenant, slackChannel.ChannelId, slackBlock); err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if err := s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsNotified(ctx, record.ID); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}
	return nil
}
