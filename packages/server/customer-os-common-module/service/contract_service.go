package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ContractService interface {
	GetById(ctx context.Context, id string) (*neo4jentity.ContractEntity, error)
}
type contractService struct {
	log      logger.Logger
	services *Services
}

func NewContractService(log logger.Logger, services *Services) ContractService {
	return &contractService{
		log:      log,
		services: services,
	}
}

func (s *contractService) GetById(ctx context.Context, contractId string) (*neo4jentity.ContractEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contractId", contractId))

	if contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, common.GetContext(ctx).Tenant, contractId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Contract with id {%s} not found", contractId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToContractEntity(contractDbNode), nil
	}
}
