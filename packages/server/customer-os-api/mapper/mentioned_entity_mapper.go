package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"reflect"
)

func MapEntityToMentionedEntity(mentionedEntity *entity.MentionedEntity) any {
	switch (*mentionedEntity).MentionedEntityLabel() {
	case entity.NodeLabel_Issue:
		return MapEntityToIssue((*mentionedEntity).(*entity.IssueEntity))
	}
	fmt.Errorf("mentioned entity of type %s not identified", reflect.TypeOf(mentionedEntity))
	return nil
}

func MapEntitiesToMentionedEntities(entities *entity.MentionedEntities) []model.MentionedEntity {
	var mentionedEntities []model.MentionedEntity
	for _, entity := range *entities {
		mentionedEntity := MapEntityToMentionedEntity(&entity)
		if mentionedEntity != nil {
			mentionedEntities = append(mentionedEntities, mentionedEntity.(model.MentionedEntity))
		}
	}
	return mentionedEntities
}
