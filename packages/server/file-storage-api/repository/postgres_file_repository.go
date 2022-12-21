package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/repository/helper"
	"gorm.io/gorm"
)

type FileRepo struct {
	db *gorm.DB
}

type FileRepository interface {
	FindById(ctx context.Context, id string, tenantId string) helper.QueryResult
	Save(ctx context.Context, file entity.File) helper.QueryResult
}

func NewFileRepo(db *gorm.DB) *FileRepo {
	return &FileRepo{db: db}
}

func (r *FileRepo) FindById(ctx context.Context, id string, tenantId string) helper.QueryResult {
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

func (r *FileRepo) Save(ctx context.Context, file entity.File) helper.QueryResult {

	result := r.db.Create(&file)

	if result.Error != nil {
		return helper.QueryResult{Error: result.Error}
	}

	return helper.QueryResult{Result: &file}
}
