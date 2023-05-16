package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type PersonService interface {
	GetPersonByEmailProvider(ctx context.Context, email string, provider string) (*entity.PersonEntity, error)
	GetPersonForUser(ctx context.Context, tenant string, userId string) (*entity.PersonEntity, error)
	GetUsers(ctx context.Context) (*entity.UserEntities, error)
	GetUsersByIdentityId(ctx context.Context, identityId string) (*entity.UserEntities, error)
	SetDefaultUser(ctx context.Context, personId string, userId string) (*entity.PersonEntity, error)
	Merge(ctx context.Context, person *entity.PersonEntity) (*entity.PersonEntity, error)
	Update(ctx context.Context, person *entity.PersonEntity) (*entity.PersonEntity, error)
}

type personService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewPersonService(repositories *repository.Repositories, service *Services) PersonService {
	return &personService{
		repositories: repositories,
		services:     service,
	}
}

func (s *personService) getDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *personService) GetPersonForUser(ctx context.Context, tenant string, userId string) (*entity.PersonEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	person, err := s.repositories.PersonRepository.GetPersonForUser(ctx, tenant, userId, entity.IDENTIFIES)
	if err != nil {
		return nil, err
	}

	personEntity := s.mapDbNodeToPersonEntity(*person)

	return personEntity, nil
}

func (s *personService) GetPersonByEmailProvider(ctx context.Context, email string, provider string) (*entity.PersonEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	person, err := s.repositories.PersonRepository.GetPersonByEmailProvider(ctx, email, provider)
	if err != nil {
		return nil, err
	}

	personEntity := s.mapDbNodeToPersonEntity(*person)

	return personEntity, nil
}

func (s *personService) GetUsers(ctx context.Context) (*entity.UserEntities, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbPerson, err := s.repositories.PersonRepository.GetPersonByIdentityId(ctx, common.GetIdentityIdFromContext(ctx))
	if err != nil {
		return nil, err
	}
	person := s.mapDbNodeToPersonEntity(*dbPerson)

	usersDb, err := s.repositories.PersonRepository.GetUsersForPerson(ctx, []string{person.Id})
	if err != nil {
		return nil, err
	}

	users := make(entity.UserEntities, 0)
	for _, dbUser := range usersDb {
		user := s.services.UserService.mapDbNodeToUserEntity(*dbUser.Node)
		s.services.UserService.addPersonDbRelationshipToUser(*dbUser.Relationship, user)
		user.Tenant = dbUser.Tenant
		users = append(users, *user)
	}

	return &users, nil

}

func (s *personService) GetUsersByIdentityId(ctx context.Context, identityId string) (*entity.UserEntities, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbPerson, err := s.repositories.PersonRepository.GetPersonByIdentityId(ctx, identityId)
	if err != nil {
		return nil, err
	}
	person := s.mapDbNodeToPersonEntity(*dbPerson)

	usersDb, err := s.repositories.PersonRepository.GetUsersForPerson(ctx, []string{person.Id})
	if err != nil {
		return nil, err
	}

	users := make(entity.UserEntities, 0)
	for _, dbUser := range usersDb {
		user := s.services.UserService.mapDbNodeToUserEntity(*dbUser.Node)
		s.services.UserService.addPersonDbRelationshipToUser(*dbUser.Relationship, user)
		user.Tenant = dbUser.Tenant
		users = append(users, *user)
	}

	return &users, nil

}

func (s *personService) setDefaultUserInDBTxWork(ctx context.Context, personId string, userId string) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		personDbNode, err := s.repositories.PersonRepository.SetDefaultUserInTx(ctx, tx, personId, userId, entity.IDENTIFIES)
		if err != nil {
			return nil, err
		}

		return personDbNode, nil
	}
}

func (s *personService) SetDefaultUser(ctx context.Context, personId string, userId string) (*entity.PersonEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	personDbNode, err := session.ExecuteWrite(ctx, s.setDefaultUserInDBTxWork(ctx, personId, userId))
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToPersonEntity(*personDbNode.(*dbtype.Node)), nil

}
func (s *personService) createUserInDBTxWork(ctx context.Context, person *entity.PersonEntity) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		personDbNode, err := s.repositories.PersonRepository.Merge(ctx, tx, person)
		if err != nil {
			return nil, err
		}

		return personDbNode, nil
	}
}

func (s *personService) Merge(ctx context.Context, person *entity.PersonEntity) (*entity.PersonEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	personDbNode, err := session.ExecuteWrite(ctx, s.createUserInDBTxWork(ctx, person))
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToPersonEntity(*personDbNode.(*dbtype.Node)), nil
}

func (s *personService) updateUserInDBTxWork(ctx context.Context, person *entity.PersonEntity) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		personDbNode, err := s.repositories.PersonRepository.Update(ctx, tx, person)
		if err != nil {
			return nil, err
		}

		return personDbNode, nil
	}
}

func (s *personService) Update(ctx context.Context, person *entity.PersonEntity) (*entity.PersonEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	personDbNode, err := session.ExecuteWrite(ctx, s.updateUserInDBTxWork(ctx, person))
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToPersonEntity(*personDbNode.(*dbtype.Node)), nil
}

func (s *personService) mapDbNodeToPersonEntity(node neo4j.Node) *entity.PersonEntity {
	props := utils.GetPropsFromNode(node)

	return &entity.PersonEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Email:         utils.GetStringPropOrEmpty(props, "email"),
		Provider:      utils.GetStringPropOrEmpty(props, "provider"),
		IdentityId:    utils.GetStringPropOrNil(props, "identityId"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
}
