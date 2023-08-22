package server

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"strings"
)

//func (server *server) runMetrics(cancel context.CancelFunc) {
//	metricsServer := echo.New()
//	go func() {
//		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
//			StackSize:         stackSize,
//			DisablePrintStack: true,
//			DisableStackAll:   true,
//		}))
//		metricsServer.GET(server.cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
//		server.log.Infof("Metrics server is running on port: {%server}", server.cfg.Probes.PrometheusPort)
//		if err := metricsServer.Init(server.cfg.Probes.PrometheusPort); err != nil {
//			server.log.Errorf("metricsServer.Init: {%v}", err)
//			cancel()
//		}
//	}()
//}

//func (server *server) getGrpcMetricsCb() func(err error) {
//	return func(err error) {
//		if err != nil {
//			server.metrics.ErrorGrpcRequests.Inc()
//		} else {
//			server.metrics.SuccessGrpcRequests.Inc()
//		}
//	}
//}

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("%s", strings.ToUpper(cfg.ServiceName))
}
