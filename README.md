# Go Web Sample

Echoフレームワークを使用したGo Webアプリケーションのサンプルプロジェクトです。

## 機能

- REST API
- データベース連携（GORM + SQLite）
- Redisキャッシュ
- JWT認証・認可
- CORS設定
- Dockerコンテナ対応
- AWS ECS Fargateデプロイ対応

## 技術スタック

- **フレームワーク**: Echo v4
- **ORM**: GORM
- **データベース**: SQLite（開発環境）、PostgreSQL/MySQL対応可能
- **キャッシュ**: Redis (go-redis/v9)
- **認証**: JWT (golang-jwt/jwt/v5)
- **設定管理**: godotenv
- **コンテナ**: Docker
- **デプロイ**: AWS ECS Fargate

## プロジェクト構造

```
go_web_sample/
├── cmd/
│   └── server/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── config/              # 設定管理
│   ├── handler/             # HTTPハンドラー
│   ├── model/               # データモデル
│   ├── repository/          # データアクセス層
│   ├── service/             # ビジネスロジック層
│   └── middleware/          # カスタムミドルウェア
├── pkg/                     # 共通パッケージ
│   └── redis/               # Redis接続管理とキャッシュユーティリティ
├── migrations/              # データベースマイグレーション
├── deploy/                  # AWSデプロイ関連
│   ├── task-definition.json # ECSタスク定義
│   ├── Dockerfile.prod      # 本番用Dockerfile
│   ├── deploy.sh            # デプロイスクリプト
│   └── README.md            # デプロイ手順
├── Dockerfile               # Dockerfile
├── docker-compose.yml       # Docker Compose設定
├── .dockerignore            # Docker除外設定
├── .env.example             # 環境変数テンプレート
└── README.md
```

## セットアップ

### 1. 前提条件

- Go 1.23以上
- Git

### 2. 環境変数の設定

`.env.example`をコピーして`.env`ファイルを作成し、必要に応じて値を変更してください。

```bash
cp .env.example .env
```

### 3. 依存関係のインストール

```bash
go mod download
```

### 4. サーバーの起動

#### ローカル実行

```bash
go run cmd/server/main.go
```

サーバーはデフォルトで `http://localhost:8080` で起動します。

#### Docker Composeで実行

```bash
# ビルドと起動
docker-compose up -d

# ログ確認
docker-compose logs -f app

# 停止
docker-compose down
```

アプリケーションとRedisがコンテナで起動します。

## API エンドポイント

### ヘルスチェック

```
GET /health
```

### API v1

```
GET /api/v1/
```

## 開発

### ビルド

```bash
go build -o bin/server cmd/server/main.go
```

### テスト

```bash
go test ./...
```

## 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| PORT | サーバーポート | 8080 |
| ENV | 環境（development/production） | development |
| DB_TYPE | データベースタイプ | sqlite |
| DB_PATH | データベースパス | ./data/app.db |
| REDIS_HOST | Redisホスト | localhost |
| REDIS_PORT | Redisポート | 6379 |
| REDIS_PASSWORD | Redisパスワード | （空） |
| REDIS_DB | Redisデータベース番号 | 0 |
| JWT_SECRET | JWT秘密鍵 | your-secret-key-change-in-production |
| JWT_EXPIRATION_HOURS | JWT有効期限（時間） | 24 |
| CORS_ALLOWED_ORIGINS | CORS許可オリジン（カンマ区切り） | <http://localhost:3000,http://localhost:8080> |

## Docker

### イメージのビルド

```bash
docker build -t go-web-sample .
```

### コンテナの実行

```bash
docker run -p 8080:8080 --env-file .env go-web-sample
```

## AWS ECS Fargate デプロイ

詳細なデプロイ手順は [`deploy/README.md`](deploy/README.md) を参照してください。

### クイックスタート

1. ECRリポジトリとECSクラスターを作成
2. `deploy/task-definition.json`を編集（プレースホルダーを実際の値に置き換え）
3. デプロイスクリプトを実行：

```bash
./deploy/deploy.sh
```

### 必要なAWSリソース

- ECRリポジトリ
- ECSクラスター
- VPCとサブネット
- セキュリティグループ
- IAMロール（タスク実行ロール、タスクロール）
- CloudWatch Logsグループ
- Secrets Manager（シークレット管理用）
- ElastiCache（Redis、オプション）
- RDS（データベース、オプション）

## ライセンス

MIT
