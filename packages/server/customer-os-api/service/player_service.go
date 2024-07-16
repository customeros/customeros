package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
)

type PlayerService interface {
	GetPlayerByAuthIdProvider(ctx context.Context, authId string, provider string) (*neo4jentity.PlayerEntity, error)
	GetPlayerForUser(ctx context.Context, tenant string, userId string) (*neo4jentity.PlayerEntity, error)
	GetUsers(ctx context.Context) (*entity.UserEntities, error)
	GetUsersByIdentityId(ctx context.Context, identityId string) (*entity.UserEntities, error)
	SetDefaultUser(ctx context.Context, playerId string, userId string) (*neo4jentity.PlayerEntity, error)
	Merge(ctx context.Context, player *neo4jentity.PlayerEntity) (*neo4jentity.PlayerEntity, error)
	Update(ctx context.Context, player *neo4jentity.PlayerEntity) (*neo4jentity.PlayerEntity, error)
}

type playerService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewPlayerService(repositories *repository.Repositories, service *Services) PlayerService {
	return &playerService{
		repositories: repositories,
		services:     service,
	}
}

func (s *playerService) getDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *playerService) GetPlayerForUser(ctx context.Context, tenant string, userId string) (*neo4jentity.PlayerEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	player, err := s.repositories.PlayerRepository.GetPlayerForUser(ctx, tenant, userId, neo4jentity.IDENTIFIES)
	if err != nil {
		return nil, err
	}

	playerEntity := neo4jmapper.MapDbNodeToPlayerEntity(*player)

	return playerEntity, nil
}

func (s *playerService) GetPlayerByAuthIdProvider(ctx context.Context, authId string, provider string) (*neo4jentity.PlayerEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	player, err := s.repositories.Neo4jRepositories.PlayerReadRepository.GetPlayerByAuthIdProvider(ctx, authId, provider)
	if err != nil {
		return nil, err
	}

	playerEntity := neo4jmapper.MapDbNodeToPlayerEntity(*player)

	return playerEntity, nil
}

func (s *playerService) GetUsers(ctx context.Context) (*entity.UserEntities, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbPlayer, err := s.repositories.PlayerRepository.GetPlayerByIdentityId(ctx, common.GetIdentityIdFromContext(ctx))
	if err != nil {
		return nil, err
	}
	player := neo4jmapper.MapDbNodeToPlayerEntity(*dbPlayer)

	usersDb, err := s.repositories.Neo4jRepositories.PlayerReadRepository.GetUsersForPlayer(ctx, []string{player.Id})
	if err != nil {
		return nil, err
	}

	users := make(entity.UserEntities, 0, len(usersDb))
	for _, dbUser := range usersDb {
		user := s.services.UserService.mapDbNodeToUserEntity(*dbUser.Node)
		s.services.UserService.addPlayerDbRelationshipToUser(*dbUser.Relationship, user)
		user.Tenant = dbUser.Tenant
		users = append(users, *user)
	}

	return &users, nil

}

func (s *playerService) GetUsersByIdentityId(ctx context.Context, identityId string) (*entity.UserEntities, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbPlayer, err := s.repositories.PlayerRepository.GetPlayerByIdentityId(ctx, identityId)
	if err != nil {
		return nil, err
	}
	player := neo4jmapper.MapDbNodeToPlayerEntity(*dbPlayer)

	usersDb, err := s.repositories.Neo4jRepositories.PlayerReadRepository.GetUsersForPlayer(ctx, []string{player.Id})
	if err != nil {
		return nil, err
	}

	users := make(entity.UserEntities, 0, len(usersDb))
	for _, dbUser := range usersDb {
		user := s.services.UserService.mapDbNodeToUserEntity(*dbUser.Node)
		s.services.UserService.addPlayerDbRelationshipToUser(*dbUser.Relationship, user)
		user.Tenant = dbUser.Tenant
		users = append(users, *user)
	}

	return &users, nil

}

func (s *playerService) setDefaultUserInDBTxWork(ctx context.Context, playerId string, userId string) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		playerDbNode, err := s.repositories.PlayerRepository.SetDefaultUserInTx(ctx, tx, playerId, userId, neo4jentity.IDENTIFIES)
		if err != nil {
			return nil, err
		}

		return playerDbNode, nil
	}
}

func (s *playerService) SetDefaultUser(ctx context.Context, playerId string, userId string) (*neo4jentity.PlayerEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	playerDbNode, err := session.ExecuteWrite(ctx, s.setDefaultUserInDBTxWork(ctx, playerId, userId))
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToPlayerEntity(*playerDbNode.(*dbtype.Node)), nil

}
func (s *playerService) createUserInDBTxWork(ctx context.Context, player *neo4jentity.PlayerEntity) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		playerDbNode, err := s.repositories.PlayerRepository.Merge(ctx, tx, player)
		if err != nil {
			return nil, err
		}

		return playerDbNode, nil
	}
}

func (s *playerService) Merge(ctx context.Context, player *neo4jentity.PlayerEntity) (*neo4jentity.PlayerEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	playerDbNode, err := session.ExecuteWrite(ctx, s.createUserInDBTxWork(ctx, player))
	if err != nil {
		return nil, err
	}

	return neo4jmapper.MapDbNodeToPlayerEntity(*playerDbNode.(*dbtype.Node)), nil
}

func (s *playerService) updateUserInDBTxWork(ctx context.Context, player *neo4jentity.PlayerEntity) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		playerDbNode, err := s.repositories.PlayerRepository.Update(ctx, tx, player)
		if err != nil {
			return nil, err
		}

		return playerDbNode, nil
	}
}

func (s *playerService) Update(ctx context.Context, player *neo4jentity.PlayerEntity) (*neo4jentity.PlayerEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	playerDbNode, err := session.ExecuteWrite(ctx, s.updateUserInDBTxWork(ctx, player))
	if err != nil {
		return nil, err
	}

	return neo4jmapper.MapDbNodeToPlayerEntity(*playerDbNode.(*dbtype.Node)), nil
}
