package commands

type ContactCommands struct {
	CreateContact CreateContactCommandHandler
	UpdateContact UpdateContactCommandHandler
	DeleteContact DeleteContactCommandHandler
}

func NewContactCommands(
	createContact CreateContactCommandHandler,
	updateContact UpdateContactCommandHandler,
	deleteContact DeleteContactCommandHandler,
) *ContactCommands {
	return &ContactCommands{
		CreateContact: createContact,
		UpdateContact: updateContact,
		DeleteContact: deleteContact,
	}
}
