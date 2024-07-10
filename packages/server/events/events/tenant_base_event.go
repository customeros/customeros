package events

type TenantBaseEvent struct {
	BaseEvent

	Tenant string `json:"tenant" validate:"required"`
}
