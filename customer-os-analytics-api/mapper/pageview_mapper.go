package mapper

import (
	"github.com.openline-ai.customer-os-analytics-api/graph/model"
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
)

func MapPageViews(pageViewEntities *entity.PageViewEntities) []*model.PageView {
	var pageViews []*model.PageView
	for _, pageViewEntity := range *pageViewEntities {
		pageViews = append(pageViews, MapPageView(&pageViewEntity))
	}
	return pageViews
}

func MapPageView(pageViewEntity *entity.PageViewEntity) *model.PageView {
	return &model.PageView{
		ID:          pageViewEntity.ID,
		Title:       pageViewEntity.Title,
		Path:        pageViewEntity.Path,
		Order:       pageViewEntity.OrderInSession,
		EngagedTime: pageViewEntity.EngagedTime,
		SessionId:   pageViewEntity.SessionID,
	}
}
