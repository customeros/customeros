package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type UserRepository interface {
	MergeUser(tenant string, syncDate time.Time, user entity.UserData) (string, error)
}

type userRepository struct {
	driver *neo4j.Driver
}

func NewUserRepository(driver *neo4j.Driver) UserRepository {
	return &userRepository{
		driver: driver,
	}
}

func (r *userRepository) MergeUser(tenant string, syncDate time.Time, user entity.UserData) (string, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MERGE (u:User)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET r.externalId=$externalId, u.id=randomUUID(), u.createdAt=$createdAt, " +
		"               u.firstName=$firstName, u.lastName=$lastName, u.readonly=$readonly, r.syncDate=$syncDate, " +
		"               u.email=$email, u:%s" +
		" ON MATCH SET u.firstName=$firstName, u.lastName=$lastName, u.readonly=$readonly, r.syncDate=$syncDate, " +
		"              u.email=$email " +
		" WITH u, t " +
		" MERGE (u)-[:USER_BELONGS_TO_TENANT]->(t)" +
		" RETURN u.id"

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "User_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": user.ExternalSystem,
				"externalId":     user.ExternalId,
				"syncDate":       syncDate,
				"firstName":      user.FirstName,
				"lastName":       user.LastName,
				"email":          user.Email,
				"createdAt":      user.CreatedAt,
				"readonly":       user.Readonly,
			})
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return "", err
	}
	return dbRecord.(string), nil
}
