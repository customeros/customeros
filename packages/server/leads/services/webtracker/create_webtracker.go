package webtracker

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

const (
	DEFAULT_CNAME_HOST  = "cos"
	CNAME_TARGET_DOMAIN = "custoscdn.com"
)

func (s *webtrackerService) CreateWebtracker(ctx context.Context, webtracker *models.WebTracker) (*models.WebTracker, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "websiteRegistrationService.CreateWebtracker")
	defer span.Finish()

	err := validateCreateWebtrackerRequest(webtracker)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	webTrackerRecord := buildWebTrackerRecord(webtracker)

	// Use transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Create webtracker
		err := s.repositories.WebTracker.CreateWithTxn(ctx, tx, webTrackerRecord)
		if err != nil {
			span.TraceError(err)
			return err
		}

		// Create outbox event
		eventPayload, err := s.buildOutboxEventPayload(ctx, webTrackerRecord)
		if err != nil {
			span.TraceError(err)
			return err
		}

		event := utils.BuildOutboxEvent(enum.EventWebtrackerCreated, eventPayload)

		err = s.repositories.Outbox.CreateWithTxn(ctx, tx, event)
		if err != nil {
			span.TraceError(err)
			return err
		}

		return nil
	})
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	return webTrackerRecord, nil
}

func (s *webtrackerService) buildOutboxEventPayload(ctx context.Context, webtracker *models.WebTracker) ([]byte, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webtrackerService.buildOutboxEventPayload")
	defer span.Finish()

	data, err := proto.Marshal(&pb.WebTrackerCreated{
		Id:          webtracker.ID,
		Domain:      webtracker.Domain,
		CnameHost:   webtracker.CNAMEHost,
		CnameTarget: webtracker.CNAMETarget,
	})
	if err != nil {
		span.TraceError(err)
		return nil, fmt.Errorf("failed to marchal webtracker event: %w", err)
	}
	return data, nil
}

func buildWebTrackerRecord(webtracker *models.WebTracker) *models.WebTracker {
	var cnameHost string

	if webtracker.CNAMEHost == "" {
		cnameHost = DEFAULT_CNAME_HOST
	} else {
		cnameHost = webtracker.CNAMEHost
	}

	return &models.WebTracker{
		ID:                utils.GenerateNanoIDWithPrefix("trkr", 16),
		Domain:            webtracker.Domain,
		CNAMEHost:         cnameHost,
		CNAMETarget:       fmt.Sprintf("%s.%s", utils.GenerateNanoID(9), CNAME_TARGET_DOMAIN),
		IsCNAMEConfigured: false,
		IsProxyActive:     false,
		IsArchived:        false,
		CreatedAt:         utils.Now(),
	}
}

func validateCreateWebtrackerRequest(webtracker *models.WebTracker) error {
	switch {
	case webtracker == nil:
		err := errors.New("Webtracker is empty")
		return err
	case webtracker.Domain == "":
		err := errors.New("Domain not set")
		return err
	default:
		return nil
	}
}
