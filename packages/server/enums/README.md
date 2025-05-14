# Customer OS Enums

This module contains shared enum types used across the Customer OS platform. It provides a central location for enum definitions to ensure consistency across different services.

## Usage

Import the enums module in your code:

```go
import "github.com/customeros/enums"
```

Then use the enum types directly:

```go
eventType := enums.EventTenantCreated
```

## Available Enum Types

### NatsEventType

`NatsEventType` contains the event types used for NATS messaging across the platform. These events are used for asynchronous communication between services.

Example:

```go
import (
    "github.com/customeros/enums"
    "github.com/nats-io/nats.go"
)

msg := nats.NewMsg(enums.EventTenantCreated.String())
```

## Contribution Guidelines

When adding new enum types:

1. Group related constants together
2. Add clear comments for each type and constant
3. Follow Go naming conventions (camelCase for internal, PascalCase for exported)
4. Implement necessary methods (e.g., `String()`, `MarshalJSON()`, etc.) 