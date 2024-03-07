package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/repository/entity"
	"reflect"
	"testing"
	"time"
)

var (
	sessionEntity = entity.SessionEntity{
		ID:             "someID",
		AppId:          "someAppId",
		Country:        "someCountry",
		Region:         "someRegion",
		City:           "someCity",
		ReferrerSource: "someReferrerSource",
		UtmCampaign:    "someUtmCampaign",
		UtmContent:     "someUtmContent",
		UtmMedium:      "someUtmMedium",
		UtmSource:      "someUtmSource",
		UtmNetwork:     "someUtmNetwork",
		UtmTerm:        "someUtmTerm",
		DeviceName:     "someDeviceName",
		DeviceBrand:    "someDeviceBrand",
		DeviceClass:    "someDeviceClass",
		AgentName:      "someAgentName",
		AgentVersion:   "someAgentVersion",
		OsFamily:       "someOperatingSystem",
		OsVersionMajor: "someOsVersionMajor",
		OsVersionMinor: "someOsVersionMinor",
		FirstPagePath:  "someFirstPagePath",
		LastPagePath:   "someLastPagePath",
		Start:          time.Now(),
		End:            time.Now(),
		EngagedTime:    123456,
	}
	expectedSession = model.AppSession{
		ID:              sessionEntity.ID,
		Country:         sessionEntity.Country,
		Region:          sessionEntity.Region,
		City:            sessionEntity.City,
		ReferrerSource:  sessionEntity.ReferrerSource,
		UtmCampaign:     sessionEntity.UtmCampaign,
		UtmContent:      sessionEntity.UtmContent,
		UtmMedium:       sessionEntity.UtmMedium,
		UtmSource:       sessionEntity.UtmSource,
		UtmTerm:         sessionEntity.UtmTerm,
		UtmNetwork:      sessionEntity.UtmNetwork,
		DeviceBrand:     sessionEntity.DeviceBrand,
		DeviceName:      sessionEntity.DeviceName,
		DeviceClass:     sessionEntity.DeviceClass,
		AgentName:       sessionEntity.AgentName,
		AgentVersion:    sessionEntity.AgentVersion,
		OperatingSystem: sessionEntity.OsFamily,
		OsVersionMajor:  sessionEntity.OsVersionMajor,
		OsVersionMinor:  sessionEntity.OsVersionMinor,
		FirstPagePath:   sessionEntity.FirstPagePath,
		LastPagePath:    sessionEntity.LastPagePath,
		StartedAt:       sessionEntity.Start,
		EndedAt:         sessionEntity.End,
		EngagedTime:     sessionEntity.EngagedTime,
	}
)

func TestMapSessionEmpty(t *testing.T) {
	emptyExpectedSession := model.AppSession{}
	newSessionEntity := entity.SessionEntity{}
	session := MapSession(&newSessionEntity)
	if !reflect.DeepEqual(session, &emptyExpectedSession) {
		t.Errorf("Expecting: %+v, received: %+v", emptyExpectedSession, session)
	}
}

func TestMapSession(t *testing.T) {
	session := MapSession(&sessionEntity)
	if !reflect.DeepEqual(session, &expectedSession) {
		t.Errorf("Expecting: %+v, received: %+v", expectedSession, session)
	}
}

func TestMapSessionsEmpty(t *testing.T) {
	sessions := MapSessions(&entity.SessionEntities{})
	if len(sessions) != 0 {
		t.Errorf("Expecting empty sessions, received: %+v", sessions)
	}
}

func TestMapSessions(t *testing.T) {
	sessionEntities := entity.SessionEntities{}
	sessionEntities = append(sessionEntities, sessionEntity)
	sessionEntities = append(sessionEntities, sessionEntity)
	sessions := MapSessions(&sessionEntities)
	if len(sessions) != 2 {
		t.Errorf("Expecting sessions with 2 elements, received: %d elements", len(sessions))
	}
	for i := range []int{1, 2} {
		if !reflect.DeepEqual(sessions[i], &expectedSession) {
			t.Errorf("Expecting: %+v, received: %+v", expectedSession, sessions[i])
		}
	}
}
