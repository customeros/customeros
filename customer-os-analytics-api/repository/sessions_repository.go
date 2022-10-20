package repository

import (
	"github.com.openline-ai.customer-os-analytics-api/graph/model"
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
	"github.com.openline-ai.customer-os-analytics-api/repository/helper"
	"gorm.io/gorm"
)

type SessionsRepo struct {
	db *gorm.DB
}

type SessionsRepository interface {
	FindAllByApplication(appIdentifier entity.ApplicationUniqueIdentifier, dataFilter []*model.AppSessionsDataFilter, page int, limit int) helper.QueryResult
}

func NewSessionsRepo(db *gorm.DB) *SessionsRepo {
	return &SessionsRepo{db: db}
}

func (r *SessionsRepo) FindAllByApplication(appIdentifier entity.ApplicationUniqueIdentifier, dataFilter []*model.AppSessionsDataFilter,
	page int, limit int) helper.QueryResult {

	var sessions entity.SessionEntities

	pagination := helper.Pagination{
		Limit: limit,
		Page:  page,
	}

	find := r.db.
		Where(&entity.SessionEntity{Tenant: appIdentifier.Tenant, AppId: appIdentifier.AppId, TrackerName: appIdentifier.TrackerName})

	if dataFilter != nil {
		for _, value := range dataFilter {
			find = helper.AddDataFilter(columnNameForField(value.Field), value.Action, value.Value, find)
		}
	}

	err := find.Scopes(helper.Paginate(sessions, &pagination, find)).
		Order("start_tstamp DESC").
		Find(&sessions).
		Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	pagination.Rows = &sessions

	return helper.QueryResult{Result: &pagination}
}

func columnNameForField(field model.AppSessionField) string {
	switch field {
	case model.AppSessionFieldCountry:
		return entity.SessionColumnName_Country
	case model.AppSessionFieldCity:
		return entity.SessionColumnName_City
	case model.AppSessionFieldRegion:
		return entity.SessionColumnName_RegionName
	default:
		return ""
	}
}
