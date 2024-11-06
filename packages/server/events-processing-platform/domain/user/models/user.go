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

func (u *User) SameUserData(fields UserDataFields, externalSystem commonmodel.ExternalSystem) bool {
	if !externalSystem.Available() {
		return false
	}

	if externalSystem.Available() && !u.HasExternalSystem(externalSystem) {
		return false
	}

	if u.Source.SourceOfTruth == externalSystem.ExternalSystemId {
		if u.Name == fields.Name &&
			u.FirstName == fields.FirstName &&
			u.LastName == fields.LastName &&
			u.Internal == fields.Internal &&
			u.Bot == fields.Bot &&
			u.Test == fields.Test &&
			u.Timezone == fields.Timezone &&
			u.ProfilePhotoUrl == fields.ProfilePhotoUrl {
			return true
		}
	} else {
		if (u.Name != "" || u.Name == fields.Name) &&
			(u.FirstName != "" || u.FirstName == fields.FirstName) &&
			(u.LastName != "" || u.LastName == fields.LastName) &&
			(u.Internal != false || u.Internal == fields.Internal) &&
			(u.Bot != false || u.Bot == fields.Bot) &&
			(u.Test != false || u.Test == fields.Test) &&
			(u.Timezone != "" || u.Timezone == fields.Timezone) &&
			(u.ProfilePhotoUrl != "" || u.ProfilePhotoUrl == fields.ProfilePhotoUrl) {
			return true
		}
	}
	return false
}

func (u *User) HasExternalSystem(externalSystem commonmodel.ExternalSystem) bool {
	for _, es := range u.ExternalSystems {
		if es.ExternalSystemId == externalSystem.ExternalSystemId &&
			es.ExternalId == externalSystem.ExternalId &&
			es.ExternalSource == externalSystem.ExternalSource &&
			es.ExternalUrl == externalSystem.ExternalUrl &&
			es.ExternalIdSecond == externalSystem.ExternalIdSecond {
			return true
		}
	}
	return false
}
