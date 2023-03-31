package commands

type PhoneNumberCommands struct {
	CreatePhoneNumber CreatePhoneNumberCommandHandler
}

func NewPhoneNumberCommands(createPhoneNumber CreatePhoneNumberCommandHandler) *PhoneNumberCommands {
	return &PhoneNumberCommands{
		CreatePhoneNumber: createPhoneNumber,
	}
}
