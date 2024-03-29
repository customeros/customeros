package resolver

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	offeringpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/offering"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_OfferingCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	offeringId := uuid.New().String()

	calledCreateOffering := false

	offeringCallbacks := events_platform.MockOfferingServiceCallbacks{
		CreateOffering: func(context context.Context, offering *offeringpb.CreateOfferingGrpcRequest) (*commonpb.IdResponse, error) {

			require.Equal(t, "1", offering.Name)
			require.Equal(t, true, offering.Active)
			require.Equal(t, "PRODUCT", offering.Type)
			require.Equal(t, "SUBSCRIPTION", offering.PricingModel)
			require.Equal(t, int64(2), offering.PricingPeriodInMonths)
			require.Equal(t, "AUD", offering.Currency)
			require.Equal(t, float64(3), offering.Price)
			require.Equal(t, true, offering.PriceCalculated)
			require.Equal(t, true, offering.Conditional)
			require.Equal(t, true, offering.Taxable)
			require.Equal(t, "REVENUE_SHARE", offering.PriceCalculationType)
			require.Equal(t, float64(4), offering.PriceCalculationRevenueSharePercentage)
			require.Equal(t, "MONTHLY", offering.ConditionalsMinimumChargePeriod)
			require.Equal(t, float64(5), offering.ConditionalsMinimumChargeAmount)

			calledCreateOffering = true

			return &commonpb.IdResponse{
				Id: offeringId,
			}, nil
		},
	}
	events_platform.SetOfferingCallbacks(&offeringCallbacks)

	rawResponse := callGraphQL(t, "offering/create_offering", map[string]interface{}{})
	require.Nil(t, rawResponse.Errors)

	var graphqlResponse struct {
		Offering_Create *string
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
	require.Nil(t, err)

	require.Equal(t, offeringId, *graphqlResponse.Offering_Create)

	require.True(t, calledCreateOffering)
}

func TestMutationResolver_OfferingUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	offeringId := uuid.New().String()

	calledUpdateOffering := false

	offeringCallbacks := events_platform.MockOfferingServiceCallbacks{
		UpdateOffering: func(context context.Context, offering *offeringpb.UpdateOfferingGrpcRequest) (*commonpb.IdResponse, error) {

			require.Equal(t, "1", offering.Name)
			require.Equal(t, true, offering.Active)
			require.Equal(t, "PRODUCT", offering.Type)
			require.Equal(t, "SUBSCRIPTION", offering.PricingModel)
			require.Equal(t, int64(2), offering.PricingPeriodInMonths)
			require.Equal(t, "AUD", offering.Currency)
			require.Equal(t, float64(3), offering.Price)
			require.Equal(t, true, offering.PriceCalculated)
			require.Equal(t, true, offering.Conditional)
			require.Equal(t, true, offering.Taxable)
			require.Equal(t, "REVENUE_SHARE", offering.PriceCalculationType)
			require.Equal(t, float64(4), offering.PriceCalculationRevenueSharePercentage)
			require.Equal(t, "MONTHLY", offering.ConditionalsMinimumChargePeriod)
			require.Equal(t, float64(5), offering.ConditionalsMinimumChargeAmount)

			//require.ElementsMatch(t, []offeringpb.OfferingFieldMask{
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_NAME,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_ACTIVE,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_TYPE,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICING_MODEL,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICING_PERIOD_IN_MONTHS,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_CURRENCY,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE_CALCULATED,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_CONDITIONAL,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_TAXABLE,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE_CALCULATION_TYPE,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_PRICE_CALCULATION_REVENUE_SHARE_PERCENTAGE,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_CONDITIONALS_MINIMUM_CHARGE_PERIOD,
			//	offeringpb.OfferingFieldMask_OFFERING_FIELD_CONDITIONALS_MINIMUM_CHARGE_AMOUNT},
			//	offering.FieldsMask)

			calledUpdateOffering = true

			return &commonpb.IdResponse{
				Id: offeringId,
			}, nil
		},
	}
	events_platform.SetOfferingCallbacks(&offeringCallbacks)

	rawResponse := callGraphQL(t, "offering/update_offering", map[string]interface{}{
		"offeringId": offeringId,
	})
	require.Nil(t, rawResponse.Errors)

	var graphqlResponse struct {
		Offering_Update *string
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
	require.Nil(t, err)

	require.Equal(t, offeringId, *graphqlResponse.Offering_Update)

	require.True(t, calledUpdateOffering)
}
