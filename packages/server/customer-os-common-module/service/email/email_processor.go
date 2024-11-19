package service

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func GetRawEmailData(rawEmail string) (EmailRawData, error) {
	rawEmailData := EmailRawData{}
	err := json.Unmarshal([]byte(rawEmail), &rawEmailData)
	if err != nil {
		logrus.Errorf("Unmarshal Raw Email Data Failed: %v", err)
		return rawEmailData, err
	}

	return rawEmailData, nil
}
