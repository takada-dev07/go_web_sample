package redis

import (
	"context"
	"fmt"

	"go_web_sample/internal/config"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// Init Redisクライアントを初期化
func Init(cfg *config.RedisConfig) error {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 接続テスト
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// GetClient Redisクライアントを取得
func GetClient() *redis.Client {
	return client
}

// Close Redis接続を閉じる
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
