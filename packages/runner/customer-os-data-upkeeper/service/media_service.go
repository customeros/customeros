package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
)

type MediaService interface {
	FetchAndStoreCompanyLogos()
	FetchAndStoreCompanyIcons()
	FetchAndStoreContactProfilePhotos()
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

func (s *mediaService) getContactHashPath(contact *postgresentity.GlobalContact) string {
	var identifier string
	if contact.LinkedInIdentifier != "" {
		identifier = contact.LinkedInIdentifier
	} else if contact.WorkEmail != "" {
		identifier = contact.WorkEmail
	} else if contact.PersonalEmail != "" {
		identifier = contact.PersonalEmail
	} else {
		return "" // No valid identifier found
	}

	// Create SHA-256 hash of the identifier
	hash := sha256.Sum256([]byte(identifier))
	// Convert to hex string
	hashStr := hex.EncodeToString(hash[:])
	return fmt.Sprintf("contacts/%s", hashStr)
}

func (s *mediaService) FetchAndStoreCompanyLogos() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "MediaService.FetchAndStoreCompanyLogos")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 5

	// get companies
	orgs, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetOrganizationsToFetchLogo(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, org := range orgs {
		// download logo
		logoPath := fmt.Sprintf("%s/%s", org.PrimaryDomain, "logo")
		clearbitUrl := "https://logo.clearbit.com/" + org.PrimaryDomain
		if !s.downloadCompanyLogo(ctx, org.ID, clearbitUrl, logoPath) && org.LogoUrl != "" {
			s.downloadCompanyLogo(ctx, org.ID, org.LogoUrl, logoPath)
		}
	}
}

func (s *mediaService) FetchAndStoreCompanyIcons() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "MediaService.FetchAndStoreCompanyIcons")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 5

	// get companies
	orgs, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetOrganizationsToFetchIcon(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, org := range orgs {
		// download icon
		iconPath := fmt.Sprintf("%s/%s", org.PrimaryDomain, "icon")
		if !s.downloadCompanyIcon(ctx, org.ID, org.IconUrl, iconPath) && org.IconUrl != "" {
			s.downloadCompanyIcon(ctx, org.ID, org.IconUrl, iconPath)
		}
	}
}

func (s *mediaService) FetchAndStoreContactProfilePhotos() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "MediaService.FetchAndStoreContactProfilePhotos")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 10

	// get contacts
	contacts, err := s.commonServices.PostgresRepositories.GlobalContactRepository.GetContactsToFetchPhoto(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, contact := range contacts {
		// download profile photo
		if contact.ProfilePhotoExternalUrl != "" {
			hashPath := s.getContactHashPath(contact)
			if hashPath == "" {
				// error, this contact has no valid identifier
				err = errors.New("no valid identifier for contact")
				tracing.TraceErr(span, err)
				continue
			}
			photoPath := fmt.Sprintf("%s/%s", hashPath, "profile")
			s.downloadContactPhoto(ctx, contact.ID, contact.ProfilePhotoExternalUrl, photoPath)
		}
	}
}

func (s *mediaService) downloadCompanyLogo(ctx context.Context, globalOrgId uint64, imageUrl, imagePath string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MediaService.downloadCompanyLogo")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	span.LogFields(log.String("imageUrl", imageUrl), log.String("imagePath", imagePath), log.Uint64("globalOrgId", globalOrgId))

	imagePath, err := s.commonServices.MediaService.DownloadImageToS3(ctx, imageUrl, BUCKET, imagePath)
	if err != nil {
		if !errors.Is(err, coserrors.ErrResourceNotFound) && !errors.Is(err, coserrors.ErrResourceForbidden) {
			tracing.TraceErr(span, err)
		}
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetDownloadStatusLogo(ctx, globalOrgId, enum.DownloadError)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set download status"))
		}
		return false
	}

	if imagePath != "" {
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetLogo(ctx, globalOrgId, imagePath)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set image path"))
			return false
		}
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetDownloadStatusLogo(ctx, globalOrgId, enum.DownloadCompleted)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set download status"))
		}
	}
	return true
}

func (s *mediaService) downloadCompanyIcon(ctx context.Context, globalOrgId uint64, imageUrl, imagePath string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MediaService.downloadCompanyIcon")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	span.LogFields(log.String("imageUrl", imageUrl), log.String("imagePath", imagePath), log.Uint64("globalOrgId", globalOrgId))

	imagePath, err := s.commonServices.MediaService.DownloadImageToS3(ctx, imageUrl, BUCKET, imagePath)
	if err != nil {
		if !errors.Is(err, coserrors.ErrResourceNotFound) && !errors.Is(err, coserrors.ErrResourceForbidden) {
			tracing.TraceErr(span, err)
		}
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetDownloadStatusIcon(ctx, globalOrgId, enum.DownloadError)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set download status"))
		}
		return false
	}

	if imagePath != "" {
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetIcon(ctx, globalOrgId, imagePath)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set image path"))
			return false
		}
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetDownloadStatusIcon(ctx, globalOrgId, enum.DownloadCompleted)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set download status"))
		}
	}
	return true
}

func (s *mediaService) downloadContactPhoto(ctx context.Context, globalContactId uint64, imageUrl, imagePath string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MediaService.downloadContactPhoto")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	span.LogFields(log.String("imageUrl", imageUrl), log.String("imagePath", imagePath), log.Uint64("globalContactId", globalContactId))

	imagePath, err := s.commonServices.MediaService.DownloadImageToS3(ctx, imageUrl, BUCKET, imagePath)
	if err != nil {
		if !errors.Is(err, coserrors.ErrResourceNotFound) && !errors.Is(err, coserrors.ErrResourceForbidden) {
			tracing.TraceErr(span, err)
		} else {
			// clear global contact profile photo external url
			globalContact, err := s.commonServices.PostgresRepositories.GlobalContactRepository.GetById(ctx, globalContactId)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			if globalContact != nil {
				globalContact.ProfilePhotoExternalUrl = ""
				_, err = s.commonServices.PostgresRepositories.GlobalContactRepository.Update(ctx, globalContact)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to update global contact"))
				}
			}
		}
		err = s.commonServices.PostgresRepositories.GlobalContactRepository.SetDownloadStatus(ctx, globalContactId, enum.DownloadError)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set download status"))
		}
		return false
	}

	if imagePath != "" {
		err = s.commonServices.PostgresRepositories.GlobalContactRepository.SetProfilePhoto(ctx, globalContactId, imagePath)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set image path"))
			return false
		}
		err = s.commonServices.PostgresRepositories.GlobalContactRepository.SetDownloadStatus(ctx, globalContactId, enum.DownloadCompleted)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to set download status"))
		}
	}
	return true
}
