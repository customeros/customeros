package enum

type NatsEvents string

const (
	EventIdentifyVisitor              NatsEvents = "webtracker.visitor.identify"
	EventAIRequestGeneric             NatsEvents = "ai.request.generic"
	EventRequestWebpageClassification NatsEvents = "ai.request.webpage_classification"
	EventRequestWebpageIntent         NatsEvents = "ai.request.webpage_intent"
	EventWebpageClassified            NatsEvents = "webpage.classified"
)

func (e NatsEvents) String() string {
	return string(e)
}
