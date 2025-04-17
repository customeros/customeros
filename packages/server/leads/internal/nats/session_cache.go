package nats_internal

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
)

// SessionCache provides IP lookup caching using NATS KV
type SessionCache struct {
	kv nats.KeyValue
}

const SESSION_CACHE_BUCKET = "session_cache"

// NewIPCache creates a new NATS KV-backed IP cache
func NewSessionCache(js nats.JetStreamContext) (*SessionCache, error) {
	// First try to get the bucket
	kv, err := js.KeyValue(SESSION_CACHE_BUCKET)
	// If bucket doesn't exist, create it
	if err != nil {
		if err == nats.ErrBucketNotFound {
			kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
				Bucket:  SESSION_CACHE_BUCKET,
				History: 1,
				TTL:     24 * time.Hour, // Safety TTL for orphaned entries
			})
			if err != nil {
				return nil, err
			}
		} else {
			// Some other error occurred
			return nil, err
		}
	}

	return &SessionCache{
		kv: kv,
	}, nil
}

func (c *SessionCache) Get(ctx context.Context, visitorID string) (string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "SessionCache.Get")
	defer span.Finish()

	entry, err := c.kv.Get(visitorID)
	if err != nil {
		if err == nats.ErrKeyNotFound {
			return "", nil // Return empty string if not found
		}
		return "", err
	}

	// Assuming the value is a string
	var sessionID string
	if err := json.Unmarshal(entry.Value(), &sessionID); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (c *SessionCache) Set(ctx context.Context, visitorID, sessionID string) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "SessionCache.Set")
	defer span.Finish()

	value, err := json.Marshal(sessionID)
	if err != nil {
		return err
	}

	_, err = c.kv.Put(visitorID, value)
	return err
}

func (c *SessionCache) Delete(visitorID string) error {
	return c.kv.Delete(visitorID)
}
