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
)

func (s WebpageCategory) String() string {
	return string(s)
}
