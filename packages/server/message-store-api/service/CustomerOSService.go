package service

import (
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	msProto "github.com/openline-ai/openline-customer-os/packages/server/message-store-api/proto/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository/entity"
	"regexp"
	"strings"
	"time"
)

type customerOSService struct {
	driver               *neo4j.Driver
	postgresRepositories *repository.PostgresRepositories
	commonStoreService   *commonStoreService
}

type CustomerOSService interface {
	ContactByIdExists(contactId string) (bool, error)
	ContactByPhoneExists(e164 string) (bool, error)

	GetUserByEmail(email string) (*User, error)
	GetContactById(id string) (*Contact, error)
	GetContactByEmail(email string) (*Contact, error)

	CreateContactWithEmail(tenant string, email string) (*Contact, error)
	CreateContactWithPhone(tenant string, phone string) (*Contact, error)

	ConversationByIdExists(tenant string, conversationId string) (bool, error)

	GetConversations(tenant string) ([]Conversation, error)
	GetConversationById(tenant string, conversationId string) (Conversation, error)

	GetWebChatConversationWithContactInitiator(tenant string, contactId string) (*Conversation, error)
	GetWebChatConversationWithUserInitiator(tenant string, userId string) (*Conversation, error)

	CreateConversation(tenant string, initiatorId string, initiatorFirstName string, initiatorLastName string, initiatorUsername string, initiatorType entity.SenderType, channel entity.EventType) (*Conversation, error)
	UpdateConversation(tenant string, conversationId string, participantId string, participantType entity.SenderType, lastSenderFirstName string, lastSenderLastName string, lastContentPreview string) (string, error)
}

type User struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
}

type Contact struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
}

type Conversation struct {
	Id        string
	StartedAt time.Time
	UpdatedAt time.Time
	Channel   string
	Status    string

	InitiatorFirstName  string
	InitiatorLastName   string
	InitiatorUsername   string
	InitiatorType       string
	LastSenderId        string
	LastSenderType      string
	LastSenderFirstName string
	LastSenderLastName  string
	LastContentPreview  string
}

func parseEmail(email string) (string, string) {
	re := regexp.MustCompile("^\"{0,1}([^\"]*)\"{0,1}[ ]*<(.*)>$")
	matches := re.FindStringSubmatch(strings.Trim(email, " "))
	if matches != nil {
		return strings.Trim(matches[1], " "), matches[2]
	}
	return "", email
}

func (s *customerOSService) ContactByIdExists(contactId string) (bool, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	params := map[string]interface{}{
		"contactId": contactId,
	}
	query := "MATCH (n:Contact {id:$contactId}) RETURN count(*)"
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(query, params)
		return nil, err
	})

	//TODO check here if the error is correct
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *customerOSService) ContactByPhoneExists(e164 string) (bool, error) {

	//graphqlRequest := graphql.NewRequest(`
	//			query ($e164: String!) {
	//				contact_ByPhone(e164: $e164){firstName,lastName,id}
	//			}
	//`)
	//
	//graphqlRequest.Var("e164", e164)
	//graphqlRequest.Header.Add("X-Openline-API-KEY", s.config.Service.CustomerOsAPIKey)
	//var graphqlResponse map[string]map[string]string
	//if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
	//	return false, err
	//}
	return true, nil
}

func (s *customerOSService) GetUserByEmail(email string) (*User, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (:Email {email:$email})<-[:HAS]-(u:User)
			RETURN u`,
			map[string]any{
				"email": email,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	if err != nil {
		return nil, err
	}
	if len(dbRecords.([]*db.Record)) == 0 {
		return nil, errors.New("user not found")
	}

	dbNode := (dbRecords.([]*db.Record))[0].Values[0].(dbtype.Node)
	return mapNodeToUser(&dbNode), nil
}

func (s *customerOSService) GetContactById(id string) (*Contact, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`MATCH (c:Contact{id: $id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c:Contact)-[:HAS]->(p:Email{primary: true})
            RETURN c.id, c.firstName, c.lastName, p.email`,
			map[string]any{
				"tenant": "openline", //TODO discuss with customerOS team
				"id":     id,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	if err != nil {
		return nil, err
	}
	if len(dbRecords.([]*db.Record)) == 0 {
		return nil, errors.New("no contact found")
	}

	idd := (dbRecords.([]*db.Record))[0].Values[0].(string)
	firstName := (dbRecords.([]*db.Record))[0].Values[1].(string)
	lastName := (dbRecords.([]*db.Record))[0].Values[2].(string)
	em := (dbRecords.([]*db.Record))[0].Values[3].(string)
	return mapNodeToContact(idd, firstName, lastName, em), nil
}

func (s *customerOSService) GetContactByEmail(email string) (*Contact, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c:Contact)-[:HAS]->(p:Email {email:$email})
            RETURN c.id, c.firstName, c.lastName, p.email`,
			map[string]any{
				"tenant": "openline", //TODO discuss with customerOS team
				"email":  email,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	if err != nil {
		return nil, err
	}
	if len(dbRecords.([]*db.Record)) == 0 {
		return nil, errors.New("no contact found")
	}

	id := (dbRecords.([]*db.Record))[0].Values[0].(string)
	var firstName string
	var lastName string
	var em string

	if (dbRecords.([]*db.Record))[0].Values[1] == nil {
		firstName = ""
	} else {
		firstName = (dbRecords.([]*db.Record))[0].Values[1].(string)
	}
	if (dbRecords.([]*db.Record))[0].Values[2] == nil {
		lastName = ""
	} else {
		lastName = (dbRecords.([]*db.Record))[0].Values[2].(string)
	}
	if (dbRecords.([]*db.Record))[0].Values[3] == nil {
		em = ""
	} else {
		em = (dbRecords.([]*db.Record))[0].Values[3].(string)
	}

	return mapNodeToContact(id, firstName, lastName, em), nil
}

func (s *customerOSService) CreateContactWithEmail(tenant string, email string) (*Contact, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	contact, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		//create the contact
		contactQuery := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (c:Contact {id:randomUUID()})-[:CONTACT_BELONGS_TO_TENANT]->(t) ON CREATE SET" +
			" c.createdAt=$createdAt, " +
			" c.updatedAt=$createdAt, " +
			" c.source=$source, " +
			" c.sourceOfTruth=$sourceOfTruth, " +
			" c.appSource=$appSource, " +
			" c:%s " +
			" RETURN c.id, c.firstName, c.lastName"

		contactQueryResult, err := tx.Run(fmt.Sprintf(contactQuery, "Contact_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"createdAt":     time.Now().UTC(),
				"source":        "openline",
				"sourceOfTruth": "openline",
				"appSource":     "message-store-api",
			})

		contact, err := contactQueryResult.Single()

		if err != nil {
			return nil, err
		}

		//create the email
		emailQuery := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
			" MERGE (c)-[r:HAS]->(e:Email {email: $email}) " +
			" ON CREATE SET e.label=$label, r.primary=$primary, e.id=randomUUID(), e.createdAt=$now, e.updatedAt=$now," +
			" e.source=$source, e.sourceOfTruth=$sourceOfTruth, e.appSource=$appSource, e:%s " +
			" ON MATCH SET e.label=$label, r.primary=$primary, e.updatedAt=$now " +
			" RETURN e"

		contactId := contact.Values[0].(string)
		_, err = tx.Run(fmt.Sprintf(emailQuery, "Email_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"email":         email,
				"label":         "WORK",
				"primary":       true,
				"source":        "openline",
				"sourceOfTruth": "openline",
				"appSource":     "message-store-api",
				"now":           time.Now().UTC(),
			})

		if err != nil {
			return nil, err
		}

		return mapNodeToContact(contactId, "", "", email), nil
	})
	if err != nil {
		return nil, err
	} else {
		return contact.(*Contact), nil
	}
}

func (s *customerOSService) CreateContactWithPhone(tenant string, phone string) (*Contact, error) {
	//graphqlRequest := graphql.NewRequest(`
	//	mutation CreateContact ($firstName: String!, $lastName: String!, $e164: String!) {
	//	  contact_Create(input: {
	//        firstName: $firstName,
	//		lastName: $lastName,
	//	    phoneNumber:{e164:  $e164, label: WORK}
	//	  }) {
	//		  id
	//	  }
	//	}
	//`)
	//
	//graphqlRequest.Var("firstName", firstName)
	//graphqlRequest.Var("lastName", lastName)
	//graphqlRequest.Var("e164", phone)
	//graphqlRequest.Header.Add("X-Openline-API-KEY", s.config.Service.CustomerOsAPIKey)
	//var graphqlResponse map[string]map[string]string
	//if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
	//	return "", err
	//}
	//return graphqlResponse["contact_Create"]["id"], nil
	return nil, nil
}

func (s *customerOSService) ConversationByIdExists(tenant string, conversationId string) (bool, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run("MATCH (c:Conversation_"+tenant+"{id:$conversationId}) RETURN c",
			map[string]any{
				"conversationId": conversationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})

	if err != nil {
		return false, err
	}

	if len(dbRecords.([]*db.Record)) == 0 {
		return false, nil
	}

	return true, nil
}

func (s *customerOSService) GetConversations(tenant string) ([]Conversation, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {

		//todo move order by as param
		if queryResult, err := tx.Run("MATCH (c:Conversation_"+tenant+") RETURN c order by c.updatedAt desc", map[string]any{
			"tenant": tenant,
		}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})

	if err != nil {
		return nil, err
	}

	var conversations []Conversation
	for _, v := range dbRecords.([]*neo4j.Record) {
		node := v.Values[0].(neo4j.Node)
		conversations = append(conversations, *mapNodeToConversation(&node))
	}

	return conversations, nil
}

func (s *customerOSService) GetConversationById(tenant string, conversationId string) (*Conversation, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	conversationNode, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {

		if queryResult, err := tx.Run("MATCH (c:Conversation_"+tenant+"{id: $conversationId}) RETURN c", map[string]any{
			"conversationId": conversationId,
		}); err != nil {
			return nil, err
		} else {
			record, err := queryResult.Single()
			if err != nil {
				return nil, err
			}
			if record == nil {
				return nil, errors.New("conversation not found")
			}

			return utils.NodePtr(record.Values[0].(neo4j.Node)), nil
		}
	})
	if err != nil {
		return nil, err
	}

	return mapNodeToConversation(conversationNode.(*dbtype.Node)), nil
}

func (s *customerOSService) GetWebChatConversationWithContactInitiator(tenant string, contactId string) (*Conversation, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	conversationNode, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`match(o:Conversation{status:"ACTIVE", channel: "WEB_CHAT"})<-[:INITIATED]-(c:Contact{id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant{name:$tenant}) return o`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			}); err != nil {
			return nil, err
		} else {
			record, err := queryResult.Single()
			if err != nil && err.Error() != "Result contains no more records" {
				return nil, err
			}
			if record != nil {
				return record.Values[0].(dbtype.Node), nil
			} else {
				return nil, nil
			}
		}
	})

	if err != nil {
		return nil, err
	}

	if conversationNode != nil {
		node := conversationNode.(dbtype.Node)
		return mapNodeToConversation(&node), nil
	} else {
		return nil, nil
	}
}

func (s *customerOSService) GetWebChatConversationWithUserInitiator(tenant string, userId string) (*Conversation, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	conversationNode, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`match(o:Conversation{status:"ACTIVE", channel: "WEB_CHAT"})<-[:INITIATED]-(u:User{id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant{name:$tenant}) return o`,
			map[string]any{
				"tenant": tenant,
				"userId": userId,
			}); err != nil {
			return nil, err
		} else {
			record, err := queryResult.Single()
			if err != nil && err.Error() != "Result contains no more records" {
				return nil, err
			}
			if record != nil {
				return record.Values[0].(dbtype.Node), nil
			} else {
				return nil, nil
			}
		}
	})

	if err != nil {
		return nil, err
	}

	if conversationNode != nil {
		node := conversationNode.(dbtype.Node)
		return mapNodeToConversation(&node), nil
	} else {
		return nil, nil
	}
}

func (s *customerOSService) CreateConversation(tenant string, initiatorId string, initiatorFirstName string, initiatorLastName string, initiatorUsername string, initiatorType entity.SenderType, channel entity.EventType) (*Conversation, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	contactIds := []string{}
	userIds := []string{}
	initiatorTypeStr := ""

	if initiatorType == entity.CONTACT {
		contactIds = append(contactIds, initiatorId)
		initiatorTypeStr = "CONTACT"
	} else if initiatorType == entity.USER {
		userIds = append(userIds, initiatorId)
		initiatorTypeStr = "USER"
	}

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		query := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (o:Conversation {id:randomUUID()}) " +
			" ON CREATE SET o.startedAt=$startedAt, o.updatedAt=$updatedAt, " +
			" o.messageCount=0, o.channel=$channel, o.status=$status, o.source=$source, o.sourceOfTruth=$sourceOfTruth, " +
			" o.initiatorFirstName=$initiatorFirstName, o.initiatorLastName=$initiatorLastName, o.initiatorUsername=$initiatorUsername, o.initiatorType=$initiatorType, " +
			" o.source=$source, o.sourceOfTruth=$sourceOfTruth, " +
			" o.lastSenderId=$lastSenderId, o.lastSenderType=$lastSenderType, o.lastSenderFirstName=$lastSenderFirstName, o.lastSenderLastName=$lastSenderLastName, o.lastContentPreview=$lastContentPreview, " +
			" o.appSource=$appSource, o:%s " +
			" %s %s " +
			" RETURN DISTINCT o"
		queryLinkWithContacts := ""
		if len(contactIds) > 0 {
			queryLinkWithContacts = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) WHERE c.id in $contactIds " +
				" MERGE (c)-[:PARTICIPATES]->(o) " +
				" WITH DISTINCT t, o " +
				" OPTIONAL MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) WHERE c.id in $contactIds " +
				" MERGE (c)-[:INITIATED]->(o) "
		}
		queryLinkWithUsers := ""
		if len(userIds) > 0 {
			queryLinkWithUsers = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (u:User)-[:USER_BELONGS_TO_TENANT]->(t) WHERE u.id in $userIds " +
				" MERGE (u)-[:PARTICIPATES]->(o) " +
				" WITH DISTINCT t, o " +
				" OPTIONAL MATCH (u:User)-[:USER_BELONGS_TO_TENANT]->(t) WHERE u.id in $userIds " +
				" MERGE (u)-[:INITIATED]->(o) "
		}
		utc := time.Now().UTC()
		queryResult, err := tx.Run(fmt.Sprintf(query, "Conversation_"+tenant, queryLinkWithContacts, queryLinkWithUsers),
			map[string]interface{}{
				"tenant":              tenant,
				"source":              "openline",
				"sourceOfTruth":       "openline",
				"appSource":           "manual",
				"status":              "ACTIVE",
				"initiatorFirstName":  initiatorFirstName,
				"initiatorLastName":   initiatorLastName,
				"initiatorUsername":   initiatorUsername,
				"initiatorType":       initiatorTypeStr,
				"startedAt":           utc,
				"updatedAt":           utc,
				"channel":             channel,
				"contactIds":          contactIds,
				"userIds":             userIds,
				"lastSenderId":        "",
				"lastSenderType":      "",
				"lastSenderFirstName": "",
				"lastSenderLastName":  "",
				"lastContentPreview":  "",
			})

		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		dbNode := result.(*dbtype.Node)
		return mapNodeToConversation(dbNode), nil
	}
}

func (s *customerOSService) UpdateConversation(tenant string, conversationId string, participantId string, participantType entity.SenderType, lastSenderFirstName string, lastSenderLastName string, lastContentPreview string) (string, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	contactIds := []string{}
	userIds := []string{}
	participantTypeString := ""

	if participantType == entity.CONTACT {
		contactIds = append(contactIds, participantId)
		participantTypeString = "CONTACT"
	} else if participantType == entity.USER {
		userIds = append(userIds, participantId)
		participantTypeString = "USER"
	}

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		query := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (o:Conversation {id:$conversationId}) " +
			" ON MATCH SET " +
			" o.messageCount=o.messageCount+1, o.updatedAt=$updatedAt, " +
			" o.lastSenderId=$lastSenderId, o.lastSenderType=$lastSenderType, o.lastSenderFirstName=$lastSenderFirstName, o.lastSenderLastName=$lastSenderLastName, o.lastContentPreview=$lastContentPreview, " +
			" o:%s " +
			" %s %s " +
			" RETURN DISTINCT o"
		queryLinkWithContacts := ""
		if len(contactIds) > 0 {
			queryLinkWithContacts = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) WHERE c.id in $contactIds " +
				" MERGE (c)-[:PARTICIPATES]->(o) "
		}
		queryLinkWithUsers := ""
		if len(userIds) > 0 {
			queryLinkWithUsers = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (u:User)-[:USER_BELONGS_TO_TENANT]->(t) WHERE u.id in $userIds " +
				" MERGE (u)-[:PARTICIPATES]->(o) "
		}
		queryResult, err := tx.Run(fmt.Sprintf(query, "Conversation_"+tenant, queryLinkWithContacts, queryLinkWithUsers),
			map[string]interface{}{
				"tenant":              tenant,
				"conversationId":      conversationId,
				"updatedAt":           time.Now().UTC(),
				"contactIds":          contactIds,
				"userIds":             userIds,
				"lastSenderId":        participantId,
				"lastSenderType":      participantTypeString,
				"lastSenderFirstName": lastSenderFirstName,
				"lastSenderLastName":  lastSenderLastName,
				"lastContentPreview":  lastContentPreview,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return "", err
	} else {
		dbNode := result.(*dbtype.Node)
		return utils.GetPropsFromNode(*dbNode)["id"].(string), err
	}
}

func mapNodeToConversation(node *dbtype.Node) *Conversation {
	if node == nil {
		return nil
	}

	props := utils.GetPropsFromNode(*node)

	conversation := new(Conversation)
	conversation.Id = utils.GetPropsFromNode(*node)["id"].(string)
	conversation.Status = utils.GetPropsFromNode(*node)["status"].(string)
	conversation.Channel = utils.GetPropsFromNode(*node)["channel"].(string)
	conversation.StartedAt = utils.GetPropsFromNode(*node)["startedAt"].(time.Time)
	conversation.UpdatedAt = utils.GetPropsFromNode(*node)["updatedAt"].(time.Time)
	conversation.InitiatorFirstName = utils.GetPropsFromNode(*node)["initiatorFirstName"].(string)
	conversation.InitiatorLastName = utils.GetPropsFromNode(*node)["initiatorLastName"].(string)
	conversation.InitiatorUsername = utils.GetPropsFromNode(*node)["initiatorUsername"].(string)
	conversation.InitiatorType = utils.GetPropsFromNode(*node)["initiatorType"].(string)

	conversation.LastSenderId = utils.GetStringPropOrEmpty(props, "lastSenderId")
	conversation.LastSenderType = utils.GetStringPropOrEmpty(props, "lastSenderType")
	conversation.LastSenderFirstName = utils.GetStringPropOrEmpty(props, "lastSenderFirstName")
	conversation.LastSenderLastName = utils.GetStringPropOrEmpty(props, "lastSenderLastName")
	conversation.LastContentPreview = utils.GetStringPropOrEmpty(props, "lastContentPreview")

	return conversation
}

func mapNodeToContact(id string, firstName string, lastName string, email string) *Contact {
	user := new(Contact)
	user.Id = id
	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email

	return user
}

func mapNodeToUser(node *dbtype.Node) *User {
	if node == nil {
		return nil
	}

	props := utils.GetPropsFromNode(*node)

	user := new(User)
	user.Id = utils.GetPropsFromNode(*node)["id"].(string)
	user.FirstName = utils.GetStringPropOrEmpty(props, "firstName")
	user.LastName = utils.GetStringPropOrEmpty(props, "lastName")
	user.Email = utils.GetStringPropOrEmpty(props, "email")

	return user
}

func (s *customerOSService) GetContactWithEmailOrCreate(tenant string, email string) (Contact, error) {
	contact, err := s.GetContactByEmail(email)
	if err != nil {
		contact, err = s.CreateContactWithEmail(tenant, email)
		if err != nil {
			return Contact{}, err
		}
		if contact == nil {
			return Contact{}, errors.New("contact not found and could not be created")
		}
		return *contact, nil
	} else {
		return *contact, nil
	}
}

func (s *customerOSService) GetActiveConversationOrCreate(
	tenant string,
	participantId string,
	firstName string,
	lastname string,
	username string,
	senderType msProto.SenderType,
	eventType entity.EventType,
) (*Conversation, error) {
	var conversation *Conversation
	var err error

	if senderType == msProto.SenderType_CONTACT {
		conversation, err = s.GetWebChatConversationWithContactInitiator(tenant, participantId)
	} else if senderType == msProto.SenderType_USER {
		conversation, err = s.GetWebChatConversationWithUserInitiator(tenant, participantId)
	}

	if err != nil {
		return nil, err
	}

	if conversation == nil {
		conversation, err = s.CreateConversation(tenant, participantId, firstName, lastname, username, s.commonStoreService.ConvertMSSenderTypeToEntitySenderType(senderType), eventType)
	}
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func NewCustomerOSService(driver *neo4j.Driver, postgresRepositories *repository.PostgresRepositories, commonStoreService *commonStoreService) *customerOSService {
	customerOsService := new(customerOSService)
	customerOsService.driver = driver
	customerOsService.postgresRepositories = postgresRepositories
	customerOsService.commonStoreService = commonStoreService
	return customerOsService
}
