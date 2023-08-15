package entity

type UserData struct {
	BaseData
	Name            string `json:"name,omitempty"`
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	Email           string `json:"email,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	ExternalOwnerId string `json:"externalOwnerId,omitempty"`
	ProfilePhotoUrl string `json:"profilePhotoUrl,omitempty"`
}

func (u *UserData) HasPhoneNumber() bool {
	return len(u.PhoneNumber) > 0
}

func (u *UserData) HasEmail() bool {
	return len(u.Email) > 0
}

func (u *UserData) Normalize() {
	u.SetTimes()
}
