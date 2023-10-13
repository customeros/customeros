package entity

type UserData struct {
	BaseData
	Name            string        `json:"name,omitempty"`
	FirstName       string        `json:"firstName,omitempty"`
	LastName        string        `json:"lastName,omitempty"`
	Email           string        `json:"email,omitempty"`
	PhoneNumbers    []PhoneNumber `json:"phoneNumbers,omitempty"`
	ProfilePhotoUrl string        `json:"profilePhotoUrl,omitempty"`
	Timezone        string        `json:"timezone,omitempty"`
}
