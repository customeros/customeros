package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"strings"
)

type SourceFields struct {
	Source    string `json:"source"`
	AppSource string `json:"appSource"`
	// Deprecated
	SourceOfTruth string `json:"sourceOfTruth"`
}

func (s SourceFields) GetSource() string {
	return GetSource(s.Source)
}

func (s SourceFields) GetAppSource() string {
	return GetAppSource(s.AppSource)
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
