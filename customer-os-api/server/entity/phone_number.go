package entity

import (
	"fmt"
)

type PhoneNumberEntity struct {
	Number string
	Label  string
}

func (phone PhoneNumberEntity) ToString() string {
	return fmt.Sprintf("number: %s\nlabel: %s", phone.Number, phone.Label)
}

type PhoneNumberEntities []PhoneNumberEntity
