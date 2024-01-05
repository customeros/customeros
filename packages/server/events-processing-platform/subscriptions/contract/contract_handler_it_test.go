package contract

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	opportunitymodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	servicelineitemmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/mocked_grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestContractEventHandler_UpdateRenewalNextCycleDate_CreateRenewalOpportunityWhenMissing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
	})

	// prepare grpc client
	calledEventsPlatformToCreateRenewalOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		CreateRenewalOpportunity: func(context context.Context, op *opportunitypb.CreateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, "", op.LoggedInUserId)
			require.Equal(t, contractId, op.ContractId)
			require.Nil(t, op.CreatedAt)
			require.Nil(t, op.UpdatedAt)
			require.Equal(t, opportunitypb.RenewalLikelihood_HIGH_RENEWAL, op.RenewalLikelihood)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			calledEventsPlatformToCreateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: "some-opportunity-id",
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToCreateRenewalOpportunity)
}

func TestContractEventHandler_UpdateRenewalNextCycleDate_MonthlyContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.AppSource)
			require.Equal(t, startOfNextMonth(utils.Now()), *utils.TimestampProtoToTimePtr(op.RenewedAt))
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)
}

func TestContractEventHandler_UpdateRenewalNextCycleDate_QuarterlyContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	yesterday := utils.Now().AddDate(0, 0, -1)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: &yesterday,
		RenewalCycle:     string(model.QuarterlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.AppSource)
			in1Quarter := yesterday.AddDate(0, 3, 0)
			require.Equal(t, in1Quarter, *utils.TimestampProtoToTimePtr(op.RenewedAt))
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)
}

func TestContractEventHandler_UpdateRenewalNextCycleDate_AnnualContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare grpc client
	calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.AppSource)
			require.Equal(t, startOfNextYear(utils.Now()), *utils.TimestampProtoToTimePtr(op.RenewedAt))
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)
}

func TestContractEventHandler_UpdateRenewalNextCycleDate_MultiAnnualContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	yesterday := utils.Now().AddDate(0, 0, -1)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: &yesterday,
		RenewalCycle:     string(model.AnnuallyRenewalCycleString),
		RenewalPeriods:   utils.Int64Ptr(10),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare grpc client
	calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunityNextCycleDate: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.OpportunityId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.AppSource)
			in10Years := yesterday.AddDate(10, 0, 0)
			require.Equal(t, in10Years, *utils.TimestampProtoToTimePtr(op.RenewedAt))
			calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunityNextCycleDate)
}

func startOfNextMonth(current time.Time) time.Time {
	year, month, _ := current.Date()
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, current.Location())

	// Handle December to January transition
	if month == time.December {
		firstOfNextMonth = time.Date(year+1, time.January, 1, 0, 0, 0, 0, current.Location())
	}
	return firstOfNextMonth
}

func startOfNextYear(current time.Time) time.Time {
	year, _, _ := current.Date()
	firstOfNextYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, current.Location())
	return firstOfNextYear
}

func TestContractEventHandler_UpdateRenewalArrForecast_OnlyOnceBilled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(10),
		Billed:    string(servicelineitemmodel.OnceBilledString),
		CreatedAt: utils.Now(),
	})

	// prepare grpc client
	calledEventsPlatformToUpdateOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateOpportunity: func(context context.Context, op *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, float64(0), op.Amount)
			require.Equal(t, float64(0), op.MaxAmount)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT}, op.FieldsMask)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, constants.SourceOpenline, op.SourceFields.Source)
			calledEventsPlatformToUpdateOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestContractEventHandler_UpdateRenewalArrForecast_MultipleServices(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000000),
		Billed:    string(servicelineitemmodel.OnceBilledString),
		CreatedAt: utils.Now(),
	})
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(10),
		Quantity:  int64(5),
		Billed:    string(servicelineitemmodel.MonthlyBilledString),
		CreatedAt: utils.Now(),
	})
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000),
		Quantity:  int64(2),
		Billed:    string(servicelineitemmodel.AnnuallyBilledString),
		CreatedAt: utils.Now(),
	})

	// prepare grpc client
	calledEventsPlatformToUpdateOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateOpportunity: func(context context.Context, op *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, float64(2600), op.Amount)
			require.Equal(t, float64(2600), op.MaxAmount)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT}, op.FieldsMask)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, constants.SourceOpenline, op.SourceFields.Source)
			calledEventsPlatformToUpdateOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestContractEventHandler_UpdateRenewalArrForecast_MediumLikelihood(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringMedium),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000),
		Quantity:  int64(4),
		Billed:    string(servicelineitemmodel.AnnuallyBilledString),
		CreatedAt: utils.Now(),
	})

	// prepare grpc client
	calledEventsPlatformToUpdateOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateOpportunity: func(context context.Context, op *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, float64(2000), op.Amount)
			require.Equal(t, float64(4000), op.MaxAmount)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT}, op.FieldsMask)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, constants.SourceOpenline, op.SourceFields.Source)
			calledEventsPlatformToUpdateOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestContractEventHandler_UpdateRenewalArrForecast_ContractEndsBeforeNextRenewal(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	in5Minutes := utils.Now().Add(time.Minute * 5)
	in10Minutes := utils.Now().Add(time.Minute * 10)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.AnnuallyRenewalCycleString),
		EndedAt:          utils.TimePtr(in5Minutes),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringMedium),
			RenewedAt:         utils.TimePtr(in10Minutes),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000),
		Quantity:  int64(1),
		Billed:    string(servicelineitemmodel.AnnuallyBilledString),
		CreatedAt: utils.Now(),
	})

	// prepare grpc client
	calledEventsPlatformToUpdateOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateOpportunity: func(context context.Context, op *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, float64(0), op.Amount)
			require.Equal(t, float64(0), op.MaxAmount)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT}, op.FieldsMask)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, constants.SourceOpenline, op.SourceFields.Source)
			calledEventsPlatformToUpdateOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestContractEventHandler_UpdateRenewalArrForecast_ContractEndsIn6Months_ProrateAmount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	in6Months := utils.Now().AddDate(0, 6, 0)
	in1Month := utils.Now().AddDate(0, 1, 0)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
		EndedAt:          utils.TimePtr(in6Months),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
			RenewedAt:         utils.TimePtr(in1Month),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000),
		Quantity:  int64(1),
		Billed:    string(servicelineitemmodel.AnnuallyBilledString),
		CreatedAt: utils.Now(),
	})

	// prepare grpc client
	calledEventsPlatformToUpdateOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateOpportunity: func(context context.Context, op *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, float64(500), op.Amount)
			require.Equal(t, float64(500), op.MaxAmount)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT}, op.FieldsMask)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, constants.SourceOpenline, op.SourceFields.Source)
			calledEventsPlatformToUpdateOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestContractEventHandler_UpdateRenewalArrForecast_ContractEndsInMoreThan12Months_FullAmount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	in13Months := utils.Now().AddDate(0, 13, 0)
	in12Months := utils.Now().AddDate(1, 0, 0)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	serviceStartedAt, _ := utils.UnmarshalDateTime("2021-01-01T00:00:00Z")
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		ServiceStartedAt: serviceStartedAt,
		RenewalCycle:     string(model.MonthlyRenewalCycleString),
		EndedAt:          utils.TimePtr(in13Months),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
			RenewedAt:         utils.TimePtr(in12Months),
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)
	neo4jt.CreateServiceLineItemForContract(ctx, testDatabase.Driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     float64(1000),
		Quantity:  int64(1),
		Billed:    string(servicelineitemmodel.AnnuallyBilledString),
		CreatedAt: utils.Now(),
	})

	// prepare grpc client
	calledEventsPlatformToUpdateOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateOpportunity: func(context context.Context, op *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, float64(1000), op.Amount)
			require.Equal(t, float64(1000), op.MaxAmount)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT}, op.FieldsMask)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, constants.SourceOpenline, op.SourceFields.Source)
			calledEventsPlatformToUpdateOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityArr(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateOpportunity)
}

func TestContractEventHandler_UpdateActiveRenewalOpportunityLikelihood_EndedContract_UpdateLikelihood(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	tomorrow := utils.Now().AddDate(0, 0, 1)
	afterTomorrow := utils.Now().AddDate(0, 0, 2)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		EndedAt:      &tomorrow,
		RenewalCycle: string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringLow),
			RenewedAt:         &afterTomorrow,
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare grpc mock
	calledEventsPlatformToUpdateRenewalOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunity: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, opportunitypb.RenewalLikelihood_ZERO_RENEWAL, op.RenewalLikelihood)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD}, op.FieldsMask)
			calledEventsPlatformToUpdateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityLikelihood(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunity)
}

func TestContractEventHandler_UpdateActiveRenewalOpportunityLikelihood_EndedContract_LikelihoodAlreadyZero(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	tomorrow := utils.Now().AddDate(0, 0, 1)
	afterTomorrow := utils.Now().AddDate(0, 0, 2)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		EndedAt:      &tomorrow,
		RenewalCycle: string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringZero),
			RenewedAt:         &afterTomorrow,
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityLikelihood(ctx, tenantName, contractId)
	require.Nil(t, err)
}

func TestContractEventHandler_UpdateActiveRenewalOpportunityLikelihood_ReinitiatedContract_UpdateLikelihood(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	afterTomorrow := utils.Now().AddDate(0, 0, 2)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle: string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringZero),
			RenewedAt:         &afterTomorrow,
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare grpc mock
	calledEventsPlatformToUpdateRenewalOpportunity := false
	opportunityCallbacks := mocked_grpc.MockOpportunityServiceCallbacks{
		UpdateRenewalOpportunity: func(context context.Context, op *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, opportunityId, op.Id)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, op.SourceFields.AppSource)
			require.Equal(t, opportunitypb.RenewalLikelihood_MEDIUM_RENEWAL, op.RenewalLikelihood)
			require.Equal(t, []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD}, op.FieldsMask)
			calledEventsPlatformToUpdateRenewalOpportunity = true
			return &opportunitypb.OpportunityIdGrpcResponse{
				Id: opportunityId,
			}, nil
		},
	}
	mocked_grpc.SetOpportunityCallbacks(&opportunityCallbacks)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityLikelihood(ctx, tenantName, contractId)
	require.Nil(t, err)

	// Check
	require.True(t, calledEventsPlatformToUpdateRenewalOpportunity)
}

func TestContractEventHandler_UpdateActiveRenewalOpportunityLikelihood_ReinitiatedContract_LikelihoodNotZero(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	afterTomorrow := utils.Now().AddDate(0, 0, 2)
	contractId := neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		RenewalCycle: string(model.AnnuallyRenewalCycleString),
	})
	opportunityId := neo4jt.CreateOpportunity(ctx, testDatabase.Driver, tenantName, entity.OpportunityEntity{
		InternalType:  string(opportunitymodel.OpportunityInternalTypeStringRenewal),
		InternalStage: string(opportunitymodel.OpportunityInternalStageStringOpen),
		RenewalDetails: entity.RenewalDetails{
			RenewalLikelihood: string(opportunitymodel.RenewalLikelihoodStringHigh),
			RenewedAt:         &afterTomorrow,
		},
	})
	neo4jt.LinkContractWithOpportunity(ctx, testDatabase.Driver, contractId, opportunityId, true)

	// prepare event handler
	handler := NewContractHandler(testLogger, testDatabase.Repositories, testMockedGrpcClient)

	// EXECUTE
	err := handler.UpdateActiveRenewalOpportunityLikelihood(ctx, tenantName, contractId)
	require.Nil(t, err)
}
