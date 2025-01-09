package entity

type AgentRegistry struct {
	ID                string  `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Name              string  `gorm:"column:name;type:varchar(255);not null;index" json:"name" binding:"required"`
	Description       string  `gorm:"column:description;type:varchar(255)" json:"description"`
	WorkDefinition    string  `gorm:"column:work_definition;type:text" json:"workDefinition"`
	UserContextSchema *string `gorm:"column:user_context_schema;type:text" json:"userContext"`
	OutputSchema      *string `gorm:"column:output_schema;type:text" json:"output"`
	ConfigSchema      *string `gorm:"column:config_schema;type:text" json:"config"`
	IsActive          string  `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
}

func (AgentRegistry) TableName() string {
	return "agent_registry"
}

// Example instance (the string fields would contain JSON-marshaled data)
// agent := AgentRegistry{
//     Name:        "Traffic Deanonymizer",
//     Description: "Enriches anonymous website traffic data with user information and optionally sends notifications to Slack",
//     WorkDefinition: `{
//         "capabilities": [
//             {
//                 "name": "deanonymize_traffic",
//                 "description": "Enriches IP addresses with user/company information",
//                 "required_permissions": ["ip_lookup", "company_data_access"]
//             },
//             {
//                 "name": "notify_slack",
//                 "description": "Sends notifications to Slack",
//                 "required_permissions": ["slack_webhook"],
//                 "optional": true
//             }
//         ],
//         "supported_content_types": ["application/json", "text/csv"],
//         "rate_limits": {
//             "requests_per_minute": 30
//         }
//     }`,
//     InputSchema: `{
//         "type": "object",
//         "properties": {
//             "traffic_data": {
//                 "type": "array",
//                 "items": {
//                     "type": "object",
//                     "properties": {
//                         "ip_address": { "type": "string" },
//                         "timestamp": { "type": "string", "format": "date-time" }
//                     },
//                     "required": ["ip_address"]
//                 }
//             }
//         }
//     }`,
//     OutputSchema: `{
//         "type": "object",
//         "properties": {
//             "enriched_traffic": {
//                 "type": "array",
//                 "items": {
//                     "type": "object",
//                     "properties": {
//                         "ip_address": { "type": "string" },
//                         "company_name": { "type": "string" },
//                         "company_size": { "type": "string" },
//                         "notification_sent": { "type": "boolean" }
//                     }
//                 }
//             }
//         }
//     }`,
//     ConfigSchema: `{
//         "type": "object",
//         "properties": {
//             "slack_enabled": {
//                 "type": "boolean",
//                 "default": false
//             },
//             "slack_channel_id": {
//                 "type": "string",
//                 "pattern": "^C[A-Z0-9]{8,}$"
//             }
//         },
//         "dependencies": {
//             "slack_enabled": {
//                 "if": { "properties": { "slack_enabled": { "const": true } } },
//                 "then": { "required": ["slack_channel_id"] }
//             }
//         }
//     }`,
//     Status: "active",
// }
