package enum

type InteractionSessionChannel string

const (
	InteractionSessionChannelEmail InteractionSessionChannel = "EMAIL"
	InteractionSessionChannelChat  InteractionSessionChannel = "CHAT"
)

func DecodeInteractionSessionChannel(s string) InteractionSessionChannel {
	switch InteractionSessionChannel(s) {
	case InteractionSessionChannelEmail, InteractionSessionChannelChat:
		return InteractionSessionChannel(s)
	}
	return ""
}

func (i InteractionSessionChannel) String() string {
	return string(i)
}
