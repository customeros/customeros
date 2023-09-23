package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type RawCalendarEventRepository interface {
	GetByProviderId(externalSystem, tenantName, usernameSource, calendarId, providerId string) (*entity.RawCalendarEvent, error)
	SaveOrUpdate(externalSystem, tenantName, usernameSource, calendarId, providerId, iCalUID, rawCalendarEvent string) error
	Update(entity entity.RawCalendarEvent) error
}

type rawCalendarEventRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewRawCalendarEventRepository(gormDb *gorm.DB) RawCalendarEventRepository {
	return &rawCalendarEventRepositoryImpl{gormDb: gormDb}
}

func (repo *rawCalendarEventRepositoryImpl) GetByProviderId(externalSystem, tenantName, usernameSource, calendarId, providerId string) (*entity.RawCalendarEvent, error) {
	result := entity.RawCalendarEvent{}
	err := repo.gormDb.Find(&result, "external_system = ? AND tenant_name = ? AND username_source = ? AND calendar_id = ? AND provider_id = ?", externalSystem, tenantName, usernameSource, calendarId, providerId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, nil
}

func (repo *rawCalendarEventRepositoryImpl) SaveOrUpdate(externalSystem, tenantName, usernameSource, calendarId, providerId, iCalUID, rawCalendarEvent string) error {
	result := entity.RawCalendarEvent{}
	err := repo.gormDb.Find(&result, "external_system = ? AND tenant_name = ? AND username_source = ? AND calendar_id = ? AND provider_id = ?", externalSystem, tenantName, usernameSource, calendarId, providerId).Error

	if err != nil {
		logrus.Errorf("Failed retrieving rawCalendarEvent: %s; %s; %s; %s; %s", externalSystem, tenantName, usernameSource, calendarId, providerId)
		return err
	}

	if result.TenantName == "" {
		result.CreatedAt = time.Now().UTC()
		result.ExternalSystem = externalSystem
		result.TenantName = tenantName
		result.UsernameSource = usernameSource
		result.CalendarId = calendarId
		result.ProviderId = providerId
		result.ICalUID = iCalUID
	}

	result.UpdatedAt = time.Now().UTC()
	result.Data = rawCalendarEvent
	result.SentToEventStoreState = "PENDING"

	err = repo.gormDb.Save(&result).Error
	if err != nil {
		logrus.Errorf("Failed storing rawEmail: %s; %s; %s; %s", externalSystem, tenantName, usernameSource, iCalUID)
		return err
	}

	return nil
}

func (repo *rawCalendarEventRepositoryImpl) Update(entity entity.RawCalendarEvent) error {
	return repo.gormDb.Save(&entity).Error
}
