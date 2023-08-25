package entity

import (
	"time"
)

type ActionEntity struct {
	Id        string
	CreatedAt time.Time

	Type ActionType

	Source    string
	AppSource string

	DataloaderKey string
}

type ActionType string

const (
	ActionNA      ActionType = ""
	ActionCreated ActionType = "CREATED"
)

var AllActionType = []ActionType{
	ActionCreated,
}

func IsValidActionType(s string) bool {
	for _, ds := range AllActionType {
		if ds == ActionType(s) {
			return true
		}
	}
	return false
}

type ActionEntities []ActionEntity
