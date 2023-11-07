package model

import "strings"

type UserData struct {
	BaseData
	Name            string        `json:"name,omitempty"`
	FirstName       string        `json:"firstName,omitempty"`
	LastName        string        `json:"lastName,omitempty"`
	Email           string        `json:"email,omitempty"`
	PhoneNumbers    []PhoneNumber `json:"phoneNumbers,omitempty"`
	ExternalOwnerId string        `json:"externalOwnerId,omitempty"` // Deprecated in favor or ExternalIdSecond, to be removed after release of sync-process is modified and released
	ProfilePhotoUrl string        `json:"profilePhotoUrl,omitempty"`
	Timezone        string        `json:"timezone,omitempty"`
	Bot             bool          `json:"bot,omitempty"`
}

func (u *UserData) HasPhoneNumbers() bool {
	return len(u.PhoneNumbers) > 0
}

func (u *UserData) HasEmail() bool {
	return u.Email != ""
}

func (u *UserData) Normalize() {
	u.SetTimes()
	u.BaseData.Normalize()

	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	for _, phoneNumber := range u.PhoneNumbers {
		phoneNumber.Normalize()
	}
	u.PhoneNumbers = GetNonEmptyPhoneNumbers(u.PhoneNumbers)
	u.PhoneNumbers = RemoveDuplicatedPhoneNumbers(u.PhoneNumbers)
}
