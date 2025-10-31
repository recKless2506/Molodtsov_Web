package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedis(addr, password string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisClient{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Добавляем токен в blacklist
func (r *RedisClient) BlacklistToken(token string, ttl time.Duration) error {
	return r.client.Set(r.ctx, "blacklist:"+token, true, ttl).Err()
}

// Проверяем токен на черный список
func (r *RedisClient) IsTokenBlacklisted(token string) (bool, error) {
	res, err := r.client.Get(r.ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return res == "1" || res == "true", nil
}
