package mapper

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

const (
	channelVoice    = "VOICE"
	channelMail     = "MAIL"
	channelChat     = "CHAT"
	channelWhatsApp = "WHATSAPP"
	channelTwitter  = "TWITTER"
	channelFacebook = "FACEBOOK"
)

var ChannelsByModel = map[model.MessageChannel]string{
	model.MessageChannelVoice:    channelVoice,
	model.MessageChannelMail:     channelMail,
	model.MessageChannelChat:     channelChat,
	model.MessageChannelWhatsapp: channelWhatsApp,
	model.MessageChannelFacebook: channelFacebook,
	model.MessageChannelTwitter:  channelTwitter,
}

var ChannelsByValue = map[string]model.MessageChannel{
	channelVoice:    model.MessageChannelVoice,
	channelMail:     model.MessageChannelMail,
	channelChat:     model.MessageChannelChat,
	channelWhatsApp: model.MessageChannelWhatsapp,
	channelFacebook: model.MessageChannelFacebook,
	channelTwitter:  model.MessageChannelTwitter,
}

func MapMessageChannelFromModel(input model.MessageChannel) string {
	return ChannelsByModel[input]
}

func MapMessageChannelToModel(input string) model.MessageChannel {
	return ChannelsByValue[input]
}
