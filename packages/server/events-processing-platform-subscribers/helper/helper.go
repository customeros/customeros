package helper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/constants"
	"strings"
)

func GetSourceOfTruth(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.SourceOpenline)
}

func GetSource(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.SourceOpenline)
}

func GetAppSource(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.AppSourceEventProcessingPlatformSubscribers)
}
