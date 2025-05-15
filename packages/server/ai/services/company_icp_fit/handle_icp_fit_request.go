package company_icp_fit

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/ai/internal/telemetry"
)

type WebageContent struct {
	Homepage string
}

type CompanyMetadata struct {
	Name          string
	Description   string
	Website       string
	EmployeeCount int
	IndustryNAICS string
}

func (s *icpFitAnalyzer) handleICPFitRequest(ctx context.Context, msg *nats.Msg) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpFitAnalyzer.handleICPFitRequest")
	defer span.Finish()

	webpageContent, err := s.getScrapedWebpagesForCompany(ctx, message.PrimaryDomain)
	if err != nil {
		span.TraceError(err)
		return err
	}

	companyMetadata, err := s.getGlobalOrgMetadataForCompany(ctx, message.PrimaryDomain)
	if err != nil {
		span.TraceError(err)
		return err
	}

	icpFitDescription, err := s.getICPFitDescription(ctx)
	if err != nil {
		span.TraceError(err)
		return err
	}

	return nil
}

func (s *icpFitAnalyzer) getICPFitDescription(ctx context.Context) (*models.ICPDescription, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpFitAnalyzer.getICPFitDescription")
	defer span.Finish()

	// TODO

	return nil, nil
}

func (s *icpFitAnalyzer) getGlobalOrgMetadataForCompany(ctx context.Context, domain string) (*CompanyMetadata, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpFitAnalyzer.getGlobalOrgMetadataForCompany")
	defer span.Finish()

	// TODO

	return nil, nil
}

func (s *icpFitAnalyzer) getScrapedWebpagesForCompany(ctx context.Context, domain string) (*WebageContent, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpFitAnalyzer.getScrapedWebpagesForCompany")
	defer span.Finish()

	// TODO

	return nil, nil
}
