package entity

type RawState string

const (
	PENDING RawState = "PENDING"
	SENT    RawState = "SENT"
	SKIPPED RawState = "SKIPPED"
	ERROR   RawState = "ERROR"
)

func (rawState RawState) String() string {
	return string(rawState)
}
