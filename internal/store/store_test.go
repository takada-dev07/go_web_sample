package store

import (
	"testing"
	"time"

	"go_web_sample/internal/config"
	"go_web_sample/internal/entity"
	"go_web_sample/pkg/redis"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedisForStore(t *testing.T) {
	// Given: テスト用のRedis設定
	// 127.0.0.1を明示的に指定してIPv4接続を確実にする
	cfg := &config.RedisConfig{
		Host:     "127.0.0.1",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	// Redis接続を初期化
	err := redis.Init(cfg)
	if err != nil {
		t.Skipf("Redis is not available: %v", err)
	}
}

func teardownTestRedisForStore(t *testing.T) {
	// テスト後にRedis接続を閉じる
	redis.Close()
}

func TestGenerateTaskID_Success(t *testing.T) {
	// Given: Redis接続を初期化
	setupTestRedisForStore(t)
	defer teardownTestRedisForStore(t)

	// When: タスクIDを生成
	id, err := GenerateTaskID()

	// Then: エラーが発生しないことを確認
	require.NoError(t, err)

	// Then: IDが空でないことを確認
	assert.NotEmpty(t, id)

	// Then: IDが数値文字列であることを確認（複数回呼び出してインクリメントされることを確認）
	id2, err := GenerateTaskID()
	require.NoError(t, err)
	assert.NotEqual(t, id, id2, "ID should be different on second call")
}

func TestSaveTask_Success(t *testing.T) {
	// Given: Redis接続を初期化
	setupTestRedisForStore(t)
	defer teardownTestRedisForStore(t)

	// Given: タスクエンティティを作成
	task := &entity.Task{
		ID:      "test-123",
		Title:   "テストタスク",
		Status:  entity.StatusTodo,
		Created: time.Now(),
	}

	// When: タスクを保存
	err := SaveTask(task)

	// Then: エラーが発生しないことを確認
	assert.NoError(t, err)

	// Then: Redisから取得して内容を確認
	var savedTask entity.Task
	err = redis.Get("task:test-123", &savedTask)
	require.NoError(t, err)
	assert.Equal(t, task.ID, savedTask.ID)
	assert.Equal(t, task.Title, savedTask.Title)
	assert.Equal(t, task.Status, savedTask.Status)
}

func TestSaveTask_MultipleTasks(t *testing.T) {
	// Given: Redis接続を初期化
	setupTestRedisForStore(t)
	defer teardownTestRedisForStore(t)

	// Given: 複数のタスクエンティティを作成
	task1 := &entity.Task{
		ID:      "test-1",
		Title:   "タスク1",
		Status:  entity.StatusTodo,
		Created: time.Now(),
	}
	task2 := &entity.Task{
		ID:      "test-2",
		Title:   "タスク2",
		Status:  entity.StatusDoing,
		Created: time.Now(),
	}

	// When: 複数のタスクを保存
	err1 := SaveTask(task1)
	err2 := SaveTask(task2)

	// Then: エラーが発生しないことを確認
	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// Then: それぞれのタスクが正しく保存されていることを確認
	var savedTask1 entity.Task
	err := redis.Get("task:test-1", &savedTask1)
	require.NoError(t, err)
	assert.Equal(t, task1.Title, savedTask1.Title)

	var savedTask2 entity.Task
	err = redis.Get("task:test-2", &savedTask2)
	require.NoError(t, err)
	assert.Equal(t, task2.Title, savedTask2.Title)
}
