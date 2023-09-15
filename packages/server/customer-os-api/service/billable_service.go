package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type BillableService interface {
	GetBillableDetails(ctx context.Context) (*model.TenantBillableInfo, error)
}

type billableService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewBillableService(log logger.Logger, repositories *repository.Repositories) BillableService {
	return &billableService{
		log:          log,
		repositories: repositories,
	}
}

func (s *billableService) GetBillableDetails(ctx context.Context) (*model.TenantBillableInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BillableService.GetBillableDetails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbRecord, err := s.repositories.ContactRepository.GetBillableContactStats(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "GetBillableDetails")
	}
	return &model.TenantBillableInfo{
		WhitelistedOrganizations: dbRecord.Values[0].(int64),
		WhitelistedContacts:      dbRecord.Values[1].(int64),
		GreylistedOrganizations:  dbRecord.Values[2].(int64),
		GreylistedContacts:       dbRecord.Values[3].(int64),
	}, nil
}
