package entity

import (
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type UserData struct {
	Id              string     `json:"id,omitempty"`
	Name            string     `json:"name,omitempty"`
	FirstName       string     `json:"firstName,omitempty"`
	LastName        string     `json:"lastName,omitempty"`
	Email           string     `json:"email,omitempty"`
	PhoneNumber     string     `json:"phoneNumber,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
	ExternalId      string     `json:"externalId,omitempty"`
	ExternalOwnerId string     `json:"externalOwnerId,omitempty"`

	ExternalSystem string `json:"externalSystem,omitempty"`
	ExternalSyncId string `json:"externalSyncId,omitempty"`
}

func (u *UserData) HasPhoneNumber() bool {
	return len(u.PhoneNumber) > 0
}

func (u *UserData) HasEmail() bool {
	return len(u.Email) > 0
}

func (u *UserData) FormatTimes() {
	if u.CreatedAt != nil {
		u.CreatedAt = common_utils.TimePtr((*u.CreatedAt).UTC())
	}
	if u.UpdatedAt != nil {
		u.UpdatedAt = common_utils.TimePtr((*u.UpdatedAt).UTC())
	}
}
