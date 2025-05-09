package enum

type NatsEvents string

const (
	EventIdentifyVisitor NatsEvents = "webtracker.visitor.identify"
)

func (e NatsEvents) String() string {
	return string(e)
}
