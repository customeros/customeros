package icp

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/enums"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	core_errors "github.com/customeros/customeros/packages/server/core-crm/errors"
	"github.com/customeros/customeros/packages/server/core-crm/internal/models"
	"github.com/customeros/customeros/packages/server/core-crm/internal/proto/pb"
	"github.com/customeros/customeros/packages/server/core-crm/internal/telemetry"
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

func (s *icpService) handleCompanyCreated(ctx context.Context, msg *nats.Msg) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpService.handleCompanyCreated")
	defer span.Finish()

	message, err := s.parseCompanyCreatedMessage(ctx, msg)
	if err != nil {
		span.TraceError(err)
		return err
	}

	err = s.requestICPFitAnalysis(ctx, message.PrimaryDomain)
	if err != nil {
		span.TraceError(err)
		return err
	}

	return nil
}

func (s *icpService) requestICPFitAnalysis(ctx context.Context, domain string) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpService.requestICPFitAnalysis")
	defer span.Finish()

	// TODO

	return nil
}

func (s *icpService) parseCompanyCreatedMessage(ctx context.Context, msg *nats.Msg) (*pb.CompanyCreated, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpService.parseCompanyCreatedMessage")
	defer span.Finish()

	message := &pb.CompanyCreated{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	if message.CompanyId == "" && message.OrganizationId == "" {
		err := errors.New("CompanyID not provided")
		span.TraceError(err)
		return nil, err
	}

	if message.Tenant == "" {
		err := core_errors.ErrTenantMissing
		span.TraceError(err)
		return nil, err
	}

	if message.PrimaryDomain == "" {
		err := errors.New("Primary Domain not provided")
		span.TraceError(err)
		return nil, err
	}

	return message, nil
}
