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
	MergeEmail(tenant, userId, email, externalSystem string, createdAt time.Time) error
	GetUserIdForExternalId(tenant, userExternalId, externalSystem string) (string, error)
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

	// Create new User if it does not exist
	// If User exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternateUser for incoming source, with a new relationship 'ALTERNATE'
	// Link User with Tenant
	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MERGE (u:User)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET r.externalId=$externalId, r.externalOwnerId=$externalOwnerId, r.syncDate=$syncDate, u.id=randomUUID(), " +
		"				u.createdAt=$createdAt, u.updatedAt=$createdAt, " +
		"               u.firstName=$firstName, u.lastName=$lastName, " +
		"               u.source=$source, u.sourceOfTruth=$sourceOfTruth, u.appSource=$appSource, " +
		"               u:%s" +
		" ON MATCH SET 	r.syncDate = CASE WHEN u.sourceOfTruth=$sourceOfTruth THEN $syncDate ELSE r.syncDate END, " +
		"				u.firstName = CASE WHEN u.sourceOfTruth=$sourceOfTruth THEN $firstName ELSE u.firstName END, " +
		"				u.lastName = CASE WHEN u.sourceOfTruth=$sourceOfTruth THEN $lastName ELSE u.lastName END, " +
		"				u.updatedAt = CASE WHEN u.sourceOfTruth=$sourceOfTruth THEN $now ELSE u.updatedAt END " +
		" WITH u, t " +
		" MERGE (u)-[:USER_BELONGS_TO_TENANT]->(t)" +
		" WITH u " +
		" FOREACH (x in CASE WHEN u.sourceOfTruth <> $sourceOfTruth THEN [u] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateUser {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.firstName=$firstName, alt.lastName=$lastName " +
		") " +
		" RETURN u.id"

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "User_"+tenant),
			map[string]interface{}{
				"tenant":          tenant,
				"externalSystem":  user.ExternalSystem,
				"externalId":      user.ExternalId,
				"externalOwnerId": user.ExternalOwnerId,
				"syncDate":        syncDate,
				"firstName":       user.FirstName,
				"lastName":        user.LastName,
				"email":           user.Email,
				"createdAt":       user.CreatedAt,
				"source":          user.ExternalSystem,
				"sourceOfTruth":   user.ExternalSystem,
				"appSource":       user.ExternalSystem,
				"now":             time.Now().UTC(),
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

func (r *userRepository) MergeEmail(tenant, userId, email, externalSystem string, createdAt time.Time) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (e:Email {email: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET " +
		"				e.id=randomUUID(), " +
		"				e.createdAt=$now, " +
		"				e.updatedAt=$now, " +
		"				e.source=$source, " +
		"				e.sourceOfTruth=$sourceOfTruth, " +
		"				e.appSource=$appSource, " +
		"				e:%s " +
		" WITH DISTINCT u, e " +
		" MERGE (u)-[rel:HAS]->(e) " +
		" ON CREATE SET rel.primary=true, " +
		"				rel.label=$label "
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Email_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"userId":        userId,
				"email":         email,
				"createdAt":     createdAt,
				"label":         "WORK",
				"source":        externalSystem,
				"sourceOfTruth": externalSystem,
				"appSource":     externalSystem,
				"now":           time.Now().UTC(),
			})
		return nil, err
	})
	return err
}

func (r *userRepository) GetUserIdForExternalId(tenant, userExternalId, externalSystem string) (string, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := " MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MATCH (u:User)-[:IS_LINKED_WITH {externalId:$userExternalId}]->(e) " +
		" RETURN u.id "
	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(query,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"userExternalId": userExternalId,
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
