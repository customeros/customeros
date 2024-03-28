package offering

type OfferingCreateEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

//func NewReminderCreateEvent(aggregate eventstore.Aggregate, content, userId, organizationId string, dismissed bool, createdAt, dueDate time.Time, sourceFields cmnmod.Source) (eventstore.Event, error) {
//	eventData := ReminderCreateEvent{
//		Tenant:         aggregate.GetTenant(),
//		Content:        content,
//		DueDate:        dueDate,
//		UserId:         userId,
//		OrganizationId: organizationId,
//		Dismissed:      dismissed,
//		CreatedAt:      createdAt,
//		SourceFields:   sourceFields,
//	}
//	if err := validator.GetValidator().Struct(eventData); err != nil {
//		return eventstore.Event{}, errors.Wrap(err, "failed to validate ReminderCreateEvent")
//	}
//	event := eventstore.NewBaseEvent(aggregate, ReminderCreateV1)
//	if err := event.SetJsonData(&eventData); err != nil {
//		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ReminderCreateEvent")
//	}
//	return event, nil
//}
