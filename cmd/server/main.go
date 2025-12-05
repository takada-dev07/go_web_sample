package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_web_sample/internal/config"
	"go_web_sample/internal/handler"
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
	setupRoutes(e, cfg)

	// サーバー起動
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("Server starting on %s (env: %s)\n", serverAddr, cfg.Server.Env)

	// シグナル受信を検知するためのcontextを作成
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// グレースフルシャットダウン
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// シグナル受信を待機
	<-ctx.Done()

	fmt.Println("Shutting down server...")

	// 環境変数からタイムアウト時間を取得
	shutdownTimeout := time.Duration(cfg.Server.ShutdownTimeout) * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("Server stopped")
}

// setupRoutes ルーティングを設定
func setupRoutes(e *echo.Echo, cfg *config.Config) {
	// ヘルスチェックエンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// HelloWorldエンドポイント
	// リクエストごとのタイムアウトを環境変数から取得
	requestTimeout := time.Duration(cfg.Server.RequestTimeout) * time.Second
	e.GET("/hello", handler.HelloWorld(requestTimeout, "DefaultUser", 1))

	// API v1
	api := e.Group("/api/v1")
	{
		api.GET("/", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"message": "Go Web API v1",
			})
		})
		api.POST("/tasks", handler.AddTask)
	}
}
