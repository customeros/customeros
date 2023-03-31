package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type PhoneNumberRepository interface {
	GetIdIfExists(ctx context.Context, tenant, phoneNumber string) (string, error)
	CreatePhoneNumber(ctx context.Context, tenant, id, rawPhoneNumber string) error
}

type phoneNumberRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPhoneNumberRepository(driver *neo4j.DriverWithContext) PhoneNumberRepository {
	return &phoneNumberRepository{
		driver: driver,
	}
}

func (r *phoneNumberRepository) GetIdIfExists(ctx context.Context, tenant string, phoneNumber string) (string, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (p:PhoneNumber_%s) WHERE p.e164 = $phoneNumber OR p.rawPhoneNumber = $phoneNumber RETURN p.id LIMIT 1"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"phoneNumber": phoneNumber,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return "", err
	}
	if len(result.([]*db.Record)) == 0 {
		return "", nil
	}
	return result.([]*db.Record)[0].Values[0].(string), err
}

func (r *phoneNumberRepository) CreatePhoneNumber(ctx context.Context, tenant, id, rawPhoneNumber string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber:PhoneNumber_%s {id:$id}) " +
		" ON CREATE SET p.rawPhoneNumber = $rawPhoneNumber"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"rawPhoneNumber": rawPhoneNumber,
				"id":             id,
				"tenant":         tenant,
			})
		return nil, err
	})
	return err
}
