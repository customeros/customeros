package mapper

import (
	"github.com.openline-ai.customer-os-analytics-api/graph/model"
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
)

func MapApplications(applicationEntities *entity.ApplicationEntities) []*model.Application {
	var applications []*model.Application
	for _, applicationEntity := range *applicationEntities {
		applications = append(applications, MapApplication(&applicationEntity))
	}
	return applications
}

func MapApplication(applicationEntity *entity.ApplicationEntity) *model.Application {
	return &model.Application{
		ID:          applicationEntity.ID,
		Name:        applicationEntity.AppId,
		TrackerName: applicationEntity.TrackerName,
		Platform:    applicationEntity.Platform,
		Tenant:      applicationEntity.Tenant,
	}
}
