package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type userRepository struct {
	driver *neo4j.Driver
}

type UserRepository interface {
	FindUserByEmail(email string) (string, string, error)
}

func NewUserRepository(driver *neo4j.Driver) UserRepository {
	return &userRepository{
		driver: driver,
	}
}

func (u *userRepository) FindUserByEmail(email string) (string, string, error) {
	session := (*u.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeRead,
			BoltLogger: neo4j.ConsoleBoltLogger(),
		},
	)
	defer session.Close()

	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
			MATCH (e:Email {email:$email})<-[:EMAIL_ASSOCIATED_WITH]-(u:User)-[:USER_BELONGS_TO_TENANT]->(t:Tenant)
			RETURN t.name, u.id`,
			map[string]interface{}{
				"email": email,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
	if err != nil {
		return "", "", err
	}
	if len(records.([]*neo4j.Record)) > 0 {
		tenant := records.([]*neo4j.Record)[0].Values[0].(string)
		userId := records.([]*neo4j.Record)[0].Values[1].(string)
		return userId, tenant, nil
	} else {
		return "", "", nil
	}
}
