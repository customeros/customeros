package service

import (
	"context"
	"fmt"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/redis/go-redis/v9"
	"log"
)

type redisService struct {
	client *redis.Client
	conf   *c.Config
}

type RedisService interface {
	GetKeyInfo(ctx context.Context, tag, key string) (bool, *string)
}

func (s *redisService) GetKeyInfo(ctx context.Context, tag, key string) (bool, *string) {
	redisErr := s.client.Ping(ctx).Err()
	if redisErr != nil {
		log.Printf("Redis ping error: %v", redisErr)
		return false, nil
	}

	data, err := s.client.HGetAll(ctx, fmt.Sprintf("%s:%s", tag, key)).Result()
	if err != nil {
		log.Printf("Redis HGetAll error: %v", err)
		return false, nil
	}

	tenant := data["tenant"]
	return data["active"] == "true", &tenant
}

func NewRedisService(redisClient *redis.Client, config *c.Config) RedisService {
	return &redisService{
		client: redisClient,
		conf:   config,
	}
}
