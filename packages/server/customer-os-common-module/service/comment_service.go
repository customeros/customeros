package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type CommentService interface {
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, commentFields data_fields.CommentFields) (string, error)
}

type commentService struct {
	log      logger.Logger
	services *Services
}

func NewCommentService(log logger.Logger, services *Services) CommentService {
	return &commentService{
		log:      log,
		services: services,
	}
}

func (s *commentService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, commentFields data_fields.CommentFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "commentFields", commentFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	commentId := ""

	if utils.IfNotNilString(commentId) == "" {
		createFlow = true
		span.LogKV("flow", "create")

		// prepare missing fields
		if commentFields.CreatedAt == nil || commentFields.CreatedAt.IsZero() {
			commentFields.CreatedAt = utils.NowPtr()
		}
		if utils.IfNotNilString(commentFields.Source) == "" {
			commentFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
		if utils.IfNotNilString(commentFields.AppSource) == "" {
			commentFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}

		// validate given data exists
		if utils.IfNotNilString(commentFields.CommentedIssueId) != "" {
			exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, *commentFields.CommentedIssueId, model.NodeLabelIssue)
			if err != nil || !exists {
				err = errors.New("commented issue not found")
				tracing.TraceErr(span, err)
				return "", err
			}
		}

		if utils.IfNotNilString(commentFields.AuthorUserId) != "" {
			exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, *commentFields.AuthorUserId, model.NodeLabelUser)
			if err != nil || !exists {
				err = errors.New("author user not found")
				tracing.TraceErr(span, err)
				return "", err
			}
		}

		commentId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelIssue)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	} else {
		span.LogKV("flow", "update")
		commentId = *id

		// validate comment exists
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, commentId, model.NodeLabelComment)
		if err != nil || !exists {
			err = errors.New("comment not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, commentId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		if createFlow {
			err := s.services.Neo4jRepositories.CommentWriteRepository.Create(ctx, txWithPostCommit.Tx, tenant, commentId, commentFields)
			if err != nil {
				s.log.Errorf("Error while saving comment %s: %s", commentId, err.Error())
				return nil, err
			}
		} else {
			err := s.services.Neo4jRepositories.CommentWriteRepository.Update(ctx, txWithPostCommit.Tx, tenant, commentId, commentFields)
			if err != nil {
				s.log.Errorf("Error while updating comment %s: %s", commentId, err.Error())
				return nil, err
			}
		}
		if commentFields.ExternalSystemAvailable() {
			err := s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, commentId, model.NodeLabelComment, *commentFields.ExternalSystem)
			if err != nil {
				s.log.Errorf("Error while link comment %s with external system %s: %s", commentId, commentFields.ExternalSystem.ExternalSystemId, err.Error())
				return nil, err
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// send events
			if createFlow {
				err := s.services.RabbitMQService.PublishEvent(ctx, commentId, model.COMMENT, dto.CreateComment{commentFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateComment"))
				}
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, commentId, model.COMMENT, utils.NewEventCompletedDetails().WithCreate())
			} else {
				err := s.services.RabbitMQService.PublishEvent(ctx, commentId, model.COMMENT, dto.UpdateComment{commentFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateComment"))
				}
				if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
					s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, commentId, model.COMMENT, utils.NewEventCompletedDetails().WithUpdate())
				}
			}
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createFlow {
		span.LogFields(log.Bool("response.commentCreated", true))
	} else {
		span.LogFields(log.Bool("response.commentCreated", true))
	}

	return commentId, nil
}
