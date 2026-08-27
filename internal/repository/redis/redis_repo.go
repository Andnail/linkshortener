package redisrepo

import (
	"context"
	"fmt"
	"linkshortener/internal/models"
	"linkshortener/internal/repository"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) repository.CacheRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) Set(ctx context.Context, link *models.Link) error {
	if err := r.client.Set(ctx, link.Origin, link.Shorten, 0).Err(); err != nil {
		return fmt.Errorf("Failed to create a redis row")
	}
	return nil
}

func (r *RedisRepository) Delete(ctx context.Context, url string) error {
	if err := r.client.Del(ctx, url).Err(); err != nil {
		return fmt.Errorf("Failed to delete link")
	}
	return nil
}

func (r *RedisRepository) GetByShorten(ctx context.Context, shortLink string) (string, error) {
	return "", nil
}

func (r *RedisRepository) GetByOrigin(ctx context.Context, originUrl string) (string, error) {
	val, err := r.client.Get(ctx, originUrl).Result()
	if err == nil {
		return val, nil
	}
	return "", fmt.Errorf("failed to get shorten url: %s", err)
}
