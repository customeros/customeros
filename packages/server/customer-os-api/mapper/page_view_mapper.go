package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
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
	}
}
