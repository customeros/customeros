package service

type ServiceOptions struct {
	SkipCompletedEvents bool
}

func PublishCompletedEvents(options ...ServiceOptions) bool {
	if options == nil {
		return true
	}
	if len(options) == 0 {
		return true
	}
	return !options[0].SkipCompletedEvents
}
