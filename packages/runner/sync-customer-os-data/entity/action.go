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
