package enum

type NatsEvent string

const (
	EventTenantCreated NatsEvent = "core.tenant.created"
)

func (e NatsEvent) String() string {
	return string(e)
}
