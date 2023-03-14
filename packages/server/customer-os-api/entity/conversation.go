package entity

import (
	"fmt"
	"time"
)

type ConversationEntity struct {
	Id                 string
	StartedAt          time.Time  `neo4jDb:"property:startedAt;lookupName:STARTED_AT;supportCaseSensitive:false"`
	UpdatedAt          time.Time  `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	EndedAt            *time.Time `neo4jDb:"property:endedAt;lookupName:ENDED_AT;supportCaseSensitive:false"`
	Status             string     `neo4jDb:"property:status;lookupName:STATUS;supportCaseSensitive:true"`
	Channel            string     `neo4jDb:"property:channel;lookupName:CHANNEL;supportCaseSensitive:true"`
	Subject            string     `neo4jDb:"property:subject;lookupName:SUBJECT;supportCaseSensitive:true"`
	MessageCount       int64      `neo4jDb:"property:messageCount;lookupName:MESSAGE_COUNT;supportCaseSensitive:false"`
	Source             DataSource `neo4jDb:"property:source;lookupName:SOURCE;supportCaseSensitive:false"`
	SourceOfTruth      DataSource `neo4jDb:"property:sourceOfTruth;lookupName:SOURCE_OF_TRUTH;supportCaseSensitive:false"`
	AppSource          string     `neo4jDb:"property:appSource;lookupName:APP_SOURCE;supportCaseSensitive:true"`
	InitiatorFirstName string     `neo4jDb:"property:initiatorFirstName;lookupName:INITIATOR_FIRST_NAME;supportCaseSensitive:true"`
	InitiatorLastName  string     `neo4jDb:"property:initiatorLastName;lookupName:INITIATOR_LAST_NAME;supportCaseSensitive:true"`
	InitiatorType      string     `neo4jDb:"property:initiatorType;lookupName:INITIATOR_TYPE;supportCaseSensitive:false"`
	InitiatorUsername  string     `neo4jDb:"property:initiatorUsername;lookupName:INITIATOR_USERNAME;supportCaseSensitive:true"`
	ThreadId           string     `neo4jDb:"property:threadId;lookupName:THREAD_ID;supportCaseSensitive:true"`
}

func (conversation ConversationEntity) ToString() string {
	return fmt.Sprintf("id: %s", conversation.Id)
}

type ConversationEntities []ConversationEntity

func (ConversationEntity) Action() {
}

func (ConversationEntity) ActionName() string {
	return NodeLabel_Conversation
}

func (ConversationEntity) TimelineEvent() {
}

func (ConversationEntity) TimelineEventName() string {
	return NodeLabel_Conversation
}

func (conversation ConversationEntity) Labels(tenant string) []string {
	return []string{"Conversation", "Action", "Conversation_" + tenant, "Action_" + tenant}
}
