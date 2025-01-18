package comment

import (
	"context"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neoRepo "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type commentService struct {
	log    logger.Logger
	neo4j  *neoRepo.Repositories
	events *events.EventsService
}

func NewCommentService(log logger.Logger, neo4j *neoRepo.Repositories, events *events.EventsService) interfaces.CommentService {
	return &commentService{
		log:    log,
		neo4j:  neo4j,
		events: events,
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

	if utils.IfNotNilString(id) == "" {
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
			exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, *commentFields.CommentedIssueId, model.NodeLabelIssue)
			if err != nil || !exists {
				err = errors.New("commented issue not found")
				tracing.TraceErr(span, err)
				return "", err
			}
		}

		if utils.IfNotNilString(commentFields.AuthorUserId) != "" {
			exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, *commentFields.AuthorUserId, model.NodeLabelUser)
			if err != nil || !exists {
				err = errors.New("author user not found")
				tracing.TraceErr(span, err)
				return "", err
			}
		}

		commentId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelIssue)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	} else {
		span.LogKV("flow", "update")
		commentId = *id

		// validate comment exists
		exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, commentId, model.NodeLabelComment)
		if err != nil || !exists {
			err = errors.New("comment not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, commentId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		if createFlow {
			err := s.neo4j.CommentWriteRepository.Create(ctx, txWithPostCommit.Tx, tenant, commentId, commentFields)
			if err != nil {
				s.log.Errorf("Error while saving comment %s: %s", commentId, err.Error())
				return nil, err
			}
		} else {
			err := s.neo4j.CommentWriteRepository.Update(ctx, txWithPostCommit.Tx, tenant, commentId, commentFields)
			if err != nil {
				s.log.Errorf("Error while updating comment %s: %s", commentId, err.Error())
				return nil, err
			}
		}
		if commentFields.ExternalSystemAvailable() {
			err := s.neo4j.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, commentId, model.NodeLabelComment, *commentFields.ExternalSystem)
			if err != nil {
				s.log.Errorf("Error while link comment %s with external system %s: %s", commentId, commentFields.ExternalSystem.ExternalSystemId, err.Error())
				return nil, err
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// send events
			if createFlow {
				err := s.events.Publisher.PublishEvent(ctx, commentId, model.COMMENT, dto.CreateComment{commentFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateComment"))
				}
				s.events.Publisher.PublishEventCompleted(ctx, tenant, commentId, model.COMMENT, utils.NewEventCompletedDetails().WithCreate())
			} else {
				err := s.events.Publisher.PublishEvent(ctx, commentId, model.COMMENT, dto.UpdateComment{commentFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateComment"))
				}
				if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
					s.events.Publisher.PublishEventCompleted(ctx, tenant, commentId, model.COMMENT, utils.NewEventCompletedDetails().WithUpdate())
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
