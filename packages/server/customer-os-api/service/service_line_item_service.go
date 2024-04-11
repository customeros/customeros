package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ServiceLineItemCreateData struct {
	ContractId        string                            `json:"contractId"`
	SliName           string                            `json:"sliName"`
	SliPrice          float64                           `json:"sliPrice"`
	SliQuantity       int64                             `json:"sliQuantity"`
	SliBilledType     neo4jenum.BilledType              `json:"sliBilledType"`
	ExternalReference *neo4jentity.ExternalSystemEntity `json:"externalReference"`
	Source            neo4jentity.DataSource            `json:"source"`
	AppSource         string                            `json:"appSource"`
	StartedAt         *time.Time                        `json:"startedAt"`
	EndedAt           *time.Time                        `json:"endedAt"`
	SliVatRate        float64                           `json:"sliVatRate"`
}

type ServiceLineItemUpdateData struct {
	Id                      string                 `json:"id"`
	IsRetroactiveCorrection bool                   `json:"isRetroactiveCorrection"`
	SliName                 string                 `json:"sliName"`
	SliPrice                float64                `json:"sliPrice"`
	SliQuantity             int64                  `json:"sliQuantity"`
	SliBilledType           neo4jenum.BilledType   `json:"sliBilledType"`
	SliComments             string                 `json:"sliComments"`
	Source                  neo4jentity.DataSource `json:"source"`
	AppSource               string                 `json:"appSource"`
	SliVatRate              float64                `json:"sliVatRate"`
	StartedAt               *time.Time             `json:"startedAt"`
}

type ServiceLineItemNewVersionData struct {
	Id        string                 `json:"id"`
	Name      string                 `json:"sliName"`
	Price     float64                `json:"sliPrice"`
	Quantity  int64                  `json:"sliQuantity"`
	Comments  string                 `json:"sliComments"`
	Source    neo4jentity.DataSource `json:"source"`
	AppSource string                 `json:"appSource"`
	VatRate   float64                `json:"sliVatRate"`
	StartedAt *time.Time             `json:"startedAt"`
}

type ServiceLineItemService interface {
	Create(ctx context.Context, serviceLineItemDetails ServiceLineItemCreateData, bulk bool) (string, error)
	Update(ctx context.Context, serviceLineItemDetails ServiceLineItemUpdateData, bulk bool) error
	Delete(ctx context.Context, serviceLineItemId string) (bool, error)
	GetById(ctx context.Context, id string) (*neo4jentity.ServiceLineItemEntity, error)
	GetServiceLineItemsForContracts(ctx context.Context, contractIds []string) (*neo4jentity.ServiceLineItemEntities, error)
	Close(ctx context.Context, serviceLineItemId string, endedAt *time.Time, bulk bool) error
	CreateOrUpdateOrCloseInBulk(ctx context.Context, contractId string, sliBulkData []*ServiceLineItemDetails) ([]string, error)
	NewVersion(ctx context.Context, data ServiceLineItemNewVersionData) (string, error)
}

type serviceLineItemService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewServiceLineItemService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) ServiceLineItemService {
	return &serviceLineItemService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *serviceLineItemService) Create(ctx context.Context, serviceLineItemDetails ServiceLineItemCreateData, bulk bool) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItem.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "serviceLineItemDetails", serviceLineItemDetails)

	// check that quantity and price are not negative
	if serviceLineItemDetails.SliQuantity < 0 || serviceLineItemDetails.SliPrice < 0 {
		err := errors.New("quantity and price must not be negative")
		tracing.TraceErr(span, err)
		return "", err
	}

	createServiceLineItemRequest := servicelineitempb.CreateServiceLineItemGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		ContractId:     serviceLineItemDetails.ContractId,
		Name:           serviceLineItemDetails.SliName,
		Quantity:       serviceLineItemDetails.SliQuantity,
		Price:          serviceLineItemDetails.SliPrice,
		VatRate:        serviceLineItemDetails.SliVatRate,
		StartedAt:      utils.ConvertTimeToTimestampPtr(serviceLineItemDetails.StartedAt),
		EndedAt:        utils.ConvertTimeToTimestampPtr(serviceLineItemDetails.EndedAt),
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    string(serviceLineItemDetails.Source),
			AppSource: utils.StringFirstNonEmpty(serviceLineItemDetails.AppSource, constants.AppSourceCustomerOsApi),
		},
	}

	billedType, err := convertBilledTypeToProto(serviceLineItemDetails.SliBilledType, span)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	createServiceLineItemRequest.Billed = billedType

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*servicelineitempb.ServiceLineItemIdGrpcResponse](func() (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
		return s.grpcClients.ServiceLineItemClient.CreateServiceLineItem(ctx, &createServiceLineItemRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	WaitForNodeCreatedInNeo4j(ctx, s.repositories, response.Id, neo4jutil.NodeLabelServiceLineItem, span)

	if !bulk {
		go func() {
			time.Sleep(3 * time.Second)

			contractNode, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, common.GetTenantFromContext(ctx), response.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error on getting contract by service line item id {%s}: %s", response.Id, err.Error())
				return
			}
			contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractNode)
			err = s.generateNextPreviewInvoice(ctx, common.GetTenantFromContext(ctx), contractEntity, span)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error on generating next preview invoice: %s", err.Error())
				return
			}
		}()
	}

	span.LogFields(log.String("output - createdServiceLineItemId", response.Id))
	return response.Id, nil
}

func (s *serviceLineItemService) NewVersion(ctx context.Context, data ServiceLineItemNewVersionData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItem.NewVersion")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "serviceLineItemDetails", data)

	if data.Id == "" {
		err := fmt.Errorf("(ServiceLineItemService.NewVersion) contract line item id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return "", err
	}

	baseServiceLineItemEntity, err := s.GetById(ctx, data.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract line item by id {%s}: %s", data.Id, err.Error())
		return "", err
	}
	contractEntity, err := s.services.ContractService.GetContractByServiceLineItem(ctx, data.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", data.Id, err.Error())
		return "", err
	}

	startedAt := utils.ToDate(utils.IfNotNilTimeWithDefault(data.StartedAt, utils.Now()))

	// Validate new version creation
	if baseServiceLineItemEntity.Billed == neo4jenum.BilledTypeOnce {
		err = fmt.Errorf("cannot create new version for one time contract line item with id {%s}", data.Id)
		tracing.TraceErr(span, err)
		return "", err
	}
	// Do not allow creating new version if there is an existing version with the same start date
	serviceLineItems, err := s.GetServiceLineItemsForContracts(ctx, []string{contractEntity.Id})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting service line items for contract {%s}: %s", contractEntity.Id, err.Error())
		return "", err
	}
	for _, sli := range *serviceLineItems {
		if utils.ToDate(sli.StartedAt).Equal(startedAt) {
			err = fmt.Errorf("contract line item with id {%s} already exists with the same start date {%s}", sli.ID, startedAt.Format(time.DateOnly))
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	// If contract was invoiced - do not allow creating new version in the past
	contractInvoiced, err := s.repositories.Neo4jRepositories.ContractReadRepository.IsContractInvoiced(ctx, common.GetTenantFromContext(ctx), contractEntity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on checking if contract was invoiced: %s", err.Error())
		return "", err
	}
	if contractInvoiced && startedAt.Before(utils.Today()) {
		err = fmt.Errorf("cannot create new version for contract line item with id {%s} in the past", data.Id)
		tracing.TraceErr(span, err)
		return "", err
	}

	createServiceLineItemRequest := servicelineitempb.CreateServiceLineItemGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		ContractId:     contractEntity.Id,
		ParentId:       baseServiceLineItemEntity.ParentID,
		Name:           utils.StringFirstNonEmpty(data.Name, baseServiceLineItemEntity.Name),
		Quantity:       data.Quantity,
		Price:          data.Price,
		VatRate:        data.VatRate,
		StartedAt:      utils.ConvertTimeToTimestampPtr(&startedAt),
		Comments:       utils.IfNotNilString(data.Comments),
		SourceFields: &commonpb.SourceFields{
			Source:    data.Source.String(),
			AppSource: utils.StringFirstNonEmpty(data.AppSource, constants.AppSourceCustomerOsApi),
		},
	}

	billedType, err := convertBilledTypeToProto(baseServiceLineItemEntity.Billed, span)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	createServiceLineItemRequest.Billed = billedType

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*servicelineitempb.ServiceLineItemIdGrpcResponse](func() (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
		return s.grpcClients.ServiceLineItemClient.CreateServiceLineItem(ctx, &createServiceLineItemRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	go func() {
		time.Sleep(2 * time.Second)
		err = s.generateNextPreviewInvoice(ctx, common.GetTenantFromContext(ctx), contractEntity, span)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error on generating next preview invoice: %s", err.Error())
			return
		}
	}()

	return response.Id, err
}

func (s *serviceLineItemService) Update(ctx context.Context, serviceLineItemDetails ServiceLineItemUpdateData, bulk bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "serviceLineItemDetails", serviceLineItemDetails)

	if serviceLineItemDetails.Id == "" {
		err := fmt.Errorf("(ServiceLineItemService.Update) contract line item id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	baseServiceLineItemEntity, err := s.GetById(ctx, serviceLineItemDetails.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract line item by id {%s}: %s", serviceLineItemDetails.Id, err.Error())
		return err
	}

	if baseServiceLineItemEntity.Canceled {
		err = fmt.Errorf("service line item with id {%s} is already ended", serviceLineItemDetails.Id)
		tracing.TraceErr(span, err)
		return err
	}

	contractEntity, err := s.services.ContractService.GetContractByServiceLineItem(ctx, serviceLineItemDetails.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", serviceLineItemDetails.Id, err.Error())
		return err
	}

	contractInvoiced, err := s.repositories.Neo4jRepositories.ContractReadRepository.IsContractInvoiced(ctx, common.GetTenantFromContext(ctx), contractEntity.Id)
	startedAt := utils.IfNotNilTimeWithDefault(serviceLineItemDetails.StartedAt, baseServiceLineItemEntity.StartedAt)

	// Do not allow updating past SLIs for invoiced contracts
	if contractInvoiced && startedAt.Before(utils.Today()) {
		err = fmt.Errorf("cannot update contract line item with id {%s} in the past", serviceLineItemDetails.Id)
		tracing.TraceErr(span, err)
		return err
	}

	// check that quantity and price are not negative
	if serviceLineItemDetails.SliQuantity < 0 || serviceLineItemDetails.SliPrice < 0 {
		err := errors.New("quantity and price must not be negative")
		tracing.TraceErr(span, err)
		return err
	}

	// check that billing cycle is not changed
	if baseServiceLineItemEntity.Billed.String() != serviceLineItemDetails.SliBilledType.String() && baseServiceLineItemEntity.Billed.String() != "" {
		err = fmt.Errorf("cannot change billing cycle for contract line item with id {%s}", serviceLineItemDetails.Id)
		tracing.TraceErr(span, err)
		return err
	}

	isRetroactiveCorrection := serviceLineItemDetails.IsRetroactiveCorrection
	if baseServiceLineItemEntity.IsOneTime() {
		isRetroactiveCorrection = true
	}

	// If no price impacted fields changed, set retroactive correction to true
	if baseServiceLineItemEntity.Price == serviceLineItemDetails.SliPrice &&
		baseServiceLineItemEntity.Quantity == serviceLineItemDetails.SliQuantity &&
		baseServiceLineItemEntity.VatRate == serviceLineItemDetails.SliVatRate {
		isRetroactiveCorrection = true
	}
	if startedAt == *serviceLineItemDetails.StartedAt {
		isRetroactiveCorrection = true
	}

	// Check SLI is not invoiced
	sliInvoiced, err := s.repositories.Neo4jRepositories.ServiceLineItemReadRepository.WasServiceLineItemInvoiced(ctx, common.GetTenantFromContext(ctx), baseServiceLineItemEntity.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on checking if service line item was invoiced: %s", err.Error())
		return err
	}

	if sliInvoiced && !isRetroactiveCorrection {
		err = fmt.Errorf("service line item with id {%s} is included in invoice and cannot be updated", serviceLineItemDetails.Id)
		tracing.TraceErr(span, err)
		return err
	}

	if isRetroactiveCorrection {
		serviceLineItemUpdateRequest := servicelineitempb.UpdateServiceLineItemGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			Id:             serviceLineItemDetails.Id,
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			Name:           serviceLineItemDetails.SliName,
			Quantity:       serviceLineItemDetails.SliQuantity,
			Price:          serviceLineItemDetails.SliPrice,
			Comments:       serviceLineItemDetails.SliComments,
			VatRate:        serviceLineItemDetails.SliVatRate,
			SourceFields: &commonpb.SourceFields{
				Source:    string(serviceLineItemDetails.Source),
				AppSource: utils.StringFirstNonEmpty(serviceLineItemDetails.AppSource, constants.AppSourceCustomerOsApi),
			},
		}

		billedType, err := convertBilledTypeToProto(serviceLineItemDetails.SliBilledType, span)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		serviceLineItemUpdateRequest.Billed = billedType

		if baseServiceLineItemEntity.ParentID == baseServiceLineItemEntity.ID {
			serviceLineItemUpdateRequest.StartedAt = utils.ConvertTimeToTimestampPtr(serviceLineItemDetails.StartedAt)
		}
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*servicelineitempb.ServiceLineItemIdGrpcResponse](func() (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
			return s.grpcClients.ServiceLineItemClient.UpdateServiceLineItem(ctx, &serviceLineItemUpdateRequest)
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error from events processing: %s", err.Error())
			return err
		}
	} else {
		createServiceLineItemRequest := servicelineitempb.CreateServiceLineItemGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			ContractId:     contractEntity.Id,
			ParentId:       baseServiceLineItemEntity.ParentID,
			Name:           utils.StringFirstNonEmpty(serviceLineItemDetails.SliName, baseServiceLineItemEntity.Name),
			Quantity:       serviceLineItemDetails.SliQuantity,
			Price:          serviceLineItemDetails.SliPrice,
			VatRate:        serviceLineItemDetails.SliVatRate,
			Comments:       utils.IfNotNilString(serviceLineItemDetails.SliComments),
			StartedAt:      utils.ConvertTimeToTimestampPtr(serviceLineItemDetails.StartedAt),
			SourceFields: &commonpb.SourceFields{
				Source:    serviceLineItemDetails.Source.String(),
				AppSource: utils.StringFirstNonEmpty(serviceLineItemDetails.AppSource, constants.AppSourceCustomerOsApi),
			},
		}

		billedType, err := convertBilledTypeToProto(baseServiceLineItemEntity.Billed, span)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		createServiceLineItemRequest.Billed = billedType

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*servicelineitempb.ServiceLineItemIdGrpcResponse](func() (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
			return s.grpcClients.ServiceLineItemClient.CreateServiceLineItem(ctx, &createServiceLineItemRequest)
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	if !bulk {
		go func() {
			time.Sleep(3 * time.Second)

			err = s.generateNextPreviewInvoice(ctx, common.GetTenantFromContext(ctx), contractEntity, span)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error on generating next preview invoice: %s", err.Error())
				return
			}
		}()
	}

	return nil
}

func (s *serviceLineItemService) Delete(ctx context.Context, serviceLineItemId string) (completed bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Delete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	sliExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), serviceLineItemId, neo4jutil.NodeLabelServiceLineItem)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error on checking if service line item exists: %s", err.Error())
		return false, err
	}
	if !sliExists {
		err := fmt.Errorf("service line item with id {%s} not found", serviceLineItemId)
		tracing.TraceErr(span, err)
		s.log.Errorf(err.Error())
		return false, err
	}

	// Check SLI is not invoiced
	sliInvoiced, err := s.repositories.Neo4jRepositories.ServiceLineItemReadRepository.WasServiceLineItemInvoiced(ctx, common.GetTenantFromContext(ctx), serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on checking if service line item was invoiced: %s", err.Error())
		return false, err
	}
	if sliInvoiced {
		err := fmt.Errorf("service line item with id {%s} is included in invoice and cannot be deleted", serviceLineItemId)
		tracing.TraceErr(span, err)
		s.log.Errorf(err.Error())
		return false, err
	}

	contractEntity, err := s.services.ContractService.GetContractByServiceLineItem(ctx, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", serviceLineItemId, err.Error())
		return false, err
	}

	deleteRequest := servicelineitempb.DeleteServiceLineItemGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             serviceLineItemId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		AppSource:      constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*servicelineitempb.ServiceLineItemIdGrpcResponse](func() (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
		return s.grpcClients.ServiceLineItemClient.DeleteServiceLineItem(ctx, &deleteRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return false, err
	}

	go func() {
		time.Sleep(3 * time.Second)
		err = s.generateNextPreviewInvoice(ctx, common.GetTenantFromContext(ctx), contractEntity, span)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error on generating next preview invoice: %s", err.Error())
			return
		}
	}()

	return false, nil
}

func (s *serviceLineItemService) Close(ctx context.Context, serviceLineItemId string, endedAt *time.Time, bulk bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Close")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))
	span.SetTag(tracing.SpanTagEntityId, serviceLineItemId)

	sli, err := s.GetById(ctx, serviceLineItemId)

	contractEntity, err := s.services.ContractService.GetContractByServiceLineItem(ctx, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", serviceLineItemId, err.Error())
		return err
	}

	contractInvoiced, err := s.repositories.Neo4jRepositories.ContractReadRepository.IsContractInvoiced(ctx, common.GetTenantFromContext(ctx), contractEntity.Id)

	if contractInvoiced || sli.IsActiveAt(utils.Now()) {
		closeRequest := servicelineitempb.CloseServiceLineItemGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			Id:             serviceLineItemId,
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			AppSource:      constants.AppSourceCustomerOsApi,
			EndedAt:        utils.ConvertTimeToTimestampPtr(endedAt),
		}

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*servicelineitempb.ServiceLineItemIdGrpcResponse](func() (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
			return s.grpcClients.ServiceLineItemClient.CloseServiceLineItem(ctx, &closeRequest)
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error from events processing: %s", err.Error())
			return err
		}
	} else {
		_, err = s.Delete(ctx, serviceLineItemId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error on deleting service line item: %s", err.Error())
			return err
		}
	}

	if !bulk {
		go func() {
			time.Sleep(3 * time.Second)
			err = s.generateNextPreviewInvoice(ctx, common.GetTenantFromContext(ctx), contractEntity, span)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error on generating next preview invoice: %s", err.Error())
				return
			}
		}()
	}

	return nil
}

func (s *serviceLineItemService) GetById(ctx context.Context, serviceLineItemId string) (*neo4jentity.ServiceLineItemEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	if sliDbNode, err := s.repositories.ServiceLineItemRepository.GetById(ctx, common.GetContext(ctx).Tenant, serviceLineItemId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("service line item with id {%s} not found", serviceLineItemId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode), nil
	}
}

func (s *serviceLineItemService) GetServiceLineItemsForContracts(ctx context.Context, contractIDs []string) (*neo4jentity.ServiceLineItemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.GetServiceLineItemsForContracts")
	defer span.Finish()
	span.LogFields(log.Object("contractIDs", contractIDs))

	serviceLineItems, err := s.repositories.ServiceLineItemRepository.GetForContracts(ctx, common.GetTenantFromContext(ctx), contractIDs)
	if err != nil {
		return nil, err
	}
	serviceLineItemEntities := make(neo4jentity.ServiceLineItemEntities, 0, len(serviceLineItems))
	for _, v := range serviceLineItems {
		serviceLineItemEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(v.Node)
		serviceLineItemEntity.DataloaderKey = v.LinkedNodeId
		serviceLineItemEntities = append(serviceLineItemEntities, *serviceLineItemEntity)
	}
	return &serviceLineItemEntities, nil
}

type ServiceLineItemDetails struct {
	Id                      string
	Name                    string
	Price                   float64
	Quantity                int64
	Billed                  neo4jenum.BilledType
	Comments                string
	IsRetroactiveCorrection bool
	VatRate                 float64
	StartedAt               *time.Time
}

func (s *serviceLineItemService) CreateOrUpdateOrCloseInBulk(ctx context.Context, contractId string, sliBulkData []*ServiceLineItemDetails) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.CreateOrUpdateOrCloseInBulk")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contractId", contractId))

	var responseIds []string

	allSliDbNodes, err := s.repositories.ServiceLineItemRepository.GetForContracts(ctx, common.GetTenantFromContext(ctx), []string{contractId})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get service line items for contract: %s", err.Error())
		return []string{}, err
	}

	var sliDbNodes []*utils.DbNodeAndId
	for _, sliDbNode := range allSliDbNodes {
		endedAt := utils.GetTimePropOrNil(utils.GetPropsFromNode(*sliDbNode.Node), "endedAt")
		if endedAt == nil {
			sliDbNodes = append(sliDbNodes, sliDbNode)
		}
	}

	var existingSliIds, inputSliIds []string
	for _, sliDbNode := range sliDbNodes {
		id := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*sliDbNode.Node), "id")
		if id != "" {
			existingSliIds = append(existingSliIds, id)
		}
	}
	for _, sli := range sliBulkData {
		if sli.Id != "" {
			inputSliIds = append(inputSliIds, sli.Id)
		}
	}
	for _, existingId := range existingSliIds {
		if !utils.Contains(inputSliIds, existingId) {
			err = s.Close(ctx, existingId, nil, true)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Failed to close service line item: %s", err.Error())
			}
		}
	}

	for _, serviceLineItem := range sliBulkData {
		// check all quantity or price are not negative
		if serviceLineItem.Quantity < 0 || serviceLineItem.Price < 0 {
			err = fmt.Errorf("quantity and price must not be negative")
			tracing.TraceErr(span, err)
			return []string{}, err
		}
	}

	for _, serviceLineItem := range sliBulkData {
		if serviceLineItem.Id == "" {
			itemId, err := s.Create(ctx, ServiceLineItemCreateData{
				ContractId:    contractId,
				SliName:       serviceLineItem.Name,
				SliPrice:      serviceLineItem.Price,
				SliQuantity:   serviceLineItem.Quantity,
				SliBilledType: serviceLineItem.Billed,
				SliVatRate:    serviceLineItem.VatRate,
				Source:        neo4jentity.DataSourceOpenline,
				AppSource:     constants.AppSourceCustomerOsApi,
				StartedAt:     serviceLineItem.StartedAt,
			}, true)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error from events processing: %s", err.Error())
				return []string{}, err
			}
			responseIds = append(responseIds, itemId)

		} else {
			responseIds = append(responseIds, serviceLineItem.Id)
			err = s.Update(ctx, ServiceLineItemUpdateData{
				Id:                      serviceLineItem.Id,
				IsRetroactiveCorrection: serviceLineItem.IsRetroactiveCorrection,
				SliName:                 serviceLineItem.Name,
				SliPrice:                serviceLineItem.Price,
				SliQuantity:             serviceLineItem.Quantity,
				SliBilledType:           serviceLineItem.Billed,
				SliComments:             serviceLineItem.Comments,
				SliVatRate:              serviceLineItem.VatRate,
				Source:                  neo4jentity.DataSourceOpenline,
				AppSource:               constants.AppSourceCustomerOsApi,
				StartedAt:               serviceLineItem.StartedAt,
			}, true)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error from events processing: %s", err.Error())
				return []string{}, err
			}
		}
	}

	go func() {
		time.Sleep(3 * time.Second)

		contractDbNode, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, common.GetTenantFromContext(ctx), contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error on getting contract by id {%s}: %s", contractId, err.Error())
			return
		}
		contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

		err = s.generateNextPreviewInvoice(ctx, common.GetTenantFromContext(ctx), contractEntity, span)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error on generating next preview invoice: %s", err.Error())
			return
		}
	}()

	return responseIds, nil
}

func (s *serviceLineItemService) generateNextPreviewInvoice(ctx context.Context, tenant string, contractEntity *neo4jentity.ContractEntity, span opentracing.Span) error {
	if contractEntity == nil || contractEntity.Id == "" {
		err := errors.New("contract entity is nil or contract id is empty")
		tracing.TraceErr(span, err)
		return err
	}
	if contractEntity.InvoicingEnabled && contractEntity.BillingCycle != neo4jenum.BillingCycleNone && contractEntity.InvoicingStartDate != nil {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err := utils.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
			return s.grpcClients.InvoiceClient.NextPreviewInvoiceForContract(ctx, &invoicepb.NextPreviewInvoiceForContractRequest{
				Tenant:     tenant,
				ContractId: contractEntity.Id,
				AppSource:  constants.AppSourceCustomerOsApi,
			})
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func MapServiceLineItemBulkItemsToData(input []*model.ServiceLineItemBulkUpdateItem) []*ServiceLineItemDetails {
	var arr []*ServiceLineItemDetails
	for _, item := range input {
		sli := MapServiceLineItemBulkItemToData(item)
		if sli != nil {
			arr = append(arr, sli)
		}
	}
	return arr
}

func MapServiceLineItemBulkItemToData(input *model.ServiceLineItemBulkUpdateItem) *ServiceLineItemDetails {
	if input == nil {
		return nil
	}
	billed := neo4jenum.BilledTypeNone
	if input.Billed != nil {
		billed = mapper.MapBilledTypeFromModel(*input.Billed)
	}
	return &ServiceLineItemDetails{
		Id:                      utils.IfNotNilString(input.ServiceLineItemID),
		Name:                    utils.IfNotNilString(input.Name),
		Price:                   utils.IfNotNilFloat64(input.Price),
		Quantity:                utils.IfNotNilInt64(input.Quantity),
		Billed:                  billed,
		Comments:                utils.IfNotNilString(input.Comments),
		IsRetroactiveCorrection: utils.IfNotNilBool(input.IsRetroactiveCorrection),
		VatRate:                 utils.IfNotNilFloat64(input.VatRate),
		StartedAt:               input.ServiceStarted,
	}
}

func convertBilledTypeToProto(billedType neo4jenum.BilledType, span opentracing.Span) (commonpb.BilledType, error) {
	switch billedType {
	case neo4jenum.BilledTypeMonthly:
		return commonpb.BilledType_MONTHLY_BILLED, nil
	case neo4jenum.BilledTypeQuarterly:
		return commonpb.BilledType_QUARTERLY_BILLED, nil
	case neo4jenum.BilledTypeAnnually:
		return commonpb.BilledType_ANNUALLY_BILLED, nil
	case neo4jenum.BilledTypeOnce:
		return commonpb.BilledType_ONCE_BILLED, nil
	case neo4jenum.BilledTypeUsage:
		return commonpb.BilledType_USAGE_BILLED, nil
	case neo4jenum.BilledTypeNone:
		err := fmt.Errorf("billed type is not set")
		tracing.TraceErr(span, err)
		return commonpb.BilledType_NONE_BILLED, err
	default:
		err := fmt.Errorf("unknown billed type: %s", billedType)
		tracing.TraceErr(span, err)
		return commonpb.BilledType_NONE_BILLED, err
	}
}
