package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapStateToGCliItem(stateEntity neo4jEntity.StateEntity) model.GCliItem {
	resultItem := model.GCliItem{}

	resultItem.ID = stateEntity.Id
	resultItem.Type = model.GCliSearchResultTypeState
	resultItem.Display = stateEntity.Name
	data := []*model.GCliAttributeKeyValuePair{}
	data = append(data, &model.GCliAttributeKeyValuePair{
		Key:   "code",
		Value: stateEntity.Code,
	})
	resultItem.Data = data

	return resultItem
}

func MapContactToGCliItem(contactEntity entity.ContactEntity) model.GCliItem {
	resultItem := model.GCliItem{}

	resultItem.ID = contactEntity.Id
	resultItem.Type = model.GCliSearchResultTypeContact

	if contactEntity.FirstName != "" {
		resultItem.Display = contactEntity.FirstName + " " + contactEntity.LastName
	} else if contactEntity.Name != "" {
		resultItem.Display = contactEntity.Name
	}

	return resultItem
}
func MapOrganizationToGCliItem(contactEntity entity.OrganizationEntity) model.GCliItem {
	resultItem := model.GCliItem{}

	resultItem.ID = contactEntity.ID
	resultItem.Type = model.GCliSearchResultTypeOrganization
	resultItem.Display = contactEntity.Name

	return resultItem
}
func MapEmailToGCliItem(emailEntity entity.EmailEntity) model.GCliItem {
	resultItem := model.GCliItem{}

	resultItem.ID = emailEntity.Id
	resultItem.Type = model.GCliSearchResultTypeOrganization
	resultItem.Display = utils.StringFirstNonEmpty(emailEntity.Email, emailEntity.RawEmail)

	return resultItem
}
