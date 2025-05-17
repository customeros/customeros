package api_sli

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"math"
	"sort"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type serviceLineItemService struct {
	log            logger.Logger
	repositories   *repository.Repositories
	sli            interfaces.ServiceLineItemService
	contract       cosapi_interfaces.ContractService
	tenantSettings interfaces.TenantSettingsService
}

func NewServiceLineItemService(
	log logger.Logger,
	repositories *repository.Repositories,
	sli interfaces.ServiceLineItemService,
	contract cosapi_interfaces.ContractService,
	tenantSettings interfaces.TenantSettingsService,
) cosapi_interfaces.ServiceLineItemService {
	return &serviceLineItemService{
		log:            log,
		repositories:   repositories,
		sli:            sli,
		contract:       contract,
		tenantSettings: tenantSettings,
	}
}

func (s *serviceLineItemService) Create(ctx context.Context, serviceLineItemDetails cosapi_interfaces.ServiceLineItemCreateData) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ServiceLineItemService.Create")
	defer spans.Finish()
	spans.LogObjectAsJson("serviceLineItemDetails", serviceLineItemDetails)

	// check that quantity is not negative
	if serviceLineItemDetails.SliQuantity < 0 {
		err := errors.New("quantity must not be negative")
		spans.TraceError(err)
		return "", err
	}
	// check that price is not negative for non-one time
	if serviceLineItemDetails.SliPrice < 0 && serviceLineItemDetails.SliBilledType != neo4jenum.BilledTypeOnce {
		err := errors.New("price must not be negative")
		spans.TraceError(err)
		return "", err
	}

	sliDataFields := data_fields.SLIFields{
		ContractId:  utils.StringPtr(serviceLineItemDetails.ContractId),
		SkuId:       utils.StringPtr(serviceLineItemDetails.SkuId),
		Description: serviceLineItemDetails.SliDescription,
		Quantity:    utils.Int64Ptr(serviceLineItemDetails.SliQuantity),
		Price:       utils.Float64Ptr(serviceLineItemDetails.SliPrice),
		TaxRate:     utils.Float64Ptr(serviceLineItemDetails.SliVatRate),
		StartedAt:   serviceLineItemDetails.StartedAt,
		EndedAt:     serviceLineItemDetails.EndedAt,
		Source:      utils.StringPtr(serviceLineItemDetails.Source.String()),
		NewVersion:  utils.BoolPtr(false),
	}

	if serviceLineItemDetails.SkuId == "" {
		err := fmt.Errorf("sku id is required for all service line items")
		spans.TraceError(err)
		return "", err
	}

	sliDataFields.BilledType = utils.ToPtr(serviceLineItemDetails.SliBilledType)

	sliId, err := s.sli.Save(ctx, nil, nil, sliDataFields)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	return sliId, nil
}

func (s *serviceLineItemService) NewVersion(ctx context.Context, data cosapi_interfaces.ServiceLineItemNewVersionData) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ServiceLineItem.NewVersion")
	defer spans.Finish()
	spans.LogObjectAsJson("serviceLineItemDetails", data)

	if data.Id == "" {
		err := fmt.Errorf("(ServiceLineItemService.NewVersion) contract line item id is missing")
		s.log.Error(err.Error())
		spans.TraceError(err)
		return "", err
	}

	var baseServiceLineItemEntity *neo4jentity.ServiceLineItemEntity

	// check that given id is parentId, then use latest version as base
	serviceLineItems, err := s.sli.GetServiceLineItemsByParentId(ctx, data.Id)
	if err != nil {
		spans.TraceError(err)
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
		// if not found by parent id, treat data.id as SLI id
		baseServiceLineItemEntity, err = s.sli.GetById(ctx, data.Id)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error on getting contract line item by id {%s}: %s", data.Id, err.Error())
			return "", err
		}
	}

	if baseServiceLineItemEntity == nil {
		err := fmt.Errorf("contract line item with id {%s} not found", data.Id)
		spans.TraceError(err)
		return "", err
	}

	contractEntity, err := s.contract.GetContractByServiceLineItem(ctx, baseServiceLineItemEntity.ID)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", baseServiceLineItemEntity.ID, err.Error())
		return "", err
	}

	startedAtDate := utils.ToDate(utils.IfNotNilTimeWithDefault(data.StartedAt, utils.Now()))

	// Check no SLI of the contract are cancelled
	for _, sli := range *serviceLineItems {
		if sli.Canceled {
			err = fmt.Errorf("contract line item with id {%s} is already ended", sli.ID)
			spans.TraceError(err)
			return "", err
		}
	}

	// Do not allow creating new version if there is an existing version with the same start date
	for _, sli := range *serviceLineItems {
		if utils.ToDate(sli.StartedAt).Equal(startedAtDate) {
			err = fmt.Errorf("contract line item with id {%s} already exists with the same start date {%s}", sli.ID, startedAtDate.Format(time.DateOnly))
			spans.TraceError(err)
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
			spans.TraceError(err)
			return "", err
		}
	}

	// Validate new version creation
	if baseServiceLineItemEntity.Billed == neo4jenum.BilledTypeOnce {
		err = fmt.Errorf("cannot create new version for one time contract line item with id {%s}", baseServiceLineItemEntity.ID)
		spans.TraceError(err)
		return "", err
	}

	// If contract was invoiced - do not allow creating new version before last invoiced date
	contractInvoiced, err := s.repositories.Neo4jRepositories.ContractReadRepository.IsContractInvoiced(ctx, common.GetTenantFromContext(ctx), contractEntity.Id)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error on checking if contract was invoiced: %s", err.Error())
		return "", err
	}
	if contractInvoiced {
		// get last issued invoice
		lastInvoice, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetLastIssuedInvoiceForContract(ctx, common.GetTenantFromContext(ctx), contractEntity.Id)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error on getting last issued invoice for contract {%s}: %s", contractEntity.Id, err.Error)
		}
		if lastInvoice != nil {
			invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(lastInvoice)
			if startedAtDate.Before(utils.ToDate(invoiceEntity.PeriodEndDate)) {
				err = fmt.Errorf("cannot create new version for contract line item with id {%s} in the past", baseServiceLineItemEntity.ID)
				spans.TraceError(err)
				return "", err
			}
		}
	}

	sliDataFields := data_fields.SLIFields{
		ContractId:  utils.StringPtr(contractEntity.Id),
		ParentId:    utils.StringPtr(baseServiceLineItemEntity.ParentID),
		SkuId:       utils.StringPtr(utils.StringFirstNonEmpty(baseServiceLineItemEntity.SkuId, data.SkuId)),
		Description: utils.StringPtr(utils.StringFirstNonEmpty(utils.IfNotNilString(data.Description), baseServiceLineItemEntity.Description)),
		Quantity:    utils.Int64Ptr(data.Quantity),
		Price:       utils.Float64Ptr(data.Price),
		TaxRate:     utils.Float64Ptr(data.VatRate),
		StartedAt:   &startedAtDate,
		Comments:    utils.StringPtr(data.Comments),
		Source:      utils.StringPtr(data.Source.String()),
		BilledType:  utils.ToPtr(baseServiceLineItemEntity.Billed),
		NewVersion:  utils.BoolPtr(true),
	}

	sliId, err := s.sli.Save(ctx, nil, nil, sliDataFields)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	return sliId, err
}

func (s *serviceLineItemService) Update(ctx context.Context, serviceLineItemDetails cosapi_interfaces.ServiceLineItemUpdateData) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ServiceLineItemService.Update")
	defer spans.Finish()
	spans.LogObjectAsJson("serviceLineItemDetails", serviceLineItemDetails)

	if serviceLineItemDetails.Id == "" {
		err := fmt.Errorf("(ServiceLineItemService.Update) contract line item id is missing")
		s.log.Error(err.Error())
		spans.TraceError(err)
		return err
	}

	baseServiceLineItemEntity, err := s.sli.GetById(ctx, serviceLineItemDetails.Id)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error on getting contract line item by id {%s}: %s", serviceLineItemDetails.Id, err.Error())
		return err
	}
	spans.TagEntity(baseServiceLineItemEntity.ID)

	contractEntity, err := s.contract.GetContractByServiceLineItem(ctx, serviceLineItemDetails.Id)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", serviceLineItemDetails.Id, err.Error())
		return err
	}

	isRetroactiveCorrection := serviceLineItemDetails.IsRetroactiveCorrection
	contractIsInvoiced, _ := s.repositories.Neo4jRepositories.ContractReadRepository.IsContractInvoiced(ctx, common.GetTenantFromContext(ctx), contractEntity.Id)
	sliIsInvoiced, _ := s.repositories.Neo4jRepositories.ServiceLineItemReadRepository.WasServiceLineItemInvoiced(ctx, common.GetTenantFromContext(ctx), baseServiceLineItemEntity.ID)
	startedAt := utils.ToDate(utils.IfNotNilTimeWithDefault(serviceLineItemDetails.StartedAt, baseServiceLineItemEntity.StartedAt))

	anyFieldChanged := (baseServiceLineItemEntity.SkuId != serviceLineItemDetails.SkuId && serviceLineItemDetails.SkuId != "") ||
		baseServiceLineItemEntity.Price != serviceLineItemDetails.SliPrice ||
		baseServiceLineItemEntity.Quantity != serviceLineItemDetails.SliQuantity ||
		baseServiceLineItemEntity.VatRate != serviceLineItemDetails.SliVatRate ||
		baseServiceLineItemEntity.Comments != serviceLineItemDetails.SliComments ||
		baseServiceLineItemEntity.Description != utils.IfNotNilString(serviceLineItemDetails.SliDescription) ||
		(baseServiceLineItemEntity.Billed != serviceLineItemDetails.SliBilledType && serviceLineItemDetails.SliBilledType != neo4jenum.BilledTypeNone)
	spans.LogKV("anyFieldChanged", anyFieldChanged)

	// If no changes recorded, return
	if !anyFieldChanged && (utils.ToDate(baseServiceLineItemEntity.StartedAt).Equal(startedAt) || utils.CloseToNow(startedAt) || sliIsInvoiced) {
		spans.LogKV("result", "No changes recorded")
		return nil
	}

	if baseServiceLineItemEntity.Canceled {
		err = fmt.Errorf("contract line item with id {%s} is already ended", serviceLineItemDetails.Id)
		spans.TraceError(err)
		return err
	}

	// price impacted fields changed
	priceImpactedFieldsChanged := baseServiceLineItemEntity.Price != serviceLineItemDetails.SliPrice ||
		baseServiceLineItemEntity.Quantity != serviceLineItemDetails.SliQuantity ||
		baseServiceLineItemEntity.VatRate != serviceLineItemDetails.SliVatRate

	// Check no SLI of the contract are cancelled
	serviceLineItemsOfSameParent, err := s.sli.GetServiceLineItemsByParentId(ctx, baseServiceLineItemEntity.ParentID)
	if err != nil {
		spans.TraceError(err)
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
				spans.LogKV("result", "No changes recorded, start date is close to current timestamp")
				return nil
			}
		}
	}

	// Check no SLI of the contract are cancelled
	for _, sli := range *serviceLineItemsOfSameParent {
		if sli.Canceled {
			err = fmt.Errorf("contract line item with id {%s} is already ended", sli.ID)
			spans.TraceError(err)
			return err
		}
	}

	// check that quantity is not negative
	if serviceLineItemDetails.SliQuantity < 0 {
		err := errors.New("quantity must not be negative")
		spans.TraceError(err)
		return err
	}
	// check that price is not negative for non-one time
	if serviceLineItemDetails.SliPrice < 0 && baseServiceLineItemEntity.Billed != neo4jenum.BilledTypeOnce {
		err := errors.New("price must not be negative")
		spans.TraceError(err)
		return err
	}

	// check that billing cycle is not changed
	if serviceLineItemDetails.SliBilledType.String() == "" {
		serviceLineItemDetails.SliBilledType = baseServiceLineItemEntity.Billed
	}
	if baseServiceLineItemEntity.Billed.String() != serviceLineItemDetails.SliBilledType.String() && baseServiceLineItemEntity.Billed.String() != "" {
		err = fmt.Errorf("cannot change billing cycle for contract line item with id {%s}", serviceLineItemDetails.Id)
		spans.TraceError(err)
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
		spans.TraceError(err)
		return err
	}

	// Do not allow updating past SLIs for invoiced contracts
	if isRetroactiveCorrection && contractIsInvoiced {
		tenantSettings, _ := s.tenantSettings.GetTenantSettings(ctx)
		isInvoicingPostpaid := tenantSettings.InvoicingPostpaid
		referenceDate := utils.Today()
		if isInvoicingPostpaid && contractEntity.NextInvoiceDate != nil {
			referenceDate = utils.ToDate(*contractEntity.NextInvoiceDate)
		}
		if startedAt.Before(referenceDate) {
			err = fmt.Errorf("cannot update contract line item with id {%s} and start date before {%s}", serviceLineItemDetails.Id, referenceDate.Format(time.DateOnly))
			spans.TraceError(err)
			return err
		}
	}

	spans.LogKV("result.isRetroactiveCorrection", isRetroactiveCorrection)

	if isRetroactiveCorrection == true {
		sliDataFields := data_fields.SLIFields{
			SkuId:       utils.StringPtrNillable(serviceLineItemDetails.SkuId),
			Description: serviceLineItemDetails.SliDescription,
			Quantity:    utils.Int64Ptr(serviceLineItemDetails.SliQuantity),
			Price:       utils.Float64Ptr(serviceLineItemDetails.SliPrice),
			TaxRate:     utils.Float64Ptr(serviceLineItemDetails.SliVatRate),
			Comments:    utils.StringPtr(serviceLineItemDetails.SliComments),
			BilledType:  utils.ToPtr(serviceLineItemDetails.SliBilledType),
			ParentId:    utils.StringPtr(baseServiceLineItemEntity.ParentID),
			NewVersion:  utils.BoolPtr(false),
		}

		// if start date is changed, validate that change is allowed
		if utils.ToDate(baseServiceLineItemEntity.StartedAt) != utils.ToDate(startedAt) {
			// Do not allow creating new version if there is an existing version with the same start date
			for _, sli := range *serviceLineItemsOfSameParent {
				if sli.ID != baseServiceLineItemEntity.ID && utils.ToDate(sli.StartedAt).Equal(startedAt) {
					err = fmt.Errorf("Other version with the same start date {%s} already exists", startedAt.Format(time.DateOnly))
					spans.TraceError(err)
					return err
				}
			}
			if contractEntity.ContractStatus != neo4jenum.ContractStatusDraft && contractEntity.NextInvoiceDate != nil {
				lastInvoice, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetLastIssuedInvoiceForContract(ctx, common.GetTenantFromContext(ctx), contractEntity.Id)
				if err != nil {
					spans.TraceError(err)
					s.log.Errorf("Error on getting last issued invoice for contract {%s}: %s", contractEntity.Id, err.Error)
				}
				if lastInvoice != nil {
					invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(lastInvoice)
					if !startedAt.After(invoiceEntity.PeriodEndDate) {
						err = fmt.Errorf("cannot update contract line item with id {%s} in the past", serviceLineItemDetails.Id)
						spans.TraceError(err)
						return err
					}
				}
			}
			sliDataFields.StartedAt = serviceLineItemDetails.StartedAt
		}

		_, err = s.sli.Save(ctx, nil, &serviceLineItemDetails.Id, sliDataFields)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error on updating service line item with id {%s}: %s", serviceLineItemDetails.Id, err.Error())
			return err
		}
	} else {
		// Create new SLI version
		_, err := s.NewVersion(ctx, cosapi_interfaces.ServiceLineItemNewVersionData{
			Id:          baseServiceLineItemEntity.ParentID,
			SkuId:       utils.StringFirstNonEmpty(serviceLineItemDetails.SkuId, baseServiceLineItemEntity.SkuId),
			Description: serviceLineItemDetails.SliDescription,
			Price:       serviceLineItemDetails.SliPrice,
			Quantity:    serviceLineItemDetails.SliQuantity,
			Comments:    utils.IfNotNilString(serviceLineItemDetails.SliComments),
			Source:      serviceLineItemDetails.Source,
			AppSource:   utils.StringFirstNonEmpty(serviceLineItemDetails.AppSource, constants.AppSourceCustomerOsApi),
			VatRate:     serviceLineItemDetails.SliVatRate,
			StartedAt:   serviceLineItemDetails.StartedAt,
		})
		if err != nil {
			spans.TraceError(err)
			return err
		}
	}

	return nil
}

func (s *serviceLineItemService) Delete(ctx context.Context, serviceLineItemId string) (completed bool, err error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "ServiceLineItemService.Delete")
	defer span.Finish()
	span.LogKV("serviceLineItemId", serviceLineItemId)

	sliEntity, err := s.sli.GetById(ctx, serviceLineItemId)
	if err != nil {
		span.TraceError(err)
		s.log.Errorf("Error on getting service line item by id {%s}: %s", serviceLineItemId, err.Error())
		return false, err
	}

	// Check SLI is not invoiced
	sliInvoiced, err := s.repositories.Neo4jRepositories.ServiceLineItemReadRepository.WasServiceLineItemInvoiced(ctx, common.GetTenantFromContext(ctx), serviceLineItemId)
	if err != nil {
		span.TraceError(err)
		s.log.Errorf("Error on checking if service line item was invoiced: %s", err.Error())
		return false, err
	}
	if sliInvoiced {
		err := fmt.Errorf("service line item with id {%s} is included in invoice and cannot be deleted", serviceLineItemId)
		span.TraceError(err)
		s.log.Errorf(err.Error())
		return false, err
	}

	// if contract is not draft prevent removing current or past SLIs
	contractEntity, err := s.contract.GetContractByServiceLineItem(ctx, serviceLineItemId)
	if err != nil {
		span.TraceError(err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", serviceLineItemId, err.Error())
		return false, err
	}
	if contractEntity.ContractStatus != neo4jenum.ContractStatusDraft && !sliEntity.StartedAt.After(utils.Today()) {
		err = fmt.Errorf("cannot delete contract line item with id {%s} in the past", serviceLineItemId)
		span.TraceError(err)
		return false, err
	}

	err = s.sli.Delete(ctx, nil, serviceLineItemId)
	if err != nil {
		span.TraceError(err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return false, err
	}

	return false, nil
}

func (s *serviceLineItemService) Close(ctx context.Context, serviceLineItemId string, endedAt *time.Time) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "ServiceLineItemService.Close")
	defer span.Finish()
	span.LogKV("serviceLineItemId", serviceLineItemId)

	contractEntity, err := s.contract.GetContractByServiceLineItem(ctx, serviceLineItemId)
	if err != nil {
		span.TraceError(err)
		s.log.Errorf("Error on getting contract by service line item id {%s}: %s", serviceLineItemId, err.Error())
		return err
	}

	// if contract is draft - delete SLI
	if contractEntity.IsDraft() {
		_, err = s.Delete(ctx, serviceLineItemId)
		return err
	}

	currentSliEntity, err := s.sli.GetById(ctx, serviceLineItemId)
	if err != nil {
		span.TraceError(err)
		s.log.Errorf("Error on getting service line item by id {%s}: %s", serviceLineItemId, err.Error())
		return err
	}

	// Future SLI to be deleted
	if currentSliEntity.StartedAt.After(utils.Today()) {
		_, err = s.Delete(ctx, serviceLineItemId)
		return err
	}

	// closing past SLIs not allowed
	if currentSliEntity.EndedAt != nil && currentSliEntity.EndedAt.Before(utils.Today()) {
		err = fmt.Errorf("contract line item with id {%s} is already closed", serviceLineItemId)
		span.TraceError(err)
		return err
	}

	// First remove any future SLI with same parent ID
	sliEntities, err := s.sli.GetServiceLineItemsByParentId(ctx, currentSliEntity.ParentID)
	if err != nil {
		span.TraceError(err)
		s.log.Errorf("Error on getting service line items by parent id {%s}: %s", currentSliEntity.ParentID, err.Error())
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

	err = s.sli.Close(ctx, nil, serviceLineItemId, utils.IfNotNilTimeWithDefault(endedAt, utils.Now()))
	if err != nil {
		span.TraceError(err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}
