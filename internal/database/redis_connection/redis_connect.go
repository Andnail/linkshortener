package redisconnection

import (
	"context"
	"linkshortener/internal/config"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*redis.Client, error) {
	logger.Info("Creating redis client")

	client := redis.NewClient(&redis.Options{
		Addr:            cfg.Redis.Addr,
		Username:        cfg.Redis.User,
		Password:        cfg.Redis.Password,
		DB:              cfg.Redis.DB,
		Protocol:        cfg.Redis.Protocol,
		MaxRetries:      cfg.Redis.MaxRetries,
		DialTimeout:     cfg.Redis.DialTimeout,
		ConnMaxLifetime: cfg.Redis.ConnMaxLifetimeout,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		return nil, err
	}

	return client, nil

}
