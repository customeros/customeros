package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type IpIntelligenceService interface {
	LookupIp(ctx context.Context, ip string) (*model.IpLookupData, error)
}

type ipIntelligenceService struct {
	config   *config.Config
	Services *Services
	log      logger.Logger
}

func NewIpIntelligenceService(config *config.Config, services *Services, log logger.Logger) IpIntelligenceService {
	return &ipIntelligenceService{
		config:   config,
		Services: services,
		log:      log,
	}
}

func (s *ipIntelligenceService) LookupIp(ctx context.Context, ip string) (*model.IpLookupData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IpIntelligenceService.LookupIp")
	defer span.Finish()
	span.LogFields(log.String("ip", ip))

	result := model.IpLookupData{
		Ip: ip,
	}

	return &result, nil
}
