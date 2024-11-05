package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ContractService interface {
	GetById(ctx context.Context, id string) (*neo4jentity.ContractEntity, error)
	Save(ctx context.Context, id *string, dataFields data_fields.ContractSaveFields) (string, error)
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

func (s *contractService) Save(ctx context.Context, id *string, dataFields data_fields.ContractSaveFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	contractId := ""

	if id == nil || *id == "" {
		createFlow = true
		span.LogKV("flow", "create")
		contractId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelContract)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		// set missing fields for create flow
		if dataFields.CreatedAt == nil {
			dataFields.CreatedAt = utils.NowPtr()
		} else {
			dataFields.CreatedAt = utils.TimePtr(utils.NowIfZero(*dataFields.CreatedAt))
		}
		if utils.IfNotNilString(dataFields.AppSource) == "" {
			dataFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}
		if utils.IfNotNilString(dataFields.Source) == "" {
			dataFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
	} else {
		span.LogKV("flow", "update")
		contractId = *id

		// validate contract exists
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, contractId, model.NodeLabelContract)
		if err != nil || !exists {
			err = errors.New("contract not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, contractId)

	if createFlow {
		err := s.services.Neo4jRepositories.ContractWriteRepository.CreateForOrganization(ctx, tenant, contractId, dataFields)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	} else {
		err := s.services.Neo4jRepositories.ContractWriteRepository.UpdateContract(ctx, tenant, contractId, dataFields)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	if createFlow {
		err = s.services.RabbitMQService.PublishEvent(ctx, contractId, model.CONTRACT, dto.CreateContract{dataFields})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateContract"))
		}
		utils.EventCompleted(ctx, tenant, model.CONTRACT.String(), contractId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
	} else {
		err = s.services.RabbitMQService.PublishEvent(ctx, contractId, model.CONTRACT, dto.UpdateContract{dataFields})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateContractOld"))
		}
		if dataFields.AppSource == nil || *dataFields.AppSource != constants.AppSourceCustomerOsApi {
			utils.EventCompleted(ctx, tenant, model.CONTRACT.String(), contractId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
		}
	}

	return contractId, nil
}
