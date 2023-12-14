package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	entityDashboard "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity/dashboard"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"math/rand"
	"time"
)

type DashboardViewOrganizationsRequest struct {
	Where *model.Filter
	Sort  *model.SortBy
	Page  int
	Limit int
}

type DashboardService interface {
	GetDashboardViewOrganizationsData(ctx context.Context, requestDetails DashboardViewOrganizationsRequest) (*utils.Pagination, error)

	GetDashboardCustomerMapData(ctx context.Context) ([]*entityDashboard.DashboardCustomerMapData, error)
	GetDashboardMRRPerCustomerData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardDashboardMRRPerCustomerData, error)
	GetDashboardGrossRevenueRetentionData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardGrossRevenueRetentionData, error)
	GetDashboardARRBreakdownData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardARRBreakdownData, error)
	GetDashboardRevenueAtRiskData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardRevenueAtRiskData, error)
	GetDashboardRetentionRateData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardRetentionRateData, error)
	GetDashboardNewCustomersData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardNewCustomersData, error)
}

type dashboardService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewDashboardService(log logger.Logger, repositories *repository.Repositories, services *Services) DashboardService {
	return &dashboardService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *dashboardService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *dashboardService) GetDashboardViewOrganizationsData(ctx context.Context, requestDetails DashboardViewOrganizationsRequest) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardViewOrganizationsData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("page", requestDetails.Page), log.Int("limit", requestDetails.Limit))
	if requestDetails.Where != nil {
		span.LogFields(log.Object("filter", *requestDetails.Where))
	}
	if requestDetails.Sort != nil {
		span.LogFields(log.Object("sort", *requestDetails.Sort))
	}

	var paginatedResult = utils.Pagination{
		Limit: requestDetails.Limit,
		Page:  requestDetails.Page,
	}

	dbNodes, err := s.repositories.DashboardRepository.GetDashboardViewOrganizationData(ctx, common.GetContext(ctx).Tenant, paginatedResult.GetSkip(), paginatedResult.GetLimit(), requestDetails.Where, requestDetails.Sort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodes.Count)

	organizationEntities := entity.OrganizationEntities{}

	for _, v := range dbNodes.Nodes {
		organizationEntities = append(organizationEntities, *s.services.OrganizationService.mapDbNodeToOrganizationEntity(*v))
	}

	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *dashboardService) GetDashboardCustomerMapData(ctx context.Context) ([]*entityDashboard.DashboardCustomerMapData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardCustomerMapData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	response := make([]*entityDashboard.DashboardCustomerMapData, 0)

	data, err := s.repositories.DashboardRepository.GetDashboardCustomerMapData(ctx, common.GetContext(ctx).Tenant)
	if err != nil {
		return nil, err
	}

	for _, record := range data {
		organizationId, _ := record["organizationId"].(string)
		oldestServiceStartedAt, _ := record["oldestServiceStartedAt"].(time.Time)
		arr, _ := record["arr"].(float64)
		state, _ := record["state"].(string)

		response = append(response, &entityDashboard.DashboardCustomerMapData{
			OrganizationId:     organizationId,
			ContractSignedDate: oldestServiceStartedAt,
			State:              mapDashboardCustomerMapStateFromString(state),
			Arr:                arr,
		})
	}

	return response, nil
}

func mapDashboardCustomerMapStateFromString(state string) entityDashboard.DashboardCustomerMapState {
	switch state {
	case "OK":
		return entityDashboard.DashboardCustomerMapStateOk
	case "AT_RISK":
		return entityDashboard.DashboardCustomerMapStateAtRisk
	case "CHURNED":
		return entityDashboard.DashboardCustomerMapStateChurned
	default:
		return ""
	}
}

func (s *dashboardService) GetDashboardMRRPerCustomerData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardDashboardMRRPerCustomerData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardMRRPerCustomerData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("start", start))
	span.LogFields(log.Object("end", end))

	response := entityDashboard.DashboardDashboardMRRPerCustomerData{}

	countCustomers, err := s.repositories.OrganizationRepository.CountCustomers(ctx, common.GetContext(ctx).Tenant)
	if err != nil {
		return nil, err
	}

	data, err := s.repositories.DashboardRepository.GetDashboardMRRPerCustomerData(ctx, common.GetContext(ctx).Tenant, start, end)
	if err != nil {
		return nil, err
	}

	for _, record := range data {
		year, _ := record["year"].(int64)
		month, _ := record["month"].(int64)
		amountPerMonth, _ := record["amountPerMonth"].(float64)

		if amountPerMonth > 0 && countCustomers > 0 {
			amountPerMonth = amountPerMonth / float64(countCustomers)
		}

		newData := &entityDashboard.DashboardDashboardMRRPerCustomerPerMonthData{
			Year:  int(year),
			Month: int(month),
			Value: amountPerMonth,
		}

		response.Months = append(response.Months, newData)
	}

	if len(response.Months) == 0 {
		response.IncreasePercentage = 0
	} else if len(response.Months) == 1 {
		response.IncreasePercentage = 0
	} else {
		//lastMonthMrrPerCustomer := response.Months[len(response.Months)-1].Value
		//previousMonthMrrPerCustomer := response.Months[len(response.Months)-2].Value

		//var percentageDifference float64
		//if previousMonthMrrPerCustomer != 0 {
		//	percentageDifference = float64((lastMonthMrrPerCustomer - previousMonthMrrPerCustomer) * 100 / previousMonthMrrPerCustomer)
		//} else {
		//	if lastMonthMrrPerCustomer != 0 {
		//		percentageDifference = float64(100)
		//	} else {
		//		percentageDifference = float64(0)
		//	}
		//}

		response.IncreasePercentage = 0
	}

	return &response, nil
}

func (s *dashboardService) GetDashboardGrossRevenueRetentionData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardGrossRevenueRetentionData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardGrossRevenueRetentionData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("start", start))
	span.LogFields(log.Object("end", end))

	response := entityDashboard.DashboardGrossRevenueRetentionData{}

	response.GrossRevenueRetention = 85
	response.IncreasePercentage = 5.4

	min := float64(0)
	max := float64(1)
	for i := 1; i <= 12; i++ {
		response.Months = append(response.Months, &entityDashboard.DashboardGrossRevenueRetentionPerMonthData{
			Month:      i,
			Percentage: rand.Float64()*(max-min) + min,
		})
	}

	return &response, nil
}

func (s *dashboardService) GetDashboardARRBreakdownData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardARRBreakdownData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardARRBreakdownData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("start", start))
	span.LogFields(log.Object("end", end))

	response := entityDashboard.DashboardARRBreakdownData{}

	response.ArrBreakdown = 0
	response.IncreasePercentage = 0

	data, err := s.repositories.DashboardRepository.GetDashboardARRBreakdownData(ctx, common.GetContext(ctx).Tenant, start, end)
	if err != nil {
		return nil, err
	}

	for _, record := range data {
		year, _ := record["year"].(int64)
		month, _ := record["month"].(int64)
		newlyContracted, _ := record["newlyContracted"].(float64)
		renewals, _ := record["renewals"].(float64)
		upsells, _ := record["upsells"].(float64)
		downgrades, _ := record["downgrades"].(float64)
		cancellations, _ := record["cancellations"].(float64)
		churned, _ := record["churned"].(float64)

		newData := &entityDashboard.DashboardARRBreakdownPerMonthData{
			Year:            int(year),
			Month:           int(month),
			NewlyContracted: newlyContracted,
			Renewals:        renewals,
			Upsells:         upsells,
			Downgrades:      downgrades,
			Cancellations:   cancellations,
			Churned:         churned,
		}

		response.Months = append(response.Months, newData)
	}

	upsells, err := s.repositories.DashboardRepository.GetDashboardARRBreakdownUpsellsAndDowngradesData(ctx, common.GetContext(ctx).Tenant, "UPSELLS", start, end)
	if err != nil {
		return nil, err
	}

	for _, record := range upsells {
		year, _ := record["year"].(int64)
		month, _ := record["month"].(int64)
		value, _ := record["value"].(float64)

		for _, monthData := range response.Months {
			if monthData.Year == int(year) && monthData.Month == int(month) {
				monthData.Upsells = value
			}
		}
	}

	downgrades, err := s.repositories.DashboardRepository.GetDashboardARRBreakdownUpsellsAndDowngradesData(ctx, common.GetContext(ctx).Tenant, "DOWNGRADES", start, end)
	if err != nil {
		return nil, err
	}

	for _, record := range downgrades {
		year, _ := record["year"].(int64)
		month, _ := record["month"].(int64)
		value, _ := record["value"].(float64)

		for _, monthData := range response.Months {
			if monthData.Year == int(year) && monthData.Month == int(month) {
				monthData.Downgrades = value
			}
		}
	}

	renewals, err := s.repositories.DashboardRepository.GetDashboardARRBreakdownRenewalsData(ctx, common.GetContext(ctx).Tenant, start, end)
	if err != nil {
		return nil, err
	}

	for _, record := range renewals {
		year, _ := record["year"].(int64)
		month, _ := record["month"].(int64)
		value, _ := record["value"].(float64)

		for _, monthData := range response.Months {
			if monthData.Year == int(year) && monthData.Month == int(month) {
				monthData.Renewals = value
			}
		}
	}

	return &response, nil
}

func (s *dashboardService) GetDashboardRevenueAtRiskData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardRevenueAtRiskData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardRevenueAtRiskData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("start", start))
	span.LogFields(log.Object("end", end))

	response := entityDashboard.DashboardRevenueAtRiskData{}

	data, err := s.repositories.DashboardRepository.GetDashboardRevenueAtRiskData(ctx, common.GetContext(ctx).Tenant, start, end)
	if err != nil {
		return nil, err
	}

	high, _ := data[0]["high"].(float64)
	atRisk, _ := data[0]["atRisk"].(float64)

	response.HighConfidence = high
	response.AtRisk = atRisk

	return &response, nil
}

func (s *dashboardService) GetDashboardRetentionRateData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardRetentionRateData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardRetentionRateData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("start", start))
	span.LogFields(log.Object("end", end))

	response := entityDashboard.DashboardRetentionRateData{}

	response.RetentionRate = 86
	response.IncreasePercentage = 2.6

	min := 2
	max := 5
	for i := 1; i <= 12; i++ {
		response.Months = append(response.Months, &entityDashboard.DashboardRetentionRatePerMonthData{
			Month:      i,
			RenewCount: i*(rand.Intn(max-min)+min) + 7,
			ChurnCount: i + 2,
		})
	}

	return &response, nil
}

func (s *dashboardService) GetDashboardNewCustomersData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardNewCustomersData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DashboardService.GetDashboardNewCustomersData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("start", start))
	span.LogFields(log.Object("end", end))

	response := entityDashboard.DashboardNewCustomersData{}

	data, err := s.repositories.DashboardRepository.GetDashboardNewCustomersData(ctx, common.GetContext(ctx).Tenant, start, end)
	if err != nil {
		return nil, err
	}

	for _, record := range data {
		year, _ := record["year"].(int64)
		month, _ := record["month"].(int64)
		count, _ := record["count"].(int64)

		newData := &entityDashboard.DashboardNewCustomerMonthData{
			Year:  int(year),
			Month: int(month),
			Count: int(count),
		}

		response.Months = append(response.Months, newData)
	}

	//currentMonthCount := response.Months[len(response.Months)-1].Count
	//previousMonthCount := response.Months[len(response.Months)-2].Count

	//var percentageDifference float64
	//if previousMonthCount != 0 {
	//	percentageDifference = float64((currentMonthCount - previousMonthCount) * 100 / previousMonthCount)
	//} else {
	//	if currentMonthCount != 0 {
	//		percentageDifference = float64(100)
	//	} else {
	//		percentageDifference = float64(0)
	//	}
	//}

	//TODO fix this when we know what we want to show
	response.ThisMonthIncreasePercentage = 0
	response.ThisMonthCount = response.Months[len(response.Months)-1].Count

	return &response, nil

}
