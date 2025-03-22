package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	enummapper "github.com/customeros/customeros/packages/server/customer-os-api/mapper/enum"
)

func MapEntityToPageView(entity *entity.PageViewEntity) *model.PageView {
	return &model.PageView{
		ID:             entity.Id,
		StartedAt:      entity.StartedAt,
		EndedAt:        entity.EndedAt,
		Application:    entity.Application,
		SessionID:      entity.SessionId,
		PageURL:        entity.PageUrl,
		PageTitle:      entity.PageTitle,
		OrderInSession: entity.OrderInSession,
		EngagedTime:    entity.EngagedTime,
		Source:         enummapper.MapDataSourceToModel(entity.Source),
		AppSource:      entity.AppSource,
	}
}
