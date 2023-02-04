package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type JobRoleService interface {
	GetAllForContact(ctx context.Context, contactId string) (*entity.JobRoleEntities, error)
	GetAllForOrganization(ctx context.Context, organizationId string) (*entity.JobRoleEntities, error)
	DeleteJobRole(ctx context.Context, contactId, roleId string) (bool, error)
	CreateJobRole(ctx context.Context, contactId string, organizationId *string, entity *entity.JobRoleEntity) (*entity.JobRoleEntity, error)
	UpdateJobRole(ctx context.Context, contactId string, organizationId *string, entity *entity.JobRoleEntity) (*entity.JobRoleEntity, error)
}

type jobRoleService struct {
	repositories *repository.Repositories
}

func NewJobRoleService(repositories *repository.Repositories) JobRoleService {
	return &jobRoleService{
		repositories: repositories,
	}
}

func (s *jobRoleService) getDriver() neo4j.Driver {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *jobRoleService) GetAllForContact(ctx context.Context, contactId string) (*entity.JobRoleEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	dbNodes, err := s.repositories.JobRoleRepository.GetJobRolesForContact(session, common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	}

	jobRoleEntities := entity.JobRoleEntities{}
	for _, dbNode := range dbNodes {
		jobRoleEntities = append(jobRoleEntities, *s.mapDbNodeToJobRoleEntity(*dbNode))
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) GetAllForOrganization(ctx context.Context, organizationId string) (*entity.JobRoleEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	dbNodes, err := s.repositories.JobRoleRepository.GetJobRolesForOrganization(session, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		return nil, err
	}

	jobRoleEntities := entity.JobRoleEntities{}
	for _, dbNode := range dbNodes {
		jobRoleEntities = append(jobRoleEntities, *s.mapDbNodeToJobRoleEntity(*dbNode))
	}
	return &jobRoleEntities, nil
}

func (s *jobRoleService) CreateJobRole(ctx context.Context, contactId string, organizationId *string, entity *entity.JobRoleEntity) (*entity.JobRoleEntity, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	dbNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		if entity.Primary == true {
			s.repositories.JobRoleRepository.SetOtherJobRolesForContactNonPrimaryInTx(tx, common.GetContext(ctx).Tenant, contactId, "")
		}

		roleDbNode, err := s.repositories.JobRoleRepository.CreateJobRole(tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		var roleId = utils.GetPropsFromNode(*roleDbNode)["id"].(string)

		if organizationId != nil {
			if err = s.repositories.JobRoleRepository.LinkWithOrganization(tx, common.GetContext(ctx).Tenant, roleId, *organizationId); err != nil {
				return nil, err
			}
		}
		return roleDbNode, nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToJobRoleEntity(*dbNode.(*dbtype.Node)), nil
}

func (s *jobRoleService) UpdateJobRole(ctx context.Context, contactId string, organizationId *string, entity *entity.JobRoleEntity) (*entity.JobRoleEntity, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	dbNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		if entity.Primary == true {
			s.repositories.JobRoleRepository.SetOtherJobRolesForContactNonPrimaryInTx(tx, common.GetContext(ctx).Tenant, contactId, entity.Id)
		}

		roleDbNode, err := s.repositories.JobRoleRepository.UpdateJobRoleDetails(tx, common.GetContext(ctx).Tenant, contactId, entity.Id, *entity)
		if err != nil {
			return nil, err
		}

		if organizationId != nil {
			if err = s.repositories.JobRoleRepository.LinkWithOrganization(tx, common.GetContext(ctx).Tenant, entity.Id, *organizationId); err != nil {
				return nil, err
			}
		}
		return roleDbNode, nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToJobRoleEntity(*dbNode.(*dbtype.Node)), nil
}

func (s *jobRoleService) DeleteJobRole(ctx context.Context, contactId, roleId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return nil, s.repositories.JobRoleRepository.DeleteJobRoleInTx(tx, common.GetContext(ctx).Tenant, contactId, roleId)
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
		Primary:             utils.GetBoolPropOrFalse(props, "primary"),
		ResponsibilityLevel: utils.GetInt64PropOrZero(props, "responsibilityLevel"),
		Source:              entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:       entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:           utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:           utils.GetTimePropOrNow(props, "createdAt"),
		UpdatedAt:           utils.GetTimePropOrNow(props, "updatedAt"),
	}
	return &result
}
