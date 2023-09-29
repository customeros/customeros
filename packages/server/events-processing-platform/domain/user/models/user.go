package models

import (
	"fmt"
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type User struct {
	ID              string                         `json:"id"`
	Name            string                         `json:"name"`
	FirstName       string                         `json:"firstName"`
	LastName        string                         `json:"lastName"`
	Internal        bool                           `json:"internal"`
	ProfilePhotoUrl string                         `json:"profilePhotoUrl"`
	Timezone        string                         `json:"timezone"`
	CreatedAt       time.Time                      `json:"createdAt"`
	UpdatedAt       time.Time                      `json:"updatedAt"`
	PhoneNumbers    map[string]UserPhoneNumber     `json:"phoneNumbers"`
	Emails          map[string]UserEmail           `json:"emails"`
	JobRoles        map[string]bool                `json:"jobRoles"`
	Source          common_models.Source           `json:"source"`
	ExternalSystems []common_models.ExternalSystem `json:"externalSystem"`
}

type UserPhoneNumber struct {
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

type UserEmail struct {
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

func (u *User) String() string {
	return fmt.Sprintf("User{ID: %s, Name: %s, FirstName: %s, LastName: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s, PhoneNumbers: %v, Emails: %v}", u.ID, u.Name, u.FirstName, u.LastName, u.Source, u.CreatedAt, u.UpdatedAt, u.PhoneNumbers, u.Emails)
}

func NewUser() *User {
	return &User{}
}
