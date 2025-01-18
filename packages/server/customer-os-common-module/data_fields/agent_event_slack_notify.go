package data_fields

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type SlackNotifyEventFields struct {
	ID            string         `json:"id"`
	Tenant        string         `json:"tenant"`
	ChannelID     string         `json:"channelId"`
	Message       string         `json:"message"`
	DomainContext *DomainContext `json:"domainContext"`
}

type DomainContext struct {
	Domain                       string `json:"domain"`
	DomainRateLimit              bool   `json:"domainRateLimit"`
	MinHoursBetweenNotifications *int   `json:"minHoursBetweenNotifications"`
}

func (f SlackNotifyEventFields) Type() string {
	return "SlackNotifyEventFields"
}

func (f SlackNotifyEventFields) GenerateID() string {
	data := fmt.Sprintf("%s%s%d", f.Tenant, f.DomainContext.Domain, utils.Now().Unix())
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
