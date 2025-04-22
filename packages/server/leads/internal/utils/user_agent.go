package utils

import (
	"regexp"
	"strings"
)

type UserAgent struct {
	Raw            string
	Browser        string
	BrowserVersion string
	Engine         string
	EngineVersion  string
	OS             string
	OSVersion      string
	Device         string
	DeviceModel    string
	IsMobile       bool
	IsTablet       bool
	IsBot          bool
	IsWebView      bool
	SocialPlatform string
}

// ParseUserAgent parses a user agent string into structured information
func ParseUserAgent(userAgentString string) UserAgent {
	ua := UserAgent{
		Raw: userAgentString,
	}

	splitUA := strings.Split(userAgentString, "|")
	if len(splitUA) == 1 {
		userAgentString = splitUA[0]
	} else {
		userAgentString = splitUA[1]
	}

	// Check for WebView
	if strings.Contains(userAgentString, "wv") || strings.Contains(userAgentString, "Dalvik") {
		ua.IsWebView = true
	}

	// Check for social platform browsers
	socialPlatforms := map[string]string{
		"Facebook":  "FBAN|FB_IAB|FBAV|FBDV",
		"Instagram": "Instagram",
		"X":         "Twitter|TwitterAndroid|",
		"LinkedIn":  "LinkedIn",
		"Pinterest": "Pinterest",
		"Snapchat":  "Snapchat",
		"TikTok":    "TikTok|musical_ly|musically|ByteDance|Bytedance",
		"WhatsApp":  "WhatsApp",
		"Telegram":  "TelegramBot|Telegram",
		"Line":      "Line",
		"WeChat":    "MicroMessenger",
		"Kakao":     "KAKAOTALK",
		"Weibo":     "Weibo",
	}

	for platform, pattern := range socialPlatforms {
		if regexContains(userAgentString, pattern) {
			ua.SocialPlatform = platform
			break
		}
	}

	// Check for bots
	botPatterns := []string{
		"Googlebot", "bingbot", "Baiduspider", "YandexBot", "facebookexternalhit",
		"Slackbot", "Twitterbot", "bot", "spider", "crawler", "Mediapartners-Google", "Google",
	}
	for _, pattern := range botPatterns {
		if strings.Contains(strings.ToLower(userAgentString), strings.ToLower(pattern)) {
			ua.IsBot = true
			ua.Browser = "Bot"
			if strings.Contains(userAgentString, pattern) {
				ua.BrowserVersion = extractVersion(userAgentString, pattern)
			}
			break
		}
	}

	// Parse browser information
	browserPatterns := map[string]string{
		"Edge":    "Edge|Edg",
		"Chrome":  "Chrome",
		"Safari":  "Safari",
		"Firefox": "Firefox",
		"Opera":   "Opera|OPR",
		"IE":      "MSIE|Trident",
	}

	for browser, pattern := range browserPatterns {
		if regexContains(userAgentString, pattern) && !ua.IsBot {
			ua.Browser = browser
			ua.BrowserVersion = extractVersion(userAgentString, pattern)
			break
		}
	}

	// Parse engine information
	enginePatterns := map[string]string{
		"Gecko":   "Gecko",
		"WebKit":  "WebKit",
		"Trident": "Trident",
		"Blink":   "Chrome",
		"Presto":  "Presto",
	}

	for engine, pattern := range enginePatterns {
		if regexContains(userAgentString, pattern) {
			ua.Engine = engine
			ua.EngineVersion = extractVersion(userAgentString, pattern)
			break
		}
	}

	// Parse OS information
	osPatterns := map[string]string{
		"Windows":  "Windows NT",
		"macOS":    "Mac OS X",
		"iOS":      "iPhone|iPad|iPod",
		"Android":  "Android",
		"Linux":    "Linux",
		"Unix":     "X11",
		"ChromeOS": "CrOS",
	}

	for os, pattern := range osPatterns {
		if regexContains(userAgentString, pattern) {
			ua.OS = os
			ua.OSVersion = extractOSVersion(userAgentString, os)
			break
		}
	}

	// Parse device type
	mobilePatterns := []string{
		"Mobile", "Android", "iPhone", "iPod", "BlackBerry",
		"Windows Phone", "IEMobile",
	}

	tabletPatterns := []string{
		"iPad", "Android", "Tablet", "Silk",
	}

	for _, pattern := range mobilePatterns {
		if regexContains(userAgentString, pattern) {
			ua.IsMobile = true

			// Special case for Android
			if ua.OS == "Android" && !regexContains(userAgentString, "Mobile") {
				// Android without "Mobile" typically means tablet
				ua.IsMobile = false
				ua.IsTablet = true
			}
			break
		}
	}

	if !ua.IsMobile {
		for _, pattern := range tabletPatterns {
			if regexContains(userAgentString, pattern) {
				ua.IsTablet = true
				break
			}
		}
	}

	// Extract device model
	extractDeviceModel(userAgentString, &ua)

	// Determine device
	if ua.IsBot {
		ua.Device = "Bot"
	} else if ua.IsMobile {
		if strings.Contains(userAgentString, "iPhone") {
			ua.Device = "iPhone"
		} else if strings.Contains(userAgentString, "iPad") {
			ua.Device = "iPad"
			ua.IsTablet = true
			ua.IsMobile = false
		} else if strings.Contains(userAgentString, "Android") {
			ua.Device = "Android Mobile"
		} else {
			ua.Device = "Mobile"
		}
	} else if ua.IsTablet {
		if strings.Contains(userAgentString, "iPad") {
			ua.Device = "iPad"
		} else if strings.Contains(userAgentString, "Android") {
			ua.Device = "Android Tablet"
		} else {
			ua.Device = "Tablet"
		}
	} else {
		ua.Device = "Desktop"
	}

	// Override browser info if this is a social platform WebView
	if ua.SocialPlatform != "" {
		if ua.Browser != "" {
			ua.Browser = ua.SocialPlatform + " (" + ua.Browser + ")"
		} else {
			ua.Browser = ua.SocialPlatform + " WebView"
		}
	}

	return ua
}

// Helper functions
func regexContains(s, pattern string) bool {
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

func extractVersion(s, pattern string) string {
	// The logic for extracting version depends on the pattern
	if pattern == "Windows NT" {
		re := regexp.MustCompile(`Windows NT (\d+\.\d+)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
	} else if pattern == "Mac OS X" {
		re := regexp.MustCompile(`Mac OS X (\d+[._]\d+[._]\d+)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return strings.Replace(matches[1], "_", ".", -1)
		}
	} else if pattern == "Chrome" {
		re := regexp.MustCompile(`Chrome\/(\d+\.\d+[\.\d]*)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
	} else if pattern == "Firefox" {
		re := regexp.MustCompile(`Firefox\/(\d+\.\d+)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
	} else if pattern == "Safari" {
		re := regexp.MustCompile(`Version\/(\d+\.\d+[\.\d]*)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
	} else if pattern == "Edge|Edg" {
		re := regexp.MustCompile(`Edge?\/(\d+\.\d+[\.\d]*)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
	} else if pattern == "MSIE|Trident" {
		re := regexp.MustCompile(`MSIE (\d+\.\d+)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
		// For IE 11
		re = regexp.MustCompile(`rv:(\d+\.\d+)`)
		matches = re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
	} else if pattern == "Android" {
		re := regexp.MustCompile(`Android (\d+[\.\d]*)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// Generic version extraction as fallback
	re := regexp.MustCompile(pattern + `[\/\s](\d+[\.\d]*)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func extractOSVersion(s, os string) string {
	switch os {
	case "Windows":
		return extractVersion(s, "Windows NT")
	case "macOS":
		return extractVersion(s, "Mac OS X")
	case "iOS":
		// Try to get iOS version
		re := regexp.MustCompile(`OS (\d+[._]\d+[._]?\d*)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return strings.Replace(matches[1], "_", ".", -1)
		}
	case "Android":
		return extractVersion(s, "Android")
	case "ChromeOS":
		re := regexp.MustCompile(`CrOS [^\s]+ (\d+[\.\d]*)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

func extractDeviceModel(userAgentString string, ua *UserAgent) {
	// Extract specific device models
	// Android devices usually have their model in the user agent
	if strings.Contains(userAgentString, "Android") {
		// Try to get the model name that's usually between "Android" and "Build"
		re := regexp.MustCompile(`Android [^;]+; ([^;)]+)(?:[ ;]Build\/|[;)])`)
		matches := re.FindStringSubmatch(userAgentString)
		if len(matches) > 1 {
			// Clean up the model name
			model := strings.TrimSpace(matches[1])
			// Remove some common manufacturer prefixes if model already has manufacturer info
			ua.DeviceModel = model
		}
	} else if strings.Contains(userAgentString, "iPhone") {
		ua.DeviceModel = "iPhone"
	} else if strings.Contains(userAgentString, "iPad") {
		ua.DeviceModel = "iPad"
	} else if strings.Contains(userAgentString, "Macintosh") {
		ua.DeviceModel = "Macintosh"
	}

	// Extract additional device info from FBDV (Facebook Device)
	if strings.Contains(userAgentString, "FBDV/") {
		re := regexp.MustCompile(`FBDV/([^;]+)`)
		matches := re.FindStringSubmatch(userAgentString)
		if len(matches) > 1 {
			ua.DeviceModel = matches[1]
		}
	}
}
