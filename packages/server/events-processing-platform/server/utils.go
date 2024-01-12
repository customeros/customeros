package server

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"strings"
)

//func (Server *Server) runMetrics(cancel context.CancelFunc) {
//	metricsServer := echo.New()
//	go func() {
//		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
//			StackSize:         stackSize,
//			DisablePrintStack: true,
//			DisableStackAll:   true,
//		}))
//		metricsServer.GET(Server.cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
//		Server.log.Infof("Metrics Server is running on port: {%Server}", Server.cfg.Probes.PrometheusPort)
//		if err := metricsServer.Init(Server.cfg.Probes.PrometheusPort); err != nil {
//			Server.log.Errorf("metricsServer.Init: {%v}", err)
//			cancel()
//		}
//	}()
//}

//func (Server *Server) getGrpcMetricsCb() func(err error) {
//	return func(err error) {
//		if err != nil {
//			Server.metrics.ErrorGrpcRequests.Inc()
//		} else {
//			Server.metrics.SuccessGrpcRequests.Inc()
//		}
//	}
//}

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("%s", strings.ToUpper(cfg.ServiceName))
}
