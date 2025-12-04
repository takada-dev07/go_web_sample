package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config アプリケーション設定
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

// ServerConfig サーバー設定
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig データベース設定
type DatabaseConfig struct {
	Type string
	Path string
}

// JWTConfig JWT設定
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// CORSConfig CORS設定
type CORSConfig struct {
	AllowedOrigins []string
}

// Load 設定を読み込む
func Load() (*Config, error) {
	// .envファイルを読み込む（エラーは無視 - 本番環境では環境変数を使用）
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Type: getEnv("DB_TYPE", "sqlite"),
			Path: getEnv("DB_PATH", "./data/app.db"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:8080"}),
		},
	}

	// バリデーション
	if config.JWT.Secret == "your-secret-key-change-in-production" && config.Server.Env == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production environment")
	}

	return config, nil
}

// getEnv 環境変数を取得（デフォルト値付き）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 環境変数を整数として取得
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsSlice 環境変数をスライスとして取得（カンマ区切り）
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	// カンマ区切りで分割
	var result []string
	start := 0
	for i, char := range valueStr {
		if char == ',' {
			if start < i {
				result = append(result, valueStr[start:i])
			}
			start = i + 1
		}
	}
	if start < len(valueStr) {
		result = append(result, valueStr[start:])
	}

	if len(result) == 0 {
		return defaultValue
	}
	return result
}

