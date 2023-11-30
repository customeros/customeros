package eventstore

type Command interface {
	GetObjectID() string
	GetTenant() string
}

type BaseCommand struct {
	ObjectID       string `json:"objectID" validate:"required"`
	Tenant         string `json:"tenant" validate:"required"`
	LoggedInUserId string `json:"loggedInUserId"`
	AppSource      string `json:"appSource"`
}

func (c BaseCommand) WithAppSource(appSource string) BaseCommand {
	c.AppSource = appSource
	return c
}

func NewBaseCommand(objectID, tenant, loggedInUserId string) BaseCommand {
	return BaseCommand{
		ObjectID:       objectID,
		Tenant:         tenant,
		LoggedInUserId: loggedInUserId,
	}
}

func (c BaseCommand) GetObjectID() string {
	return c.ObjectID
}

func (c BaseCommand) GetTenant() string {
	return c.Tenant
}

func (c BaseCommand) GetLoggedInUserId() string {
	return c.LoggedInUserId
}

func (c BaseCommand) GetAppSource() string {
	return c.AppSource
}
