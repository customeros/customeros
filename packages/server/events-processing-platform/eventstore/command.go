package eventstore

type Command interface {
	GetObjectID() string
	GetTenant() string
}

type BaseCommand struct {
	ObjectID string `json:"objectID" validate:"required"`
	Tenant   string `json:"tenant"`
}

func NewBaseCommand(objectID, tenant string) BaseCommand {
	return BaseCommand{
		ObjectID: objectID,
		Tenant:   tenant,
	}
}

func (c *BaseCommand) GetObjectID() string {
	return c.ObjectID
}

func (c *BaseCommand) GetTenant() string {
	return c.Tenant
}
