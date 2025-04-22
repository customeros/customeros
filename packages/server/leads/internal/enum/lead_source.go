package enum

type SearchEngine string

const (
	SearchEngineAOL            SearchEngine = "aol"
	SearchEngineAppleSpotlight SearchEngine = "applespotlight"
	SearchEngineAsk            SearchEngine = "ask"
	SearchEngineBaidu          SearchEngine = "baidu"
	SearchEngineBing           SearchEngine = "bing"
	SearchEngineBrave          SearchEngine = "brave"
	SearchEngineCocCoc         SearchEngine = "coccoc"
	SearchEngineDuckDuckGo     SearchEngine = "duckduckgo"
	SearchEngineEcosia         SearchEngine = "ecosia"
	SearchEngineGibiru         SearchEngine = "gibiru"
	SearchEngineGoogle         SearchEngine = "google"
	SearchEngineKagi           SearchEngine = "kagi"
	SearchEngineMetaGer        SearchEngine = "metager"
	SearchEngineMojeek         SearchEngine = "mojeek"
	SearchEngineMyWay          SearchEngine = "myway"
	SearchEngineNaver          SearchEngine = "naver"
	SearchEngineNeeva          SearchEngine = "neeva"
	SearchEnginePerplexity     SearchEngine = "perplexity"
	SearchEngineQwant          SearchEngine = "qwant"
	SearchEngineSearchEncrypt  SearchEngine = "searchencrypt"
	SearchEngineSeznam         SearchEngine = "seznam"
	SearchEngineSogou          SearchEngine = "sogou"
	SearchEngineStartpage      SearchEngine = "startpage"
	SearchEngineSwissCows      SearchEngine = "swisscows"
	SearchEngineYahoo          SearchEngine = "yahoo"
	SearchEngineYandex         SearchEngine = "yandex"
	SearchEngineGeneric        SearchEngine = "generic"
)

type SocialPlatform string

const (
	SocialPlatformFacebook  SocialPlatform = "facebook"
	SocialPlatformInstagram SocialPlatform = "instagram"
	SocialPlatformX         SocialPlatform = "x"
	SocialPlatformLinkedIn  SocialPlatform = "linkedin"
	SocialPlatformPinterest SocialPlatform = "pinterest"
	SocialPlatformSnapchat  SocialPlatform = "snapchat"
	SocialPlatformTikTok    SocialPlatform = "tiktok"
	SocialPlatformWhatsApp  SocialPlatform = "whatsapp"
	SocialPlatformTelegram  SocialPlatform = "telegram"
	SocialPlatformYoutube   SocialPlatform = "youtube"
	SocialPlatformLine      SocialPlatform = "line"
	SocialPlatformWeChat    SocialPlatform = "wechat"
	SocialPlatformKakao     SocialPlatform = "kakao"
	SocialPlatformWeibo     SocialPlatform = "weibo"
	SocialPlatformReddit    SocialPlatform = "reddit"
	SocialPlatformThreads   SocialPlatform = "threads"
	SocialPlatformQuora     SocialPlatform = "quora"
	SocialPlatformDiscord   SocialPlatform = "discord"
	SocialPlatformMedium    SocialPlatform = "medium"
	SocialPlatformGithub    SocialPlatform = "github"
	SocialPlatformGitlab    SocialPlatform = "gitlab"
	SocialPlatformSlack     SocialPlatform = "slack"
	SocialPlatformTeams     SocialPlatform = "teams"
)

type AdPlatform string

const (
	AdsGoogle    AdPlatform = "google ads"
	AdsFacebook  AdPlatform = "facebook ads"
	AdsMicrosoft AdPlatform = "microsoft ads"
	AdsLinkedIn  AdPlatform = "linkedin ads"
	AdsX         AdPlatform = "x ads"
	AdsPinterest AdPlatform = "pinterest ads"
	AdsSnapchat  AdPlatform = "snapchat ads"
	AdsAmazon    AdPlatform = "amazon ads"
	AdsApple     AdPlatform = "apple ads"
)

type DeviceType string

const (
	DeviceMobile  DeviceType = "mobile"
	DeviceTablet  DeviceType = "tablet"
	DeviceDesktop DeviceType = "desktop"
	DeviceOther   DeviceType = "other"
)

type Channel string

const (
	ChannelPaidSearch      Channel = "paid search"
	ChannelPaidSocial      Channel = "paid social"
	ChannelOrganicSearch   Channel = "organic search"
	ChannelOrganicSocial   Channel = "organic social"
	ChannelReferral        Channel = "organic referral"
	ChannelPartnerReferral Channel = "partner referral"
	ChannelEmail           Channel = "email"
	ChannelDirect          Channel = "direct"
)
