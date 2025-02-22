package service

import (
	"context"
	"fmt"

	service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
)

type MediaService interface {
	FetchAndStoreCompanyLogos()
}

type mediaService struct {
	log            logger.Logger
	commonServices *service.CommonServices
}

func NewMediaService(log logger.Logger, commonServices *service.CommonServices) MediaService {
	return &mediaService{
		log:            log,
		commonServices: commonServices,
	}
}

const (
	BUCKET = "customer-os-images"
)

func (s *mediaService) FetchAndStoreCompanyLogos() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "MediaService.FetchAndStoreCompanyLogos")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 25

	// get companies
	orgs, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetOrganizationsToFetchLogo(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, org := range orgs {
		// download logo
		logoPath := fmt.Sprintf("%s/%s", org.PrimaryDomain, "logo")
		logoPath, err = s.commonServices.MediaService.DownloadImageToS3(ctx, org.LogoUrl, BUCKET, logoPath)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetLogo(ctx, org.ID, logoPath)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}

		// download icon
		iconPath := fmt.Sprintf("%s/%s", org.PrimaryDomain, "icon")
		iconPath, err = s.commonServices.MediaService.DownloadImageToS3(ctx, org.IconUrl, BUCKET, iconPath)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetIcon(ctx, org.ID, iconPath)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}
	}
}
