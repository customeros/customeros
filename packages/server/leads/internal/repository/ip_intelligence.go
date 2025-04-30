package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/leads/internal/database"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
)

type IPIntelligenceRepository interface {
	Create(ctx context.Context, intel *models.IPIntelligence) error
	FindByIP(ctx context.Context, ipAddress string) (*models.IPIntelligence, error)
	Update(ctx context.Context, intel *models.IPIntelligence) error
	SetDomain(ctx context.Context, id uint, domain string) error
	SetEmailAddress(ctx context.Context, id uint, emailAddress string) error
}

type ipIntelligenceRepository struct {
	read  *gorm.DB
	write *gorm.DB
}

func NewIPIntelligenceRepository(db *database.DbConnections) IPIntelligenceRepository {
	return &ipIntelligenceRepository{
		read:  db.ReadDB,
		write: db.WriteDB,
	}
}

func (r *ipIntelligenceRepository) Create(ctx context.Context, intel *models.IPIntelligence) error {
	span, ctx := telemetry.StartPostgresSpan(ctx, "ipIntelligenceRepository.Create")
	defer span.Finish()

	return r.write.WithContext(ctx).Create(intel).Error
}

func (r *ipIntelligenceRepository) FindByIP(ctx context.Context, ipAddress string) (*models.IPIntelligence, error) {
	span, ctx := telemetry.StartPostgresSpan(ctx, "ipIntelligenceRepository.FindByIP")
	defer span.Finish()

	var intel models.IPIntelligence
	// specifically using write db connection here instead of read to ensure updates
	// are captured in real-time as this is called right after an update in Session Manager
	err := r.write.WithContext(ctx).
		Where("ip_address = ?", ipAddress).
		First(&intel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil when not found
		}
		return nil, err
	}

	return &intel, nil
}

func (r *ipIntelligenceRepository) Update(ctx context.Context, intel *models.IPIntelligence) error {
	span, ctx := telemetry.StartPostgresSpan(ctx, "ipIntelligenceRepository.Update")
	defer span.Finish()

	now := time.Now()
	intel.UpdatedAt = &now

	return r.write.WithContext(ctx).Save(intel).Error
}

func (r *ipIntelligenceRepository) SetDomain(ctx context.Context, id uint, domain string) error {
	span, ctx := telemetry.StartPostgresSpan(ctx, "ipIntelligenceRepository.SetDomain")
	defer span.Finish()

	now := time.Now()
	return r.write.WithContext(ctx).
		Model(&models.IPIntelligence{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"domain":     domain,
			"updated_at": now,
		}).Error
}

func (r *ipIntelligenceRepository) SetEmailAddress(ctx context.Context, id uint, emailAddress string) error {
	span, ctx := telemetry.StartPostgresSpan(ctx, "ipIntelligenceRepository.SetEmailAddress")
	defer span.Finish()

	now := time.Now()
	return r.write.WithContext(ctx).
		Model(&models.IPIntelligence{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"email_address": emailAddress,
			"updated_at":    now,
		}).Error
}
