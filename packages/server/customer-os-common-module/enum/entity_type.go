package enum

import "fmt"

type EntityType string

const (
	EntityContact      EntityType = "contact"
	EntityMeeting      EntityType = "meeting"
	EntityOrganization EntityType = "organization"
)

func (t EntityType) String() string {
	return string(t)
}

func GetEntityType(s string) (EntityType, error) {
	switch EntityType(s) {
	case
		EntityContact,
		EntityMeeting,
		EntityOrganization:
		return EntityType(s), nil

	default:
		return "", fmt.Errorf("invalid EntityType: %s", s)
	}
}
