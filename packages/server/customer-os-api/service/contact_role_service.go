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
	DeleteContactRole(ctx context.Context, contactId, roleId string) (bool, error)
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
		contactRoleEntities = append(contactRoleEntities, *s.mapDbNodeToContactRoleEntity(dbNode))
	}
	return &contactRoleEntities, nil
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

func (s *contactRoleService) mapDbNodeToContactRoleEntity(node *dbtype.Node) *entity.ContactRoleEntity {
	props := utils.GetPropsFromNode(*node)
	result := entity.ContactRoleEntity{
		Id:       utils.GetStringPropOrEmpty(props, "id"),
		JobTitle: utils.GetStringPropOrEmpty(props, "jobTitle"),
		Primary:  utils.GetBoolPropOrFalse(props, "primary"),
	}
	return &result
}
