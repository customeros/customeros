package enum

type InteractionEventChannel string

const (
	InteractionEventChannelEmail InteractionEventChannel = "EMAIL"
	InteractionEventChannelChat  InteractionEventChannel = "CHAT"
	InteractionEventChannelVoice InteractionEventChannel = "VOICE"
)

func DecodeInteractionEventChannel(s string) InteractionEventChannel {
	switch InteractionEventChannel(s) {
	case InteractionEventChannelEmail, InteractionEventChannelChat, InteractionEventChannelVoice:
		return InteractionEventChannel(s)
	}
	return ""
}

func (i InteractionEventChannel) String() string {
	return string(i)
}
