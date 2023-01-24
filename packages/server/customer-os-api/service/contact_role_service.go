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

type ContactRoleService interface {
	FindAllForContact(ctx context.Context, contactId string) (*entity.ContactRoleEntities, error)
	FindAllForOrganization(ctx context.Context, organizationId string) (*entity.ContactRoleEntities, error)
	DeleteContactRole(ctx context.Context, contactId, roleId string) (bool, error)
	CreateContactRole(ctx context.Context, contactId string, organizationId *string, entity *entity.ContactRoleEntity) (*entity.ContactRoleEntity, error)
	UpdateContactRole(ctx context.Context, contactId, roleId string, organizationId *string, entity *entity.ContactRoleEntity) (*entity.ContactRoleEntity, error)
}

type contactRoleService struct {
	repositories *repository.Repositories
}

func NewContactRoleService(repositories *repository.Repositories) ContactRoleService {
	return &contactRoleService{
		repositories: repositories,
	}
}

func (s *contactRoleService) getDriver() neo4j.Driver {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *contactRoleService) FindAllForContact(ctx context.Context, contactId string) (*entity.ContactRoleEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	dbNodes, err := s.repositories.ContactRoleRepository.GetRolesForContact(session, common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	}

	contactRoleEntities := entity.ContactRoleEntities{}
	for _, dbNode := range dbNodes {
		contactRoleEntities = append(contactRoleEntities, *s.mapDbNodeToContactRoleEntity(*dbNode))
	}
	return &contactRoleEntities, nil
}

func (s *contactRoleService) FindAllForOrganization(ctx context.Context, organizationId string) (*entity.ContactRoleEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	dbNodes, err := s.repositories.ContactRoleRepository.GetRolesForOrganization(session, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		return nil, err
	}

	contactRoleEntities := entity.ContactRoleEntities{}
	for _, dbNode := range dbNodes {
		contactRoleEntities = append(contactRoleEntities, *s.mapDbNodeToContactRoleEntity(*dbNode))
	}
	return &contactRoleEntities, nil
}

func (s *contactRoleService) CreateContactRole(ctx context.Context, contactId string, organizationId *string, entity *entity.ContactRoleEntity) (*entity.ContactRoleEntity, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	dbNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		if entity.Primary == true {
			s.repositories.ContactRoleRepository.SetOtherRolesNonPrimaryInTx(tx, common.GetContext(ctx).Tenant, contactId, "")
		}

		roleDbNode, err := s.repositories.ContactRoleRepository.CreateContactRole(tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		var roleId = utils.GetPropsFromNode(*roleDbNode)["id"].(string)

		if organizationId != nil {
			if err = s.repositories.ContactRoleRepository.LinkWithOrganization(tx, common.GetContext(ctx).Tenant, roleId, *organizationId); err != nil {
				return nil, err
			}
		}
		return roleDbNode, nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactRoleEntity(*dbNode.(*dbtype.Node)), nil
}

func (s *contactRoleService) UpdateContactRole(ctx context.Context, contactId, roleId string, organizationId *string, entity *entity.ContactRoleEntity) (*entity.ContactRoleEntity, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	dbNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		if entity.Primary == true {
			s.repositories.ContactRoleRepository.SetOtherRolesNonPrimaryInTx(tx, common.GetContext(ctx).Tenant, contactId, roleId)
		}

		roleDbNode, err := s.repositories.ContactRoleRepository.UpdateContactRoleDetails(tx, common.GetContext(ctx).Tenant, contactId, roleId, *entity)
		if err != nil {
			return nil, err
		}

		if organizationId != nil {
			if err = s.repositories.ContactRoleRepository.LinkWithOrganization(tx, common.GetContext(ctx).Tenant, roleId, *organizationId); err != nil {
				return nil, err
			}
		}
		return roleDbNode, nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactRoleEntity(*dbNode.(*dbtype.Node)), nil
}

func (s *contactRoleService) DeleteContactRole(ctx context.Context, contactId, roleId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return nil, s.repositories.ContactRoleRepository.DeleteContactRoleInTx(tx, common.GetContext(ctx).Tenant, contactId, roleId)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *contactRoleService) mapDbNodeToContactRoleEntity(node dbtype.Node) *entity.ContactRoleEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.ContactRoleEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		JobTitle:      utils.GetStringPropOrEmpty(props, "jobTitle"),
		Primary:       utils.GetBoolPropOrFalse(props, "primary"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		CreatedAt:     utils.GetTimePropOrNow(props, "createdAt"),
	}
	return &result
}
