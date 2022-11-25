package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/repository/helper"
	"gorm.io/gorm"
)

type PageViewRepo struct {
	db *gorm.DB
}

type PageViewRepository interface {
	GetAllBySessionIds(sessionIds []string) helper.QueryResult
}

func NewPageViewRepo(db *gorm.DB) *PageViewRepo {
	return &PageViewRepo{db: db}
}

func (r *PageViewRepo) GetAllBySessionIds(sessionIds []string) helper.QueryResult {
	var pageViews entity.PageViewEntities

	err := r.db.Where("domain_sessionid IN ?", sessionIds).
		Order("page_view_in_session_index ASC").
		Find(&pageViews).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &pageViews}
}
