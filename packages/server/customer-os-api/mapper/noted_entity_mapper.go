package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"reflect"
)

func MapEntityToNotedEntity(notedEntity *entity.NotedEntity) any {
	switch (*notedEntity).NotedEntityLabel() {
	case entity.NodeLabel_Organization:
		return MapEntityToOrganization((*notedEntity).(*entity.OrganizationEntity))
	case entity.NodeLabel_Contact:
		return MapEntityToContact((*notedEntity).(*entity.ContactEntity))
	}
	fmt.Errorf("noted entity of type %s not identified", reflect.TypeOf(notedEntity))
	return nil
}

func MapEntitiesToNotedEntities(entities *entity.NotedEntities) []model.NotedEntity {
	var notedEntities []model.NotedEntity
	for _, entity := range *entities {
		notedEntity := MapEntityToNotedEntity(&entity)
		if notedEntity != nil {
			notedEntities = append(notedEntities, notedEntity.(model.NotedEntity))
		}
	}
	return notedEntities
}
