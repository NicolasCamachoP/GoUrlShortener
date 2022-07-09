package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

//DbIface redis implementation
type RedisHandler struct {
	redisClient *redis.Client
	ctx         context.Context
}

type DbOptions struct {
	Domain     string
	PortNumber int
	Password   string
}

func NewRedisHandler(dbOpts *DbOptions) (*RedisHandler, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", dbOpts.Domain, dbOpts.PortNumber),
		Password: dbOpts.Password,
		DB:       0,
	})
	ctx := context.Background()

	pong, err := redisClient.Ping(ctx).Result()

	if err != nil {
		log.Println("[ERROR] [DB] - DB not responding", pong)
		return nil, fmt.Errorf("error while pinging redis: %w", err)
	}
	return &RedisHandler{redisClient, ctx}, nil
}

func (rc *RedisHandler) GetItem(key string) (interface{}, error) {
	value, err := rc.redisClient.Get(rc.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("error in Get operation: %w", err)
	}
	return value, nil
}

func (rc *RedisHandler) SetItem(key string, value interface{}, expiration time.Duration) error {
	err := rc.redisClient.Set(rc.ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("error while saving value: %w", err)
	}
	return nil
}

func (rc *RedisHandler) Exists(key string) bool {
	return rc.redisClient.Exists(rc.ctx, key).Val() > 0
}
