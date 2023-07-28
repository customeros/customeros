package model

import "time"

type CallEventParty struct {
	Tel    *string `json:"tel,omitempty"`
	Stir   *string `json:"stir,omitempty"`
	Mailto *string `json:"mailto,omitempty"`
	Name   *string `json:"name,omitempty"`
}

type CallEvent struct {
	Version       string          `json:"version,default=1.0"`
	CorrelationId string          `json:"correlation_id"`
	Event         string          `json:"event"`
	From          *CallEventParty `json:"from"`
	To            *CallEventParty `json:"to"`
}

type CallEventStart struct {
	CallEvent
	StartTime time.Time `json:"start_time"`
}

type CallEventAnswered struct {
	CallEvent
	StartTime    time.Time `json:"start_time"`
	AnsweredTime time.Time `json:"answered_time"`
}

type CallEventEnd struct {
	CallEvent
	StartTime    *time.Time `json:"start_time,omitempty"`
	AnsweredTime *time.Time `json:"answered_time,omitempty"`
	EndTime      time.Time  `json:"end_time"`
	Duration     int64      `json:"duration"`
	FromCaller   bool       `json:"from_caller"`
}
