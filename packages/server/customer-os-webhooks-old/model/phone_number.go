package model

import "strings"

type PhoneNumber struct {
	Number  string `json:"number"`
	Primary bool   `json:"primary,omitempty"`
	Label   string `json:"label,omitempty"`
}

func RemoveDuplicatedPhoneNumbers(numbers []PhoneNumber) []PhoneNumber {
	var uniqueNumbers []PhoneNumber
	uniqueNumbersMap := make(map[string]bool)
	for _, number := range numbers {
		if _, ok := uniqueNumbersMap[number.Number]; !ok {
			uniqueNumbersMap[number.Number] = true
			uniqueNumbers = append(uniqueNumbers, number)
		}
	}
	return uniqueNumbers
}

func GetNonEmptyPhoneNumbers(phoneNumbers []PhoneNumber) []PhoneNumber {
	var nonEmptyPhoneNumbers []PhoneNumber
	for _, phoneNumber := range phoneNumbers {
		if phoneNumber.Number != "" {
			nonEmptyPhoneNumbers = append(nonEmptyPhoneNumbers, phoneNumber)
		}
	}
	return nonEmptyPhoneNumbers
}

func (p *PhoneNumber) Normalize() {
	p.Number = strings.TrimSpace(p.Number)
}
