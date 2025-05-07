package enum

type APIVendor string

const (
	VendorNotSet        APIVendor = "Not Set"
	VendorIPData        APIVendor = "IP Data"
	VendorSnitcher      APIVendor = "Snitcher"
	VendorScrapin       APIVendor = "Scrapin"
	VendorEnrow         APIVendor = "Enrow"
	VendorBetterContact APIVendor = "BetterContact"
	VendorScrubbyIo     APIVendor = "ScrubbyIo"
	VendorGoogle        APIVendor = "Google"
	VendorMicrosoft     APIVendor = "Microsoft"
	VendorSlack         APIVendor = "Slack"
)
