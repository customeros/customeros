package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type JobRoleService interface {
	GetAllForContact(ctx context.Context, contactId string) (*neo4jentity.JobRoleEntities, error)
	GetAllForContacts(ctx context.Context, contactIds []string) (*neo4jentity.JobRoleEntities, error)
	GetAllForOrganization(ctx context.Context, organizationId string) (*neo4jentity.JobRoleEntities, error)
	GetAllForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.JobRoleEntities, error)
	DeleteJobRole(ctx context.Context, contactId, roleId string) (bool, error)
	CreateJobRole(ctx context.Context, contactId string, organizationId *string, entity *neo4jentity.JobRoleEntity) (*neo4jentity.JobRoleEntity, error)
	UpdateJobRole(ctx context.Context, contactId string, organizationId *string, entity *neo4jentity.JobRoleEntity) (*neo4jentity.JobRoleEntity, error)
	GetAllForUsers(ctx context.Context, userIds []string) (*neo4jentity.JobRoleEntities, error)
}

type jobRoleService struct {
	services *Services
}

func NewJobRoleService(services *Services) JobRoleService {
	return &jobRoleService{
		services: services,
	}
}

func (s *jobRoleService) getDriver() neo4j.DriverWithContext {
	return *s.services.Neo4jRepositories.Neo4jDriver
}

func (s *jobRoleService) GetAllForContact(ctx context.Context, contactId string) (*neo4jentity.JobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.GetAllForContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))

	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbNodes, err := s.services.Neo4jRepositories.JobRoleReadRepository.GetAllForContact(ctx, session, common.GetContext(ctx).Tenant, contactId)
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
	jobRoles, err := s.services.Neo4jRepositories.JobRoleReadRepository.GetAllForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
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

	jobRoles, err := s.services.Neo4jRepositories.JobRoleReadRepository.GetAllForUsers(ctx, common.GetTenantFromContext(ctx), userIds)
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

	dbNodes, err := s.services.Neo4jRepositories.JobRoleReadRepository.GetAllForOrganization(ctx, session, common.GetContext(ctx).Tenant, organizationId)
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

	jobRoles, err := s.services.Neo4jRepositories.JobRoleReadRepository.GetAllForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
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

func (s *jobRoleService) CreateJobRole(ctx context.Context, contactId string, organizationId *string, entity *neo4jentity.JobRoleEntity) (*neo4jentity.JobRoleEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.CreateJobRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)

	defer session.Close(ctx)
	dbNode, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if entity.Primary == true {
			err := s.services.Neo4jRepositories.JobRoleWriteRepository.SetOtherJobRolesForContactNonPrimaryInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, "")
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		roleDbNode, err := s.services.Neo4jRepositories.JobRoleWriteRepository.CreateJobRoleInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		var roleId = utils.GetPropsFromNode(*roleDbNode)["id"].(string)

		if organizationId != nil {
			if err = s.services.Neo4jRepositories.JobRoleWriteRepository.LinkWithOrganization(ctx, tx, common.GetContext(ctx).Tenant, roleId, *organizationId); err != nil {
				return nil, err
			}
		}
		return roleDbNode, nil
	})
	if err != nil {
		return nil, err
	}

	//TODO EDI

	//if organizationId != nil {
	//	s.services.OrganizationService.UpdateLastTouchpoint(ctx, *organizationId)
	//}

	return neo4jmapper.MapDbNodeToJobRoleEntity(dbNode.(*dbtype.Node)), nil
}

func (s *jobRoleService) UpdateJobRole(ctx context.Context, contactId string, organizationId *string, entity *neo4jentity.JobRoleEntity) (*neo4jentity.JobRoleEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.UpdateJobRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))
	if organizationId != nil {
		span.LogFields(log.String("organizationId", *organizationId))
	}

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	dbNode, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if entity.Primary == true {
			err := s.services.Neo4jRepositories.JobRoleWriteRepository.SetOtherJobRolesForContactNonPrimaryInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, entity.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		roleDbNode, err := s.services.Neo4jRepositories.JobRoleWriteRepository.UpdateJobRoleDetails(ctx, tx, common.GetContext(ctx).Tenant, contactId, entity.Id, *entity)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if organizationId != nil {
			err := s.services.Neo4jRepositories.JobRoleWriteRepository.LinkWithOrganization(ctx, tx, common.GetContext(ctx).Tenant, entity.Id, *organizationId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}
		return roleDbNode, nil
	})
	if err != nil {
		return nil, err
	}

	//TODO EDI
	//if organizationId != nil {
	//	s.services.OrganizationService.UpdateLastTouchpoint(ctx, *organizationId)
	//}

	return neo4jmapper.MapDbNodeToJobRoleEntity(dbNode.(*dbtype.Node)), nil
}

func (s *jobRoleService) DeleteJobRole(ctx context.Context, contactId, roleId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleService.DeleteJobRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.String("roleId", roleId))

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, s.services.Neo4jRepositories.JobRoleWriteRepository.DeleteJobRoleInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, roleId)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
