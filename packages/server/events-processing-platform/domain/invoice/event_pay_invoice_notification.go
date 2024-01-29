package invoice

//type InvoicePayNotificationEvent struct {
//	Tenant string `json:"tenant" validate:"required"`
//}
//
//func NewInvoicePayNotificationEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
//	eventData := InvoicePaidEvent{
//		Tenant: aggregate.GetTenant(),
//	}
//
//	if err := validator.GetValidator().Struct(eventData); err != nil {
//		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoicePayNotificationEvent")
//	}
//
//	event := eventstore.NewBaseEvent(aggregate, asdfaskdjhfaskdjf)
//	if err := event.SetJsonData(&eventData); err != nil {
//		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoicePayNotificationEvent")
//	}
//
//	return event, nil
//}
