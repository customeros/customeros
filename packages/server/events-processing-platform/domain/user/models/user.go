package models

import (
	"fmt"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"time"
)

type User struct {
	ID              string                       `json:"id"`
	Name            string                       `json:"name"`
	FirstName       string                       `json:"firstName"`
	LastName        string                       `json:"lastName"`
	Internal        bool                         `json:"internal"`
	Bot             bool                         `json:"bot"`
	Test            bool                         `json:"test"`
	ProfilePhotoUrl string                       `json:"profilePhotoUrl"`
	Timezone        string                       `json:"timezone"`
	CreatedAt       time.Time                    `json:"createdAt"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
	PhoneNumbers    map[string]UserPhoneNumber   `json:"phoneNumbers"`
	JobRoles        map[string]bool              `json:"jobRoles"`
	Source          commonmodel.Source           `json:"source"`
	ExternalSystems []commonmodel.ExternalSystem `json:"externalSystems"`
	Roles           []string                     `json:"roles"`
}

type UserPhoneNumber struct {
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

func (u *User) String() string {
	return fmt.Sprintf("User{ID: %s, Name: %s, FirstName: %s, LastName: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s, PhoneNumbers: %v}", u.ID, u.Name, u.FirstName, u.LastName, u.Source, u.CreatedAt, u.UpdatedAt, u.PhoneNumbers)
}

func (u *User) HasPhoneNumber(phoneNumberId, label string) bool {
	if len(u.PhoneNumbers) == 0 {
		return false
	}
	if phoneNumber, ok := u.PhoneNumbers[phoneNumberId]; ok {
		return phoneNumber.Label == label
	}
	return false
}
