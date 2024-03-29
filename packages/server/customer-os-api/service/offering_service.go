package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	offeringpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/offering"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OfferingService interface {
	CreateOffering(ctx context.Context, input *model.OfferingCreateInput) (string, error)
	UpdateOffering(ctx context.Context, input *model.OfferingUpdateInput) error
	GetOfferings(ctx context.Context) (*neo4jentity.OfferingEntities, error)
	GetOffering(ctx context.Context, id string) (*neo4jentity.OfferingEntity, error)
}

type offeringService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewOfferingService(log logger.Logger, repository *repository.Repositories, grpcClients *grpc_client.Clients) OfferingService {
	return &offeringService{
		log:          log,
		repositories: repository,
		grpcClients:  grpcClients,
	}
}

func (s *offeringService) GetOfferings(ctx context.Context) (*neo4jentity.OfferingEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OfferingService.GetTenantOfferings")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)

	dbNodes, err := s.repositories.Neo4jRepositories.OfferingReadRepository.GetOfferings(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantOfferings: %s", err.Error())
	}

	tenantOfferings := neo4jentity.OfferingEntities{}
	for _, dbNode := range dbNodes {
		tenantOfferings = append(tenantOfferings, *neo4jmapper.MapDbNodeToOfferingEntity(dbNode))
	}

	return &tenantOfferings, nil
}

func (s *offeringService) GetOffering(ctx context.Context, id string) (*neo4jentity.OfferingEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OfferingService.GetOffering")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.String("offeringId", id))

	dbNode, err := s.repositories.Neo4jRepositories.OfferingReadRepository.GetOfferingById(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetOffering: %s", err.Error())
	}

	return neo4jmapper.MapDbNodeToOfferingEntity(dbNode), nil
}

func (s *offeringService) CreateOffering(ctx context.Context, input *model.OfferingCreateInput) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OfferingService.CreateOffering")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	tracing.LogObjectAsJson(span, "input", input)

	grpcRequest := offeringpb.CreateOfferingGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
		Name:                                   utils.IfNotNilString(input.Name),
		Active:                                 utils.IfNotNilBool(input.Active),
		PricingPeriodInMonths:                  utils.IfNotNilInt64(input.PricingPeriodInMonths, func() int64 { return 1 }),
		Price:                                  utils.IfNotNilFloat64(input.Price),
		PriceCalculated:                        utils.IfNotNilBool(input.PriceCalculated),
		Conditional:                            utils.IfNotNilBool(input.Conditional),
		Taxable:                                utils.IfNotNilBool(input.Taxable),
		PriceCalculationRevenueSharePercentage: utils.IfNotNilFloat64(input.PriceCalculationRevenueSharePercentage),
		ConditionalsMinimumChargeAmount:        utils.IfNotNilFloat64(input.ConditionalsMinimumChargeAmount),
	}
	if input.Type != nil {
		grpcRequest.Type = input.Type.String()
	}
	if input.PricingModel != nil {
		grpcRequest.PricingModel = input.PricingModel.String()
	}
	if input.Currency != nil {
		grpcRequest.Currency = input.Currency.String()
	}
	if input.PriceCalculationType != nil {
		grpcRequest.PriceCalculationType = input.PriceCalculationType.String()
	}
	if input.ConditionalsMinimumChargePeriod != nil {
		grpcRequest.ConditionalsMinimumChargePeriod = input.ConditionalsMinimumChargePeriod.String()
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := CallEventsPlatformGRPCWithRetry[*commonpb.IdResponse](func() (*commonpb.IdResponse, error) {
		return s.grpcClients.OfferingClient.CreateOffering(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	return response.Id, nil
}

func (s *offeringService) UpdateOffering(ctx context.Context, input *model.OfferingUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OfferingService.UpdateOffering")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	if input.ID == "" {
		err := fmt.Errorf("offering id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	var fieldsMask []offeringpb.OfferingFieldMask
	updateRequest := offeringpb.UpdateOfferingGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             input.ID,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		AppSource:      constants.AppSourceCustomerOsApi,
	}

	if input.Name != nil {
		updateRequest.Name = *input.Name
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_NAME)
	}
	if input.Active != nil {
		updateRequest.Active = *input.Active
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_ACTIVE)
	}
	if input.Type != nil {
		updateRequest.Type = input.Type.String()
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_TYPE)
	}
	if input.PricingModel != nil {
		updateRequest.PricingModel = input.PricingModel.String()
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICING_MODEL)
	}
	if input.PricingPeriodInMonths != nil {
		updateRequest.PricingPeriodInMonths = *input.PricingPeriodInMonths
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICING_PERIOD_IN_MONTHS)
	}
	if input.Price != nil {
		updateRequest.Price = *input.Price
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE)
	}
	if input.Currency != nil {
		updateRequest.Currency = input.Currency.String()
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_CURRENCY)
	}
	if input.PriceCalculated != nil {
		updateRequest.PriceCalculated = *input.PriceCalculated
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE_CALCULATED)
	}
	if input.Conditional != nil {
		updateRequest.Conditional = *input.Conditional
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_CONDITIONAL)
	}
	if input.Taxable != nil {
		updateRequest.Taxable = *input.Taxable
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_TAXABLE)
	}
	if input.PriceCalculationRevenueSharePercentage != nil {
		updateRequest.PriceCalculationRevenueSharePercentage = *input.PriceCalculationRevenueSharePercentage
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE_CALCULATION_REVENUE_SHARE_PERCENTAGE)
	}
	if input.PriceCalculationType != nil {
		updateRequest.PriceCalculationType = input.PriceCalculationType.String()
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE_CALCULATION_TYPE)
	}
	if input.ConditionalsMinimumChargeAmount != nil {
		updateRequest.ConditionalsMinimumChargeAmount = *input.ConditionalsMinimumChargeAmount
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_CONDITIONALS_MINIMUM_CHARGE_AMOUNT)
	}
	if input.ConditionalsMinimumChargePeriod != nil {
		updateRequest.ConditionalsMinimumChargePeriod = input.ConditionalsMinimumChargePeriod.String()
		fieldsMask = append(fieldsMask, offeringpb.OfferingFieldMask_OFFERING_FIELD_CONDITIONALS_MINIMUM_CHARGE_PERIOD)
	}

	if len(fieldsMask) == 0 {
		span.LogFields(log.String("result", "No fields to update"))
		return nil
	}
	updateRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := CallEventsPlatformGRPCWithRetry[*commonpb.IdResponse](func() (*commonpb.IdResponse, error) {
		return s.grpcClients.OfferingClient.UpdateOffering(ctx, &updateRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}
