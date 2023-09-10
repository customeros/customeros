package entity

type ReferencedEntity interface {
	GetReferencedEntityType() ReferencedEntityType
	Available() bool
}

type ReferencedEntityType string

const (
	ReferencedEntityTypeUnknown      ReferencedEntityType = "unknown"
	ReferencedEntityTypeSession      ReferencedEntityType = "session"
	ReferencedEntityTypeIssue        ReferencedEntityType = "issue"
	ReferencedEntityTypeJobRole      ReferencedEntityType = "jobRole"
	ReferencedEntityTypeUser         ReferencedEntityType = "user"
	ReferencedEntityTypeContact      ReferencedEntityType = "contact"
	ReferencedEntityTypeOrganization ReferencedEntityType = "organization"
)
