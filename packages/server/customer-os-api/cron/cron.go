package cron

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"sync"
)

var jobLock1 sync.Mutex

func StartCronJobs(config *config.Config, services *service.Services) *cron.Cron {
	c := cron.New()

	lockAndRunJob(&jobLock1, config, services, refreshApiCache)

	//every 1 hour
	err := c.AddFunc("*/30 * * * * *", func() {
		go func(jobLock *sync.Mutex) {
			lockAndRunJob(jobLock, config, services, refreshApiCache)
		}(&jobLock1)
	})
	if err != nil {
		logrus.Fatalf("Could not add cron job: %v", err.Error())
	}

	c.Start()

	return c
}

func lockAndRunJob(jobLock *sync.Mutex, config *config.Config, services *service.Services, job func(config *config.Config, services *service.Services)) {
	jobLock.Lock()
	defer jobLock.Unlock()

	job(config, services)
}

func refreshApiCache(config *config.Config, services *service.Services) {
	ctx, cancel := utils.GetContextWithTimeout(context.Background(), utils.QuarterOfHourDuration)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "Cron.refreshApiCache")
	defer span.Finish()

	apiCacheList, err := services.Repositories.PostgresRepositories.ApiCacheRepository.GetAll(ctx)
	if err != nil {
		logrus.Errorf("failed to get tenants: %v", err)
		return
	}

	for _, apiCache := range apiCacheList {

		var resp neo4jEntity.OrganizationEntities
		err := json.Unmarshal([]byte(apiCache.Data), &resp)
		if err != nil {
			logrus.Errorf("failed to unmarshal data: %v", err)
			return
		}

		organizations := mapper.MapEntitiesToOrganizations(&resp)

		services.Caches.SetOrganizations(apiCache.Tenant, organizations)
	}

	span.LogFields(log.Object("tenant.count", len(apiCacheList)))
}
