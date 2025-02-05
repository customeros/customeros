package service

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	commonservice "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
)

type WebSesssionService interface {
	ProcessIntentSignals()
}

type webSessionService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonservice.CommonServices
}

func NewWebSessionService(cfg *config.Config, log logger.Logger, s *commonservice.CommonServices) WebSesssionService {
	return &webSessionService{
		cfg:            cfg,
		log:            log,
		commonServices: s,
	}
}

func (s *webSessionService) ProcessIntentSignals() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	span, ctx := tracing.StartTracerSpan(ctx, "WebSessionService.ProcessIntentSignals")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	// find closed events that have not been analyzed for intent
	sessions, err := s.commonServices.PostgresRepositories.WebSessionRepository.FindAllSessionsForIntentAnalysis(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	if len(sessions) == 0 {
		return
	}

	for _, session := range sessions {
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    session.Tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})
		err = s.commonServices.WebVisitProcessor.Process(innerCtx, session)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}
	}
}
