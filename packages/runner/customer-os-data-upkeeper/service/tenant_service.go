package service

import (
	"context"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
)

type TenantService interface {
	CheckOnboarding()
}

type tenantService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonService.CommonServices
}

func NewTenantService(cfg *config.Config, log logger.Logger, commonServices *commonService.CommonServices) TenantService {
	return &tenantService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *tenantService) CheckOnboarding() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "TenantService.CheckOnboarding")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 50
	delayFromLastCheckHours := 24

	tenantDbNodes, err := s.commonServices.Neo4jRepositories.TenantReadRepository.GetTenantsForOnboardingCheck(ctx, limit, delayFromLastCheckHours)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting tenants for onboarding check"))
		s.log.Errorf("Error getting tenants for onboarding check: %s", err.Error())
		return
	}

	if len(tenantDbNodes) == 0 {
		return
	}

	for _, tenantDbNode := range tenantDbNodes {
		recordSpan, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.CheckOnboarding.Record")
		defer recordSpan.Finish()

		tenantEntity := neo4jmapper.MapDbNodeToTenantEntity(tenantDbNode)
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenantEntity.Name,
			AppSource: constants.AppSourceDataUpkeeper,
		})
		tracing.TagTenant(recordSpan, tenantEntity.Name)

		// mark tenant as checked
		err = s.commonServices.Neo4jRepositories.TenantWriteRepository.MarkOnboardingChecked(innerCtx, tenantEntity.Name)
		if err != nil {
			tracing.TraceErr(recordSpan, errors.Wrap(err, "error marking tenant as checked"))
			s.log.Errorf("Error marking tenant as checked: %s", err.Error())
			continue
		}

		// check web visit agents
		s.checkWebVisitorAgents(innerCtx, tenantEntity.Name)
	}
}

func (s *tenantService) checkWebVisitorAgents(ctx context.Context, tenant string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.CheckOnboarding.CheckWebVisitorAgents")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	tracing.TagTenant(span, tenant)

	// get web visitor agents
	webVisitorAgents, err := s.commonServices.PostgresRepositories.AgentsRepository.GetAllAgentsByTypes(ctx, []enum.AgentType{enum.AgentVisitorID})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting web visitor agents"))
		s.log.Errorf("Error getting web visitor agents: %s", err.Error())
		return
	}

	if len(webVisitorAgents) > 0 {
		span.LogFields(log.Bool("result.onboarded", true))
		return
	}

	// create web visitor agents
	_, err = s.commonServices.AgentService.CreateAgent(ctx, enum.AgentVisitorID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error creating web visitor agent"))
		s.log.Errorf("Error creating web visitor agents: %s", err.Error())
	}
}
