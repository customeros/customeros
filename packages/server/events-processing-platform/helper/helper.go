package helper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"strings"
)

func GetSourceOfTruth(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.SourceOpenline)
}

func GetSource(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.SourceOpenline)
}

func GetAppSource(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.AppSourceEventProcessingPlatform)
}
