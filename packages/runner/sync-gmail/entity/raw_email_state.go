package entity

type RawEmailState string

const (
	PENDING RawEmailState = "PENDING"
	SENT    RawEmailState = "SENT"
	SKIPPED RawEmailState = "SKIPPED"
	ERROR   RawEmailState = "ERROR"
)

func (rawEmailState RawEmailState) String() string {
	return string(rawEmailState)
}
