package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go_web_sample/internal/config"
	"go_web_sample/pkg/redis"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 設定を読み込む
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Redis接続を初期化
	if err := redis.Init(&cfg.Redis); err != nil {
		fmt.Printf("Failed to initialize Redis: %v\n", err)
		// Redis接続失敗時もサーバーは起動（オプショナル）
	} else {
		fmt.Println("Redis connected successfully")
		defer redis.Close()
	}

	// Echoインスタンスを作成
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// ルーティング設定
	setupRoutes(e)

	// サーバー起動
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("Server starting on %s (env: %s)\n", serverAddr, cfg.Server.Env)

	// グレースフルシャットダウン
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("Server stopped")
}

// setupRoutes ルーティングを設定
func setupRoutes(e *echo.Echo) {
	// ヘルスチェックエンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// HelloWorldエンドポイント
	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "HelloWorld",
		})
	})

	// API v1
	api := e.Group("/api/v1")
	{
		api.GET("/", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"message": "Go Web API v1",
			})
		})
	}
}
