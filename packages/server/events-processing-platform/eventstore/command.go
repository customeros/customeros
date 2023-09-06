package eventstore

type Command interface {
	GetObjectID() string
	GetTenant() string
}

type BaseCommand struct {
	ObjectID string `json:"objectID" validate:"required"`
	Tenant   string `json:"tenant" validate:"required"`
	UserID   string `json:"userId"`
}

func NewBaseCommand(objectID, tenant, userID string) BaseCommand {
	return BaseCommand{
		ObjectID: objectID,
		Tenant:   tenant,
		UserID:   userID,
	}
}

func (c *BaseCommand) GetObjectID() string {
	return c.ObjectID
}

func (c *BaseCommand) GetTenant() string {
	return c.Tenant
}

func (c *BaseCommand) GetUserID() string {
	return c.UserID
}
