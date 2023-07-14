package entity

import (
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

/*
{
  "name": "John Doe",
  "firstName": "John",
  "lastName": "Doe",
  "email": "john@email.com",
  "phoneNumber": "123-456-7890",
  "externalOwnerId": "user-123",

  "skip": false,
  "skipReason": "draft data",
  "id": "1234",
  "externalId": "abcd1234",
  "externalSystem": "HubSpot",
  "createdAt": "2022-02-28T19:52:05Z",
  "updatedAt": "2022-03-01T11:23:45Z",
  "syncId": "sync_1234"
}
*/

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
