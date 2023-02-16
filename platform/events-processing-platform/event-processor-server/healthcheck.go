package server

import (
	"context"
	"github.com/AleksK1NG/es-microservice/pkg/constants"
	"github.com/heptiolabs/healthcheck"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func (server *server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()

	mux := http.NewServeMux()
	server.ps = &http.Server{
		Handler:      mux,
		Addr:         server.cfg.Probes.Port,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
	mux.HandleFunc(server.cfg.Probes.LivenessPath, health.LiveEndpoint)
	mux.HandleFunc(server.cfg.Probes.ReadinessPath, health.ReadyEndpoint)

	server.configureHealthCheckEndpoints(ctx, health)

	go func() {
		server.log.Infof("(%server) Kubernetes probes listening on port: {%server}", server.cfg.ServiceName, server.cfg.Probes.Port)
		if err := server.ps.ListenAndServe(); err != nil {
			server.log.Errorf("(ListenAndServe) err: {%v}", err)
		}
	}()
}

func (server *server) configureHealthCheckEndpoints(ctx context.Context, health healthcheck.Handler) {

	health.AddReadinessCheck(constants.MongoDB, healthcheck.AsyncWithContext(ctx, func() error {
		if err := server.mongoClient.Ping(ctx, nil); err != nil {
			server.log.Warnf("(MongoDB Readiness Check) err: {%v}", err)
			return err
		}
		return nil
	}, time.Duration(server.cfg.Probes.CheckIntervalSeconds)*time.Second))

	health.AddLivenessCheck(constants.MongoDB, healthcheck.AsyncWithContext(ctx, func() error {
		if err := server.mongoClient.Ping(ctx, nil); err != nil {
			server.log.Warnf("(MongoDB Liveness Check) err: {%v}", err)
			return err
		}
		return nil
	}, time.Duration(server.cfg.Probes.CheckIntervalSeconds)*time.Second))

	health.AddReadinessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
		_, _, err := server.elasticClient.Ping(server.cfg.Elastic.URL).Do(ctx)
		if err != nil {
			server.log.Warnf("(ElasticSearch Readiness Check) err: {%v}", err)
			return errors.Wrap(err, "client.Ping")
		}
		return nil
	}, time.Duration(server.cfg.Probes.CheckIntervalSeconds)*time.Second))

	health.AddLivenessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
		_, _, err := server.elasticClient.Ping(server.cfg.Elastic.URL).Do(ctx)
		if err != nil {
			server.log.Warnf("(ElasticSearch Liveness Check) err: {%v}", err)
			return errors.Wrap(err, "client.Ping")
		}
		return nil
	}, time.Duration(server.cfg.Probes.CheckIntervalSeconds)*time.Second))
}

func (server *server) shutDownHealthCheckServer(ctx context.Context) error {
	return server.ps.Shutdown(ctx)
}
