package service

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"

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

	limit := 10

	// get companies
	orgs, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetOrganizationsToFetchLogo(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, org := range orgs {
		// download icon
		if org.IconUrl != "" {
			iconPath := fmt.Sprintf("%s/%s", org.PrimaryDomain, "icon")
			s.downloadImage(ctx, org.ID, org.IconUrl, iconPath, "icon")
		}

		// download logo
		logoPath := fmt.Sprintf("%s/%s", org.PrimaryDomain, "logo")
		clearbitUrl := "https://logo.clearbit.com/" + org.PrimaryDomain
		if !s.downloadImage(ctx, org.ID, clearbitUrl, logoPath, "logo") && org.LogoUrl != "" {
			s.downloadImage(ctx, org.ID, org.LogoUrl, logoPath, "logo")
		}

	}
}

func (s *mediaService) downloadImage(ctx context.Context, globalOrgId uint64, imageUrl, imagePath, imageType string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MediaService.downloadImage")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	span.LogFields(log.String("imageUrl", imageUrl), log.String("imagePath", imagePath), log.String("imageType", imageType), log.Uint64("globalOrgId", globalOrgId))

	imagePath, err := s.commonServices.MediaService.DownloadImageToS3(ctx, imageUrl, BUCKET, imagePath)
	if err != nil {
		if !errors.Is(err, coserrors.ErrResourceNotFound) {
			tracing.TraceErr(span, err)
		}
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetDownloadStatus(ctx, globalOrgId, enum.DownloadError)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set download status"))
		}
		return false
	}

	if imagePath != "" {
		if imageType == "icon" {
			err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetIcon(ctx, globalOrgId, imagePath)
		} else {
			err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetLogo(ctx, globalOrgId, imagePath)
		}
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set image path"))
			return false
		}
	}
	return true
}
