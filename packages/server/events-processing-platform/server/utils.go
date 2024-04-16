package server

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"strings"
)

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("%s", strings.ToUpper(cfg.ServiceName))
}
