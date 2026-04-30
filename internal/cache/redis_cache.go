package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisCache) Get(code string) (string, error) {
	val, err := r.client.Get(r.ctx, code).Result()
	if err == redis.Nil {
		return "", nil // cache miss
	}
	return val, err
}

func (r *RedisCache) Set(code string, url string) error {
	return r.client.Set(r.ctx, code, url, 0).Err()
}

func (r *RedisCache) Delete(code string) error {
	return r.client.Del(r.ctx, code).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
