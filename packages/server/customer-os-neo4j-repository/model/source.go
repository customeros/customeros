package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"strings"
)

type Source struct {
	Source        string `json:"source"`
	SourceOfTruth string `json:"sourceOfTruth"`
	AppSource     string `json:"appSource"`
}

func GetSourceOfTruth(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.SourceOpenline)
}

func GetSource(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.SourceOpenline)
}

func GetAppSource(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.AppSourceCustomerOsApi)
}
