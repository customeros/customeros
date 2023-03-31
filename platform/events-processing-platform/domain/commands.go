package domain

import (
	contactService "github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/contact/service"
	phoneNumberService "github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number/service"
)

type Commands struct {
	ContactCommandService     *contactService.ContactCommandsService
	PhoneNumberCommandService *phoneNumberService.PhoneNumberCommandsService
}
