package entity

//
// import "time"
//
// // The main flow definition
// type Flow struct {
//    ID            string    `json:"id"`
//    TenantID      string    `json:"tenantId"`
//    Name          string    `json:"name"`
//    Description   string    `json:"description"`
//    Status        FlowStatus `json:"status"`
//    ListenerEvent string    `json:"listenerEvent"` // Event that triggers this flow
//    CreatedAt     time.Time `json:"createdAt"`
//    UpdatedAt     time.Time `json:"updatedAt"`
//    CreatedBy     string    `json:"createdBy"`
//    LastModifiedBy string   `json:"lastModifiedBy"`
//
//    // Statistics
//    TotalExecutions      int64 `json:"totalExecutions"`
//    SuccessfulExecutions int64 `json:"successfulExecutions"`
//    FailedExecutions     int64 `json:"failedExecutions"`
// }
//
// // Node in the flow
// type FlowNode struct {
//    ID          string    `json:"id"`
//    FlowID      string    `json:"flowId"`
//    Type        NodeType  `json:"type"`    // action, condition
//    Config      JSON      `json:"config"`   // Node-specific configuration
//    Position    Position  `json:"position"` // UI positioning
// }
//
// // Connection between nodes
// type FlowEdge struct {
//    ID          string    `json:"id"`
//    FlowID      string    `json:"flowId"`
//    Source      string    `json:"source"`      // Source node ID
//    Target      string    `json:"target"`      // Target node ID
//    Condition   string    `json:"condition"`   // For branching paths
// }
//
// type Position struct {
//    X int `json:"x"`
//    Y int `json:"y"`
// }
//
// type FlowStatus string
//
// const (
//    FlowStatusDraft    FlowStatus = "DRAFT"
//    FlowStatusActive   FlowStatus = "ACTIVE"
//    FlowStatusArchived FlowStatus = "ARCHIVED"
// )
//
// type NodeType string
//
// const (
//    NodeTypeAction    NodeType = "ACTION"
//    NodeTypeCondition NodeType = "CONDITION"
//    NodeTypeTrigger   NodeType = "TRIGGER"
// )
//
// // Optional configuration types that could be stored in Node.Config
// type ActionConfig struct {
//    ActionType  string       `json:"actionType"`
//    Schedule    *Schedule    `json:"schedule,omitempty"`
//    RateLimit   *RateLimit   `json:"rateLimit,omitempty"`
//    // Action specific settings
//    Settings    JSON         `json:"settings"`
// }
//
// type ConditionConfig struct {
//    Type        string   `json:"type"`         // equals, contains, greater_than, etc
//    Field       string   `json:"field"`        // field to evaluate
//    Value       any      `json:"value"`        // value to compare against
//    Branches    []string `json:"branches"`     // possible outcome paths
// }
//
// type Schedule struct {
//    Delay      *time.Duration `json:"delay,omitempty"`      // Delay after previous step
//    ExecuteAt  *string        `json:"executeAt,omitempty"`  // Cron expression or specific time
// }
//
// type RateLimit struct {
//    Type      string `json:"type"`      // daily, hourly, etc.
//    Limit     int    `json:"limit"`     // Maximum number of executions
//    EntityID  string `json:"entityId"`  // ID of entity being rate limited
// }
