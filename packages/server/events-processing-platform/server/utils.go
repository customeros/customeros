package server

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/events-processing-platform/config"
	"strings"
)

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("%s", strings.ToUpper(cfg.ServiceName))
}
