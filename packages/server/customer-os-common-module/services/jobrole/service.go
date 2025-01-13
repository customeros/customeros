package jobrole

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type jobRoleService struct {
	neo4j  *neoRepo.Repositories
	events *events.EventsService
	org    interfaces.OrganizationService
}

func NewJobRoleService(neo4j *neoRepo.Repositories, events *events.EventsService, org interfaces.OrganizationService) interfaces.JobRoleService {
	return &jobRoleService{
		neo4j:  neo4j,
		events: events,
		org:    org,
	}
}

func (s *jobRoleService) SetOrganizationService(org interfaces.OrganizationService) {
	s.org = org
}

func (s *jobRoleService) IsInitialized() bool {
	if s.neo4j == nil || s.events == nil || s.org == nil {
		return false
	}
	return true
}

func (s *jobRoleService) getDriver() neo4j.DriverWithContext {
	return *s.neo4j.Neo4jDriver
}

func (s *jobRoleService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id, contactId, organizationId *string, dataFields data_fields.JobRoleFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("dataFields", dataFields))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	jobRoleId := ""
	var jobRoleEntity *neo4jentity.JobRoleEntity

	// identify if it is a create flow
	if utils.IfNotNilString(id) != "" {
		jobRoleId = *id
		jobRoleEntity, err = s.GetById(ctx, jobRoleId)
		if err != nil {
			return "", err
		}
	} else {
		// if job role is missing, contact is mandatory
		if utils.IfNotNilString(contactId) == "" {
			return "", errors.New("contactId is mandatory")
		}
		jobRoleEntity, err = s.IdentifyJobRole(ctx, *contactId, utils.IfNotNilString(organizationId))
		if err != nil {
			return "", err
		}
		if jobRoleEntity == nil {
			createFlow = true
		}
	}

	if createFlow {
		span.LogFields(log.String("flow", "create"))

		// validate organization exists
		if utils.IfNotNilString(organizationId) != "" {
			err = s.org.ValidateOrganizationExists(ctx, utils.IfNotNilString(organizationId))
			if err != nil {
				tracing.TraceErr(span, err)
				return "", err
			}
		}

		// set default values
		if dataFields.Source == nil {
			dataFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
		if dataFields.AppSource == nil {
			dataFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}
		// generate new id
		jobRoleId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelJobRole)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	} else {
		span.LogFields(log.String("flow", "update"))
	}

	tracing.TagEntity(span, jobRoleId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		if createFlow {
			// if organization is provided, link job role with organization
			if utils.IfNotNilString(organizationId) != "" {
				if err = s.neo4j.JobRoleWriteRepository.LinkContactWithOrganization(ctx, txWithPostCommit.Tx, common.GetContext(ctx).Tenant, jobRoleId, utils.IfNotNilString(contactId), utils.IfNotNilString(organizationId), dataFields); err != nil {
					return "", err
				}
			} else {
				// else create job role for contact
				if err = s.neo4j.JobRoleWriteRepository.CreateJobRoleForContact(ctx, txWithPostCommit.Tx, common.GetContext(ctx).Tenant, jobRoleId, utils.IfNotNilString(contactId), dataFields); err != nil {
					return "", err
				}
			}
		} else {
			if err = s.neo4j.JobRoleWriteRepository.UpdateJobRoleDetails(ctx, txWithPostCommit.Tx, common.GetContext(ctx).Tenant, jobRoleId, dataFields); err != nil {
				return "", err
			}
		}

		if utils.IfNotNilBool(dataFields.Primary) == true && utils.IfNotNilString(contactId) != "" {
			err = s.neo4j.JobRoleWriteRepository.SetOtherJobRolesForContactNonPrimaryInTx(ctx, txWithPostCommit.Tx, common.GetContext(ctx).Tenant, utils.IfNotNilString(contactId), jobRoleId)
			if err != nil {
				return nil, err
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.events.Publisher.PublishEvent(ctx, jobRoleId, model.JOB_ROLE, dto.SaveJobRole{dataFields})
			if err != nil {
				tracing.TraceErr(span, err)
			}

			if utils.IfNotNilString(contactId) != "" {
				s.events.Publisher.PublishEventCompleted(ctx, tenant, utils.IfNotNilString(contactId), model.CONTACT, utils.NewEventCompletedDetails().WithUpdate())
			}
			if utils.IfNotNilString(organizationId) != "" {
				s.events.Publisher.PublishEventCompleted(ctx, tenant, utils.IfNotNilString(organizationId), model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
			}

			return nil
		})

		return nil, nil
	})

	return jobRoleId, nil
}

func (s *jobRoleService) GetAllForContact(ctx context.Context, contactId string) (*neo4jentity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))

	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbNodes, err := s.neo4j.JobRoleReadRepository.GetAllForContact(ctx, session, common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	}

	jobRoleEntities := neo4jentity.JobRoleEntities{}
	for _, dbNode := range dbNodes {
		entity := neo4jmapper.MapDbNodeToJobRoleEntity(dbNode)
		jobRoleEntities = append(jobRoleEntities, *entity)
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) GetAllForContacts(ctx context.Context, contactIds []string) (*neo4jentity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForContacts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contactIds", contactIds))

	jobRoles, err := s.neo4j.JobRoleReadRepository.GetAllForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
	if err != nil {
		return nil, err
	}
	jobRoleEntities := neo4jentity.JobRoleEntities{}
	for _, v := range jobRoles {
		jobRoleEntity := neo4jmapper.MapDbNodeToJobRoleEntity(v.Node)
		jobRoleEntity.DataloaderKey = v.LinkedNodeId
		jobRoleEntities = append(jobRoleEntities, *jobRoleEntity)
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) GetAllForUsers(ctx context.Context, userIds []string) (*neo4jentity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForUsers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("userIds", userIds))

	jobRoles, err := s.neo4j.JobRoleReadRepository.GetAllForUsers(ctx, common.GetTenantFromContext(ctx), userIds)
	if err != nil {
		return nil, err
	}
	jobRoleEntities := neo4jentity.JobRoleEntities{}
	for _, v := range jobRoles {
		jobRoleEntity := neo4jmapper.MapDbNodeToJobRoleEntity(v.Node)
		jobRoleEntity.DataloaderKey = v.LinkedNodeId
		jobRoleEntities = append(jobRoleEntities, *jobRoleEntity)
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) GetAllForOrganization(ctx context.Context, organizationId string) (*neo4jentity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbNodes, err := s.neo4j.JobRoleReadRepository.GetAllForOrganization(ctx, session, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		return nil, err
	}

	jobRoleEntities := neo4jentity.JobRoleEntities{}
	for _, dbNode := range dbNodes {
		jobRoleEntity := neo4jmapper.MapDbNodeToJobRoleEntity(dbNode)
		jobRoleEntities = append(jobRoleEntities, *jobRoleEntity)
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) GetAllForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("organizationIds", organizationIds))

	jobRoles, err := s.neo4j.JobRoleReadRepository.GetAllForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	jobRoleEntities := neo4jentity.JobRoleEntities{}
	for _, v := range jobRoles {
		jobRoleEntity := neo4jmapper.MapDbNodeToJobRoleEntity(v.Node)
		jobRoleEntity.DataloaderKey = v.LinkedNodeId
		jobRoleEntities = append(jobRoleEntities, *jobRoleEntity)
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) DeleteJobRole(ctx context.Context, contactId, roleId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.DeleteJobRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("roleId", roleId))

	session := utils.NewNeo4jWriteSession(ctx, *s.neo4j.Neo4jDriver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, s.neo4j.JobRoleWriteRepository.DeleteJobRoleInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, roleId)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *jobRoleService) GetJobRolesByIds(ctx context.Context, ids []string) (*neo4jentity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetJobRolesByIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	jobRoleDbNodes, err := s.neo4j.JobRoleReadRepository.GetByIds(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	jobRoleEntities := neo4jentity.JobRoleEntities{}
	for _, v := range jobRoleDbNodes {
		jobRoleEntities = append(jobRoleEntities, *neo4jmapper.MapDbNodeToJobRoleEntity(v))
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) IdentifyJobRole(ctx context.Context, contactId, organizationId string) (*neo4jentity.JobRoleEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.IdentifyJobRoleId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("organizationId", organizationId))

	if contactId == "" {
		return nil, errors.New("contactId is mandatory")
	}

	var jobRoleDbNode *dbtype.Node
	var err error
	if organizationId != "" {
		jobRoleDbNode, err = s.neo4j.JobRoleReadRepository.GetJobRoleForContactAndOrganization(ctx, common.GetTenantFromContext(ctx), contactId, organizationId)
	} else {
		jobRoleDbNode, err = s.neo4j.JobRoleReadRepository.GetJobRoleForContactWithoutOrganization(ctx, common.GetTenantFromContext(ctx), contactId)
	}
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if jobRoleDbNode == nil {
		return nil, nil
	}
	return neo4jmapper.MapDbNodeToJobRoleEntity(jobRoleDbNode), nil
}

func (s *jobRoleService) GetById(ctx context.Context, jobRoleId string) (*neo4jentity.JobRoleEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("jobRoleId", jobRoleId))

	dbNode, err := s.neo4j.JobRoleReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), jobRoleId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return neo4jmapper.MapDbNodeToJobRoleEntity(dbNode), nil
}
