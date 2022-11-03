package entity

import (
	"fmt"
)

type PhoneNumberEntity struct {
	Id      string
	Number  string
	Label   string
	Primary bool
}

func (phone PhoneNumberEntity) ToString() string {
	return fmt.Sprintf("id: %s\nnumber: %s\nlabel: %s", phone.Id, phone.Number, phone.Label)
}

type PhoneNumberEntities []PhoneNumberEntity
