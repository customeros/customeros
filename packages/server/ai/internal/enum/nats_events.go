package enum

type NatsEvents string

const (
	EventIdentifyVisitor  NatsEvents = "webtracker.visitor.identify"
	EventAIRequestGeneric NatsEvents = "ai.request.generic"
)

func (e NatsEvents) String() string {
	return string(e)
}
