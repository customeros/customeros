package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/repository/helper"
	"gorm.io/gorm"
)

type FileRepository interface {
	FindById(tenantId string, id string) helper.QueryResult
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

func (r *fileRepo) FindById(tenantId string, id string) helper.QueryResult {
	var file entity.File

	err := r.db.
		Where("id = ?", id).
		Where("tenant_id = ?", tenantId).
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
