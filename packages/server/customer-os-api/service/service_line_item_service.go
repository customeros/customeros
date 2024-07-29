package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"math"
	"sort"
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
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
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
	Create(ctx context.Context, serviceLineItemDetails ServiceLineItemCreateData) (string, error)
	Update(ctx context.Context, serviceLineItemDetails ServiceLineItemUpdateData) error
	Delete(ctx context.Context, serviceLineItemId string) (bool, error)
	Close(ctx context.Context, serviceLineItemId string, endedAt *time.Time) error
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

func (s *serviceLineItemService) Create(ctx context.Context, serviceLineItemDetails ServiceLineItemCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItem.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "serviceLineItemDetails", serviceLineItemDetails)

	// check that quantity is not negative
	if serviceLineItemDetails.SliQuantity < 0 {
		err := errors.New("quantity must not be negative")
		tracing.TraceErr(span, err)
		return "", err
	}
	// check that price is not negative for non-one time
	if serviceLineItemDetails.SliPrice < 0 && serviceLineItemDetails.SliBilledType != neo4jenum.BilledTypeOnce {
		err := errors.New("price must not be negative")
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

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, model2.NodeLabelServiceLineItem, span)

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

	var baseServiceLineItemEntity *neo4jentity.ServiceLineItemEntity

	// check that given id is parentId, then use latest version as base
	serviceLineItems, err := s.services.CommonServices.ServiceLineItemService.GetServiceLineItemsByParentId(ctx, data.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting service line items by parent id {%s}: %s", data.Id, err.Error())
		return "", err
	}
	// sort by startedAt descending and select first one
	if len(*serviceLineItems) > 0 {
		sort.Slice(*serviceLineItems, func(i, j int) bool {
			return (*serviceLineItems)[i].StartedAt.After((*serviceLineItems)[j].StartedAt)
		})
		baseServiceLineItemEntity = &(*serviceLineItems)[0]
	} else {
		//if not found by parent id, treat data.id as SLI id
		baseServiceLineItemEntity, err = s.services.CommonServices.ServiceLineItemService.GetById(ctx, data.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error on getting contract line item by id {%s}: %s", data.Id, err.Error())
			return "", err
		}
	}

	if baseServiceLineItemEntity == nil {
		err := fmt.Errorf("contract line item with id {%s} not found", data.Id)
		tracing.TraceErr(span, err)
		return "", err
	}

	contractEntity, err := s.services.ContractService.GetContractByServiceLineItem(ctx, baseServiceLineItemEntity.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", baseServiceLineItemEntity.ID, err.Error())
		return "", err
	}

	startedAtDate := utils.ToDate(utils.IfNotNilTimeWithDefault(data.StartedAt, utils.Now()))

	// Check no SLI of the contract are cancelled
	for _, sli := range *serviceLineItems {
		if sli.Canceled {
			err = fmt.Errorf("contract line item with id {%s} is already ended", sli.ID)
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	// Do not allow creating new version if there is an existing version with the same start date
	for _, sli := range *serviceLineItems {
		if utils.ToDate(sli.StartedAt).Equal(startedAtDate) {
			err = fmt.Errorf("contract line item with id {%s} already exists with the same start date {%s}", sli.ID, startedAtDate.Format(time.DateOnly))
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	// For non-draft contracts do not allow creating new version before current active version
	if contractEntity.ContractStatus != neo4jenum.ContractStatusDraft {
		// identify live version
		var liveSli *neo4jentity.ServiceLineItemEntity
		for _, sli := range *serviceLineItems {
			if !utils.ToDate(sli.StartedAt).After(utils.Today()) && (sli.EndedAt == nil || sli.EndedAt.After(utils.Today())) {
				liveSli = &sli
				break
			}
		}
		if liveSli != nil && startedAtDate.Before(utils.ToDate(liveSli.StartedAt)) {
			err = fmt.Errorf("cannot create new version before current active version {%s}", liveSli.StartedAt.Format(time.DateOnly))
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	// Validate new version creation
	if baseServiceLineItemEntity.Billed == neo4jenum.BilledTypeOnce {
		err = fmt.Errorf("cannot create new version for one time contract line item with id {%s}", baseServiceLineItemEntity.ID)
		tracing.TraceErr(span, err)
		return "", err
	}

	// If contract was invoiced - do not allow creating new version before last invoiced date
	contractInvoiced, err := s.repositories.Neo4jRepositories.ContractReadRepository.IsContractInvoiced(ctx, common.GetTenantFromContext(ctx), contractEntity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on checking if contract was invoiced: %s", err.Error())
		return "", err
	}
	if contractInvoiced {
		// get last issued invoice
		lastInvoice, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetLastIssuedInvoiceForContract(ctx, common.GetTenantFromContext(ctx), contractEntity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error on getting last issued invoice for contract {%s}: %s", contractEntity.Id, err.Error)
		}
		if lastInvoice != nil {
			invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(lastInvoice)
			if startedAtDate.Before(utils.ToDate(invoiceEntity.PeriodEndDate)) {
				err = fmt.Errorf("cannot create new version for contract line item with id {%s} in the past", baseServiceLineItemEntity.ID)
				tracing.TraceErr(span, err)
				return "", err
			}
		}
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
		StartedAt:      utils.ConvertTimeToTimestampPtr(&startedAtDate),
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

	return response.Id, err
}

func (s *serviceLineItemService) Update(ctx context.Context, serviceLineItemDetails ServiceLineItemUpdateData) error {
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

	baseServiceLineItemEntity, err := s.services.CommonServices.ServiceLineItemService.GetById(ctx, serviceLineItemDetails.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract line item by id {%s}: %s", serviceLineItemDetails.Id, err.Error())
		return err
	}

	contractEntity, err := s.services.ContractService.GetContractByServiceLineItem(ctx, serviceLineItemDetails.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", serviceLineItemDetails.Id, err.Error())
		return err
	}

	isRetroactiveCorrection := serviceLineItemDetails.IsRetroactiveCorrection
	contractIsInvoiced, _ := s.repositories.Neo4jRepositories.ContractReadRepository.IsContractInvoiced(ctx, common.GetTenantFromContext(ctx), contractEntity.Id)
	sliIsInvoiced, _ := s.repositories.Neo4jRepositories.ServiceLineItemReadRepository.WasServiceLineItemInvoiced(ctx, common.GetTenantFromContext(ctx), baseServiceLineItemEntity.ID)
	startedAt := utils.ToDate(utils.IfNotNilTimeWithDefault(serviceLineItemDetails.StartedAt, baseServiceLineItemEntity.StartedAt))

	anyFieldChanged := baseServiceLineItemEntity.Name != serviceLineItemDetails.SliName ||
		baseServiceLineItemEntity.Price != serviceLineItemDetails.SliPrice ||
		baseServiceLineItemEntity.Quantity != serviceLineItemDetails.SliQuantity ||
		baseServiceLineItemEntity.VatRate != serviceLineItemDetails.SliVatRate ||
		baseServiceLineItemEntity.Comments != serviceLineItemDetails.SliComments ||
		baseServiceLineItemEntity.Billed != serviceLineItemDetails.SliBilledType

	// If no changes recorded, return
	if !anyFieldChanged && (utils.ToDate(baseServiceLineItemEntity.StartedAt).Equal(startedAt) || sliIsInvoiced) {
		span.LogFields(log.String("result", "No changes recorded"))
		return nil
	}

	if baseServiceLineItemEntity.Canceled {
		err = fmt.Errorf("contract line item with id {%s} is already ended", serviceLineItemDetails.Id)
		tracing.TraceErr(span, err)
		return err
	}

	// price impacted fields changed
	priceImpactedFieldsChanged := baseServiceLineItemEntity.Price != serviceLineItemDetails.SliPrice ||
		baseServiceLineItemEntity.Quantity != serviceLineItemDetails.SliQuantity ||
		baseServiceLineItemEntity.VatRate != serviceLineItemDetails.SliVatRate

	// Check no SLI of the contract are cancelled
	serviceLineItemsOfSameParent, err := s.services.CommonServices.ServiceLineItemService.GetServiceLineItemsByParentId(ctx, baseServiceLineItemEntity.ParentID)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting service line items for contract {%s}: %s", contractEntity.Id, err.Error())
		return err
	}

	// sort serviceLineItemsOfSameParent by startedAt ascending
	sort.Slice(*serviceLineItemsOfSameParent, func(i, j int) bool {
		return (*serviceLineItemsOfSameParent)[i].StartedAt.Before((*serviceLineItemsOfSameParent)[j].StartedAt)
	})
	// get latest version
	lastVersion := (*serviceLineItemsOfSameParent)[len(*serviceLineItemsOfSameParent)-1]

	// Do not update SLI if last version is updated without any changes and start date is close to current timestamp
	if lastVersion.ID == baseServiceLineItemEntity.ID {
		if !anyFieldChanged && !priceImpactedFieldsChanged {
			// diff between passed date and now in seconds
			diff := utils.IfNotNilTimeWithDefault(serviceLineItemDetails.StartedAt, utils.Now()).Unix() - utils.Now().Unix()
			diffAbs := math.Abs(float64(diff))
			// if diff is less than 10 min, skip update
			if diffAbs < float64(600) {
				span.LogFields(log.String("result", "No changes recorded, start date is close to current timestamp"))
				return nil
			}
		}
	}

	// Check no SLI of the contract are cancelled
	for _, sli := range *serviceLineItemsOfSameParent {
		if sli.Canceled {
			err = fmt.Errorf("contract line item with id {%s} is already ended", sli.ID)
			tracing.TraceErr(span, err)
			return err
		}
	}

	// check that quantity is not negative
	if serviceLineItemDetails.SliQuantity < 0 {
		err := errors.New("quantity must not be negative")
		tracing.TraceErr(span, err)
		return err
	}
	// check that price is not negative for non-one time
	if serviceLineItemDetails.SliPrice < 0 && baseServiceLineItemEntity.Billed != neo4jenum.BilledTypeOnce {
		err := errors.New("price must not be negative")
		tracing.TraceErr(span, err)
		return err
	}

	// check that billing cycle is not changed
	if serviceLineItemDetails.SliBilledType.String() == "" {
		serviceLineItemDetails.SliBilledType = baseServiceLineItemEntity.Billed
	}
	if baseServiceLineItemEntity.Billed.String() != serviceLineItemDetails.SliBilledType.String() && baseServiceLineItemEntity.Billed.String() != "" {
		err = fmt.Errorf("cannot change billing cycle for contract line item with id {%s}", serviceLineItemDetails.Id)
		tracing.TraceErr(span, err)
		return err
	}

	if baseServiceLineItemEntity.IsOneTime() ||
		utils.ToDate(baseServiceLineItemEntity.StartedAt) == startedAt {
		isRetroactiveCorrection = true
	}
	if !priceImpactedFieldsChanged {
		isRetroactiveCorrection = true
	}

	// Do not allow changing price impacting data for invoiced SLIs
	if isRetroactiveCorrection && priceImpactedFieldsChanged && sliIsInvoiced {
		err = fmt.Errorf("service line item with id {%s} is included in invoice and cannot be updated", serviceLineItemDetails.Id)
		tracing.TraceErr(span, err)
		return err
	}

	// Do not allow updating past SLIs for invoiced contracts
	if isRetroactiveCorrection && contractIsInvoiced {
		tenantSettings, _ := s.services.CommonServices.TenantService.GetTenantSettings(ctx)
		isInvoicingPostpaid := tenantSettings.InvoicingPostpaid
		referenceDate := utils.Today()
		if isInvoicingPostpaid && contractEntity.NextInvoiceDate != nil {
			referenceDate = utils.ToDate(*contractEntity.NextInvoiceDate)
		}
		if startedAt.Before(referenceDate) {
			err = fmt.Errorf("cannot update contract line item with id {%s} and start date before {%s}", serviceLineItemDetails.Id, referenceDate.Format(time.DateOnly))
			tracing.TraceErr(span, err)
			return err
		}
	}

	span.LogFields(log.Bool("result.isRetroactiveCorrection", isRetroactiveCorrection))

	if isRetroactiveCorrection == true {
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

		// if start date is changed, validate that change is allowed
		if utils.ToDate(baseServiceLineItemEntity.StartedAt) != utils.ToDate(startedAt) {
			// Do not allow creating new version if there is an existing version with the same start date
			for _, sli := range *serviceLineItemsOfSameParent {
				if sli.ID != baseServiceLineItemEntity.ID && utils.ToDate(sli.StartedAt).Equal(startedAt) {
					err = fmt.Errorf("Other version with the same start date {%s} already exists", startedAt.Format(time.DateOnly))
					tracing.TraceErr(span, err)
					return err
				}
			}
			if contractEntity.ContractStatus != neo4jenum.ContractStatusDraft && !startedAt.After(utils.Today()) {
				err = fmt.Errorf("cannot update contract line item with id {%s} in the past", serviceLineItemDetails.Id)
				tracing.TraceErr(span, err)
				return err
			}
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
		// Create new SLI version
		_, err := s.NewVersion(ctx, ServiceLineItemNewVersionData{
			Id:        baseServiceLineItemEntity.ParentID,
			Name:      utils.StringFirstNonEmpty(serviceLineItemDetails.SliName, baseServiceLineItemEntity.Name),
			Price:     serviceLineItemDetails.SliPrice,
			Quantity:  serviceLineItemDetails.SliQuantity,
			Comments:  utils.IfNotNilString(serviceLineItemDetails.SliComments),
			Source:    serviceLineItemDetails.Source,
			AppSource: utils.StringFirstNonEmpty(serviceLineItemDetails.AppSource, constants.AppSourceCustomerOsApi),
			VatRate:   serviceLineItemDetails.SliVatRate,
			StartedAt: serviceLineItemDetails.StartedAt,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *serviceLineItemService) Delete(ctx context.Context, serviceLineItemId string) (completed bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Delete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	sliEntity, err := s.services.CommonServices.ServiceLineItemService.GetById(ctx, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting service line item by id {%s}: %s", serviceLineItemId, err.Error())
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

	// if contract is not draft prevent removing current or past SLIs
	contractEntity, err := s.services.ContractService.GetContractByServiceLineItem(ctx, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", serviceLineItemId, err.Error())
		return false, err
	}
	if contractEntity.ContractStatus != neo4jenum.ContractStatusDraft && !sliEntity.StartedAt.After(utils.Today()) {
		err = fmt.Errorf("cannot delete contract line item with id {%s} in the past", serviceLineItemId)
		tracing.TraceErr(span, err)
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

	return false, nil
}

func (s *serviceLineItemService) Close(ctx context.Context, serviceLineItemId string, endedAt *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Close")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))
	span.SetTag(tracing.SpanTagEntityId, serviceLineItemId)

	contractEntity, err := s.services.ContractService.GetContractByServiceLineItem(ctx, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", serviceLineItemId, err.Error())
		return err
	}

	// if contract is draft - delete SLI
	if contractEntity.IsDraft() {
		_, err = s.Delete(ctx, serviceLineItemId)
		return err
	}

	sliEntity, err := s.services.CommonServices.ServiceLineItemService.GetById(ctx, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting service line item by id {%s}: %s", serviceLineItemId, err.Error())
		return err
	}

	// Future SLI to be deleted
	if sliEntity.StartedAt.After(utils.Today()) {
		_, err = s.Delete(ctx, serviceLineItemId)
		return err
	}

	// closing past SLIs not allowed
	if sliEntity.EndedAt != nil && sliEntity.EndedAt.Before(utils.Today()) {
		err = fmt.Errorf("cannot close contract line item with id {%s} in the past", serviceLineItemId)
		tracing.TraceErr(span, err)
		return err
	}

	// First remove any future SLI
	sliEntities, err := s.services.CommonServices.ServiceLineItemService.GetServiceLineItemsForContract(ctx, contractEntity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on getting service line items for contract {%s}: %s", contractEntity.Id, err.Error())
		return err
	}
	for _, sli := range *sliEntities {
		if sli.StartedAt.After(utils.Today()) {
			_, err = s.Delete(ctx, sli.ID)
			if err != nil {
				return err
			}
		}
	}

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

	return nil
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
	CloseVersion            bool
	NewVersion              bool
}

func (s *serviceLineItemService) CreateOrUpdateOrCloseInBulk(ctx context.Context, contractId string, sliBulkData []*ServiceLineItemDetails) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.CreateOrUpdateOrCloseInBulk")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contractId", contractId))

	var responseIds []string

	for _, serviceLineItem := range sliBulkData {
		// check all quantities are not negative
		if serviceLineItem.Quantity < 0 {
			err := fmt.Errorf("quantity must not be negative")
			tracing.TraceErr(span, err)
			return []string{}, err
		}
	}

	for _, serviceLineItem := range sliBulkData {
		if serviceLineItem.Id == "" && !serviceLineItem.CloseVersion && !serviceLineItem.NewVersion {
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
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error from events processing: %s", err.Error())
				return []string{}, err
			}
			responseIds = append(responseIds, itemId)
		} else if serviceLineItem.CloseVersion && serviceLineItem.Id != "" {
			err := s.Close(ctx, serviceLineItem.Id, nil)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Failed to close service line item: %s", err.Error())
			}
		} else if serviceLineItem.NewVersion && !serviceLineItem.CloseVersion {
			itemId, err := s.NewVersion(ctx, ServiceLineItemNewVersionData{
				Id:        serviceLineItem.Id,
				Name:      serviceLineItem.Name,
				Price:     serviceLineItem.Price,
				Quantity:  serviceLineItem.Quantity,
				Comments:  serviceLineItem.Comments,
				Source:    neo4jentity.DataSourceOpenline,
				AppSource: constants.AppSourceCustomerOsApi,
				VatRate:   serviceLineItem.VatRate,
				StartedAt: serviceLineItem.StartedAt,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error from events processing: %s", err.Error())
				return []string{}, err
			}
			responseIds = append(responseIds, itemId)
		} else if serviceLineItem.Id != "" && !serviceLineItem.CloseVersion && !serviceLineItem.NewVersion {
			err := s.Update(ctx, ServiceLineItemUpdateData{
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
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error from events processing: %s", err.Error())
				return []string{}, err
			}
			responseIds = append(responseIds, serviceLineItem.Id)
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
		StartedAt:               input.ServiceStarted,
		CloseVersion:            utils.IfNotNilBool(input.CloseVersion),
		NewVersion:              utils.IfNotNilBool(input.NewVersion),
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
