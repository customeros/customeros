package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/mailstack"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
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

	spans, ctx := telemetry.StartCronSpan(ctx, "TenantService.CheckOnboarding")
	defer spans.Finish()

	limit := 50
	delayFromLastCheckHours := 24

	tenantDbNodes, err := s.commonServices.Neo4jRepositories.TenantReadRepository.GetTenantsForOnboardingCheck(ctx, limit, delayFromLastCheckHours)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error getting tenants for onboarding check"))
		s.log.Errorf("Error getting tenants for onboarding check: %s", err.Error())
		return
	}

	if len(tenantDbNodes) == 0 {
		return
	}

	for _, tenantDbNode := range tenantDbNodes {
		func(tenantDbNode *dbtype.Node) {
			recordSpans, ctx := telemetry.StartCronSpan(ctx, "TenantService.CheckOnboarding.Record", telemetry.WithNewRoot())
			defer recordSpans.Finish()

			tenantEntity := neo4jmapper.MapDbNodeToTenantEntity(tenantDbNode)
			innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    tenantEntity.Name,
				AppSource: constants.AppSourceDataUpkeeper,
			})
			recordSpans.TagTenant(tenantEntity.Name)

			// mark tenant as checked
			err = s.commonServices.Neo4jRepositories.TenantWriteRepository.MarkOnboardingChecked(innerCtx, tenantEntity.Name)
			if err != nil {
				recordSpans.TraceError(errors.Wrap(err, "error marking tenant as checked"))
				s.log.Errorf("Error marking tenant as checked: %s", err.Error())
				return
			}

			//s.checkWebVisitorAgents(innerCtx, tenantEntity.Name) // temporary disabled
			//s.checkIcpQualificationAgents(innerCtx, tenantEntity.Name) // temporary disabled
			s.checkTestMailbox(innerCtx, tenantEntity.Name)
			s.checkTenantPresentInWorkspaceSwitcher(innerCtx, tenantEntity.Name)
		}(tenantDbNode)
	}
}

func (s *tenantService) checkWebVisitorAgents(ctx context.Context, tenant string) {
	spans, ctx := telemetry.StartCronSpan(ctx, "TenantService.checkWebVisitorAgents")
	defer spans.Finish()

	// get web visitor agents
	webVisitorAgents, err := s.commonServices.PostgresRepositories.AgentRepository.GetAllAgentsByTypes(ctx, []enum.AgentType{enum.AgentWebVisitorIdentifier})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error getting web visitor agents"))
		s.log.Errorf("Error getting web visitor agents: %s", err.Error())
		return
	}

	if len(webVisitorAgents) > 0 {
		spans.LogKV("result.onboarded", true)
		return
	}

	// create web visitor agents
	_, err = s.commonServices.AgentService.CreateAgent(ctx, enum.AgentWebVisitorIdentifier)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error creating web visitor agent"))
		s.log.Errorf("Error creating web visitor agents: %s", err.Error())
	}
}

func (s *tenantService) checkIcpQualificationAgents(ctx context.Context, tenant string) {
	spans, ctx := telemetry.StartCronSpan(ctx, "TenantService.checkIcpQualificationAgents")
	defer spans.Finish()

	// get icp qualification agents
	icpQualificationAgents, err := s.commonServices.PostgresRepositories.AgentRepository.GetAllAgentsByTypes(ctx, []enum.AgentType{enum.AgentICPQualifier})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error getting icp qualification agents"))
		s.log.Errorf("Error getting icp qualification agents: %s", err.Error())
		return
	}

	if len(icpQualificationAgents) > 0 {
		spans.LogKV("result.onboarded", true)
		return
	}

	// create icp qualification agents
	_, err = s.commonServices.AgentService.CreateAgent(ctx, enum.AgentICPQualifier)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error creating icp qualification agent"))
		s.log.Errorf("Error creating icp qualification agents: %s", err.Error())
	}
}

func (s *tenantService) checkTestMailbox(ctx context.Context, tenant string) {
	spans, ctx := telemetry.StartCronSpan(ctx, "TenantService.checkTestMailbox")
	defer spans.Finish()

	// Skip if tenant has uppercase letters

	// Skip if tenant has uppercase letters
	if tenant != strings.ToLower(tenant) {
		spans.LogKV("result.skipped", true)
		spans.LogKV("reason", "tenant name was auto generated")
		return
	}

	mailboxAddress := strings.ToLower(fmt.Sprintf("%s@%s", tenant, mailstack.TEST_MAILBOX_DOMAIN))

	// Check if mailbox exists using GetMailboxes from mailstack service
	statusCode, errMsg, mailboxes, err := s.commonServices.MailstackService.GetMailboxes(ctx, tenant, mailstack.TEST_MAILBOX_DOMAIN, "")
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to get mailboxes"))
		s.log.Errorf("Error checking test mailbox: %s", err.Error())
		return
	}

	// Handle non-200 responses
	if statusCode != http.StatusOK {
		err = errors.New(errMsg)
		spans.TraceError(errors.Wrap(err, "failed to get mailboxes"))
		s.log.Errorf("Error checking test mailbox: %s", err.Error())
		return
	}

	// Check if our mailbox exists in the returned list
	mailboxExists := false
	for _, mailbox := range mailboxes {
		if strings.EqualFold(mailbox.Email, mailboxAddress) {
			mailboxExists = true
			break
		}
	}

	if mailboxExists {
		spans.LogKV("result.exists", true)
		return
	}

	testUserSetup, err := s.commonServices.RegistrationService.ConfigureTestMailbox(ctx)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to configure test mailbox"))
		s.log.Errorf("Error configuring test mailbox: %s", err.Error())
		return
	}

	if testUserSetup == nil {
		spans.LogKV("result.created", false)
		return
	}

	spans.LogKV("result.created", true)
}

func (s *tenantService) checkTenantPresentInWorkspaceSwitcher(ctx context.Context, tenant string) {
	spans, ctx := telemetry.StartCronSpan(ctx, "TenantService.checkTenantPresentInWorkspaceSwitcher")
	defer spans.Finish()

	// just assign the tenant to the platform owners, method is idempotent
	err := s.commonServices.RegistrationService.ProvideAccessToPlatformOwners(ctx, tenant)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to provide access to platform owners"))
		s.log.Errorf("Error providing access to platform owners: %s", err.Error())
	}
}
