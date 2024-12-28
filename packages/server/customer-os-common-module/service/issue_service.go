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

type IssueService interface {
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, issueFields data_fields.IssueFields) (string, error)
	AddUserAssignee(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error
	RemoveUserAssignee(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error
	AddUserFollower(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error
	RemoveUserFollower(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error
}

type issueService struct {
	log      logger.Logger
	services *Services
}

func NewIssueService(log logger.Logger, services *Services) IssueService {
	return &issueService{
		log:      log,
		services: services,
	}
}

func (s *issueService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, issueFields data_fields.IssueFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "issueFields", issueFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	issueId := ""

	if utils.IfNotNilString(id) == "" {
		createFlow = true
		span.LogKV("flow", "create")

		// prepare missing fields
		if issueFields.CreatedAt == nil || issueFields.CreatedAt.IsZero() {
			issueFields.CreatedAt = utils.NowPtr()
		}
		if utils.IfNotNilString(issueFields.Source) == "" {
			issueFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
		if utils.IfNotNilString(issueFields.AppSource) == "" {
			issueFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}

		// validate given data exists
		if utils.IfNotNilString(issueFields.ReportedByOrganizationId) != "" {
			exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, *issueFields.ReportedByOrganizationId, model.NodeLabelOrganization)
			if err != nil || !exists {
				err = errors.New("reported by organization not found")
				tracing.TraceErr(span, err)
				return "", err
			}
		}

		if utils.IfNotNilString(issueFields.SubmittedByOrganizationId) != "" {
			exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, *issueFields.SubmittedByOrganizationId, model.NodeLabelOrganization)
			if err != nil || !exists {
				err = errors.New("submitted by organization not found")
				tracing.TraceErr(span, err)
				return "", err
			}
		}

		if utils.IfNotNilString(issueFields.SubmittedByUserId) != "" {
			exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, *issueFields.SubmittedByUserId, model.NodeLabelUser)
			if err != nil || !exists {
				err = errors.New("submitted by user not found")
				tracing.TraceErr(span, err)
				return "", err
			}
		}

		issueId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelIssue)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	} else {
		span.LogKV("flow", "update")
		issueId = *id

		// validate issue exists
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, issueId, model.NodeLabelIssue)
		if err != nil || !exists {
			err = errors.New("issue not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, issueId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		if createFlow {
			innerErr := s.services.Neo4jRepositories.IssueWriteRepository.Create(ctx, txWithPostCommit.Tx, tenant, issueId, issueFields)
			if innerErr != nil {
				s.log.Errorf("Error while saving issue %s: %s", issueId, err.Error())
				return nil, innerErr
			}
		} else {
			innerErr := s.services.Neo4jRepositories.IssueWriteRepository.Update(ctx, txWithPostCommit.Tx, tenant, issueId, issueFields)
			if innerErr != nil {
				s.log.Errorf("Error while updating issue %s: %s", issueId, err.Error())
				return nil, innerErr
			}
		}
		if issueFields.ExternalSystemAvailable() {
			innerErr := s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, issueId, model.NodeLabelIssue, *issueFields.ExternalSystem)
			if err != nil {
				s.log.Errorf("Error while link issue %s with external system %s: %s", issueId, issueFields.ExternalSystem.ExternalSystemId, err.Error())
				return nil, innerErr
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// send events
			if createFlow {
				if utils.IfNotNilString(issueFields.ReportedByOrganizationId) != "" {
					err = s.services.OrganizationService.RequestRefreshLastTouchpoint(ctx, *issueFields.ReportedByOrganizationId)
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "unable to request refresh last touchpoint"))
					}
				}
				err = s.services.RabbitMQService.PublishEvent(ctx, issueId, model.ISSUE, dto.CreateIssue{issueFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateIssue"))
				}
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, issueId, model.ISSUE, utils.NewEventCompletedDetails().WithCreate())
			} else {
				err = s.services.RabbitMQService.PublishEvent(ctx, issueId, model.ISSUE, dto.UpdateIssue{issueFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateIssue"))
				}
				if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
					s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, issueId, model.ISSUE, utils.NewEventCompletedDetails().WithUpdate())
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
		span.LogFields(log.Bool("response.issueCreated", true))
	} else {
		span.LogFields(log.Bool("response.issueUpdated", true))
	}

	return issueId, nil
}

func (s *issueService) AddUserAssignee(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.AddUserAssignee")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("userId", userId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate issue exists
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, issueId, model.NodeLabelIssue)
	if err != nil || !exists {
		err = errors.New("issue not found")
		tracing.TraceErr(span, err)
		return err
	}
	tracing.TagEntity(span, issueId)

	// validate user exists
	exists, err = s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, userId, model.NodeLabelUser)
	if err != nil || !exists {
		err = errors.New("user not found")
		tracing.TraceErr(span, err)
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		err = s.services.Neo4jRepositories.IssueWriteRepository.AddUserAssignee(ctx, txWithPostCommit.Tx, tenant, issueId, userId)
		if err != nil {
			s.log.Errorf("Error while adding user assignee %s to issue %s: %s", userId, issueId, err.Error())
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.services.RabbitMQService.PublishEvent(ctx, issueId, model.ISSUE, dto.AddUserAssigneeToIssue{UserID: userId})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddUserAssigneeToIssue"))
			}
			if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, issueId, model.ISSUE, utils.NewEventCompletedDetails().WithUpdate())
			}
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *issueService) RemoveUserAssignee(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.RemoveUserAssignee")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("userId", userId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate issue exists
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, issueId, model.NodeLabelIssue)
	if err != nil || !exists {
		err = errors.New("issue not found")
		tracing.TraceErr(span, err)
		return err
	}
	tracing.TagEntity(span, issueId)

	// validate user exists
	exists, err = s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, userId, model.NodeLabelUser)
	if err != nil || !exists {
		err = errors.New("user not found")
		tracing.TraceErr(span, err)
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		err = s.services.Neo4jRepositories.IssueWriteRepository.RemoveUserAssignee(ctx, txWithPostCommit.Tx, tenant, issueId, userId)
		if err != nil {
			s.log.Errorf("Error while adding user assignee %s to issue %s: %s", userId, issueId, err.Error())
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.services.RabbitMQService.PublishEvent(ctx, issueId, model.ISSUE, dto.RemoveUserAssigneeFromIssue{UserID: userId})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message RemoveUserAssigneeFromIssue"))
			}
			if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, issueId, model.ISSUE, utils.NewEventCompletedDetails().WithUpdate())
			}
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *issueService) AddUserFollower(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.AddUserFollower")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("userId", userId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate issue exists
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, issueId, model.NodeLabelIssue)
	if err != nil || !exists {
		err = errors.New("issue not found")
		tracing.TraceErr(span, err)
		return err
	}
	tracing.TagEntity(span, issueId)

	// validate user exists
	exists, err = s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, userId, model.NodeLabelUser)
	if err != nil || !exists {
		err = errors.New("user not found")
		tracing.TraceErr(span, err)
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		err = s.services.Neo4jRepositories.IssueWriteRepository.AddUserFollower(ctx, txWithPostCommit.Tx, tenant, issueId, userId)
		if err != nil {
			s.log.Errorf("Error while adding user assignee %s to issue %s: %s", userId, issueId, err.Error())
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.services.RabbitMQService.PublishEvent(ctx, issueId, model.ISSUE, dto.AddUserFollowerToIssue{UserID: userId})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddUserFollowerToIssue"))
			}
			if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, issueId, model.ISSUE, utils.NewEventCompletedDetails().WithUpdate())
			}
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *issueService) RemoveUserFollower(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.RemoveUserFollower")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("userId", userId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate issue exists
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, issueId, model.NodeLabelIssue)
	if err != nil || !exists {
		err = errors.New("issue not found")
		tracing.TraceErr(span, err)
		return err
	}
	tracing.TagEntity(span, issueId)

	// validate user exists
	exists, err = s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, userId, model.NodeLabelUser)
	if err != nil || !exists {
		err = errors.New("user not found")
		tracing.TraceErr(span, err)
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		err = s.services.Neo4jRepositories.IssueWriteRepository.RemoveUserFollower(ctx, txWithPostCommit.Tx, tenant, issueId, userId)
		if err != nil {
			s.log.Errorf("Error while adding user assignee %s to issue %s: %s", userId, issueId, err.Error())
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.services.RabbitMQService.PublishEvent(ctx, issueId, model.ISSUE, dto.RemoveUserFollowerFromIssue{UserID: userId})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message RemoveUserFollowerFromIssue"))
			}
			if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, issueId, model.ISSUE, utils.NewEventCompletedDetails().WithUpdate())
			}
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
