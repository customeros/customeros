package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/repository/helper"
	"gorm.io/gorm"
)

type FileRepository interface {
	FindById(tenantName string, id string) helper.QueryResult
	Save(file *entity.File) helper.QueryResult
}

type fileRepo struct {
	db *gorm.DB
}

func NewFileRepo(db *gorm.DB) FileRepository {
	return &fileRepo{
		db: db,
	}
}

func (r *fileRepo) FindById(tenantName string, id string) helper.QueryResult {
	var file entity.File

	err := r.db.
		Where("id = ?", id).
		Where("tenant_name = ?", tenantName).
		First(&file).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &file}
}

func (r *fileRepo) Save(file *entity.File) helper.QueryResult {

	result := r.db.Create(&file)

	if result.Error != nil {
		return helper.QueryResult{Error: result.Error}
	}

	return helper.QueryResult{Result: &file}
}
