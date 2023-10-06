package entity

import "strings"

type UserData struct {
	BaseData
	Name            string        `json:"name,omitempty"`
	FirstName       string        `json:"firstName,omitempty"`
	LastName        string        `json:"lastName,omitempty"`
	Email           string        `json:"email,omitempty"`
	PhoneNumbers    []PhoneNumber `json:"phoneNumbers,omitempty"`
	ExternalOwnerId string        `json:"externalOwnerId,omitempty"` // Deprecated, remove once release latest webhooks and sync-customer-os-data
	ProfilePhotoUrl string        `json:"profilePhotoUrl,omitempty"`
	Timezone        string        `json:"timezone,omitempty"`
}

func (u *UserData) HasPhoneNumbers() bool {
	return len(u.PhoneNumbers) > 0
}

func (u *UserData) HasEmail() bool {
	return len(u.Email) > 0
}

func (u *UserData) Normalize() {
	u.SetTimes()

	u.Email = strings.ToLower(u.Email)

	u.PhoneNumbers = GetNonEmptyPhoneNumbers(u.PhoneNumbers)
	u.PhoneNumbers = RemoveDuplicatedPhoneNumbers(u.PhoneNumbers)
}
