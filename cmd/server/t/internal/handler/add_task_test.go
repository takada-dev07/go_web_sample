package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_web_sample/internal/config"
	"go_web_sample/internal/entity"
	"go_web_sample/internal/handler"
	"go_web_sample/pkg/redis"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) {
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

func teardownTestRedis(t *testing.T) {
	// テスト後にRedis接続を閉じる
	redis.Close()
}

func TestAddTask_Success_WithStatusTodo(t *testing.T) {
	// Given: Redis接続を初期化
	setupTestRedis(t)
	defer teardownTestRedis(t)

	// Given: テスト用のEchoインスタンスとリクエストを作成
	e := echo.New()
	reqBody := map[string]string{
		"title":  "テストタスク",
		"status": entity.StatusTodo,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// When: ハンドラーを実行
	err := handler.AddTask(c)

	// Then: エラーが発生しないことを確認
	assert.NoError(t, err)

	// Then: HTTPステータスコードが200であることを確認
	assert.Equal(t, http.StatusOK, rec.Code)

	// Then: レスポンスボディにIDが含まれることを確認
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response["id"])
}

func TestAddTask_Success_WithStatusDoing(t *testing.T) {
	// Given: Redis接続を初期化
	setupTestRedis(t)
	defer teardownTestRedis(t)

	// Given: テスト用のEchoインスタンスとリクエストを作成
	e := echo.New()
	reqBody := map[string]string{
		"title":  "進行中のタスク",
		"status": entity.StatusDoing,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// When: ハンドラーを実行
	err := handler.AddTask(c)

	// Then: エラーが発生しないことを確認
	assert.NoError(t, err)

	// Then: HTTPステータスコードが200であることを確認
	assert.Equal(t, http.StatusOK, rec.Code)

	// Then: レスポンスボディにIDが含まれることを確認
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response["id"])
}

func TestAddTask_Success_WithStatusDone(t *testing.T) {
	// Given: Redis接続を初期化
	setupTestRedis(t)
	defer teardownTestRedis(t)

	// Given: テスト用のEchoインスタンスとリクエストを作成
	e := echo.New()
	reqBody := map[string]string{
		"title":  "完了したタスク",
		"status": entity.StatusDone,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// When: ハンドラーを実行
	err := handler.AddTask(c)

	// Then: エラーが発生しないことを確認
	assert.NoError(t, err)

	// Then: HTTPステータスコードが200であることを確認
	assert.Equal(t, http.StatusOK, rec.Code)

	// Then: レスポンスボディにIDが含まれることを確認
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response["id"])
}

func TestAddTask_Success_WithoutStatus(t *testing.T) {
	// Given: Redis接続を初期化
	setupTestRedis(t)
	defer teardownTestRedis(t)

	// Given: テスト用のEchoインスタンスとリクエストを作成（status省略）
	e := echo.New()
	reqBody := map[string]string{
		"title": "デフォルトステータスのタスク",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// When: ハンドラーを実行
	err := handler.AddTask(c)

	// Then: エラーが発生しないことを確認
	assert.NoError(t, err)

	// Then: HTTPステータスコードが200であることを確認
	assert.Equal(t, http.StatusOK, rec.Code)

	// Then: レスポンスボディにIDが含まれることを確認
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response["id"])
}

func TestAddTask_InvalidJSON(t *testing.T) {
	// Given: テスト用のEchoインスタンスと不正なJSONリクエストを作成
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString("{invalid json}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// When: ハンドラーを実行
	err := handler.AddTask(c)

	// Then: エラーが発生しないことを確認（ハンドラー内でエラーレスポンスを返す）
	assert.NoError(t, err)

	// Then: HTTPステータスコードが400であることを確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Then: エラーメッセージが含まれることを確認
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid request body", response["error"])
}

func TestAddTask_EmptyBody(t *testing.T) {
	// Given: テスト用のEchoインスタンスと空のリクエストボディを作成
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// When: ハンドラーを実行
	err := handler.AddTask(c)

	// Then: エラーが発生しないことを確認（ハンドラー内でエラーレスポンスを返す）
	assert.NoError(t, err)

	// Then: HTTPステータスコードが400であることを確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Then: エラーメッセージが含まれることを確認
	// 空のボディは空のJSONオブジェクトとして解釈されるため、title is required エラーが返される
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "title is required", response["error"])
}

func TestAddTask_EmptyTitle(t *testing.T) {
	// Given: Redis接続を初期化
	setupTestRedis(t)
	defer teardownTestRedis(t)

	// Given: テスト用のEchoインスタンスと空のtitleを含むリクエストを作成
	e := echo.New()
	reqBody := map[string]string{
		"title":  "",
		"status": entity.StatusTodo,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// When: ハンドラーを実行
	err := handler.AddTask(c)

	// Then: エラーが発生しないことを確認（ハンドラー内でエラーレスポンスを返す）
	assert.NoError(t, err)

	// Then: HTTPステータスコードが400であることを確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Then: エラーメッセージが含まれることを確認
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "title is required", response["error"])
}

func TestAddTask_WhitespaceOnlyTitle(t *testing.T) {
	// Given: Redis接続を初期化
	setupTestRedis(t)
	defer teardownTestRedis(t)

	// Given: テスト用のEchoインスタンスと空白のみのtitleを含むリクエストを作成
	e := echo.New()
	reqBody := map[string]string{
		"title":  "   ",
		"status": entity.StatusTodo,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// When: ハンドラーを実行
	err := handler.AddTask(c)

	// Then: エラーが発生しないことを確認（ハンドラー内でエラーレスポンスを返す）
	assert.NoError(t, err)

	// Then: HTTPステータスコードが400であることを確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Then: エラーメッセージが含まれることを確認
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "title is required", response["error"])
}

func TestAddTask_InvalidStatus(t *testing.T) {
	// Given: Redis接続を初期化
	setupTestRedis(t)
	defer teardownTestRedis(t)

	// Given: テスト用のEchoインスタンスと無効なstatusを含むリクエストを作成
	e := echo.New()
	reqBody := map[string]string{
		"title":  "テストタスク",
		"status": "invalid_status",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// When: ハンドラーを実行
	err := handler.AddTask(c)

	// Then: エラーが発生しないことを確認（ハンドラー内でエラーレスポンスを返す）
	assert.NoError(t, err)

	// Then: HTTPステータスコードが400であることを確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Then: エラーメッセージが含まれることを確認
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "status must be one of: todo, doing, done", response["error"])
}
