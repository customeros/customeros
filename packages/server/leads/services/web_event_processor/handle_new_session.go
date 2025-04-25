package web_event_processor

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

const REQUEST_TIMEOUT = 60 * time.Second

type LeadSource struct {
	Channel        enum.Channel
	SearchEngine   enum.SearchEngine
	SocialPlatform enum.SocialPlatform
	AdPlatform     enum.AdPlatform
	ReferrerHost   string
	ReferrerPath   string
	Campaign       *CampaignMetadata
}

type CampaignMetadata struct {
	IsPaid      bool
	AdPlatform  enum.AdPlatform
	UTMSource   string
	UTMMedium   string
	UTMCampiagn string
	UTMTerm     string
	UTMContent  string
}

func (s *webEventProcessor) processNewSession(ctx context.Context, message *pb.WebTrackerEvent) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.processNewSession")
	defer spans.Finish()

	// attempt to identify company
	sessionIdentity, err := s.identifyWebSession(ctx, message.Ip)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	if sessionIdentity == nil {
		return nil
	}

	// determime lead source
	_, err = s.determineLeadSource(ctx, message.Href, message.Referrer)

	// lookup company to determine if new or existing lead

	// publish new or exisitng lead event
	// lead.new.webtracker
	// lead.existing.webtracker

	return nil
}

func (s *webEventProcessor) determineLeadSource(ctx context.Context, href, referrer string) (LeadSource, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.determineLeadSource")
	defer spans.Finish()

	leadSource := LeadSource{}

	if referrer == "" {
		leadSource.Channel = enum.ChannelDirect
		return leadSource, nil
	}

	parsedReferrer, err := utils.ParseURL(referrer)
	if err != nil {
		spans.TraceError(err)
		return leadSource, nil
	}

	isSearch, searchEngine := isSearchEngine(parsedReferrer.Host)
	isSocial, socialPlatform := isSocialPlatform(parsedReferrer)
	isEmail := isEmail(parsedReferrer.Host)

	parsedHref, err := utils.ParseURL(href)
	if err != nil {
		spans.TraceError(err)
		return leadSource, nil
	}

	campaignMetadata := parseCampaignMetadata(parsedHref)

	isPaid := campaignMetadata != nil && campaignMetadata.IsPaid

	switch {
	case isSearch && isPaid:
		leadSource.AdPlatform = campaignMetadata.AdPlatform
		leadSource.Channel = enum.ChannelPaidSearch

	case isSearch && !isPaid:
		leadSource.SearchEngine = searchEngine
		leadSource.Channel = enum.ChannelOrganicSearch

	case isSocial && isPaid:
		leadSource.AdPlatform = campaignMetadata.AdPlatform
		leadSource.Channel = enum.ChannelPaidSocial

	case isSocial && !isPaid:
		leadSource.SocialPlatform = socialPlatform
		leadSource.Channel = enum.ChannelOrganicSocial

	case isEmail && !isPaid:
		leadSource.Channel = enum.ChannelEmail

	case isPaid:
		leadSource.Channel = enum.ChannelPaidSocial
		leadSource.AdPlatform = campaignMetadata.AdPlatform

	default:
		leadSource.ReferrerHost = parsedReferrer.Host
		leadSource.ReferrerPath = parsedReferrer.Path
		leadSource.Channel = enum.ChannelReferral
	}

	leadSource.Campaign = campaignMetadata

	return leadSource, nil
}

func parseCampaignMetadata(href *utils.URLComponents) *CampaignMetadata {
	metadata := &CampaignMetadata{}

	if href == nil || len(href.QueryParams) == 0 {
		return nil
	}

	utmSource, ok := href.QueryParams["utm_source"]
	if ok && len(utmSource) > 0 {
		metadata.UTMSource = utmSource[0]
	}

	utmMedium, ok := href.QueryParams["utm_medium"]
	if ok && len(utmMedium) > 0 {
		metadata.UTMMedium = utmMedium[0]
	}

	utmCampaign, ok := href.QueryParams["utm_campaign"]
	if ok && len(utmCampaign) > 0 {
		metadata.UTMCampiagn = utmCampaign[0]
	}

	utmTerm, ok := href.QueryParams["utm_term"]
	if ok && len(utmTerm) > 0 {
		metadata.UTMTerm = utmTerm[0]
	}

	utmContent, ok := href.QueryParams["utm_content"]
	if ok && len(utmContent) > 0 {
		metadata.UTMContent = utmContent[0]
	}

	gclid, ok := href.QueryParams["gclid"]
	if ok && len(gclid) > 0 {
		metadata.AdPlatform = enum.AdsGoogle
		metadata.IsPaid = true
	}

	wbraid, ok := href.QueryParams["wbraid"]
	if ok && len(wbraid) > 0 {
		metadata.AdPlatform = enum.AdsGoogle
		metadata.IsPaid = true
	}

	fbclid, ok := href.QueryParams["fbclid"]
	if ok && len(fbclid) > 0 {
		metadata.AdPlatform = enum.AdsFacebook
		metadata.IsPaid = true
	}

	msclkid, ok := href.QueryParams["msclkid"]
	if ok && len(msclkid) > 0 {
		metadata.AdPlatform = enum.AdsMicrosoft
		metadata.IsPaid = true
	}

	lifatid, ok := href.QueryParams["li_fat_id"]
	if ok && len(lifatid) > 0 {
		metadata.AdPlatform = enum.AdsLinkedIn
		metadata.IsPaid = true
	}

	twclid, ok := href.QueryParams["twclid"]
	if ok && len(twclid) > 0 {
		metadata.AdPlatform = enum.AdsX
		metadata.IsPaid = true
	}

	pinclk, ok := href.QueryParams["pinclk"]
	if ok && len(pinclk) > 0 {
		metadata.AdPlatform = enum.AdsPinterest
		metadata.IsPaid = true
	}

	snap, ok := href.QueryParams["sc_referrer"]
	if ok && len(snap) > 0 {
		metadata.AdPlatform = enum.AdsSnapchat
		metadata.IsPaid = true
	}

	amazonAds, ok := href.QueryParams["aaxitk"]
	if ok && len(amazonAds) > 0 {
		metadata.AdPlatform = enum.AdsAmazon
		metadata.IsPaid = true
	}

	// Check for Apple Search Ads identifiers
	appleSearchAds, ok := href.QueryParams["ct"]
	if ok && len(appleSearchAds) > 0 {
		campaigns, ok := href.QueryParams["cp"]
		if ok && len(campaigns) > 0 {
			if strings.HasPrefix(campaigns[0], "app_store") {
				metadata.AdPlatform = enum.AdsApple
				metadata.IsPaid = true
			}
		}
	}

	return metadata
}

func (s *webEventProcessor) identifyWebSession(ctx context.Context, ipAddress string) (*pb.IdentifyVisitorResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.identifyWebSession")
	defer spans.Finish()

	if ipAddress == "" {
		return nil, nil
	}

	request := &pb.IdentifyVisitorRequest{}

	resp, err := s.sendIdentifyVisitorRequest(ctx, request)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if resp == nil {
		err = errors.New("Unable to identify web session")
		spans.TraceError(err)
		return nil, err
	}

	return resp, nil
}

func (s *webEventProcessor) sendIdentifyVisitorRequest(ctx context.Context, request *pb.IdentifyVisitorRequest) (*pb.IdentifyVisitorResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.sendIdentifyVisitorRequest")
	defer spans.Finish()

	// Marshal request to protobuf
	reqData, err := proto.Marshal(request)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Send request to service
	msg := nats.NewMsg(enum.EventAskSnitcher.String())
	msg.Header = nats.Header{
		enum.TENANT_HEADER:  []string{utils.GetTenantFromContext(ctx)},
		enum.USER_ID_HEADER: []string{utils.GetUserIdFromContext(ctx)},
	}
	msg.Data = reqData

	resp, err := s.natsConn.Conn.RequestMsg(msg, REQUEST_TIMEOUT)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Unmarshal response
	response := &pb.IdentifyVisitorResponse{}
	if err := proto.Unmarshal(resp.Data, response); err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return response, nil
}

func isSearchEngine(referrer string) (bool, enum.SearchEngine) {
	lowerRef := strings.ToLower(referrer)

	switch {
	case strings.Contains(lowerRef, "google."):
		return true, enum.SearchEngineGoogle
	case strings.Contains(lowerRef, "googlequicksearchbox"):
		return true, enum.SearchEngineGoogle
	case strings.Contains(lowerRef, "bing."):
		return true, enum.SearchEngineBing
	case strings.Contains(lowerRef, "yahoo."):
		return true, enum.SearchEngineYahoo
	case strings.Contains(lowerRef, "brave."):
		return true, enum.SearchEngineBrave
	case strings.Contains(lowerRef, "yandex."):
		return true, enum.SearchEngineYandex
	case strings.Contains(lowerRef, "baidu."):
		return true, enum.SearchEngineBaidu
	case strings.Contains(lowerRef, "duckduckgo."):
		return true, enum.SearchEngineDuckDuckGo
	case strings.Contains(lowerRef, "ecosia."):
		return true, enum.SearchEngineEcosia
	case strings.Contains(lowerRef, "ask."):
		return true, enum.SearchEngineAsk
	case strings.Contains(lowerRef, "qwant."):
		return true, enum.SearchEngineQwant
	case strings.Contains(lowerRef, "aol."):
		return true, enum.SearchEngineAOL
	case strings.Contains(lowerRef, "naver."):
		return true, enum.SearchEngineNaver
	case strings.Contains(lowerRef, "search."):
		return true, enum.SearchEngineGeneric
	case strings.Contains(lowerRef, "seznam."):
		return true, enum.SearchEngineSeznam
	case strings.Contains(lowerRef, "sogou."):
		return true, enum.SearchEngineSogou
	case strings.Contains(lowerRef, "startpage."):
		return true, enum.SearchEngineStartpage
	case strings.Contains(lowerRef, "coccoc."):
		return true, enum.SearchEngineCocCoc
	case strings.Contains(lowerRef, "searchencrypt."):
		return true, enum.SearchEngineSearchEncrypt
	case strings.Contains(lowerRef, "metager."):
		return true, enum.SearchEngineMetaGer
	case strings.Contains(lowerRef, "swisscows."):
		return true, enum.SearchEngineSwissCows
	case strings.Contains(lowerRef, "mojeek."):
		return true, enum.SearchEngineMojeek
	case strings.Contains(lowerRef, "gibiru."):
		return true, enum.SearchEngineGibiru
	case strings.Contains(lowerRef, "neeva."):
		return true, enum.SearchEngineNeeva
	case strings.Contains(lowerRef, "perplexity."):
		return true, enum.SearchEnginePerplexity
	case strings.Contains(lowerRef, "kagi."):
		return true, enum.SearchEngineKagi
	case strings.Contains(lowerRef, "spotlight."):
		return true, enum.SearchEngineAppleSpotlight
	case strings.Contains(lowerRef, "search.myway."):
		return true, enum.SearchEngineMyWay
	}

	return false, ""
}

func isSocialPlatform(referrer *utils.URLComponents) (bool, enum.SocialPlatform) {
	lowerRef := strings.ToLower(referrer.URL)

	// Extract hostname to avoid partial matches
	hostname := ""
	if referrer.Host != "" {
		hostname = referrer.Host
	} else {
		// If unable to parse, use the full string but with caution
		hostname = lowerRef
	}

	// Check for app referrers first
	if strings.HasPrefix(lowerRef, "android-app://") {
		appPackage := strings.TrimPrefix(lowerRef, "android-app://")
		if idx := strings.Index(appPackage, "/"); idx > 0 {
			appPackage = appPackage[:idx]
		}

		// Match app packages
		switch {
		case strings.Contains(appPackage, "facebook") || strings.Contains(appPackage, "fb"):
			return true, enum.SocialPlatformFacebook
		case strings.Contains(appPackage, "instagram"):
			return true, enum.SocialPlatformInstagram
		case strings.Contains(appPackage, "twitter") || strings.Contains(appPackage, ".x."):
			return true, enum.SocialPlatformX
		case strings.Contains(appPackage, "linkedin"):
			return true, enum.SocialPlatformLinkedIn
		case strings.Contains(appPackage, "pinterest"):
			return true, enum.SocialPlatformPinterest
		case strings.Contains(appPackage, "snapchat"):
			return true, enum.SocialPlatformSnapchat
		case strings.Contains(appPackage, "tiktok") || strings.Contains(appPackage, "musical.ly"):
			return true, enum.SocialPlatformTikTok
		case strings.Contains(appPackage, "whatsapp"):
			return true, enum.SocialPlatformWhatsApp
		case strings.Contains(appPackage, "telegram"):
			return true, enum.SocialPlatformTelegram
		case strings.Contains(appPackage, "youtube") || strings.Contains(appPackage, "ytshorts"):
			return true, enum.SocialPlatformYoutube
		case strings.Contains(appPackage, "line.android"):
			return true, enum.SocialPlatformLine
		case strings.Contains(appPackage, "wechat") || strings.Contains(appPackage, "micromessenger"):
			return true, enum.SocialPlatformWeChat
		case strings.Contains(appPackage, "kakao"):
			return true, enum.SocialPlatformKakao
		case strings.Contains(appPackage, "weibo"):
			return true, enum.SocialPlatformWeibo
		case strings.Contains(appPackage, "reddit"):
			return true, enum.SocialPlatformReddit
		case strings.Contains(appPackage, "quora"):
			return true, enum.SocialPlatformQuora
		case strings.Contains(appPackage, "discord"):
			return true, enum.SocialPlatformDiscord
		case strings.Contains(appPackage, "medium"):
			return true, enum.SocialPlatformMedium
		case strings.Contains(appPackage, "github"):
			return true, enum.SocialPlatformGithub
		case strings.Contains(appPackage, "gitlab"):
			return true, enum.SocialPlatformGitlab
		case strings.Contains(appPackage, "slack"):
			return true, enum.SocialPlatformSlack
		case strings.Contains(appPackage, "teams"):
			return true, enum.SocialPlatformTeams

		}
	}

	// Check for iOS app schemes
	if strings.HasPrefix(lowerRef, "fb://") {
		return true, enum.SocialPlatformFacebook
	} else if strings.HasPrefix(lowerRef, "twitter://") || strings.HasPrefix(lowerRef, "x://") {
		return true, enum.SocialPlatformX
	} else if strings.HasPrefix(lowerRef, "instagram://") {
		return true, enum.SocialPlatformInstagram
	} else if strings.HasPrefix(lowerRef, "pinterest://") {
		return true, enum.SocialPlatformPinterest
	} else if strings.HasPrefix(lowerRef, "snapchat://") {
		return true, enum.SocialPlatformSnapchat
	} else if strings.HasPrefix(lowerRef, "tiktok://") {
		return true, enum.SocialPlatformTikTok
	} else if strings.HasPrefix(lowerRef, "whatsapp://") {
		return true, enum.SocialPlatformWhatsApp
	} else if strings.HasPrefix(lowerRef, "tg://") {
		return true, enum.SocialPlatformTelegram
	} else if strings.HasPrefix(lowerRef, "youtube://") {
		return true, enum.SocialPlatformYoutube
	} else if strings.HasPrefix(lowerRef, "line://") {
		return true, enum.SocialPlatformLine
	} else if strings.HasPrefix(lowerRef, "wechat://") || strings.HasPrefix(lowerRef, "weixin://") {
		return true, enum.SocialPlatformWeChat
	} else if strings.HasPrefix(lowerRef, "kakao://") {
		return true, enum.SocialPlatformKakao
	} else if strings.HasPrefix(lowerRef, "reddit://") {
		return true, enum.SocialPlatformReddit
	} else if strings.HasPrefix(lowerRef, "quora://") {
		return true, enum.SocialPlatformQuora
	} else if strings.HasPrefix(lowerRef, "discord://") {
		return true, enum.SocialPlatformDiscord
	} else if strings.HasPrefix(lowerRef, "medium://") {
		return true, enum.SocialPlatformMedium
	} else if strings.HasPrefix(lowerRef, "github://") {
		return true, enum.SocialPlatformGithub
	} else if strings.HasPrefix(lowerRef, "gitlab://") {
		return true, enum.SocialPlatformGitlab
	}

	// Check web referrers
	switch {
	case strings.Contains(hostname, "facebook.") || hostname == "fb.com" || strings.Contains(hostname, "fbcdn."):
		return true, enum.SocialPlatformFacebook
	case strings.Contains(hostname, "instagram.") || hostname == "ig.me":
		return true, enum.SocialPlatformInstagram
	case hostname == "x.com" || strings.Contains(hostname, "twitter.") || hostname == "t.co":
		return true, enum.SocialPlatformX
	case strings.Contains(hostname, "linkedin.") || hostname == "lnkd.in":
		return true, enum.SocialPlatformLinkedIn
	case strings.Contains(hostname, "pinterest.") || hostname == "pin.it":
		return true, enum.SocialPlatformPinterest
	case strings.Contains(hostname, "snapchat."):
		return true, enum.SocialPlatformSnapchat
	case strings.Contains(hostname, "tiktok.") || strings.Contains(hostname, "musical.ly"):
		return true, enum.SocialPlatformTikTok
	case strings.Contains(hostname, "whatsapp."):
		return true, enum.SocialPlatformWhatsApp
	case strings.Contains(hostname, "telegram.") || strings.Contains(hostname, "t.me"):
		return true, enum.SocialPlatformTelegram
	case strings.Contains(hostname, "youtube.") || hostname == "youtu.be":
		return true, enum.SocialPlatformYoutube
	case strings.Contains(hostname, "line.me"):
		return true, enum.SocialPlatformLine
	case strings.Contains(hostname, "wechat.") || strings.Contains(hostname, "weixin."):
		return true, enum.SocialPlatformWeChat
	case strings.Contains(hostname, "kakao."):
		return true, enum.SocialPlatformKakao
	case strings.Contains(hostname, "weibo."):
		return true, enum.SocialPlatformWeibo
	case strings.Contains(hostname, "reddit.") || hostname == "redd.it":
		return true, enum.SocialPlatformReddit
	case hostname == "threads.net":
		return true, enum.SocialPlatformThreads
	case hostname == "slack.":
		return true, enum.SocialPlatformSlack
	case hostname == "teams.":
		return true, enum.SocialPlatformTeams
	case strings.Contains(hostname, "quora.") || hostname == "qr.ae":
		return true, enum.SocialPlatformQuora
	case strings.Contains(hostname, "discord.") || strings.Contains(hostname, "discordapp."):
		return true, enum.SocialPlatformDiscord
	case strings.Contains(hostname, "medium.") || hostname == "med.io":
		return true, enum.SocialPlatformMedium
	case strings.Contains(hostname, "github.") || hostname == "github.dev" || hostname == "gist.github.com":
		return true, enum.SocialPlatformGithub
	case strings.Contains(hostname, "gitlab.") || strings.Contains(hostname, "gitlab-static."):
		return true, enum.SocialPlatformGitlab
	}

	return false, ""
}

func isEmail(referrer string) bool {
	lowerRef := strings.ToLower(referrer)

	// Android app email clients
	if strings.HasPrefix(lowerRef, "android-app://") {
		switch {
		case strings.Contains(lowerRef, "com.google.android.gm"):
			return true
		case strings.Contains(lowerRef, "com.microsoft.office.outlook"):
			return true
		case strings.Contains(lowerRef, "com.yahoo.mobile.client.android.mail"):
			return true
		case strings.Contains(lowerRef, "com.samsung.android.email"):
			return true
		case strings.Contains(lowerRef, "com.mobincube.android.sc_bbv5m"):
			return true
		case strings.Contains(lowerRef, "ch.protonmail.android"):
			return true
		case strings.Contains(lowerRef, "com.readdle.spark"):
			return true
		case strings.Contains(lowerRef, "com.boxer.email"):
			return true
		case strings.Contains(lowerRef, "com.my.mail"):
			return true
		case strings.Contains(lowerRef, ".mail.") || strings.Contains(lowerRef, ".email."):
			return true
		}
	}

	// iOS mail app URI schemes
	if strings.HasPrefix(lowerRef, "message:") ||
		strings.HasPrefix(lowerRef, "mailto:") ||
		strings.HasPrefix(lowerRef, "ms-outlook:") {
		return true
	}

	// Web-based email clients
	hostname := ""
	if u, err := url.Parse(lowerRef); err == nil {
		hostname = u.Hostname()
	} else {
		hostname = lowerRef
	}

	switch {
	case strings.Contains(hostname, "mail.google.com"):
		return true
	case strings.Contains(hostname, "outlook.live.com") ||
		strings.Contains(hostname, "outlook.office365.com") ||
		strings.Contains(hostname, "outlook.office.com"):
		return true
	case strings.Contains(hostname, "mail.yahoo.com"):
		return true
	case strings.Contains(hostname, "protonmail.com") ||
		strings.Contains(hostname, "mail.proton.me"):
		return true
	case strings.Contains(hostname, "mail.aol.com"):
		return true
	case strings.Contains(hostname, "zoho.com/mail"):
		return true
	case strings.Contains(hostname, "mail.") ||
		strings.Contains(hostname, "webmail."):
		return true

	// Email marketing platforms
	case strings.Contains(hostname, "mailchimp.com") ||
		strings.Contains(hostname, "list-manage.com"):
		return true
	case strings.Contains(hostname, "sendgrid.net"):
		return true
	case strings.Contains(hostname, "constantcontact.com"):
		return true
	case strings.Contains(hostname, "mailerlite.com"):
		return true
	case strings.Contains(hostname, "campaign-archive.com"):
		return true
	case strings.Contains(hostname, "hubspotemail.net"):
		return true
	case strings.Contains(hostname, "convertkit-mail.com"):
		return true
	case strings.Contains(hostname, "klaviyo.com"):
		return true
	case strings.Contains(hostname, "aweber.com"):
		return true
	case strings.Contains(hostname, "getresponse.com"):
		return true
	case strings.Contains(hostname, "sendinblue.com") ||
		strings.Contains(hostname, "brevo.com"):
		return true
	}

	// Check for common email tracking parameters
	if strings.Contains(lowerRef, "utm_medium=email") ||
		strings.Contains(lowerRef, "utm_source=newsletter") ||
		strings.Contains(lowerRef, "mc_eid=") ||
		strings.Contains(lowerRef, "emlid=") ||
		strings.Contains(lowerRef, "email_id=") {
		return true
	}

	return false
}
