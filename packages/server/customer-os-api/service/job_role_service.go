package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type JobRoleService interface {
	GetAllForContact(ctx context.Context, contactId string) (*entity.JobRoleEntities, error)
	GetAllForContacts(ctx context.Context, contactIds []string) (*entity.JobRoleEntities, error)
	GetAllForOrganization(ctx context.Context, organizationId string) (*entity.JobRoleEntities, error)
	GetAllForOrganizations(ctx context.Context, organizationIds []string) (*entity.JobRoleEntities, error)
	DeleteJobRole(ctx context.Context, contactId, roleId string) (bool, error)
	CreateJobRole(ctx context.Context, contactId string, organizationId *string, entity *entity.JobRoleEntity) (*entity.JobRoleEntity, error)
	UpdateJobRole(ctx context.Context, contactId string, organizationId *string, entity *entity.JobRoleEntity) (*entity.JobRoleEntity, error)
	GetAllForUsers(ctx context.Context, userIds []string) (*entity.JobRoleEntities, error)
	mapDbNodeToJobRoleEntity(node dbtype.Node) *entity.JobRoleEntity
}

type jobRoleService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewJobRoleService(log logger.Logger, repositories *repository.Repositories, services *Services) JobRoleService {
	return &jobRoleService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *jobRoleService) getDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *jobRoleService) GetAllForContact(ctx context.Context, contactId string) (*entity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))

	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbNodes, err := s.repositories.JobRoleRepository.GetAllForContact(ctx, session, common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	}

	jobRoleEntities := entity.JobRoleEntities{}
	for _, dbNode := range dbNodes {
		jobRoleEntities = append(jobRoleEntities, *s.mapDbNodeToJobRoleEntity(*dbNode))
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) GetAllForContacts(ctx context.Context, contactIds []string) (*entity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForContacts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contactIds", contactIds))
	jobRoles, err := s.repositories.JobRoleRepository.GetAllForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
	if err != nil {
		return nil, err
	}
	jobRoleEntities := entity.JobRoleEntities{}
	for _, v := range jobRoles {
		jobRoleEntity := s.mapDbNodeToJobRoleEntity(*v.Node)
		jobRoleEntity.DataloaderKey = v.LinkedNodeId
		jobRoleEntities = append(jobRoleEntities, *jobRoleEntity)
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) GetAllForUsers(ctx context.Context, userIds []string) (*entity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForUsers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("userIds", userIds))

	jobRoles, err := s.repositories.JobRoleRepository.GetAllForUsers(ctx, common.GetTenantFromContext(ctx), userIds)
	if err != nil {
		return nil, err
	}
	jobRoleEntities := entity.JobRoleEntities{}
	for _, v := range jobRoles {
		jobRoleEntity := s.mapDbNodeToJobRoleEntity(*v.Node)
		jobRoleEntity.DataloaderKey = v.LinkedNodeId
		jobRoleEntities = append(jobRoleEntities, *jobRoleEntity)
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) GetAllForOrganization(ctx context.Context, organizationId string) (*entity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbNodes, err := s.repositories.JobRoleRepository.GetAllForOrganization(ctx, session, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		return nil, err
	}

	jobRoleEntities := entity.JobRoleEntities{}
	for _, dbNode := range dbNodes {
		jobRoleEntities = append(jobRoleEntities, *s.mapDbNodeToJobRoleEntity(*dbNode))
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) GetAllForOrganizations(ctx context.Context, organizationIds []string) (*entity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("organizationIds", organizationIds))

	jobRoles, err := s.repositories.JobRoleRepository.GetAllForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	jobRoleEntities := entity.JobRoleEntities{}
	for _, v := range jobRoles {
		jobRoleEntity := s.mapDbNodeToJobRoleEntity(*v.Node)
		jobRoleEntity.DataloaderKey = v.LinkedNodeId
		jobRoleEntities = append(jobRoleEntities, *jobRoleEntity)
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) CreateJobRole(ctx context.Context, contactId string, organizationId *string, entity *entity.JobRoleEntity) (*entity.JobRoleEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.CreateJobRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))

	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)

	defer session.Close(ctx)
	dbNode, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if entity.Primary == true {
			s.repositories.JobRoleRepository.SetOtherJobRolesForContactNonPrimaryInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, "")
		}

		roleDbNode, err := s.repositories.JobRoleRepository.CreateJobRole(ctx, tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		var roleId = utils.GetPropsFromNode(*roleDbNode)["id"].(string)

		if organizationId != nil {
			if err = s.repositories.JobRoleRepository.LinkWithOrganization(ctx, tx, common.GetContext(ctx).Tenant, roleId, *organizationId); err != nil {
				return nil, err
			}
		}
		return roleDbNode, nil
	})
	if err != nil {
		return nil, err
	}

	if organizationId != nil {
		s.services.OrganizationService.UpdateLastTouchpointSync(ctx, *organizationId)
	}

	return s.mapDbNodeToJobRoleEntity(*dbNode.(*dbtype.Node)), nil
}

func (s *jobRoleService) UpdateJobRole(ctx context.Context, contactId string, organizationId *string, entity *entity.JobRoleEntity) (*entity.JobRoleEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.UpdateJobRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))
	if organizationId != nil {
		span.LogFields(log.String("organizationId", *organizationId))
	}

	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	dbNode, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if entity.Primary == true {
			s.repositories.JobRoleRepository.SetOtherJobRolesForContactNonPrimaryInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, entity.Id)
		}

		roleDbNode, err := s.repositories.JobRoleRepository.UpdateJobRoleDetails(ctx, tx, common.GetContext(ctx).Tenant, contactId, entity.Id, *entity)
		if err != nil {
			return nil, err
		}

		if organizationId != nil {
			if err = s.repositories.JobRoleRepository.LinkWithOrganization(ctx, tx, common.GetContext(ctx).Tenant, entity.Id, *organizationId); err != nil {
				return nil, err
			}
		}
		return roleDbNode, nil
	})
	if err != nil {
		return nil, err
	}

	if organizationId != nil {
		s.services.OrganizationService.UpdateLastTouchpointSync(ctx, *organizationId)
	}

	return s.mapDbNodeToJobRoleEntity(*dbNode.(*dbtype.Node)), nil
}

func (s *jobRoleService) DeleteJobRole(ctx context.Context, contactId, roleId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.DeleteJobRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("roleId", roleId))

	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, s.repositories.JobRoleRepository.DeleteJobRoleInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, roleId)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *jobRoleService) mapDbNodeToJobRoleEntity(node dbtype.Node) *entity.JobRoleEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.JobRoleEntity{
		Id:                  utils.GetStringPropOrEmpty(props, "id"),
		JobTitle:            utils.GetStringPropOrEmpty(props, "jobTitle"),
		Description:         utils.GetStringPropOrNil(props, "description"),
		Company:             utils.GetStringPropOrNil(props, "company"),
		Primary:             utils.GetBoolPropOrFalse(props, "primary"),
		ResponsibilityLevel: utils.GetInt64PropOrZero(props, "responsibilityLevel"),
		Source:              entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:       entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:           utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:           utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:           utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:           utils.GetTimePropOrNil(props, "startedAt"),
		EndedAt:             utils.GetTimePropOrNil(props, "endedAt"),
	}
	return &result
}
