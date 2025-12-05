package store

import (
	"context"
	"fmt"
	"strconv"

	"go_web_sample/internal/entity"
	"go_web_sample/pkg/redis"
)

var ctx = context.Background()

// GenerateTaskID Redis INCRを使用してタスクIDを生成
func GenerateTaskID() (string, error) {
	client := redis.GetClient()
	if client == nil {
		return "", fmt.Errorf("redis client is not initialized")
	}

	// task:counter キーを使用してインクリメント
	id, err := client.Incr(ctx, "task:counter").Result()
	if err != nil {
		return "", fmt.Errorf("failed to generate task ID: %w", err)
	}

	return strconv.FormatInt(id, 10), nil
}

// SaveTask タスクをRedisに保存
func SaveTask(task *entity.Task) error {
	// Redisキー形式: task:{id}
	key := fmt.Sprintf("task:%s", task.ID)

	// TTLなし（無期限）で保存
	// expiration に 0 を指定すると無期限になる
	return redis.Set(key, task, 0)
}
