package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// contextキーの型定義
type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
)

// RequestIDKey RequestIDのcontextキーを取得
func RequestIDKey() contextKey {
	return requestIDKey
}

// UserIDKey UserIDのcontextキーを取得
func UserIDKey() contextKey {
	return userIDKey
}

// HelloWorld HelloWorldを返すハンドラー
// timeout: リクエストのタイムアウト時間
// name: ユーザー名
// count: カウント値
func HelloWorld(timeout time.Duration, name string, count int) echo.HandlerFunc {
	return func(c echo.Context) error {
		// リクエストごとにタイムアウト付きのcontextを作成
		ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
		defer cancel()

		// リクエストIDとユーザーIDをcontextに設定
		ctx = context.WithValue(ctx, requestIDKey, c.Request().Header.Get("X-Request-ID"))
		if ctx.Value(requestIDKey) == nil || ctx.Value(requestIDKey) == "" {
			ctx = context.WithValue(ctx, requestIDKey, "default-request")
		}
		ctx = context.WithValue(ctx, userIDKey, c.Request().Header.Get("X-User-ID"))
		if ctx.Value(userIDKey) == nil || ctx.Value(userIDKey) == "" {
			ctx = context.WithValue(ctx, userIDKey, "default-user")
		}

		// キャンセルまたはタイムアウトを監視するチャネル
		done := make(chan error, 1)

		// 実際の処理をゴルーチンで実行
		go func() {
			// contextのキャンセルをチェック
			select {
			case <-ctx.Done():
				// キャンセルされた場合はエラーを返す
				done <- ctx.Err()
				return
			default:
			}

			// context.Contextの値を標準出力に出力
			fmt.Printf("Context: %v\n", ctx)
			fmt.Printf("Request ID: %v\n", ctx.Value(requestIDKey))
			fmt.Printf("User ID: %v\n", ctx.Value(userIDKey))

			// 引数の値を標準出力に出力
			fmt.Printf("Name: %s\n", name)
			fmt.Printf("Count: %d\n", count)

			// 処理が完了したことを通知
			done <- nil
		}()

		// タイムアウトまたはキャンセルを待機
		select {
		case <-ctx.Done():
			// タイムアウトまたはキャンセルが発生
			fmt.Printf("Request cancelled or timed out: %v\n", ctx.Err())
			if ctx.Err() == context.DeadlineExceeded {
				return c.JSON(http.StatusRequestTimeout, map[string]string{
					"error": "Request timeout",
				})
			}
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": "Request cancelled",
			})
		case err := <-done:
			// 処理が正常に完了
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, map[string]string{
				"message": "HelloWorld3",
			})
		}
	}
}
