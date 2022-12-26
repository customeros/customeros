package entity

import (
	"fmt"
)

type AddressEntity struct {
	Id       string
	Source   string
	Country  string
	State    string
	City     string
	Address  string
	Address2 string
	Zip      string
	Phone    string
	Fax      string
}

func (address AddressEntity) ToString() string {
	return fmt.Sprintf("id: %s", address.Id)
}

type AddressEntities []AddressEntity

func (address AddressEntity) Labels() []string {
	return []string{"Address"}
}
