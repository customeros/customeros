package service

import (
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository/entity"
	"regexp"
	"strings"
	"time"
)

type customerOSService struct {
	driver               *neo4j.Driver
	postgresRepositories *repository.PostgresRepositories
}

type CustomerOSService interface {
	ContactByIdExists(contactId string) (bool, error)
	ContactByPhoneExists(e164 string) (bool, error)

	GetUserByEmail(email string) (*UserInfo, error)
	GetContactByEmail(email string) (*ContactInfo, error)

	CreateContactWithEmail(tenant string, email string) (string, error)
	CreateContactWithPhone(phone string) (string, error)

	ConversationByIdExists(tenant string, conversationId string) (bool, error)

	GetConversations(tenant string) ([]Conversation, error)
	GetConversationById(tenant string, conversationId string) (Conversation, error)

	GetWebChatConversationWithContactInitiator(tenant string, contactId string) (string, error)
	GetWebChatConversationWithUserInitiator(tenant string, userId string) (string, error)

	CreateConversation(tenant string, initiatorId string, initiatorType entity.SenderType, channel entity.EventType) (string, error)
	UpdateConversation(tenant string, conversationId string, participantId string, participantType entity.SenderType) (string, error)
}

type UserInfo struct {
	Id string
}

type ContactInfo struct {
	Id string
}

type Conversation struct {
	Id                string
	StartedAt         time.Time
	Channel           string
	Status            string
	InitiatorUsername string
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

func (s *customerOSService) GetUserByEmail(email string) (*UserInfo, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`MATCH (u:User{email:$email}) return u`,
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
	props := utils.GetPropsFromNode(dbNode)

	userInfo := UserInfo{Id: utils.GetStringPropOrEmpty(props, "id")}

	return &userInfo, nil
}

func (s *customerOSService) GetContactByEmail(email string) (*ContactInfo, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c:Contact)-[:EMAILED_AT]->(p:Email {email:$email})
            RETURN c`,
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
	dbNode := (dbRecords.([]*db.Record))[0].Values[0].(dbtype.Node)
	props := utils.GetPropsFromNode(dbNode)

	contactInfo := ContactInfo{Id: utils.GetStringPropOrEmpty(props, "id")}

	return &contactInfo, nil
}

func (s *customerOSService) CreateContactWithEmail(tenant string, email string) (string, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		//create the contact
		contactQuery := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (c:Contact {id:randomUUID()})-[:CONTACT_BELONGS_TO_TENANT]->(t) ON CREATE SET" +
			" c.createdAt=$createdAt, " +
			" c.source=$source, " +
			" c.sourceOfTruth=$sourceOfTruth, " +
			" c:%s " +
			" RETURN c"

		contactQueryResult, err := tx.Run(fmt.Sprintf(contactQuery, "Contact_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"createdAt":     time.Now().UTC(),
				"source":        "openline",
				"sourceOfTruth": "openline",
			})

		node, err := utils.ExtractSingleRecordFirstValueAsNode(contactQueryResult, err)
		if err != nil {
			return nil, err
		}

		//create the email
		emailQuery := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
			" MERGE (c)-[r:EMAILED_AT]->(e:Email {email: $email}) " +
			" ON CREATE SET e.label=$label, r.primary=$primary, e.id=randomUUID(), e:%s " +
			" ON MATCH SET e.label=$label, r.primary=$primary " +
			" RETURN e, r"

		contactId := utils.GetPropsFromNode(*node)["id"].(string)
		_, err = tx.Run(fmt.Sprintf(emailQuery, "Email_"+tenant),
			map[string]interface{}{
				"tenant":    tenant,
				"contactId": contactId,
				"email":     email,
				"label":     "WORK",
				"primary":   true,
			})

		if err != nil {
			return nil, err
		}

		return contactId, nil
	})
	if err != nil {
		return "", err
	} else {
		return result.(string), nil
	}
}

func (s *customerOSService) CreateContactWithPhone(phone string) (string, error) {
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
	return "", nil
}

func (s *customerOSService) ConversationByIdExists(tenant string, conversationId string) (bool, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`MATCH (c:Conversation {id:$conversationId}) RETURN c`, //TODO need to filter by tenant
			map[string]any{
				"tenant":         tenant,
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

		if queryResult, err := tx.Run("MATCH (c:Conversation_"+tenant+") RETURN c", map[string]any{
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
		node := utils.NodePtr(v.Values[0].(neo4j.Node))
		conversation := new(Conversation)
		conversation.Id = utils.GetPropsFromNode(*node)["id"].(string)
		conversation.Status = utils.GetPropsFromNode(*node)["status"].(string)
		conversation.Channel = utils.GetPropsFromNode(*node)["channel"].(string)
		conversation.StartedAt = utils.GetPropsFromNode(*node)["startedAt"].(time.Time)
		conversations = append(conversations, *conversation)
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
			"tenant":         tenant,
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

// returns the conversation if exists otherwise nil
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

// returns the conversation if exists otherwise nil
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

func (s *customerOSService) CreateConversation(tenant string, initiatorId string, initiatorUsername string, initiatorType entity.SenderType, channel entity.EventType) (*Conversation, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	contactIds := []string{}
	userIds := []string{}

	if initiatorType == entity.CONTACT {
		contactIds = append(contactIds, initiatorId)
	} else if initiatorType == entity.USER {
		userIds = append(userIds, initiatorId)
	}

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		query := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (o:Conversation {id:randomUUID()}) " +
			" ON CREATE SET o.initiatorUsername=$initiatorUsername, o.startedAt=$startedAt, o.messageCount=0, o.channel=$channel, o.status=$status, o.source=$source, o.sourceOfTruth=$sourceOfTruth, o.appSource=$appSource, o:%s " +
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
		queryResult, err := tx.Run(fmt.Sprintf(query, "Conversation_"+tenant, queryLinkWithContacts, queryLinkWithUsers),
			map[string]interface{}{
				"tenant":            tenant,
				"source":            "openline",
				"sourceOfTruth":     "openline",
				"appSource":         "manual",
				"status":            "ACTIVE",
				"initiatorUsername": initiatorUsername,
				"startedAt":         time.Now().UTC(),
				"channel":           channel,
				"contactIds":        contactIds,
				"userIds":           userIds,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		dbNode := result.(*dbtype.Node)
		return mapNodeToConversation(dbNode), nil
	}
}

func (s *customerOSService) UpdateConversation(tenant string, conversationId string, participantId string, participantType entity.SenderType) (string, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	contactIds := []string{}
	userIds := []string{}

	if participantType == entity.CONTACT {
		contactIds = append(contactIds, participantId)
	} else if participantType == entity.USER {
		userIds = append(userIds, participantId)
	}

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		query := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (o:Conversation {id:$conversationId}) " +
			" ON MATCH SET o.messageCount=o.messageCount+1, o:%s " +
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
				"tenant":         tenant,
				"conversationId": conversationId,
				"contactIds":     contactIds,
				"userIds":        userIds,
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

	conversation := new(Conversation)
	conversation.Id = utils.GetPropsFromNode(*node)["id"].(string)
	conversation.Status = utils.GetPropsFromNode(*node)["status"].(string)
	conversation.Channel = utils.GetPropsFromNode(*node)["channel"].(string)
	conversation.StartedAt = utils.GetPropsFromNode(*node)["startedAt"].(time.Time)
	conversation.InitiatorUsername = utils.GetPropsFromNode(*node)["initiatorUsername"].(string)

	return conversation
}

func NewCustomerOSService(driver *neo4j.Driver, postgresRepositories *repository.PostgresRepositories) *customerOSService {
	customerOsService := new(customerOSService)
	customerOsService.driver = driver
	customerOsService.postgresRepositories = postgresRepositories
	return customerOsService
}
