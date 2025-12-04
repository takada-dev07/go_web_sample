package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go_web_sample/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	// Given: テスト用のEchoインスタンスとリクエストを作成
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set("X-Request-ID", "test-request-123")
	req.Header.Set("X-User-ID", "test-user-456")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Given: ハンドラー関数を作成（タイムアウトと引数を指定）
	timeout := 5 * time.Second
	handlerFunc := handler.HelloWorld(timeout, "TestUser", 42)

	// When: ハンドラーを実行
	err := handlerFunc(c)

	// Then: エラーが発生しないことを確認
	assert.NoError(t, err)

	// Then: HTTPステータスコードが200であることを確認
	assert.Equal(t, http.StatusOK, rec.Code)

	// Then: レスポンスボディが期待値と一致することを確認
	expectedBody := `{"message":"HelloWorld"}`
	assert.JSONEq(t, expectedBody, rec.Body.String())
}

func TestHelloWorld_WithDifferentContext(t *testing.T) {
	// Given: 異なるヘッダー値でテスト
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set("X-Request-ID", "another-request-789")
	req.Header.Set("X-User-ID", "another-user-012")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Given: 異なる引数でハンドラーを作成
	timeout := 5 * time.Second
	handlerFunc := handler.HelloWorld(timeout, "AnotherUser", 100)

	// When: ハンドラーを実行
	err := handlerFunc(c)

	// Then: エラーが発生しないことを確認
	assert.NoError(t, err)

	// Then: HTTPステータスコードが200であることを確認
	assert.Equal(t, http.StatusOK, rec.Code)

	// Then: レスポンスボディが期待値と一致することを確認
	assert.JSONEq(t, `{"message":"HelloWorld"}`, rec.Body.String())
}

func TestHelloWorld_Timeout(t *testing.T) {
	// Given: 非常に短いタイムアウトでテスト
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Given: 非常に短いタイムアウト（1ナノ秒 - 実質的に即座にタイムアウト）
	timeout := 1 * time.Nanosecond
	handlerFunc := handler.HelloWorld(timeout, "TestUser", 42)

	// When: ハンドラーを実行（少し待機してタイムアウトを確実に発生させる）
	time.Sleep(10 * time.Millisecond)
	err := handlerFunc(c)

	// Then: エラーが発生しないことを確認（タイムアウトは正常な動作）
	assert.NoError(t, err)

	// Then: HTTPステータスコードが408（Request Timeout）であることを確認
	assert.Equal(t, http.StatusRequestTimeout, rec.Code)
}
