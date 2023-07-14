package entity

import (
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type UserData struct {
	BaseData
	Name            string `json:"name,omitempty"`
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	Email           string `json:"email,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	ExternalOwnerId string `json:"externalOwnerId,omitempty"`
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
	} else {
		u.CreatedAt = common_utils.TimePtr(common_utils.Now())
	}
	if u.UpdatedAt != nil {
		u.UpdatedAt = common_utils.TimePtr((*u.UpdatedAt).UTC())
	} else {
		u.UpdatedAt = common_utils.TimePtr(common_utils.Now())
	}
}

func (u *UserData) Normalize() {
	u.FormatTimes()
}
