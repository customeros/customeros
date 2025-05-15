package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/core-crm/internal/database"
	"github.com/customeros/customeros/packages/server/core-crm/internal/models"
)

type ICPDescriptionRepository interface {
	Create(icp *models.ICPDescription) error
	GetByTenant(tenant string) (*models.ICPDescription, error)
	UpdateProfile(id, tenant, profile string) error
	UpdateQualifyingAttributes(id, tenant, attributes string) error
	UpdateDisqualifyingAttributes(id, tenant, attributes string) error
}

type icpDescriptionRepo struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewICPDescriptionRepository(openlineDB *database.DatabaseConnection) ICPDescriptionRepository {
	return &icpDescriptionRepo{
		readDB:  openlineDB.ReadDB,
		writeDB: openlineDB.WriteDB,
	}
}

func (r *icpDescriptionRepo) Create(icp *models.ICPDescription) error {
	return r.writeDB.Create(icp).Error
}

func (r *icpDescriptionRepo) GetByTenant(tenant string) (*models.ICPDescription, error) {
	var icp models.ICPDescription
	if err := r.readDB.Where("tenant = ?", tenant).First(&icp).Error; err != nil {
		return nil, err
	}
	return &icp, nil
}

func (r *icpDescriptionRepo) UpdateProfile(id, tenant, profile string) error {
	result := r.writeDB.Model(&models.ICPDescription{}).
		Where("id = ? AND tenant = ?", id, tenant).
		Update("profile", profile)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("icp description not found")
	}

	return nil
}

func (r *icpDescriptionRepo) UpdateQualifyingAttributes(id, tenant, attributes string) error {
	result := r.writeDB.Model(&models.ICPDescription{}).
		Where("id = ? AND tenant = ?", id, tenant).
		Update("qualifying_attributes", attributes)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("icp description not found")
	}

	return nil
}

func (r *icpDescriptionRepo) UpdateDisqualifyingAttributes(id, tenant, attributes string) error {
	result := r.writeDB.Model(&models.ICPDescription{}).
		Where("id = ? AND tenant = ?", id, tenant).
		Update("disqualifying_attributes", attributes)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("icp description not found")
	}

	return nil
}
