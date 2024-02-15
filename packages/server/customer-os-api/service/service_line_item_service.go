package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
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
	ContractId        string                       `json:"contractId"`
	SliName           string                       `json:"sliName"`
	SliPrice          float64                      `json:"sliPrice"`
	SliQuantity       int64                        `json:"sliQuantity"`
	SliBilledType     neo4jenum.BilledType         `json:"sliBilledType"`
	ExternalReference *entity.ExternalSystemEntity `json:"externalReference"`
	Source            neo4jentity.DataSource       `json:"source"`
	AppSource         string                       `json:"appSource"`
	StartedAt         *time.Time                   `json:"startedAt"`
	EndedAt           *time.Time                   `json:"endedAt"`
	SliVatRate        float64                      `json:"sliVatRate"`
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
}

type ServiceLineItemService interface {
	Create(ctx context.Context, serviceLineItemDetails ServiceLineItemCreateData) (string, error)
	Update(ctx context.Context, serviceLineItemDetails ServiceLineItemUpdateData) error
	Delete(ctx context.Context, serviceLineItemId string) (bool, error)
	GetById(ctx context.Context, id string) (*neo4jentity.ServiceLineItemEntity, error)
	GetServiceLineItemsForContracts(ctx context.Context, contractIds []string) (*neo4jentity.ServiceLineItemEntities, error)
	Close(ctx context.Context, serviceLineItemId string, endedAt *time.Time) error
	CreateOrUpdateInBulk(ctx context.Context, contractId string, sliBulkData []*ServiceLineItemDetails) ([]string, error)
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

func (s *serviceLineItemService) Create(ctx context.Context, serviceLineItemDetails ServiceLineItemCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItem.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "serviceLineItemDetails", serviceLineItemDetails)

	serviceLineItemId, err := s.createServiceLineItemWithEvents(ctx, &serviceLineItemDetails)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	span.LogFields(log.String("output - createdServiceLineItemId", serviceLineItemId))
	return serviceLineItemId, nil
}

func (s *serviceLineItemService) createServiceLineItemWithEvents(ctx context.Context, serviceLineItemDetails *ServiceLineItemCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.createServiceLineItemWithEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

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

	switch serviceLineItemDetails.SliBilledType {
	case neo4jenum.BilledTypeMonthly:
		createServiceLineItemRequest.Billed = commonpb.BilledType_MONTHLY_BILLED
	case neo4jenum.BilledTypeQuarterly:
		createServiceLineItemRequest.Billed = commonpb.BilledType_QUARTERLY_BILLED
	case neo4jenum.BilledTypeAnnually:
		createServiceLineItemRequest.Billed = commonpb.BilledType_ANNUALLY_BILLED
	case neo4jenum.BilledTypeOnce:
		createServiceLineItemRequest.Billed = commonpb.BilledType_ONCE_BILLED
	case neo4jenum.BilledTypeUsage:
		createServiceLineItemRequest.Billed = commonpb.BilledType_USAGE_BILLED
	default:
		createServiceLineItemRequest.Billed = commonpb.BilledType_NONE_BILLED
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := s.grpcClients.ServiceLineItemClient.CreateServiceLineItem(ctx, &createServiceLineItemRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	WaitForObjectCreationAndLogSpan(ctx, s.repositories, response.Id, neo4jutil.NodeLabelServiceLineItem, span)
	return response.Id, err
}

func (s *serviceLineItemService) Update(ctx context.Context, serviceLineItemDetails ServiceLineItemUpdateData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "serviceLineItemDetails", serviceLineItemDetails)

	if serviceLineItemDetails.Id == "" {
		err := fmt.Errorf("(ServiceLineItemService.Update) service line item id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	serviceLineItem, err := s.GetById(ctx, serviceLineItemDetails.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting service line item by id {%s}: %s", serviceLineItemDetails.Id, err.Error())
		return err
	}

	isRetroactiveCorrection := serviceLineItemDetails.IsRetroactiveCorrection
	// If no price impacted fields changed, set retroactive correction to true
	if serviceLineItem.Price == serviceLineItemDetails.SliPrice &&
		serviceLineItem.Quantity == serviceLineItemDetails.SliQuantity &&
		serviceLineItem.Billed == serviceLineItemDetails.SliBilledType &&
		serviceLineItem.VatRate == serviceLineItemDetails.SliVatRate {
		isRetroactiveCorrection = true
	}

	serviceLineItemUpdateRequest := servicelineitempb.UpdateServiceLineItemGrpcRequest{
		Tenant:                  common.GetTenantFromContext(ctx),
		Id:                      serviceLineItemDetails.Id,
		LoggedInUserId:          common.GetUserIdFromContext(ctx),
		Name:                    serviceLineItemDetails.SliName,
		Quantity:                serviceLineItemDetails.SliQuantity,
		Price:                   serviceLineItemDetails.SliPrice,
		Comments:                serviceLineItemDetails.SliComments,
		IsRetroactiveCorrection: isRetroactiveCorrection,
		VatRate:                 serviceLineItemDetails.SliVatRate,
		SourceFields: &commonpb.SourceFields{
			Source:    string(serviceLineItemDetails.Source),
			AppSource: utils.StringFirstNonEmpty(serviceLineItemDetails.AppSource, constants.AppSourceCustomerOsApi),
		},
	}
	switch serviceLineItemDetails.SliBilledType {
	case neo4jenum.BilledTypeMonthly:
		serviceLineItemUpdateRequest.Billed = commonpb.BilledType_MONTHLY_BILLED
	case neo4jenum.BilledTypeQuarterly:
		serviceLineItemUpdateRequest.Billed = commonpb.BilledType_QUARTERLY_BILLED
	case neo4jenum.BilledTypeAnnually:
		serviceLineItemUpdateRequest.Billed = commonpb.BilledType_ANNUALLY_BILLED
	case neo4jenum.BilledTypeOnce:
		serviceLineItemUpdateRequest.Billed = commonpb.BilledType_ONCE_BILLED
	case neo4jenum.BilledTypeUsage:
		serviceLineItemUpdateRequest.Billed = commonpb.BilledType_USAGE_BILLED
	default:
		serviceLineItemUpdateRequest.Billed = commonpb.BilledType_NONE_BILLED
	}
	// set contract id if it's not a retroactive correction
	if !isRetroactiveCorrection {
		contractDbNode, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, common.GetTenantFromContext(ctx), serviceLineItemDetails.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error on getting contract by service line item id {%s}: %s", serviceLineItemDetails.Id, err.Error())
			return err
		}
		if contractDbNode == nil {
			err := fmt.Errorf("contract not found for service line item id {%s}", serviceLineItemDetails.Id)
			tracing.TraceErr(span, err)
			s.log.Errorf(err.Error())
			return err
		}
		serviceLineItemUpdateRequest.ContractId = utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*contractDbNode), "id")
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = s.grpcClients.ServiceLineItemClient.UpdateServiceLineItem(ctx, &serviceLineItemUpdateRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
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

	deleteRequest := servicelineitempb.DeleteServiceLineItemGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             serviceLineItemId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		AppSource:      constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = s.grpcClients.ServiceLineItemClient.DeleteServiceLineItem(ctx, &deleteRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return false, err
	}

	// wait for service line item to be deleted from graph db
	for i := 1; i <= constants.MaxRetriesCheckDataInNeo4jAfterEventRequest; i++ {
		serviceLineItemFound, findErr := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), serviceLineItemId, neo4jutil.NodeLabelServiceLineItem)
		if findErr != nil {
			tracing.TraceErr(span, findErr)
			s.log.Errorf("error on checking if service line item exists: %s", findErr.Error())
		} else if !serviceLineItemFound {
			span.LogFields(log.Bool("serviceLineItemDeletedFromGraphDb", true))
			return true, nil
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}

	return false, nil
}

func (s *serviceLineItemService) Close(ctx context.Context, serviceLineItemId string, endedAt *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Close")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	sliExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), serviceLineItemId, neo4jutil.NodeLabelServiceLineItem)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error on checking if service line item exists: %s", err.Error())
		return err
	}
	if !sliExists {
		err := fmt.Errorf("service line item with id {%s} not found", serviceLineItemId)
		tracing.TraceErr(span, err)
		s.log.Errorf(err.Error())
		return err
	}

	closeRequest := servicelineitempb.CloseServiceLineItemGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             serviceLineItemId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		AppSource:      constants.AppSourceCustomerOsApi,
		EndedAt:        utils.ConvertTimeToTimestampPtr(endedAt),
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = s.grpcClients.ServiceLineItemClient.CloseServiceLineItem(ctx, &closeRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
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
}

func (s *serviceLineItemService) CreateOrUpdateInBulk(ctx context.Context, contractId string, sliBulkData []*ServiceLineItemDetails) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.CreateOrUpdateInBulk")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contractId", contractId))

	if len(sliBulkData) == 0 {
		return []string{}, nil
	}
	var responseIds []string

	sliDbNodes, err := s.repositories.ServiceLineItemRepository.GetForContracts(ctx, common.GetTenantFromContext(ctx), []string{contractId})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get service line items for contract: %s", err.Error())
		return []string{}, err
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
			err = s.Close(ctx, existingId, nil)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Failed to close service line item: %s", err.Error())
			}
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
			})
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
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error from events processing: %s", err.Error())
				return []string{}, err
			}
		}
	}

	return responseIds, nil
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
	}
}
