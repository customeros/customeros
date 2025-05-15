package enums

// NatsStream represents a NATS stream name
type NatsStream string

const (
	StreamTenant       NatsStream = "tenant"
	StreamOrganization NatsStream = "organization"
	StreamContact      NatsStream = "contact"
	StreamWeb          NatsStream = "web"
	StreamRequest      NatsStream = "request"
	StreamICP          NatsStream = "icp"
	StreamWebtracker   NatsStream = "webtracker"
	StreamProxy        NatsStream = "proxy"
	StreamLead         NatsStream = "lead"
	StreamDLQ          NatsStream = "dlq"
)

// String returns the string representation of the NatsStream
func (s NatsStream) String() string {
	return string(s)
}

// GetAllStreams returns all available NATS streams
func GetAllStreams() []NatsStream {
	return []NatsStream{
		StreamTenant,
		StreamOrganization,
		StreamContact,
		StreamWeb,
		StreamRequest,
		StreamICP,
		StreamWebtracker,
		StreamProxy,
		StreamLead,
		StreamDLQ,
	}
}
