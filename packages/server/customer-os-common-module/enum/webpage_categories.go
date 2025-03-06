package enum

type WebpageCategory string

const (
	WebpageAbout        WebpageCategory = "about"
	WebpageAccount      WebpageCategory = "account"
	WebpageContact      WebpageCategory = "contact"
	WebpageHelp         WebpageCategory = "help"
	WebpageLegal        WebpageCategory = "legal"
	WebpagePartner      WebpageCategory = "partner"
	WebpagePricing      WebpageCategory = "pricing"
	WebpageProduct      WebpageCategory = "product"
	WebpageResources    WebpageCategory = "resources"
	WebpageSuccessStory WebpageCategory = "success story"
	WebpageOther        WebpageCategory = "other"
	WebpageUnknown      WebpageCategory = ""
)

func (s WebpageCategory) String() string {
	return string(s)
}

func GetWebpageCategory(s string) WebpageCategory {
	switch WebpageCategory(s) {
	case
		WebpageAbout,
		WebpageAccount,
		WebpageContact,
		WebpageHelp,
		WebpageLegal,
		WebpagePartner,
		WebpagePricing,
		WebpageProduct,
		WebpageResources,
		WebpageSuccessStory,
		WebpageOther:
		return WebpageCategory(s)
	default:
		return WebpageUnknown
	}
}
