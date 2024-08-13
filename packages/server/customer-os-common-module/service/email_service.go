package service

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type emailService struct {
	services *Services
}

type EmailService interface {
	Merge(c context.Context, input neo4jentity.EmailEntity, linkWith *LinkWith) error
}

func NewEmailService(services *Services) EmailService {
	return &emailService{
		services: services,
	}
}

func (h *emailService) Merge(c context.Context, input neo4jentity.EmailEntity, linkWith *LinkWith) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EmailService.Merge")
	defer span.Finish()

	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("input", input))

	tenant := common.GetTenantFromContext(ctx)
	emailId := ""

	if input.Id == "" && input.Email == "" {
		return errors.New("email id or email is required")
	}

	if input.Id != "" {
		span.LogFields(log.String("email.id", input.Id))

		emailById, err := h.services.Neo4jRepositories.EmailReadRepository.GetById(ctx, tenant, input.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if emailById == nil {
			span.LogFields(log.Bool("email.found", false))
			return errors.New("email not found")
		}
		emailId = input.Id
		span.LogFields(log.Bool("email.found", true))
	} else {

		emailId = utils.NewUUIDIfEmpty("")
		err := h.services.Neo4jRepositories.EmailWriteRepository.CreateEmail(ctx, tenant, emailId, neo4jrepository.EmailCreateFields{
			RawEmail:  input.Email,
			CreatedAt: utils.Now(),
			SourceFields: neo4jmodel.Source{
				AppSource: input.AppSource,
			},
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		span.LogFields(log.Bool("email.created", true))
		span.LogFields(log.String("email.id", emailId))
	}

	if linkWith != nil && linkWith.Id != "" && linkWith.Type != "" && linkWith.Relationship != "" {
		if linkWith.Type.String() == neo4jenum.CONTACT.String() {
			err := h.services.Neo4jRepositories.EmailWriteRepository.LinkWithContact(ctx, tenant, linkWith.Id, emailId, "Work", true)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		} else if linkWith.Type.String() == neo4jenum.USER.String() {
			err := h.services.Neo4jRepositories.EmailWriteRepository.LinkWithUser(ctx, tenant, linkWith.Id, emailId, "Work", true)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
		//TODO continue and generify
	}

	return nil
}
