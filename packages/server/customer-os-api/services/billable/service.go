package api_billable

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type billableService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewBillableService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.BillableService {
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
