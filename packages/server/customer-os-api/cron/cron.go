package cron

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
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

	//every 2 minutes
	err := c.AddFunc("* */2 * * * *", func() {
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

	for _, tenantApiCache := range apiCacheList {

		var tenantCachedData []*commonService.ApiCacheOrganization

		err = json.Unmarshal([]byte(tenantApiCache.Data), &tenantCachedData)
		if err != nil {
			logrus.Errorf("failed to unmarshal data: %v", err)
			return
		}

		graphData := make([]*model.Organization, 0)

		for _, rowData := range tenantCachedData {
			organizationGraph := mapper.MapEntityToOrganization(rowData.Organization)

			if rowData.Contacts != nil && len(rowData.Contacts) > 0 {
				organizationGraph.Contacts = &model.ContactsPage{}
				for _, contactId := range rowData.Contacts {
					organizationGraph.Contacts.Content = append(organizationGraph.Contacts.Content, mapper.MapEntityToContact(&neo4jentity.ContactEntity{
						Id: *contactId,
					}))
				}
			}

			if rowData.Tags != nil && len(rowData.Tags) > 0 {
				organizationGraph.Tags = make([]*model.Tag, 0)
				for _, tagId := range rowData.Tags {
					organizationGraph.Tags = append(organizationGraph.Tags, mapper.MapEntityToTag(neo4jentity.TagEntity{
						Id: *tagId,
					}))
				}
			}

			if rowData.SocialMedia != nil && len(rowData.SocialMedia) > 0 {
				organizationGraph.SocialMedia = make([]*model.Social, 0)
				for _, socialId := range rowData.SocialMedia {
					organizationGraph.SocialMedia = append(organizationGraph.SocialMedia, mapper.MapEntityToSocial(&neo4jentity.SocialEntity{
						Id: *socialId,
					}))
				}
			}

			if rowData.ParentCompanies != nil && len(rowData.ParentCompanies) > 0 {
				organizationGraph.ParentCompanies = make([]*model.LinkedOrganization, 0)
				for _, parentId := range rowData.ParentCompanies {
					organizationGraph.ParentCompanies = append(organizationGraph.ParentCompanies, mapper.MapEntityToLinkedOrganization(&neo4jentity.OrganizationEntity{
						ID: *parentId,
					}))
				}
			}

			if rowData.Subsidiaries != nil && len(rowData.Subsidiaries) > 0 {
				organizationGraph.Subsidiaries = make([]*model.LinkedOrganization, 0)
				for _, subsidiaryId := range rowData.Subsidiaries {
					organizationGraph.Subsidiaries = append(organizationGraph.Subsidiaries, mapper.MapEntityToLinkedOrganization(&neo4jentity.OrganizationEntity{
						ID: *subsidiaryId,
					}))
				}
			}

			if rowData.Owner != nil {
				organizationGraph.Owner = mapper.MapEntityToUser(&entity.UserEntity{
					Id: *rowData.Owner,
				})
			}

			graphData = append(graphData, organizationGraph)

		}

		services.Caches.SetOrganizations(tenantApiCache.Tenant, graphData)
	}

	span.LogFields(log.Object("tenant.count", len(apiCacheList)))
}
