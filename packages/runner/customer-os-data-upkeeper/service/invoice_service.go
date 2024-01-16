package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/events_processing_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
)

type InvoiceService interface {
	GenerateInvoices()
}

type invoiceService struct {
	cfg                    *config.Config
	log                    logger.Logger
	repositories           *repository.Repositories
	eventsProcessingClient *events_processing_client.Client
}

func NewInvoiceService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, client *events_processing_client.Client) InvoiceService {
	return &invoiceService{
		cfg:                    cfg,
		log:                    log,
		repositories:           repositories,
		eventsProcessingClient: client,
	}
}

func (s *invoiceService) GenerateInvoices() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "InvoiceService.GenerateInvoices")
	defer span.Finish()

	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil. Will not update next cycle date.")
		return
	}

	//now := utils.Now()

	//organizationsForInvoicing, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsForInvoicing(ctx, now)
	//if err != nil {
	//	s.log.Error("failed to get organizations for invoicing", err)
	//	return
	//}
	//
	//for _, dbNode := range organizationsForInvoicing {
	//	props := utils.GetPropsFromNode(*dbNode)
	//
	//	tenant := utils.GetStringPropOrEmpty(props, "tenant")
	//	organizationId := utils.GetStringPropOrEmpty(props, "organizationId")
	//
	//	_, err = s.eventsProcessingClient.InvoiceClient.NewInvoice(ctx, &invoicepb.NewInvoiceRequest{
	//		Tenant:         tenant,
	//		OrganizationId: organizationId,
	//		SourceFields: &commonpb.SourceFields{
	//			AppSource: constants.AppSourceDataUpkeeper,
	//		},
	//	})
	//	if err != nil {
	//		tracing.TraceErr(span, err)
	//		s.log.Errorf("Error invoicing organization {%s} in tenant {%s}: %s", organizationId, tenant, err.Error())
	//	}
	//}
	//
	//s.repositories.Neo4jRepositories.OrganizationWriteRepository.UpdateInvoicingActive(ctx, tenant, organizationId, true, now)
}
