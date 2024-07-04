package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

func OrganizationsCacheHandler(serviceContainer *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/stream/organizations-cache", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys[security.KEY_TENANT_NAME].(string)

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.WriteHeader(http.StatusOK)

		// Get the http.Flusher interface to flush data to the client
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			tracing.TraceErr(span, fmt.Errorf("Streaming unsupported!"))
			c.String(http.StatusInternalServerError, "Streaming unsupported!")
			return
		}

		organizations := serviceContainer.Caches.GetOrganizations(tenant)

		span.LogFields(log.Int("count", len(organizations)))

		// Stream data in chunks
		for i := 0; i < len(organizations); i++ {

			data, err := json.Marshal(organizations[i])
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.Bool("streamed", false))
				return
			}
			jsonStr := string(data)

			_, err = fmt.Fprintf(c.Writer, "%s\n", jsonStr)
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.Bool("streamed", false))
				return
			}
			flusher.Flush()
		}

		span.LogFields(log.Bool("streamed", true))
	}
}

func OrganizationsPatchesCacheHandler(serviceContainer *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/stream/organizations-cache", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys[security.KEY_TENANT_NAME].(string)

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.WriteHeader(http.StatusOK)

		// Get the http.Flusher interface to flush data to the client
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			tracing.TraceErr(span, fmt.Errorf("Streaming unsupported!"))
			c.String(http.StatusInternalServerError, "Streaming unsupported!")
			return
		}

		patches, err := serviceContainer.CommonServices.ApiCacheService.GetPatchesForApiCache(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			span.LogFields(log.Bool("streamed", false))
			return
		}

		span.LogFields(log.Int("count", len(patches)))

		// Stream data in chunks
		for i := 0; i < len(patches); i++ {

			organizationGraph := mapper.MapEntityToOrganization(patches[i].Organization)

			if patches[i].Contacts != nil && len(patches[i].Contacts) > 0 {
				organizationGraph.Contacts = &model.ContactsPage{}
				for _, contactId := range patches[i].Contacts {
					organizationGraph.Contacts.Content = append(organizationGraph.Contacts.Content, mapper.MapEntityToContact(&neo4jentity.ContactEntity{
						Id: *contactId,
					}))
				}
			}

			if patches[i].Tags != nil && len(patches[i].Tags) > 0 {
				organizationGraph.Tags = make([]*model.Tag, 0)
				for _, tagId := range patches[i].Tags {
					organizationGraph.Tags = append(organizationGraph.Tags, mapper.MapEntityToTag(neo4jentity.TagEntity{
						Id: *tagId,
					}))
				}
			}

			if patches[i].SocialMedia != nil && len(patches[i].SocialMedia) > 0 {
				organizationGraph.SocialMedia = make([]*model.Social, 0)
				for _, socialId := range patches[i].SocialMedia {
					organizationGraph.SocialMedia = append(organizationGraph.SocialMedia, mapper.MapEntityToSocial(&neo4jentity.SocialEntity{
						Id: *socialId,
					}))
				}
			}

			if patches[i].ParentCompanies != nil && len(patches[i].ParentCompanies) > 0 {
				organizationGraph.ParentCompanies = make([]*model.LinkedOrganization, 0)
				for _, parentId := range patches[i].ParentCompanies {
					organizationGraph.ParentCompanies = append(organizationGraph.ParentCompanies, mapper.MapEntityToLinkedOrganization(&neo4jentity.OrganizationEntity{
						ID: *parentId,
					}))
				}
			}

			if patches[i].Subsidiaries != nil && len(patches[i].Subsidiaries) > 0 {
				organizationGraph.Subsidiaries = make([]*model.LinkedOrganization, 0)
				for _, subsidiaryId := range patches[i].Subsidiaries {
					organizationGraph.Subsidiaries = append(organizationGraph.Subsidiaries, mapper.MapEntityToLinkedOrganization(&neo4jentity.OrganizationEntity{
						ID: *subsidiaryId,
					}))
				}
			}

			if patches[i].Owner != nil {
				organizationGraph.Owner = mapper.MapEntityToUser(&entity.UserEntity{
					Id: *patches[i].Owner,
				})
			}

			data, err := json.Marshal(organizationGraph)
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.Bool("streamed", false))
				return
			}
			jsonStr := string(data)

			_, err = fmt.Fprintf(c.Writer, "%s\n", jsonStr)
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.Bool("streamed", false))
				return
			}
			flusher.Flush()
		}

		span.LogFields(log.Bool("streamed", true))
	}
}
