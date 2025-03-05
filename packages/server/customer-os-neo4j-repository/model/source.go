package model

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/constants"
	"strings"
)

type SourceFields struct {
	Source    string `json:"source"`
	AppSource string `json:"appSource"`
}

func (s SourceFields) GetSource() string {
	return GetSource(s.Source)
}

func (s SourceFields) GetAppSource() string {
	return GetAppSource(s.AppSource)
}

func GetSource(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.SourceOpenline)
}

func GetAppSource(input string) string {
	return utils.StringFirstNonEmpty(strings.TrimSpace(input), constants.AppSourceCustomerOsApi)
}
