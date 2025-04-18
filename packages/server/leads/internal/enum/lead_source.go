package enum

type LeadSource string

const (
	LeadSourceGoogleAds LeadSource = "google ads"
)

type SocialPlatform string

const (
	SocialPlatformGoogle SocialPlatform = "google"
)

type DeviceType string

const (
	DeviceMobile DeviceType = "mobile"
)

type Channel string

const (
	ChannelPaidSearch      Channel = "paid search"
	ChannelOrganicSearch   Channel = "organic search"
	ChannelSocial          Channel = "social"
	ChannelEmail           Channel = "email"
	ChannelDirect          Channel = "direct"
	ChannelPartnerReferral Channel = "partner referral"
)
