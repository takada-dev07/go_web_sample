package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Set キャッシュに値を設定
func Set(key string, value interface{}, expiration time.Duration) error {
	if client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return client.Set(ctx, key, data, expiration).Err()
}

// Get キャッシュから値を取得
func Get(key string, dest interface{}) error {
	if client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	data, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return fmt.Errorf("failed to get value: %w", err)
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete キャッシュから値を削除
func Delete(key string) error {
	if client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	return client.Del(ctx, key).Err()
}

// Exists キーが存在するか確認
func Exists(key string) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("redis client is not initialized")
	}

	count, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}

	return count > 0, nil
}

// SetNX キーが存在しない場合のみ値を設定（排他制御）
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("redis client is not initialized")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return client.SetNX(ctx, key, data, expiration).Result()
}
