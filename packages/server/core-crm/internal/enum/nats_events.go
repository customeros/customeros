package enum

type NatsEvents string

const (
	EventAskIPData       NatsEvents = "ipaddress.verify.ipdata"
	EventAskSnitcher     NatsEvents = "ipaddress.identify.snitcher"
	EventIdentifyVisitor NatsEvents = "webtracker.visitor.identify"
)

func (e NatsEvents) String() string {
	return string(e)
}
