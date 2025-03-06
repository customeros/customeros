package enum

import (
	"fmt"
	"strings"
)

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

// WebpageCategoryValidator implements the EnumValidator interface for WebpageCategory
type WebpageCategoryValidator struct{}

func GetWebpageCategoryValidator() *WebpageCategoryValidator {
	return &WebpageCategoryValidator{}
}

func (v *WebpageCategoryValidator) IsValid(value string) bool {
	cleanValue := strings.TrimSpace(strings.ToLower(value))

	switch WebpageCategory(cleanValue) {
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

		return true
	default:
		return false
	}
}

func (v *WebpageCategoryValidator) ValidValues() []string {
	return []string{
		string(WebpageAbout),
		string(WebpageAccount),
		string(WebpageContact),
		string(WebpageHelp),
		string(WebpageLegal),
		string(WebpagePartner),
		string(WebpagePricing),
		string(WebpageProduct),
		string(WebpageResources),
		string(WebpageSuccessStory),
		string(WebpageOther),
	}
}

func (v *WebpageCategoryValidator) ParseCategory(value string) (WebpageCategory, error) {
	cleanValue := strings.TrimSpace(strings.ToLower(value))
	category := WebpageCategory(cleanValue)

	if !v.IsValid(cleanValue) {
		return WebpageUnknown, fmt.Errorf("invalid webpage category: %s", value)
	}

	return category, nil
}
