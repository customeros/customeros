package enum

type APIVendor string

const (
	VendorAnthropic     APIVendor = "Anthropic"
	VendorBetterContact APIVendor = "BetterContact"
	VendorCrustData     APIVendor = "Crust Data"
	VendorDeepseek      APIVendor = "Deepseek"
	VendorEnrow         APIVendor = "Enrow"
	VendorGemini        APIVendor = "Gemini"
	VendorGoogle        APIVendor = "Google"
	VendorGroq          APIVendor = "Groq"
	VendorIPData        APIVendor = "IP Data"
	VendorJina          APIVendor = "Jina"
	VendorMicrosoft     APIVendor = "Microsoft"
	VendorQuickbooks    APIVendor = "Quickbooks"
	VendorScrapin       APIVendor = "Scrapin"
	VendorScrubbyIo     APIVendor = "ScrubbyIo"
	VendorSlack         APIVendor = "Slack"
	VendorSnitcher      APIVendor = "Snitcher"

	VendorNotSet APIVendor = "Not Set"
)
