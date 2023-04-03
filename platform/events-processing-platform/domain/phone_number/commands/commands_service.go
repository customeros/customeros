package commands

type PhoneNumberCommands struct {
	CreatePhoneNumber CreatePhoneNumberCommandHandler
	//DeletePhoneNumber DeletePhoneNumberCommandHandler
}

func NewPhoneNumberCommands(createPhoneNumber CreatePhoneNumberCommandHandler) *PhoneNumberCommands {
	return &PhoneNumberCommands{
		CreatePhoneNumber: createPhoneNumber,
	}
}
