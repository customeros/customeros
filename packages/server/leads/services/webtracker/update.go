package webtracker

import (
	"context"

	"github.com/customeros/customeros/packages/server/leads/internal/models"
)

func (s *webtrackerService) UpdateCNAMEHost(ctx context.Context, cnameHost string) (*models.WebTracker, error) {
	return nil, nil
}

func (s *webtrackerService) IsCNAMEConfigured(ctx context.Context, webtrackerID string) (bool, error) {
	return false, nil
}

func (s *webtrackerService) ArchiveWebtracker(ctx context.Context, webtrackerID string) error {
	return nil
}

func (s *webtrackerService) GetWebtracker(ctx context.Context, webtrackerID string) (*models.WebTracker, error) {
	return nil, nil
}

func (s *webtrackerService) GetActiveWebtrackers(ctx context.Context) ([]models.WebTracker, error) {
	return nil, nil
}
