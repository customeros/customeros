package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/repository/entity"
	"reflect"
	"testing"
	"time"
)

var (
	applicationEntity = entity.ApplicationEntity{
		ID:          "someID",
		AppId:       "someAppId",
		Platform:    "somePlatform",
		TrackerName: "someTrackerName",
		Tenant:      "theTenant",
		UpdatedOn:   time.Now(),
	}
	expectedApplication = model.Application{
		ID:          applicationEntity.ID,
		Platform:    applicationEntity.Platform,
		Name:        applicationEntity.AppId,
		TrackerName: applicationEntity.TrackerName,
		Tenant:      applicationEntity.Tenant,
	}
)

func TestMapApplicationEmpty(t *testing.T) {
	emptyExpectedApplication := model.Application{}
	newApplicationEntity := entity.ApplicationEntity{}
	application := MapApplication(&newApplicationEntity)
	if !reflect.DeepEqual(application, &emptyExpectedApplication) {
		t.Errorf("Expecting: %+v, received: %+v", emptyExpectedApplication, application)
	}
}

func TestMapApplication(t *testing.T) {
	application := MapApplication(&applicationEntity)
	if !reflect.DeepEqual(application, &expectedApplication) {
		t.Errorf("Expecting: %+v, received: %+v", expectedApplication, application)
	}
}

func TestMapApplicationsEmpty(t *testing.T) {
	applications := MapApplications(&entity.ApplicationEntities{})
	if len(applications) != 0 {
		t.Errorf("Expecting empty applications, received: %+v", applications)
	}
}

func TestMapApplications(t *testing.T) {
	applicationEntities := entity.ApplicationEntities{}
	applicationEntities = append(applicationEntities, applicationEntity)
	applicationEntities = append(applicationEntities, applicationEntity)
	applications := MapApplications(&applicationEntities)
	if len(applications) != 2 {
		t.Errorf("Expecting applications with 2 elements, received: %d elements", len(applications))
	}
	for i := range []int{1, 2} {
		if !reflect.DeepEqual(applications[i], &expectedApplication) {
			t.Errorf("Expecting: %+v, received: %+v", expectedApplication, applications[i])
		}
	}
}
