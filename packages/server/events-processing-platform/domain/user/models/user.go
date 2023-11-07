package models

import (
	"fmt"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type User struct {
	ID              string                       `json:"id"`
	Name            string                       `json:"name"`
	FirstName       string                       `json:"firstName"`
	LastName        string                       `json:"lastName"`
	Internal        bool                         `json:"internal"`
	Bot             bool                         `json:"bot"`
	ProfilePhotoUrl string                       `json:"profilePhotoUrl"`
	Timezone        string                       `json:"timezone"`
	CreatedAt       time.Time                    `json:"createdAt"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
	PhoneNumbers    map[string]UserPhoneNumber   `json:"phoneNumbers"`
	Emails          map[string]UserEmail         `json:"emails"`
	JobRoles        map[string]bool              `json:"jobRoles"`
	Source          commonmodel.Source           `json:"source"`
	ExternalSystems []commonmodel.ExternalSystem `json:"externalSystems"`
	Players         []PlayerInfo                 `json:"players"`
	Roles           []string                     `json:"roles"`
}

type PlayerInfo struct {
	Provider   string `json:"provider"`
	AuthId     string `json:"authId"`
	IdentityId string `json:"identityId"`
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

func (u *User) HasEmail(emailId, label string, primary bool) bool {
	if len(u.Emails) == 0 {
		return false
	}
	if email, ok := u.Emails[emailId]; ok {
		return email.Label == label && email.Primary == primary
	}
	return false
}

func (u *User) SameData(fields UserDataFields, externalSystem commonmodel.ExternalSystem) bool {
	if externalSystem.Available() && !u.HasExternalSystem(externalSystem) {
		return false
	}
	if u.Name == fields.Name &&
		u.FirstName == fields.FirstName &&
		u.LastName == fields.LastName &&
		u.Internal == fields.Internal &&
		u.Bot == fields.Bot &&
		u.Timezone == fields.Timezone &&
		u.ProfilePhotoUrl == fields.ProfilePhotoUrl {
		return true
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
