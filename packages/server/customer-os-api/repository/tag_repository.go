package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type TagRepository interface {
	Merge(tenant string, tag entity.TagEntity) (*dbtype.Node, error)
	Update(tenant string, tag entity.TagEntity) (*dbtype.Node, error)
	// FIXME alexb refactor
	Delete(tenant string, id string) error
	// FIXME alexb refactor
	FindAll(tenant string) ([]*dbtype.Node, error)
	// FIXME alexb refactor
	FindForContact(tenant, contactId string) (*dbtype.Node, error)
}

type tagRepository struct {
	driver *neo4j.Driver
}

func NewTagRepository(driver *neo4j.Driver) TagRepository {
	return &tagRepository{
		driver: driver,
	}
}

func (r *tagRepository) Merge(tenant string, tag entity.TagEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:$name}) " +
		" ON CREATE SET " +
		"  tag.id=randomUUID(), " +
		"  tag.createdAt=$now, " +
		"  tag.updatedAt=$now, " +
		"  tag.source=$source, " +
		"  tag.appSource=$appSource, " +
		"  tag:%s" +
		" RETURN tag"

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "Tag_"+tenant),
			map[string]any{
				"tenant":    tenant,
				"name":      tag.Name,
				"source":    tag.Source,
				"appSource": tag.AppSource,
				"now":       utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *tagRepository) Update(tenant string, tag entity.TagEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
			SET tag.name=$name, tag.updatedAt=$now
			RETURN tag`,
			map[string]any{
				"tenant": tenant,
				"id":     tag.Id,
				"name":   tag.Name,
				"now":    utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *tagRepository) Delete(tenant string, id string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[r:TAG_BELONGS_TO_TENANT]-(c:Tag {id:$id})
			DELETE r, c`,
			map[string]any{
				"tenant": tenant,
				"id":     id,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *tagRepository) FindAll(tenant string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(c:Tag)
			RETURN c ORDER BY c.name`,
			map[string]any{
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	tagDbNodes := []*dbtype.Node{}
	for _, v := range records.([]*neo4j.Record) {
		tagDbNodes = append(tagDbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
	}
	return tagDbNodes, err
}

func (r *tagRepository) FindForContact(tenant, contactId string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})-[:IS_OF_TYPE]->(o:Tag)
			RETURN o`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	if err != nil {
		return nil, err
	} else if len(dbRecords.([]*neo4j.Record)) == 0 {
		return nil, nil
	} else {
		return utils.NodePtr(dbRecords.([]*neo4j.Record)[0].Values[0].(dbtype.Node)), nil
	}
}
