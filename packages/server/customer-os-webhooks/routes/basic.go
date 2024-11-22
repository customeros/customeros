package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (r *Registry) registerBasicRoutes() {
	r.engine.GET("/health", r.healthCheck)
	r.engine.GET("/readiness", r.readiness)
	r.engine.GET("/", r.root)

	if r.config.ApiPort == r.config.MetricsPort {
		r.engine.GET(r.config.Metrics.PrometheusPath, r.metrics)
	}
}

func (r *Registry) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func (r *Registry) readiness(c *gin.Context) {
	c.JSON(200, gin.H{"status": "READY"})
}

func (r *Registry) root(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Customer OS Webhooks",
	})
}

func (r *Registry) metrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
