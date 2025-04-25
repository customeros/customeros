package webtracker

import (
	"context"

	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
)

type WebtrackerService interface {
	CreateWebtracker(ctx context.Context, webtracker *models.WebTracker) (*models.WebTracker, error)
	// UpdateWebtracker(ctx context.Context, updateWebtracker *dto.WebTrackerUpdate) (*models.WebTracker, error)
	// UpdateCNAMEHost(ctx context.Context, cnameHost string) (*models.WebTracker, error)
	// ArchiveWebtracker(ctx context.Context, webtrackerID string) error
	// GetWebtracker(ctx context.Context, webtrackerID string) (*models.WebTracker, error)
	GetWebtrackerByOrigin(ctx context.Context, origin string) (*models.WebTracker, error)
	// GetActiveWebtrackers(ctx context.Context) ([]models.WebTracker, error)
	// IsCNAMEConfigured(ctx context.Context, webtrackerID string) (bool, error)
}

type webtrackerService struct {
	db           *gorm.DB
	repositories *repository.Repositories
}

func NewWebtrackerService(leadsDB *gorm.DB, repos *repository.Repositories) WebtrackerService {
	return &webtrackerService{
		db:           leadsDB,
		repositories: repos,
	}
}
