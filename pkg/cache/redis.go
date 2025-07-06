// pkg/cache/redis.go
package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"zleeper-be/config"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisClient(config config.RedisConfig) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, jsonVal, expiration).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) DeleteAll(ctx context.Context, key string) error {

	iter := r.client.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.Delete(ctx, iter.Val()); err != nil {
			log.Printf("Failed to delete key %s: %v", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}