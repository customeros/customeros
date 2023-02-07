package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntitiesToDashboardViewItems(entities []*entity.DashboardViewResultEntity) []*model.DashboardViewItem {
	var items []*model.DashboardViewItem
	for _, entity := range entities {
		item := new(model.DashboardViewItem)
		if entity.Organization != nil {
			item.Organization = MapEntityToOrganization(entity.Organization)
		}
		if entity.Contact != nil {
			item.Contact = MapEntityToContact(entity.Contact)
		}
		items = append(items, item)
	}
	return items
}
