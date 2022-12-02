package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"reflect"
)

func MapEntityToAction(actionEntity *entity.Action) any {
	switch (*actionEntity).ActionName() {
	case entity.ActionName_PageView:
		pageViewActionEntityPtr := (*actionEntity).(*entity.PageViewActionEntity)
		return MapEntityToPageViewAction(pageViewActionEntityPtr)
	}
	fmt.Errorf("action of type %s not identified", reflect.TypeOf(actionEntity))
	return nil
}

func MapEntitiesToActions(entities *entity.ActionEntities) []model.Action {
	var actions []model.Action
	for _, actionEntity := range *entities {
		action := MapEntityToAction(&actionEntity)
		if action != nil {
			actions = append(actions, action.(model.Action))
		}
	}
	return actions
}
