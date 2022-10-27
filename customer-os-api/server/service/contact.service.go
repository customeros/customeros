package service

import (
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
)

type ContactService interface {
	FindAll() ([]entity.ContactNode, error)
	FindAllByName(name string) ([]entity.ContactNode, error)
	Create(contact entity.ContactNode) (entity.ContactNode, error)
}

type neo4jContactService struct {
	cfg *config.Config
}

func NewContactService(cfg *config.Config) ContactService {
	return &neo4jContactService{
		cfg: cfg,
	}
}

func (cs *neo4jContactService) Create(aContact entity.ContactNode) (entity.ContactNode, error) {
	contact := entity.ContactNode{}
	driver, err := neo4j.NewDriver(cs.cfg.Neo4j.Target, neo4j.BasicAuth(cs.cfg.Neo4j.User, cs.cfg.Neo4j.Pwd, cs.cfg.Neo4j.Realm))
	if err != nil {
		return contact, err
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			CREATE (c:ContactNode {
				  id: randomUuid(),
				  firstName: $firstName,
				  lastName: $lastName,
				  label: $label,
				  contactType: $contactType
			})
			RETURN c { .id, .firstName, .lastName, .label, .contactType } as c`,
			map[string]interface{}{
				"firstName":   aContact.LastName,
				"lastName":    aContact.LastName,
				"label":       aContact.Label,
				"contactType": aContact.ContactType,
			})

		// Extract safe properties from the aContact node (`c`) in the first row
		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		contact, _ := record.Get("c")
		return contact, nil
	})
	if err != nil {
		return contact, err
	}

	contactData := result.(map[string]interface{})

	mapstructure.Decode(contactData, &contact)

	if err != nil {
		return contact, err
	}
	return contact, err
}

func (cs *neo4jContactService) FindAll() ([]entity.ContactNode, error) {

	driver, err := neo4j.NewDriver(cs.cfg.Neo4j.Target, neo4j.BasicAuth(cs.cfg.Neo4j.User, cs.cfg.Neo4j.Pwd, cs.cfg.Neo4j.Realm))
	// Open a new Session
	session := driver.NewSession(neo4j.SessionConfig{})

	// Close the session once this function has completed
	defer driver.Close()

	// Execute a query in a new Read Transaction
	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Get an array of IDs for the User's favorite movies

		// Retrieve a list of movies

		result, err := tx.Run("MATCH (c:Contact) RETURN c { .* } AS contact", map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		// Get a list of Movies from the Result
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		var results []map[string]interface{}
		for _, record := range records {
			person, _ := record.Get("contact")
			results = append(results, person.(map[string]interface{}))
		}
		return results, nil
	})

	if err != nil {
		return nil, err
	}
	return results.([]entity.ContactNode), nil
}

func (n neo4jContactService) FindAllByName(name string) ([]entity.ContactNode, error) {
	//TODO implement me
	panic("implement me")
}
