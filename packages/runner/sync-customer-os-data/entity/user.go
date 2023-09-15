package entity

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"

type UserData struct {
	BaseData
	Name            string   `json:"name,omitempty"`
	FirstName       string   `json:"firstName,omitempty"`
	LastName        string   `json:"lastName,omitempty"`
	Email           string   `json:"email,omitempty"`
	PhoneNumbers    []string `json:"phoneNumbers,omitempty"`
	ExternalOwnerId string   `json:"externalOwnerId,omitempty"`
	ProfilePhotoUrl string   `json:"profilePhotoUrl,omitempty"`
	Timezone        string   `json:"timezone,omitempty"`
}

func (u *UserData) HasPhoneNumbers() bool {
	return len(u.PhoneNumbers) > 0
}

func (u *UserData) HasEmail() bool {
	return len(u.Email) > 0
}

func (u *UserData) Normalize() {
	u.SetTimes()

	u.PhoneNumbers = utils.FilterEmpty(u.PhoneNumbers)
	u.PhoneNumbers = utils.RemoveDuplicates(u.PhoneNumbers)
}
